package video

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/freewilll/apple2/mmu"
)

const (
	// ScreenSizeFactor is the fFactor by which the whole screen is resized
	ScreenSizeFactor = 1

	textVideoMemory = 0x400 // Base location of page 1 text video memory
	flashFrames     = 11    // Number of frames when FLASH mode is toggled
)

// drawTextLoresByte is a function definition used for mixed text/lores rendering
type drawTextLoresByte func(*ebiten.Image, int, int, uint8) error

var (
	charMap      *ebiten.Image     // Character map for text screen
	flashCounter int               // Counter used for flashing characters on the text screen
	flashOn      bool              // Are we currently flashing?
	loresSquares [16]*ebiten.Image // Colored blocks for lores rendering

	// ShowFPS determines if the FPS is shown in the corner of the video
	ShowFPS bool
)

// initTextCharMap initializes the text character map
func initTextCharMap() {
	// The character map pr-latin1.png was downloaded from
	// http://www.kreativekorp.com/software/fonts/apple2.shtml

	var err error
	charMap, _, err = ebitenutil.NewImageFromFile("video/pr-latin1.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}
}

// initLoresSquares creates 16 colored squares for the lores renderer
func initLoresSquares() {
	var err error
	for i := 0; i < 16; i++ {
		loresSquares[i], err = ebiten.NewImage(7, 4, ebiten.FilterNearest)
		if err != nil {
			panic(err)
		}
	}

	// From
	// https://mrob.com/pub/xgithub.com/freewilll/apple2/colors.html
	// https://archive.org/details/IIgs_2523063_Master_Color_Values
	alpha := uint8(0xff)
	loresSquares[0x00].Fill(color.NRGBA{0, 0, 0, alpha})
	loresSquares[0x01].Fill(color.NRGBA{221, 0, 51, alpha})
	loresSquares[0x02].Fill(color.NRGBA{0, 0, 153, alpha})
	loresSquares[0x03].Fill(color.NRGBA{221, 34, 221, alpha})
	loresSquares[0x04].Fill(color.NRGBA{0, 119, 34, alpha})
	loresSquares[0x05].Fill(color.NRGBA{85, 85, 85, alpha})
	loresSquares[0x06].Fill(color.NRGBA{34, 34, 255, alpha})
	loresSquares[0x07].Fill(color.NRGBA{102, 170, 255, alpha})
	loresSquares[0x08].Fill(color.NRGBA{136, 85, 0, alpha})
	loresSquares[0x09].Fill(color.NRGBA{255, 102, 0, alpha})
	loresSquares[0x0A].Fill(color.NRGBA{170, 170, 170, alpha})
	loresSquares[0x0B].Fill(color.NRGBA{255, 153, 136, alpha})
	loresSquares[0x0C].Fill(color.NRGBA{17, 221, 0, alpha})
	loresSquares[0x0D].Fill(color.NRGBA{255, 255, 0, alpha})
	loresSquares[0x0E].Fill(color.NRGBA{68, 255, 153, alpha})
	loresSquares[0x0F].Fill(color.NRGBA{255, 255, 255, alpha})
}

// Init the video data structures used for rendering
func Init() {
	ShowFPS = false

	initTextCharMap()
	initLoresSquares()
}

// drawText draws a single text character at x, y. The characters are either normal, inverted or flashing
func drawText(screen *ebiten.Image, x int, y int, value uint8) error {
	// Determine if the character is inverted
	inverted := false

	if (value & 0xc0) == 0 {
		// Inverted
		inverted = true
	} else if (value & 0x80) == 0 {
		// Flashing
		value = value & 0x3f
		inverted = flashOn
	}

	// Convert the value to a index for the charMap
	if !inverted {
		value = value & 0x7f
	}

	if value < 0x20 {
		value += 0x40
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(ScreenSizeFactor, ScreenSizeFactor)
	op.GeoM.Translate(ScreenSizeFactor*7*float64(x), ScreenSizeFactor*8*float64(y))

	// Grab the image from the font image
	fontRow := value % 16
	fontCol := value / 16
	var fontX = (int)(15 + fontCol*12)
	var fontY = (int)(32 + fontRow*11)
	r := image.Rect(fontX, fontY, fontX+7, fontY+8)
	op.SourceRect = &r

	// The charMap is already inverted. Invert it back if we ourselves aren't inverted.
	if !inverted {
		op.ColorM.Scale(-1, -1, -1, 1)
		op.ColorM.Translate(1, 1, 1, 0)
	}

	// Make it look greenish
	op.ColorM.Scale(0.20, 0.75, 0.20, 1)

	return screen.DrawImage(charMap, op)
}

// drawLores draws two colored lores squares at the equivalent text location x,y.
func drawLores(screen *ebiten.Image, x int, y int, value uint8) error {
	// Convert the 8 bit value to two 4 bit values
	var values = [2]uint8{value & 0xf, value >> 4}

	// Render top & bottom squares
	for i := 0; i < 2; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(ScreenSizeFactor, ScreenSizeFactor)
		op.GeoM.Translate(ScreenSizeFactor*7*float64(x), ScreenSizeFactor*8*float64(y)+float64(i)*4)
		if err := screen.DrawImage(loresSquares[values[i]], op); err != nil {
			return err
		}
	}

	return nil
}

// drawTextLoresBlock draws a number of lines of text or lores from start to end
func drawTextLoresBlock(screen *ebiten.Image, start int, end int, drawer drawTextLoresByte) error {
	for y := start; y < end; y++ {
		base := 128*(y%8) + 40*(y/8)

		// Flip to the 2nd page if so toggled
		if mmu.Page2 {
			base += 0x400
		}

		for x := 0; x < 40; x++ {
			offset := textVideoMemory + base + x
			value := mmu.ReadPageTable[offset>>8][offset&0xff]
			if err := drawer(screen, x, y, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// drawTextBlock draws a number of lines of text from start to end
func drawTextBlock(screen *ebiten.Image, start int, end int) error {
	drawTextLoresBlock(screen, start, end, drawText)
	return nil
}

// drawTextBlock draws a number of lores lines from the equivalent text start to end line
func drawLoresBlock(screen *ebiten.Image, start int, end int) error {
	drawTextLoresBlock(screen, start, end, drawLores)
	return nil
}

// drawTextOrLoresScreen draws a text and/or lores screen depending on the VideoState
func drawTextOrLoresScreen(screen *ebiten.Image) error {
	topHalfIsLowRes := !mmu.VideoState.TextMode
	bottomHalfIsLowRes := !mmu.VideoState.TextMode && !mmu.VideoState.Mixed

	if !topHalfIsLowRes {
		drawTextBlock(screen, 0, 20)
	} else {
		drawLoresBlock(screen, 0, 20)
	}

	if !bottomHalfIsLowRes {
		drawTextBlock(screen, 20, 24)
	} else {
		drawLoresBlock(screen, 20, 24)
	}

	return nil
}

// drawHiresScreen draws an entire hires screen. If it's in mixed mode, the lower end is drawn in text.
func drawHiresScreen(screen *ebiten.Image) error {
	if ScreenSizeFactor != 1 {
		panic("Hires mode for ScreenSizeFactor != 1 not implemented")
	}

	pixels := make([]byte, 280*192*4)

	// Loop over all hires lines
	for y := 0; y < 192; y++ {
		if mmu.VideoState.Mixed && y >= 160 {
			continue
		}

		// Woz is a genius
		yOffset := 0x2000 - (0x3d8)*(y>>6) + 0x80*(y>>3) + 0x400*(y&0x7)

		// Flip to the 2nd page if so toggled
		if mmu.Page2 {
			yOffset += 0x2000
		}

		// For each byte, flip the 7 bits and write it to the pixels array
		for x := 0; x < 40; x++ {
			offset := yOffset + x
			value := mmu.ReadPageTable[offset>>8][offset&0xff]
			value &= 0x7f

			for bit := 0; bit < 7; bit++ {
				b := float64(value & 1)
				value = value >> 1
				p := (y*280 + x*7 + bit) * 4

				pixels[p+0] = byte(0xff * float64(0.20) * b)
				pixels[p+1] = byte(0xff * float64(0.75) * b)
				pixels[p+2] = byte(0xff * float64(0.20) * b)
				pixels[p+3] = 0xff
			}
		}
	}

	// The hires pixels are read, flush them to the screen
	screen.ReplacePixels(pixels)

	// Draw text bit at the bottom
	if mmu.VideoState.Mixed {
		drawTextBlock(screen, 20, 24)
	}

	return nil
}

// DrawScreen draws a text, lores, hires or combination screen
func DrawScreen(screen *ebiten.Image) error {
	flashCounter--
	if flashCounter < 0 {
		flashCounter = flashFrames
		flashOn = !flashOn
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	if !mmu.VideoState.HiresMode {
		drawTextOrLoresScreen(screen)
	} else {
		drawHiresScreen(screen)
	}

	if ShowFPS {
		msg := fmt.Sprintf(`FPS: %0.2f`, ebiten.CurrentFPS())
		ebitenutil.DebugPrint(screen, msg)
	}

	return nil
}

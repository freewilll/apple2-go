package video

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/freewilll/apple2-go/mmu"
)

const (
	textVideoMemory = 0x400 // Base location of page 1 text video memory
	flashFrames     = 11    // Number of frames when FLASH mode is toggled
)

// drawTextLoresByte is a function definition used for mixed text/lores rendering
type drawTextLoresByte func(*ebiten.Image, int, int, uint8) error

var (
	flashCounter           int               // Counter used for flashing characters on the text screen
	flashOn                bool              // Are we currently flashing?
	monochromeLoresSquares [16]*ebiten.Image // Monochrome blocks for lores rendering
	colorLoresSquares      [16]*ebiten.Image // Colored blocks for lores rendering
	colors                 [16]color.NRGBA   // 4-bit Colors

	// ShowFPS determines if the FPS is shown in the corner of the video
	ShowFPS    bool
	Monochrome bool
)

// initLoresSquares creates 16 colored squares for the lores renderer
func initLoresSquares() {
	var err error
	for i := 0; i < 16; i++ {
		monochromeLoresSquares[i], err = ebiten.NewImage(7, 4, ebiten.FilterNearest)
		if err != nil {
			panic(err)
		}

		colorLoresSquares[i], err = ebiten.NewImage(7, 4, ebiten.FilterNearest)
		if err != nil {
			panic(err)
		}
	}

	// From
	// https://mrob.com/pub/xgithub.com/freewilll/apple2/colors.html
	// https://archive.org/details/IIgs_2523063_Master_Color_Values
	alpha := uint8(0xff)

	colors[0x00] = color.NRGBA{0, 0, 0, alpha}
	colors[0x01] = color.NRGBA{221, 0, 51, alpha}
	colors[0x02] = color.NRGBA{0, 0, 153, alpha}
	colors[0x03] = color.NRGBA{221, 34, 221, alpha}
	colors[0x04] = color.NRGBA{0, 119, 34, alpha}
	colors[0x05] = color.NRGBA{85, 85, 85, alpha}
	colors[0x06] = color.NRGBA{34, 34, 255, alpha}
	colors[0x07] = color.NRGBA{102, 170, 255, alpha}
	colors[0x08] = color.NRGBA{136, 85, 0, alpha}
	colors[0x09] = color.NRGBA{255, 102, 0, alpha}
	colors[0x0A] = color.NRGBA{170, 170, 170, alpha}
	colors[0x0B] = color.NRGBA{255, 153, 136, alpha}
	colors[0x0C] = color.NRGBA{17, 221, 0, alpha}
	colors[0x0D] = color.NRGBA{255, 255, 0, alpha}
	colors[0x0E] = color.NRGBA{68, 255, 153, alpha}
	colors[0x0F] = color.NRGBA{255, 255, 255, alpha}

	for i := 0; i < 0x10; i++ {
		colorLoresSquares[i].Fill(colors[i])
		avgIntensity := float64(int(colors[i].R)+int(colors[i].G)+int(colors[i].B)) / 3
		avgColor := color.NRGBA{byte(avgIntensity * 0.2), byte(avgIntensity * 0.75), byte(avgIntensity * 0.2), alpha}
		monochromeLoresSquares[i].Fill(avgColor)
	}
}

// Init the video data structures used for rendering
func Init() {
	ShowFPS = false
	Monochrome = true

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
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(2*7*float64(x), 2*8*float64(y))

	r := image.Rect(0, 0, 7, 8)
	op.SourceRect = &r

	// The charMap is already inverted. Invert it back if we ourselves aren't inverted.
	if inverted {
		op.ColorM.Scale(-1, -1, -1, 1)
		op.ColorM.Translate(1, 1, 1, 0)
	}

	if Monochrome {
		// Make it look greenish
		op.ColorM.Scale(0.20, 0.75, 0.20, 1)
	}

	return screen.DrawImage(charMap[value], op)
}

// drawLores draws two colored lores squares at the equivalent text location x,y.
func drawLores(screen *ebiten.Image, x int, y int, value uint8) error {
	// Convert the 8 bit value to two 4 bit values
	var values = [2]uint8{value & 0xf, value >> 4}

	// Render top & bottom squares
	for i := 0; i < 2; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(2*7*float64(x), 2*8*float64(y)+2*float64(i)*4)

		var loresSquare *ebiten.Image
		if Monochrome {
			loresSquare = monochromeLoresSquares[values[i]]
		} else {
			loresSquare = colorLoresSquares[values[i]]
		}

		if err := screen.DrawImage(loresSquare, op); err != nil {
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
	pixels := make([]byte, 560*384*4)
	halfPixels := make([]byte, 14)

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

		// Initialize 4-bit 4-color and 14 half-pixels
		var color uint8      // Current 4-bit color
		colorPos := uint8(0) // Current half-pixel in the 4-bit color
		for i := 0; i < 14; i++ {
			halfPixels[i] = 0
		}

		// For each byte, expand the 7 bits to the 14 half-pixels array
		// If the high bit is set, shift one half pixel over.
		// Don't shift half-bits in monochrome mode
		for x := 0; x < 40; x++ {
			offset := yOffset + x
			value := mmu.ReadPageTable[offset>>8][offset&0xff]

			phaseShifted := value >> 7

			var hp uint8
			if Monochrome {
				hp = 0
			} else {
				hp = phaseShifted
			}

			halfPixels[0] = halfPixels[13] // Rotate the last phase shifted pixel in

			// Double up the pixels into half pixels starting at offset hp
			for bit := 0; bit < 7; bit++ {
				halfPixels[hp] = value & 1
				hp = hp + 1
				if hp < 14 {
					halfPixels[hp] = value & 1
					hp = hp + 1
				}
				value = value >> 1
			}

			for hp = 0; hp < 14; hp++ {
				// Update the color bit in colorPos with the half pixel value
				color &= ((1 << colorPos) ^ 0xf)
				color |= halfPixels[hp] << colorPos
				colorPos = (colorPos + 1) & 3

				// Draw two lines at a time
				for rowDouble := 0; rowDouble < 2; rowDouble++ {
					p := ((y*2+rowDouble)*560 + x*2*7 + int(hp)) * 4

					if Monochrome {
						b := float64(halfPixels[hp])
						pixels[p+0] = byte(0xff * float64(0.20) * b)
						pixels[p+1] = byte(0xff * float64(0.75) * b)
						pixels[p+2] = byte(0xff * float64(0.20) * b)
						pixels[p+3] = 0xff
					} else {
						pixels[p+0] = colors[color].R
						pixels[p+1] = colors[color].G
						pixels[p+2] = colors[color].B
						pixels[p+3] = 0xff
					}
				}
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

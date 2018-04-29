package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenSizeFactor = 1     // Factor by which the whole screen is resized
	textVideoMemory  = 0x400 // Base location of page 1 text video memory
	flashFrames      = 8     // Number of frames when FLASH mode is toggled
)

var (
	ebitenImage  *ebiten.Image
	memory       [0x100000]byte
	flashCounter int
	flashOn      bool
)

func drawTextScreen(screen *ebiten.Image) error {
	flashCounter--
	if flashCounter < 0 {
		flashCounter = flashFrames
		flashOn = !flashOn
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	for y := 0; y < 24; y++ {
		base := 128*(y%8) + 40*(y/8)
		for x := 0; x < 40; x++ {
			offset := textVideoMemory + base + x
			value := memory[offset]
			inverted := false

			if (value & 0xc0) == 0 {
				inverted = true
			} else if (value & 0x80) == 0 {
				value = value & 0x3f
				inverted = flashOn
			}

			if !inverted {
				value = value & 0x7f
			}

			if value < 0x20 {
				value += 0x40
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(screenSizeFactor, screenSizeFactor)
			op.GeoM.Translate(screenSizeFactor*7*float64(x), screenSizeFactor*8*float64(y))

			fontRow := value % 16
			fontCol := value / 16
			var fontX = (int)(15 + fontCol*12)
			var fontY = (int)(32 + fontRow*11)
			r := image.Rect(fontX, fontY, fontX+7, fontY+8)
			op.SourceRect = &r

			if !inverted {
				op.ColorM.Scale(-1, -1, -1, 1)
				op.ColorM.Translate(1, 1, 1, 0)
			}

			op.ColorM.Scale(0.20, 0.75, 0.20, 1)

			if err := screen.DrawImage(ebitenImage, op); err != nil {
				return err
			}
		}
	}

	return nil
}

func update(screen *ebiten.Image) error {
	return drawTextScreen(screen)
}

func addTestTextScreenData() {
	// Clear screen
	for i := 0; i < 0x400; i++ {
		memory[textVideoMemory+i] = 160
	}

	for i := 0; i < 255; i++ {
		memory[textVideoMemory+i] = byte(i)
	}
}

func main() {
	addTestTextScreenData()

	var err error
	ebitenImage, _, err = ebitenutil.NewImageFromFile("./pr-latin1.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.Run(update, 280*screenSizeFactor, 192*screenSizeFactor, 2, "Apple //")
}

package video

import (
	"image"

	"mos6502go/mmu"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenSizeFactor = 1     // Factor by which the whole screen is resized
	textVideoMemory  = 0x400 // Base location of page 1 text video memory
	flashFrames      = 11    // Number of frames when FLASH mode is toggled
)

var (
	charMap      *ebiten.Image
	flashCounter int
	flashOn      bool
)

func Init() {
	var err error
	charMap, _, err = ebitenutil.NewImageFromFile("video/pr-latin1.png", ebiten.FilterNearest)
	if err != nil {
		panic(err)
	}

}

func DrawTextScreen(pageTable *mmu.PageTable, screen *ebiten.Image) error {
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
			value := (*pageTable)[uint8(offset>>8)][uint8(offset&0xff)]
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

			if err := screen.DrawImage(charMap, op); err != nil {
				return err
			}
		}
	}

	return nil
}

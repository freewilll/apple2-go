package main

import (
	"flag"

	"github.com/hajimehoshi/ebiten"

	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/video"
)

const (
	screenSizeFactor = 1     // Factor by which the whole screen is resized
	textVideoMemory  = 0x400 // Base location of page 1 text video memory
	flashFrames      = 8     // Number of frames when FLASH mode is toggled
)

var showInstructions *bool
var disableBell *bool
var resetKeysDown bool

func reset() {
	bootVector := 0xfffc
	lsb := mmu.PageTable[bootVector>>8][bootVector&0xff] // TODO move readMemory to mmu
	msb := mmu.PageTable[(bootVector+1)>>8][(bootVector+1)&0xff]
	cpu.State.PC = uint16(lsb) + uint16(msb)<<8
}

// checkResetKeys check ctrl-alt-R has been pressed. Releasing the R does a warm reset
func checkResetKeys() {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyR) {
		resetKeysDown = true
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && !ebiten.IsKeyPressed(ebiten.KeyR) && resetKeysDown {
		resetKeysDown = false
		reset()
	} else {
		resetKeysDown = false
	}
}

func update(screen *ebiten.Image) error {
	keyboard.Poll()
	checkResetKeys()

	cpu.Run(*showInstructions, nil, *disableBell, 1024000/60)
	return video.DrawTextScreen(screen)
}

func main() {
	showInstructions = flag.Bool("show-instructions", false, "Show instructions code while running")
	disableBell = flag.Bool("disable-bell", false, "Disable bell")
	flag.Parse()

	cpu.InitDisasm()
	mmu.InitRAM()
	mmu.InitApple2eROM()

	cpu.Init()

	keyboard.Init()
	video.Init()

	reset()

	ebiten.Run(update, 280*screenSizeFactor, 192*screenSizeFactor, 2, "Apple //e")
}

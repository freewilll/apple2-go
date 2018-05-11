package main

import (
	"flag"

	"github.com/hajimehoshi/ebiten"

	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/utils"
	"mos6502go/video"
)

const (
	screenSizeFactor = 1     // Factor by which the whole screen is resized
	textVideoMemory  = 0x400 // Base location of page 1 text video memory
	flashFrames      = 8     // Number of frames when FLASH mode is toggled
)

var showInstructions *bool
var disableFirmwareWait *bool
var resetKeysDown bool
var breakAddress *uint16

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

	cpu.Run(*showInstructions, breakAddress, *disableFirmwareWait, 1024000/60)
	return video.DrawTextScreen(screen)
}

func main() {
	showInstructions = flag.Bool("show-instructions", false, "Show instructions code while running")
	disableFirmwareWait = flag.Bool("disable-wait", false, "Ignore JSRs to firmware wait at $FCA8")
	breakAddressString := flag.String("break", "", "Break on address")
	diskImage := flag.String("image", "", "Disk Image")

	flag.Parse()

	breakAddress = utils.DecodeCmdLineAddress(breakAddressString)

	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()

	if *diskImage != "" {
		mmu.ReadDiskImage(*diskImage)
	}

	cpu.Init()

	keyboard.Init()
	video.Init()

	reset()

	ebiten.Run(update, 280*screenSizeFactor, 192*screenSizeFactor, 2, "Apple //e")
}

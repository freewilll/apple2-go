package main

import (
	"flag"
	"fmt"

	"github.com/hajimehoshi/ebiten"

	"mos6502go/audio"
	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/system"
	"mos6502go/utils"
	"mos6502go/video"
)

var showInstructions *bool
var disableFirmwareWait *bool
var resetKeysDown bool
var fpsKeysDown bool
var breakAddress *uint16

// checkSpecialKeys checks
// - ctrl-alt-R has been pressed. Releasing the R does a warm reset
// - ctrl-alt-F has been pressed, toggling FPS display
func checkSpecialKeys() {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyR) {
		resetKeysDown = true
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && !ebiten.IsKeyPressed(ebiten.KeyR) && resetKeysDown {
		resetKeysDown = false
		cpu.Reset()

	} else {
		resetKeysDown = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyF) {
		fpsKeysDown = true
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && !ebiten.IsKeyPressed(ebiten.KeyF) && fpsKeysDown {
		fpsKeysDown = false
		video.ShowFPS = !video.ShowFPS
		fmt.Println("Toggled")

	} else {
		fpsKeysDown = false
	}
}

func update(screen *ebiten.Image) error {
	keyboard.Poll()
	checkSpecialKeys()

	system.FrameCycles = 0
	system.LastAudioCycles = 0
	exitAtBreak := true
	cpu.Run(*showInstructions, breakAddress, exitAtBreak, *disableFirmwareWait, system.CpuFrequency/60)
	audio.ForwardToFrameCycle()
	system.Cycles += system.FrameCycles
	system.FrameCycles = 0

	return video.DrawScreen(screen)
}

func main() {
	showInstructions = flag.Bool("show-instructions", false, "Show instructions code while running")
	disableFirmwareWait = flag.Bool("disable-wait", false, "Ignore JSRs to firmware wait at $FCA8")
	breakAddressString := flag.String("break", "", "Break on address")
	mute := flag.Bool("mute", false, "Mute sound")
	diskImage := flag.String("image", "", "Disk Image")
	flag.Parse()

	breakAddress = utils.DecodeCmdLineAddress(breakAddressString)

	ebiten.SetRunnableInBackground(true)

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
	audio.InitEbiten()
	audio.Mute = *mute
	system.Init()
	cpu.Reset()

	ebiten.Run(update, 280*video.ScreenSizeFactor, 192*video.ScreenSizeFactor, 2, "Apple //e")

	mmu.FlushImage()
}

package main

// Main emulator executable. This integrates all the components of the emulator
// and contains the ebiten main loop.

import (
	"flag"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten"

	"github.com/freewilll/apple2-go/audio"
	"github.com/freewilll/apple2-go/cpu"
	"github.com/freewilll/apple2-go/disk"
	"github.com/freewilll/apple2-go/keyboard"
	"github.com/freewilll/apple2-go/mmu"
	"github.com/freewilll/apple2-go/system"
	"github.com/freewilll/apple2-go/utils"
	"github.com/freewilll/apple2-go/video"
)

var (
	showInstructions    *bool   // Display all instructions as they are executed
	disableFirmwareWait *bool   // Disable the WAIT function at $fca8
	disableDosDelay     *bool   // Disable DOS delay functions
	breakAddress        *uint16 // Break address from the command line
	scale               float64 // Scale

	resetKeysDown      bool // Keep track of ctrl-alt-R key down state
	fpsKeysDown        bool // Keep track of ctrl-alt-F key down state
	monochromeKeysDown bool // Keep track of ctrl-alt-M key down state
)

// checkSpecialKeys checks
// - ctrl-alt-R has been pressed. Releasing the R does a warm reset
// - ctrl-alt-F has been pressed, toggling FPS display
func checkSpecialKeys() {
	// Check for ctrl-alt-R, and if released, do a warm CPU reset
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyR) {
		resetKeysDown = true
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && !ebiten.IsKeyPressed(ebiten.KeyR) && resetKeysDown {
		resetKeysDown = false
		cpu.Reset()
	} else {
		resetKeysDown = false
	}

	// Check for ctrl-alt-F and toggle FPS display
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyF) {
		fpsKeysDown = true
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && !ebiten.IsKeyPressed(ebiten.KeyF) && fpsKeysDown {
		fpsKeysDown = false
		video.ShowFPS = !video.ShowFPS
	} else {
		fpsKeysDown = false
	}

	// Check for ctrl-alt-M and toggle FPS display
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyM) {
		monochromeKeysDown = true
	} else if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && !ebiten.IsKeyPressed(ebiten.KeyM) && monochromeKeysDown {
		monochromeKeysDown = false
		video.Monochrome = !video.Monochrome
	} else {
		monochromeKeysDown = false
	}
}

// update is the main ebiten loop
func update(screen *ebiten.Image) error {

	checkSpecialKeys() // Poll the keyboard and check for R and F keys

	if !(fpsKeysDown || monochromeKeysDown) {
		keyboard.Poll() // Convert ebiten's keyboard state to an interal value
	}

	system.FrameCycles = 0     // Reset cycles processed this frame
	system.LastAudioCycles = 0 // Reset processed audio cycles
	exitAtBreak := true        // Die if a BRK instruction is seen

	// Run for 1/60 of a second, the duration of an ebiten frame
	cpu.Run(*showInstructions, breakAddress, exitAtBreak, *disableFirmwareWait, *disableDosDelay, system.CPUFrequency/60)

	// Process any audio speaker clicks from this frame
	audio.ForwardToFrameCycle()

	// Updated the cycle accounting
	system.Cycles += system.FrameCycles

	// Finally render the screen
	return video.DrawScreen(screen)
}

func main() {
	var Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Synopsis: %s [disk image file]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Options\n")
		flag.PrintDefaults()
	}
	flag.Usage = Usage

	showInstructions = flag.Bool("show-instructions", false, "Show instructions code while running")
	disableFirmwareWait = flag.Bool("disable-wait", false, "Ignore JSRs to firmware wait at $FCA8")
	disableDosDelay = flag.Bool("disable-dos-delay", false, "Ignore DOS ARM move and motor on waits")
	breakAddressString := flag.String("break", "", "Break on address")
	mute := flag.Bool("mute", false, "Mute sound")
	scale := flag.Float64("scale", 2, "Video scale")
	clickWhenDriveHeadMoves := flag.Bool("drive-head-click", false, "Click speaker when drive head moves")
	flag.Parse()

	breakAddress = utils.DecodeCmdLineAddress(breakAddressString)

	cpu.InitInstructionDecoder() // Init the instruction decoder data structures
	mmu.InitRAM()                // Set all switches to bootup values and initialize the page tables
	mmu.InitApple2eROM()         // Load the ROM and init page tables
	mmu.InitIO()                 // Init slots, video and disk image statuses

	// If there is a disk image on the command line, load it
	diskImages := flag.Args()
	if len(diskImages) > 0 {
		disk.ReadDiskImage(diskImages[0])
	}

	cpu.Init()         // Init the CPU registers, interrupts and disable testing code
	keyboard.Init()    // Init the keyboard state and ebiten translation tables
	video.Init()       // Init the video data structures used for rendering
	audio.InitEbiten() // Initialize the audio sets up the ebiten output stream

	audio.Mute = *mute
	audio.ClickWhenDriveHeadMoves = *clickWhenDriveHeadMoves

	system.Init()           // Initialize the system-wide state
	cpu.SetColdStartReset() // Prepare memory to ensure a cold reset
	cpu.Reset()             // Set the CPU and memory states so that a next call to cpu.Run() calls the firmware reset code

	// Start the ebiten main loop
	ebiten.SetRunnableInBackground(true)
	ebiten.Run(update, 560, 384, *scale, "Apple //e")

	// The main loop has ended, flush any data to the disk image if any writes have been done.
	disk.FlushImage()
}

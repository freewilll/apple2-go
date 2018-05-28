package main

import (
	"fmt"
	"testing"

	"github.com/freewilll/apple2/cpu"
	"github.com/freewilll/apple2/mmu"
	"github.com/freewilll/apple2/system"
)

func testBellCycles(delay int) {
	// Add some code to $800
	mmu.WriteMemory(0x800, 0xa9)         // LDA #delay
	mmu.WriteMemory(0x801, uint8(delay)) //
	mmu.WriteMemory(0x802, 0x20)         // JSR $fca8 BELL
	mmu.WriteMemory(0x803, 0xa8)         //
	mmu.WriteMemory(0x804, 0xfc)         //
	mmu.WriteMemory(0x805, 0x00)         // BRK

	// Run the code until the BRK instruction and count the cycles
	showInstructions := false
	breakAddress := uint16(0x805)
	exitAtBreak := false
	disableFirmwareWait := false
	cpu.State.PC = 0x800
	cpu.Run(showInstructions, &breakAddress, exitAtBreak, disableFirmwareWait, system.CPUFrequency*1000)

	// See http://apple2.org.za/gswv/a2zine/GS.WorldView/Resources/USEFUL.TABLES/WAIT.DELAY.CR.txt
	expectedCycles := (26 + 27*delay + 5*delay*delay) / 2

	gotCycles := int(system.FrameCycles - 2) // Exclude the cycles taken by the LDA

	fmt.Printf("Delay %3d ", delay)
	if gotCycles == expectedCycles {
		fmt.Println("OK")

	} else {
		fmt.Printf("Failed expected %6d, got %6d\n", expectedCycles, gotCycles)
	}
}

// TestBell tests the nunber of cycles in the system BELL loop for different
// values of the accumulator. This test was mainly used to diagnose a bug
// related to sound frequencies being incorrect due to invalid cycle
// housekeeping in the CPU branch code.
func TestBell(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	cpu.Init()
	system.Init()

	testBellCycles(1)
	testBellCycles(2)
	testBellCycles(3)
	testBellCycles(4)
	testBellCycles(12)
	testBellCycles(0x10)
	testBellCycles(0x20)
	testBellCycles(0x40)
	testBellCycles(0x80)
	testBellCycles(0xc0)
	testBellCycles(0xff)
}

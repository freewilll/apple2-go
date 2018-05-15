package main

import (
	"fmt"
	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/system"
	"mos6502go/video"
	"testing"
)

func testBellCycles(delay int) {
	cpu.State.PC = 0x800
	mmu.WriteMemory(0x800, 0xa9) // LDA #$xx
	mmu.WriteMemory(0x801, uint8(delay))
	mmu.WriteMemory(0x802, 0x20) // JSR $fca8
	mmu.WriteMemory(0x803, 0xa8)
	mmu.WriteMemory(0x804, 0xfc)
	mmu.WriteMemory(0x805, 0x00) // Break address

	system.FrameCycles = 0
	showInstructions := false
	breakAddress := uint16(0x805)
	exitAtBreak := false
	disableFirmwareWait := false
	cpu.Run(showInstructions, &breakAddress, exitAtBreak, disableFirmwareWait, system.CpuFrequency*1000)

	// See http://apple2.org.za/gswv/a2zine/GS.WorldView/Resources/USEFUL.TABLES/WAIT.DELAY.CR.txt
	expectedCycles := (26 + 27*delay + 5*delay*delay) / 2

	gotCycles := int(system.FrameCycles - 2)
	fmt.Printf("Delay %3d ", delay)
	if gotCycles == expectedCycles {
		fmt.Println("OK")

	} else {
		fmt.Printf("Failed expected %6d, got %6d\n", expectedCycles, gotCycles)
	}
}

func TestBell(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	cpu.Init()
	keyboard.Init()
	video.Init()
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

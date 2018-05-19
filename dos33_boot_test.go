package main

import (
	"fmt"
	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/system"
	"mos6502go/utils"
	"mos6502go/video"
	"testing"
	"time"
)

const dosDiskImage = "dos33.dsk"

func TestDOS33Boot(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	mmu.ReadDiskImage(dosDiskImage)
	cpu.Init()
	keyboard.Init()
	video.Init()
	system.Init()
	cpu.Reset()

	t0 := time.Now()

	// Run until BASIC would execute the program.
	utils.RunUntilBreakPoint(t, 0xd7d2, 5, false, "BASIC NEWSTT")

	elapsed := float64(time.Since(t0) / time.Millisecond)
	fmt.Printf("CPU Cycles:    %d\n", system.FrameCycles)
	fmt.Printf("Time elapsed:  %0.2f ms\n", elapsed)
	fmt.Printf("Speed:         %0.2f cycles/ms\n", float64(system.FrameCycles)/elapsed)
}

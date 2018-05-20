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

const prodosDiskImage = "prodos.dsk"

func TestProdosBoot(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	mmu.ReadDiskImage(prodosDiskImage)
	cpu.Init()
	keyboard.Init()
	video.Init()
	system.Init()
	cpu.SetColdStartReset()
	cpu.Reset()

	t0 := time.Now()

	utils.RunUntilBreakPoint(t, 0xc600, 2, false, "Boot ROM")
	utils.RunUntilBreakPoint(t, 0x0801, 2, false, "Loader")
	utils.RunUntilBreakPoint(t, 0x2000, 3, false, "Relocator")
	utils.RunUntilBreakPoint(t, 0x0080, 1, false, "AUX RAM test")
	utils.RunUntilBreakPoint(t, 0x2932, 1, false, "Relocation done")
	utils.RunUntilBreakPoint(t, 0x21f3, 1, false, "The first JSR $bf00 - ONLINE - get names of one or all online volumes")
	utils.RunUntilBreakPoint(t, 0xd000, 1, false, "First call to MLI kernel")
	// utils.RunUntilBreakPoint(t, 0x0800, 1, false, "BI loader")

	elapsed := float64(time.Since(t0) / time.Millisecond)
	fmt.Printf("CPU Cycles:    %d\n", system.FrameCycles)
	fmt.Printf("Time elapsed:  %0.2f ms\n", elapsed)
	fmt.Printf("Speed:         %0.2f cycles/ms\n", float64(system.FrameCycles)/elapsed)
}

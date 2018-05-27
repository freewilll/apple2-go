package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/freewilll/apple2/cpu"
	"github.com/freewilll/apple2/disk"
	"github.com/freewilll/apple2/keyboard"
	"github.com/freewilll/apple2/mmu"
	"github.com/freewilll/apple2/system"
	"github.com/freewilll/apple2/utils"
	"github.com/freewilll/apple2/video"
)

const prodosDiskImage = "prodos19.dsk"

func TestProdosBoot(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	disk.ReadDiskImage(prodosDiskImage)
	cpu.Init()
	keyboard.Init()
	video.Init()
	system.Init()
	cpu.SetColdStartReset()
	cpu.Reset()

	t0 := time.Now()

	utils.RunUntilBreakPoint(t, 0xc600, 2, false, "Boot ROM")
	utils.RunUntilBreakPoint(t, 0x0801, 2, false, "Loader")
	utils.RunUntilBreakPoint(t, 0x2000, 3, false, "Kernel Relocator")
	utils.RunUntilBreakPoint(t, 0x0080, 1, false, "AUX RAM test")
	utils.RunUntilBreakPoint(t, 0x2932, 1, false, "Relocation done")
	utils.RunUntilBreakPoint(t, 0x21f3, 1, false, "The first JSR $bf00 - ONLINE - get names of one or all online volumes")
	utils.RunUntilBreakPoint(t, 0xd000, 1, false, "First call to MLI kernel")
	utils.RunUntilBreakPoint(t, 0x0800, 2, false, "BI loader")
	utils.RunUntilBreakPoint(t, 0x2000, 2, false, "BI Relocator")
	utils.RunUntilBreakPoint(t, 0xbe00, 52, false, "BI Start")

	elapsed := float64(time.Since(t0) / time.Millisecond)
	fmt.Printf("CPU Cycles:    %d\n", system.FrameCycles)
	fmt.Printf("Time elapsed:  %0.2f ms\n", elapsed)
	fmt.Printf("Speed:         %0.2f cycles/ms\n", float64(system.FrameCycles)/elapsed)
}

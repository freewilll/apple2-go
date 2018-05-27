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

const dosDiskImage = "dos33.dsk"

func runDos33Boot(t *testing.T) {
	// Boot up DOS3.3
	utils.RunUntilBreakPoint(t, 0x0801, 2, false, "Boot0")
	utils.RunUntilBreakPoint(t, 0x3700, 1, false, "Boot1") // $3700 is for master disk, $b700 for slave
	utils.RunUntilBreakPoint(t, 0x9d84, 3, false, "Boot2")
	utils.RunUntilBreakPoint(t, 0xd7d2, 2, false, "JMP to basic interpreter NEWSTT")
}

func TestDOS33Boot(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	disk.ReadDiskImage(dosDiskImage)
	cpu.Init()
	keyboard.Init()
	video.Init()
	system.Init()
	cpu.SetColdStartReset()
	cpu.Reset()

	t0 := time.Now()

	runDos33Boot(t)

	elapsed := float64(time.Since(t0) / time.Millisecond)
	fmt.Printf("CPU Cycles:    %d\n", system.FrameCycles)
	fmt.Printf("Time elapsed:  %0.2f ms\n", elapsed)
	fmt.Printf("Speed:         %0.2f cycles/ms\n", float64(system.FrameCycles)/elapsed)
}

package main

import (
	"fmt"
	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/system"
	"mos6502go/video"
	"testing"
	"time"
)

const DiskImage = "dos33_disk.dsk"

func TestDOS33Boot(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	mmu.ReadDiskImage(DiskImage)
	cpu.Init()
	keyboard.Init()
	video.Init()
	system.Init()
	cpu.Reset()

	system.FrameCycles = 0
	system.LastAudioCycles = 0
	showInstructions := false
	var breakAddress uint16
	disableFirmwareWait := false
	t0 := time.Now()
	cpu.Run(showInstructions, &breakAddress, disableFirmwareWait, system.CpuFrequency*1000)

	elapsed := float64(time.Since(t0) / time.Millisecond)
	fmt.Printf("CPU Cycles:    %d\n", system.FrameCycles)
	fmt.Printf("Time elapsed:  %0.2f ms\n", elapsed)
	fmt.Printf("Speed:         %0.2f cycles/ms\n", float64(system.FrameCycles)/elapsed)
}

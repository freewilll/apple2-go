package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"apple2/cpu"
	"apple2/keyboard"
	"apple2/mmu"
	"apple2/system"
	"apple2/video"
)

func TestIoBankSwitching(t *testing.T) {
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	cpu.Init()
	keyboard.Init()
	video.Init()
	system.Init()
	cpu.SetColdStartReset()
	cpu.Reset()

	mmu.MapFirstHalfOfIO()
	assert.Equal(t, uint8(0xa2), mmu.ReadMemory(0xc600)) // read from Primary Slot 6 ROM
	mmu.MapSecondHalfOfIO()
	assert.Equal(t, uint8(0x8d), mmu.ReadMemory(0xc600)) // read from Primary Slot 6 ROM
}

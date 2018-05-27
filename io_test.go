package main

import (
	"testing"

	"github.com/freewilll/apple2/cpu"
	"github.com/freewilll/apple2/keyboard"
	"github.com/freewilll/apple2/mmu"
	"github.com/freewilll/apple2/system"
	"github.com/freewilll/apple2/video"
	"github.com/stretchr/testify/assert"
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

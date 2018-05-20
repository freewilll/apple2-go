package main

import (
	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/system"
	"mos6502go/video"
	"testing"

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

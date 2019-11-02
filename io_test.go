package main

import (
	"testing"

	"github.com/freewilll/apple2-go/cpu"
	"github.com/freewilll/apple2-go/keyboard"
	"github.com/freewilll/apple2-go/mmu"
	"github.com/freewilll/apple2-go/system"
	"github.com/freewilll/apple2-go/video"
	"github.com/stretchr/testify/assert"
)

// TestIoBankSwitching tests the switching of the IO memory ROM at $c000-$c7ff
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

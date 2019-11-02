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

func assertMemoryConfiguration(t *testing.T, address uint16, upperRAMReadOnly bool, upperReadMappedToROM bool, d000Bank int) {
	mmu.WriteMemory(address, 0x00)
	assert.Equal(t, upperRAMReadOnly, mmu.UpperRAMReadOnly)
	assert.Equal(t, upperReadMappedToROM, mmu.UpperReadMappedToROM)
	assert.Equal(t, d000Bank, mmu.D000Bank)
}

// TestBankSwitching tests the area starting at $d000 and managed by $c08x.
// First the initial settings are checked. Then a bunch of assertions on the
// internal code. Then writes to $c08x.
func TestBankSwitching(t *testing.T) {
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

	// Sanity test that what we expect from the apple //e ROM is correct
	assert.Equal(t, uint8(0x6f), mmu.ReadMemory(0xd000)) // read from ROM
	assert.Equal(t, uint8(0xc3), mmu.ReadMemory(0xffff)) // read from ROM

	// Verify ROM & RAM settings at startup
	mmu.WipeRAM()
	assert.Equal(t, uint8(0xc3), mmu.ReadMemory(0xffff))                // read from ROM
	mmu.WriteMemory(0xffff, 0xff)                                       // write to $ffff
	assert.Equal(t, uint8(0xc3), mmu.ReadMemory(0xffff))                // ROM value is the same
	assert.Equal(t, uint8(0xff), mmu.PhysicalMemory.MainMemory[0xffff]) // RAM has been updated
	mmu.WriteMemory(0xd000, 0xfe)                                       // write to $d000
	assert.Equal(t, uint8(0x00), mmu.PhysicalMemory.MainMemory[0xc000]) // bank #1 RAM
	assert.Equal(t, uint8(0xfe), mmu.PhysicalMemory.MainMemory[0xd000]) // bank #2 RAM

	// Switch bank to 1, write and check physical memory
	mmu.SetD000Bank(1)
	mmu.SetUpperReadMappedToROM(false)
	mmu.WriteMemory(0xd000, 0xfd)                                       // write to $d000
	assert.Equal(t, uint8(0xfd), mmu.PhysicalMemory.MainMemory[0xc000]) // bank #1 RAM
	assert.Equal(t, uint8(0xfe), mmu.PhysicalMemory.MainMemory[0xd000]) // bank #2 RAM

	// Enable RAM area for reading and check values
	mmu.SetUpperReadMappedToROM(false)
	assert.Equal(t, uint8(0xfd), mmu.ReadMemory(0xd000)) // read from bank #1 RAM
	mmu.SetD000Bank(2)
	assert.Equal(t, uint8(0xfe), mmu.ReadMemory(0xd000)) // read from bank #1 RAM

	// Enable ROM area for reading and check values
	mmu.SetUpperReadMappedToROM(true)
	assert.Equal(t, uint8(0x6f), mmu.ReadMemory(0xd000)) // read from ROM
	assert.Equal(t, uint8(0xc3), mmu.ReadMemory(0xffff)) // read from ROM

	// Set d000 RAM to bank 1, RAM to read only and attempt writes
	mmu.SetD000Bank(1)
	mmu.SetUpperRAMReadOnly(true)
	assert.Equal(t, uint8(0xfd), mmu.PhysicalMemory.MainMemory[0xc000]) // bank #1 RAM
	assert.Equal(t, uint8(0xfe), mmu.PhysicalMemory.MainMemory[0xd000]) // bank #2 RAM
	mmu.WriteMemory(0xd000, 0x01)                                       // attempt to write to read only RAM
	mmu.WriteMemory(0xffff, 0x02)                                       // attempt to write to read only RAM
	assert.Equal(t, uint8(0xfd), mmu.PhysicalMemory.MainMemory[0xc000]) // bank #1 RAM is unchanged
	assert.Equal(t, uint8(0xfe), mmu.PhysicalMemory.MainMemory[0xd000]) // bank #2 RAM is unchanged
	assert.Equal(t, uint8(0xff), mmu.PhysicalMemory.MainMemory[0xffff]) // top of RAM is unchanged

	// Set RAM to write and write to it
	mmu.SetUpperRAMReadOnly(false)
	mmu.WriteMemory(0xd000, 0xfc)                                       // write to RAM
	mmu.WriteMemory(0xffff, 0xfb)                                       // write to RAM
	assert.Equal(t, uint8(0xfc), mmu.PhysicalMemory.MainMemory[0xc000]) // bank #1 RAM has been updated
	assert.Equal(t, uint8(0xfe), mmu.PhysicalMemory.MainMemory[0xd000]) // bank #2 RAM is untouched
	assert.Equal(t, uint8(0xfb), mmu.PhysicalMemory.MainMemory[0xffff]) // top of RAM has been updated

	// Enable ROM area for reading and check values
	mmu.SetUpperReadMappedToROM(true)
	assert.Equal(t, uint8(0x6f), mmu.ReadMemory(0xd000)) // read from ROM
	assert.Equal(t, uint8(0xc3), mmu.ReadMemory(0xffff)) // read from ROM

	// Test writes to 0xc08x lead to correct memory configurations
	assertMemoryConfiguration(t, 0xc080, true, false, 2)
	assertMemoryConfiguration(t, 0xc081, false, true, 2)
	assertMemoryConfiguration(t, 0xc082, true, true, 2)
	assertMemoryConfiguration(t, 0xc083, false, false, 2)
	assertMemoryConfiguration(t, 0xc084, true, false, 2)
	assertMemoryConfiguration(t, 0xc085, false, true, 2)
	assertMemoryConfiguration(t, 0xc086, true, true, 2)
	assertMemoryConfiguration(t, 0xc087, false, false, 2)
	assertMemoryConfiguration(t, 0xc088, true, false, 1)
	assertMemoryConfiguration(t, 0xc089, false, true, 1)
	assertMemoryConfiguration(t, 0xc08a, true, true, 1)
	assertMemoryConfiguration(t, 0xc08b, false, false, 1)
	assertMemoryConfiguration(t, 0xc08c, true, false, 1)
	assertMemoryConfiguration(t, 0xc08d, false, true, 1)
	assertMemoryConfiguration(t, 0xc08e, true, true, 1)
	assertMemoryConfiguration(t, 0xc08f, false, false, 1)
}

package mmu

import (
	"fmt"
	"io/ioutil"
)

const RomPath = "apple2e.rom"

const StackPage = 1

// https://mirrors.apple2.org.za/apple.cabi.net/Languages.Programming/MemoryMap.IIe.64K.128K.txt

const (
	KEYBOARD = 0xC000 // keyboard data (latched) (RD-only)
	CLR80COL = 0xC000 // use 80-column memory mapping (WR-only)
	SET80COL = 0xC001
	CLRAUXRD = 0xC002 // read from auxilliary 48K
	SETAUXRD = 0xC003
	CLRAUXWR = 0xC004 // write to auxilliary 48K
	SETAUXWR = 0xC005
	CLRCXROM = 0xC006 // use external slot ROM
	SETCXROM = 0xC007
	CLRAUXZP = 0xC008 // use auxilliary ZP, stack, & LC
	SETAUXZP = 0xC009
	CLRC3ROM = 0xC00A // use external slot C3 ROM
	SETC3ROM = 0xC00B
	CLR80VID = 0xC00C // use 80-column display mode
	SET80VID = 0xC00D
	CLRALTCH = 0xC00E // use alternate character set ROM
	SETALTCH = 0xC00F
	STROBE   = 0xC010 // strobe (unlatch) keyboard data
)

type PhysicalMemory struct {
	MainMemory [0xc000]uint8
	UpperROM   [0x3000]uint8
	RomC1      [0x1000]uint8
	RomC2      [0x1000]uint8
}

type PageTable [0x100][]uint8

type Memory struct {
	PageTable      PageTable
	PhysicalMemory PhysicalMemory
}

func MapFirstHalfOfIO(m *Memory) {
	for i := 0x1; i < 0x10; i++ {
		m.PageTable[i+0xc0] = m.PhysicalMemory.RomC1[i*0x100 : i*0x100+0x100]
	}
}

func MapSecondHalfOfIO(m *Memory) {
	for i := 0x1; i < 0x10; i++ {
		m.PageTable[i+0xc0] = m.PhysicalMemory.RomC2[i*0x100 : i*0x100+0x100]
	}
}

func readApple2eROM(m *Memory) {
	bytes, err := ioutil.ReadFile(RomPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to read ROM: %s", err))
	}

	// Copy both I/O areas over c000-cfff, including unused c000-c0ff
	for i := 0x0000; i < 0x1000; i++ {
		m.PhysicalMemory.RomC1[i] = bytes[i]
		m.PhysicalMemory.RomC2[i] = bytes[i+0x4000]
	}

	// Copy ROM over for 0xd000-0xffff area
	for i := 0x0; i < 0x3000; i++ {
		m.PhysicalMemory.UpperROM[i] = bytes[i+0x1000]
	}
}

func InitApple2eROM(m *Memory) {
	readApple2eROM(m)

	// Map 0xc100-0xcfff
	MapFirstHalfOfIO(m)

	// Map 0xd000-0xffff
	for i := 0x0; i < 0x30; i++ {
		m.PageTable[i+0xd0] = m.PhysicalMemory.UpperROM[i*0x100 : i*0x100+0x100]
	}
}

func InitRAM() (memory *Memory) {
	memory = new(Memory)

	// Map main RAM
	for i := 0x0; i < 0xc0; i++ {
		memory.PageTable[i] = memory.PhysicalMemory.MainMemory[i*0x100 : i*0x100+0x100]
	}

	return
}

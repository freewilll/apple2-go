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

	RDLCBNK2 = 0xC011 // reading from LC bank $Dx 2
	RDLCRAM  = 0xC012 // reading from LC RAM
	RDRAMRD  = 0xC013 // reading from auxilliary 48K
	RDRAMWR  = 0xC014 // writing to auxilliary 48K
	RDCXROM  = 0xC015 // using external slot ROM
	RDAUXZP  = 0xC016 // using auxilliary ZP, stack, & LC
	RDC3ROM  = 0xC017 // using external slot C3 ROM
	RD80COL  = 0xC018 // using 80-column memory mapping
	RDVBLBAR = 0xC019 // not VBL (VBL signal low)
	RDTEXT   = 0xC01A // using text mode
	RDMIXED  = 0xC01B // using mixed mode
	RDPAGE2  = 0xC01C // using text/graphics page2
	RDHIRES  = 0xC01D // using Hi-res graphics mode
	RDALTCH  = 0xC01E // using alternate character set ROM
	RD80VID  = 0xC01F // using 80-column display mode

)

var PhysicalMemory struct {
	MainMemory [0xc000]uint8
	UpperROM   [0x3000]uint8
	RomC1      [0x1000]uint8
	RomC2      [0x1000]uint8
}

var PageTable [0x100][]uint8

func MapFirstHalfOfIO() {
	for i := 0x1; i < 0x10; i++ {
		PageTable[i+0xc0] = PhysicalMemory.RomC1[i*0x100 : i*0x100+0x100]
	}
}

func MapSecondHalfOfIO() {
	for i := 0x1; i < 0x10; i++ {
		PageTable[i+0xc0] = PhysicalMemory.RomC2[i*0x100 : i*0x100+0x100]
	}
}

// emptySlot zeroes all RAM for a slot
func emptySlot(slot int) {
	for i := slot * 0x100; i < (slot+1)*0x100; i++ {
		PhysicalMemory.RomC1[i] = 0
		PhysicalMemory.RomC2[i] = 0
	}
}

func readApple2eROM() {
	bytes, err := ioutil.ReadFile(RomPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to read ROM: %s", err))
	}

	// Copy both I/O areas over c000-cfff, including unused c000-c0ff
	for i := 0x0000; i < 0x1000; i++ {
		PhysicalMemory.RomC1[i] = bytes[i]
		PhysicalMemory.RomC2[i] = bytes[i+0x4000]
	}

	// Copy ROM over for 0xd000-0xffff area
	for i := 0x0; i < 0x3000; i++ {
		PhysicalMemory.UpperROM[i] = bytes[i+0x1000]
	}

	// Empty slots that aren't yet implemented
	emptySlot(3)
	emptySlot(4)
	emptySlot(6)
	emptySlot(7)
}

func InitApple2eROM() {
	readApple2eROM()

	// Map 0xc100-0xcfff
	MapFirstHalfOfIO()

	// Map 0xd000-0xffff
	for i := 0x0; i < 0x30; i++ {
		PageTable[i+0xd0] = PhysicalMemory.UpperROM[i*0x100 : i*0x100+0x100]
	}
}

func InitRAM() {
	// Map main RAM
	for i := 0x0; i < 0xc0; i++ {
		PageTable[i] = PhysicalMemory.MainMemory[i*0x100 : i*0x100+0x100]
	}

	return
}

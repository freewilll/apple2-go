package mmu

import (
	"fmt"
	"io/ioutil"
	"mos6502go/system"
)

const RomPath = "apple2e.rom"
const StackPage = 1

var PhysicalMemory struct {
	MainMemory [0x10000]uint8
	UpperROM   [0x3000]uint8
	RomC1      [0x1000]uint8
	RomC2      [0x1000]uint8
}

var ReadPageTable [0x100][]uint8
var WritePageTable [0x100][]uint8

// Memory mapping states
var (
	D000Bank             int  // one maps to $c000, two maps to $d000
	UsingExternalSlotRom bool // Which IO ROM is being used
	UpperReadMappedToROM bool // Do reads go to the RAM or ROM
	UpperRamReadOnly     bool // Is the upper RAM read only
)

func ApplyMemoryConfiguration() {
	// Map main RAM for read/write
	for i := 0x0; i < 0xc0; i++ {
		ReadPageTable[i] = PhysicalMemory.MainMemory[i*0x100 : i*0x100+0x100]
		WritePageTable[i] = PhysicalMemory.MainMemory[i*0x100 : i*0x100+0x100]
	}

	// Map $c000
	var ioRom *[0x1000]uint8
	if UsingExternalSlotRom {
		ioRom = &PhysicalMemory.RomC2
	} else {
		ioRom = &PhysicalMemory.RomC1
	}

	for i := 0x1; i < 0x10; i++ {
		ReadPageTable[0xc0+i] = (*ioRom)[i*0x100 : i*0x100+0x100]
		WritePageTable[0xc0+i] = nil
	}

	// Map $d000
	for i := 0xd0; i < 0xe0; i++ {
		base := i*0x100 + D000Bank*0x1000 - 0x2000
		if !UpperReadMappedToROM {
			ReadPageTable[i] = PhysicalMemory.MainMemory[base : base+0x100]
		}

		if UpperRamReadOnly {
			WritePageTable[i] = nil
		} else {
			WritePageTable[i] = PhysicalMemory.MainMemory[base : base+0x100]
		}
	}

	// Map 0xe00 to 0xffff
	for i := 0xe0; i < 0x100; i++ {
		base := i * 0x100
		if !UpperReadMappedToROM {
			ReadPageTable[i] = PhysicalMemory.MainMemory[base : base+0x100]
		}
		if UpperRamReadOnly {
			WritePageTable[i] = nil
		} else {
			WritePageTable[i] = PhysicalMemory.MainMemory[base : base+0x100]
		}
	}

	if UpperReadMappedToROM {
		for i := 0x00; i < 0x30; i++ {
			ReadPageTable[i+0xd0] = PhysicalMemory.UpperROM[i*0x100 : i*0x100+0x100]
		}
	}

}

// Map 0xc100-0xcfff for reading from RomC1
func MapFirstHalfOfIO() {
	UsingExternalSlotRom = false
	ApplyMemoryConfiguration()
}

// Map 0xc100-0xcfff for reading from RomC2
func MapSecondHalfOfIO() {
	UsingExternalSlotRom = true
	ApplyMemoryConfiguration()
}

// emptySlot zeroes all RAM for a slot
func emptySlot(slot int) {
	for i := slot * 0x100; i < (slot+1)*0x100; i++ {
		PhysicalMemory.RomC1[i] = 0
		PhysicalMemory.RomC2[i] = 0
	}
}

func loadApple2eROM() {
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
}

func InitApple2eROM() {
	loadApple2eROM()
	MapFirstHalfOfIO() // Map 0xc100-0xcfff for reading
	InitROM()          // Map 0xd000-0xffff for reading
}

func InitROM() {
	UpperReadMappedToROM = true
	ApplyMemoryConfiguration()
}

func SetUpperReadMappedToROM(value bool) {
	UpperReadMappedToROM = value
	ApplyMemoryConfiguration()
}

func SetUpperRamReadOnly(value bool) {
	UpperRamReadOnly = value
	ApplyMemoryConfiguration()
}

func SetD000Bank(value int) {
	D000Bank = value
	ApplyMemoryConfiguration()
}

func InitRAM() {
	UpperRamReadOnly = false
	D000Bank = 2
	ApplyMemoryConfiguration()
}

func WipeRAM() {
	for i := 0; i < 0x10000; i++ {
		PhysicalMemory.MainMemory[i] = 0
	}
}

func SetMemoryMode(mode uint8) {
	// mode corresponds to a read/write to $c080 with
	// $c080 mode=$00
	// $c08f mode=$0f

	if (mode & 1) == 0 {
		UpperRamReadOnly = true
	} else {
		UpperRamReadOnly = false
	}

	if (((mode & 2) >> 1) ^ (mode & 1)) == 0 {
		UpperReadMappedToROM = false

	} else {
		UpperReadMappedToROM = true
	}

	if (mode & 8) == 0 {
		D000Bank = 2
	} else {
		D000Bank = 1
	}

	ApplyMemoryConfiguration()
}

func ReadMemory(address uint16) uint8 {
	if (address >= 0xc000) && (address < 0xc100) {
		return ReadIO(address)
	} else {
		return ReadPageTable[address>>8][address&0xff]
	}
}

func WriteMemory(address uint16, value uint8) {
	if (address >= 0xc000) && (address < 0xc100) {
		WriteIO(address, value)
		return
	}

	if system.RunningInterruptTests && address == 0xbffc {
		oldValue := ReadMemory(address)
		system.WriteInterruptTestOpenCollector(address, oldValue, value)
		WritePageTable[uint8(address>>8)][uint8(address&0xff)] = value
		return
	}

	memory := WritePageTable[address>>8]
	// If memory is nil, then it's read only. The write is ignored.
	if memory != nil {
		memory[uint8(address&0xff)] = value
	}

	if system.RunningFunctionalTests && address == 0x200 {
		testNumber := ReadMemory(0x200)
		if testNumber == 0xf0 {
			fmt.Println("Opcode testing completed")
		} else {
			fmt.Printf("Test %d OK\n", ReadMemory(0x200))
		}
	}
}

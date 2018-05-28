package mmu

import (
	"fmt"
	"io/ioutil"

	"github.com/freewilll/apple2/system"
)

const RomPath = "apple2e.rom" // So far only one ROM is supported and it's loaded at startup
const StackPage = 1           // The 6502 stack is at 0x100

// PhysicalMemory contains all the unmapped memory, ROM and RAM
var PhysicalMemory struct {
	MainMemory [0x10000]uint8 // Main RAM
	UpperROM   [0x3000]uint8  // $c000-$ffff ROM area
	RomC1      [0x1000]uint8  // First half of IO ROM
	RomC2      [0x1000]uint8  // Second half of IO ROM
}

// Page tables for read & write
var ReadPageTable [0x100][]uint8
var WritePageTable [0x100][]uint8

// Memory mapping states
var (
	D000Bank             int  // one maps to $c000, two maps to $d000
	UsingExternalSlotRom bool // Which IO ROM is being used
	UpperReadMappedToROM bool // Do reads go to the RAM or ROM
	UpperRamReadOnly     bool // Is the upper RAM read only
	FakeAuxMemoryRead    bool // Aux memory isn't implemented
	FakeAuxMemoryWrite   bool // Aux memory isn't implemented
	FakeAltZP            bool // Aux memory isn't implemented
	FakePage2            bool // Aux memory isn't implemented
	Col80                bool // 80 Column card is on (not implemented)
	Store80              bool // 80 Column card is on (not implemented)
	Page2                bool // Main memory Page2 is selected
)

// Make page tables for current RAM, ROM and IO configuration
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

// emptySlot zeroes all RAM for a slot, effectively disabling the slot
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

// Set upper memory area for reading from ROM
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

// Set d000 bank to map to $c000 or $d000 in the physical  memory
func SetD000Bank(value int) {
	D000Bank = value
	ApplyMemoryConfiguration()
}

// Aux memory hasn't been implemented. If aux memory is selected, and a read
// is attempted, then nonsense must be returned.
func SetFakeAuxMemoryRead(value bool) {
	FakeAuxMemoryRead = value
	ApplyMemoryConfiguration()
}

// Aux memory hasn't been implemented. If aux memory is selected, and a write
// is attempted, then it must be ignored.
func SetFakeAuxMemoryWrite(value bool) {
	FakeAuxMemoryWrite = value
	ApplyMemoryConfiguration()
}

// Alternate zero page isn't implemented
func SetFakeAltZP(value bool) {
	FakeAltZP = value
	ApplyMemoryConfiguration()
}

// 80 column card isn't implemented
func SetCol80(value bool) {
	Col80 = value
	// No changes are needed when this is toggled
}

// Page switching is only implemented for the main memory
func SetPage2(value bool) {
	// If the 80 column card is enabled, then this toggles aux memory
	// Otherwise, page1/page2 is toggled in the main memory
	if Col80 {
		FakePage2 = value
		ApplyMemoryConfiguration()
	} else {
		Page2 = value
	}
}

// 80 column card isn't implemented
func SetStore80(value bool) {
	Store80 = value
	FakePage2 = value
	ApplyMemoryConfiguration()
}

// InitRAM sets all default RAM memory settings and resets the page tables
func InitRAM() {
	UpperRamReadOnly = false
	D000Bank = 2
	FakeAuxMemoryRead = false  // Aux memory isn't implemented
	FakeAuxMemoryWrite = false // Aux memory isn't implemented
	FakeAltZP = false          // Aux memory isn't implemented
	FakePage2 = false          // Aux memory isn't implemented
	Col80 = false              // Aux memory isn't implemented
	Page2 = false
	ApplyMemoryConfiguration()
}

func WipeRAM() {
	for i := 0; i < 0x10000; i++ {
		PhysicalMemory.MainMemory[i] = 0
	}
}

// SetMemoryMode is used to set UpperRamReadOnly, UpperReadMappedToROM and D000Bank number
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

// ReadMemory reads the ROM or RAM page table
func ReadMemory(address uint16) uint8 {
	if (address >= 0xc000) && (address < 0xc100) {
		return ReadIO(address)
	}

	if FakePage2 && (address >= 0x400 && address < 0x800) {
		// Return nothingness
		return uint8(0x00)
	}

	if FakeAuxMemoryRead {
		if address >= 0x200 {
			// Return nothingness
			return uint8(0x00)
		} else {
			if FakeAltZP {
				return uint8(0x00)
			}
		}
	}

	// Implicit else, we're reading the main non-IO RAM
	return ReadPageTable[address>>8][address&0xff]
}

// ReadMemory writes to the ROM or RAM page table
func WriteMemory(address uint16, value uint8) {
	if (address >= 0xc000) && (address < 0xc100) {
		WriteIO(address, value)
		return
	}

	// Magic routine to trigger an interrupt, used in the CPU interrupt tests
	if system.RunningInterruptTests && address == 0xbffc {
		oldValue := ReadMemory(address)
		system.WriteInterruptTestOpenCollector(address, oldValue, value)
		WritePageTable[uint8(address>>8)][uint8(address&0xff)] = value
		return
	}

	if FakePage2 && (address >= 0x400 && address < 0x800) {
		// Do nothing
		return
	}

	if FakeAuxMemoryWrite {
		// If there is no aux memory, then the write is ignored.
		return
	}

	memory := WritePageTable[address>>8]

	// If memory is nil, then it's read only. The write is ignored.
	if memory != nil {
		memory[uint8(address&0xff)] = value
	}

	// If doing CPU functional tests, 0x200 has the test number in it. A write to
	// it means a test passed or the tests are complete.
	if system.RunningFunctionalTests && address == 0x200 {
		testNumber := ReadMemory(0x200)
		if testNumber == 0xf0 {
			fmt.Println("Opcode testing completed")
		} else {
			fmt.Printf("Test %d OK\n", ReadMemory(0x200))
		}
	}
}

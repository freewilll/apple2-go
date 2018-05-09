package mmu

import (
	"fmt"
	"io/ioutil"
	"mos6502go/system"
)

const RomPath = "apple2e.rom"
const StackPage = 1

var PhysicalMemory struct {
	MainMemory [0xc000]uint8
	UpperROM   [0x3000]uint8
	RomC1      [0x1000]uint8
	RomC2      [0x1000]uint8
}

var PageTable [0x100][]uint8

var UsingExternalSlotRom bool

func MapFirstHalfOfIO() {
	UsingExternalSlotRom = false

	for i := 0x1; i < 0x10; i++ {
		PageTable[i+0xc0] = PhysicalMemory.RomC1[i*0x100 : i*0x100+0x100]
	}
}

func MapSecondHalfOfIO() {
	UsingExternalSlotRom = true

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

	UsingExternalSlotRom = true

	return
}

func ReadMemory(address uint16) uint8 {
	if (address >= 0xc000) && (address < 0xc100) {
		return ReadIO(address)
	} else {
		return PageTable[address>>8][address&0xff]
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
		PageTable[uint8(address>>8)][uint8(address&0xff)] = value
		return
	}

	PageTable[uint8(address>>8)][uint8(address&0xff)] = value

	if system.RunningFunctionalTests && address == 0x200 {
		testNumber := ReadMemory(0x200)
		if testNumber == 0xf0 {
			fmt.Println("Opcode testing completed")
		} else {
			fmt.Printf("Test %d OK\n", ReadMemory(0x200))
		}
	}
}

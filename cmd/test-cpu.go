package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"mos6502go/cpu"
	"mos6502go/mmu"
	"mos6502go/utils"
)

func main() {
	showInstructions := flag.Bool("show-instructions", false, "Show instructions code while running")
	skipTest0 := flag.Bool("skip-functional-test", false, "Skip functional test")
	skipTest1 := flag.Bool("skip-interrupt-test", false, "Skip interrupt test")
	breakAddressString := flag.String("break", "", "Break on address")
	flag.Parse()

	cpu.InitDisasm()
	memory := mmu.InitRAM()

	var Roms = []string{
		"6502_functional_test.bin.gz",
		"6502_interrupt_test.bin.gz",
	}

	for i, rom := range Roms {
		if (i == 0) && *skipTest0 {
			continue
		}

		if (i == 1) && *skipTest1 {
			continue
		}

		fmt.Printf("Running %s\n", rom)

		var s cpu.State
		s.Init()
		s.PC = 0x800
		cpu.RunningTests = true

		if i == 0 {
			cpu.RunningFunctionalTests = true
		}

		if i == 1 {
			cpu.RunningInterruptTests = true
		}

		bytes, err := utils.ReadMemoryFromGzipFile(rom)
		if err != nil {
			panic(err)
		}

		// Copy main RAM area 0x0000-0xbfff
		for i := 0; i < 0xc000; i++ {
			memory.PhysicalMemory.MainMemory[i] = bytes[i]
		}

		// Map writable RAM area in 0xc000-0xffff
		var RomPretendingToBeRAM [0x4000]uint8
		for i := 0x0; i < 0x4000; i++ {
			RomPretendingToBeRAM[i] = bytes[0xc000+i]
		}
		for i := 0x0; i < 0x40; i++ {
			memory.PageTable[0xc0+i] = RomPretendingToBeRAM[i*0x100 : i*0x100+0x100]
		}

		s.Memory = memory
		s.PageTable = &memory.PageTable

		var breakAddress *uint16
		if *breakAddressString != "" {
			breakAddressValue, err := hex.DecodeString(*breakAddressString)
			if err != nil {
				panic(err)
			}

			var foo uint16
			if len(breakAddressValue) == 1 {
				foo = uint16(breakAddressValue[0])
			} else if len(breakAddressValue) == 2 {
				foo = uint16(breakAddressValue[0])*uint16(0x100) + uint16(breakAddressValue[1])
			} else {
				panic("Invalid break address")
			}
			breakAddress = &foo
		}

		cpu.Run(&s, *showInstructions, breakAddress, false, 0)
		fmt.Printf("Finished running %s\n\n", rom)
	}
}

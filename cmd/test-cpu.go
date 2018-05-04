package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"mos6502go/cpu"
	"mos6502go/utils"
)

func main() {
	cpu.InitDisasm()

	showInstructions := flag.Bool("show-instructions", false, "Show instructions code while running")
	skipTest0 := flag.Bool("skip-functional-test", false, "Skip functional test")
	skipTest1 := flag.Bool("skip-interrupt-test", false, "Skip interrupt test")
	breakAddressString := flag.String("break", "", "Break on address")
	flag.Parse()

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
		cpu.RunningTests = true

		if i == 0 {
			cpu.RunningFunctionalTests = true
		}

		if i == 1 {
			cpu.RunningInterruptTests = true
		}

		bytes, err := utils.ReadMemoryFromFile(rom)
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(bytes); i++ {
			s.Memory[i] = bytes[i]
		}

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

		cpu.Run(&s, *showInstructions, breakAddress)
		fmt.Printf("Finished running %s\n\n", rom)
	}
}

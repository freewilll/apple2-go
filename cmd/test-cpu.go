package main

import (
	"encoding/hex"
	"flag"
	"mos6502go/cpu"
	"mos6502go/utils"
)

func main() {
	cpu.InitDisasm()
	var s cpu.State
	s.Init()

	bytes, err := utils.ReadMemoryFromFile("6502_functional_test.bin.gz")
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(bytes); i++ {
		s.Memory[i] = bytes[i]
	}

	showInstructions := flag.Bool("show-instructions", false, "Show instructions code while running")
	breakAddressString := flag.String("break", "", "Break on address")
	flag.Parse()

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
}

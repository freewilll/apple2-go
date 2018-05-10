package main

import (
	"flag"

	"mos6502go/cpu"
	"mos6502go/mmu"
	"mos6502go/utils"
)

func main() {
	startString := flag.String("start", "", "Start address")
	endString := flag.String("end", "", "End address")
	flag.Parse()

	startAddress := utils.DecodeCmdLineAddress(startString)
	endAddress := utils.DecodeCmdLineAddress(endString)

	if startAddress == nil {
		panic("Must include -start")
	}

	if endAddress == nil {
		e := uint16(0xffff)
		endAddress = &e
	}

	cpu.InitInstructionDecoder()
	mmu.InitApple2eROM()

	cpu.State.PC = *startAddress
	for cpu.State.PC <= *endAddress {
		cpu.PrintInstruction(false)
		cpu.AdvanceInstruction()
	}
}

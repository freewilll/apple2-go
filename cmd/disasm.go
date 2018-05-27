package main

import (
	"flag"

	"apple2/cpu"
	"apple2/mmu"
	"apple2/utils"
)

func main() {
	startString := flag.String("start", "", "Start address")
	endString := flag.String("end", "", "End address")
	flag.Parse()

	start := utils.DecodeCmdLineAddress(startString)
	end := utils.DecodeCmdLineAddress(endString)

	if start == nil {
		panic("Must include -start")
	}

	if end == nil {
		e := uint16(0xffff)
		end = &e
	}

	cpu.InitInstructionDecoder()
	mmu.InitApple2eROM()
	utils.Disassemble(*start, *end)
}

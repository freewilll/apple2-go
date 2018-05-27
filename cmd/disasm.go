package main

// Command line tool to disassemble a range of addresses

import (
	"flag"
	"fmt"
	"os"

	"github.com/freewilll/apple2/cpu"
	"github.com/freewilll/apple2/mmu"
	"github.com/freewilll/apple2/utils"
)

func main() {
	var Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Synopsis:\n    %s -start ADDRESS [-end ADDRESS]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Example:\n    %s -start d000 -end ffff\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Usage = Usage

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

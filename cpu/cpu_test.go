package cpu_test

// Test the CPU using the functional and interrupt tests defined in the *.a65
// files and compiled to bin.gz files. The cpu package is aware of tests being run and
// will exit or bail on success and failure certain conditions.

import (
	"flag"
	"fmt"
	"testing"

	"github.com/freewilll/apple2-go/cpu"
	"github.com/freewilll/apple2-go/mmu"
	"github.com/freewilll/apple2-go/system"
	"github.com/freewilll/apple2-go/utils"
)

func TestCPU(t *testing.T) {
	showInstructions := flag.Bool("show-instructions", false, "Show instructions code while running")
	skipTest0 := flag.Bool("skip-functional-test", false, "Skip functional test")
	skipTest1 := flag.Bool("skip-interrupt-test", false, "Skip interrupt test")
	breakAddressString := flag.String("break", "", "Break on address")
	flag.Parse()

	breakAddress := utils.DecodeCmdLineAddress(breakAddressString)

	cpu.InitInstructionDecoder()

	mmu.InitRAM()

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

		cpu.Init()
		cpu.State.PC = 0x800
		system.RunningTests = true

		if i == 0 {
			system.RunningFunctionalTests = true
		}

		if i == 1 {
			system.RunningInterruptTests = true
		}

		bytes, err := utils.ReadMemoryFromGzipFile(rom)
		if err != nil {
			panic(err)
		}

		// Copy main RAM area 0x0000-0xbfff
		for i := 0; i < 0xc000; i++ {
			mmu.PhysicalMemory.MainMemory[i] = bytes[i]
		}

		// Map writable RAM area in 0xc000-0xffff
		var RomPretendingToBeRAM [0x4000]uint8
		for i := 0x0; i < 0x4000; i++ {
			RomPretendingToBeRAM[i] = bytes[0xc000+i]
		}
		for i := 0x0; i < 0x40; i++ {
			mmu.ReadPageTable[0xc0+i] = RomPretendingToBeRAM[i*0x100 : i*0x100+0x100]
			mmu.WritePageTable[0xc0+i] = RomPretendingToBeRAM[i*0x100 : i*0x100+0x100]
		}

		cpu.Run(*showInstructions, breakAddress, true, false, 0)
		fmt.Printf("Finished running %s\n\n", rom)
	}
}

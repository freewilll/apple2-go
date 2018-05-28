package utils

import (
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/freewilll/apple2/cpu"
	"github.com/freewilll/apple2/system"
)

// ReadMemoryFromGzipFile just reads and uncompresses a gzip file
func ReadMemoryFromGzipFile(filename string) (data []byte, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err = ioutil.ReadAll(reader)

	return
}

// DecodeCmdLineAddress decodes a 4 byte string hex address
func DecodeCmdLineAddress(s *string) (result *uint16) {
	if *s != "" {
		breakAddressValue, err := hex.DecodeString(*s)
		if err != nil {
			panic(err)
		}

		var value uint16
		if len(breakAddressValue) == 1 {
			value = uint16(breakAddressValue[0])
		} else if len(breakAddressValue) == 2 {
			value = uint16(breakAddressValue[0])*uint16(0x100) + uint16(breakAddressValue[1])
		} else {
			panic("Invalid break address")
		}
		result = &value
	}

	return result
}

// RunUntilBreakPoint runs the CPU until it either hits a breakpoint or a time
// has expired. An assertion is done at the end to ensure the breakpoint has
// been reached.
func RunUntilBreakPoint(t *testing.T, breakAddress uint16, seconds int, showInstructions bool, message string) {
	fmt.Printf("Running until %#04x: %s \n", breakAddress, message)
	system.LastAudioCycles = 0
	exitAtBreak := false
	disableFirmwareWait := false
	cpu.Run(showInstructions, &breakAddress, exitAtBreak, disableFirmwareWait, uint64(system.CPUFrequency*seconds))
	if cpu.State.PC != breakAddress {
		t.Fatalf("Did not reach breakpoint at %04x. Got to %04x", breakAddress, cpu.State.PC)
	}
}

// Disassemble disassembles and prints the code in memory between start and end
func Disassemble(start uint16, end uint16) {
	oldPC := cpu.State.PC

	cpu.State.PC = start
	for cpu.State.PC <= end {
		cpu.PrintInstruction(false)
		cpu.AdvanceInstruction()
	}

	cpu.State.PC = oldPC
}

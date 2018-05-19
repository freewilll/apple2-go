package utils

import (
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"mos6502go/cpu"
	"mos6502go/system"
	"os"
	"testing"
)

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

func RunUntilBreakPoint(t *testing.T, breakAddress uint16, seconds int, showInstructions bool, message string) {
	fmt.Printf("Running until %#04x: %s \n", breakAddress, message)
	system.FrameCycles = 0
	system.LastAudioCycles = 0
	exitAtBreak := false
	disableFirmwareWait := false
	cpu.Run(showInstructions, &breakAddress, exitAtBreak, disableFirmwareWait, uint64(system.CpuFrequency*seconds))
	if cpu.State.PC != breakAddress {
		t.Fatalf("Did not reach breakpoint at %04x. Got to %04x", breakAddress, cpu.State.PC)
	}
}

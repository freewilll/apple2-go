package test_cpu

import (
	"mos6502go/cpu"
	"mos6502go/utils"
	"testing"
)

func TestFunctionalTests(*testing.T) {
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

	cpu.Run(&s, false, nil)
}

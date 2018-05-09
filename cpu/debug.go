package cpu

import (
	"fmt"
	"strings"
)

func printFlag(p byte, flag uint8, code string) {
	if (p & flag) == 0 {
		fmt.Print(code)
	} else {
		fmt.Printf("%s", strings.ToUpper(code))
	}
}

func printInstruction(s *State, instruction string) {
	fmt.Printf("%04x   %-24s A=%02x X=%02x Y=%02x S=%02x P=%02x ",
		s.PC,
		instruction,
		s.A,
		s.X,
		s.Y,
		s.SP,
		s.P,
	)

	printFlag(s.P, CpuFlagN, "n")
	printFlag(s.P, CpuFlagV, "v")
	fmt.Print("-") // CpuFlagR flag that's always 1
	printFlag(s.P, CpuFlagB, "b")
	printFlag(s.P, CpuFlagD, "d")
	printFlag(s.P, CpuFlagI, "i")
	printFlag(s.P, CpuFlagZ, "z")
	printFlag(s.P, CpuFlagC, "c")

	fmt.Println("")
}

func PrintInstruction(s *State) {
	opcodeValue := s.PageTable[(s.PC)>>8][(s.PC)&0xff]
	opcode := OpCodes[opcodeValue]
	mnemonic := opcode.Mnemonic
	size := opcode.AddressingMode.OperandSize
	stringFormat := opcode.AddressingMode.StringFormat

	var value uint16
	if size == 0 {
		printInstruction(s, fmt.Sprintf("%02x        %s", opcodeValue, mnemonic))
		return
	}

	var opcodes string
	var suffix string

	if opcode.AddressingMode.Mode == AmRelative {
		value = uint16(s.PageTable[(s.PC+1)>>8][(s.PC+1)&0xff])
		var relativeAddress uint16
		if (value & 0x80) == 0 {
			relativeAddress = s.PC + 2 + uint16(value)
		} else {
			relativeAddress = s.PC + 2 + uint16(value) - 0x100
		}

		suffix = fmt.Sprintf(stringFormat, relativeAddress)
		opcodes = fmt.Sprintf("%02x %02x    ", opcodeValue, value)
	} else if size == 1 {
		value = uint16(s.PageTable[(s.PC+1)>>8][(s.PC+1)&0xff])
		suffix = fmt.Sprintf(stringFormat, value)
		opcodes = fmt.Sprintf("%02x %02x    ", opcodeValue, value)
	} else if size == 2 {
		lsb := s.PageTable[(s.PC+1)>>8][(s.PC+1)&0xff]
		msb := s.PageTable[(s.PC+2)>>8][(s.PC+2)&0xff]
		value = uint16(lsb) + uint16(msb)*0x100
		suffix = fmt.Sprintf(stringFormat, value)
		opcodes = fmt.Sprintf("%02x %02x %02x ", opcodeValue, lsb, msb)
	}

	printInstruction(s, fmt.Sprintf("%s %s %s", opcodes, mnemonic, suffix))
}

func DumpMemory(s *State, offset uint16) {
	var i uint16
	for i = 0; i < 0x100; i++ {
		if (i & 0xf) == 8 {
			fmt.Print(" ")
		}
		if (i & 0xf) == 0 {
			if i > 0 {
				fmt.Print("\n")
			}
			fmt.Printf("%04x  ", offset+i)
		}
		fmt.Printf(" %02x", s.PageTable[(offset+i)>>8][(offset+i)&0xff])
	}
	fmt.Print("\n")
}

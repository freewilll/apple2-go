package cpu

import (
	"fmt"
	"mos6502go/mmu"
	"strings"
)

func printFlag(p byte, flag uint8, code string) {
	if (p & flag) == 0 {
		fmt.Print(code)
	} else {
		fmt.Printf("%s", strings.ToUpper(code))
	}
}

func printInstruction(instruction string, showRegisters bool) {
	fmt.Printf("%04x-   %-24s", State.PC, instruction)

	if showRegisters {
		fmt.Printf("     A=%02x X=%02x Y=%02x S=%02x P=%02x ",
			State.A,
			State.X,
			State.Y,
			State.SP,
			State.P,
		)

		printFlag(State.P, CpuFlagN, "n")
		printFlag(State.P, CpuFlagV, "v")
		fmt.Print("-") // CpuFlagR flag that's always 1
		printFlag(State.P, CpuFlagB, "b")
		printFlag(State.P, CpuFlagD, "d")
		printFlag(State.P, CpuFlagI, "i")
		printFlag(State.P, CpuFlagZ, "z")
		printFlag(State.P, CpuFlagC, "c")
	}

	fmt.Println("")
}

func PrintInstruction(showRegisters bool) {
	opcodeValue := mmu.ReadPageTable[(State.PC)>>8][(State.PC)&0xff]
	opcode := OpCodes[opcodeValue]
	mnemonic := opcode.Mnemonic
	size := opcode.AddressingMode.OperandSize
	stringFormat := opcode.AddressingMode.StringFormat

	var value uint16
	if size == 0 {
		printInstruction(fmt.Sprintf("%02x           %s", opcodeValue, mnemonic), showRegisters)
		return
	}

	var opcodes string
	var suffix string

	if opcode.AddressingMode.Mode == AmRelative {
		value = uint16(mmu.ReadPageTable[(State.PC+1)>>8][(State.PC+1)&0xff])
		var relativeAddress uint16
		if (value & 0x80) == 0 {
			relativeAddress = State.PC + 2 + uint16(value)
		} else {
			relativeAddress = State.PC + 2 + uint16(value) - 0x100
		}

		suffix = fmt.Sprintf(stringFormat, relativeAddress)
		opcodes = fmt.Sprintf("%02x %02x       ", opcodeValue, value)
	} else if size == 1 {
		value = uint16(mmu.ReadPageTable[(State.PC+1)>>8][(State.PC+1)&0xff])
		suffix = fmt.Sprintf(stringFormat, value)
		opcodes = fmt.Sprintf("%02x %02x       ", opcodeValue, value)
	} else if size == 2 {
		lsb := mmu.ReadPageTable[(State.PC+1)>>8][(State.PC+1)&0xff]
		msb := mmu.ReadPageTable[(State.PC+2)>>8][(State.PC+2)&0xff]
		value = uint16(lsb) + uint16(msb)*0x100
		suffix = fmt.Sprintf(stringFormat, value)
		opcodes = fmt.Sprintf("%02x %02x %02x    ", opcodeValue, lsb, msb)
	}

	printInstruction(fmt.Sprintf("%s %s %s", opcodes, mnemonic, suffix), showRegisters)
}

func AdvanceInstruction() {
	opcodeValue := mmu.ReadPageTable[(State.PC)>>8][(State.PC)&0xff]
	opcode := OpCodes[opcodeValue]
	size := opcode.AddressingMode.OperandSize + 1
	State.PC += uint16(size)
}

func DumpMemory(offset uint16) {
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
		fmt.Printf(" %02x", mmu.ReadPageTable[(offset+i)>>8][(offset+i)&0xff])
	}
	fmt.Print("\n")
}

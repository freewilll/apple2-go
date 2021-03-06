package cpu

import (
	"fmt"
	"strings"

	"github.com/freewilll/apple2-go/mmu"
)

// printFlag prints a lower or uppercase letter depending on the state of the flag
func printFlag(p byte, flag uint8, code string) {
	if (p & flag) == 0 {
		fmt.Print(code)
	} else {
		fmt.Printf("%s", strings.ToUpper(code))
	}
}

// printInstruction prings a single instruction and optionally also registers
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

		printFlag(State.P, cpuFlagN, "n")
		printFlag(State.P, cpuFlagV, "v")
		fmt.Print("-") // cpuFlagR flag that's always 1
		printFlag(State.P, cpuFlagB, "b")
		printFlag(State.P, cpuFlagD, "d")
		printFlag(State.P, cpuFlagI, "i")
		printFlag(State.P, cpuFlagZ, "z")
		printFlag(State.P, cpuFlagC, "c")
	}

	fmt.Println("")
}

// PrintInstruction prints the instruction at the current PC
func PrintInstruction(showRegisters bool) {
	opcodeValue := mmu.ReadPageTable[(State.PC)>>8][(State.PC)&0xff]
	opcode := opCodes[opcodeValue]
	mnemonic := opcode.mnemonic
	size := opcode.addressingMode.operandSize
	stringFormat := opcode.addressingMode.stringFormat

	var value uint16
	if size == 0 {
		printInstruction(fmt.Sprintf("%02x           %s", opcodeValue, mnemonic), showRegisters)
		return
	}

	var opcodes string
	var suffix string

	if opcode.addressingMode.mode == amRelative {
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

// AdvanceInstruction goes forward one instruction without executing anything
func AdvanceInstruction() {
	opcodeValue := mmu.ReadPageTable[(State.PC)>>8][(State.PC)&0xff]
	opcode := opCodes[opcodeValue]
	size := opcode.addressingMode.operandSize + 1
	State.PC += uint16(size)
}

// DumpMemory dumps $100 bytes of memory
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

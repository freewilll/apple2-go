package cpu

import (
	"fmt"
	"mos6502go/mmu"
	"os"
)

const (
	CpuFlagC byte = 1 << iota // 0x01
	CpuFlagZ                  // 0x02
	CpuFlagI                  // 0x04
	CpuFlagD                  // 0x08
	CpuFlagB                  // 0x10
	CpuFlagR                  // 0x20
	CpuFlagV                  // 0x40
	CpuFlagN                  // 0x80
)

var (
	RunningTests           bool
	RunningFunctionalTests bool
	RunningInterruptTests  bool
)

type State struct {
	Memory           *mmu.Memory
	MemoryMap        *mmu.MemoryMap // For easy access, this is a shortcut for Memory.MemoryMap
	pendingInterrupt bool
	pendingNMI       bool
	A                uint8
	X                uint8
	Y                uint8
	PC               uint16
	SP               uint8
	P                uint8
}

func (s *State) Init() {
	RunningTests = false
	RunningFunctionalTests = false
	RunningInterruptTests = false

	s.A = 0
	s.X = 0
	s.Y = 0
	s.P = CpuFlagR | CpuFlagB | CpuFlagZ
	s.SP = 0xff
	s.pendingInterrupt = false
	s.pendingNMI = false
}

func (s *State) setC(value bool) {
	if value {
		s.P |= CpuFlagC
	} else {
		s.P &= ^CpuFlagC
	}
}

func (s *State) setV(value bool) {
	if value {
		s.P |= CpuFlagV
	} else {
		s.P &= ^CpuFlagV
	}
}

func (s *State) setN(value uint8) {
	if (value & 0x80) != 0 {
		s.P |= CpuFlagN
	} else {
		s.P &= ^CpuFlagN
	}
}

func (s *State) setZ(value uint8) {
	if value == 0 {
		s.P |= CpuFlagZ
	} else {
		s.P &= ^CpuFlagZ
	}
}

func (s *State) isC() bool {
	return (s.P & CpuFlagC) != 0
}

func (s *State) isZ() bool {
	return (s.P & CpuFlagZ) != 0
}

func (s *State) isD() bool {
	return (s.P & CpuFlagD) != 0
}

func (s *State) isV() bool {
	return (s.P & CpuFlagV) != 0
}

func (s *State) isN() bool {
	return (s.P & CpuFlagN) != 0
}

func push8(s *State, value uint8) {
	(*s.MemoryMap)[mmu.StackPage][s.SP] = value
	s.SP -= 1
	s.SP &= 0xff
}

func push16(s *State, value uint16) {
	(*s.MemoryMap)[mmu.StackPage][s.SP] = uint8(value >> 8)
	(*s.MemoryMap)[mmu.StackPage][s.SP-1] = uint8(value & 0xff)
	s.SP -= 2
	s.SP &= 0xff
}

func pop8(s *State) uint8 {
	s.SP += 1
	s.SP &= 0xff
	return (*s.MemoryMap)[mmu.StackPage][s.SP]
}

func pop16(s *State) uint16 {
	s.SP += 2
	s.SP &= 0xff
	msb := uint16((*s.MemoryMap)[mmu.StackPage][s.SP])
	lsb := uint16((*s.MemoryMap)[mmu.StackPage][s.SP-1])
	return lsb + msb<<8
}

func readMemory(s *State, address uint16) uint8 {
	if (address >= 0xc000) && (address < 0xc100) {
		fmt.Printf("TODO read %04x\n", address)
		return 0
	}

	return (*s.MemoryMap)[uint8(address>>8)][uint8(address&0xff)]
}

// Handle a write to a magic test address that triggers an interrupt and/or an NMI
func writeInterruptTestOpenCollector(s *State, address uint16, value uint8) {
	oldValue := readMemory(s, address)

	oldInterrupt := (oldValue & 0x1) == 0x1
	oldNMI := (oldValue & 0x2) == 0x2

	interrupt := (value & 0x1) == 0x1
	NMI := (value & 0x2) == 0x2

	if oldInterrupt != interrupt {
		s.pendingInterrupt = interrupt
	}

	if oldNMI != NMI {
		s.pendingNMI = NMI
	}

	(*s.MemoryMap)[uint8(address>>8)][uint8(address&0xff)] = value
}

func writeMemory(s *State, address uint16, value uint8) {
	if RunningInterruptTests && address == 0xbffc {
		writeInterruptTestOpenCollector(s, address, value)
		return
	}

	if address >= 0xc000 {
		if address == mmu.CLRCXROM {
			mmu.MapFirstHalfOfIO(s.Memory)
		} else if address == mmu.SETCXROM {
			mmu.MapSecondHalfOfIO(s.Memory)
		} else {
			fmt.Printf("TODO write %04x\n", address)
		}
		return
	}

	if address >= 0x400 && address < 0x800 {
		fmt.Printf("Text page write %04x: %02x\n", address, value)
	}

	(*s.MemoryMap)[uint8(address>>8)][uint8(address&0xff)] = value

	if RunningFunctionalTests && address == 0x200 {
		testNumber := readMemory(s, 0x200)
		if testNumber == 0xf0 {
			fmt.Println("Opcode testing completed")
		} else {
			fmt.Printf("Test %d OK\n", readMemory(s, 0x200))
		}
	}
}

func branch(s *State, cycles *int, instructionName string, doBranch bool) {
	value := readMemory(s, s.PC+1)

	var relativeAddress uint16
	if (value & 0x80) == 0 {
		relativeAddress = s.PC + uint16(value) + 2
	} else {
		relativeAddress = s.PC + uint16(value) + 2 - 0x100
	}

	*cycles += 2
	if doBranch {
		if RunningTests && s.PC == relativeAddress {
			fmt.Printf("Trap at $%04x\n", relativeAddress)
			os.Exit(0)
		}

		samePage := (s.PC & 0xff00) != (relativeAddress & 0xff00)
		if samePage {
			*cycles += 1
		} else {
			*cycles += 2
		}
		s.PC = relativeAddress
	} else {
		s.PC += 2
	}
}

func getAddressFromAddressMode(s *State, addressMode byte) (result uint16, pageBoundaryCrossed bool) {
	switch addressMode {
	case AmZeroPage:
		result = uint16(readMemory(s, s.PC+1))
	case AmZeroPageX:
		result = (uint16(readMemory(s, s.PC+1)) + uint16(s.X)) & 0xff
	case AmZeroPageY:
		result = (uint16(readMemory(s, s.PC+1)) + uint16(s.Y)) & 0xff
	case AmAbsolute:
		result = uint16(readMemory(s, s.PC+1)) + uint16(readMemory(s, s.PC+2))<<8
	case AmAbsoluteX:
		value := uint16(readMemory(s, s.PC+1)) + uint16(readMemory(s, s.PC+2))<<8
		pageBoundaryCrossed = (value & 0xff00) != ((value + uint16(s.X)) & 0xff00)
		result = value + uint16(s.X)
	case AmAbsoluteY:
		value := uint16(readMemory(s, s.PC+1)) + uint16(readMemory(s, s.PC+2))<<8
		pageBoundaryCrossed = (value & 0xff00) != ((value + uint16(s.Y)) & 0xff00)
		result = value + uint16(s.Y)
	case AmIndirectX:
		zeroPageAddress := (readMemory(s, s.PC+1) + s.X) & 0xff
		result = uint16(readMemory(s, uint16(zeroPageAddress))) + uint16(readMemory(s, uint16(zeroPageAddress)+1))<<8
	case AmIndirectY:
		address := uint16(readMemory(s, s.PC+1))
		lsb := uint16(readMemory(s, address))
		msb := uint16(readMemory(s, address+1))
		value := lsb + msb<<8
		pageBoundaryCrossed = (value & 0xff00) != ((value + uint16(s.Y)) & 0xff00)
		result = value + uint16(s.Y)
	default:
		panic(fmt.Sprintf("Unknown address mode %d in getAddressFromAddressMode()", addressMode))
	}

	return result, pageBoundaryCrossed
}

func readMemoryWithAddressMode(s *State, addressMode byte) (result uint8, pageBoundaryCrossed bool) {
	switch addressMode {
	case AmImmediate:
		result = readMemory(s, s.PC+1)
		s.PC += 2
	case AmZeroPage:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 2
	case AmZeroPageX:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 2
	case AmZeroPageY:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 2
	case AmAbsolute:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 3
	case AmAbsoluteX:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 3
	case AmAbsoluteY:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 3
	case AmIndirectX:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 2
	case AmIndirectY:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(s, addressMode)
		result = readMemory(s, address)
		s.PC += 2
	default:
		result = 0
		s.PC++
	}

	return result, pageBoundaryCrossed
}

// STA, STX and STY
func store(s *State, cycles *int, regValue uint8, addressMode byte) {
	address, _ := getAddressFromAddressMode(s, addressMode)
	writeMemory(s, address, regValue)

	switch addressMode {
	case AmZeroPage:
		s.PC += 2
		*cycles += 3
	case AmZeroPageX:
		s.PC += 2
		*cycles += 4
	case AmZeroPageY:
		s.PC += 2
		*cycles += 4
	case AmAbsolute:
		s.PC += 3
		*cycles += 4
	case AmAbsoluteX:
		s.PC += 3
		*cycles += 5
	case AmAbsoluteY:
		s.PC += 3
		*cycles += 5
	case AmIndirect:
		s.PC += 2
		*cycles += 6
	case AmIndirectX:
		s.PC += 2
		*cycles += 6
	case AmIndirectY:
		s.PC += 2
		*cycles += 6
	default:
		panic(fmt.Sprintf("Unknown address mode %d in store()", addressMode))
	}
}

// These instructions take the same amount of cycles
func advanceCyclesForAcculumatorOperation(cycles *int, addressMode byte, pageBoundaryCrossed bool) {
	extraCycle := 0
	if pageBoundaryCrossed {
		extraCycle = 1
	}

	switch addressMode {
	case AmImmediate:
		*cycles += 2
	case AmZeroPage:
		*cycles += 3
	case AmZeroPageX:
		*cycles += 4
	case AmZeroPageY:
		*cycles += 4
	case AmAbsolute:
		*cycles += 4
	case AmAbsoluteX:
		*cycles += 4 + extraCycle
	case AmAbsoluteY:
		*cycles += 4 + extraCycle
	case AmIndirectX:
		*cycles += 6
	case AmIndirectY:
		*cycles += 5 + extraCycle
	default:
		panic(fmt.Sprintf("Unknown address mode %d in advanceCyclesForAcculumatorOperation()", addressMode))
	}
}

func load(s *State, cycles *int, addressMode byte) uint8 {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(s, addressMode)
	s.setN(value)
	s.setZ(value)
	advanceCyclesForAcculumatorOperation(cycles, addressMode, pageBoundaryCrossed)
	return value
}

func cmp(s *State, cycles *int, regValue uint8, addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(s, addressMode)
	var result uint16
	result = uint16(regValue) - uint16(value)
	s.setC(result < 0x100)
	s.setN(uint8(result))
	s.setZ(uint8(result & 0xff))
	advanceCyclesForAcculumatorOperation(cycles, addressMode, pageBoundaryCrossed)
}

func ora(s *State, cycles *int, addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(s, addressMode)
	s.A |= value
	s.setN(s.A)
	s.setZ(s.A)
	advanceCyclesForAcculumatorOperation(cycles, addressMode, pageBoundaryCrossed)
}

func and(s *State, cycles *int, addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(s, addressMode)
	s.A &= value
	s.setN(s.A)
	s.setZ(s.A)
	advanceCyclesForAcculumatorOperation(cycles, addressMode, pageBoundaryCrossed)
}

func eor(s *State, cycles *int, addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(s, addressMode)
	s.A ^= value
	s.setN(s.A)
	s.setZ(s.A)
	advanceCyclesForAcculumatorOperation(cycles, addressMode, pageBoundaryCrossed)
}

func adc(s *State, cycles *int, addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(s, addressMode)

	var temp uint16
	temp = uint16(s.A) + uint16(value)

	var carry uint8
	if s.isC() {
		carry = 1
	}

	if carry > 0 {
		temp++
	}

	// This is not valid in decimal mode
	s.setZ(uint8(temp & 0xff))

	if s.isD() {
		if ((s.A & 0xf) + (value & 0xf) + carry) > 9 {
			temp += 6
		}

		s.setN(uint8(temp))
		s.setV((((s.A ^ value) & 0x80) == 0) && (((s.A ^ uint8(temp)) & 0x80) != 0))

		if temp > 0x99 {
			temp += 96
		}
		s.setC(temp > 0x99)
	} else {
		s.setN(uint8(temp))
		s.setV((((s.A ^ value) & 0x80) == 0) && (((s.A ^ uint8(temp)) & 0x80) != 0))
		s.setC(temp > 0xff)
	}

	s.A = uint8(temp & 0xff)

	s.setN(s.A)
	s.setZ(s.A)
	advanceCyclesForAcculumatorOperation(cycles, addressMode, pageBoundaryCrossed)
}

func sbc(s *State, cycles *int, addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(s, addressMode)

	var temp uint16
	temp = uint16(s.A) - uint16(value)

	var carry uint8
	if s.isC() {
		carry = 0
	} else {
		carry = 1
	}

	if carry > 0 {
		temp--
	}

	s.setN(uint8(temp))

	// This is not valid in decimal mode
	s.setZ(uint8(temp & 0xff))

	s.setV((((s.A ^ uint8(temp)) & 0x80) != 0) && (((s.A ^ value) & 0x80) != 0))

	if s.isD() {
		if ((int8(s.A) & 0xf) - int8(carry)) < (int8(value) & 0xf) {
			temp -= 6
		}

		if temp > 0x99 {
			temp -= 96
		}
	}

	s.setC(temp < 0x100)
	s.A = uint8(temp & 0xff)
	advanceCyclesForAcculumatorOperation(cycles, addressMode, pageBoundaryCrossed)
}

func bit(s *State, address uint16) {
	value := readMemory(s, address)
	s.setN(value)
	s.setV((value & 0x40) != 0)
	s.setZ(value & s.A)
}

// Read the address/value for an ASL, LSR, ROR, ROL
func preProcessShift(s *State, cycles *int, addressMode byte) (address uint16, value uint8) {
	if addressMode == AmAccumulator {
		value = s.A
	} else {
		address, _ = getAddressFromAddressMode(s, addressMode)
		value = readMemory(s, address)
	}

	if addressMode == AmAccumulator {
		value = s.A
	} else {
		address, _ = getAddressFromAddressMode(s, addressMode)
		value = readMemory(s, address)
	}

	return
}

// Store the result of a ASL, LSR, ROR, ROL and advance PC and cycles
func postProcessShift(s *State, cycles *int, addressMode byte, address uint16, value uint8) {
	switch addressMode {
	case AmAccumulator:
		s.A = value
		s.PC += 1
		*cycles += 2
	case AmZeroPage:
		writeMemory(s, address, value)
		s.PC += 2
		*cycles += 5
	case AmZeroPageX:
		writeMemory(s, address, value)
		s.PC += 2
		*cycles += 6
	case AmAbsolute:
		writeMemory(s, address, value)
		s.PC += 3
		*cycles += 6
	case AmAbsoluteX:
		writeMemory(s, address, value)
		s.PC += 3
		*cycles += 7
	default:
		panic(fmt.Sprintf("Unknown address mode %d in postProcessShift()", addressMode))
	}
}

func postProcessIncDec(s *State, cycles *int, addressMode byte) {
	switch addressMode {
	case AmZeroPage:
		s.PC += 2
		*cycles += 5
	case AmZeroPageX:
		s.PC += 2
		*cycles += 6
	case AmAbsolute:
		s.PC += 3
		*cycles += 6
	case AmAbsoluteX:
		s.PC += 3
		*cycles += 7
	default:
		panic(fmt.Sprintf("Unknown address mode %d in INC", addressMode))
	}
}

func brk(s *State, cycles *int) {
	push16(s, s.PC+2)
	s.P |= CpuFlagB
	push8(s, s.P)
	s.P |= CpuFlagI
	s.PC = uint16(readMemory(s, 0xffff))<<8 + uint16(readMemory(s, 0xfffe))
	*cycles += 7
}

func irq(s *State, cycles *int) {
	push16(s, s.PC)
	s.P &= ^CpuFlagB
	push8(s, s.P)
	s.P |= CpuFlagI
	s.PC = uint16(readMemory(s, 0xffff))<<8 + uint16(readMemory(s, 0xfffe))
	*cycles += 7
}

func nmi(s *State, cycles *int) {
	push16(s, s.PC)
	s.P &= ^CpuFlagB
	push8(s, s.P)
	s.P |= CpuFlagI
	s.PC = uint16(readMemory(s, 0xfffb))<<8 + uint16(readMemory(s, 0xfffa))
	*cycles += 7
}

func Run(s *State, showInstructions bool, breakAddress *uint16, disableBell bool, wantedCycles int) {
	cycles := 0

	for {
		if (wantedCycles != 0) && (cycles >= wantedCycles) {
			return
		}

		if RunningTests && (s.PC == 0x3869) {
			fmt.Println("Functional tests passed")
			return
		}

		if RunningTests && (s.PC == 0x0af5) {
			fmt.Println("Interrupt tests passed")
			return
		}

		if s.pendingInterrupt && ((s.P & CpuFlagI) == 0) {
			irq(s, &cycles)
			s.pendingInterrupt = false
			continue
		}

		if s.pendingNMI {
			nmi(s, &cycles)
			s.pendingNMI = false
			continue
		}

		if showInstructions {
			PrintInstruction(s)
		}

		opcode := readMemory(s, s.PC)
		addressMode := OpCodes[opcode].AddressingMode.Mode

		if breakAddress != nil && s.PC == *breakAddress {
			fmt.Printf("Break at $%04x\n", *breakAddress)
			os.Exit(0)
		}

		switch opcode {

		case 0x4c: // JMP $0000
			value := uint16(readMemory(s, s.PC+1)) + uint16(readMemory(s, s.PC+2))<<8
			if RunningTests && s.PC == value {
				fmt.Printf("Trap at $%04x\n", value)
				os.Exit(0)
			}
			s.PC = value
			cycles += 3
		case 0x6c: // JMP ($0000)
			value := uint16(readMemory(s, s.PC+1)) + uint16(readMemory(s, s.PC+2))<<8
			s.PC = uint16(readMemory(s, value)) + uint16(readMemory(s, value+1))<<8
			cycles += 5

		case 0x20: // JSR $0000
			value := uint16(readMemory(s, s.PC+1)) + uint16(readMemory(s, s.PC+2))<<8
			cycles += 6

			if disableBell && value == 0xfca8 {
				s.PC += 3
				continue
			}

			push16(s, s.PC+2)
			s.PC = value

		case 0x60: // RTS
			value := pop16(s)
			s.PC = value + 1
			cycles += 6

		case 0xa9, 0xa5, 0xb5, 0xad, 0xbd, 0xb9, 0xa1, 0xb1: // LDA
			s.A = load(s, &cycles, addressMode)
		case 0xa2, 0xa6, 0xb6, 0xae, 0xbe: // LDX
			s.X = load(s, &cycles, addressMode)
		case 0xa0, 0xa4, 0xb4, 0xac, 0xbc: // LDY
			s.Y = load(s, &cycles, addressMode)

		case 0x85, 0x95, 0x8d, 0x9d, 0x99, 0x81, 0x91: //STA
			store(s, &cycles, s.A, addressMode)
		case 0x86, 0x96, 0x8e: // STX
			store(s, &cycles, s.X, addressMode)
		case 0x84, 0x94, 0x8c: //STY
			store(s, &cycles, s.Y, addressMode)

		case 0xc9, 0xc5, 0xd5, 0xcd, 0xdd, 0xd9, 0xc1, 0xd1: // CMP
			cmp(s, &cycles, s.A, addressMode)
		case 0xe0, 0xe4, 0xeC: // CPX
			cmp(s, &cycles, s.X, addressMode)
		case 0xc0, 0xc4, 0xcc: // CPY
			cmp(s, &cycles, s.Y, addressMode)
		case 0x09, 0x05, 0x15, 0x0d, 0x1d, 0x19, 0x01, 0x11: // ORA
			ora(s, &cycles, addressMode)
		case 0x29, 0x25, 0x35, 0x2d, 0x3d, 0x39, 0x21, 0x31: // AND
			and(s, &cycles, addressMode)
		case 0x49, 0x45, 0x55, 0x4d, 0x5d, 0x59, 0x41, 0x51: // EOR
			eor(s, &cycles, addressMode)
		case 0x69, 0x65, 0x75, 0x6d, 0x7d, 0x79, 0x61, 0x71: // ADC
			adc(s, &cycles, addressMode)
		case 0xe9, 0xe5, 0xf5, 0xed, 0xfd, 0xf9, 0xe1, 0xf1: // SBC
			sbc(s, &cycles, addressMode)

		// Register transfers
		case 0xaa: // TAX
			s.X = s.A
			s.setN(s.X)
			s.setZ(s.X)
			s.PC++
			cycles += 2
		case 0xa8: // TAY
			s.Y = s.A
			s.setN(s.Y)
			s.setZ(s.Y)
			s.PC++
			cycles += 2
		case 0xba: // TSX
			s.X = s.SP
			s.setN(s.X)
			s.setZ(s.X)
			s.PC++
			cycles += 2
		case 0x8a: // TXA
			s.A = s.X
			s.setN(s.A)
			s.setZ(s.A)
			s.PC++
			cycles += 2
		case 0x9a: // TXS
			s.SP = s.X
			s.PC++
			cycles += 2
		case 0x98: // TYA
			s.A = s.Y
			s.setN(s.A)
			s.setZ(s.A)
			s.PC++
			cycles += 2

		case 0xE8:
			s.X = (s.X + 1) & 0xff
			s.setN(s.X)
			s.setZ(s.X)
			s.PC++
			cycles += 2
		case 0xC8:
			s.Y = (s.Y + 1) & 0xff
			s.setN(s.Y)
			s.setZ(s.Y)
			s.PC++
			cycles += 2
		case 0xca:
			s.X = (s.X - 1) & 0xff
			s.setN(s.X)
			s.setZ(s.X)
			s.PC++
			cycles += 2
		case 0x88:
			s.Y = (s.Y - 1) & 0xff
			s.setN(s.Y)
			s.setZ(s.Y)
			s.PC++
			cycles += 2

		// Branch instructions
		case 0x10:
			branch(s, &cycles, "BPL", !s.isN())
		case 0x30:
			branch(s, &cycles, "BMI", s.isN())
		case 0x50:
			branch(s, &cycles, "BVC", !s.isV())
		case 0x70:
			branch(s, &cycles, "BVS", s.isV())
		case 0x90:
			branch(s, &cycles, "BCC", !s.isC())
		case 0xb0:
			branch(s, &cycles, "BCS", s.isC())
		case 0xd0:
			branch(s, &cycles, "BNE", !s.isZ())
		case 0xf0:
			branch(s, &cycles, "BEQ", s.isZ())

		// Flag setting
		case 0x18:
			s.setC(false)
			s.PC++
			cycles += 2
		case 0x38:
			s.setC(true)
			s.PC++
			cycles += 2
		case 0x58:
			s.P &= ^CpuFlagI
			s.PC++
			cycles += 2
		case 0x78:
			s.P |= CpuFlagI
			s.PC++
			cycles += 2
		case 0xb8:
			s.P &= ^CpuFlagV
			s.PC++
			cycles += 2
		case 0xd8:
			s.P &= ^CpuFlagD
			s.PC++
			cycles += 2
		case 0xf8:
			s.P |= CpuFlagD
			s.PC++
			cycles += 2

		case 0x48: // PHA
			push8(s, s.A)
			s.PC++
			cycles += 3
		case 0x68: // PLA
			s.A = pop8(s)
			s.setN(s.A)
			s.setZ(s.A)
			s.PC++
			cycles += 4
		case 0x08: // PHP
			// From http://visual6502.org/wiki/index.php?title=6502_BRK_and_B_bit#the_B_flag_and_the_various_mechanisms
			// software instructions BRK & PHP will push the B flag as being 1.
			push8(s, s.P|CpuFlagB)
			s.PC++
			cycles += 3
		case 0x28: // PLP
			// CpuFlagR is always supposed to be 1
			s.P = pop8(s) | CpuFlagR
			s.PC++
			cycles += 4
		case 0xea:
			s.PC++
			cycles += 2

		case 0x00: // BRK
			brk(s, &cycles)
		case 0x40: // RTI
			s.P = pop8(s) | CpuFlagR
			value := pop16(s)
			s.PC = value
			cycles += 6

		case 0x24: // BIT $00
			address := readMemory(s, s.PC+1)
			bit(s, uint16(address))
			s.PC += 2
			cycles += 3
		case 0x2C: // BIT $0000
			address := uint16(readMemory(s, s.PC+1)) + uint16(readMemory(s, s.PC+2))<<8
			bit(s, address)
			s.PC += 3
			cycles += 4

		case 0x0a, 0x06, 0x16, 0x0e, 0x1e: // ASL
			address, value := preProcessShift(s, &cycles, addressMode)
			s.setC((value & 0x80) != 0)
			value = (value << 1) & 0xff
			s.setZ(value)
			s.setN(value)
			postProcessShift(s, &cycles, addressMode, address, value)
		case 0x4a, 0x46, 0x56, 0x4e, 0x5e: // LSR
			address, value := preProcessShift(s, &cycles, addressMode)
			s.setC((value & 0x01) != 0)
			value >>= 1
			s.setZ(value)
			s.setN(value)
			postProcessShift(s, &cycles, addressMode, address, value)
		case 0x2a, 0x26, 0x36, 0x2e, 0x3e: // ROL
			address, value := preProcessShift(s, &cycles, addressMode)
			value16 := uint16(value)
			value16 <<= 1
			if (s.P & CpuFlagC) != 0 {
				value16 |= 0x01
			}
			s.setC((value16 & 0x100) != 0)
			value = uint8(value16 & 0xff)
			s.setZ(value)
			s.setN(value)
			postProcessShift(s, &cycles, addressMode, address, value)
		case 0x6a, 0x66, 0x76, 0x6e, 0x7e: // ROR
			address, value := preProcessShift(s, &cycles, addressMode)
			value16 := uint16(value)
			if (s.P & CpuFlagC) != 0 {
				value16 |= 0x100
			}
			s.setC((value16 & 0x01) != 0)
			value = uint8(value16 >> 1)
			s.setZ(value)
			s.setN(value)
			postProcessShift(s, &cycles, addressMode, address, value)

		case 0xe6, 0xf6, 0xee, 0xfe: // INC
			address, _ := getAddressFromAddressMode(s, addressMode)
			value := readMemory(s, address)
			value = (value + 1) & 0xff
			s.setZ(value)
			s.setN(value)
			writeMemory(s, address, value)
			postProcessIncDec(s, &cycles, addressMode)

		case 0xc6, 0xd6, 0xce, 0xde: // DEC
			address, _ := getAddressFromAddressMode(s, addressMode)
			value := readMemory(s, address)
			value = (value - 1) & 0xff
			s.setZ(value)
			s.setN(value)
			writeMemory(s, address, value)
			postProcessIncDec(s, &cycles, addressMode)

		default:
			fmt.Printf("Unknown opcode $%02x\n", opcode)
			return
		}
	}
}

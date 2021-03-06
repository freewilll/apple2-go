package cpu

import (
	"fmt"
	"os"

	"github.com/freewilll/apple2-go/mmu"
	"github.com/freewilll/apple2-go/system"
)

const (
	cpuFlagC byte = 1 << iota // 0x01 carry
	cpuFlagZ                  // 0x02 zero
	cpuFlagI                  // 0x04 interrupt disable
	cpuFlagD                  // 0x08 decimal mode
	cpuFlagB                  // 0x10 break
	cpuFlagR                  // 0x20 reserved (unused)
	cpuFlagV                  // 0x40 overflow
	cpuFlagN                  // 0x80 sign/negative
)

// State contains the CPU state
var State struct {
	A  uint8  // accumulator
	X  uint8  // X register
	Y  uint8  // Y register
	PC uint16 // program counter
	SP uint8  // stack pointer
	P  uint8  // processor flags
}

// Init sets up the CPU registers, interrupts and disable testing code
func Init() {
	system.RunningTests = false
	system.RunningFunctionalTests = false
	system.RunningInterruptTests = false

	system.PendingInterrupt = false
	system.PendingNMI = false

	State.A = 0
	State.X = 0
	State.Y = 0
	State.P = cpuFlagR | cpuFlagB | cpuFlagZ
	State.SP = 0xff
}

// setC sets the carry flag
func setC(value bool) {
	if value {
		State.P |= cpuFlagC
	} else {
		State.P &= ^cpuFlagC
	}
}

// setV sets the overflow flag
func setV(value bool) {
	if value {
		State.P |= cpuFlagV
	} else {
		State.P &= ^cpuFlagV
	}
}

// setN sets the sign/negative flag if the value is negative (>=0x80)
func setN(value uint8) {
	if (value & 0x80) != 0 {
		State.P |= cpuFlagN
	} else {
		State.P &= ^cpuFlagN
	}
}

// setZ sets the zero flag if the value is zero
func setZ(value uint8) {
	if value == 0 {
		State.P |= cpuFlagZ
	} else {
		State.P &= ^cpuFlagZ
	}
}

func isC() bool {
	return (State.P & cpuFlagC) != 0
}

func isZ() bool {
	return (State.P & cpuFlagZ) != 0
}

func isD() bool {
	return (State.P & cpuFlagD) != 0
}

func isV() bool {
	return (State.P & cpuFlagV) != 0
}

func isN() bool {
	return (State.P & cpuFlagN) != 0
}

// push8 pushes an 8 bit value to the stack
func push8(value uint8) {
	mmu.WritePageTable[mmu.StackPage][State.SP] = value
	State.SP--
	State.SP &= 0xff
}

// push16 pushes a 16 bit value to the stack
func push16(value uint16) {
	mmu.WritePageTable[mmu.StackPage][State.SP] = uint8(value >> 8)
	mmu.WritePageTable[mmu.StackPage][State.SP-1] = uint8(value & 0xff)
	State.SP -= 2
	State.SP &= 0xff
}

// pop8 pulls an 8 bit value from the stack
func pop8() uint8 {
	State.SP++
	State.SP &= 0xff
	return mmu.ReadPageTable[mmu.StackPage][State.SP]
}

// pop16 pulls a 16 bit value from the stack
func pop16() uint16 {
	State.SP += 2
	State.SP &= 0xff
	msb := uint16(mmu.ReadPageTable[mmu.StackPage][State.SP])
	lsb := uint16(mmu.ReadPageTable[mmu.StackPage][State.SP-1])
	return lsb + msb<<8
}

// branch handles a branch instruction
func branch(doBranch bool) {
	value := mmu.ReadMemory(State.PC + 1)

	var relativeAddress uint16
	if (value & 0x80) == 0 {
		relativeAddress = State.PC + uint16(value) + 2
	} else {
		relativeAddress = State.PC + uint16(value) + 2 - 0x100
	}

	system.FrameCycles += 2
	if doBranch {
		if system.RunningTests && State.PC == relativeAddress {
			// Catch an infinite loop and exit
			fmt.Printf("Trap at $%04x\n", relativeAddress)
			os.Exit(0)
		}

		// The number of cycles depends on if a page boundary was crossed
		samePage := (State.PC & 0xff00) == (relativeAddress & 0xff00)
		if samePage {
			system.FrameCycles++
		} else {
			system.FrameCycles += 2
		}
		State.PC = relativeAddress
	} else {
		State.PC += 2
	}
}

// getAddressFromAddressMode gets the address an instruction is referring to
func getAddressFromAddressMode(addressMode byte) (result uint16, pageBoundaryCrossed bool) {
	switch addressMode {
	case amZeroPage:
		result = uint16(mmu.ReadMemory(State.PC + 1))
	case amZeroPageX:
		result = (uint16(mmu.ReadMemory(State.PC+1)) + uint16(State.X)) & 0xff
	case amZeroPageY:
		result = (uint16(mmu.ReadMemory(State.PC+1)) + uint16(State.Y)) & 0xff
	case amAbsolute:
		result = uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
	case amAbsoluteX:
		value := uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
		pageBoundaryCrossed = (value & 0xff00) != ((value + uint16(State.X)) & 0xff00)
		result = value + uint16(State.X)
	case amAbsoluteY:
		value := uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
		pageBoundaryCrossed = (value & 0xff00) != ((value + uint16(State.Y)) & 0xff00)
		result = value + uint16(State.Y)
	case amIndirectX:
		zeroPageAddress := (mmu.ReadMemory(State.PC+1) + State.X) & 0xff
		result = uint16(mmu.ReadMemory(uint16(zeroPageAddress))) + uint16(mmu.ReadMemory(uint16(zeroPageAddress)+1))<<8
	case amIndirectY:
		address := uint16(mmu.ReadMemory(State.PC + 1))
		lsb := uint16(mmu.ReadMemory(address))
		msb := uint16(mmu.ReadMemory(address + 1))
		value := lsb + msb<<8
		pageBoundaryCrossed = (value & 0xff00) != ((value + uint16(State.Y)) & 0xff00)
		result = value + uint16(State.Y)
	default:
		panic(fmt.Sprintf("Unknown address mode %d in getAddressFromAddressMode()", addressMode))
	}

	return result, pageBoundaryCrossed
}

// readMemoryWithAddressMode reads memory using a particular address mode
func readMemoryWithAddressMode(addressMode byte) (result uint8, pageBoundaryCrossed bool) {
	switch addressMode {
	case amImmediate:
		result = mmu.ReadMemory(State.PC + 1)
		State.PC += 2
	case amZeroPage:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 2
	case amZeroPageX:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 2
	case amZeroPageY:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 2
	case amAbsolute:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 3
	case amAbsoluteX:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 3
	case amAbsoluteY:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 3
	case amIndirectX:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 2
	case amIndirectY:
		var address uint16
		address, pageBoundaryCrossed = getAddressFromAddressMode(addressMode)
		result = mmu.ReadMemory(address)
		State.PC += 2
	default:
		result = 0
		State.PC++
	}

	return result, pageBoundaryCrossed
}

// store handles STA, STX and STY
func store(regValue uint8, addressMode byte) {
	address, _ := getAddressFromAddressMode(addressMode)
	mmu.WriteMemory(address, regValue)

	switch addressMode {
	case amZeroPage:
		State.PC += 2
		system.FrameCycles += 3
	case amZeroPageX:
		State.PC += 2
		system.FrameCycles += 4
	case amZeroPageY:
		State.PC += 2
		system.FrameCycles += 4
	case amAbsolute:
		State.PC += 3
		system.FrameCycles += 4
	case amAbsoluteX:
		State.PC += 3
		system.FrameCycles += 5
	case amAbsoluteY:
		State.PC += 3
		system.FrameCycles += 5
	case amIndirect:
		State.PC += 2
		system.FrameCycles += 6
	case amIndirectX:
		State.PC += 2
		system.FrameCycles += 6
	case amIndirectY:
		State.PC += 2
		system.FrameCycles += 6
	default:
		panic(fmt.Sprintf("Unknown address mode %d in store()", addressMode))
	}
}

// advanceCyclesForAcculumatorOperation advances the number of cycles for common accumulator operations
func advanceCyclesForAcculumatorOperation(addressMode byte, pageBoundaryCrossed bool) {
	extraCycle := uint64(0)
	if pageBoundaryCrossed {
		extraCycle = 1
	}

	switch addressMode {
	case amImmediate:
		system.FrameCycles += 2
	case amZeroPage:
		system.FrameCycles += 3
	case amZeroPageX:
		system.FrameCycles += 4
	case amZeroPageY:
		system.FrameCycles += 4
	case amAbsolute:
		system.FrameCycles += 4
	case amAbsoluteX:
		system.FrameCycles += 4 + extraCycle
	case amAbsoluteY:
		system.FrameCycles += 4 + extraCycle
	case amIndirectX:
		system.FrameCycles += 6
	case amIndirectY:
		system.FrameCycles += 5 + extraCycle
	default:
		panic(fmt.Sprintf("Unknown address mode %d in advanceCyclesForAcculumatorOperation()", addressMode))
	}
}

// LDA, LDX, LDY
func load(addressMode byte) uint8 {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(addressMode)
	setN(value)
	setZ(value)
	advanceCyclesForAcculumatorOperation(addressMode, pageBoundaryCrossed)
	return value
}

// CMP, CPX, CPY
func cmp(regValue uint8, addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(addressMode)
	var result uint16
	result = uint16(regValue) - uint16(value)
	setC(result < 0x100)
	setN(uint8(result))
	setZ(uint8(result & 0xff))
	advanceCyclesForAcculumatorOperation(addressMode, pageBoundaryCrossed)
}

func ora(addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(addressMode)
	State.A |= value
	setN(State.A)
	setZ(State.A)
	advanceCyclesForAcculumatorOperation(addressMode, pageBoundaryCrossed)
}

func and(addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(addressMode)
	State.A &= value
	setN(State.A)
	setZ(State.A)
	advanceCyclesForAcculumatorOperation(addressMode, pageBoundaryCrossed)
}

func eor(addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(addressMode)
	State.A ^= value
	setN(State.A)
	setZ(State.A)
	advanceCyclesForAcculumatorOperation(addressMode, pageBoundaryCrossed)
}

func adc(addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(addressMode)

	var temp uint16
	temp = uint16(State.A) + uint16(value)

	var carry uint8
	if isC() {
		carry = 1
	}

	if carry > 0 {
		temp++
	}

	// This is not valid in decimal mode
	setZ(uint8(temp & 0xff))

	if isD() {
		if ((State.A & 0xf) + (value & 0xf) + carry) > 9 {
			temp += 6
		}

		setN(uint8(temp))
		setV((((State.A ^ value) & 0x80) == 0) && (((State.A ^ uint8(temp)) & 0x80) != 0))

		if temp > 0x99 {
			temp += 96
		}
		setC(temp > 0x99)
	} else {
		setN(uint8(temp))
		setV((((State.A ^ value) & 0x80) == 0) && (((State.A ^ uint8(temp)) & 0x80) != 0))
		setC(temp > 0xff)
	}

	State.A = uint8(temp & 0xff)

	setN(State.A)
	setZ(State.A)
	advanceCyclesForAcculumatorOperation(addressMode, pageBoundaryCrossed)
}

func sbc(addressMode byte) {
	value, pageBoundaryCrossed := readMemoryWithAddressMode(addressMode)

	var temp uint16
	temp = uint16(State.A) - uint16(value)

	var carry uint8
	if isC() {
		carry = 0
	} else {
		carry = 1
	}

	if carry > 0 {
		temp--
	}

	setN(uint8(temp))

	// This is not valid in decimal mode
	setZ(uint8(temp & 0xff))

	setV((((State.A ^ uint8(temp)) & 0x80) != 0) && (((State.A ^ value) & 0x80) != 0))

	if isD() {
		if ((int8(State.A) & 0xf) - int8(carry)) < (int8(value) & 0xf) {
			temp -= 6
		}

		if temp > 0x99 {
			temp -= 96
		}
	}

	setC(temp < 0x100)
	State.A = uint8(temp & 0xff)
	advanceCyclesForAcculumatorOperation(addressMode, pageBoundaryCrossed)
}

func bit(address uint16) {
	value := mmu.ReadMemory(address)
	setN(value)
	setV((value & 0x40) != 0)
	setZ(value & State.A)
}

// Read the address/value for an ASL, LSR, ROR, ROL
func preProcessShift(addressMode byte) (address uint16, value uint8) {
	if addressMode == amAccumulator {
		value = State.A
	} else {
		address, _ = getAddressFromAddressMode(addressMode)
		value = mmu.ReadMemory(address)
	}

	if addressMode == amAccumulator {
		value = State.A
	} else {
		address, _ = getAddressFromAddressMode(addressMode)
		value = mmu.ReadMemory(address)
	}

	return
}

// Store the result of a ASL, LSR, ROR, ROL and advance PC and system.FrameCycles
func postProcessShift(addressMode byte, address uint16, value uint8) {
	switch addressMode {
	case amAccumulator:
		State.A = value
		State.PC++
		system.FrameCycles += 2
	case amZeroPage:
		mmu.WriteMemory(address, value)
		State.PC += 2
		system.FrameCycles += 5
	case amZeroPageX:
		mmu.WriteMemory(address, value)
		State.PC += 2
		system.FrameCycles += 6
	case amAbsolute:
		mmu.WriteMemory(address, value)
		State.PC += 3
		system.FrameCycles += 6
	case amAbsoluteX:
		mmu.WriteMemory(address, value)
		State.PC += 3
		system.FrameCycles += 7
	default:
		panic(fmt.Sprintf("Unknown address mode %d in postProcessShift()", addressMode))
	}
}

func postProcessIncDec(addressMode byte) {
	switch addressMode {
	case amZeroPage:
		State.PC += 2
		system.FrameCycles += 5
	case amZeroPageX:
		State.PC += 2
		system.FrameCycles += 6
	case amAbsolute:
		State.PC += 3
		system.FrameCycles += 6
	case amAbsoluteX:
		State.PC += 3
		system.FrameCycles += 7
	default:
		panic(fmt.Sprintf("Unknown address mode %d in INC", addressMode))
	}
}

func brk() {
	push16(State.PC + 2)
	State.P |= cpuFlagB
	push8(State.P)
	State.P |= cpuFlagI
	State.PC = uint16(mmu.ReadMemory(0xffff))<<8 + uint16(mmu.ReadMemory(0xfffe))
	system.FrameCycles += 7
}

func irq() {
	push16(State.PC)
	State.P &= ^cpuFlagB
	push8(State.P)
	State.P |= cpuFlagI
	State.PC = uint16(mmu.ReadMemory(0xffff))<<8 + uint16(mmu.ReadMemory(0xfffe))
	system.FrameCycles += 7
}

func nmi() {
	push16(State.PC)
	State.P &= ^cpuFlagB
	push8(State.P)
	State.P |= cpuFlagI
	State.PC = uint16(mmu.ReadMemory(0xfffb))<<8 + uint16(mmu.ReadMemory(0xfffa))
	system.FrameCycles += 7
}

// Run runs the CPU until either wantedCycles has been reached (if non-zero) or the program counter reaches breakAddress.
// system.FrameCycles is the amount of cycles executed so far.
func Run(showInstructions bool, breakAddress *uint16, exitAtBreak bool, disableFirmwareWait bool, disableDosDelay bool, wantedCycles uint64) {
	system.FrameCycles = 0

	for {
		// Exit if wantedCycles is set and has been reached
		if (wantedCycles != 0) && (system.FrameCycles >= wantedCycles) {
			return
		}

		// Exit if the magic address of the functional tests has been reached
		if system.RunningTests && (State.PC == 0x3869) {
			fmt.Println("Functional tests passed")
			return
		}

		// Exit if the magic address of the interupt tests has been reached
		if system.RunningTests && (State.PC == 0x0af5) {
			fmt.Println("Interrupt tests passed")
			return
		}

		// Handle an IRQ f there is one pending and interrupts are enabled
		if system.PendingInterrupt && ((State.P & cpuFlagI) == 0) {
			irq()
			system.PendingInterrupt = false
			continue
		}

		// Handle an NMI if there is one pending
		if system.PendingNMI {
			nmi()
			system.PendingNMI = false
			continue
		}

		if showInstructions {
			PrintInstruction(true)
		}

		// Handle case of breakAddress being set and being been reached
		if breakAddress != nil && State.PC == *breakAddress {
			if exitAtBreak {
				// Exit the process completely
				fmt.Printf("Break at $%04x\n", *breakAddress)
				PrintInstruction(true)
				os.Exit(0)
			} else {
				// Exit politely
				return
			}
		}

		// Decode opcode
		opcode := mmu.ReadMemory(State.PC)
		addressMode := opCodes[opcode].addressingMode.mode

		switch opcode {

		case 0x4c: // JMP $0000
			value := uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
			if system.RunningTests && State.PC == value {
				// Check for an infinite loop and exit if so
				fmt.Printf("Trap at $%04x\n", value)
				os.Exit(0)
			}
			State.PC = value
			system.FrameCycles += 3
		case 0x6c: // JMP ($0000)
			value := uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
			State.PC = uint16(mmu.ReadMemory(value)) + uint16(mmu.ReadMemory(value+1))<<8
			system.FrameCycles += 5

		case 0x20: // JSR $0000
			value := uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
			system.FrameCycles += 6

			if disableFirmwareWait && value == 0xfca8 {
				// Don't call the firmware wait, just move forward and pretend it happened.
				State.PC += 3
				State.A = 0
				continue
			} else if disableDosDelay && value == 0xba00 {
				// Don't call the delay, just move forward and pretend it happened.
				State.PC += 3
				continue
			}

			push16(State.PC + 2)
			State.PC = value

		case 0x60: // RTS
			value := pop16()
			State.PC = value + 1
			system.FrameCycles += 6

		case 0xa9, 0xa5, 0xb5, 0xad, 0xbd, 0xb9, 0xa1, 0xb1: // LDA
			State.A = load(addressMode)
		case 0xa2, 0xa6, 0xb6, 0xae, 0xbe: // LDX
			State.X = load(addressMode)
		case 0xa0, 0xa4, 0xb4, 0xac, 0xbc: // LDY
			if disableDosDelay {
				pc := uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
				if pc == 0xbd9e {
					// Don't do delay, just move forward and pretend it happened.
					State.PC += 13
					continue
				}
			}

			State.Y = load(addressMode)

		case 0x85, 0x95, 0x8d, 0x9d, 0x99, 0x81, 0x91: //STA
			store(State.A, addressMode)
		case 0x86, 0x96, 0x8e: // STX
			store(State.X, addressMode)
		case 0x84, 0x94, 0x8c: //STY
			store(State.Y, addressMode)

		case 0xc9, 0xc5, 0xd5, 0xcd, 0xdd, 0xd9, 0xc1, 0xd1: // CMP
			cmp(State.A, addressMode)
		case 0xe0, 0xe4, 0xeC: // CPX
			cmp(State.X, addressMode)
		case 0xc0, 0xc4, 0xcc: // CPY
			cmp(State.Y, addressMode)
		case 0x09, 0x05, 0x15, 0x0d, 0x1d, 0x19, 0x01, 0x11: // ORA
			ora(addressMode)
		case 0x29, 0x25, 0x35, 0x2d, 0x3d, 0x39, 0x21, 0x31: // AND
			and(addressMode)
		case 0x49, 0x45, 0x55, 0x4d, 0x5d, 0x59, 0x41, 0x51: // EOR
			eor(addressMode)
		case 0x69, 0x65, 0x75, 0x6d, 0x7d, 0x79, 0x61, 0x71: // ADC
			adc(addressMode)
		case 0xe9, 0xe5, 0xf5, 0xed, 0xfd, 0xf9, 0xe1, 0xf1: // SBC
			sbc(addressMode)

		// Register transfers
		case 0xaa: // TAX
			State.X = State.A
			setN(State.X)
			setZ(State.X)
			State.PC++
			system.FrameCycles += 2
		case 0xa8: // TAY
			State.Y = State.A
			setN(State.Y)
			setZ(State.Y)
			State.PC++
			system.FrameCycles += 2
		case 0xba: // TSX
			State.X = State.SP
			setN(State.X)
			setZ(State.X)
			State.PC++
			system.FrameCycles += 2
		case 0x8a: // TXA
			State.A = State.X
			setN(State.A)
			setZ(State.A)
			State.PC++
			system.FrameCycles += 2
		case 0x9a: // TXS
			State.SP = State.X
			State.PC++
			system.FrameCycles += 2
		case 0x98: // TYA
			State.A = State.Y
			setN(State.A)
			setZ(State.A)
			State.PC++
			system.FrameCycles += 2

		case 0xE8: // INX
			State.X = (State.X + 1) & 0xff
			setN(State.X)
			setZ(State.X)
			State.PC++
			system.FrameCycles += 2
		case 0xC8: // INY
			State.Y = (State.Y + 1) & 0xff
			setN(State.Y)
			setZ(State.Y)
			State.PC++
			system.FrameCycles += 2
		case 0xca: // DEX
			State.X = (State.X - 1) & 0xff
			setN(State.X)
			setZ(State.X)
			State.PC++
			system.FrameCycles += 2
		case 0x88: // DEY
			State.Y = (State.Y - 1) & 0xff
			setN(State.Y)
			setZ(State.Y)
			State.PC++
			system.FrameCycles += 2

		// Branch instructions
		case 0x10:
			branch(!isN())
		case 0x30:
			branch(isN())
		case 0x50:
			branch(!isV())
		case 0x70:
			branch(isV())
		case 0x90:
			branch(!isC())
		case 0xb0:
			branch(isC())
		case 0xd0:
			branch(!isZ())
		case 0xf0:
			branch(isZ())

		// Flag setting
		case 0x18: // CLC
			setC(false)
			State.PC++
			system.FrameCycles += 2
		case 0x38: // SEC
			setC(true)
			State.PC++
			system.FrameCycles += 2
		case 0x58: // CLI
			State.P &= ^cpuFlagI
			State.PC++
			system.FrameCycles += 2
		case 0x78: // SEI
			State.P |= cpuFlagI
			State.PC++
			system.FrameCycles += 2
		case 0xb8: // CLV
			State.P &= ^cpuFlagV
			State.PC++
			system.FrameCycles += 2
		case 0xd8: // CLD
			State.P &= ^cpuFlagD
			State.PC++
			system.FrameCycles += 2
		case 0xf8: //SED
			State.P |= cpuFlagD
			State.PC++
			system.FrameCycles += 2

		// Stack operations
		case 0x48: // PHA
			push8(State.A)
			State.PC++
			system.FrameCycles += 3
		case 0x68: // PLA
			State.A = pop8()
			setN(State.A)
			setZ(State.A)
			State.PC++
			system.FrameCycles += 4
		case 0x08: // PHP
			// From http://visual6502.org/wiki/index.php?title=6502_BRK_and_B_bit#the_B_flag_and_the_various_mechanisms
			// software instructions BRK & PHP will push the B flag as being 1.
			push8(State.P | cpuFlagB)
			State.PC++
			system.FrameCycles += 3
		case 0x28: // PLP
			// cpuFlagR is always supposed to be 1
			State.P = pop8() | cpuFlagR
			State.PC++
			system.FrameCycles += 4
		case 0xea:
			State.PC++
			system.FrameCycles += 2

		case 0x00: // BRK
			brk()
		case 0x40: // RTI
			State.P = pop8() | cpuFlagR
			value := pop16()
			State.PC = value
			system.FrameCycles += 6

		case 0x24: // BIT $00
			address := mmu.ReadMemory(State.PC + 1)
			bit(uint16(address))
			State.PC += 2
			system.FrameCycles += 3
		case 0x2C: // BIT $0000
			address := uint16(mmu.ReadMemory(State.PC+1)) + uint16(mmu.ReadMemory(State.PC+2))<<8
			bit(address)
			State.PC += 3
			system.FrameCycles += 4

		// Shifts and rotations
		case 0x0a, 0x06, 0x16, 0x0e, 0x1e: // ASL
			address, value := preProcessShift(addressMode)
			setC((value & 0x80) != 0)
			value = (value << 1) & 0xff
			setZ(value)
			setN(value)
			postProcessShift(addressMode, address, value)
		case 0x4a, 0x46, 0x56, 0x4e, 0x5e: // LSR
			address, value := preProcessShift(addressMode)
			setC((value & 0x01) != 0)
			value >>= 1
			setZ(value)
			setN(value)
			postProcessShift(addressMode, address, value)
		case 0x2a, 0x26, 0x36, 0x2e, 0x3e: // ROL
			address, value := preProcessShift(addressMode)
			value16 := uint16(value)
			value16 <<= 1
			if (State.P & cpuFlagC) != 0 {
				value16 |= 0x01
			}
			setC((value16 & 0x100) != 0)
			value = uint8(value16 & 0xff)
			setZ(value)
			setN(value)
			postProcessShift(addressMode, address, value)
		case 0x6a, 0x66, 0x76, 0x6e, 0x7e: // ROR
			address, value := preProcessShift(addressMode)
			value16 := uint16(value)
			if (State.P & cpuFlagC) != 0 {
				value16 |= 0x100
			}
			setC((value16 & 0x01) != 0)
			value = uint8(value16 >> 1)
			setZ(value)
			setN(value)
			postProcessShift(addressMode, address, value)

		case 0xe6, 0xf6, 0xee, 0xfe: // INC
			address, _ := getAddressFromAddressMode(addressMode)
			value := mmu.ReadMemory(address)
			value = (value + 1) & 0xff
			setZ(value)
			setN(value)
			mmu.WriteMemory(address, value)
			postProcessIncDec(addressMode)

		case 0xc6, 0xd6, 0xce, 0xde: // DEC
			address, _ := getAddressFromAddressMode(addressMode)
			value := mmu.ReadMemory(address)
			value = (value - 1) & 0xff
			setZ(value)
			setN(value)
			mmu.WriteMemory(address, value)
			postProcessIncDec(addressMode)

		default:
			fmt.Printf("Unknown opcode $%02x at %04x\n", opcode, State.PC)
			return
		}
	}
}

// SetColdStartReset nukes the checkum byte for the reset vector. When this is called, the apple boot firmware will
// conclude that the reset vector is invalid and do a cold start.
func SetColdStartReset() {
	mmu.WriteMemory(uint16(0x3f4), uint8(0))
}

// Reset sets the CPU and memory states so that a next call to cpu.Run() calls the firmware reset code
func Reset() {
	mmu.InitROM() // Set upper memory area for reading from ROM
	mmu.InitRAM()

	bootVector := 0xfffc
	lsb := mmu.ReadPageTable[bootVector>>8][bootVector&0xff]
	msb := mmu.ReadPageTable[(bootVector+1)>>8][(bootVector+1)&0xff]
	State.PC = uint16(lsb) + uint16(msb)<<8
}

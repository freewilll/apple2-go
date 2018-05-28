package cpu

// Address mode constants
const (
	amNone        byte = 1 + iota
	amAccumulator      //
	amImplied          //
	amRelative         //
	amExpansion        //
	amImmediate        // #$00
	amZeroPage         // $00
	amZeroPageX        // $00,X
	amZeroPageY        // $00,Y
	amAbsolute         // $0000
	amAbsoluteX        // $0000,X
	amAbsoluteY        // $0000,Y
	amIndirect         // ($0000)
	amIndirectX        // ($00,X)
	amIndirectY        // ($00),Y
)

// addressingMode is a struct that describes a single addressing mode
type addressingMode struct {
	mode         byte   // One of the Am* constants
	operandSize  byte   // Number of bytes
	stringFormat string // Format string for the disassembler
}

type opCode struct {
	mnemonic       string         // 3-letter mnemonic, e.g. LDA, STA
	addressingMode addressingMode // Addressing mode
}

var addressingModes map[byte]addressingMode
var opCodes [0x100]opCode

func initAddressingModes() {
	addressingModes = make(map[byte]addressingMode)
	addressingModes[amAccumulator] = addressingMode{mode: amAccumulator, operandSize: 0, stringFormat: ""}
	addressingModes[amImplied] = addressingMode{mode: amImplied, operandSize: 0, stringFormat: ""}
	addressingModes[amRelative] = addressingMode{mode: amRelative, operandSize: 1, stringFormat: "$%04x"}
	addressingModes[amExpansion] = addressingMode{mode: amExpansion, operandSize: 0, stringFormat: ""}
	addressingModes[amImmediate] = addressingMode{mode: amImmediate, operandSize: 1, stringFormat: "#$%02x"}
	addressingModes[amZeroPage] = addressingMode{mode: amZeroPage, operandSize: 1, stringFormat: "$%02x"}
	addressingModes[amZeroPageX] = addressingMode{mode: amZeroPageX, operandSize: 1, stringFormat: "$%02x,X"}
	addressingModes[amZeroPageY] = addressingMode{mode: amZeroPageY, operandSize: 1, stringFormat: "$%02x,Y"}
	addressingModes[amAbsolute] = addressingMode{mode: amAbsolute, operandSize: 2, stringFormat: "$%04x"}
	addressingModes[amAbsoluteX] = addressingMode{mode: amAbsoluteX, operandSize: 2, stringFormat: "$%04x,X"}
	addressingModes[amAbsoluteY] = addressingMode{mode: amAbsoluteY, operandSize: 2, stringFormat: "$%04x,Y"}
	addressingModes[amIndirect] = addressingMode{mode: amIndirect, operandSize: 2, stringFormat: "($%04x)"}
	addressingModes[amIndirectX] = addressingMode{mode: amIndirectX, operandSize: 1, stringFormat: "($%02x,X)"}
	addressingModes[amIndirectY] = addressingMode{mode: amIndirectY, operandSize: 1, stringFormat: "($%02x),Y"}
}

func initOpCodes() {
	opCodes[0x00] = opCode{mnemonic: "BRK", addressingMode: addressingModes[amNone]}
	opCodes[0x01] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amIndirectX]}
	opCodes[0x02] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x03] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x04] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x05] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x06] = opCode{mnemonic: "ASL", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x07] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x08] = opCode{mnemonic: "PHP", addressingMode: addressingModes[amNone]}
	opCodes[0x09] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amImmediate]}
	opCodes[0x0A] = opCode{mnemonic: "ASL", addressingMode: addressingModes[amAccumulator]}
	opCodes[0x0B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x0C] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x0D] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x0E] = opCode{mnemonic: "ASL", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x0F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x10] = opCode{mnemonic: "BPL", addressingMode: addressingModes[amRelative]}
	opCodes[0x11] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amIndirectY]}
	opCodes[0x12] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x13] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x14] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x15] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x16] = opCode{mnemonic: "ASL", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x17] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x18] = opCode{mnemonic: "CLC", addressingMode: addressingModes[amNone]}
	opCodes[0x19] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0x1A] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x1B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x1C] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x1D] = opCode{mnemonic: "ORA", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x1E] = opCode{mnemonic: "ASL", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x1F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x20] = opCode{mnemonic: "JSR", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x21] = opCode{mnemonic: "AND", addressingMode: addressingModes[amIndirectX]}
	opCodes[0x22] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x23] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x24] = opCode{mnemonic: "BIT", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x25] = opCode{mnemonic: "AND", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x26] = opCode{mnemonic: "ROL", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x27] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x28] = opCode{mnemonic: "PLP", addressingMode: addressingModes[amNone]}
	opCodes[0x29] = opCode{mnemonic: "AND", addressingMode: addressingModes[amImmediate]}
	opCodes[0x2A] = opCode{mnemonic: "ROL", addressingMode: addressingModes[amAccumulator]}
	opCodes[0x2B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x2C] = opCode{mnemonic: "BIT", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x2D] = opCode{mnemonic: "AND", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x2E] = opCode{mnemonic: "ROL", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x2F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x30] = opCode{mnemonic: "BMI", addressingMode: addressingModes[amRelative]}
	opCodes[0x31] = opCode{mnemonic: "AND", addressingMode: addressingModes[amIndirectY]}
	opCodes[0x32] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x33] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x34] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x35] = opCode{mnemonic: "AND", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x36] = opCode{mnemonic: "ROL", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x37] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x38] = opCode{mnemonic: "SEC", addressingMode: addressingModes[amNone]}
	opCodes[0x39] = opCode{mnemonic: "AND", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0x3A] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x3B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x3C] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x3D] = opCode{mnemonic: "AND", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x3E] = opCode{mnemonic: "ROL", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x3F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x40] = opCode{mnemonic: "RTI", addressingMode: addressingModes[amNone]}
	opCodes[0x41] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amIndirectX]}
	opCodes[0x42] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x43] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x44] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x45] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x46] = opCode{mnemonic: "LSR", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x47] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x48] = opCode{mnemonic: "PHA", addressingMode: addressingModes[amNone]}
	opCodes[0x49] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amImmediate]}
	opCodes[0x4A] = opCode{mnemonic: "LSR", addressingMode: addressingModes[amAccumulator]}
	opCodes[0x4B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x4C] = opCode{mnemonic: "JMP", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x4D] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x4E] = opCode{mnemonic: "LSR", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x4F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x50] = opCode{mnemonic: "BVC", addressingMode: addressingModes[amRelative]}
	opCodes[0x51] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amIndirectY]}
	opCodes[0x52] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x53] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x54] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x55] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x56] = opCode{mnemonic: "LSR", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x57] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x58] = opCode{mnemonic: "CLI", addressingMode: addressingModes[amNone]}
	opCodes[0x59] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0x5A] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x5B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x5C] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x5D] = opCode{mnemonic: "EOR", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x5E] = opCode{mnemonic: "LSR", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x5F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x60] = opCode{mnemonic: "RTS", addressingMode: addressingModes[amNone]}
	opCodes[0x61] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amIndirectX]}
	opCodes[0x62] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x63] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x64] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x65] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x66] = opCode{mnemonic: "ROR", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x67] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x68] = opCode{mnemonic: "PLA", addressingMode: addressingModes[amNone]}
	opCodes[0x69] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amImmediate]}
	opCodes[0x6A] = opCode{mnemonic: "ROR", addressingMode: addressingModes[amAccumulator]}
	opCodes[0x6B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x6C] = opCode{mnemonic: "JMP", addressingMode: addressingModes[amIndirect]}
	opCodes[0x6D] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x6E] = opCode{mnemonic: "ROR", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x6F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x70] = opCode{mnemonic: "BVS", addressingMode: addressingModes[amRelative]}
	opCodes[0x71] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amIndirectY]}
	opCodes[0x72] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x73] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x74] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x75] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x76] = opCode{mnemonic: "ROR", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x77] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x78] = opCode{mnemonic: "SEI", addressingMode: addressingModes[amNone]}
	opCodes[0x79] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0x7A] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x7B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x7C] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x7D] = opCode{mnemonic: "ADC", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x7E] = opCode{mnemonic: "ROR", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x7F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x80] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x81] = opCode{mnemonic: "STA", addressingMode: addressingModes[amIndirectX]}
	opCodes[0x82] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x83] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x84] = opCode{mnemonic: "STY", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x85] = opCode{mnemonic: "STA", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x86] = opCode{mnemonic: "STX", addressingMode: addressingModes[amZeroPage]}
	opCodes[0x87] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x88] = opCode{mnemonic: "DEY", addressingMode: addressingModes[amNone]}
	opCodes[0x89] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x8A] = opCode{mnemonic: "TXA", addressingMode: addressingModes[amNone]}
	opCodes[0x8B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x8C] = opCode{mnemonic: "STY", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x8D] = opCode{mnemonic: "STA", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x8E] = opCode{mnemonic: "STX", addressingMode: addressingModes[amAbsolute]}
	opCodes[0x8F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x90] = opCode{mnemonic: "BCC", addressingMode: addressingModes[amRelative]}
	opCodes[0x91] = opCode{mnemonic: "STA", addressingMode: addressingModes[amIndirectY]}
	opCodes[0x92] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x93] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x94] = opCode{mnemonic: "STY", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x95] = opCode{mnemonic: "STA", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0x96] = opCode{mnemonic: "STX", addressingMode: addressingModes[amZeroPageY]}
	opCodes[0x97] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x98] = opCode{mnemonic: "TYA", addressingMode: addressingModes[amNone]}
	opCodes[0x99] = opCode{mnemonic: "STA", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0x9A] = opCode{mnemonic: "TXS", addressingMode: addressingModes[amNone]}
	opCodes[0x9B] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x9C] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x9D] = opCode{mnemonic: "STA", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0x9E] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0x9F] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xA0] = opCode{mnemonic: "LDY", addressingMode: addressingModes[amImmediate]}
	opCodes[0xA1] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amIndirectX]}
	opCodes[0xA2] = opCode{mnemonic: "LDX", addressingMode: addressingModes[amImmediate]}
	opCodes[0xA3] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xA4] = opCode{mnemonic: "LDY", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xA5] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xA6] = opCode{mnemonic: "LDX", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xA7] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xA8] = opCode{mnemonic: "TAY", addressingMode: addressingModes[amNone]}
	opCodes[0xA9] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amImmediate]}
	opCodes[0xAA] = opCode{mnemonic: "TAX", addressingMode: addressingModes[amNone]}
	opCodes[0xAB] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xAC] = opCode{mnemonic: "LDY", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xAD] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xAE] = opCode{mnemonic: "LDX", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xAF] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xB0] = opCode{mnemonic: "BCS", addressingMode: addressingModes[amRelative]}
	opCodes[0xB1] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amIndirectY]}
	opCodes[0xB2] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xB3] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xB4] = opCode{mnemonic: "LDY", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0xB5] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0xB6] = opCode{mnemonic: "LDX", addressingMode: addressingModes[amZeroPageY]}
	opCodes[0xB7] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xB8] = opCode{mnemonic: "CLV", addressingMode: addressingModes[amNone]}
	opCodes[0xB9] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0xBA] = opCode{mnemonic: "TSX", addressingMode: addressingModes[amNone]}
	opCodes[0xBB] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xBC] = opCode{mnemonic: "LDY", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0xBD] = opCode{mnemonic: "LDA", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0xBE] = opCode{mnemonic: "LDX", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0xBF] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xC0] = opCode{mnemonic: "CPY", addressingMode: addressingModes[amImmediate]}
	opCodes[0xC1] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amIndirectX]}
	opCodes[0xC2] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xC3] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xC4] = opCode{mnemonic: "CPY", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xC5] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xC6] = opCode{mnemonic: "DEC", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xC7] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xC8] = opCode{mnemonic: "INY", addressingMode: addressingModes[amNone]}
	opCodes[0xC9] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amImmediate]}
	opCodes[0xCA] = opCode{mnemonic: "DEX", addressingMode: addressingModes[amNone]}
	opCodes[0xCB] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xCC] = opCode{mnemonic: "CPY", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xCD] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xCE] = opCode{mnemonic: "DEC", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xCF] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xD0] = opCode{mnemonic: "BNE", addressingMode: addressingModes[amRelative]}
	opCodes[0xD1] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amIndirectY]}
	opCodes[0xD2] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xD3] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xD4] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xD5] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0xD6] = opCode{mnemonic: "DEC", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0xD7] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xD8] = opCode{mnemonic: "CLD", addressingMode: addressingModes[amNone]}
	opCodes[0xD9] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0xDA] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xDB] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xDC] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xDD] = opCode{mnemonic: "CMP", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0xDE] = opCode{mnemonic: "DEC", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0xDF] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xE0] = opCode{mnemonic: "CPX", addressingMode: addressingModes[amImmediate]}
	opCodes[0xE1] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amIndirectX]}
	opCodes[0xE2] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xE3] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xE4] = opCode{mnemonic: "CPX", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xE5] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xE6] = opCode{mnemonic: "INC", addressingMode: addressingModes[amZeroPage]}
	opCodes[0xE7] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xE8] = opCode{mnemonic: "INX", addressingMode: addressingModes[amNone]}
	opCodes[0xE9] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amImmediate]}
	opCodes[0xEA] = opCode{mnemonic: "NOP", addressingMode: addressingModes[amNone]}
	opCodes[0xEB] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xEC] = opCode{mnemonic: "CPX", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xED] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xEE] = opCode{mnemonic: "INC", addressingMode: addressingModes[amAbsolute]}
	opCodes[0xEF] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xF0] = opCode{mnemonic: "BEQ", addressingMode: addressingModes[amRelative]}
	opCodes[0xF1] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amIndirectY]}
	opCodes[0xF2] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xF3] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xF4] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xF5] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0xF6] = opCode{mnemonic: "INC", addressingMode: addressingModes[amZeroPageX]}
	opCodes[0xF7] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xF8] = opCode{mnemonic: "SED", addressingMode: addressingModes[amNone]}
	opCodes[0xF9] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amAbsoluteY]}
	opCodes[0xFA] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xFB] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xFC] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
	opCodes[0xFD] = opCode{mnemonic: "SBC", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0xFE] = opCode{mnemonic: "INC", addressingMode: addressingModes[amAbsoluteX]}
	opCodes[0xFF] = opCode{mnemonic: "???", addressingMode: addressingModes[amExpansion]}
}

func InitInstructionDecoder() {
	initAddressingModes()
	initOpCodes()
}

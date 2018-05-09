package cpu

const (
	AmNone        byte = 1 + iota
	AmAccumulator      //
	AmImplied          //
	AmRelative         //
	AmExpansion        //
	AmImmediate        // #$00
	AmZeroPage         // $00
	AmZeroPageX        // $00,X
	AmZeroPageY        // $00,Y
	AmAbsolute         // $0000
	AmAbsoluteX        // $0000,X
	AmAbsoluteY        // $0000,Y
	AmIndirect         // ($0000)
	AmIndirectX        // ($00,X)
	AmIndirectY        // ($00),Y
)

type AddressingMode struct {
	Mode         byte
	OperandSize  byte
	StringFormat string
}

type OpCode struct {
	Mnemonic       string
	AddressingMode AddressingMode
}

var AddressingModes map[byte]AddressingMode
var OpCodes [0x100]OpCode

func InitAddressingModes() {
	AddressingModes = make(map[byte]AddressingMode)
	AddressingModes[AmAccumulator] = AddressingMode{Mode: AmAccumulator, OperandSize: 0, StringFormat: ""}
	AddressingModes[AmImplied] = AddressingMode{Mode: AmImplied, OperandSize: 0, StringFormat: ""}
	AddressingModes[AmRelative] = AddressingMode{Mode: AmRelative, OperandSize: 1, StringFormat: "$%04x"}
	AddressingModes[AmExpansion] = AddressingMode{Mode: AmExpansion, OperandSize: 0, StringFormat: ""}
	AddressingModes[AmImmediate] = AddressingMode{Mode: AmImmediate, OperandSize: 1, StringFormat: "#$%02x"}
	AddressingModes[AmZeroPage] = AddressingMode{Mode: AmZeroPage, OperandSize: 1, StringFormat: "$%02x"}
	AddressingModes[AmZeroPageX] = AddressingMode{Mode: AmZeroPageX, OperandSize: 1, StringFormat: "$%02x,X"}
	AddressingModes[AmZeroPageY] = AddressingMode{Mode: AmZeroPageY, OperandSize: 1, StringFormat: "$%02x,Y"}
	AddressingModes[AmAbsolute] = AddressingMode{Mode: AmAbsolute, OperandSize: 2, StringFormat: "$%04x"}
	AddressingModes[AmAbsoluteX] = AddressingMode{Mode: AmAbsoluteX, OperandSize: 2, StringFormat: "$%04x,X"}
	AddressingModes[AmAbsoluteY] = AddressingMode{Mode: AmAbsoluteY, OperandSize: 2, StringFormat: "$%04x,Y"}
	AddressingModes[AmIndirect] = AddressingMode{Mode: AmIndirect, OperandSize: 2, StringFormat: "($%04x)"}
	AddressingModes[AmIndirectX] = AddressingMode{Mode: AmIndirectX, OperandSize: 1, StringFormat: "($%02x,X)"}
	AddressingModes[AmIndirectY] = AddressingMode{Mode: AmIndirectY, OperandSize: 1, StringFormat: "($%02x),Y"}
}

func InitOpCodes() {
	OpCodes[0x00] = OpCode{Mnemonic: "BRK", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x01] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0x02] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x03] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x04] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x05] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x06] = OpCode{Mnemonic: "ASL", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x07] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x08] = OpCode{Mnemonic: "PHP", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x09] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0x0A] = OpCode{Mnemonic: "ASL", AddressingMode: AddressingModes[AmAccumulator]}
	OpCodes[0x0B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x0C] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x0D] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x0E] = OpCode{Mnemonic: "ASL", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x0F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x10] = OpCode{Mnemonic: "BPL", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0x11] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0x12] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x13] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x14] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x15] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x16] = OpCode{Mnemonic: "ASL", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x17] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x18] = OpCode{Mnemonic: "CLC", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x19] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0x1A] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x1B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x1C] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x1D] = OpCode{Mnemonic: "ORA", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x1E] = OpCode{Mnemonic: "ASL", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x1F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x20] = OpCode{Mnemonic: "JSR", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x21] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0x22] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x23] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x24] = OpCode{Mnemonic: "BIT", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x25] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x26] = OpCode{Mnemonic: "ROL", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x27] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x28] = OpCode{Mnemonic: "PLP", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x29] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0x2A] = OpCode{Mnemonic: "ROL", AddressingMode: AddressingModes[AmAccumulator]}
	OpCodes[0x2B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x2C] = OpCode{Mnemonic: "BIT", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x2D] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x2E] = OpCode{Mnemonic: "ROL", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x2F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x30] = OpCode{Mnemonic: "BMI", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0x31] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0x32] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x33] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x34] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x35] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x36] = OpCode{Mnemonic: "ROL", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x37] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x38] = OpCode{Mnemonic: "SEC", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x39] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0x3A] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x3B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x3C] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x3D] = OpCode{Mnemonic: "AND", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x3E] = OpCode{Mnemonic: "ROL", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x3F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x40] = OpCode{Mnemonic: "RTI", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x41] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0x42] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x43] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x44] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x45] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x46] = OpCode{Mnemonic: "LSR", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x47] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x48] = OpCode{Mnemonic: "PHA", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x49] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0x4A] = OpCode{Mnemonic: "LSR", AddressingMode: AddressingModes[AmAccumulator]}
	OpCodes[0x4B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x4C] = OpCode{Mnemonic: "JMP", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x4D] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x4E] = OpCode{Mnemonic: "LSR", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x4F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x50] = OpCode{Mnemonic: "BVC", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0x51] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0x52] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x53] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x54] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x55] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x56] = OpCode{Mnemonic: "LSR", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x57] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x58] = OpCode{Mnemonic: "CLI", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x59] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0x5A] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x5B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x5C] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x5D] = OpCode{Mnemonic: "EOR", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x5E] = OpCode{Mnemonic: "LSR", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x5F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x60] = OpCode{Mnemonic: "RTS", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x61] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0x62] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x63] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x64] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x65] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x66] = OpCode{Mnemonic: "ROR", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x67] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x68] = OpCode{Mnemonic: "PLA", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x69] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0x6A] = OpCode{Mnemonic: "ROR", AddressingMode: AddressingModes[AmAccumulator]}
	OpCodes[0x6B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x6C] = OpCode{Mnemonic: "JMP", AddressingMode: AddressingModes[AmIndirect]}
	OpCodes[0x6D] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x6E] = OpCode{Mnemonic: "ROR", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x6F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x70] = OpCode{Mnemonic: "BVS", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0x71] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0x72] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x73] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x74] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x75] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x76] = OpCode{Mnemonic: "ROR", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x77] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x78] = OpCode{Mnemonic: "SEI", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x79] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0x7A] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x7B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x7C] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x7D] = OpCode{Mnemonic: "ADC", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x7E] = OpCode{Mnemonic: "ROR", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x7F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x80] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x81] = OpCode{Mnemonic: "STA", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0x82] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x83] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x84] = OpCode{Mnemonic: "STY", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x85] = OpCode{Mnemonic: "STA", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x86] = OpCode{Mnemonic: "STX", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0x87] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x88] = OpCode{Mnemonic: "DEY", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x89] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x8A] = OpCode{Mnemonic: "TXA", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x8B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x8C] = OpCode{Mnemonic: "STY", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x8D] = OpCode{Mnemonic: "STA", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x8E] = OpCode{Mnemonic: "STX", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0x8F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x90] = OpCode{Mnemonic: "BCC", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0x91] = OpCode{Mnemonic: "STA", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0x92] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x93] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x94] = OpCode{Mnemonic: "STY", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x95] = OpCode{Mnemonic: "STA", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0x96] = OpCode{Mnemonic: "STX", AddressingMode: AddressingModes[AmZeroPageY]}
	OpCodes[0x97] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x98] = OpCode{Mnemonic: "TYA", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x99] = OpCode{Mnemonic: "STA", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0x9A] = OpCode{Mnemonic: "TXS", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0x9B] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x9C] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x9D] = OpCode{Mnemonic: "STA", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0x9E] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0x9F] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xA0] = OpCode{Mnemonic: "LDY", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0xA1] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0xA2] = OpCode{Mnemonic: "LDX", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0xA3] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xA4] = OpCode{Mnemonic: "LDY", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xA5] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xA6] = OpCode{Mnemonic: "LDX", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xA7] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xA8] = OpCode{Mnemonic: "TAY", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xA9] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0xAA] = OpCode{Mnemonic: "TAX", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xAB] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xAC] = OpCode{Mnemonic: "LDY", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xAD] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xAE] = OpCode{Mnemonic: "LDX", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xAF] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xB0] = OpCode{Mnemonic: "BCS", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0xB1] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0xB2] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xB3] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xB4] = OpCode{Mnemonic: "LDY", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0xB5] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0xB6] = OpCode{Mnemonic: "LDX", AddressingMode: AddressingModes[AmZeroPageY]}
	OpCodes[0xB7] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xB8] = OpCode{Mnemonic: "CLV", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xB9] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0xBA] = OpCode{Mnemonic: "TSX", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xBB] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xBC] = OpCode{Mnemonic: "LDY", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0xBD] = OpCode{Mnemonic: "LDA", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0xBE] = OpCode{Mnemonic: "LDX", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0xBF] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xC0] = OpCode{Mnemonic: "CPY", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0xC1] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0xC2] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xC3] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xC4] = OpCode{Mnemonic: "CPY", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xC5] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xC6] = OpCode{Mnemonic: "DEC", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xC7] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xC8] = OpCode{Mnemonic: "INY", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xC9] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0xCA] = OpCode{Mnemonic: "DEX", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xCB] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xCC] = OpCode{Mnemonic: "CPY", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xCD] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xCE] = OpCode{Mnemonic: "DEC", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xCF] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xD0] = OpCode{Mnemonic: "BNE", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0xD1] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0xD2] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xD3] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xD4] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xD5] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0xD6] = OpCode{Mnemonic: "DEC", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0xD7] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xD8] = OpCode{Mnemonic: "CLD", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xD9] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0xDA] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xDB] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xDC] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xDD] = OpCode{Mnemonic: "CMP", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0xDE] = OpCode{Mnemonic: "DEC", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0xDF] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xE0] = OpCode{Mnemonic: "CPX", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0xE1] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmIndirectX]}
	OpCodes[0xE2] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xE3] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xE4] = OpCode{Mnemonic: "CPX", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xE5] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xE6] = OpCode{Mnemonic: "INC", AddressingMode: AddressingModes[AmZeroPage]}
	OpCodes[0xE7] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xE8] = OpCode{Mnemonic: "INX", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xE9] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmImmediate]}
	OpCodes[0xEA] = OpCode{Mnemonic: "NOP", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xEB] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xEC] = OpCode{Mnemonic: "CPX", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xED] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xEE] = OpCode{Mnemonic: "INC", AddressingMode: AddressingModes[AmAbsolute]}
	OpCodes[0xEF] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xF0] = OpCode{Mnemonic: "BEQ", AddressingMode: AddressingModes[AmRelative]}
	OpCodes[0xF1] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmIndirectY]}
	OpCodes[0xF2] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xF3] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xF4] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xF5] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0xF6] = OpCode{Mnemonic: "INC", AddressingMode: AddressingModes[AmZeroPageX]}
	OpCodes[0xF7] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xF8] = OpCode{Mnemonic: "SED", AddressingMode: AddressingModes[AmNone]}
	OpCodes[0xF9] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmAbsoluteY]}
	OpCodes[0xFA] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xFB] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xFC] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
	OpCodes[0xFD] = OpCode{Mnemonic: "SBC", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0xFE] = OpCode{Mnemonic: "INC", AddressingMode: AddressingModes[AmAbsoluteX]}
	OpCodes[0xFF] = OpCode{Mnemonic: "???", AddressingMode: AddressingModes[AmExpansion]}
}

func InitInstructionDecoder() {
	InitAddressingModes()
	InitOpCodes()
}

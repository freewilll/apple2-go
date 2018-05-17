package main

import (
	"mos6502go/cpu"
	"mos6502go/keyboard"
	"mos6502go/mmu"
	"mos6502go/system"
	"mos6502go/utils"
	"mos6502go/video"

	"testing"
)

const diskImage = "dos33_disk.dsk"

func writeBytes(address int, data []uint8) {
	for i := 0; i < len(data); i++ {
		mmu.WriteMemory(uint16(address)+uint16(i), data[i])
	}
}

func TestDos33RtsWriteRead(t *testing.T) {
	// Test writing and reading a sector using DOS 3.3's RWTS
	cpu.InitInstructionDecoder()
	mmu.InitRAM()
	mmu.InitApple2eROM()
	mmu.InitIO()
	mmu.ReadDiskImage(diskImage)
	cpu.Init()
	keyboard.Init()
	video.Init()
	system.Init()
	cpu.Reset()

	// Boot up DOS3.3
	utils.RunUntilBreakPoint(t, 0xd7d2, 5, false, "BASIC NEWSTT")

	// Write a sector from 0x2000 to track 35, sector 14
	start := 0x800
	writeBuffer := 0x2000
	readBuffer := 0x2100

	// Put some test data in
	for i := uint16(0); i < 0x100; i++ {
		mmu.WriteMemory(uint16(writeBuffer)+i, uint8(i)^0xaa)
	}

	writeBytes(start+0x00, []uint8{0x20, 0xe3, 0x03})                // JSR $03E3      LOCRPL = LOCATE RWTS PARAM LIST
	writeBytes(start+0x03, []uint8{0x84, 0x00})                      // STY $00
	writeBytes(start+0x05, []uint8{0x85, 0x01})                      // STA $01
	writeBytes(start+0x07, []uint8{0xa9, 0x22})                      // LDA #$22       track 34
	writeBytes(start+0x09, []uint8{0xa0, 0x04})                      // LDY #$04
	writeBytes(start+0x0b, []uint8{0x91, 0x00})                      // STA ($00),Y
	writeBytes(start+0x0d, []uint8{0xa9, 0x0e})                      // LDA #$0e       sector 14
	writeBytes(start+0x0f, []uint8{0xa0, 0x05})                      // LDY #$05
	writeBytes(start+0x11, []uint8{0x91, 0x00})                      // STA ($00),Y
	writeBytes(start+0x13, []uint8{0xa9, uint8(writeBuffer & 0xff)}) // LDA            writeBuffer lsb
	writeBytes(start+0x15, []uint8{0xa0, 0x08})                      // LDY #$08
	writeBytes(start+0x17, []uint8{0x91, 0x00})                      // STA ($00),Y
	writeBytes(start+0x19, []uint8{0xa9, uint8(writeBuffer >> 8)})   // LDA            writeBuffer msb
	writeBytes(start+0x1b, []uint8{0xa0, 0x09})                      // LDY #$09
	writeBytes(start+0x1d, []uint8{0x91, 0x00})                      // STA ($00),Y
	writeBytes(start+0x1f, []uint8{0xa9, 0x02})                      // LDA #$02       command=2 (write)
	writeBytes(start+0x21, []uint8{0xa0, 0x0c})                      // LDY #$0c
	writeBytes(start+0x23, []uint8{0x91, 0x00})                      // STA ($00),Y
	writeBytes(start+0x25, []uint8{0xa9, 0x00})                      // LDA #$00       any volume will do
	writeBytes(start+0x27, []uint8{0xa0, 0x03})                      // LDY #$03
	writeBytes(start+0x29, []uint8{0x91, 0x00})                      // STA ($00),Y
	writeBytes(start+0x2b, []uint8{0x20, 0xe3, 0x03})                // JSR $03E3      Relocate pointer to parms
	writeBytes(start+0x2f, []uint8{0x20, 0xd9, 0x03})                // JSR $03D9      RWTS
	writeBytes(start+0x32, []uint8{0x00})                            // BRK

	cpu.State.PC = uint16(start)
	utils.RunUntilBreakPoint(t, 0xb944, 128, false, "RWTS RDADDR")
	utils.RunUntilBreakPoint(t, 0xb82a, 8, false, "RWTS WRITESEC")
	utils.RunUntilBreakPoint(t, 0xb7ba, 8, false, "RWTS ENTERWTS")
	utils.RunUntilBreakPoint(t, uint16(start+0x32), 1, false, "Write routine break")

	// Now run some modified code to read the same track/sector
	writeBytes(start+0x13, []uint8{0xa9, uint8(readBuffer & 0xff)}) // LDA             readBuffer lsb
	writeBytes(start+0x15, []uint8{0xa0, 0x08})                     // LDY #$08
	writeBytes(start+0x17, []uint8{0x91, 0x00})                     // STA ($00),Y
	writeBytes(start+0x19, []uint8{0xa9, uint8(readBuffer >> 8)})   // LDA             readBuffer msb
	writeBytes(start+0x1b, []uint8{0xa0, 0x09})                     // LDY #$09
	writeBytes(start+0x1d, []uint8{0x91, 0x00})                     // STA ($00),Y
	writeBytes(start+0x1f, []uint8{0xa9, 0x01})                     // LDA #$01        command=1 (read)
	writeBytes(start+0x1b, []uint8{0xa0, 0x09})                     // LDY #$09
	writeBytes(start+0x1d, []uint8{0x91, 0x00})                     // STA ($00),Y

	cpu.State.PC = uint16(start)
	utils.RunUntilBreakPoint(t, uint16(start+0x32), 1, false, "Read routine break")

	// Check the read bytes match the witten ones
	for i := 0; i < 0x100; i++ {
		b1 := mmu.ReadMemory(uint16(readBuffer + i))
		b2 := mmu.ReadMemory(uint16(writeBuffer + i))
		if b1 != b2 {
			t.Fatalf("Mismatch at %02x: %02x vs %02x", readBuffer+i, b1, b2)
		}
	}
}

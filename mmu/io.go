package mmu

import (
	"fmt"
	"mos6502go/keyboard"
)

// https://mirrors.apple2.org.za/apple.cabi.net/Languages.Programming/MemoryMap.IIe.64K.128K.txt

const (
	KEYBOARD = 0xC000 // keyboard data (latched) (RD-only)
	CLR80COL = 0xC000 // use 80-column memory mapping (WR-only)
	SET80COL = 0xC001
	CLRAUXRD = 0xC002 // read from auxilliary 48K
	SETAUXRD = 0xC003
	CLRAUXWR = 0xC004 // write to auxilliary 48K
	SETAUXWR = 0xC005
	CLRCXROM = 0xC006 // use external slot ROM
	SETCXROM = 0xC007
	CLRAUXZP = 0xC008 // use auxilliary ZP, stack, & LC
	SETAUXZP = 0xC009
	CLRC3ROM = 0xC00A // use external slot C3 ROM
	SETC3ROM = 0xC00B
	CLR80VID = 0xC00C // use 80-column display mode
	SET80VID = 0xC00D
	CLRALTCH = 0xC00E // use alternate character set ROM
	SETALTCH = 0xC00F
	STROBE   = 0xC010 // strobe (unlatch) keyboard data

	RDLCBNK2 = 0xC011 // reading from LC bank $Dx 2
	RDLCRAM  = 0xC012 // reading from LC RAM
	RDRAMRD  = 0xC013 // reading from auxilliary 48K
	RDRAMWR  = 0xC014 // writing to auxilliary 48K
	RDCXROM  = 0xC015 // using external slot ROM
	RDAUXZP  = 0xC016 // using auxilliary ZP, stack, & LC
	RDC3ROM  = 0xC017 // using external slot C3 ROM
	RD80COL  = 0xC018 // using 80-column memory mapping
	RDVBLBAR = 0xC019 // not VBL (VBL signal low)
	RDTEXT   = 0xC01A // using text mode
	RDMIXED  = 0xC01B // using mixed mode
	RDPAGE2  = 0xC01C // using text/graphics page2
	RDHIRES  = 0xC01D // using Hi-res graphics mode
	RDALTCH  = 0xC01E // using alternate character set ROM
	RD80VID  = 0xC01F // using 80-column display mode
	SPEAKER  = 0xC030 // toggle speaker diaphragm

	CLRTEXT  = 0xC050 // enable text-only mode
	SETTEXT  = 0xC051
	CLRMIXED = 0xC052 // enable graphics/text mixed mode
	SETMIXED = 0xC053
	TXTPAGE1 = 0xC054 // select page1/2 (or page1/1x)
	TXTPAGE2 = 0xC055
	CLRHIRES = 0xC056 // enable Hi-res graphics
	SETHIRES = 0xC057

	SETAN0 = 0xC058 // 4-bit annunciator inputs
	CLRAN0 = 0xC059
	SETAN1 = 0xC05A
	CLRAN1 = 0xC05B
	SETAN2 = 0xC05C
	CLRAN2 = 0xC05D
	SETAN3 = 0xC05E
	CLRAN3 = 0xC05F

	OPNAPPLE = 0xC061 // open apple (command) key data
	CLSAPPLE = 0xC062 // closed apple (option) key data

	PDLTRIG = 0xC070 // trigger paddles
)

// Handle soft switch addresses where both a read and a write has a side
// effect and the return value is meaningless
func readWrite(address uint16) bool {
	switch address {
	case CLRTEXT:
		panic("CLRTEXT not implemented")
	case SETTEXT:
		return true
	case TXTPAGE1:
		return true
	case TXTPAGE2:
		return true
		fmt.Println("TXTPAGE2 not implemented")
		// panic("TXTPAGE2 not implemented")
		return true
	case CLRHIRES:
		return true
	case SETHIRES:
		panic("SETIRES not implemented")
	default:
		return false
	}
}

func ReadIO(address uint16) uint8 {
	if readWrite(address) {
		return 0
	}

	switch address {
	case KEYBOARD, STROBE:
		keyBoardData, strobe := keyboard.Read()
		if address == KEYBOARD {
			return keyBoardData
		} else {
			keyboard.ResetStrobe()
			return strobe
		}
	case RDCXROM:
		if UsingExternalSlotRom {
			return 0x8d
		} else {
			return 0x0d
		}
	case RD80VID:
		// using 80-column display mode not implemented
		return 0x0d

	// 4-bit annunciator inputs
	case SETAN0, CLRAN0, SETAN1, CLRAN1, SETAN2, CLRAN2, SETAN3, CLRAN3:
		// Annunciators not implemented
	case OPNAPPLE:
		// Open apple key not implemented
		return 0
	case CLSAPPLE:
		// Closed apple key not implemented
	case RD80COL:
		// RD80COL not implemented
		return 0x0d
	case RDPAGE2:
		// RDPAGE2 not implemented
		return 0x0d
	case RDALTCH:
		// RDALTCH not implemented
		return 0x0d
	case SPEAKER:
		// Speaker not implemented
		// Not printing anything since this will generate a lot of noise
	default:
		panic(fmt.Sprintf("TODO read %04x\n", address))
	}

	return 0
}

func WriteIO(address uint16, value uint8) {
	if readWrite(address) {
		return
	}

	switch address {
	case STROBE:
		keyboard.ResetStrobe()
	case CLRCXROM:
		MapFirstHalfOfIO()
	case SETCXROM:
		MapSecondHalfOfIO()
	case CLRALTCH:
		return
	case SETALTCH:
		panic("SETALTCH not implemented")
	case CLR80COL:
		// CLR80COL not implemented
		return
	case SET80COL:
		// SET80COL not implemented
	case CLRC3ROM:
		// CLRC3ROM not implemented
	case SETC3ROM:
		// SETC3ROM not implemented
	default:
		panic(fmt.Sprintf("TODO write %04x\n", address))
	}

	return
}

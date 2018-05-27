package mmu

import (
	"fmt"

	"apple2/audio"
	"apple2/disk"
	"apple2/keyboard"
	"apple2/system"
)

// Adapted from
// https://mirrors.apple2.org.za/apple.cabi.net/Languages.Programming/MemoryMap.IIe.64K.128K.txt
// https://github.com/cmosher01/Apple-II-Platform/blob/master/asminclude/iorom.asm

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
	STATEREG = 0xC068 // Has no effect on //e

	PDLTRIG = 0xC070 // trigger paddles

	// Slot 6 Drive IO
	S6CLRDRVP0 = 0xC0E0 // stepper phase 0  (Q0)
	S6SETDRVP0 = 0xC0E1 //
	S6CLRDRVP1 = 0xC0E2 // stepper phase 1  (Q1)
	S6SETDRVP1 = 0xC0E3 //
	S6CLRDRVP2 = 0xC0E4 // stepper phase 2  (Q2)
	S6SETDRVP2 = 0xC0E5 //
	S6CLRDRVP3 = 0xC0E6 // stepper phase 3  (Q3)
	S6SETDRVP3 = 0xC0E7 //
	S6MOTOROFF = 0xC0E8 // drive motor      (Q4)
	S6MOTORON  = 0xC0E9 //
	S6SELDRV1  = 0xC0EA // drive select     (Q5)
	S6SELDRV2  = 0xC0EB //
	S6Q6L      = 0xC0EC // read             (Q6)
	S6Q6H      = 0xC0ED // WP sense
	S6Q7L      = 0xC0EE // WP sense/read    (Q7)
	S6Q7H      = 0xC0EF // write
)

var VideoState struct {
	TextMode  bool
	HiresMode bool
	Mixed     bool
}

func InitIO() {
	// Empty slots that aren't yet implemented
	emptySlot(3)
	emptySlot(4)
	emptySlot(7)

	// Initialize slot 6 drive
	system.DriveState.Drive = 1
	system.DriveState.Spinning = false
	system.DriveState.Phase = 0
	system.DriveState.BytePosition = 0
	system.DriveState.Q6 = false
	system.DriveState.Q7 = false

	VideoState.TextMode = true
	VideoState.HiresMode = false
	VideoState.Mixed = false

	disk.InitDiskImage()
}

func driveIsreadSequencing() bool {
	return (!system.DriveState.Q6) && (!system.DriveState.Q7)
}

// Handle soft switch addresses where both a read and a write has a side
// effect and the return value is meaningless
func readWrite(address uint16, isRead bool) bool {
	lsb := address & 0xff
	if lsb >= 0x80 && lsb < 0x90 {
		SetMemoryMode(uint8(lsb - 0x80))
		return true
	}

	switch address {
	case CLRAUXRD:
		SetFakeAuxMemoryRead(false)
		return true
	case SETAUXRD:
		SetFakeAuxMemoryRead(true)
		return true
	case CLRAUXWR:
		SetFakeAuxMemoryWrite(false)
		return true
	case SETAUXWR:
		SetFakeAuxMemoryWrite(true)
		return true
	case CLRAUXZP:
		SetFakeAltZP(false)
		return true
	case SETAUXZP:
		SetFakeAltZP(true)
		return true
	case CLR80VID:
		SetCol80(false)
		return true
	case SET80VID:
		SetCol80(true)
		return true
	case TXTPAGE1:
		SetPage2(false)
		return true
	case TXTPAGE2:
		SetPage2(true)
		return true
	case CLRTEXT:
		VideoState.TextMode = false
		return true
	case SETTEXT:
		VideoState.TextMode = true
		return true
	case CLRMIXED:
		VideoState.Mixed = false
		return true
	case SETMIXED:
		VideoState.Mixed = true
		return true
	case CLRHIRES:
		VideoState.HiresMode = false
		return true
	case SETHIRES:
		VideoState.HiresMode = true
		return true
	case CLR80COL:
		if !isRead {
			SetStore80(false)
			return true
		} else {
			return false
		}
	case SET80COL:
		SetStore80(true)
		return true
	case STATEREG:
		// Ignore not implemented memory management reg
		return true

	// Drive stepper motor phase change
	case S6CLRDRVP0, S6SETDRVP0, S6CLRDRVP1, S6SETDRVP1, S6CLRDRVP2, S6SETDRVP2, S6CLRDRVP3, S6SETDRVP3:
		magnet := (address - S6CLRDRVP0) / 2
		on := ((address - S6CLRDRVP0) % 2) == 1

		if on {
			system.DriveState.Phases |= (1 << magnet)

			// Move head if a neighboring magnet is on and all others are off
			direction := int8(0)
			if (system.DriveState.Phases & (1 << uint8((system.DriveState.Phase+1)&3))) != 0 {
				direction += 1
			}
			if (system.DriveState.Phases & (1 << uint8((system.DriveState.Phase+3)&3))) != 0 {
				direction -= 1
			}

			if direction != 0 {
				system.DriveState.Phase += direction

				if system.DriveState.Phase < 0 {
					system.DriveState.Phase = 0
				}
				if system.DriveState.Phase == 80 {
					system.DriveState.Phase = 79
				}

				disk.MakeTrackData(uint8(system.DriveState.Phase))

				if audio.ClickWhenDriveHeadMoves {
					audio.Click()
				}
			}
		} else {
			system.DriveState.Phases &= ^(1 << magnet)
		}

		return true

	case S6MOTOROFF:
		system.DriveState.Spinning = false
		return true
	case S6MOTORON:
		system.DriveState.Spinning = true
		return true
	case S6SELDRV1:
		system.DriveState.Drive = 1
		return true
	case S6SELDRV2:
		system.DriveState.Drive = 2
		return true
	case S6Q6L:
		if !isRead {
			system.DriveState.Q6 = false
			return true
		}
		return false
	case S6Q6H:
		if isRead {
			system.DriveState.Q6 = true
			return true
		}
		return false
	case S6Q7L:
		system.DriveState.Q7 = false
		return true
	case S6Q7H:
		system.DriveState.Q7 = true
		return true

	default:
		return false
	}
}

func ReadIO(address uint16) uint8 {
	if readWrite(address, true) {
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
	case RDRAMRD, RDRAMWR, RDAUXZP:
		panic("Read/write aux memory not implemented")
		return 0x0d
	case RDCXROM:
		if UsingExternalSlotRom {
			return 0x8d
		} else {
			return 0x0d
		}
	case RD80VID:
		// using 80-column display mode not implemented
		return 0x0d

	case RDPAGE2:
		if Page2 {
			return 0x8d
		} else {
			return 0x0d
		}

	// 4-bit annunciator inputs
	case SETAN0, CLRAN0, SETAN1, CLRAN1, SETAN2, CLRAN2, SETAN3, CLRAN3:
		// Annunciators not implemented
	case OPNAPPLE:
		// Open apple key not implemented
		return 0
	case CLSAPPLE:
		// Closed apple key not implemented
	case RD80COL:
		if Store80 {
			return 0x8d
		} else {
			return 0x0d
		}
	case RDALTCH:
		// RDALTCH not implemented
		return 0x0d
	case SPEAKER:
		audio.Click()
	case S6Q6L:
		return disk.ReadTrackData()
	default:
		panic(fmt.Sprintf("TODO read %04x\n", address))
	}

	return 0
}

func WriteIO(address uint16, value uint8) {
	if readWrite(address, false) {
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
	case CLRC3ROM:
		// CLRC3ROM not implemented
	case SETC3ROM:
		// SETC3ROM not implemented
	case S6Q6H:
		disk.WriteTrackData(value)
	default:
		panic(fmt.Sprintf("TODO write %04x\n", address))
	}

	return
}

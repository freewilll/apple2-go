package mmu

import (
	"fmt"

	"github.com/freewilll/apple2-go/audio"
	"github.com/freewilll/apple2-go/disk"
	"github.com/freewilll/apple2-go/keyboard"
	"github.com/freewilll/apple2-go/system"
)

// Adapted from
// https://mirrors.apple2.org.za/apple.cabi.net/Languages.Programming/MemoryMap.IIe.64K.128K.txt
// https://github.com/cmosher01/Apple-II-Platform/blob/master/asminclude/iorom.asm

const (
	mKEYBOARD = 0xC000 // keyboard data (latched) (RD-only)
	mCLR80COL = 0xC000 // use 80-column memory mapping (WR-only)
	mSET80COL = 0xC001
	mCLRAUXRD = 0xC002 // read from auxilliary 48K
	mSETAUXRD = 0xC003
	mCLRAUXWR = 0xC004 // write to auxilliary 48K
	mSETAUXWR = 0xC005
	mCLRCXROM = 0xC006 // use external slot ROM
	mSETCXROM = 0xC007
	mCLRAUXZP = 0xC008 // use auxilliary ZP, stack, & LC
	mSETAUXZP = 0xC009
	mCLRC3ROM = 0xC00A // use external slot C3 ROM
	mSETC3ROM = 0xC00B
	mCLR80VID = 0xC00C // use 80-column display mode
	mSET80VID = 0xC00D
	mCLRALTCH = 0xC00E // use alternate character set ROM
	mSETALTCH = 0xC00F
	mSTROBE   = 0xC010 // strobe (unlatch) keyboard data

	mRDLCBNK2 = 0xC011 // reading from LC bank $Dx 2
	mRDLCRAM  = 0xC012 // reading from LC RAM
	mRDRAMRD  = 0xC013 // reading from auxilliary 48K
	mRDRAMWR  = 0xC014 // writing to auxilliary 48K
	mRDCXROM  = 0xC015 // using external slot ROM
	mRDAUXZP  = 0xC016 // using auxilliary ZP, stack, & LC
	mRDC3ROM  = 0xC017 // using external slot C3 ROM
	mRD80COL  = 0xC018 // using 80-column memory mapping
	mRDVBLBAR = 0xC019 // not VBL (VBL signal low)
	mRDTEXT   = 0xC01A // using text mode
	mRDMIXED  = 0xC01B // using mixed mode
	mRDPAGE2  = 0xC01C // using text/graphics page2
	mRDHIRES  = 0xC01D // using Hi-res graphics mode
	mRDALTCH  = 0xC01E // using alternate character set ROM
	mRD80VID  = 0xC01F // using 80-column display mode
	mSPEAKER  = 0xC030 // toggle speaker diaphragm

	mCLRTEXT  = 0xC050 // enable text-only mode
	mSETTEXT  = 0xC051
	mCLRMIXED = 0xC052 // enable graphics/text mixed mode
	mSETMIXED = 0xC053
	mTXTPAGE1 = 0xC054 // select page1/2 (or page1/1x)
	mTXTPAGE2 = 0xC055
	mCLRHIRES = 0xC056 // enable Hi-res graphics
	mSETHIRES = 0xC057

	mSETAN0 = 0xC058 // 4-bit annunciator inputs
	mCLRAN0 = 0xC059
	mSETAN1 = 0xC05A
	mCLRAN1 = 0xC05B
	mSETAN2 = 0xC05C
	mCLRAN2 = 0xC05D
	mSETAN3 = 0xC05E
	mCLRAN3 = 0xC05F

	mOPNAPPLE = 0xC061 // open apple (command) key data
	mCLSAPPLE = 0xC062 // closed apple (option) key data
	mSTATEREG = 0xC068 // Has no effect on //e

	mPDLTRIG = 0xC070 // trigger paddles

	// Slot 6 Drive IO
	mS6CLRDRVP0 = 0xC0E0 // stepper phase 0  (Q0)
	mS6SETDRVP0 = 0xC0E1 //
	mS6CLRDRVP1 = 0xC0E2 // stepper phase 1  (Q1)
	mS6SETDRVP1 = 0xC0E3 //
	mS6CLRDRVP2 = 0xC0E4 // stepper phase 2  (Q2)
	mS6SETDRVP2 = 0xC0E5 //
	mS6CLRDRVP3 = 0xC0E6 // stepper phase 3  (Q3)
	mS6SETDRVP3 = 0xC0E7 //
	mS6MOTOROFF = 0xC0E8 // drive motor      (Q4)
	mS6MOTORON  = 0xC0E9 //
	mS6SELDRV1  = 0xC0EA // drive select     (Q5)
	mS6SELDRV2  = 0xC0EB //
	mS6Q6L      = 0xC0EC // read             (Q6)
	mS6Q6H      = 0xC0ED // WP sense
	mS6Q7L      = 0xC0EE // WP sense/read    (Q7)
	mS6Q7H      = 0xC0EF // write
)

// VideoState has 3 booleans which determine the video configuration:
//                    TextMode HiresMode Mixed
// text				  1        0 		 N/A
// lores + text		  0        0         1
// lores              0        0         0
// hires              N/A      1         0
// hires + text       N/A      1         1
var VideoState struct {
	TextMode  bool
	HiresMode bool
	Mixed     bool
}

// InitIO resets all IO states
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

	// Initialize video
	VideoState.TextMode = true
	VideoState.HiresMode = false
	VideoState.Mixed = false

	disk.InitDiskImage()
}

// Handle soft switch addresses between $c000-$c0ff where both a read and a write has a side
// effect. Returns true if the read/write has been handled.
func readWrite(address uint16, isRead bool) bool {
	lsb := address & 0xff
	if lsb >= 0x80 && lsb < 0x90 {
		SetMemoryMode(uint8(lsb - 0x80))
		return true
	}

	switch address {
	case mCLRAUXRD:
		SetFakeAuxMemoryRead(false)
		return true
	case mSETAUXRD:
		SetFakeAuxMemoryRead(true)
		return true

	case mCLRAUXWR:
		SetFakeAuxMemoryWrite(false)
		return true
	case mSETAUXWR:
		SetFakeAuxMemoryWrite(true)
		return true

	case mCLRAUXZP:
		SetFakeAltZP(false)
		return true
	case mSETAUXZP:
		SetFakeAltZP(true)
		return true

	case mCLR80VID:
		SetCol80(false)
		return true
	case mSET80VID:
		SetCol80(true)
		return true

	case mTXTPAGE1:
		SetPage2(false)
		return true
	case mTXTPAGE2:
		SetPage2(true)
		return true

	case mCLRTEXT:
		VideoState.TextMode = false
		return true
	case mSETTEXT:
		VideoState.TextMode = true
		return true

	case mCLRMIXED:
		VideoState.Mixed = false
		return true
	case mSETMIXED:
		VideoState.Mixed = true
		return true

	case mCLRHIRES:
		VideoState.HiresMode = false
		return true
	case mSETHIRES:
		VideoState.HiresMode = true
		return true

	case mCLR80COL:
		if !isRead {
			SetStore80(false)
			return true
		}
		return false

	case mSET80COL:
		SetStore80(true)
		return true

	case mSTATEREG:
		// Ignore not implemented memory management reg
		return true

	// Drive stepper motor phase change
	case mS6CLRDRVP0, mS6SETDRVP0, mS6CLRDRVP1, mS6SETDRVP1, mS6CLRDRVP2, mS6SETDRVP2, mS6CLRDRVP3, mS6SETDRVP3:
		magnet := (address - mS6CLRDRVP0) / 2
		on := ((address - mS6CLRDRVP0) % 2) == 1
		if !on {
			// Turn off the magnet in Phases
			system.DriveState.Phases &= ^(1 << magnet)
			return true
		}

		// Implicit else, a magnet has been switched on
		system.DriveState.Phases |= (1 << magnet)

		// Move head if a neighboring magnet is on and all others are off
		direction := int8(0)
		if (system.DriveState.Phases & (1 << uint8((system.DriveState.Phase+1)&3))) != 0 {
			direction++
		}
		if (system.DriveState.Phases & (1 << uint8((system.DriveState.Phase+3)&3))) != 0 {
			direction--
		}

		// Move the head
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

		return true

	case mS6MOTOROFF:
		system.DriveState.Spinning = false
		return true
	case mS6MOTORON:
		system.DriveState.Spinning = true
		return true

	case mS6SELDRV1:
		system.DriveState.Drive = 1
		return true
	case mS6SELDRV2:
		system.DriveState.Drive = 2
		return true

	case mS6Q6L:
		if !isRead {
			system.DriveState.Q6 = false
			return true
		}
		return false
	case mS6Q6H:
		if isRead {
			system.DriveState.Q6 = true
			return true
		}
		return false

	case mS6Q7L:
		system.DriveState.Q7 = false
		return true
	case mS6Q7H:
		system.DriveState.Q7 = true
		return true

	default:
		return false
	}
}

// ReadIO does a read in the $c000-$c0ff area
func ReadIO(address uint16) uint8 {
	// Try the generic readWrite and return if it has handled the read
	if readWrite(address, true) {
		return 0
	}

	switch address {

	case mKEYBOARD, mSTROBE:
		keyBoardData, strobe := keyboard.Read()
		if address == mKEYBOARD {
			return keyBoardData
		}
		keyboard.ResetStrobe()
		return strobe

	case mRDRAMRD, mRDRAMWR, mRDAUXZP:
		panic("Read/write aux memory not implemented")
		return 0x0d

	case mRDCXROM:
		if UsingExternalSlotRom {
			return 0x8d
		}
		return 0x0d

	case mRD80VID:
		// using 80-column display mode not implemented
		return 0x0d

	case mRDPAGE2:
		if Page2 {
			return 0x8d
		}
		return 0x0d

	// 4-bit annunciator inputs
	case mSETAN0, mCLRAN0, mSETAN1, mCLRAN1, mSETAN2, mCLRAN2, mSETAN3, mCLRAN3:
		// Annunciators not implemented

	case mOPNAPPLE:
		// Open apple key not implemented
		return 0

	case mCLSAPPLE:
		// Closed apple key not implemented

	case mRD80COL:
		if Store80 {
			return 0x8d
		}
		return 0x0d

	case mRDALTCH:
		// RDALTCH not implemented, but it's also used, so don't fail on it.
		return 0x0d

	case mSPEAKER:
		audio.Click()
		return 0

	case mS6Q6L:
		// A read from disk
		return disk.ReadTrackData()

	default:
		panic(fmt.Sprintf("TODO read %04x\n", address))
	}

	return 0
}

// WriteIO does a write in the $c000-$c0ff area
func WriteIO(address uint16, value uint8) {
	// Try the generic readWrite and return if it has handled the write
	if readWrite(address, false) {
		return
	}

	switch address {

	case mSTROBE:
		keyboard.ResetStrobe()

	case mCLRCXROM:
		MapFirstHalfOfIO()
	case mSETCXROM:
		MapSecondHalfOfIO()

	case mCLRALTCH:
		return
	case mSETALTCH:
		panic("SETALTCH not implemented")

	case mCLR80COL:
		// CLR80COL not implemented
		return

	case mCLRC3ROM:
		// CLRC3ROM not implemented
	case mSETC3ROM:
		// SETC3ROM not implemented

	case mS6Q6H:
		// A write to disk
		disk.WriteTrackData(value)

	default:
		panic(fmt.Sprintf("TODO write %04x\n", address))
	}

	return
}

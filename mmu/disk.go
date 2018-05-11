package mmu

import (
	"fmt"
	"io/ioutil"
)

const tracksPerDisk = 35
const sectorsPerTrack = 16
const imageLength = tracksPerDisk * sectorsPerTrack * 0x100

// Each sector has
// Address field prologue               0x003 bytes
// Volume, Track, Sector, Checksum      0x008 bytes
// Address field epilogue               0x003 bytes
// Data prologue                        0x003 bytes
// 2-bits                               0x056 bytes
// 6-bits                               0x100 bytes
// checksum                             0x001 byte
// Data epilogue                        0x003 bytes
const diskSectorBytes = 3 + 8 + 3 + 3 + 0x56 + 0x100 + 1 + 3
const trackDataBytes = sectorsPerTrack * diskSectorBytes

var dos33SectorInterleaving [16]uint8
var translateTable62 [0x40]uint8 // Conversion of a 6 bit byte to a 8 bit "disk" byte

type sector struct {
	data [0x100]uint8
}

type track struct {
	sectors [sectorsPerTrack]sector
}

type disk struct {
	tracks [tracksPerDisk]track
}

var image disk                      // A loaded disk image
var TrackData [trackDataBytes]uint8 // Converted image data as it it returned by the disk controller for a single track

func InitDiskImage() {
	// Map DOS 3.3 sector interleaving
	// Physical sector -> Logical sector
	dos33SectorInterleaving[0x0] = 0x0
	dos33SectorInterleaving[0x1] = 0x7
	dos33SectorInterleaving[0x2] = 0xe
	dos33SectorInterleaving[0x3] = 0x6
	dos33SectorInterleaving[0x4] = 0xd
	dos33SectorInterleaving[0x5] = 0x5
	dos33SectorInterleaving[0x6] = 0xc
	dos33SectorInterleaving[0x7] = 0x4
	dos33SectorInterleaving[0x8] = 0xb
	dos33SectorInterleaving[0x9] = 0x3
	dos33SectorInterleaving[0xa] = 0xa
	dos33SectorInterleaving[0xb] = 0x2
	dos33SectorInterleaving[0xc] = 0x9
	dos33SectorInterleaving[0xd] = 0x1
	dos33SectorInterleaving[0xe] = 0x8
	dos33SectorInterleaving[0xf] = 0xf

	// Zero disk image data
	for t := 0; t < tracksPerDisk; t++ {
		for s := 0; s < sectorsPerTrack; s++ {
			for i := 0; i < 0x100; i++ {
				image.tracks[t].sectors[s].data[i] = 0
			}
		}
	}

	// Convert a 6 bit "byte" to a 8 bit "disk" byte
	translateTable62 = [0x40]uint8{
		0x96, 0x97, 0x9a, 0x9b, 0x9d, 0x9e, 0x9f, 0xa6,
		0xa7, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb2, 0xb3,
		0xb4, 0xb5, 0xb6, 0xb7, 0xb9, 0xba, 0xbb, 0xbc,
		0xbd, 0xbe, 0xbf, 0xcb, 0xcd, 0xce, 0xcf, 0xd3,
		0xd6, 0xd7, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde,
		0xdf, 0xe5, 0xe6, 0xe7, 0xe9, 0xea, 0xeb, 0xec,
		0xed, 0xee, 0xef, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6,
		0xf7, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
	}
}

func ReadDiskImage(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Unable to read ROM: %s", err))
	}

	if len(bytes) != imageLength {
		panic(fmt.Sprintf("Disk image has invalid length %d, expected %d", len(bytes), imageLength))
	}

	pos := 0
	for t := 0; t < tracksPerDisk; t++ {
		for s := 0; s < sectorsPerTrack; s++ {
			for i := 0; i < 0x100; i++ {
				image.tracks[t].sectors[s].data[i] = bytes[pos]
				pos++
			}
		}
	}
}

// Encode a byte into two 4-bit bytes with odd-even encoding. This is used
// for the sector and data headers
func oddEvenEncode(data uint8) (uint8, uint8) {
	bit0 := (data & 0x01) >> 0
	bit1 := (data & 0x02) >> 1
	bit2 := (data & 0x04) >> 2
	bit3 := (data & 0x08) >> 3
	bit4 := (data & 0x10) >> 4
	bit5 := (data & 0x20) >> 5
	bit6 := (data & 0x40) >> 6
	bit7 := (data & 0x80) >> 7

	l := 0xaa | (bit7 << 6) | (bit5 << 4) | (bit3 << 2) | (bit1)
	h := 0xaa | (bit6 << 6) | (bit4 << 4) | (bit2 << 2) | (bit0)
	return l, h
}

// Convert 8 bit bytes to 0x56 2-bit sections and 0x100 6 bit sections
func makeSectorData(s sector) (data [0x56 + 0x100]uint8) {
	twoBitPos := 0x0
	for i := 0; i < 0x100; i++ {
		b := s.data[i]
		bit0 := b & 0x1
		bit1 := (b & 0x2) >> 1
		data[twoBitPos] = (data[twoBitPos] >> 2) | (bit0 << 5) | (bit1 << 4)
		data[i+0x56] = b >> 2

		twoBitPos++
		if twoBitPos == 0x56 {
			twoBitPos = 0x0
		}
	}

	// Make sure the bits for 2 remainders from the 256 divide by 3 are in the right place.
	data[0x54] = (data[0x54] >> 2)
	data[0x55] = (data[0x55] >> 2)

	return
}

func clearTrackData() {
	for i := 0; i < trackDataBytes; i++ {
		TrackData[i] = 0
	}
}

func MakeTrackData(armPosition uint8) {
	// Tracks are present on even arm positions.
	track := armPosition / 2

	// If it's an odd arm position or a track beyond the image, zero the data
	if (armPosition >= (tracksPerDisk * 2)) || ((armPosition % 2) == 1) {
		clearTrackData()
		return
	}

	DriveState.BytePosition = 0 // Point the head at the first sector

	// For each sector, encode the data and add it to TrackData
	for physicalSector := 0; physicalSector < sectorsPerTrack; physicalSector++ {
		logicalSector := dos33SectorInterleaving[physicalSector]
		offset := int(physicalSector) * diskSectorBytes

		volume := uint8(254) // Volume numbers aren't implemented
		checksum := volume ^ track ^ uint8(physicalSector)

		volL, volH := oddEvenEncode(volume)
		trL, trH := oddEvenEncode(track)
		seL, seH := oddEvenEncode(uint8(physicalSector))
		csL, csH := oddEvenEncode(checksum)

		// Address field prologue
		TrackData[offset+0] = 0xd5
		TrackData[offset+1] = 0xaa
		TrackData[offset+2] = 0x96

		// Volume, track, sector and checksum
		TrackData[offset+3] = volL
		TrackData[offset+4] = volH
		TrackData[offset+5] = trL
		TrackData[offset+6] = trH
		TrackData[offset+7] = seL
		TrackData[offset+8] = seH
		TrackData[offset+9] = csL
		TrackData[offset+10] = csH

		// Address epilogue
		TrackData[offset+11] = 0xde
		TrackData[offset+12] = 0xaa
		TrackData[offset+13] = 0xeb

		// Data field prologue
		TrackData[offset+14] = 0xd5
		TrackData[offset+15] = 0xaa
		TrackData[offset+16] = 0xad

		sectorData := makeSectorData(image.tracks[track].sectors[logicalSector])

		// a is the previous byte's value
		a := uint8(0)
		for i := 0; i < 0x56+0x100; i++ {
			a ^= sectorData[i]
			b := translateTable62[a]
			TrackData[offset+17+i] = b
			a = sectorData[i]
		}

		// Set the checksum byte
		TrackData[offset+17+0x56+0x100] = translateTable62[a]

		// Data epilogue
		TrackData[offset+17+0x56+0x100+1] = 0xde
		TrackData[offset+17+0x56+0x100+2] = 0xaa
		TrackData[offset+17+0x56+0x100+3] = 0xeb

	}
}

// Read a byte from the disk head and spin the disk along
func ReadTrackData() (result uint8) {
	result = TrackData[DriveState.BytePosition]

	DriveState.BytePosition++
	if DriveState.BytePosition == trackDataBytes {
		DriveState.BytePosition = 0
	}

	return
}

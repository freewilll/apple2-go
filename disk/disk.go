package disk

import (
	"fmt"
	"io/ioutil"

	"github.com/freewilll/apple2/system"
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
var sixTwoEncoding [0x40]uint8  // Conversion of a 6 bit byte to a 8 bit "disk" byte
var sixTwoDecoding [0x100]uint8 // Conversion of a 8 bit "disk" byte to a 6 bit byte

type sector struct {
	data [0x100]uint8
}

type track struct {
	sectors [sectorsPerTrack]sector
}

type disk struct {
	tracks [tracksPerDisk]track
}

var imagePath string                // Loaded disk image path
var image disk                      // A loaded disk image
var imageIsDirty bool               // If an image has been written to and needs a flush
var trackData [trackDataBytes]uint8 // Converted image data as it it returned by the disk controller for a single track

// vars to keep track of writes
const (
	waitingForDataPrologue byte = 1 + iota
	receivingData
)

const rawDataBufferSize = diskSectorBytes + 16

type addressField struct {
	volume uint8
	track  uint8
	sector uint8
}

var lastReadAddress addressField
var lastReadSectorDataPosition int

var sectorWriteState struct {
	State           byte
	RawData         [rawDataBufferSize]uint8
	RawDataPosition uint16
	Address         addressField
}

func resetsectorWriteState() {
	sectorWriteState.State = waitingForDataPrologue
	sectorWriteState.RawDataPosition = 0
}

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
	sixTwoEncoding = [0x40]uint8{
		0x96, 0x97, 0x9a, 0x9b, 0x9d, 0x9e, 0x9f, 0xa6,
		0xa7, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb2, 0xb3,
		0xb4, 0xb5, 0xb6, 0xb7, 0xb9, 0xba, 0xbb, 0xbc,
		0xbd, 0xbe, 0xbf, 0xcb, 0xcd, 0xce, 0xcf, 0xd3,
		0xd6, 0xd7, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde,
		0xdf, 0xe5, 0xe6, 0xe7, 0xe9, 0xea, 0xeb, 0xec,
		0xed, 0xee, 0xef, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6,
		0xf7, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
	}

	for i := uint8(0); i < 0x40; i++ {
		sixTwoDecoding[sixTwoEncoding[i]] = i
	}

	resetsectorWriteState()
}

func ReadDiskImage(path string) {
	imagePath = path

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Unable to read disk image: %s", err))
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

	imageIsDirty = false
}

func writeDiskImage() {
	bytes := make([]byte, tracksPerDisk*sectorsPerTrack*0x100)

	pos := 0
	for t := 0; t < tracksPerDisk; t++ {
		for s := 0; s < sectorsPerTrack; s++ {
			for i := 0; i < 0x100; i++ {
				bytes[pos] = byte(image.tracks[t].sectors[s].data[i])
				pos++
			}
		}
	}

	err := ioutil.WriteFile(imagePath, bytes, 0644)
	if err != nil {
		panic(fmt.Sprintf("Unable to write disk image: %s", err))
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

// Merge the two bytes together produce by oddEvenEncode
func oddEvenDecode(data0 byte, data1 byte) uint8 {
	return ((data0 << 1) | 1) & data1
}

// Convert 8 bit bytes to 0x56 2-bit sections and 0x100 6 bit sections
func sectorDataEncode(s sector) (data [0x56 + 0x100]uint8) {
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

func sectorDataDecode(data []uint8) (sector [0x100]uint8) {
	for i := 0; i < 0x100; i++ {
		sector[i] = data[i+0x56]
	}

	twoBitPos := 0x00
	for i := 0; i < 0x100; i++ {
		twoBit := data[twoBitPos]
		sector[i] = (sector[i] << 2) + ((twoBit & 1) << 1) + ((twoBit & 2) >> 1)
		data[twoBitPos] >>= 2

		twoBitPos++
		if twoBitPos == 0x56 {
			twoBitPos = 0x0
		}
	}

	return
}

func clearTrackData() {
	for i := 0; i < trackDataBytes; i++ {
		trackData[i] = 0
	}
}

func makeSectorData(track uint8, physicalSector uint8) {
	logicalSector := dos33SectorInterleaving[physicalSector]
	offset := int(physicalSector) * diskSectorBytes

	volume := uint8(254) // Volume numbers aren't implemented
	checksum := volume ^ track ^ uint8(physicalSector)

	volL, volH := oddEvenEncode(volume)
	trL, trH := oddEvenEncode(track)
	seL, seH := oddEvenEncode(uint8(physicalSector))
	csL, csH := oddEvenEncode(checksum)

	// Address field prologue
	trackData[offset+0] = 0xd5
	trackData[offset+1] = 0xaa
	trackData[offset+2] = 0x96

	// Volume, track, sector and checksum
	trackData[offset+3] = volL
	trackData[offset+4] = volH
	trackData[offset+5] = trL
	trackData[offset+6] = trH
	trackData[offset+7] = seL
	trackData[offset+8] = seH
	trackData[offset+9] = csL
	trackData[offset+10] = csH

	// Address epilogue
	trackData[offset+11] = 0xde
	trackData[offset+12] = 0xaa
	trackData[offset+13] = 0xeb

	// Data field prologue
	trackData[offset+14] = 0xd5
	trackData[offset+15] = 0xaa
	trackData[offset+16] = 0xad

	sectorData := sectorDataEncode(image.tracks[track].sectors[logicalSector])

	// a is the previous byte's value
	a := uint8(0)
	for i := 0; i < 0x56+0x100; i++ {
		a ^= sectorData[i]
		b := sixTwoEncoding[a]
		trackData[offset+17+i] = b
		a = sectorData[i]
	}

	// Set the checksum byte
	trackData[offset+17+0x56+0x100] = sixTwoEncoding[a]

	// Data epilogue
	trackData[offset+17+0x56+0x100+1] = 0xde
	trackData[offset+17+0x56+0x100+2] = 0xaa
	trackData[offset+17+0x56+0x100+3] = 0xeb
}

func MakeTrackData(armPosition uint8) {
	// Tracks are present on even arm positions.
	track := uint8(armPosition / 2)

	// If it's an odd arm position or a track beyond the image, zero the data
	if (armPosition >= (tracksPerDisk * 2)) || ((armPosition % 2) == 1) {
		clearTrackData()
		return
	}

	system.DriveState.BytePosition = 0 // Point the head at the first sector

	// For each sector, encode the data and add it to trackData
	for physicalSector := uint8(0); physicalSector < sectorsPerTrack; physicalSector++ {
		makeSectorData(track, physicalSector)
	}
}

func decodeAddressField(data []uint8) addressField {
	var af addressField
	af.volume = oddEvenDecode(data[0], data[1])
	af.track = oddEvenDecode(data[2], data[3])
	af.sector = oddEvenDecode(data[4], data[5])
	return af
}

// Read a byte from the disk head and spin the disk along
func ReadTrackData() (result uint8) {
	result = trackData[system.DriveState.BytePosition]

	if system.DriveState.BytePosition >= 9 {
		if trackData[system.DriveState.BytePosition-9] == 0xd5 &&
			trackData[system.DriveState.BytePosition-8] == 0xaa &&
			trackData[system.DriveState.BytePosition-7] == 0x96 {
			var addressData []uint8
			addressData = trackData[system.DriveState.BytePosition-6 : system.DriveState.BytePosition]
			lastReadAddress = decodeAddressField(addressData)
			lastReadSectorDataPosition = system.DriveState.BytePosition + 8
		}
	}

	system.DriveState.BytePosition++
	if system.DriveState.BytePosition == trackDataBytes {
		system.DriveState.BytePosition = 0
	}

	return
}

// WriteTrackData gets called whenever a byte is written to the write address.
// Reads are done at the same time by the OS to await the drive to be in the
// right position. The last read address determines the track and sector. The expeted sequence of writes are:
// - up to 5 bytes of 0xff padding (ignored)
// - data prologue d5 aa ad
// - 0x56 bytes of 2-bit data
// - 0x56 bytes of 6-bit data
// - checksum byte (ignored)
// - data epilogue (ignored)
//
// The sector is decoded and updated in memory once the 0x156 data  bytes have
// been read. The image is flagged as dirty and flushed on exit.
func WriteTrackData(value uint8) {
	if sectorWriteState.State == waitingForDataPrologue {
		if sectorWriteState.RawDataPosition >= 16 {
			resetsectorWriteState()
			return
		}

		sectorWriteState.RawData[sectorWriteState.RawDataPosition] = value
		sectorWriteState.RawDataPosition += 1

		// Check for address prologue
		if sectorWriteState.RawDataPosition > 2 && sectorWriteState.RawData[sectorWriteState.RawDataPosition-3] == 0xd5 &&
			sectorWriteState.RawData[sectorWriteState.RawDataPosition-2] == 0xaa &&
			sectorWriteState.RawData[sectorWriteState.RawDataPosition-1] == 0xad {

			// We got it, record the last read address field and reset RawDataPosition
			sectorWriteState.State = receivingData
			sectorWriteState.Address = lastReadAddress
			sectorWriteState.RawDataPosition = 0
			return
		}

	} else if sectorWriteState.State == receivingData {
		sectorWriteState.RawData[sectorWriteState.RawDataPosition] = value
		sectorWriteState.RawDataPosition += 1

		if sectorWriteState.RawDataPosition == 0x56+0x100 {
			// We have the full sector data
			physicalSector := lastReadAddress.sector
			logicalSector := dos33SectorInterleaving[physicalSector]

			// transform the data from disk bytes to 6-bytes and EOR it
			a := uint8(0)
			for i := 0; i < 0x56+0x100; i++ {
				b := sixTwoDecoding[sectorWriteState.RawData[i]]
				a ^= b
				sectorWriteState.RawData[i] = a
			}

			// Transform the 0x156 bytes into the final 0x100 bytes
			sectorData := sectorDataDecode(sectorWriteState.RawData[0:0x156])

			// Save the data to memory & recreate the raw sector data
			image.tracks[lastReadAddress.track].sectors[logicalSector].data = sectorData
			makeSectorData(lastReadAddress.track, physicalSector)

			resetsectorWriteState()
			imageIsDirty = true
		}
	}
}

func FlushImage() {
	if imageIsDirty {
		writeDiskImage()
	}
}

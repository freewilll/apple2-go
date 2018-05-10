package utils

import (
	"compress/gzip"
	"encoding/hex"
	"io/ioutil"
	"os"
)

func ReadMemoryFromGzipFile(filename string) (data []byte, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err = ioutil.ReadAll(reader)

	return
}

func DecodeCmdLineAddress(s *string) (result *uint16) {
	if *s != "" {
		breakAddressValue, err := hex.DecodeString(*s)
		if err != nil {
			panic(err)
		}

		var value uint16
		if len(breakAddressValue) == 1 {
			value = uint16(breakAddressValue[0])
		} else if len(breakAddressValue) == 2 {
			value = uint16(breakAddressValue[0])*uint16(0x100) + uint16(breakAddressValue[1])
		} else {
			panic("Invalid break address")
		}
		result = &value
	}

	return result
}

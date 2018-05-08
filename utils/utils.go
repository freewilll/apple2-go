package utils

import (
	"compress/gzip"
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

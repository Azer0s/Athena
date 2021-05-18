package util

import (
	"encoding/binary"
	"io"
)

func ReadBytesLen(f io.Reader) ([]byte, error) {
	valLen, err := ReadLen(f)
	if err != nil {
		return []byte{}, err
	}

	valBytes := make([]byte, valLen)
	_, err = f.Read(valBytes)
	if err != nil {
		return []byte{}, err
	}

	return valBytes, nil
}

func ReadLen(f io.Reader) (uint32, error) {
	valLenBytes := make([]byte, 4)
	_, err := f.Read(valLenBytes)
	if err != nil {
		return 0, err
	}

	valLen := binary.BigEndian.Uint32(valLenBytes)
	return valLen, nil
}

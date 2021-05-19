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

func ReadType(f io.Reader) ([]byte, error) {
	typeBytes := make([]byte, 2)
	_, err := f.Read(typeBytes)
	if err != nil {
		return []byte{}, err
	}

	return typeBytes, nil
}

func ReadLen(f io.Reader) (uint16, error) {
	valLenBytes := make([]byte, 2)
	_, err := f.Read(valLenBytes)
	if err != nil {
		return 0, err
	}

	valLen := binary.BigEndian.Uint16(valLenBytes)
	return valLen, nil
}

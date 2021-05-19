package util

import (
	"bytes"
	"encoding/binary"
)

func StrToLenBytes(str string, buff *bytes.Buffer) {
	//write the length of the string as 2 bytes
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(len(str)))
	buff.Write(lenBytes)

	//write the string as bytes
	buff.Write([]byte(str))
}

func BytesToLenBytes(val []byte, buff *bytes.Buffer) {
	//write the length of the value
	valLenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(valLenBytes, uint16(len(val)))
	buff.Write(valLenBytes)

	//write the value
	buff.Write(val)
}

package util

import (
	"bytes"
	"encoding/binary"
)

func StrToLenBytes(str string, buff *bytes.Buffer) {
	//write the length of the string as 4 bytes
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(str)))
	buff.Write(lenBytes)

	//write the string as bytes
	buff.Write([]byte(str))
}

func BytesToLenBytes(val []byte, buff *bytes.Buffer) {
	//write the length of the value
	valLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(valLenBytes, uint32(len(val)))
	buff.Write(valLenBytes)

	//write the value
	buff.Write(val)
}

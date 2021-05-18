package marshaller

import (
	"bytes"
	"encoding/binary"
	"github.com/Azer0s/athena/util"
	"reflect"
)

func init() {
	Register(stringSliceMarshaller{})
}

type stringSliceMarshaller struct {
}

func (s stringSliceMarshaller) Value(b []byte) (interface{}, error) {
	buff := bytes.NewBuffer(b)
	readLen, err := util.ReadLen(buff)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for i := uint32(0); i < readLen; i++ {
		strBytes, err := util.ReadBytesLen(buff)
		if err != nil {
			return nil, err
		}

		res = append(res, string(strBytes))
	}
	return res, nil
}

func (s stringSliceMarshaller) Bytes(i interface{}) ([]byte, error) {
	val := i.([]string)

	buff := &bytes.Buffer{}

	//write the length of the string as 4 bytes
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(val)))
	buff.Write(lenBytes)

	for _, s := range val {
		util.StrToLenBytes(s, buff)
	}

	return buff.Bytes(), nil
}

func (s stringSliceMarshaller) Type() reflect.Type {
	return reflect.TypeOf([]string{})
}

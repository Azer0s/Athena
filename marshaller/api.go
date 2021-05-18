package marshaller

import (
	"bytes"
	"encoding/binary"
	"github.com/Azer0s/athena/util"
	"reflect"
)

type AthenaMarshaller interface {
	Value([]byte) (interface{}, error)
	Bytes(interface{}) ([]byte, error)
	Type() reflect.Type
}

var marshallerMap = make(map[string]AthenaMarshaller)

func typeToString(p reflect.Type) string {
	if p.Name() == "" {
		return "/" + p.String()
	}

	return p.PkgPath() + "/" + p.Name()
}

func Register(marshaller AthenaMarshaller) {
	marshallerMap[typeToString(marshaller.Type())] = marshaller
}

func Marshal(value interface{}, buff *bytes.Buffer) error {
	t := typeToString(reflect.TypeOf(value))

	util.StrToLenBytes(t, buff)

	b, err := marshallerMap[t].Bytes(value)

	if err != nil {
		return err
	}

	util.BytesToLenBytes(b, buff)

	return nil
}

func Unmarshal(b []byte) (interface{}, error) {
	tLen := binary.BigEndian.Uint32(b[:4])
	t := b[4 : 4+tLen]

	valLen := binary.BigEndian.Uint32(b[4+tLen : 8+tLen])

	val := b[8+tLen : 8+tLen+valLen]

	return UnmarshalWithoutLen(t, val)
}

func UnmarshalWithoutLen(typeBytes []byte, valBytes []byte) (interface{}, error) {
	return marshallerMap[string(typeBytes)].Value(valBytes)
}

package marshaller

import (
	"reflect"
)

func init() {
	Register(boolMarshaller{})
}

type boolMarshaller struct {
}

func (b boolMarshaller) Value(bytes []byte) (interface{}, error) {
	return bytes[0] == 0x1, nil
}

func (b boolMarshaller) Bytes(i interface{}) ([]byte, error) {
	if i.(bool) {
		return []byte{0x1}, nil
	}

	return []byte{0x0}, nil
}

func (b boolMarshaller) Type() uint16 {
	return TypeIdBool
}

func (b boolMarshaller) ReflectType() reflect.Type {
	return reflect.TypeOf(true)
}

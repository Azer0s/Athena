package marshaller

import (
	"reflect"
)

func init() {
	Register(stringMarshaller{})
}

type stringMarshaller struct {
}

func (s stringMarshaller) Value(bytes []byte) (interface{}, error) {
	return string(bytes), nil
}

func (s stringMarshaller) Bytes(i interface{}) ([]byte, error) {
	return []byte(i.(string)), nil
}

func (s stringMarshaller) Type() reflect.Type {
	return reflect.TypeOf("")
}

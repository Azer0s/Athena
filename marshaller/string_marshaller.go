package marshaller

import "reflect"

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

func (s stringMarshaller) Type() uint16 {
	return TypeIdString
}

func (s stringMarshaller) ReflectType() reflect.Type {
	return reflect.TypeOf("")
}

package marshaller

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

func init() {
	Register(intMarshaller{})
}

type intMarshaller struct {
}

func (i intMarshaller) Value(bytes []byte) (interface{}, error) {
	val := binary.BigEndian.Uint32(bytes)
	return int(*(*int32)(unsafe.Pointer(&val))), nil
}

func (i intMarshaller) Bytes(v interface{}) ([]byte, error) {
	val := int32(v.(int))
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, *(*uint32)(unsafe.Pointer(&val)))
	return bytes, nil
}

func (i intMarshaller) Type() uint16 {
	return TypeIdInt
}

func (i intMarshaller) ReflectType() reflect.Type {
	return reflect.TypeOf(0)
}

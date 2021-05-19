package marshaller

import (
	"encoding/binary"
	"math"
	"reflect"
	"unsafe"
)

func init() {
	Register(float32Marshaller{})
}

type float32Marshaller struct {
}

func (f float32Marshaller) Value(bytes []byte) (interface{}, error) {
	val := binary.BigEndian.Uint32(bytes)
	u := *(*float32)(unsafe.Pointer(&val))

	return u, nil
}

func (f float32Marshaller) Bytes(i interface{}) ([]byte, error) {
	val := i.(float32)

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(val))

	return buf[:], nil
}

func (f float32Marshaller) Type() uint16 {
	return TypeIdFloat32
}

func (f float32Marshaller) ReflectType() reflect.Type {
	return reflect.TypeOf(float32(0))
}

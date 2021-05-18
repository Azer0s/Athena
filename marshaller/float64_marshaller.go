package marshaller

import (
	"encoding/binary"
	"math"
	"reflect"
	"unsafe"
)

func init() {
	Register(float64Marshaller{})
}

type float64Marshaller struct {
}

func (f float64Marshaller) Value(bytes []byte) (interface{}, error) {
	val := binary.BigEndian.Uint64(bytes)
	u := *(*float64)(unsafe.Pointer(&val))

	return u, nil
}

func (f float64Marshaller) Bytes(i interface{}) ([]byte, error) {
	val := i.(float64)

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(val))

	return buf[:], nil
}

func (f float64Marshaller) Type() reflect.Type {
	return reflect.TypeOf(float64(0))
}

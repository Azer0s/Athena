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
	Type() uint16
	ReflectType() reflect.Type
}

var marshallerMap = make(map[uint16]AthenaMarshaller)
var marshallerMapStrings = make(map[string]AthenaMarshaller)

func typeToString(p reflect.Type) string {
	if p.PkgPath() == "" {
		return "athena/" + p.Name()
	}

	return p.PkgPath() + "/" + p.Name()
}

func reflectValueToInterface(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.Elem().Interface()
	default:
		return v.Interface()
	}
}

func Register(marshaller AthenaMarshaller) {
	marshallerMap[marshaller.Type()] = marshaller
	marshallerMapStrings[typeToString(marshaller.ReflectType())] = marshaller
}

func marshalSlice(val reflect.Value, buff *bytes.Buffer) error {
	buff.Write(TypeIdSlice)

	valBuff := &bytes.Buffer{}

	//write the length of the slice as 2 bytes
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(val.Len()))
	valBuff.Write(lenBytes)

	for i := 0; i < val.Len(); i++ {
		err := Marshal(reflectValueToInterface(val.Index(i)), valBuff)
		if err != nil {
			return err
		}
	}

	util.BytesToLenBytes(valBuff.Bytes(), buff)

	return nil
}

func marshalMap(val reflect.Value, buff *bytes.Buffer) error {
	buff.Write(TypeIdMap)

	valBuff := &bytes.Buffer{}

	keys := val.MapKeys()

	//write the length of the map as 2 bytes
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(len(keys)))
	valBuff.Write(lenBytes)

	for _, key := range keys {
		err := Marshal(reflectValueToInterface(key), valBuff)
		if err != nil {
			return err
		}

		err = Marshal(reflectValueToInterface(val.MapIndex(key)), valBuff)
		if err != nil {
			return err
		}
	}

	util.BytesToLenBytes(valBuff.Bytes(), buff)

	return nil
}

func Marshal(value interface{}, buff *bytes.Buffer) error {
	typ := reflect.ValueOf(value)
	switch typ.Kind() {
	case reflect.Slice:
		err := marshalSlice(typ, buff)
		if err != nil {
			return err
		}

	case reflect.Map:
		err := marshalMap(typ, buff)
		if err != nil {
			return err
		}

	default:
		t := typeToString(reflect.TypeOf(value))
		marshaller := marshallerMapStrings[t]

		typeBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(typeBytes, marshaller.Type())
		buff.Write(typeBytes)

		b, err := marshaller.Bytes(value)
		if err != nil {
			return err
		}
		util.BytesToLenBytes(b, buff)
	}

	return nil
}

func unmarshalSlice(valBytes []byte) ([]interface{}, error) {
	buff := bytes.NewBuffer(valBytes)

	readLen, err := util.ReadLen(buff)
	if err != nil {
		return nil, err
	}

	res := []interface{}{}

	for i := uint16(0); i < readLen; i++ {
		typeBytes, err := util.ReadType(buff)
		if err != nil {
			return nil, err
		}

		vBytes, err := util.ReadBytesLen(buff)
		if err != nil {
			return nil, err
		}

		val, err := UnmarshalWithoutLen(typeBytes, vBytes)
		if err != nil {
			return nil, err
		}

		res = append(res, val)
	}

	return res, nil
}

func unmarshalMap(valBytes []byte) (map[interface{}]interface{}, error) {
	buff := bytes.NewBuffer(valBytes)

	readLen, err := util.ReadLen(buff)
	if err != nil {
		return nil, err
	}

	res := make(map[interface{}]interface{})

	for i := uint16(0); i < readLen; i++ {
		typeBytes, err := util.ReadType(buff)
		if err != nil {
			return nil, err
		}

		kValBytes, err := util.ReadBytesLen(buff)
		if err != nil {
			return nil, err
		}

		kVal, err := UnmarshalWithoutLen(typeBytes, kValBytes)
		if err != nil {
			return nil, err
		}

		typeBytes, err = util.ReadType(buff)
		if err != nil {
			return nil, err
		}

		vValBytes, err := util.ReadBytesLen(buff)
		if err != nil {
			return nil, err
		}

		vVal, err := UnmarshalWithoutLen(typeBytes, vValBytes)
		if err != nil {
			return nil, err
		}

		res[kVal] = vVal
	}

	return res, nil
}

func Unmarshal(b []byte) (interface{}, error) {
	val := b[4:]
	return UnmarshalWithoutLen(b[:2], val)
}

func UnmarshalWithoutLen(typeBytes []byte, valBytes []byte) (interface{}, error) {
	if typeBytes[0] == 0xFF {
		switch typeBytes[1] {
		case 0xFF:
			return unmarshalSlice(valBytes)
		case 0xFE:
			return unmarshalMap(valBytes)
		}
	}

	return marshallerMap[binary.BigEndian.Uint16(typeBytes)].Value(valBytes)
}

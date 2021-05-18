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
	marshallerMap[typeToString(marshaller.Type())] = marshaller
}

func marshalSlice(val reflect.Value, buff *bytes.Buffer) error {
	util.StrToLenBytes("athena/[]", buff)

	valBuff := &bytes.Buffer{}

	//write the length of the string as 4 bytes
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(val.Len()))
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
	util.StrToLenBytes("athena/{}", buff)

	valBuff := &bytes.Buffer{}

	keys := val.MapKeys()

	//write the length of the string as 4 bytes
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(keys)))
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
		util.StrToLenBytes(t, buff)

		b, err := marshallerMap[t].Bytes(value)
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

	res := make([]interface{}, 0)

	for i := uint32(0); i < readLen; i++ {
		typeBytes, err := util.ReadBytesLen(buff)
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

	for i := uint32(0); i < readLen; i++ {
		typeBytes, err := util.ReadBytesLen(buff)
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

		typeBytes, err = util.ReadBytesLen(buff)
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
	tLen := binary.BigEndian.Uint32(b[:4])
	t := b[4 : 4+tLen]

	valLen := binary.BigEndian.Uint32(b[4+tLen : 8+tLen])

	val := b[8+tLen : 8+tLen+valLen]

	return UnmarshalWithoutLen(t, val)
}

func UnmarshalWithoutLen(typeBytes []byte, valBytes []byte) (interface{}, error) {
	switch string(typeBytes) {
	case "athena/[]":
		return unmarshalSlice(valBytes)
	case "athena/{}":
		return unmarshalMap(valBytes)
	default:
		return marshallerMap[string(typeBytes)].Value(valBytes)
	}
}

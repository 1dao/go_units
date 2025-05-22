package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// 定义枚举类型
const (
	TypeInt8 = iota + 1
	TypeByte = iota + 1
	TypeInt16
	TypeInt32
	TypeInt64
	TypeInt
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeUint
	TypeFloat32
	TypeFloat64
	TypeBool
	TypeString
	TypeBytes
)

// 打包并返回结果
func GPackerPack(values ...interface{}) []byte {
	buf := new(bytes.Buffer)

	for _, v := range values {
		switch v := v.(type) {
		case int8:
			binary.Write(buf, binary.BigEndian, uint8(TypeInt8))
			binary.Write(buf, binary.BigEndian, v)
		case int16:
			binary.Write(buf, binary.BigEndian, uint8(TypeInt16))
			binary.Write(buf, binary.BigEndian, v)
		case int32:
			binary.Write(buf, binary.BigEndian, uint8(TypeInt32))
			binary.Write(buf, binary.BigEndian, v)
		case int64:
			binary.Write(buf, binary.BigEndian, uint8(TypeInt64))
			binary.Write(buf, binary.BigEndian, v)
		case uint8:
			binary.Write(buf, binary.BigEndian, uint8(TypeUint8))
			binary.Write(buf, binary.BigEndian, v)
		case uint16:
			binary.Write(buf, binary.BigEndian, uint8(TypeUint16))
			binary.Write(buf, binary.BigEndian, v)
		case uint32:
			binary.Write(buf, binary.BigEndian, uint8(TypeUint32))
			binary.Write(buf, binary.BigEndian, v)
		case uint64:
			binary.Write(buf, binary.BigEndian, uint8(TypeUint64))
			binary.Write(buf, binary.BigEndian, v)
		case float32:
			binary.Write(buf, binary.BigEndian, uint8(TypeFloat32))
			binary.Write(buf, binary.BigEndian, v)
		case float64:
			binary.Write(buf, binary.BigEndian, uint8(TypeFloat64))
			binary.Write(buf, binary.BigEndian, v)
		case int:
			binary.Write(buf, binary.BigEndian, uint8(TypeInt))
			binary.Write(buf, binary.BigEndian, int32(v))
		case uint:
			binary.Write(buf, binary.BigEndian, uint8(TypeUint))
			binary.Write(buf, binary.BigEndian, uint32(v))
		case bool:
			binary.Write(buf, binary.BigEndian, uint8(TypeBool))
			binary.Write(buf, binary.BigEndian, v)
		case string:
			binary.Write(buf, binary.BigEndian, uint8(TypeString))
			strBytes := []byte(v)
			binary.Write(buf, binary.BigEndian, uint16(len(strBytes)))
			binary.Write(buf, binary.BigEndian, strBytes)
		case []byte:
			binary.Write(buf, binary.BigEndian, uint8(TypeBytes))
			binary.Write(buf, binary.BigEndian, uint16(len(v)))
			binary.Write(buf, binary.BigEndian, v)
		default:
			fmt.Printf("Unknown type: %T, value: %v\n", v, v)
		}
	}

	return buf.Bytes()
}

func rawUnpack(data []byte) ([]interface{}, error) {
	buf := bytes.NewReader(data)
	var results []interface{}

	for buf.Len() > 0 {
		var typeID uint8
		if err := binary.Read(buf, binary.BigEndian, &typeID); err != nil {
			return nil, err
		}

		switch typeID {
		case TypeInt8:
			var v int8
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeInt16:
			var v int16
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeInt32:
			var v int32
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeInt64:
			var v int64
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeUint8:
			var v uint8
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeUint16:
			var v uint16
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeUint32:
			var v uint32
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeUint64:
			var v uint64
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeInt:
			var v int32
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeUint:
			var v uint32
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeFloat32:
			var v float32
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeFloat64:
			var v float64
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeBool:
			var v bool
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nil, err
			}
			results = append(results, v)
		case TypeString:
			var strLen uint16
			if err := binary.Read(buf, binary.BigEndian, &strLen); err != nil {
				return nil, err
			}
			str := make([]byte, strLen)
			if err := binary.Read(buf, binary.BigEndian, &str); err != nil {
				return nil, err
			}
			results = append(results, string(str))
		case TypeBytes:
			var byteLen uint16
			if err := binary.Read(buf, binary.BigEndian, &byteLen); err != nil {
				return nil, err
			}
			bytes := make([]byte, byteLen)
			if err := binary.Read(buf, binary.BigEndian, &bytes); err != nil {
				return nil, err
			}
			results = append(results, bytes)
		default:
			return nil, fmt.Errorf("unknown type ID: %d", typeID)
		}
	}

	return results, nil
}

// 解包并返回结果
func GPackerUnPack(data []byte) ([]interface{}, error) {
	return rawUnpack(data)
}

func GPackerUnPackArgs(data []byte) ([]reflect.Value, error) {
	if data == nil {
		return nil, fmt.Errorf("no data to unpack")
	}
	if len(data) == 0 {
		return nil, nil
	}

	res, err := rawUnpack(data)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("no data to unpack")
	}

	ts := ""
	fmt.Println("res:", ts)
	args := make([]reflect.Value, len(res))
	for i, v := range res {
		args[i] = reflect.ValueOf(v)
	}
	return args, nil
}

// func main() {
// 	fmt.Println("Packing data...", 100, 200, "hello", []byte{1, 2, 3})
// 	data := GPackerPack(100, int16(200), "hello", []byte{1, 2, 3})
// 	fmt.Println("Packed data:", data)

// 	fmt.Println("Unpacking data...")
// 	vars, err := GPackerUnPack(data)
// 	if err != nil {
// 		fmt.Println("Error unpacking data:", err)
// 		return
// 	}
// 	fmt.Println("Unpacked variables:", vars)
// }

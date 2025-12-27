package serializer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// 支持的数据类型
const (
	TypeInt = iota
	TypeFloat32
	TypeFloat64
	TypeString
	TypeBool
)

// MAGIC 魔数
var MAGIC = [...]byte{23, 37, 111, 51}

type MySerializer struct {
}

// MarshalArguments 将一批参数序列化
func MarshalArguments(arguments ...any) ([]byte, error) {
	types := make([]byte, 0, len(arguments))  // 每个参数的类型, 一个 byte 可以放下
	lengths := make([]int, 0, len(arguments)) // 每个参数的长度, 用 int, 后面转 byte
	buf := make([]byte, 0, 256)
	buffer := bytes.NewBuffer(buf) // 将参数往里写

	// 进行序列化
	for i, argument := range arguments {
		// 序列化单个参数
		argumentType, argumentLength, err := marshalArgument(argument, buffer)
		if err != nil {
			return nil, fmt.Errorf("序列化第 %d 个参数失败\n, error:%s", i, err.Error())
		}

		// 将当前参数的信息追加
		types = append(types, argumentType)
		lengths = append(lengths, argumentLength)
	}

	/*
		总长度 =
		魔数长度 (每个魔数用 1 个 Byte 即可表示) +
		参数个数 (用 8 个 Byte 即可表示) +
		每个参数的类型 (每个类型用 1 个 Byte 即可表示) +
		每个参数的长度 (每个长度用 8 个 Byte 即可表示) +
		具体参数内容的长度
	*/
	resultLength := len(MAGIC) + 8 + len(types) + 8*len(types) + buffer.Len()
	result := make([]byte, 0, resultLength)
	resultBuffer := bytes.NewBuffer(result)

	// 1. 写魔数
	resultBuffer.Write(MAGIC[:]) // 转切片

	// 2. 写参数个数
	resultBuffer.Write(IntToBytes(len(types)))

	// 3. 写每个参数的类型
	resultBuffer.Write(types)

	// 4. 写每个参数的长度
	for _, length := range lengths {
		resultBuffer.Write(IntToBytes(length))
	}

	// 5. 写具体参数内容
	resultBuffer.Write(buffer.Bytes())

	// 返回
	return resultBuffer.Bytes(), nil
}

// UnmarshalArguments 反序列化一个数据流
func UnmarshalArguments(bs []byte) ([]any, error) {
	pos := 0 // 当前指向

	// 参数校验
	if len(bs) < len(MAGIC) {
		return nil, errors.New("数据流长度过短")
	}

	// 1. 检查魔数
	if !bytes.Equal(MAGIC[:], bs[pos:len(MAGIC)]) {
		return nil, errors.New("魔数校验失败")
	}
	pos += len(MAGIC)

	// 参数校验
	if len(bs) < len(MAGIC)+8 {
		return nil, errors.New("数据流长度过短")
	}

	// 2. 获取参数个数
	cnt := BytesToInt(bs[pos : pos+8])
	pos += 8
	if cnt <= 0 {
		return nil, nil
	}

	// 参数校验
	if len(bs) < len(MAGIC)+8+cnt+8*cnt {
		return nil, errors.New("数据流长度过短")
	}

	// 3. 获取每个参数类型
	types := make([]byte, 0, cnt)
	for i := 0; i < cnt; i++ {
		types = append(types, bs[pos])
		pos++
	}

	// 4. 获取每个参数的长度
	totalLength := 0 // 参数总长度
	lengths := make([]int, 0, cnt)
	for i := 0; i < cnt; i++ {
		length := BytesToInt(bs[pos : pos+8])
		lengths = append(lengths, length)
		totalLength += length
		pos += 8
	}

	// 参数校验
	if len(bs[pos:]) < totalLength {
		return nil, errors.New("数据流长度过短")
	}

	// 5. 获取每个参数
	arguments := make([]any, 0, cnt)

	for i := 0; i < cnt; i++ {
		argumentType := types[i]                   // 当前参数类型
		argumentLength := lengths[i]               // 当前参数长度
		argumentBS := bs[pos : pos+argumentLength] // 当前参数 Byte 切片

		argument, err := unmarshallArgument(argumentType, argumentBS)
		if err != nil {
			return nil, fmt.Errorf("第 %d 个参数反序列化失败, error:%s", i, err.Error())
		}
		// 更新 pos
		pos += argumentLength
		arguments = append(arguments, argument)
	}

	return arguments, nil
}

// IntToBytes Int 转 []Byte
func IntToBytes(n int) []byte {
	x := int64(n)
	buffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(buffer, binary.BigEndian, x)
	return buffer.Bytes()
}

// BytesToInt []Byte 转 Int
func BytesToInt(bs []byte) int {
	buffer := bytes.NewBuffer(bs)
	var x int64
	_ = binary.Read(buffer, binary.BigEndian, &x)
	return int(x)
}

// 序列化单个参数
func marshalArgument(argument any, buffer *bytes.Buffer) (byte, int, error) {
	switch v := argument.(type) {
	case int:
		err := binary.Write(buffer, binary.BigEndian, int64(v))
		if err != nil {
			return 0, 0, errors.New(fmt.Sprintf("序列化失败，数据类型:%s", reflect.TypeOf(v).Name()))
		}
		return TypeInt, 8, nil // 统一转成 8 个字节
	case float32:
		err := binary.Write(buffer, binary.BigEndian, v)
		if err != nil {
			return 0, 0, errors.New(fmt.Sprintf("序列化失败，数据类型:%s", reflect.TypeOf(v).Name()))
		}
		return TypeFloat32, 4, nil
	case float64:
		err := binary.Write(buffer, binary.BigEndian, v)
		if err != nil {
			return 0, 0, errors.New(fmt.Sprintf("序列化失败，数据类型:%s", reflect.TypeOf(v).Name()))
		}
		return TypeFloat64, 8, nil
	case string:
		_, err := buffer.WriteString(v)
		if err != nil {
			return 0, 0, errors.New(fmt.Sprintf("序列化失败，数据类型:%s", reflect.TypeOf(v).Name()))
		}
		return TypeString, len(v), nil // 断定 String 的最大长度不超过 MaxInt
	case bool:
		err := binary.Write(buffer, binary.BigEndian, v)
		if err != nil {
			return 0, 0, errors.New(fmt.Sprintf("序列化失败，数据类型:%s", reflect.TypeOf(v).Name()))
		}
		return TypeBool, 1, nil
	default:
		return 0, 0, errors.New("当前类型不支持序列化")
	}
}

// 反序列化单个参数
func unmarshallArgument(argumentType byte, argumentBS []byte) (any, error) {
	switch argumentType {
	case TypeInt:
		return BytesToInt(argumentBS), nil
	case TypeFloat32:
		var x float32
		buffer := bytes.NewBuffer(argumentBS)
		err := binary.Read(buffer, binary.BigEndian, &x)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("参数序列化失败，类型为 Float32"))
		}
		return x, nil
	case TypeFloat64:
		var x float64
		buffer := bytes.NewBuffer(argumentBS)
		err := binary.Read(buffer, binary.BigEndian, &x)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("参数序列化失败，类型为 Float64"))
		}
		return x, nil
	case TypeString:
		return string(argumentBS), nil
	case TypeBool:
		var x bool
		buffer := bytes.NewBuffer(argumentBS)
		err := binary.Read(buffer, binary.BigEndian, &x)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("参数序列化失败，类型为 Bool"))
		}
		return x, nil
	default:
		return nil, errors.New(fmt.Sprintf("参数序列化失败，不支持该类型"))
	}
}

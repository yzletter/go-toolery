package serializer

// Serializer 定义接口
type Serializer interface {
	Marshal(object any) ([]byte, error)
	Unmarshal(buffer []byte, object any) error
}

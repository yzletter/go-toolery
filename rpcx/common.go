package rpcx

type RpcxData struct {
	A  int
	B  float32
	C  bool
	D  float64
	E  string
	f  int // 可导出的字段才会被序列化，该字段不可导出
	Id string
}

package main

import "github.com/yzletter/go-toolery/rpcx/serializer"

func main() {
	var s = serializer.MySerializer{}
	server := NewServer(5678, s)
	server.Serve()
}

package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/yzletter/go-toolery/rpcx"
	"github.com/yzletter/go-toolery/rpcx/serializer"
)

type EchoServer struct {
	conn *net.UDPConn          // 面向报文
	s    serializer.Serializer // 序列化接口
}

// NewServer 构造函数
func NewServer(port int, s serializer.Serializer) *EchoServer {
	// UDP 地址
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
		return nil
	}

	// 建立连接
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("服务端启动成功")
	return &EchoServer{
		conn: conn,
		s:    s,
	}
}

func (server EchoServer) Serve() {
	defer server.conn.Close()
	for {
		request := make([]byte, 4096) // 设置最大长度, 防止 flood attack, 注意不要写成 make([]byte, 0, 4096)

		// 不断读取请求
		n, remoteAddr, err := server.conn.ReadFromUDP(request)
		if err != nil {
			log.Println(err)
			return
		}

		// 并发进行业务处理
		go server.handle(request[:n], remoteAddr)
	}

}

// 业务处理
func (server EchoServer) handle(request []byte, remoteAddr *net.UDPAddr) {
	// 提示信息
	log.Printf("接收到来自 %s 的调用请求\n", remoteAddr.String())

	// 反序列化请求
	var target rpcx.RpcxData
	err := server.s.Unmarshal(request, &target)
	if err != nil {
		log.Println(err)
		return
	}

	// 模拟不同业务的处理时长不同
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))

	// 这里用简单的业务逻辑 直接把请求还发给对方
	stream, err := server.s.Marshal(target)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = server.conn.WriteToUDP(stream, remoteAddr)
	if err != nil {
		log.Println(err)
		return
	}

	// 提示信息
	log.Printf("已将请求结果发给 %s \n", remoteAddr.String())
}

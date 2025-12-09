package udp_logger

import (
	"fmt"
	"net"
)

type LogProducer struct {
	conn    net.Conn
	buffer  chan string
	flushed chan struct{} // 用于 Close
}

func NewLogProducer(collectorAddr string, bufferSize int) (*LogProducer, error) {
	if conn, err := net.Dial("udp", collectorAddr); err == nil {
		producer := &LogProducer{
			conn:    conn,
			buffer:  make(chan string, bufferSize),
			flushed: make(chan struct{}, 1),
		}

		go producer.daemonSend()
		return producer, nil
	} else {
		fmt.Println("NewLogProducer failed")
		return nil, err
	}
}

func (p *LogProducer) daemonSend() {
	for {
		if log, ok := <-p.buffer; ok {
			_, err := p.conn.Write([]byte(log))
			if err != nil {
				fmt.Printf("producer write log : %s failed err : %s", log, err)
			}
		} else {
			p.flushed <- struct{}{} // buffer 已经关闭且清空
			break
		}
	}
}

// Send 在往本地文件写 log 的时候同时向 Producer 写一份
func (p *LogProducer) Send(log string) (err error) {
	defer func() {
		if obj := recover(); obj != nil {
			err = fmt.Errorf("%v", obj)
		}
	}()
	p.buffer <- log
	return
}

func (p *LogProducer) Close() {
	close(p.buffer) // 关闭 buffer 管道, 不允许再往里面发送内容
	<-p.flushed     // 等待 buffer 管道被清空
	p.conn.Close()  // 关闭 UDP 连接
}

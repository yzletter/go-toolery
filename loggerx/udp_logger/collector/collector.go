package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type LogCollector struct {
	conn   *net.UDPConn  // UDP 连接
	fout   *os.File      // 聚合的目标文件
	writer *bufio.Writer // 带缓存的 Writer
}

func NewLogCollector(port int, file string) (*LogCollector, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	fout, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(fout)

	collector := &LogCollector{
		conn:   conn,
		fout:   fout,
		writer: writer,
	}

	// 每秒刷新一次缓冲区
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			collector.writer.Flush()
		}
	}()

	return collector, nil
}

func (c *LogCollector) Receive() {
	content := make([]byte, 4<<20)
	for {
		n, address, err := c.conn.ReadFromUDP(content)
		if err != nil {
			fmt.Printf("receive log failed: %v\n", err)
		} else {
			// 给每条日志前加上来源地址和端口号
			log := fmt.Sprintf("[%s:%d] %s", address.IP.String(), address.Port, string(content[:n]))
			if _, err := c.writer.Write([]byte(log)); err != nil {
				fmt.Printf("write log <%s> to file fail: %v\n", log, err)
			} else {
				c.writer.WriteString("\n")

			}
		}
	}
}

func (c *LogCollector) Close() {
	c.conn.Close()
	c.writer.Flush()
	c.fout.Close()
}

func listenTerminal(collector *LogCollector) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) // 注册信号 2 和 15

	sig := <-c
	fmt.Printf("receive signal %s\n", sig.String())

	if collector != nil {
		collector.Close()
	}

	os.Exit(0)
}

func main() {
	port := flag.Int("port", 4321, "udp server port")
	logFile := flag.String("log", "./loggerx/log/udp_collect", "log sink")
	flag.Parse()
	collector, err := NewLogCollector(*port, *logFile+".log")
	if err != nil {
		panic(err)
	}
	go listenTerminal(collector)

	collector.Receive()
}

// go run ./logger/udp_logger/collector -port=4321 -log=log/collect.log

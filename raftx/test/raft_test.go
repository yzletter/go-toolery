package test

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"testing"
	"time"

	"github.com/yzletter/go-toolery/raftx"
)

func init() {
	fout, err := os.OpenFile("../raftx.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	handler := slog.NewTextHandler(fout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func TestRaft(t *testing.T) {
	size := 7 //集群大小
	servers := make([]*raftx.Server, 0, size)

	// 启动各台 Raft server
	for i := 0; i < size; i++ {
		port := 7000 + i
		transporter := raftx.NewHttpTransporter("/raftx", 2*time.Second)
		server := raftx.NewServer(fmt.Sprintf("http://127.0.0.1:%d", port), port, nil, transporter)
		server.Start(false)
		defer server.Stop()
		servers = append(servers, server)
	}

	// 给每个 server 添加 peer
	for _, server := range servers {
		for _, peer := range servers {
			server.AddPeer(&raftx.Peer{ID: peer.ID, ConnectionString: peer.ConnectionString})
		}
	}

	// 随机停掉一台server
	downServer := servers[rand.IntN(size)]
	downServer.Stop()

	// 等待 Leader 产生
	var leader *raftx.Server
	for {
		leader = getLeader(servers)
		if leader != nil {
			fmt.Printf("Leader Elevated %s\n", leader.GetID()) // 打印 Leader 是谁
			break
		}
		time.Sleep(1 * time.Second) //等待1秒，进入下一轮循环
	}

	// 产生 300 条 NoopCommand
	for i := 0; i < 3*raftx.MaxLogEntriesPerHeartBeat; i++ {
		leader.Do(raftx.NoopCommand{})
		time.Sleep(50 * time.Millisecond)
	}

	// 重启down server，它需要接收 300 条 NoopCommand，一次心跳发不完
	downServer.Start(true)

	// 等所有日志复制完成
	time.Sleep(5 * time.Second)

	// 最终，集群中所有节点的LastLogIndex应该都是3*raftx.MaxLogEntriesPerHeartbeat
}

// 并不是每个人心目中的leader都是一个，得找到票数最多的leader
func getLeader(servers []*raftx.Server) *raftx.Server {
	if len(servers) == 0 {
		return nil
	}
	countMap := make(map[string]int, len(servers))
	for _, server := range servers {
		countMap[server.LeaderID()] = countMap[server.LeaderID()] + 1
	}
	var leaderID string
	var maxVotes int = 1
	for k, v := range countMap {
		slog.Debug("vote count", k, v)
		if v > maxVotes {
			maxVotes = v
			leaderID = k
		}
	}
	for _, server := range servers {
		if leaderID == server.GetID() {
			return server
		}
	}
	return nil
}

// go test -v ./raftx/test -run=^TestRaft$ -count=1 -timeout=5m
/*
yzletter@192 go-toolery % go test -v ./raftx/test -run=^TestRaft$ -count=1 -timeout=5m
=== RUN   TestRaft
Leader Elevated d5fsrom1n8d6peii27dg
Raftx Cluster Information Leader d5fsrom1n8d6peii27dg LastLogIndex 40 CommitIndex 40
Raftx Cluster Information Follower d5fsrom1n8d6peii27cg LastLogIndex 40
Raftx Cluster Information Follower d5fsrom1n8d6peii27d0 LastLogIndex 40
Raftx Cluster Information Follower d5fsrom1n8d6peii27e0 LastLogIndex 40
Raftx Cluster Information Follower d5fsrom1n8d6peii27eg LastLogIndex 40
Raftx Cluster Information Follower d5fsrom1n8d6peii27f0 LastLogIndex 0
Raftx Cluster Information Follower d5fsrom1n8d6peii27fg LastLogIndex 40
---------------------------------------
Raftx Cluster Information Leader d5fsrom1n8d6peii27dg LastLogIndex 99 CommitIndex 99
Raftx Cluster Information Follower d5fsrom1n8d6peii27cg LastLogIndex 99
Raftx Cluster Information Follower d5fsrom1n8d6peii27d0 LastLogIndex 99
Raftx Cluster Information Follower d5fsrom1n8d6peii27e0 LastLogIndex 99
Raftx Cluster Information Follower d5fsrom1n8d6peii27eg LastLogIndex 99
Raftx Cluster Information Follower d5fsrom1n8d6peii27f0 LastLogIndex 0
Raftx Cluster Information Follower d5fsrom1n8d6peii27fg LastLogIndex 99
---------------------------------------
Raftx Cluster Information Leader d5fsrom1n8d6peii27dg LastLogIndex 158 CommitIndex 158
Raftx Cluster Information Follower d5fsrom1n8d6peii27cg LastLogIndex 158
Raftx Cluster Information Follower d5fsrom1n8d6peii27d0 LastLogIndex 158
Raftx Cluster Information Follower d5fsrom1n8d6peii27e0 LastLogIndex 158
Raftx Cluster Information Follower d5fsrom1n8d6peii27eg LastLogIndex 158
Raftx Cluster Information Follower d5fsrom1n8d6peii27f0 LastLogIndex 0
Raftx Cluster Information Follower d5fsrom1n8d6peii27fg LastLogIndex 158
---------------------------------------
Raftx Cluster Information Leader d5fsrom1n8d6peii27dg LastLogIndex 217 CommitIndex 217
Raftx Cluster Information Follower d5fsrom1n8d6peii27cg LastLogIndex 217
Raftx Cluster Information Follower d5fsrom1n8d6peii27d0 LastLogIndex 217
Raftx Cluster Information Follower d5fsrom1n8d6peii27e0 LastLogIndex 217
Raftx Cluster Information Follower d5fsrom1n8d6peii27eg LastLogIndex 217
Raftx Cluster Information Follower d5fsrom1n8d6peii27f0 LastLogIndex 0
Raftx Cluster Information Follower d5fsrom1n8d6peii27fg LastLogIndex 217
---------------------------------------
Raftx Cluster Information Leader d5fsrom1n8d6peii27dg LastLogIndex 276 CommitIndex 276
Raftx Cluster Information Follower d5fsrom1n8d6peii27cg LastLogIndex 276
Raftx Cluster Information Follower d5fsrom1n8d6peii27d0 LastLogIndex 276
Raftx Cluster Information Follower d5fsrom1n8d6peii27e0 LastLogIndex 276
Raftx Cluster Information Follower d5fsrom1n8d6peii27eg LastLogIndex 276
Raftx Cluster Information Follower d5fsrom1n8d6peii27f0 LastLogIndex 0
Raftx Cluster Information Follower d5fsrom1n8d6peii27fg LastLogIndex 276
---------------------------------------
Raftx Cluster Information Leader d5fsrom1n8d6peii27dg LastLogIndex 300 CommitIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27cg LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27d0 LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27e0 LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27eg LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27f0 LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27fg LastLogIndex 300
---------------------------------------
Raftx Cluster Information Leader d5fsrom1n8d6peii27dg LastLogIndex 300 CommitIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27cg LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27d0 LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27e0 LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27eg LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27f0 LastLogIndex 300
Raftx Cluster Information Follower d5fsrom1n8d6peii27fg LastLogIndex 300
---------------------------------------
--- PASS: TestRaft (21.23s)
PASS
ok      github.com/yzletter/go-toolery/raftx/test       21.648s
*/

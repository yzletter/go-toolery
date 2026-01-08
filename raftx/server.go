package raftx

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/yzletter/go-toolery/errx"
)

type Peer struct {
	ID               string
	ConnectionString string // IP 和端口号
}

// Server Raft 集群节点
type Server struct {
	sync.RWMutex
	Peer
	port         int
	term         int64
	leaderID     string
	votedFor     string
	state        State
	log          *Log
	peers        []*Peer // 集群中其他所有节点, 进行 RPC 通信
	prevLogIndex map[string]int64
	fsm          FSM
	transporter  Transporter
	shutdownCh   chan struct{}
	rpcCh        chan RPC
	routineGroup sync.WaitGroup
}

func NewServer(connString string, port int, fsm FSM, transporter Transporter) *Server {
	server := &Server{
		port:         port,
		peers:        make([]*Peer, 0, 8),
		prevLogIndex: make(map[string]int64, 8),
		fsm:          fsm,
		transporter:  transporter,
		shutdownCh:   make(chan struct{}, 1),
		rpcCh:        make(chan RPC, 100),
	}
	server.log = NewLog(server)
	server.ID = xid.New().String()
	server.ConnectionString = connString
	return server
}

// LeaderID 当前leader的id
func (server *Server) LeaderID() string {
	return server.leaderID
}

func (server *Server) GetID() string {
	return server.ID
}

// QuorumSize 返回超过一半的数量是多少
func (server *Server) QuorumSize() int {
	return (len(server.peers)+1)/2 + 1
}

func (server *Server) GetState() State {
	server.RLock()
	defer server.RUnlock()
	return server.state
}

func (server *Server) AddPeer(peer *Peer) {
	if peer == nil {
		return
	}

	if len(peer.ID) == 0 || len(peer.ConnectionString) == 0 {
		return
	}

	if peer.ID == server.ID { // 排查自己
		return
	}

	server.peers = append(server.peers, peer)
}

func (server *Server) upgradeTerm(term int64, leaderID string) {
	server.term = term
	server.votedFor = ""
	server.leaderID = leaderID
	server.SetState(Follower) // 降级
}

func (server *Server) SetState(state State) {
	server.Lock()
	defer server.Unlock()

	server.state = state
}

func (server *Server) Start(restart bool) {
	// 判断是否是重启, 只有第一次启动时需要启动 HTTP Server
	if !restart {
		go server.transporter.Start(server.port, server)
	} else {
		// 重启只需要初始化相关成员即可
		server.state = Follower
		server.shutdownCh = make(chan struct{}, 1)
		server.rpcCh = make(chan RPC, 100)
		server.routineGroup = sync.WaitGroup{}
	}

	go server.print()

	server.routineGroup.Add(1)
	go func() {
		defer server.routineGroup.Done()
		for server.GetState() != Stopped {
			// 根据不同状态进行不同循环
			switch server.GetState() {
			case Follower:
				server.FollowerLoop()
			case Candidate:
				server.CandidateLoop()
			case Leader:
				server.LeaderLoop()
			default:
				panic("unhandled default case")
			}
		}
	}()
}

func (server *Server) Stop() {
	// 已经停止了
	if server.GetState() == Stopped {
		return
	}

	server.SetState(Stopped)
	server.shutdownCh <- struct{}{}

	close(server.shutdownCh)
	close(server.rpcCh)

	// 等待所有异步任务
	server.routineGroup.Wait()
	slog.Info("Server Shut Down", "id", server.ID)
}

func (server *Server) FollowerLoop() {
	slog.Info("Run As Follower", "id", server.ID)
	electionTimer := randomTimeout(ElectionTimeout) // 开始选举倒计时

	for server.GetState() == Follower {
		select {
		case <-server.shutdownCh:
			return
		case <-electionTimer: // 心跳超时
			server.SetState(Candidate)
		case rpc := <-server.rpcCh: // 把AppendEntriesRequest和VoteRequest放到一个等待队列里，串行执行，防止中间状态错乱
			switch data := rpc.Command.(type) {
			case NoopCommand:
				slog.Warn("Follower Receive Command")
				rpc.Respond(nil, errx.ErrNotLeader)
			case *AppendEntriesRequest:
				electionTimer = randomTimeout(ElectionTimeout) // 重置计时器

				// 处理 AppendEntriesRequest
				resp := server.processAppendEntriesRequest(data)

				rpc.Respond(resp, nil)
			case *VoteRequest:
				// 处理 VoteRequest
				resp := server.processVoteRequest(data)

				// 必须在给对方投票的前提下，才能重置 ElectionTimeout 计时器
				if resp.Granted {
					electionTimer = randomTimeout(ElectionTimeout)
				}

				rpc.Respond(resp, nil)
			}
		}
	}

}

func (server *Server) CandidateLoop() {
	slog.Info("Run As Candidate", "id", server.ID)
	var leaderChangeTimer <-chan time.Time // 竞选倒计时
	doVote := true                         // 本次 for 循环是否要发起投票
	voteGranted := 0                       // 获得的票数

	for server.GetState() == Candidate {
		if doVote {
			// Term + 1
			server.term++

			// 给自己投票
			server.votedFor = server.ID
			voteGranted++

			// 让其他节点投票
			lastLogIndex, lastLogTerm := server.log.LastLogInfo()
			req := VoteRequest{
				CandidateID:  server.ID,
				Term:         server.term,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}
			for _, peer := range server.peers {
				server.routineGroup.Add(1)
				go func(peer *Peer) {
					defer server.routineGroup.Done()
					resp, err := server.transporter.RequestVote(peer, &req)
					if err == nil {
						rpc := RPC{
							Command:      resp,
							ResponseChan: nil,
						}
						server.rpcCh <- rpc
					}
				}(peer)
			}

			doVote = false
			leaderChangeTimer = randomTimeout(LeaderChangeTimeout) // 竞选倒计时开始
		}

		// 选举成功
		if voteGranted >= server.QuorumSize() {
			server.SetState(Leader)
			return
		}

		select {
		case <-server.shutdownCh:
			return
		case <-leaderChangeTimer: // 选举超时
			doVote = true
		case rpc := <-server.rpcCh:
			switch data := rpc.Command.(type) {
			case NoopCommand:
				rpc.Respond(nil, errx.ErrNotLeader)
			case *AppendEntriesRequest:
				resp := server.processAppendEntriesRequest(data)
				rpc.Respond(resp, nil)
			case *VoteResponse:
				// 对于任何RPC请求或响应，只要对方发过来的Term比自己的大，就无条件地用对方的Term覆盖自己的Term，并把自己降为Follower
				if data.Term > server.term {
					// 升级term，把 votedFor 清空，把自己降为 Follower
					server.upgradeTerm(data.Term, "")
					return
				}
				if data.Granted {
					voteGranted++
				}
			case *VoteRequest: // 也可能会收到其他 Candidate 的投票请求
				resp := server.processVoteRequest(data)
				rpc.Respond(resp, nil)
			case *AppendEntriesResponse:
			}
		}
	}
}

func (server *Server) LeaderLoop() {
	slog.Info("Run As Leader", "id", server.ID)
	heartbeatTicker := time.NewTicker(HeartBeatInterval) // 定时发送心跳

	// 对于新上任的 Leader, 认为所有 Follower 的 LastLogIndex (即 Leader 的 PrevLogIndex 集合) 与自己相同
	// 通过一轮 AE 的返回结果将 PrevLogIndex 改为真正的值
	lastLogIndex := server.log.LastLogIndex()
	for _, peer := range server.peers {
		server.prevLogIndex[peer.ID] = lastLogIndex
	}

	// 成为 Leader 后, 立即发一个 HeartBeat, 让其他 Candidate 放弃
	server.doHeartBeat()

	for server.GetState() == Leader {
		select {
		case <-server.shutdownCh:
			return
		case <-heartbeatTicker.C:
			server.doHeartBeat()
		case rpc := <-server.rpcCh:
			switch data := rpc.Command.(type) {
			case NoopCommand:
				server.processCommand(data)
			case *AppendEntriesRequest:
				resp := server.processAppendEntriesRequest(data)
				rpc.Respond(resp, nil)
			case *AppendEntriesResponse:
				server.processAppendEntriesResponse(data)
			case *VoteRequest:
				resp := server.processVoteRequest(data)
				rpc.Respond(resp, nil)
			case *VoteResponse:
			default:
				slog.Warn("rpc.NoopCommand", "type", reflect.TypeOf(data).Name())
			}
		}
	}

}

// 打印集群的 LastLogIndex 和 CommitIndex
func (server *Server) print() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		<-ticker.C
		if server.GetState() == Leader {
			fmt.Println("raft cluster info", "leader", server.ID, "log index", server.log.LastLogIndex(), "commit index", server.log.CommitIndex())

			for _, peer := range server.peers {
				prevLogIndex := server.prevLogIndex[peer.ID]
				fmt.Println("raft cluster info", "follower", peer.ID, "log index", prevLogIndex)
			}
		}
	}
}

func (server *Server) processAppendEntriesRequest(req *AppendEntriesRequest) *AppendEntriesResponse {

}

func (server *Server) processVoteRequest(req *VoteRequest) *VoteResponse {

}

func (leader *Server) processAppendEntriesResponse(resp *AppendEntriesResponse) {
}

func (leader *Server) processCommand(command NoopCommand) {
}

func (leader *Server) doHeartBeat() {
}

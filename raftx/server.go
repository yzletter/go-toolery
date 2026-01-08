package raftx

import (
	"fmt"
	"log/slog"
	"reflect"
	"sort"
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
	if req.Term < server.term {
		return &AppendEntriesResponse{
			FollowerID:   server.ID,
			Term:         server.term,
			CommitIndex:  server.log.CommitIndex(),
			LastLogIndex: server.log.LastLogIndex(),
			Success:      false,
		}
	}

	// 升级term，把votedFor清空，重置leaderId，把自己降为follower
	server.upgradeTerm(req.Term, req.LeaderID)

	if len(req.LogEntries) == 0 {
		return &AppendEntriesResponse{
			FollowerID:   server.ID,
			Term:         server.term,
			CommitIndex:  server.log.CommitIndex(),
			LastLogIndex: server.log.LastLogIndex(),
			Success:      false,
		}
	}

	ok := server.log.AppendEntries(req.PrevLogIndex, req.PrevLogTerm, req.LogEntries)

	// 先更新 CommitIndex 再返回, 因为但只接收了部分
	commitIndex := min(server.log.LastLogIndex(), req.CommitIndex)
	server.log.SetCommitIndex(commitIndex)

	if !ok {
		return &AppendEntriesResponse{
			FollowerID:   server.ID,
			Term:         server.term,
			CommitIndex:  server.log.CommitIndex(),
			LastLogIndex: server.log.LastLogIndex(),
			Success:      false,
		}
	}

	return &AppendEntriesResponse{
		FollowerID:   server.ID,
		Term:         server.term,
		CommitIndex:  server.log.CommitIndex(),
		LastLogIndex: server.log.LastLogIndex(),
		Success:      true,
	}
}

func (server *Server) processVoteRequest(req *VoteRequest) *VoteResponse {
	if req.Term > server.term {
		server.upgradeTerm(req.Term, "")
	} else if req.Term < server.term { // 拒绝投票
		return &VoteResponse{Term: server.term, Granted: false}
	} else if len(server.votedFor) > 0 && server.votedFor != req.CandidateID {
		// 同 Term 内给别的 Candidate 投过票了
		return &VoteResponse{Term: server.term, Granted: false}
	}

	lastLogIndex, lastLogTerm := server.log.LastLogInfo()
	if req.LastLogIndex < lastLogIndex || req.LastLogTerm < lastLogTerm {
		return &VoteResponse{Term: server.term, Granted: false}
	}

	// 进行投票
	server.upgradeTerm(req.Term, req.CandidateID)
	server.votedFor = req.CandidateID
	return &VoteResponse{Term: server.term, Granted: true}
}

// Leader 操作

// 处理 AE 响应
func (leader *Server) processAppendEntriesResponse(resp *AppendEntriesResponse) {
	if resp.Term > leader.term {
		leader.upgradeTerm(resp.Term, "") // 升级term，把votedFor清空，重置leaderId，把自己降为follower。此时还不能确定谁是Leader，先把LeaderId置空，收到心跳后就知道谁是Leader了。目前代码里也没有使用leaderId
		return
	}

	// 更新 PrevLogIndex
	leader.prevLogIndex[resp.FollowerID] = resp.LastLogIndex

	// Follower 没有接收AE请求
	if !resp.Success {
		return
	}

	// Follower 接收了 AE 请求, 更新 Leader 的 CommitIndex
	var indices []int64

	// 将 Leader 自己的 LastLogIndex 加进去
	indices = append(indices, leader.log.LastLogIndex())

	// 将所有 Follower 的 LastLogIndex 加进去
	for _, follower := range leader.peers {
		indices = append(indices, leader.prevLogIndex[follower.ID])
	}

	// 将所有节点的 LastLogIndex 进行排序
	sort.Slice(indices, func(i, j int) bool {
		return indices[i] < indices[j]
	})

	// 计算新的 CommitIndex
	commitIndex := indices[leader.QuorumSize()-1]

	// 判断是否需要更新 CommitIndex
	if commitIndex > leader.log.commitIndex {
		leader.log.SetCommitIndex(commitIndex)
	}
}

// 处理 Client Command
func (leader *Server) processCommand(command NoopCommand) {
	leader.log.CreateEntry(command)
}

// 发送心跳
func (leader *Server) doHeartBeat() {
	for _, peer := range leader.peers {
		leader.routineGroup.Add(1)
		go func(peer *Peer) {
			defer leader.routineGroup.Done()

			prevLogIndex := leader.prevLogIndex[peer.ID]
			entries, prevLogTerm := leader.log.GetEntriesAfter(prevLogIndex)
			req := AppendEntriesRequest{
				LeaderID:     leader.ID,
				Term:         leader.term,
				CommitIndex:  leader.log.CommitIndex(),
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				LogEntries:   entries,
			}

			resp, err := leader.transporter.AppendEntries(peer, &req)
			if err == nil {
				rpc := RPC{
					Command: resp,
				}
				leader.rpcCh <- rpc // 响应放入 ch
			}

		}(peer)
	}
}

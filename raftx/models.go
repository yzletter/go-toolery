package raftx

// FSM Finite State Machine 有限状态机, 需要和具体业务结合
type FSM interface {
}

// NoopCommand no operation 空 Command
type NoopCommand struct {
}

func (NoopCommand) Apply(sm FSM) error {
	return nil
}

// VoteRequest 投票请求
type VoteRequest struct {
	CandidateID  string
	Term         int64
	LastLogIndex int64
	LastLogTerm  int64
}

// VoteResponse 投票响应
type VoteResponse struct {
	Term    int64
	Granted bool
}

// AppendEntriesRequest 追加日志请求
type AppendEntriesRequest struct {
	LeaderID     string
	Term         int64
	CommitIndex  int64
	PrevLogIndex int64
	PrevLogTerm  int64
	LogEntries   []*LogEntry
}

// AppendEntriesResponse 追加日志响应
type AppendEntriesResponse struct {
	FollowerID   string
	Term         int64
	CommitIndex  int64
	LastLogIndex int64
	Success      bool
}

// State 节点所处状态
type State uint32

const (
	Follower State = iota
	Candidate
	Leader
	Stopped
)

// Go 打印变量时候, 默认会调用 String 方法
func (s State) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

package raftx

import "sync"

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
	peers        []*Peer
	prevLogIndex map[string]int64
	fsm          FSM
	transporter  Transporter

	shutdownCh chan struct{}
	rpcCh      chan RPC
}

type Peer struct {
	ID               string
	ConnectionString string // IP 和端口号
}

func (server *Server) GetState() State {
	return server.state
}

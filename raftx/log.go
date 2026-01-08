package raftx

import (
	"log/slog"
	"sync"
)

type Log struct {
	sync.RWMutex
	server  *Server
	entries []*LogEntry // 只维护近期的日志
	// file os.File 	// 日志写内存的同时就要写磁盘, 宕机重启后从磁盘回复日志
	startIndex  int64 // 从 1 开始, 0 表示日志集合为空
	commitIndex int64 // 初始为 0
}

// LogEntry 一条日志
type LogEntry struct {
	Term        int64
	Index       int64
	NoopCommand NoopCommand
}

func NewLog(s *Server) *Log {
	return &Log{
		RWMutex:     sync.RWMutex{},
		server:      s,
		entries:     make([]*LogEntry, 0, 100),
		startIndex:  1, // 当你需要创建 Log 的时候, 意味着你至少有一个 Log Entry 要添加
		commitIndex: 0,
	}
}

// CreateEntry 创建一条日志
func (log *Log) CreateEntry(command NoopCommand) *LogEntry {
	// 上锁
	log.Lock()
	defer log.Unlock()

	// 获取 lastIndex
	var lastIndex int64
	if len(log.entries) > 0 {
		entry := log.entries[len(log.entries)-1] // 取当前 Entries 中最后一条 Entry
		lastIndex = entry.Index
	}

	entry := &LogEntry{
		Term:        log.server.term, // Leader 的任期
		Index:       lastIndex + 1,
		NoopCommand: command,
	}

	log.entries = append(log.entries, entry)
	return entry
}

// LastLogIndex 返回最后一条日志的 Index
func (log *Log) LastLogIndex() int64 {
	log.RLock()
	defer log.RUnlock()

	if len(log.entries) == 0 {
		return 0
	}

	return log.entries[len(log.entries)-1].Index
}

func (log *Log) LastLogInfo() (int64, int64) {
	log.RLock()
	defer log.RUnlock()

	if len(log.entries) == 0 {
		return 0, 0
	}

	entry := log.entries[len(log.entries)-1]
	return entry.Index, entry.Term
}

func (log *Log) CommitIndex() int64 {
	log.RLock()
	defer log.RUnlock()

	return log.commitIndex
}

// 查找 Index 为 target 的 Entry
func (log *Log) findByIndex(target int64) (int, *LogEntry) {
	// log index 是递增的, 但不一定连续
	l, r := 0, len(log.entries)-1
	for l < r {
		mid := (l + r) / 2
		if log.entries[mid].Index < target {
			l = mid + 1
		} else {
			r = mid
		}
	}

	if log.entries[l].Index != target {
		return -1, nil
	}
	return l, log.entries[l]
}

// SetCommitIndex 提交第 idx 条日志
func (log *Log) SetCommitIndex(idx int64) int64 {
	log.Lock()
	defer log.Unlock()

	// 查找 idx 日志
	i, entry := log.findByIndex(idx)
	if i == -1 || entry == nil {
		slog.Warn("Find Log Entry Failed", "log index", idx)
		return -1
	}

	// 只能提交本 Term 的 Index
	if log.server.state == Leader && log.server.term != entry.Term {
		return -1
	}

	// CommitIndex 只能比之前的大
	if idx <= log.commitIndex {
		return -1
	}

	// 把 [prevCommitIndex + 1, commitIndex] 之间的 entry 提交
	prevCommitIndex := log.commitIndex
	log.commitIndex = idx

	j, _ := log.findByIndex(prevCommitIndex)
	if j < 0 {
		return -1
	}

	// 下标 [j + 1, i]
	for k := j + 1; k <= i; k++ {
		entry := log.entries[k]
		entry.NoopCommand.Apply(log.server.fsm)
	}

	return idx
}

// AppendEntries 从 log 集合中找到 <prevLogIndex, prevLogTerm> 这条日志, 然后把 Entries 追加到后面
func (log *Log) AppendEntries(prevLogIndex, prevLogTerm int64, entries []*LogEntry) bool {
	// 上锁
	log.Lock()
	defer log.Unlock()

	if len(entries) == 0 {
		return false
	}

	total := len(log.entries)
	if prevLogIndex == 0 {
		// 全部覆盖
		log.entries = entries
		return true
	}

	i, entry := log.findByIndex(prevLogIndex)
	if i < 0 || entry == nil {
		return false
	}

	// CommitIndex 过的日志不能被覆盖
	if entry.Index < log.startIndex || entry.Index < log.commitIndex {
		return false
	}

	// index 相同但 term 不同
	if entry.Term != prevLogTerm {
		log.entries = log.entries[:i] // 舍弃后半部分
		return false
	}

	if i < total-1 { // 删除当前位置后面的所有日志
		log.entries = log.entries[:i+1]
	}
	log.entries = append(log.entries, entries...) // 追加

	return true
}

// GetEntriesAfter 获取 prevLogIndex 之后的日志传和最后一条日志的 Term
func (log *Log) GetEntriesAfter(prevLogIndex int64) ([]*LogEntry, int64) {
	log.Lock()
	defer log.Unlock()

	if len(log.entries) == 0 {
		return nil, 0
	}

	// 比 Leader 的最后一个 idx 还要大
	if prevLogIndex >= log.entries[len(log.entries)-1].Index {
		return nil, 0
	}

	total := len(log.entries)

	var entries []*LogEntry
	if prevLogIndex < log.startIndex {
		// 直接把所有日志返回即可
		entries = log.entries
	}

	i, _ := log.findByIndex(prevLogIndex)
	if i < total-1 {
		entries = log.entries[i+1:]
	}

	// 超过单次上限
	if len(entries) > MaxLogEntriesPerHeartBeat {
		entries = entries[:MaxLogEntriesPerHeartBeat]
	}

	if len(entries) == 0 {
		return nil, 0
	}

	return entries, entries[len(entries)-1].Term
}

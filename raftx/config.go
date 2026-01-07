package raftx

import "time"

const (
	HeatBeatInterval          = 20 * time.Millisecond  // 心跳间隔, 由于收到 Log Entry 后还需要写磁盘, 所以时间不能太短, 建议 0.5 - 20ms
	ElectionTimeout           = 250 * time.Millisecond // 心跳等待超时, Follower 多久接收不到心跳就变为 Candidate, 这个时间要明显大于 HeatBeatInterval
	LeaderChangeTimeout       = 500 * time.Millisecond // 选主等待超时, 建议 10 - 500ms, 要明显大于 HeatBeatInterval
	MaxLogEntriesPerHeartBeat = 100                    // 每次 AE 最多发送多少条日志, 太多了会超过 rpc data size 的限制, 也会让 leader 阻塞更长时间
)

// 检查全局参数合法性
func init() {
	if MaxLogEntriesPerHeartBeat < 1 {
		panic("MaxLogEntriesPerHeartBeat Must More Than One")
	}

	if LeaderChangeTimeout <= 2*HeatBeatInterval {
		panic("LeaderChangeTimeout Too Short")
	}

	if ElectionTimeout <= 2*HeatBeatInterval {
		panic("ElectionTimeout Too Short")
	}
}

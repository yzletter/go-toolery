package utilx

import (
	"math"
	"math/rand"
	"sync/atomic"

	"github.com/yzletter/go-toolery/errx"
)

// MinimumConcurrencyBalancer 负载均衡算法 —— 基于最小并发度的负载均衡算法
type MinimumConcurrencyBalancer struct {
	endPoints   []string // 服务器地址
	concurrency []int32  // 并发度情况
}

// NewMinimumConcurrencyBalancer 构造函数
func NewMinimumConcurrencyBalancer(endPoints []string, concurrency []int32) *MinimumConcurrencyBalancer {
	return &MinimumConcurrencyBalancer{
		endPoints:   endPoints,
		concurrency: concurrency,
	}
}

// Take 取出一个 EndPoint, 返回下标和地址
func (balancer *MinimumConcurrencyBalancer) Take() (int, string) {
	// 计算 EndPoints 长度
	length := len(balancer.endPoints)

	// 没有节点或只有一个节点
	if length <= 0 {
		return -1, ""
	} else if length == 1 {
		return 0, balancer.endPoints[0]
	}

	// 取 EndPoint
	index, minm := -1, int32(math.MaxInt32)

	begin := rand.Intn(length) // 从随机起始处开始
	for offset := 0; offset < length; offset++ {
		now := (begin + offset) % length // 当前下标位置
		c := atomic.LoadInt32(&balancer.concurrency[now])
		if c < minm { // 若 c == minm, 则不会选, 所以开始要随机开始
			minm = c
			index = now
		}
	}

	// 当前 EndPoint 并发度 + 1
	atomic.AddInt32(&balancer.concurrency[index], 1)
	return index, balancer.endPoints[index]
}

// Return 将 EndPoint 归还, 返回可能的错误
func (balancer *MinimumConcurrencyBalancer) Return(index int) error {
	// 参数校验
	if index < 0 || index > len(balancer.endPoints) {
		return errx.ErrMinConcurrencyInvalidParam
	}

	// 并发度 - 1
	atomic.AddInt32(&balancer.concurrency[index], -1)
	return nil
}

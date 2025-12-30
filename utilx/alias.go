package utilx

import (
	"math/rand"

	"github.com/yzletter/go-toolery/datastructurex/stackx"
	"github.com/yzletter/go-toolery/errx"
)

// AliasSampler Alias 采样
type AliasSampler struct {
	accept []float64 // 属于自己的概率是多少
	alias  []int     // 其他部分从谁那补来的
}

// NewAliasSampler 构造函数
func NewAliasSampler(weight []float64) (*AliasSampler, error) {
	length := len(weight)

	// 参数校验
	if length <= 0 {
		return nil, errx.ErrAliasInvalidParam
	}

	sum := 0.0 // 求和
	for _, prob := range weight {
		sum += prob
	}

	// 构造两个栈, 存放下标
	largeStk := stackx.NewStack[int]()
	smallStk := stackx.NewStack[int]()

	// 计算每个事件的概率再乘以事件个数
	probs := make([]float64, length)
	for i, _ := range probs {
		probs[i] = weight[i] / sum
	}

	for i, _ := range probs {
		probs[i] *= float64(length)
		if probs[i] < 1 {
			smallStk.Push(i)
		} else {
			largeStk.Push(i)
		}
	}

	accept := make([]float64, length)
	alias := make([]int, length)
	for i := 0; i < length; i++ {
		accept[i] = 1.0
		alias[i] = -1
	}

	for smallStk.Size() > 0 && largeStk.Size() > 0 {
		// 取堆顶, 弹出堆顶
		smallIndex, _ := smallStk.Top()
		_ = smallStk.Pop()
		largeIndex, _ := largeStk.Top()
		_ = largeStk.Pop()

		// 计算差额
		delta := 1 - probs[smallIndex]

		// 进行补差
		probs[largeIndex] -= delta
		if probs[largeIndex] > 1 {
			// 补给 small 还 大于 1, 继续放大栈
			largeStk.Push(largeIndex)
		} else if probs[largeIndex] < 1 {
			// 补给 small 导致自己小于 1, 放小栈
			smallStk.Push(largeIndex)
		}

		// 记录 accept 和 alias
		accept[smallIndex] = probs[smallIndex] // 自己的概率
		alias[smallIndex] = largeIndex         // 补给自己的是谁
	}

	return &AliasSampler{
		accept: accept,
		alias:  alias,
	}, nil
}

// Sample 进行采样, 支持并发
func (sampler *AliasSampler) Sample() int {
	idx := rand.Intn(len(sampler.accept)) // 生成 [0, n - 1] 的随机整数
	f := rand.Float64()                   // 生成 [0, 1) 的随机小数

	if f < sampler.accept[idx] {
		// 若 f < accept[idx] 则采样 idx
		return idx
	} else {
		// 否则采样 alias[idx]
		return sampler.alias[idx]
	}
}

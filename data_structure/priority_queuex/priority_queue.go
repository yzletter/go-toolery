package priority_queuex

import (
	"cmp"

	"github.com/yzletter/go-toolery/errx"
)

// PriorityQueue 堆
type PriorityQueue[T cmp.Ordered] struct {
	Data        []T
	compareFunc func(a, b T) bool
}

// NewPriorityQueue 构造一个空堆, 传入比较函数
func NewPriorityQueue[T cmp.Ordered](compareFunc func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		Data:        make([]T, 0),
		compareFunc: compareFunc,
	}
}

// Push 新元素入堆
func (heap *PriorityQueue[T]) Push(val T) {
	heap.Data = append(heap.Data, val) // 将新元素放到堆底
	heap.pushUp()
}

// Pop 弹出堆顶, 返回可能的错误
func (heap *PriorityQueue[T]) Pop() error {
	if heap.Size() == 0 {
		return errx.ErrPriorityQueueEmpty
	}
	// 将最后一个元素放到堆顶
	heap.Data[0] = heap.Data[heap.Size()-1]
	heap.Data = heap.Data[:heap.Size()-1]
	// 堆顶元素沉底
	heap.pushDown()
	return nil
}

// Top 返回堆顶元素和可能的错误
func (heap *PriorityQueue[T]) Top() (T, error) {
	if heap.Size() == 0 {
		var t T
		return t, errx.ErrPriorityQueueEmpty
	}
	return heap.Data[0], nil
}

// Size 返回堆的大小
func (heap *PriorityQueue[T]) Size() int {
	return len(heap.Data)
}

// 将堆底元素向上更新
func (heap *PriorityQueue[T]) pushUp() {
	child := heap.Size() - 1
	parent := (child - 1) >> 1
	// 向上更新
	for parent >= 0 && heap.compareFunc(heap.Data[child], heap.Data[parent]) {
		heap.Data[child], heap.Data[parent] = heap.Data[parent], heap.Data[child]
		child = parent
		parent = (child - 1) / 2
	}
}

// 将堆顶元素向下更新
func (heap *PriorityQueue[T]) pushDown() {
	// now 表示当前需要更新的节点, target 表示要将 now 与哪个节点换
	now, target := 0, 0
	for {

		// 比较当前节点和左右儿子, 将当前节点与 target 进行交换
		leftChild, rightChild := 2*now+1, 2*now+2 // 左右儿子
		if leftChild <= heap.Size()-1 && heap.compareFunc(heap.Data[leftChild], heap.Data[target]) {
			target = leftChild
		}
		if rightChild <= heap.Size()-1 && heap.compareFunc(heap.Data[rightChild], heap.Data[target]) {
			target = rightChild
		}
		if target == now { // 左右儿子都不满足交换条件, 无需更换
			break
		}
		// 交换
		heap.Data[now], heap.Data[target] = heap.Data[target], heap.Data[now]
		now = target
	}
}

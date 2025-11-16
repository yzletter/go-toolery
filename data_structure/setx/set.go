package setx

// EmptyStruct 空结构体
type EmptyStruct struct{}

// Set 集合
type Set[T comparable] map[T]EmptyStruct

// NewSet 构造函数
func NewSet[T comparable]() Set[T] {
	return make(map[T]EmptyStruct)
}

// Insert 插入一个新元素
func (hash Set[T]) Insert(val T) {
	hash[val] = EmptyStruct{}
}

// Delete 删除一个元素
func (hash Set[T]) Delete(val T) {
	delete(hash, val)
}

// Exist 查询一个元素
func (hash Set[T]) Exist(val T) bool {
	_, ok := hash[val]
	return ok
}

// Size 返回集合中元素个数
func (hash Set[T]) Size() int {
	return len(hash)
}

// Values 返回集合中所有元素的切片
func (hash Set[T]) Values() []T {
	ans := make([]T, hash.Size())
	idx := 0
	for v, _ := range hash {
		ans[idx] = v
		idx++
	}
	return ans
}

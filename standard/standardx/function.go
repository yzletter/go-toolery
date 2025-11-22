package standardx

// Ternary 三目运算符, 传入 bool 和可能返回的两个变量
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

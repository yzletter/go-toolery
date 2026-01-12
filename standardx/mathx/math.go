package mathx

import (
	"math"

	"github.com/yzletter/go-toolery/errx"
)

// QMI 快速幂求 a ^ k % p
func QMI(a, k, p int) int {
	res := 1
	for k > 0 {
		if k&1 == 1 {
			res = res * a % p
		}
		a = a * a % p
		k >>= 1
	}
	return res
}

// NormVector 向量归一化
func NormVector(vector []float64) ([]float64, error) {
	// 检查参数
	if vector == nil || len(vector) == 0 {
		return nil, errx.ErrMathInvalidParam
	}

	// 计算模长
	sum := 0.
	for _, degree := range vector {
		sum += degree * degree
	}
	norm := math.Sqrt(sum)

	for i := range vector {
		vector[i] /= norm
	}

	return vector, nil
}

// InnerProduct 计算向量内积
func InnerProduct(vector1, vector2 []float64) (float64, error) {
	if vector1 == nil || vector2 == nil || len(vector1) != len(vector2) || len(vector1) == 0 {
		return 0, errx.ErrMathInvalidParam
	}

	sum := 0.

	for i := range vector1 {
		sum += vector1[i] * vector2[i]
	}

	return sum, nil
}

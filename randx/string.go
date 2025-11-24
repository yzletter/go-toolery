package randx

import (
	"math/rand"
	"strings"
)

// RandString 根据传入的种子, 生成长度为 len 的随机字符串
func RandString(seed string, length int) string {
	res := strings.Builder{}
	letterCollection := []rune(seed)
	for i := 0; i < length; i++ {
		randIndex := rand.Intn(len(letterCollection))
		res.WriteRune(letterCollection[randIndex])
	}
	return res.String()
}

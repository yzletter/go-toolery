package standardx

import (
	"crypto/md5"
	"encoding/hex"
)

// Ternary 三目运算符, 传入 bool 和可能返回的两个变量
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// Hash 返回字符串 MD5 哈希后 32 位的十六进制编码结果
func Hash(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	digest := hasher.Sum(nil)
	return hex.EncodeToString(digest)
}

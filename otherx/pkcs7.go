package otherx

import (
	"bytes"

	"github.com/yzletter/go-toolery/errx"
)

// Padding 将数据填充, 填充至 blockSize 的整数倍, 返回填充好的切片
func Padding(src []byte, blockSize int) []byte {
	srcLen := len(src)                     // 源数据长度
	padLen := blockSize - srcLen%blockSize // 需要填充的长度

	suffix := bytes.Repeat([]byte{byte(padLen)}, padLen) // 将 padLen 重复 padLen 个
	res := append(src, suffix...)
	return res
}

// UnPadding 将数据还原, 返回填充好的切片和可能存在的错误
func UnPadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src) // 源数据长度
	// 判断合法
	if srcLen%blockSize != 0 || srcLen < blockSize {
		return nil, errx.ErrPKCS7InvalidParam
	}
	padLen := int(src[srcLen-1]) // 填充的长度
	return src[:srcLen-padLen], nil
}

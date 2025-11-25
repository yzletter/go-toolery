package jwtx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/yzletter/go-toolery/errx"
)

// GenJWT 根据 payload 和 secret 生成 JWT, 返回生成的 JWT token 和可能的错误
func GenJWT(payload JwtPayload, secret string) (string, error) {
	// 参数校验
	if secret == "" {
		return "", errx.ErrJwtInvalidParam
	}

	// 1. header 转成 json, 再用 base64 编码, 得到 JWT 第一部分
	header := DefaultHeader
	part1, err := marshalBase64Encode(header)
	if err != nil {
		return "", err
	}

	// 2. payload 转成 json, 再用 base64 编码, 得到 JWT 第二部分
	part2, err := marshalBase64Encode(payload)
	if err != nil {
		return "", err
	}

	// 3. 根据 msg 使用 secret 进行加密得到签名 signature
	jwtMsg := part1 + "." + part2                                       // JWT 信息部分
	hash := hmac.New(sha256.New, []byte(secret))                        // 根据 secret 生成 sha256 哈希算法器
	hash.Write([]byte(jwtMsg))                                          // 写入 msg
	jwtSignature := base64.RawURLEncoding.EncodeToString(hash.Sum(nil)) // 哈希结果进行 base64 编码, 得到第三部分

	return jwtMsg + "." + jwtSignature, nil
}

// VerifyJWT 根据传入的 JWT token 和 secret 校验 JWT 的合法性
func VerifyJWT(jwtToken string, secret string) (*JwtPayload, error) {
	// 参数校验
	if jwtToken == "" || secret == "" {
		return nil, errx.ErrJwtInvalidParam
	}
	parts := strings.SplitN(jwtToken, ".", 3)
	if len(parts) != 3 {
		// 传入的 JWT 格式有误
		return nil, errx.ErrJwtInvalidParam
	}

	// 获得 msg 和 signature 部分
	jwtMsg := parts[0] + "." + parts[1]
	jwtSignature := parts[2]

	// 1. 签名校验
	// 对 jwtMsg 加密得到 thisSignature 判断与 jwtSignature 是否相同
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(jwtMsg))
	thisSignature := base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
	if thisSignature != jwtSignature {
		// 签名校验失败
		return nil, errx.ErrJwtInvalidParam
	}

	// 2. 反解出 header 和 payload
	var (
		header  JwtHeader
		payload JwtPayload
	)
	err := base64DecodeUnmarshal(parts[0], &header)
	if err != nil {
		return nil, err
	}
	err = base64DecodeUnmarshal(parts[1], &payload)
	if err != nil {
		return nil, err
	}

	// 3. 时间校验
	now := time.Now()
	if payload.IssueAt > 0 && now.Add(defaultLeeway).Unix() < payload.IssueAt {
		// 当前时间(加上漂移量) < 签名时间, 签在未来
		return nil, errx.ErrJwtInvalidTime
	}
	if payload.NotBefore > 0 && now.Add(defaultLeeway).Unix() < payload.NotBefore {
		// 当前时间(加上漂移量) > 生效时间, 还未生效
		return nil, errx.ErrJwtInvalidTime
	}
	if payload.Expiration > 0 && now.Add(-defaultLeeway).Unix() > payload.Expiration {
		// 当前时间(减去漂移量) > 过期时间，已经过期
		return nil, errx.ErrJwtInvalidTime
	}

	return &payload, nil
}

// 对结构体依次进行 json 序列化和 base64 编码
func marshalBase64Encode(v any) (string, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return "", errx.ErrJwtMarshalFailed
	} else {
		return base64.RawURLEncoding.EncodeToString(bs), nil
	}
}

// 对字符串依次进行 base64 解码和 json 反序列化
func base64DecodeUnmarshal(s string, v any) error {
	bs, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return errx.ErrJwtBase64DecodeFailed
	}
	// 将 bs 反序列化到 v 中
	err = json.Unmarshal(bs, v)
	if err != nil {
		return errx.ErrJwtUnMarshalFailed
	}
	return nil
}

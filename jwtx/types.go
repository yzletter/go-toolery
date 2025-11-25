package jwtx

import "time"

// 可容忍的时间漂移
const defaultLeeway = 5 * time.Second

var (
	// DefaultHeader 默认的 JWT Header
	DefaultHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

type JwtHeader struct {
	Algo string `json:"alg"` // 哈希算法, HS256
	Type string `json:"typ"` // JWT
}

type JwtPayload struct {
	ID          string         `json:"jti"` // JWT ID
	Issue       string         `json:"iss"` // 签发者
	Audience    string         `json:"aud"` // 受众
	Subject     string         `json:"sub"` // 主题
	IssueAt     int64          `json:"iat"` // 签发时间（秒）
	NotBefore   int64          `json:"nbf"` // 生效时间（秒）
	Expiration  int64          `json:"exp"` // 过期时间（秒），0=永不过期
	UserDefined map[string]any `json:"ud"`  // 自定义字段
}

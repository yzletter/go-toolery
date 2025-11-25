package jwtx_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/yzletter/go-toolery/jwtx"
)

func TestJWT(t *testing.T) {
	secret := "abcdefg"
	payload := jwtx.JwtPayload{
		ID:         "adsidajiddasd",
		Issue:      "go_postery",
		Audience:   "anyone",
		Subject:    "shopping",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(2 * time.Hour).Unix(), // 两小时过期
		UserDefined: map[string]any{
			"name": strings.Repeat("yzletter ", 100), // 信息量很大时，jwt长度可能会超过4K
		},
	}

	if token, err := jwtx.GenJWT(payload, secret); err != nil {
		fmt.Printf("生成 JWT 失败: %v\n", err)
	} else {
		fmt.Printf("生成 JWT 成功: \n%s\n", token)
		if p, err := jwtx.VerifyJWT(token, secret); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("JWT 验证通过，欢迎 %s !\n", p.UserDefined["name"])
		}
	}
}

// go test -v ./jwtx -run=^TestJWT$ -count=1

package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

const JwtSecret = "ioclubhjr" // jwt 密匙，不能外泄

var (
	DefautHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

type JwtHeader struct {
	Algo string `json:"alg"` // 哈希算法
	Type string `json:"typ"` // 令牌类型
}

type JwtPayload struct {
	ID          string         `json:"jti"` // JWT ID 用于标识该 JWT
	Issue       string         `json:"iss"` // 发行人。比如微信
	Audience    string         `json:"aud"` // 受众人。比如王者荣耀
	Subject     string         `json:"sub"` // 主题
	IssueAt     int64          `json:"iat"` // 发布时间，精确到秒
	NotBefore   int64          `json:"nbf"` // 在此之前不可用，精确到秒
	Expiration  int64          `json:"exp"` // 到期时间，精确到秒
	UserDefined map[string]any `json:"ud"`  // 自定义的其他字段
}

func GenJwt(header JwtHeader, payLoad JwtPayload) (string, error) {
	var head, load, signature string
	// Generate the header (jwt identifier)
	if part, err := json.Marshal(header); err != nil {
		return "", err
	} else {
		head = base64.RawURLEncoding.EncodeToString(part) // Use RawURLEncoding to remove special characters like =+/ in the URL
	}
	// Generate the payload
	if part, err := json.Marshal(payLoad); err != nil {
		return "", err
	} else {
		load = base64.RawURLEncoding.EncodeToString(part)
	}
	// Generate the signature using the header and payload
	h := hmac.New(sha256.New, []byte(JwtSecret))
	h.Write([]byte(head + "." + load))
	signature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return head + "." + load + "." + signature, nil
}

func VerifyJwt(token string) (*JwtHeader, *JwtPayload, error) {
	part := strings.Split(token, ".")
	if len(part) != 3 {
		return nil, nil, errors.New("fake token")
	}
	h := hmac.New(sha256.New, []byte(JwtSecret))
	h.Write([]byte(part[0] + "." + part[1]))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if signature != part[2] {
		return nil, nil, errors.New("token verification failed")
	}
	// Try to parse the header
	var header JwtHeader
	part1, err := base64.RawURLEncoding.DecodeString(part[0])
	if err != nil {
		return nil, nil, err
	}
	if err = json.Unmarshal(part1, &header); err != nil {
		return nil, nil, err
	}
	// Try to parse the payload
	var payload JwtPayload
	part2, err := base64.RawURLEncoding.DecodeString(part[1])
	if err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal(part2, &payload); err != nil {
		return nil, nil, err
	}
	return &header, &payload, nil
}

package jwtutil

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"lostfound/pkg/redisdb"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secrectKey string

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

func Init(secret string) error {
	if secret == "" {
		return errors.New("密钥不能为空")
	}
	secrectKey = secret
	return nil
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

// randomJTI 生成随机 jti（JWT ID），保证同用户在同一秒内多次签发 token 也互不相同，
// 避免 refresh 轮换时并发插入相同的 token_hash 触发唯一索引冲突
func randomJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format("150405.000000000")
	}
	return hex.EncodeToString(b)
}

func GenerateAccessToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        randomJTI(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 设置过期时间为2小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    // 设置签发时间为当前时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secrectKey))
}

func GenerateRefreshToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        randomJTI(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 设置过期时间为7天
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 设置签发时间为当前时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secrectKey))
}

func ParseToken(tokenString string, tokenType TokenType) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secrectKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token无效")
	}
	if claims.Type != string(tokenType) {
		return nil, errors.New("token类型不匹配")
	}

	// 黑名单检查:Redis 挂时 fail-open,refresh token 由上层 service 走 DB 兜底
	revoked, _ := redisdb.IsRevokedToken(tokenString)
	if revoked {
		return nil, errors.New("token无效(黑名单)")
	}

	return claims, nil
}

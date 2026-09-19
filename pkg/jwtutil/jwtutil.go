package jwtutil

import (
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

func GenerateAccessToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
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

	n, err := redisdb.ChackToken(tokenString)
	if err != nil {
		return nil, err
	} else if n > 0 {
		return nil, errors.New("token无效(黑名单)")
	}

	return claims, nil
}

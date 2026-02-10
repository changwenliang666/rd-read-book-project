package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	jwt.RegisteredClaims
}

var secretKey = []byte("go-demo")

// 生成token
func GenerateToken(userId int, username string) (string, error) {
	claims := CustomClaims{
		Id:       userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// 解析token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

package common

import (
	"gin_demo/src/entity"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var sign_key = []byte("maplesukid")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// ReleaseToken 生成Token
func ReleaseToken(user entity.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "gin learning",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sign_key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return sign_key, nil
	})
	return token, claims, err
}

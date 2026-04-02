package utils

import (
	"bbsgo/config"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string) (string, error) {
	secret := config.GetConfig("jwt_secret")
	if secret == "" {
		secret = "change-this-secret-in-production"
	}
	
	expireDays := config.GetConfigInt("jwt_expire_days", 7)
	
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, expireDays)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	secret := config.GetConfig("jwt_secret")
	if secret == "" {
		secret = "change-this-secret-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("解析Token失败, token: %s, error: %v", tokenString, err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	log.Printf("Token无效, token: %s", tokenString)
	return nil, errors.New("invalid token")
}

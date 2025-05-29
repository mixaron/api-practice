package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type TokenService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type TokenServiceImpl struct {
	secret string
}

func NewTokenService(secret string) TokenService {
	return &TokenServiceImpl{secret}
}

func (t TokenServiceImpl) GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	return tokenString, err
}

func (t TokenServiceImpl) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT secret is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

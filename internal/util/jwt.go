package util

import (
	"auth/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJwt(key string, user *repository.User, expired time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.Username,
		"exp":       expired.Unix(),
		"role":      user.Role,
		"createdAt": user.CreatedAt.Unix(),
	})

	return token.SignedString(key)
}

func ParseJwt(key string, tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return claims, err
	} else if !tkn.Valid {
		return claims, errors.New("invalid token")
	} else {
		return claims, nil
	}
}

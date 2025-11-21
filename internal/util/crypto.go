package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//claims hold the signed users data
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

//func to generate JWT token
func GenerateJWT(userID string, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	
	return token.SignedString([]byte(secret))
}

//validate func
func ValidateJWT(tokenStr string, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected string method")
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return "", err
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}
	
	return claims.UserID, nil
}
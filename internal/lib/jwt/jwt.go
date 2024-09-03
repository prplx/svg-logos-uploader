package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string, JWTSecret string) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	})
	return token.SignedString([]byte(JWTSecret))
}

func ParseToken(tokenString string, JWTSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")

	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

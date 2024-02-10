package util

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func CreateToken(id int) (string, error) {
	secretKey := os.Getenv("secret_key")

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

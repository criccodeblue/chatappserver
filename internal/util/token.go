package util

import (
	"chatappserver/internal/model"
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func CreateToken(id int) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	jwtClaims := model.JWTClaims{
		UserId: id,
		Expiry: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (claims *model.JWTClaims, err error) {
	secretKey := os.Getenv("SECRET_KEY")

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		err = errors.New("invalid token")
		return nil, err
	}

	claims = token.Claims.(*model.JWTClaims)

	if claims.Expiry < time.Now().Unix() {
		err = errors.New("invalid token")
		return nil, err
	}

	return claims, nil
}

package util

import (
	"chatappserver/internal/model"
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

func CreateToken(id int, privateKey *rsa.PrivateKey) (string, error) {
	jwtClaims := model.JWTClaims{
		UserId: id,
		Expiry: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string, privateKey *rsa.PrivateKey) (claims *model.JWTClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey.Public(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		err = errors.New("invalid token")
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		err = errors.New("failed to parse claims")
		return nil, err
	}

	if claims.Expiry < time.Now().Unix() {
		err = errors.New("invalid token")
		return nil, err
	}

	return claims, nil
}

func GetBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	// Check if the Authorization header is in the format "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	return parts[1], nil
}

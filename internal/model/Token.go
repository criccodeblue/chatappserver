package model

import "github.com/golang-jwt/jwt"

type JWTClaims struct {
	jwt.StandardClaims
	UserId int
	Expiry int64
}

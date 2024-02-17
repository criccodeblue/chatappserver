package auth

import (
	"chatappserver/internal/util"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
)

type Auth struct {
	PrivateKey *rsa.PrivateKey
}

func NewAuth() (*Auth, error) {
	auth := &Auth{}

	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	auth.PrivateKey = privateKey
	return auth, nil
}

func (auth *Auth) VerifyAuth(endpointFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := util.GetBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		_, err = util.VerifyToken(tokenString, auth.PrivateKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		endpointFunc(w, r)
	}
}

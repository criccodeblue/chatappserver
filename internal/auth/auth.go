package auth

import (
	"chatappserver/internal/util"
	"net/http"
)

func VerifyAuth(endpointFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			tokenString := r.Header["Token"][0]

			_, err := util.VerifyToken(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			endpointFunc(w, r)
		} else {
			http.Error(w, "Token not found", http.StatusForbidden)
		}
	}
}

package util

import (
	. "chatappserver/internal/model"
	"encoding/json"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, errMessage string) {
	if jsonErr := json.NewEncoder(w).Encode(ErrorMessage{ErrorMessage: errMessage}); jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
}

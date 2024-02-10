package model

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

package model

type ErrorResponse struct {
	Message string `json:"message" example:"ERROR_MSG"`
	Code    string `json:"code" example:"BAD_REQUEST"`
}

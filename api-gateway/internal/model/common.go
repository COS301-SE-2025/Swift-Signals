package model

type ErrorResponse struct {
	Message string `json:"message" example:"resource not found"`
	Code    int    `json:"code"    example:"404"`
}

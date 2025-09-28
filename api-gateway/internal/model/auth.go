package model

type LoginRequest struct {
	Email    string `json:"email"    example:"testuser@example.com" binding:"required" validate:"required,email"`
	Password string `json:"password" example:"testpass1234"         binding:"required" validate:"required"`
}

type LoginResponse struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token"   example:"header.payload.signature"`
}

type RegisterRequest struct {
	Username string `json:"username" example:"tester"               binding:"required" validate:"required,min=3,max=32"`
	Email    string `json:"email"    example:"testuser@example.com" binding:"required" validate:"required,email"`
	Password string `json:"password" example:"testpass1234"         binding:"required" validate:"required,min=8,max=64"`
}

type RegisterResponse struct {
	UserID string `json:"user_id" example:"1"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Logout successful"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" example:"user@example.com" binding:"required"`
}

type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset instructions sent to your email."`
}

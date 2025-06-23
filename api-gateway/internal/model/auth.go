package model

type LoginRequest struct {
	Email    string `json:"email"    example:"user@example.com"  binding:"required"`
	Password string `json:"password" example:"StrongPassword123" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"               binding:"required"`
	Email    string `json:"email"    example:"newuser@example.com"   binding:"required"`
	Password string `json:"password" example:"VeryStrongPassword456" binding:"required"`
}

type AuthResponse struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token"   example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
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

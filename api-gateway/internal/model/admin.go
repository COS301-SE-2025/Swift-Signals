package model

type GetAllUsersRequest struct {
	Page     int `json:"page"      example:"1"  validate:"required,min=1"`
	PageSize int `json:"page_size" example:"10" validate:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	Username string `json:"username" example:"newusername"      validate:"omitempty,min=3,max=32"`
	Email    string `json:"email"    example:"user@example.com" validate:"omitempty,email"`
}

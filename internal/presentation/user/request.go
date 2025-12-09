package user

// CreateUserRequest はユーザー作成のリクエスト
type CreateUserRequest struct {
	Name      string  `json:"name" validate:"required"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
}

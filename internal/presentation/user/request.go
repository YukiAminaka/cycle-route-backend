package user

// CreateUserRequest はユーザー作成のリクエスト
type CreateUserRequest struct {
	KratosID  string  `json:"kratos_id" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
}

// UpdateUserProfileRequest はユーザープロフィール更新のリクエスト
type UpdateUserProfileRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
}

// UpdateUserLocationRequest はユーザー位置情報更新のリクエスト
type UpdateUserLocationRequest struct {
	Locality           string `json:"locality" validate:"required"`
	AdministrativeArea string `json:"administrative_area" validate:"required"`
	CountryCode        string `json:"country_code" validate:"required"`
	PostalCode         string `json:"postal_code" validate:"required"`
	Geom               string `json:"geom" validate:"required"`
}

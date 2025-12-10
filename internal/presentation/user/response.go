package user

// userResponse はユーザー取得のレスポンス
type userResponse struct {
	User userResponseModel `json:"user"`
}



// userResponseModel はユーザー情報のレスポンスモデル
type userResponseModel struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	HighlightedPhotoID *int64     `json:"highlighted_photo_id,omitempty"`
	Locale             *string    `json:"locale,omitempty"`
	Description        *string    `json:"description,omitempty"`
	Locality           *string    `json:"locality,omitempty"`
	AdministrativeArea *string    `json:"administrative_area,omitempty"`
	CountryCode        *string    `json:"country_code,omitempty"`
	PostalCode         *string    `json:"postal_code,omitempty"`
	Geom               *string 	  `json:"geom,omitempty"`
	FirstName          *string    `json:"first_name,omitempty"`
	LastName           *string    `json:"last_name,omitempty"`
	Email              *string    `json:"email,omitempty"`
	HasSetLocation     bool       `json:"has_set_location"`
}

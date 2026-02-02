package user

import (
	"github.com/YukiAminaka/cycle-route-backend/internal/pkg/geometry"
	"github.com/YukiAminaka/cycle-route-backend/internal/presentation/response"
	"github.com/YukiAminaka/cycle-route-backend/internal/presentation/validator"
	userUsecase "github.com/YukiAminaka/cycle-route-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb/geojson"
)

// Handler はユーザー関連のHTTPハンドラー
type Handler struct {
	createUserUsecase userUsecase.ICreateUserUsecase
	getUserUsecase    userUsecase.IGetUserByIDUsecase
}

// NewHandler はHandlerを作成する
func NewHandler(
	createUserUsecase userUsecase.ICreateUserUsecase,
	getUserUsecase userUsecase.IGetUserByIDUsecase,
) *Handler {
	return &Handler{
		createUserUsecase: createUserUsecase,
		getUserUsecase:    getUserUsecase,
	}
}

// GetUserByID godoc
//	@Summary	ユーザーを取得する
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	UserResponse
//	@Failure	400	{object}	response.ErrorResponse
//	@Failure	404	{object}	response.ErrorResponse
//	@Failure	500	{object}	response.ErrorResponse
//	@Router		/users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	dto, err := h.getUserUsecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.ReturnStatusInternalServerError(c, err)
		return
	}

	// GEOJSON 形式に変換
	var rawJSON string
	if dto.Geom != nil {
		fc := geojson.NewFeatureCollection()
		fc.Append(geojson.NewFeature(dto.Geom))
		b, _ := fc.MarshalJSON()
		rawJSON = string(b)
	}

	res := UserResponse{
		User: UserResponseModel{
			ID:                 dto.ID,
			Name:               dto.Name,
			HighlightedPhotoID: dto.HighlightedPhotoID,
			Locale:             dto.Locale,
			Description:        dto.Description,
			Locality:           dto.Locality,
			AdministrativeArea: dto.AdministrativeArea,
			CountryCode:        dto.CountryCode,
			PostalCode:         dto.PostalCode,
			Geom:               &rawJSON,
			FirstName:          dto.FirstName,
			LastName:           dto.LastName,
			Email:              dto.Email,
			HasSetLocation:     dto.HasSetLocation,
		},
	}

	response.ReturnStatusOK(c, res)
}

// CreateUser godoc
//	@Summary	ユーザーを作成する
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		request	body		CreateUserRequest	true	"Create User Request"
//	@Success	201		{object}	UserResponse
//	@Failure	400		{object}	response.ErrorResponse
//	@Failure	500		{object}	response.ErrorResponse
//	@Router		/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnBadRequest(c, err)
		return
	}

	// バリデーション
	validate := validator.GetValidator()
	if err := validate.Struct(req); err != nil {
		response.ReturnStatusBadRequest(c, err)
		return
	}

	input := userUsecase.CreateUserUseCaseInputDto{
		KratosID:  req.KratosID,
		Name:      req.Name,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	dto, err := h.createUserUsecase.CreateUser(c, input)
	if err != nil {
		response.ReturnStatusInternalServerError(c, err)
		return
	}


	res := UserResponse{
		User: UserResponseModel{
			ID:                 dto.ID,
			Name:               dto.Name,
			HighlightedPhotoID: dto.HighlightedPhotoID,
			Locale:             dto.Locale,
			Description:        dto.Description,
			Locality:           dto.Locality,
			AdministrativeArea: dto.AdministrativeArea,
			CountryCode:        dto.CountryCode,
			PostalCode:         dto.PostalCode,
			Geom:               geometry.GeometryToGeoJSON(dto.Geom),
			FirstName:          dto.FirstName,
			LastName:           dto.LastName,
			Email:              dto.Email,
			HasSetLocation:     dto.HasSetLocation,
		},
	}
	response.ReturnStatusCreated(c, res)
}

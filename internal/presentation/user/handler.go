package user

import (
	"github.com/YukiAminaka/cycle-route-backend/internal/presentation/settings"
	"github.com/YukiAminaka/cycle-route-backend/internal/presentation/validator"
	userUsecase "github.com/YukiAminaka/cycle-route-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
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
// @Summary ユーザーを取得する
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} GetUserResponse
// @Failure 400 {object} presenter.Response
// @Failure 404 {object} presenter.Response
// @Failure 500 {object} presenter.Response
// @Router /v1/users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	dto, err := h.getUserUsecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		settings.ReturnError(c, err)
		return
	}

	res := userResponse{
		User: userResponseModel{
			ID:                 dto.ID,
			Name:               dto.Name,
			HighlightedPhotoID: dto.HighlightedPhotoID,
			Locale:             dto.Locale,
			Description:        dto.Description,
			Locality:           dto.Locality,
			AdministrativeArea: dto.AdministrativeArea,
			CountryCode:        dto.CountryCode,
			PostalCode:         dto.PostalCode,
			Geom:               dto.Geom,
			FirstName:          dto.FirstName,
			LastName:           dto.LastName,
			Email:              dto.Email,
			HasSetLocation:     dto.HasSetLocation,
		},
	}

	settings.ReturnStatusOK(c, res)
}

// CreateUser godoc
// @Summary ユーザーを作成する
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create User Request"
// @Success 201 {object} presenter.Response
// @Failure 400 {object} presenter.Response
// @Failure 500 {object} presenter.Response
// @Router /v1/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		settings.ReturnBadRequest(c, err)
		return
	}

	// バリデーション
	validate := validator.GetValidator()
	if err := validate.Struct(req); err != nil {
		settings.ReturnStatusBadRequest(c, err)
		return
	}

	input := userUsecase.CreateUserUseCaseInputDto{
		Name:      req.Name,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	dto, err := h.createUserUsecase.CreateUser(c, input)
	if err != nil {
		settings.ReturnError(c, err)
		return
	}
	response := userResponse{
		userResponseModel{
			ID:                 dto.ID,
			Name:               dto.Name,
			HighlightedPhotoID: dto.HighlightedPhotoID,
			Locale:             dto.Locale,
			Description:        dto.Description,
			Locality:           dto.Locality,
			AdministrativeArea: dto.AdministrativeArea,
			CountryCode:        dto.CountryCode,
			PostalCode:         dto.PostalCode,
			Geom:               dto.Geom,
			FirstName:          dto.FirstName,
			LastName:           dto.LastName,
			Email:              dto.Email,
			HasSetLocation:     dto.HasSetLocation,
		},
	}
	settings.ReturnStatusCreated(c, response)
}

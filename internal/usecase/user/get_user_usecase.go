package usecase

import (
	"context"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/paulmach/orb"
)

// IUserUsecase はユーザーユースケースのインターフェース
type IGetUserByIDUsecase interface {
	GetUserByID(ctx context.Context, id string) (*GetUserByIDUseCaseDto, error)
}

type getUserByIDUsecase struct {
	userRepo userDomain.IUserRepository
}

// ユーザーユースケースを作成する
func NewGetUserByIDUsecase(userRepo userDomain.IUserRepository) IGetUserByIDUsecase {
	return &getUserByIDUsecase{
		userRepo: userRepo,
	}
}

// GetUserByIDUseCaseDto は GetUserByID ユースケースの出力DTO
type GetUserByIDUseCaseDto struct {
	ID                 string
	Name               string
	HighlightedPhotoID *int64
	Locale             *string
	Description        *string
	Locality           *string
	AdministrativeArea *string
	CountryCode        *string
	PostalCode         *string
	Geom               *orb.Point
	FirstName          *string
	LastName           *string
	Email              *string
	HasSetLocation     bool
}

func (u *getUserByIDUsecase) GetUserByID(ctx context.Context, id string) (*GetUserByIDUseCaseDto, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetUserByIDUseCaseDto{
		ID:                 user.ID().String(),
		Name:               user.Name(),
		HighlightedPhotoID: user.HighlightedPhotoID(),
		Locale:             user.Locale(),
		Description:        user.Description(),
		Locality:           user.Locality(),
		AdministrativeArea: user.AdministrativeArea(),
		CountryCode:        user.CountryCode(),
		PostalCode:         user.PostalCode(),
		Geom: func() *orb.Point {
			if user.Geom() != nil {
				if point, ok := user.Geom().Geometry.(orb.Point); ok {
					return &point
				}
			}
			return nil
		}(),
		FirstName:      user.FirstName(),
		LastName:       user.LastName(),
		Email:          user.Email(),
		HasSetLocation: user.HasSetLocation(),
	}, nil
}


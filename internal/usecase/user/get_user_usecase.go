package usecase

import (
	"context"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/paulmach/orb"
)

// IUserUsecase はユーザーユースケースのインターフェース
type IGetUserByIDUsecase interface {
	GetUserByID(ctx context.Context, id int64) (*GetUserByIDUseCaseDto, error)
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
	id                 string
	name               string
	highlightedPhotoID *int64
	locale             *string
	description        *string
	locality           *string
	administrativeArea *string
	countryCode        *string
	postalCode         *string
	geom               *orb.Point
	firstName          *string
	lastName           *string
	email              *string
	hasSetLocation     bool
}

func (u *getUserByIDUsecase) GetUserByID(ctx context.Context, id int64) (*GetUserByIDUseCaseDto, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &GetUserByIDUseCaseDto{
		id:                 user.ID().String(),
		name:               user.Name(),
		highlightedPhotoID: user.HighlightedPhotoID(),
		locale:             user.Locale(),
		description:        user.Description(),
		locality:           user.Locality(),
		administrativeArea: user.AdministrativeArea(),
		countryCode:        user.CountryCode(),
		postalCode:         user.PostalCode(),
		geom: func() *orb.Point {
			if user.Geom() != nil {
				if point, ok := user.Geom().Geometry.(orb.Point); ok {
					return &point
				}
			}
			return nil
		}(),
		firstName:      user.FirstName(),
		lastName:       user.LastName(),
		email:          user.Email(),
		hasSetLocation: user.HasSetLocation(),
	}, nil
}


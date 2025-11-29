package usecase

import (
	"context"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/paulmach/orb"
)

// IUserUsecase はユーザーユースケースのインターフェース
type ICreateUserUsecase interface {
	CreateUser(ctx context.Context, user *userDomain.User) (*CreateUserUseCaseDto, error)
}

type createUserUsecase struct {
	userRepo userDomain.IUserRepository
}

// ユーザーユースケースを作成する
func NewUserUsecase(userRepo userDomain.IUserRepository) ICreateUserUsecase {
	return &createUserUsecase{
		userRepo: userRepo,
	}
}

type CreateUserUseCaseDto struct {
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

func (u *createUserUsecase) CreateUser(ctx context.Context, user *userDomain.User) (*CreateUserUseCaseDto, error) {
	
	user, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &CreateUserUseCaseDto{
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

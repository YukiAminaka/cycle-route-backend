package usecase

import (
	"context"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/paulmach/orb"
)

// IUserUsecase はユーザーユースケースのインターフェース
type ICreateUserUsecase interface {
	CreateUser(ctx context.Context, user CreateUserUseCaseInputDto) (*CreateUserUseCaseOutputDto, error)
}

type createUserUsecase struct {
	userRepo userDomain.IUserRepository
}

// ユーザーユースケースを作成する
func NewCreateUserUsecase(userRepo userDomain.IUserRepository) ICreateUserUsecase {
	return &createUserUsecase{
		userRepo: userRepo,
	}
}

type CreateUserUseCaseInputDto struct {
	KratosID  string
	Name      string
	Email     *string
	FirstName *string
	LastName  *string
}

type CreateUserUseCaseOutputDto struct {
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

func (u *createUserUsecase) CreateUser(ctx context.Context, dto CreateUserUseCaseInputDto) (*CreateUserUseCaseOutputDto, error) {
	newUser, err := userDomain.NewUser(
		dto.KratosID,
		dto.Name,
		dto.Email,
		dto.FirstName,
		dto.LastName,
	)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &CreateUserUseCaseOutputDto{
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

package usecase

import (
	"context"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/paulmach/orb"
)

// IUserUsecase はユーザーユースケースのインターフェース
type IUpdateUserUsecase interface {
	UpdateUser(ctx context.Context, id string, dto UpdateUserUseCaseInputDto) (*UpdateUserUseCaseOutputDto, error)
}

type updateUserUsecase struct {
	userRepo userDomain.IUserRepository
}

// ユーザーユースケースを作成する
func NewUpdateUserUsecase(userRepo userDomain.IUserRepository) IUpdateUserUsecase {
	return &updateUserUsecase{
		userRepo: userRepo,
	}
}

type UpdateUserUseCaseInputDto struct {
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

type UpdateUserUseCaseOutputDto struct {
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

func (u *updateUserUsecase) UpdateUser(ctx context.Context, KratosID string, dto UpdateUserUseCaseInputDto) (*UpdateUserUseCaseOutputDto, error) {
	userEntity, err := u.userRepo.GetUserByKratosID(ctx, KratosID)
	if err != nil {
		return nil, err
	}

	updatedUser, err := userDomain.ReconstructUser(
		userEntity.ID(),
		userEntity.KratosID(),
		dto.Name,
		dto.HighlightedPhotoID,
		dto.Locale,
		dto.Description,
		dto.Locality,
		dto.AdministrativeArea,
		dto.CountryCode,
		dto.PostalCode,
		func() *userDomain.Geometry {
			if dto.Geom == nil {
				return nil
			}
			return &userDomain.Geometry{Geometry: dto.Geom}
		}(),
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.HasSetLocation,
	)
	if err != nil {
		return nil, err
	}

	user,err := u.userRepo.UpdateUser(ctx, updatedUser)
	if err != nil {
		return nil, err
	}

	return &UpdateUserUseCaseOutputDto{
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

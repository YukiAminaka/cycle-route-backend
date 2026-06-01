package user

import (
	"context"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
)

type IUpdateUserUsecase interface {
	UpdateUserProfile(ctx context.Context, kratosID string, dto UpdateUserUseCaseInputDto) error
	UpdateUserLocation(ctx context.Context, kratosID string, dto UpdateUserLocationUseCaseInputDto) error
}

type updateUserUsecase struct {
	userRepo userDomain.IUserRepository
}

func NewUpdateUserUsecase(userRepo userDomain.IUserRepository) IUpdateUserUsecase {
	return &updateUserUsecase{
		userRepo: userRepo,
	}
}

type UpdateUserUseCaseInputDto struct {
	Name        *string
	Description *string
	FirstName   *string
	LastName    *string
}

type UpdateUserLocationUseCaseInputDto struct {
	Locality           string
	AdministrativeArea string
	CountryCode        string
	PostalCode         string
	Geom               userDomain.Geometry
}

func (u *updateUserUsecase) UpdateUserProfile(ctx context.Context, kratosID string, dto UpdateUserUseCaseInputDto) error {
	userEntity, err := u.userRepo.GetUserByKratosID(ctx, kratosID)
	if err != nil {
		return err
	}

	if err := userEntity.UpdateProfile(dto.Name, dto.Description, dto.FirstName, dto.LastName); err != nil {
		return err
	}

	if err := u.userRepo.UpdateUserProfile(ctx, userEntity); err != nil {
		return err
	}

	return nil
}

func (u *updateUserUsecase) UpdateUserLocation(ctx context.Context, kratosID string, dto UpdateUserLocationUseCaseInputDto) error {
	userEntity, err := u.userRepo.GetUserByKratosID(ctx, kratosID)
	if err != nil {
		return err
	}

	if err := userEntity.SetLocation(dto.Locality, dto.AdministrativeArea, dto.CountryCode, dto.PostalCode, dto.Geom); err != nil {
		return err
	}


	if err := u.userRepo.UpdateUserLocation(ctx, userEntity); err != nil {
		return err
	}

	return nil
}

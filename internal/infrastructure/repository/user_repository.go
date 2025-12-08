package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/jackc/pgx/v5"
)

type userRepositoryImpl struct {
	queries *dbgen.Queries
}

// ユーザーリポジトリの実装
func NewUserRepository(queries *dbgen.Queries) user.IUserRepository {
	return &userRepositoryImpl{queries: queries}
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	u, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	// ドメインモデルのUserに変換して返す
	var geom *user.Geometry
	if u.Geom != nil {
		geom = &user.Geometry{Geometry: u.Geom.Geometry}
	}

	ud, err := user.ReconstructUser(
		user.UserID(u.Ulid),
		u.Name,
		u.HighlightedPhotoID,
		u.Locale,
		u.Description,
		u.Locality,
		u.AdministrativeArea,
		u.CountryCode,
		u.PostalCode,
		geom,
		u.FirstName,
		u.LastName,
		u.Email,
		u.HasSetLocation,
	)
	if err != nil {
		return nil, err
	}
	return ud, nil
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, userDomain *user.User) (*user.User, error) {
	u, err := r.queries.CreateUser(ctx, dbgen.CreateUserParams{
		Ulid:               string(userDomain.ID()),
		Name:               userDomain.Name(),
		HighlightedPhotoID: userDomain.HighlightedPhotoID(),
		Locale:             userDomain.Locale(),
		Description:        userDomain.Description(),
		Locality:           userDomain.Locality(),
		AdministrativeArea: userDomain.AdministrativeArea(),
		CountryCode:        userDomain.CountryCode(),
		PostalCode:         userDomain.PostalCode(),
		Geom: func() *dbgen.OrbGeometry {
			if userDomain.Geom() == nil {
				return nil
			}
			return &dbgen.OrbGeometry{Geometry: userDomain.Geom().Geometry}
		}(),
		FirstName:      userDomain.FirstName(),
		LastName:       userDomain.LastName(),
		Email:          userDomain.Email(),
		HasSetLocation: userDomain.HasSetLocation(),	
	})
	if err != nil {
		return nil, err
	}
	// ドメインモデルのUserに変換して返す
	var geom *user.Geometry
	if u.Geom != nil {
		geom = &user.Geometry{Geometry: u.Geom.Geometry}
	}

	ud, err := user.ReconstructUser(
		user.UserID(u.Ulid),
		u.Name,
		u.HighlightedPhotoID,
		u.Locale,
		u.Description,
		u.Locality,
		u.AdministrativeArea,
		u.CountryCode,
		u.PostalCode,
		geom,
		u.FirstName,
		u.LastName,
		u.Email,
		u.HasSetLocation,
	)
	if err != nil {
		return nil, err
	}

	return ud, nil
}

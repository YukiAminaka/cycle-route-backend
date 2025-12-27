package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/google/uuid"
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
	// UUIDに変換
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	u, err := r.queries.GetUserByID(ctx, uuid)
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
		user.UserID(u.ID.String()),
		u.KratosID.String(),
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
	// UserIDをuuid.UUIDに変換
	id, err := uuid.Parse(userDomain.ID().String())
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	// KratosIDをuuid.UUIDに変換
	kratosID, err := uuid.Parse(userDomain.KratosID())
	if err != nil {
		return nil, fmt.Errorf("invalid kratos id: %w", err)
	}

	u, err := r.queries.CreateUser(ctx, dbgen.CreateUserParams{
		ID:                 id,
		KratosID:           kratosID,
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
		user.UserID(u.ID.String()),
		u.KratosID.String(),
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

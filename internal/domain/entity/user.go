package entity

import (
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID                     int64              `json:"id"`
	Name                   string             `json:"name"`
	HighlightedPhotoID     pgtype.Int8        `json:"highlighted_photo_id"`
	Locale                 pgtype.Text        `json:"locale"`
	Description            pgtype.Text        `json:"description"`
	Locality               pgtype.Text        `json:"locality"`
	AdministrativeArea     pgtype.Text        `json:"administrative_area"`
	CountryCode            pgtype.Text        `json:"country_code"`
	PostalCode             pgtype.Text        `json:"postal_code"`
	TotalTripDistance      pgtype.Float8      `json:"total_trip_distance"`
	TotalTripDuration      pgtype.Float8      `json:"total_trip_duration"`
	TotalTripElevationGain pgtype.Float8      `json:"total_trip_elevation_gain"`
	Geom                   *sqlc.OrbGeometry  `json:"geom"`
	FirstName              pgtype.Text        `json:"first_name"`
	LastName               pgtype.Text        `json:"last_name"`
	Email                  pgtype.Text        `json:"email"`
	HasSetLocation         pgtype.Bool        `json:"has_set_location"`
}

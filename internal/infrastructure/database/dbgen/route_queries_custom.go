package dbgen

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	// import the dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type SearchRoutesByUserIDParams struct {
	UserID       uuid.UUID `json:"user_id"`
	NameKeywords []string  `json:"name_keywords"`
	Visibility   *int16    `json:"visibility"`
	MinDistance  *float64  `json:"min_distance"`
	MaxDistance  *float64  `json:"max_distance"`
}

type ExploreRoutesParams struct {
	RadiusM      *float64    `json:"radius_m"`
	Location     OrbGeometry `json:"location"`
	NameKeywords []string    `json:"name_keywords"`
	MinDistance  *float64    `json:"min_distance"`
	MaxDistance  *float64    `json:"max_distance"`
	OffsetCount  int32       `json:"offset_count"`
	LimitCount   int32       `json:"limit_count"`
}

type ExploreRoutesRow struct {
	ID                 uuid.UUID   `json:"id"`
	UserID             uuid.UUID   `json:"user_id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	HighlightedPhotoID *int64      `json:"highlighted_photo_id"`
	Distance           float64     `json:"distance"`
	Duration           float64     `json:"duration"`
	ElevationGain      float64     `json:"elevation_gain"`
	ElevationLoss      float64     `json:"elevation_loss"`
	PathGeom           OrbGeometry `json:"path_geom"`
	Bbox               OrbGeometry `json:"bbox"`
	FirstPoint         OrbGeometry `json:"first_point"`
	LastPoint          OrbGeometry `json:"last_point"`
	Polyline           string      `json:"polyline"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	Visibility         int16       `json:"visibility"`
	TotalCount         int64       `json:"total_count"`
	UserName           string      `json:"user_name"`
}

// SELECT * FROM routes
// WHERE user_id = $1
//   AND name ILIKE $2 OR name ILIKE $3
//   AND visibility = $3
//   AND distance >= $4
//   AND distance <= $5;

func (q *Queries) SearchRoutesByUserID(ctx context.Context, arg SearchRoutesByUserIDParams) ([]Route, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.
		Select(
			"id",
			"user_id",
			"name",
			"description",
			"highlighted_photo_id",
			"distance",
			"duration",
			"elevation_gain",
			"elevation_loss",
			"path_geom",
			"bbox",
			"first_point",
			"last_point",
			"polyline",
			"created_at",
			"updated_at",
			"visibility",
		).
		From("routes").
		Where(goqu.C("user_id").Eq(arg.UserID)).
		Prepared(true)

	if len(arg.NameKeywords) > 0 {
		orExprs := make([]goqu.Expression, len(arg.NameKeywords))
		for i, keyword := range arg.NameKeywords {
			orExprs[i] = goqu.C("name").ILike(keyword)
		}
		ds = ds.Where(goqu.Or(orExprs...))
	}

	if arg.Visibility != nil {
		ds = ds.Where(goqu.C("visibility").Eq(*arg.Visibility))
	}

	if arg.MinDistance != nil {
		ds = ds.Where(goqu.C("distance").Gte(*arg.MinDistance))
	}

	if arg.MaxDistance != nil {
		ds = ds.Where(goqu.C("distance").Lte(*arg.MaxDistance))
	}

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := q.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Route
	for rows.Next() {
		var i Route
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Description,
			&i.HighlightedPhotoID,
			&i.Distance,
			&i.Duration,
			&i.ElevationGain,
			&i.ElevationLoss,
			&i.PathGeom,
			&i.Bbox,
			&i.FirstPoint,
			&i.LastPoint,
			&i.Polyline,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Visibility,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// SELECT
//   "filtered_routes"."id",
//   "filtered_routes"."user_id",
//   "filtered_routes"."name",
//   "filtered_routes"."description",
//   "filtered_routes"."highlighted_photo_id",
//   "filtered_routes"."distance",
//   "filtered_routes"."duration",
//   "filtered_routes"."elevation_gain",
//   "filtered_routes"."elevation_loss",
//   "filtered_routes"."path_geom",
//   "filtered_routes"."bbox",
//   "filtered_routes"."first_point",
//   "filtered_routes"."last_point",
//   "filtered_routes"."polyline",
//   "filtered_routes"."created_at",
//   "filtered_routes"."updated_at",
//   "filtered_routes"."visibility",
//   "filtered_routes"."total_count",
//   "users"."name" AS "user_name"
// FROM (
//   SELECT
//     "id",
//     "user_id",
//     "name",
//     "description",
//     "highlighted_photo_id",
//     "distance",
//     "duration",
//     "elevation_gain",
//     "elevation_loss",
//     "path_geom",
//     "bbox",
//     "first_point",
//     "last_point",
//     "polyline",
//     "created_at",
//     "updated_at",
//     "visibility",
//     COUNT(*) OVER () AS "total_count"
//   FROM "routes"
//   WHERE (("visibility" = $1) AND ST_DWithin(first_point :: geography, ST_GeomFromEWKB($2) :: geography, $3) AND (("name" ILIKE $4) OR ("name" ILIKE $5)) AND ("distance" >= $6) AND ("distance" <= $7))
// ) AS "filtered_routes"
// INNER JOIN "users"
// ON ("filtered_routes"."user_id" = "users"."id")
// ORDER BY ST_Distance(filtered_routes.first_point :: geography, ST_GeomFromEWKB($8) :: geography) ASC
// LIMIT $9
// OFFSET $10

func (q *Queries) ExploreRoutes(ctx context.Context, arg ExploreRoutesParams) ([]ExploreRoutesRow, error) {
	dialect := goqu.Dialect("postgres")

	hasLocation := arg.RadiusM != nil && arg.Location.Geometry != nil

	filtered := dialect.
		Select(
			"id",
			"user_id",
			"name",
			"description",
			"highlighted_photo_id",
			"distance",
			"duration",
			"elevation_gain",
			"elevation_loss",
			"path_geom",
			"bbox",
			"first_point",
			"last_point",
			"polyline",
			"created_at",
			"updated_at",
			"visibility",
			goqu.COUNT("*").Over(goqu.W()).As("total_count"),
		).
		From("routes").
		Where(goqu.C("visibility").Eq(1))

	if hasLocation {
		filtered = filtered.Where(goqu.L(
			"ST_DWithin(first_point::geography, ST_GeomFromEWKB(?)::geography, ?)",
			arg.Location, *arg.RadiusM,
		))
	}

	if len(arg.NameKeywords) > 0 {
		orExprs := make([]goqu.Expression, len(arg.NameKeywords))
		for i, keyword := range arg.NameKeywords {
			orExprs[i] = goqu.C("name").ILike(keyword)
		}
		filtered = filtered.Where(goqu.Or(orExprs...))
	}

	if arg.MinDistance != nil {
		filtered = filtered.Where(goqu.C("distance").Gte(*arg.MinDistance))
	}

	if arg.MaxDistance != nil {
		filtered = filtered.Where(goqu.C("distance").Lte(*arg.MaxDistance))
	}

	ds := dialect.
		From(filtered.As("filtered_routes")).
		InnerJoin(goqu.T("users"), goqu.On(goqu.I("filtered_routes.user_id").Eq(goqu.I("users.id")))).
		Select(
			"filtered_routes.id",
			"filtered_routes.user_id",
			"filtered_routes.name",
			"filtered_routes.description",
			"filtered_routes.highlighted_photo_id",
			"filtered_routes.distance",
			"filtered_routes.duration",
			"filtered_routes.elevation_gain",
			"filtered_routes.elevation_loss",
			"filtered_routes.path_geom",
			"filtered_routes.bbox",
			"filtered_routes.first_point",
			"filtered_routes.last_point",
			"filtered_routes.polyline",
			"filtered_routes.created_at",
			"filtered_routes.updated_at",
			"filtered_routes.visibility",
			"filtered_routes.total_count",
			goqu.I("users.name").As("user_name"),
		).
		Limit(uint(arg.LimitCount)).
		Offset(uint(arg.OffsetCount)).
		Prepared(true)

	if hasLocation {
		ds = ds.Order(goqu.L(
			"ST_Distance(filtered_routes.first_point::geography, ST_GeomFromEWKB(?)::geography)",
			arg.Location,
		).Asc())
	}

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := q.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ExploreRoutesRow
	for rows.Next() {
		var i ExploreRoutesRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Description,
			&i.HighlightedPhotoID,
			&i.Distance,
			&i.Duration,
			&i.ElevationGain,
			&i.ElevationLoss,
			&i.PathGeom,
			&i.Bbox,
			&i.FirstPoint,
			&i.LastPoint,
			&i.Polyline,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Visibility,
			&i.TotalCount,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

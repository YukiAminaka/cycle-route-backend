package dbgen

import (
	"context"

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

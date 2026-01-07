package route

import (
	"context"
)

type IRouteRepository interface {
	GetRoutesByUserID(ctx context.Context, userID string) ([]*Route, error)
	CountRoutesByUserID(ctx context.Context, userID string) (int64, error)
	GetRouteByID(ctx context.Context, id string) (*Route, error)
	SaveRoute(ctx context.Context, route *Route) error
	DeleteRoute(ctx context.Context, id string) error
	UpdateRoute(ctx context.Context, route *Route) error
}

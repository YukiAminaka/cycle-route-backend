package route

import (
	"context"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
)

type IRouteRepository interface {
	GetRoutesByUserID(ctx context.Context, userID user.UserID) ([]*Route, error)
	CountRoutesByUserID(ctx context.Context, userID user.UserID) (int64, error)
	GetRouteByID(ctx context.Context, id RouteID) (*Route, error)
	SaveRoute(ctx context.Context, route *Route) error
	DeleteRoute(ctx context.Context, id RouteID) error
	UpdateRoute(ctx context.Context, route *Route) error
}

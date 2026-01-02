package route

import "context"

type IRouteRepository interface {
	GetRouteByID(ctx context.Context, id RouteID) (*Route, error)
	SaveRoute(ctx context.Context, route *Route) error
	DeleteRoute(ctx context.Context, id RouteID) error
	UpdateRoute(ctx context.Context, route *Route) error
}

package route

import (
	"context"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	"github.com/paulmach/orb"
)

type IGetRouteUsecase interface {
	GetRouteByID(ctx context.Context, routeID string) (*GetRouteUseCaseOutputDto, error)
	GetRoutesByUserID(ctx context.Context, userID string) ([]*GetRouteUseCaseOutputDto, error)
}

type getRouteUsecase struct {
	routeRepo routeDomain.IRouteRepository
}

func NewGetRouteUsecase(routeRepo routeDomain.IRouteRepository) IGetRouteUsecase {
	return &getRouteUsecase{
		routeRepo: routeRepo,
	}
}

type CoursePointOutput struct {
	ID            string
	StepOrder     int32
	SegDistM      *float64
	CumDistM      *float64
	Duration      *float64
	Instruction   *string
	RoadName      *string
	ManeuverType  *string
	Modifier      *string
	Location      *orb.Point
	BearingBefore *int32
	BearingAfter  *int32
}

type WaypointOutput struct {
	ID       string
	Location orb.Point
}

type GetRouteUseCaseOutputDto struct {
	ID                 string
	UserID             string
	Name               string
	Description        string
	HighlightedPhotoID *int64
	Distance           float64
	Duration           int32
	ElevationGain      float64
	ElevationLoss      float64
	PathGeom           orb.LineString
	Bbox               orb.Polygon
	FirstPoint         orb.Point
	LastPoint          orb.Point
	Visibility         int16
	CoursePoints       []CoursePointOutput
	Waypoints          []WaypointOutput
}

func (u *getRouteUsecase) GetRouteByID(ctx context.Context, routeID string) (*GetRouteUseCaseOutputDto, error) {
	route, err := u.routeRepo.GetRouteByID(ctx, routeID)
	if err != nil {
		return nil, err
	}

	return u.convertToOutputDto(route), nil
}

func (u *getRouteUsecase) GetRoutesByUserID(ctx context.Context, userID string) ([]*GetRouteUseCaseOutputDto, error) {
	routes, err := u.routeRepo.GetRoutesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	outputs := make([]*GetRouteUseCaseOutputDto, len(routes))
	for i, route := range routes {
		outputs[i] = u.convertToOutputDto(route)
	}

	return outputs, nil
}

func (u *getRouteUsecase) convertToOutputDto(route *routeDomain.Route) *GetRouteUseCaseOutputDto {
	// CoursePointsの変換
	coursePoints := make([]CoursePointOutput, len(route.CoursePoints()))
	for i, cp := range route.CoursePoints() {
		var location *orb.Point
		if cp.Location() != nil {
			if point, ok := cp.Location().Geometry.(orb.Point); ok {
				location = &point
			}
		}

		coursePoints[i] = CoursePointOutput{
			ID:            cp.ID(),
			StepOrder:     cp.StepOrder(),
			SegDistM:      cp.SegDistM(),
			CumDistM:      cp.CumDistM(),
			Duration:      cp.Duration(),
			Instruction:   cp.Instruction(),
			RoadName:      cp.RoadName(),
			ManeuverType:  cp.ManeuverType(),
			Modifier:      cp.Modifier(),
			Location:      location,
			BearingBefore: cp.BearingBefore(),
			BearingAfter:  cp.BearingAfter(),
		}
	}

	// Waypointsの変換
	waypoints := make([]WaypointOutput, len(route.Waypoints()))
	for i, wp := range route.Waypoints() {
		waypoints[i] = WaypointOutput{
			ID:       wp.ID(),
			Location: wp.Location().Geometry.(orb.Point),
		}
	}

	return &GetRouteUseCaseOutputDto{
		ID:                 route.ID(),
		UserID:             route.UserID(),
		Name:               route.Name(),
		Description:        route.Description(),
		HighlightedPhotoID: route.HighlightedPhotoID(),
		Distance:           route.Distance(),
		Duration:           route.Duration(),
		ElevationGain:      route.ElevationGain(),
		ElevationLoss:      route.ElevationLoss(),
		PathGeom:           route.PathGeom().Geometry.(orb.LineString),//DBから取得したデータはドメイン層で定義した型を満たしていると仮定。外部システムからDBに直接書き込む可能性がある場合はバリデーションが必要
		Bbox:               route.Bbox().Geometry.(orb.Polygon),
		FirstPoint:         route.FirstPoint().Geometry.(orb.Point),
		LastPoint:          route.LastPoint().Geometry.(orb.Point),
		Visibility:         route.Visibility(),
		CoursePoints:       coursePoints,
		Waypoints:          waypoints,
	}
}

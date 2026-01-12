package route

import (
	"context"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/repository"
	"github.com/YukiAminaka/cycle-route-backend/internal/usecase/transaction"
	"github.com/paulmach/orb"
)

type ICreateRouteUsecase interface {
	CreateRoute(ctx context.Context, dto CreateRouteUseCaseInputDto) (*CreateRouteUseCaseOutputDto, error)
}

type createRouteUsecase struct {
	userRepository user.IUserRepository
	txManager      transaction.TransactionManager
}

func NewCreateRouteUsecase(userRepository user.IUserRepository, txManager transaction.TransactionManager) ICreateRouteUsecase {
	return &createRouteUsecase{
		userRepository: userRepository,
		txManager:      txManager,
	}
}

type CoursePointInput struct {
	SegDistM      *float64
	CumDistM      *float64
	Duration      *float64
	Instruction   *string
	RoadName      *string
	ManeuverType  *string
	Modifier      *string
	Location      orb.Point
	BearingBefore *int32
	BearingAfter  *int32
}

type WaypointInput struct {
	Location orb.Point
}

type CreateRouteUseCaseInputDto struct {
	KratosID           string
	Name               string
	Description        string
	HighlightedPhotoID *int64
	Distance           float64
	Duration           int32
	ElevationGain      float64
	ElevationLoss      float64
	PathGeom           orb.LineString
	FirstPoint         orb.Point
	LastPoint          orb.Point
	Visibility         int16
	CoursePoints       []CoursePointInput
	Waypoints          []WaypointInput
}

type CreateRouteUseCaseOutputDto struct {
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
}

func (u *createRouteUsecase) CreateRoute(ctx context.Context, dto CreateRouteUseCaseInputDto) (*CreateRouteUseCaseOutputDto, error) {
	// KratosIDからユーザー情報を取得
	userEntity, err := u.userRepository.GetUserByKratosID(ctx, dto.KratosID)
	if err != nil {
		return nil, err
	}

	// ドメインモデルの作成
	route, err := routeDomain.NewRoute(
		userEntity.ID().String(),
		dto.Name,
		dto.Description,
		dto.HighlightedPhotoID,
		dto.Distance,
		dto.Duration,
		dto.ElevationGain,
		dto.ElevationLoss,
		routeDomain.Geometry{Geometry: dto.PathGeom},
		routeDomain.Geometry{Geometry: dto.FirstPoint},
		routeDomain.Geometry{Geometry: dto.LastPoint},
		dto.Visibility,
	)
	if err != nil {
		return nil, err
	}

	// Waypointsの追加
	for _, wp := range dto.Waypoints {
		if err := route.AddWaypoint(routeDomain.Geometry{Geometry: wp.Location}); err != nil {
			return nil, err
		}
	}

	// CoursePointsの追加
	for _, cp := range dto.CoursePoints {
		location := routeDomain.Geometry{Geometry: cp.Location}
		if err := route.AddCoursePoint(
			cp.SegDistM,
			cp.CumDistM,
			cp.Duration,
			cp.Instruction,
			cp.RoadName,
			cp.ManeuverType,
			cp.Modifier,
			&location,
			cp.BearingBefore,
			cp.BearingAfter,
		); err != nil {
			return nil, err
		}
	}

	// トランザクション内でリポジトリ操作を実行
	err = u.txManager.RunInTransaction(ctx, func(q *dbgen.Queries) error {
		// トランザクション用のQueriesでリポジトリを作成
		routeRepo := repository.NewRouteRepository(q)
		return routeRepo.SaveRoute(ctx, route)
	})

	if err != nil {
		return nil, err
	}

	// 出力DTOの作成
	return &CreateRouteUseCaseOutputDto{
		ID:                 route.ID(),
		UserID:             route.UserID(),
		Name:               route.Name(),
		Description:        route.Description(),
		HighlightedPhotoID: route.HighlightedPhotoID(),
		Distance:           route.Distance(),
		Duration:           route.Duration(),
		ElevationGain:      route.ElevationGain(),
		ElevationLoss:      route.ElevationLoss(),
		PathGeom:           route.PathGeom().Geometry.(orb.LineString),
		Bbox:               route.Bbox().Geometry.(orb.Polygon),
		FirstPoint:         route.FirstPoint().Geometry.(orb.Point),
		LastPoint:          route.LastPoint().Geometry.(orb.Point),
		Visibility:         route.Visibility(),
	}, nil
}


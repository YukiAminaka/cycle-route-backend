package route

import (
	"context"
	"errors"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/repository"
	"github.com/YukiAminaka/cycle-route-backend/internal/usecase/transaction"
	"github.com/paulmach/orb"
)

type IUpdateRouteUsecase interface {
	UpdateRoute(ctx context.Context,dto UpdateRouteUseCaseInputDto) error
}

type updateRouteUsecase struct {
	txManager transaction.TransactionManager
	routeRepo routeDomain.IRouteRepository
}

func NewUpdateRouteUsecase(txManager transaction.TransactionManager, routeRepo routeDomain.IRouteRepository) IUpdateRouteUsecase {
	return &updateRouteUsecase{
		txManager: txManager,
		routeRepo: routeRepo,
	}
}

type UpdatedCoursePointInput struct {
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

type UpdatedWaypointInput struct {
	Location orb.Point
}

type UpdateRouteUseCaseInputDto struct {
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
	FirstPoint         orb.Point
	LastPoint          orb.Point
	Visibility         int16
	CoursePoints       []UpdatedCoursePointInput
	Waypoints          []UpdatedWaypointInput
}

func (u *updateRouteUsecase) UpdateRoute(ctx context.Context, dto UpdateRouteUseCaseInputDto) error {
	// 既存のルートを取得
	route, err := u.routeRepo.GetRouteByID(ctx, dto.ID)
	if err != nil {
		return err
	}

	// 権限確認: ルートの所有者とリクエストのユーザーIDが一致するか
	if route.UserID() != dto.UserID {
		return errors.New("unauthorized: user does not own the route")
	}

	// 基本情報の更新
	if err := route.UpdateBasicInfo(
		dto.Name,
		dto.Description,
		dto.HighlightedPhotoID,
		dto.Visibility,
	); err != nil {
		return err
	}

	// ジオメトリ情報の更新
	if err := route.UpdateRouteGeometry(
		dto.Distance,
		dto.Duration,
		dto.ElevationGain,
		dto.ElevationLoss,
		routeDomain.Geometry{Geometry: dto.PathGeom},
		routeDomain.Geometry{Geometry: dto.FirstPoint},
		routeDomain.Geometry{Geometry: dto.LastPoint},
	); err != nil {
		return err
	}

	// コースポイントとウェイポイントをクリアして再設定
	route.ClearCoursePointsAndWaypoints()

	// Waypointsの追加
	for _, wp := range dto.Waypoints {
		if err := route.AddWaypoint(routeDomain.Geometry{Geometry: wp.Location}); err != nil {
			return err
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
			return err
		}
	}

	// トランザクション内でリポジトリ操作を実行
	err = u.txManager.RunInTransaction(ctx, func(q *dbgen.Queries) error {
		// トランザクション用のQueriesでリポジトリを作成
		routeRepo := repository.NewRouteRepository(q)
		return routeRepo.SaveRoute(ctx, route)
	})

	if err != nil {
		return err
	}

	return nil
}


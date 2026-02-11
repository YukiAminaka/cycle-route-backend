package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type routeRepositoryImpl struct {
	queries *dbgen.Queries
}

// ルートリポジトリの実装
func NewRouteRepository(queries *dbgen.Queries) route.IRouteRepository {
	return &routeRepositoryImpl{queries: queries}
}

func (r *routeRepositoryImpl) GetRouteByID(ctx context.Context, id string) (*route.Route, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid route id: %w", err)
	}

	rd, err := r.queries.GetRouteByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("route not found")
		}
		return nil, err
	}

	routeModel, err := route.ReconstructRoute(
		rd.ID.String(),
		rd.UserID.String(),
		rd.Name,
		rd.Description,
		rd.HighlightedPhotoID,
		rd.Distance,
		rd.Duration,
		rd.ElevationGain,
		rd.ElevationLoss,
		route.Geometry{Geometry: rd.PathGeom.Geometry},
		route.Geometry{Geometry: rd.Bbox.Geometry},
		route.Geometry{Geometry: rd.FirstPoint.Geometry},
		route.Geometry{Geometry: rd.LastPoint.Geometry},
		rd.Polyline,
		rd.Visibility,
		rd.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		rd.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	)
	if err != nil {
		return nil, err
	}

	// コースポイントを取得
	coursePointsData, err := r.queries.GetCoursePointsByRouteID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get course points: %w", err)
	}

	// CoursePoint構造体へのポインタのスライスを作成
	coursePoints := make([]*route.CoursePoint, 0, len(coursePointsData))
	for _, cp := range coursePointsData {
		var location *route.Geometry
		if cp.Location != nil {
			location = &route.Geometry{Geometry: cp.Location.Geometry}
		}
		coursePoints = append(coursePoints, route.ReconstructCoursePoint(
			cp.ID.String(),
			cp.RouteID.String(),
			cp.StepOrder,
			cp.SegDistM,
			cp.CumDistM,
			cp.Duration,
			cp.Instruction,
			cp.RoadName,
			cp.ManeuverType,
			cp.Modifier,
			location,
			cp.BearingBefore,
			cp.BearingAfter,
		))
	}
	routeModel.SetCoursePoints(coursePoints)

	// ウェイポイントを取得
	waypointsData, err := r.queries.GetWaypointsByRouteID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get waypoints: %w", err)
	}

	// Waypoint構造体へのポインタのスライスを作成
	waypoints := make([]*route.Waypoint, 0, len(waypointsData))
	for _, wp := range waypointsData {
		var location route.Geometry
		if wp.Location != nil {
			location = route.Geometry{Geometry: wp.Location.Geometry}
		}
		waypoints = append(waypoints, route.ReconstructWaypoint(
			wp.ID.String(),
			wp.RouteID.String(),
			location,
		))
	}
	routeModel.SetWaypoints(waypoints)

	return routeModel, nil
}

func (r *routeRepositoryImpl) GetRoutesByUserID(ctx context.Context, userID string) ([]*route.Route, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	rows, err := r.queries.GetRoutesByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := make([]*route.Route, 0, len(rows))
	for _, rd := range rows {
		routeModel, err := route.ReconstructRoute(
			rd.ID.String(),
			rd.UserID.String(),
			rd.Name,
			rd.Description,
			rd.HighlightedPhotoID,
			rd.Distance,
			rd.Duration,
			rd.ElevationGain,
			rd.ElevationLoss,
			route.Geometry{Geometry: rd.PathGeom.Geometry},
			route.Geometry{Geometry: rd.Bbox.Geometry},
			route.Geometry{Geometry: rd.FirstPoint.Geometry},
			route.Geometry{Geometry: rd.LastPoint.Geometry},
			rd.Polyline,
			rd.Visibility,
			rd.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			rd.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		)
		if err != nil {
			return nil, err
		}
		result = append(result, routeModel)
	}
	return result, nil
}

func (r *routeRepositoryImpl) CountRoutesByUserID(ctx context.Context, userID string) (int64, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user id: %w", err)
	}
	count, err := r.queries.CountRoutesByUserID(ctx, uid)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *routeRepositoryImpl) SaveRoute(ctx context.Context, rt *route.Route) error {
	// ルートIDとユーザーIDをUUIDに変換
	routeID, err := uuid.Parse(rt.ID())
	if err != nil {
		return fmt.Errorf("invalid route id: %w", err)
	}
	userID, err := uuid.Parse(rt.UserID())
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	// path_geomからbboxを自動計算
	bbox := CalculateBbox(rt.PathGeom().Geometry)

	// ルート本体を保存
	err = r.queries.CreateRoute(ctx, dbgen.CreateRouteParams{
		ID:                 routeID,
		UserID:             userID,
		Name:               rt.Name(),
		Description:        rt.Description(),
		HighlightedPhotoID: rt.HighlightedPhotoID(),
		Distance:           rt.Distance(),
		Duration:           rt.Duration(),
		ElevationGain:      rt.ElevationGain(),
		ElevationLoss:      rt.ElevationLoss(),
		PathGeom:           dbgen.OrbGeometry{Geometry: rt.PathGeom().Geometry},
		Bbox:               bbox,
		FirstPoint:         dbgen.OrbGeometry{Geometry: rt.FirstPoint().Geometry},
		LastPoint:          dbgen.OrbGeometry{Geometry: rt.LastPoint().Geometry},
		Visibility:         rt.Visibility(),
	})
	if err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}

	// コースポイントを保存
	for _, cp := range rt.CoursePoints() {
		cpID, err := uuid.Parse(cp.ID())
		if err != nil {
			return fmt.Errorf("invalid course point id: %w", err)
		}

		var location *dbgen.OrbGeometry
		if cp.Location() != nil {
			location = &dbgen.OrbGeometry{Geometry: cp.Location().Geometry}
		}

		err = r.queries.CreateCoursePoint(ctx, dbgen.CreateCoursePointParams{
			ID:            cpID,
			RouteID:       routeID,
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
		})
		if err != nil {
			return fmt.Errorf("failed to create course point: %w", err)
		}
	}

	// ウェイポイントを保存
	for _, wp := range rt.Waypoints() {
		wpID, err := uuid.Parse(wp.ID())
		if err != nil {
			return fmt.Errorf("invalid waypoint id: %w", err)
		}

		err = r.queries.CreateWaypoint(ctx, dbgen.CreateWaypointParams{
			ID:       wpID,
			RouteID:  routeID,
			Location: &dbgen.OrbGeometry{Geometry: wp.Location().Geometry},
		})
		if err != nil {
			return fmt.Errorf("failed to create waypoint: %w", err)
		}
	}

	return nil
}

func (r *routeRepositoryImpl) DeleteRoute(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid route id: %w", err)
	}

	// ルートを削除（カスケード削除でコースポイントとウェイポイントも削除される）
	_, err = r.queries.DeleteRoute(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("route not found")
		}
		return fmt.Errorf("failed to delete route: %w", err)
	}

	return nil
}

func (r *routeRepositoryImpl) UpdateRoute(ctx context.Context, rt *route.Route) error {
	// ルートIDをUUIDに変換
	routeID, err := uuid.Parse(rt.ID())
	if err != nil {
		return fmt.Errorf("invalid route id: %w", err)
	}

	// path_geomからbboxを自動計算
	bbox := CalculateBbox(rt.PathGeom().Geometry)

	// ルート本体を更新
	err = r.queries.UpdateRoute(ctx, dbgen.UpdateRouteParams{
		ID:                 routeID,
		Name:               rt.Name(),
		Description:        rt.Description(),
		HighlightedPhotoID: rt.HighlightedPhotoID(),
		Distance:           rt.Distance(),
		Duration:           rt.Duration(),
		ElevationGain:      rt.ElevationGain(),
		ElevationLoss:      rt.ElevationLoss(),
		PathGeom:           dbgen.OrbGeometry{Geometry: rt.PathGeom().Geometry},
		Bbox:               bbox,
		FirstPoint:         dbgen.OrbGeometry{Geometry: rt.FirstPoint().Geometry},
		LastPoint:          dbgen.OrbGeometry{Geometry: rt.LastPoint().Geometry},
		Visibility:         rt.Visibility(),
	})
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}

	// 既存のコースポイントとウェイポイントを削除
	err = r.queries.DeleteCoursePointsByRouteID(ctx, routeID)
	if err != nil {
		return fmt.Errorf("failed to delete course points: %w", err)
	}

	err = r.queries.DeleteWaypointsByRouteID(ctx, routeID)
	if err != nil {
		return fmt.Errorf("failed to delete waypoints: %w", err)
	}

	// 新しいコースポイントを保存
	for _, cp := range rt.CoursePoints() {
		cpID, err := uuid.Parse(cp.ID())
		if err != nil {
			return fmt.Errorf("invalid course point id: %w", err)
		}

		var location *dbgen.OrbGeometry
		if cp.Location() != nil {
			location = &dbgen.OrbGeometry{Geometry: cp.Location().Geometry}
		}

		err = r.queries.CreateCoursePoint(ctx, dbgen.CreateCoursePointParams{
			ID:            cpID,
			RouteID:       routeID,
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
		})
		if err != nil {
			return fmt.Errorf("failed to create course point: %w", err)
		}
	}

	// 新しいウェイポイントを保存
	for _, wp := range rt.Waypoints() {
		wpID, err := uuid.Parse(wp.ID())
		if err != nil {
			return fmt.Errorf("invalid waypoint id: %w", err)
		}

		err = r.queries.CreateWaypoint(ctx, dbgen.CreateWaypointParams{
			ID:       wpID,
			RouteID:  routeID,
			Location: &dbgen.OrbGeometry{Geometry: wp.Location().Geometry},
		})
		if err != nil {
			return fmt.Errorf("failed to create waypoint: %w", err)
		}
	}

	return nil
}
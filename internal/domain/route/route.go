package route

import (
	"time"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
)

// ルートのコースポイント
type CoursePoint struct {
	id            string
	routeID       string
	stepOrder     int32
	segDistM      *float64       // セグメント距離(メートル)
	cumDistM      *float64       // 累積距離(メートル)
	duration      *float64       // 所要時間(秒)
	instruction   *string        // ナビゲーション指示
	roadName      *string        // 道路名
	maneuverType  *string        // 操作タイプ（右折、左折など）
	modifier      *string        // 操作修飾子
	location      *user.Geometry // 地点の位置情報
	bearingBefore *int32         // 進入前の方位角
	bearingAfter  *int32         // 進入後の方位角
}

func NewCoursePoint(
	id string,
	routeID string,
	stepOrder int32,
	segDistM *float64,
	cumDistM *float64,
	duration *float64,
	instruction *string,
	roadName *string,
	maneuverType *string,
	modifier *string,
	location *user.Geometry,
	bearingBefore *int32,
	bearingAfter *int32) (*CoursePoint, error) {
		return &CoursePoint{
			id:            id,
			routeID:       routeID,
			stepOrder:     stepOrder,
			segDistM:      segDistM,
			cumDistM:      cumDistM,
			duration:      duration,
			instruction:   instruction,
			roadName:      roadName,
			maneuverType:  maneuverType,
			modifier:      modifier,
			location:      location,
			bearingBefore: bearingBefore,
			bearingAfter:  bearingAfter,
		}, nil
}

func (cp *CoursePoint) ID() string {
	return cp.id
}

func (cp *CoursePoint) RouteID() string {
	return cp.routeID
}

func (cp *CoursePoint) StepOrder() int32 {
	return cp.stepOrder
}

func (cp *CoursePoint) SegDistM() *float64 {
	return cp.segDistM
}

func (cp *CoursePoint) CumDistM() *float64 {
	return cp.cumDistM
}

func (cp *CoursePoint) Duration() *float64 {
	return cp.duration
}

func (cp *CoursePoint) Instruction() *string {
	return cp.instruction
}

func (cp *CoursePoint) RoadName() *string {
	return cp.roadName
}

func (cp *CoursePoint) ManeuverType() *string {
	return cp.maneuverType
}

func (cp *CoursePoint) Modifier() *string {
	return cp.modifier
}

func (cp *CoursePoint) Location() *user.Geometry {
	return cp.location
}

func (cp *CoursePoint) BearingBefore() *int32 {
	return cp.bearingBefore
}

func (cp *CoursePoint) BearingAfter() *int32 {
	return cp.bearingAfter
}

type Route struct {
	id                 string
	author             string
	name               string
	description        string
	highlightedPhotoID *int64
	hasCoursePoints    bool
	distance           float64
	duration           *int32
	elevationGain      float64
	elevationLoss      float64
	pathGeom           user.Geometry
	bbox               user.Geometry
	firstPoint         user.Geometry
	lastPoint          user.Geometry
	createdAt          time.Time
	visibility         int16
}

func NewRoute(
	id string,
	author string,
	name string,
	description string,
	highlightedPhotoID *int64,
	hasCoursePoints bool,
	distance float64,
	duration *int32,
	elevationGain float64,
	elevationLoss float64,
	pathGeom user.Geometry,
	bbox user.Geometry,
	firstPoint user.Geometry,
	lastPoint user.Geometry,
	createdAt time.Time,
	visibility int16) (*Route, error) {
		return &Route{
			id:                 id,
			author: 		   author,
			name:               name,
			description:        description,
			highlightedPhotoID: highlightedPhotoID,
			hasCoursePoints:    hasCoursePoints,
			distance:           distance,
			duration:           duration,
			elevationGain:      elevationGain,
			elevationLoss:      elevationLoss,
			pathGeom:           pathGeom,
			bbox:               bbox,
			firstPoint:        firstPoint,
			lastPoint:         lastPoint,
			createdAt:          createdAt,
			visibility:         visibility,
		}, nil
}

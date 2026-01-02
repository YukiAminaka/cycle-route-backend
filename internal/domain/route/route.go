package route

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)


type RouteID string
type CoursePointID string
type WaypointID string


func NewRouteID() RouteID {
	uuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return RouteID(uuid.String())
}

func NewCoursePointID() CoursePointID {
	uuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return CoursePointID(uuid.String())
}

func NewWaypointID() WaypointID {
	uuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return WaypointID(uuid.String())
}

func (id RouteID) String() string {
	return string(id)
}

func (id CoursePointID) String() string {
	return string(id)
}

func (id WaypointID) String() string {
	return string(id)
}

// Geometry はドメイン層でのジオメトリ型（PostGISのgeometryに対応）
type Geometry struct {
	orb.Geometry
}

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
	location      *Geometry // 地点の位置情報
	bearingBefore *int32         // 進入前の方位角
	bearingAfter  *int32         // 進入後の方位角
}

// CoursePointは集約ルート(Route)を通してのみ作成されます

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

func (cp *CoursePoint) Location() *Geometry {
	return cp.location
}

func (cp *CoursePoint) BearingBefore() *int32 {
	return cp.bearingBefore
}

func (cp *CoursePoint) BearingAfter() *int32 {
	return cp.bearingAfter
}

type Waypoint struct {
	id        string
	routeID   string
	location  Geometry
}

// Waypointは集約ルート(Route)を通してのみ作成されます

func (w *Waypoint) ID() string {
	return w.id
}

func (w *Waypoint) RouteID() string {
	return w.routeID
}

func (w *Waypoint) Location() Geometry {
	return w.location
}


// Route は集約ルート
type Route struct {
    id                 string
    userID             string
    name               string
    description        string
    highlightedPhotoID *int64
    distance           float64
    duration           int32
    elevationGain      float64
    elevationLoss      float64
    pathGeom           Geometry
    bbox               Geometry
    firstPoint         Geometry
    lastPoint          Geometry
    visibility         int16
    
    // 集約内のエンティティコレクション
    coursePoints []*CoursePoint
    waypoints    []*Waypoint
}

// NewRoute は新しいルートを作成（Mapbox Direction APIからの情報を基に作成）
func NewRoute(
    userID string,
    name string,
    description string,
    highlightedPhotoID *int64,
    distance float64,
    duration int32,
    elevationGain float64,
    elevationLoss float64,
    pathGeom Geometry,
    bbox Geometry,
    firstPoint Geometry,
    lastPoint Geometry,
    visibility int16) (*Route, error) {

    // バリデーション
    if userID == "" {
        return nil, errors.New("userID is required")
    }
    if name == "" {
        return nil, errors.New("name is required")
    }
    if pathGeom.Geometry == nil {
        return nil, errors.New("pathGeom is required")
    }
    if pathGeom.Geometry.GeoJSONType() != "LineString" {
        return nil, errors.New("pathGeom must be a LineString")
    }
    if bbox.Geometry == nil {
        return nil, errors.New("bbox is required")
    }
    if bbox.Geometry.GeoJSONType() != "Polygon" {
        return nil, errors.New("bbox must be a Polygon")
    }
    if firstPoint.Geometry == nil {
        return nil, errors.New("firstPoint is required")
    }
    if firstPoint.Geometry.GeoJSONType() != "Point" {
        return nil, errors.New("firstPoint must be a Point")
    }
    if lastPoint.Geometry == nil {
        return nil, errors.New("lastPoint is required")
    }
    if lastPoint.Geometry.GeoJSONType() != "Point" {
        return nil, errors.New("lastPoint must be a Point")
    }
    if distance < 0 {
        return nil, errors.New("distance must be non-negative")
    }
    if duration < 0 {
        return nil, errors.New("duration must be non-negative")
    }

    // IDの生成
    id := NewRouteID().String()

    return &Route{
        id:                 id,
        userID:             userID,
        name:               name,
        description:        description,
        highlightedPhotoID: highlightedPhotoID,
        distance:           distance,
        duration:           duration,
        elevationGain:      elevationGain,
        elevationLoss:      elevationLoss,
        pathGeom:           pathGeom,
        bbox:               bbox,
        firstPoint:         firstPoint,
        lastPoint:          lastPoint,
        visibility:         visibility,
        coursePoints:       []*CoursePoint{},
        waypoints:          []*Waypoint{},
    }, nil
}


// 集約ルートを通してのみCoursePointを追加
func (r *Route) AddCoursePoint(
    segDistM *float64,
    cumDistM *float64,
    duration *float64,
    instruction *string,
    roadName *string,
    maneuverType *string,
    modifier *string,
    location *Geometry,
    bearingBefore *int32,
    bearingAfter *int32) error {
    
    // ビジネスルール: stepOrderは自動採番
    nextOrder := int32(len(r.coursePoints))
    
    // バリデーション
    if location == nil {
        return errors.New("location is required")
    }

	if location.Geometry.GeoJSONType() != "Point" {
		return errors.New("location must be a Point geometry")
	}
    
    // IDは内部で生成（UUIDなど）
    id := NewCoursePointID().String()
    
    cp := &CoursePoint{
        id:            id,
        routeID:       r.id,
        stepOrder:     nextOrder,
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
    }
    
    r.coursePoints = append(r.coursePoints, cp)
    
    // ルート全体の距離・時間を再計算
    r.recalculateMetrics()
    
    return nil
}

// Waypointの追加も同様に集約ルート経由
func (r *Route) AddWaypoint(location Geometry) error {
    if location.Geometry == nil {
        return errors.New("location is required")
    }

	// 	GeometryがPoint型であることを確認
	if location.Geometry.GeoJSONType() != "Point" {
		return errors.New("location must be a Point geometry")
	}

    
    id := NewWaypointID().String()
    
    wp := &Waypoint{
        id:       id,
        routeID:  r.id,
        location: location,
    }
    
    r.waypoints = append(r.waypoints, wp)
    
    return nil
}

// CoursePointsを取得（イミュータブルなコピーを返す）
func (r *Route) CoursePoints() []*CoursePoint {
    // 防御的コピー
    points := make([]*CoursePoint, len(r.coursePoints))
    copy(points, r.coursePoints)
    return points
}

func (r *Route) Waypoints() []*Waypoint {
    points := make([]*Waypoint, len(r.waypoints))
    copy(points, r.waypoints)
    return points
}

// ビジネスロジック: ルート全体のメトリクスを再計算
func (r *Route) recalculateMetrics() {
    if len(r.coursePoints) == 0 {
        return
    }
    
    // 累積距離の再計算など
    totalDistance := 0.0
    totalDuration := 0.0
    
    for _, cp := range r.coursePoints {
        if cp.segDistM != nil {
            totalDistance += *cp.segDistM
        }
        if cp.duration != nil {
            totalDuration += *cp.duration
        }
    }
    
    r.distance = totalDistance
    r.duration = int32(totalDuration)
}


func ReconstructRoute(
    id string,
    userID string,
    name string,
    description string,
    highlightedPhotoID *int64,
    distance float64,
    duration int32,
    elevationGain float64,
    elevationLoss float64,
    pathGeom Geometry,
    bbox Geometry,
    firstPoint Geometry,
    lastPoint Geometry,
    createdAt time.Time,
    visibility int16,
    coursePoints []*CoursePoint,
    waypoints []*Waypoint) *Route {
    
    return &Route{
        id:                 id,
        userID:             userID,
        name:               name,
        description:        description,
        highlightedPhotoID: highlightedPhotoID,
        distance:           distance,
        duration:           duration,
        elevationGain:      elevationGain,
        elevationLoss:      elevationLoss,
        pathGeom:           pathGeom,
        bbox:               bbox,
        firstPoint:         firstPoint,
        lastPoint:          lastPoint,
        visibility:         visibility,
        coursePoints:       coursePoints,
        waypoints:          waypoints,
    }
}

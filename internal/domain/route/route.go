package route

import (
	"errors"

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
	segDistM      *float64  // セグメント距離(メートル)
	cumDistM      *float64  // 累積距離(メートル)
	duration      *float64  // 所要時間(秒)
	instruction   *string   // ナビゲーション指示
	roadName      *string   // 道路名
	maneuverType  *string   // 操作タイプ（右折、左折など）
	modifier      *string   // 操作修飾子
	location      *Geometry // 地点の位置情報
	bearingBefore *int32    // 進入前の方位角
	bearingAfter  *int32    // 進入後の方位角
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
	id       string
	routeID  string
	location Geometry
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
	createdAt          string
	updatedAt          string

	// 集約内のエンティティコレクション
	coursePoints []*CoursePoint
	waypoints    []*Waypoint
}

// newRoute は新しいルートを作成（Mapbox Direction APIからの情報を基に作成）
func newRoute(
	userID string,
	name string,
	description string,
	highlightedPhotoID *int64,
	distance float64,
	duration int32,
	elevationGain float64,
	elevationLoss float64,
	pathGeom Geometry,
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

	// bboxは空のGeometryで初期化（リポジトリ層でpathGeomから計算する）
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
		bbox:               Geometry{}, // 空のGeometry
		firstPoint:         firstPoint,
		lastPoint:          lastPoint,
		visibility:         visibility,
		coursePoints:       []*CoursePoint{},
		waypoints:          []*Waypoint{},
	}, nil
}

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
	firstPoint Geometry,
	lastPoint Geometry,
	visibility int16) (*Route, error) {
	return newRoute(
		userID,
		name,
		description,
		highlightedPhotoID,
		distance,
		duration,
		elevationGain,
		elevationLoss,
		pathGeom,
		firstPoint,
		lastPoint,
		visibility)
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

// Routeのゲッターメソッド
func (r *Route) ID() string {
	return r.id
}

func (r *Route) UserID() string {
	return r.userID
}

func (r *Route) Name() string {
	return r.name
}

func (r *Route) Description() string {
	return r.description
}

func (r *Route) HighlightedPhotoID() *int64 {
	return r.highlightedPhotoID
}

func (r *Route) Distance() float64 {
	return r.distance
}

func (r *Route) Duration() int32 {
	return r.duration
}

func (r *Route) ElevationGain() float64 {
	return r.elevationGain
}

func (r *Route) ElevationLoss() float64 {
	return r.elevationLoss
}

func (r *Route) PathGeom() Geometry {
	return r.pathGeom
}

func (r *Route) Bbox() Geometry {
	return r.bbox
}

func (r *Route) FirstPoint() Geometry {
	return r.firstPoint
}

func (r *Route) LastPoint() Geometry {
	return r.lastPoint
}

func (r *Route) Visibility() int16 {
	return r.visibility
}

func (r *Route) CreatedAt() string {
	return r.createdAt
}

func (r *Route) UpdatedAt() string {
	return r.updatedAt
}


// コースポイントとウェイポイントをクリア（更新時に使用）
func (r *Route) ClearCoursePointsAndWaypoints() {
	r.coursePoints = []*CoursePoint{}
	r.waypoints = []*Waypoint{}
}

// CoursePointを再構築（リポジトリ層からの復元用）
func ReconstructCoursePoint(
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
	location *Geometry,
	bearingBefore *int32,
	bearingAfter *int32) *CoursePoint {
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
	}
}

// Waypointを再構築（リポジトリ層からの復元用）
func ReconstructWaypoint(id string, routeID string, location Geometry) *Waypoint {
	return &Waypoint{
		id:       id,
		routeID:  routeID,
		location: location,
	}
}

// Routeを再構築（リポジトリ層からの復元用）
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
	visibility int16,
	createdAt string,
	updatedAt string,
	) (*Route, error) {
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
		createdAt:          createdAt,
		updatedAt:          updatedAt,
		coursePoints:       []*CoursePoint{},
		waypoints:          []*Waypoint{},
	}, nil
}

// CoursePointsを直接設定（リポジトリ層での復元用）
func (r *Route) SetCoursePoints(coursePoints []*CoursePoint) {
	r.coursePoints = coursePoints
}

// Waypointsを直接設定（リポジトリ層での復元用）
func (r *Route) SetWaypoints(waypoints []*Waypoint) {
	r.waypoints = waypoints
}

// 作成したルートの基本情報を更新する（名前、説明、写真など）
func (r *Route) UpdateBasicInfo(
	name string,
	description string,
	highlightedPhotoID *int64,
	visibility int16) error {

	if name == "" {
		return errors.New("name is required")
	}

	r.name = name
	r.description = description
	r.highlightedPhotoID = highlightedPhotoID
	r.visibility = visibility

	return nil
}

// ルートのジオメトリ情報を更新する（ルート編集時に使用）
func (r *Route) UpdateRouteGeometry(
	distance float64,
	duration int32,
	elevationGain float64,
	elevationLoss float64,
	pathGeom Geometry,
	firstPoint Geometry,
	lastPoint Geometry) error {

	// バリデーション
	if pathGeom.Geometry == nil {
		return errors.New("pathGeom is required")
	}
	if pathGeom.Geometry.GeoJSONType() != "LineString" {
		return errors.New("pathGeom must be a LineString")
	}
	if firstPoint.Geometry == nil {
		return errors.New("firstPoint is required")
	}
	if firstPoint.Geometry.GeoJSONType() != "Point" {
		return errors.New("firstPoint must be a Point")
	}
	if lastPoint.Geometry == nil {
		return errors.New("lastPoint is required")
	}
	if lastPoint.Geometry.GeoJSONType() != "Point" {
		return errors.New("lastPoint must be a Point")
	}
	if distance < 0 {
		return errors.New("distance must be non-negative")
	}
	if duration < 0 {
		return errors.New("duration must be non-negative")
	}

	r.distance = distance
	r.duration = duration
	r.elevationGain = elevationGain
	r.elevationLoss = elevationLoss
	r.pathGeom = pathGeom
	r.firstPoint = firstPoint
	r.lastPoint = lastPoint

	return nil
}

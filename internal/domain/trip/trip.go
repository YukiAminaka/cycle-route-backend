package trip

import (
	"errors"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type TripID string

func NewTripID() TripID {
	uuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return TripID(uuid.String())
}

func (id TripID) String() string {
	return string(id)
}

// Geometry はドメイン層でのジオメトリ型（PostGISのgeometryに対応）
type Geometry struct {
	orb.Geometry
}

// Trip は活動（ライド記録）の集約ルート
type Trip struct {
	id                 string
	userID             string
	name               string
	description        string
	visibility         int16
	highlightedPhotoID int64

	// 位置情報
	pathGeom   *Geometry
	firstPoint *Geometry
	lastPoint  *Geometry
	bboxGeom   *Geometry

	// 計測系
	distance      *float64
	duration      *int32
	movingTime    *int32
	elevationGain *float64
	elevationLoss *float64
	avgSpeed      *float64
	maxSpeed      *float64

	// センサー/パワー関連
	avgCad             *float64
	maxCad             *float64
	minCad             *float64
	maxHr              *int32
	minHr              *int32
	avgWatts           *float64
	maxWatts           *float64
	minWatts           *float64
	avgWattsEstimated  *bool
	avgPowerEstimated  *float64
	calories           *float64

	// 種別/状態
	isGPS        bool
	isStationary bool
	processed    bool

	// 時刻/タイムゾーン
	createdAt  string
	updatedAt  string
	deletedAt  *string
	departedAt *string
	timeZone   *string
	utcOffset  *int32

	// アクティビティ
	activityTypeID int32

	// 付帯情報
	pace       *float64
	movingPace *float64
}

func newTrip(
	userID string,
	name string,
	description string,
	visibility int16,
	activityTypeID int32,
) (*Trip, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	id := NewTripID().String()

	return &Trip{
		id:             id,
		userID:         userID,
		name:           name,
		description:    description,
		visibility:     visibility,
		activityTypeID: activityTypeID,
	}, nil
}

func NewTrip(
	userID string,
	name string,
	description string,
	visibility int16,
	activityTypeID int32,
) (*Trip, error) {
	return newTrip(userID, name, description, visibility, activityTypeID)
}

// ReconstructTrip はリポジトリ層からの復元用
func ReconstructTrip(
	id string,
	userID string,
	name string,
	description string,
	visibility int16,
	highlightedPhotoID int64,
	pathGeom *Geometry,
	firstPoint *Geometry,
	lastPoint *Geometry,
	bboxGeom *Geometry,
	distance *float64,
	duration *int32,
	movingTime *int32,
	elevationGain *float64,
	elevationLoss *float64,
	avgSpeed *float64,
	maxSpeed *float64,
	avgCad *float64,
	maxCad *float64,
	minCad *float64,
	maxHr *int32,
	minHr *int32,
	avgWatts *float64,
	maxWatts *float64,
	minWatts *float64,
	avgWattsEstimated *bool,
	avgPowerEstimated *float64,
	calories *float64,
	isGPS bool,
	isStationary bool,
	processed bool,
	createdAt string,
	updatedAt string,
	deletedAt *string,
	departedAt *string,
	timeZone *string,
	utcOffset *int32,
	activityTypeID int32,
	pace *float64,
	movingPace *float64,
) *Trip {
	return &Trip{
		id:                 id,
		userID:             userID,
		name:               name,
		description:        description,
		visibility:         visibility,
		highlightedPhotoID: highlightedPhotoID,
		pathGeom:           pathGeom,
		firstPoint:         firstPoint,
		lastPoint:          lastPoint,
		bboxGeom:           bboxGeom,
		distance:           distance,
		duration:           duration,
		movingTime:         movingTime,
		elevationGain:      elevationGain,
		elevationLoss:      elevationLoss,
		avgSpeed:           avgSpeed,
		maxSpeed:           maxSpeed,
		avgCad:             avgCad,
		maxCad:             maxCad,
		minCad:             minCad,
		maxHr:              maxHr,
		minHr:              minHr,
		avgWatts:           avgWatts,
		maxWatts:           maxWatts,
		minWatts:           minWatts,
		avgWattsEstimated:  avgWattsEstimated,
		avgPowerEstimated:  avgPowerEstimated,
		calories:           calories,
		isGPS:              isGPS,
		isStationary:       isStationary,
		processed:          processed,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
		deletedAt:          deletedAt,
		departedAt:         departedAt,
		timeZone:           timeZone,
		utcOffset:          utcOffset,
		activityTypeID:     activityTypeID,
		pace:               pace,
		movingPace:         movingPace,
	}
}

// UpdateBasicInfo は名前・説明・公開範囲などの基本情報を更新する
func (t *Trip) UpdateBasicInfo(
	name string,
	description string,
	visibility int16,
	highlightedPhotoID int64,
) error {
	t.name = name
	t.description = description
	t.visibility = visibility
	t.highlightedPhotoID = highlightedPhotoID
	return nil
}

// SetMetrics はGPSファイルや計測データから取得したメトリクスをセットする
func (t *Trip) SetMetrics(
	pathGeom *Geometry,
	firstPoint *Geometry,
	lastPoint *Geometry,
	bboxGeom *Geometry,
	distance *float64,
	duration *int32,
	movingTime *int32,
	elevationGain *float64,
	elevationLoss *float64,
	avgSpeed *float64,
	maxSpeed *float64,
	departedAt *string,
	timeZone *string,
	utcOffset *int32,
	pace *float64,
	movingPace *float64,
) error {
	if pathGeom != nil && pathGeom.Geometry != nil && pathGeom.Geometry.GeoJSONType() != "LineString" {
		return errors.New("pathGeom must be a LineString")
	}
	if firstPoint != nil && firstPoint.Geometry != nil && firstPoint.Geometry.GeoJSONType() != "Point" {
		return errors.New("firstPoint must be a Point")
	}
	if lastPoint != nil && lastPoint.Geometry != nil && lastPoint.Geometry.GeoJSONType() != "Point" {
		return errors.New("lastPoint must be a Point")
	}
	if distance != nil && *distance < 0 {
		return errors.New("distance must be non-negative")
	}
	if duration != nil && *duration < 0 {
		return errors.New("duration must be non-negative")
	}
	if movingTime != nil && *movingTime < 0 {
		return errors.New("movingTime must be non-negative")
	}

	t.pathGeom = pathGeom
	t.firstPoint = firstPoint
	t.lastPoint = lastPoint
	t.bboxGeom = bboxGeom
	t.distance = distance
	t.duration = duration
	t.movingTime = movingTime
	t.elevationGain = elevationGain
	t.elevationLoss = elevationLoss
	t.avgSpeed = avgSpeed
	t.maxSpeed = maxSpeed
	t.departedAt = departedAt
	t.timeZone = timeZone
	t.utcOffset = utcOffset
	t.pace = pace
	t.movingPace = movingPace
	t.isGPS = true
	return nil
}

// SetSensorData はセンサー・パワーデータをセットする
func (t *Trip) SetSensorData(
	avgCad *float64,
	maxCad *float64,
	minCad *float64,
	maxHr *int32,
	minHr *int32,
	avgWatts *float64,
	maxWatts *float64,
	minWatts *float64,
	avgWattsEstimated *bool,
	avgPowerEstimated *float64,
	calories *float64,
) {
	t.avgCad = avgCad
	t.maxCad = maxCad
	t.minCad = minCad
	t.maxHr = maxHr
	t.minHr = minHr
	t.avgWatts = avgWatts
	t.maxWatts = maxWatts
	t.minWatts = minWatts
	t.avgWattsEstimated = avgWattsEstimated
	t.avgPowerEstimated = avgPowerEstimated
	t.calories = calories
}

// ゲッター
func (t *Trip) ID() string               { return t.id }
func (t *Trip) UserID() string           { return t.userID }
func (t *Trip) Name() string             { return t.name }
func (t *Trip) Description() string      { return t.description }
func (t *Trip) Visibility() int16        { return t.visibility }
func (t *Trip) HighlightedPhotoID() int64 { return t.highlightedPhotoID }

func (t *Trip) PathGeom() *Geometry   { return t.pathGeom }
func (t *Trip) FirstPoint() *Geometry { return t.firstPoint }
func (t *Trip) LastPoint() *Geometry  { return t.lastPoint }
func (t *Trip) BboxGeom() *Geometry   { return t.bboxGeom }

func (t *Trip) Distance() *float64      { return t.distance }
func (t *Trip) Duration() *int32        { return t.duration }
func (t *Trip) MovingTime() *int32      { return t.movingTime }
func (t *Trip) ElevationGain() *float64 { return t.elevationGain }
func (t *Trip) ElevationLoss() *float64 { return t.elevationLoss }
func (t *Trip) AvgSpeed() *float64      { return t.avgSpeed }
func (t *Trip) MaxSpeed() *float64      { return t.maxSpeed }

func (t *Trip) AvgCad() *float64            { return t.avgCad }
func (t *Trip) MaxCad() *float64            { return t.maxCad }
func (t *Trip) MinCad() *float64            { return t.minCad }
func (t *Trip) MaxHr() *int32               { return t.maxHr }
func (t *Trip) MinHr() *int32               { return t.minHr }
func (t *Trip) AvgWatts() *float64          { return t.avgWatts }
func (t *Trip) MaxWatts() *float64          { return t.maxWatts }
func (t *Trip) MinWatts() *float64          { return t.minWatts }
func (t *Trip) AvgWattsEstimated() *bool    { return t.avgWattsEstimated }
func (t *Trip) AvgPowerEstimated() *float64 { return t.avgPowerEstimated }
func (t *Trip) Calories() *float64          { return t.calories }

func (t *Trip) IsGPS() bool        { return t.isGPS }
func (t *Trip) IsStationary() bool { return t.isStationary }
func (t *Trip) Processed() bool    { return t.processed }

func (t *Trip) CreatedAt() string   { return t.createdAt }
func (t *Trip) UpdatedAt() string   { return t.updatedAt }
func (t *Trip) DeletedAt() *string  { return t.deletedAt }
func (t *Trip) DepartedAt() *string { return t.departedAt }
func (t *Trip) TimeZone() *string   { return t.timeZone }
func (t *Trip) UtcOffset() *int32   { return t.utcOffset }

func (t *Trip) ActivityTypeID() int32 { return t.activityTypeID }
func (t *Trip) Pace() *float64        { return t.pace }
func (t *Trip) MovingPace() *float64  { return t.movingPace }

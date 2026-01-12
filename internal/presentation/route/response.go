package route

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type RouteResponse struct {
	Route RouteResponseModel `json:"route"`
}

type RouteListResponse struct {
	Routes []RouteResponseModel `json:"routes"`
}


type RouteResponseModel struct {
	ID                 string   `json:"id"`
	UserID             string   `json:"user_id"`
	UserName           string   `json:"user_name"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	HighlightedPhotoID *int64   `json:"highlighted_photo_id"`
	Distance           float64  `json:"distance"`
	Duration           int32    `json:"duration"`
	ElevationGain      float64  `json:"elevation_gain"`
	ElevationLoss      float64  `json:"elevation_loss"`
	Visibility         int16    `json:"visibility"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
	PathGeom           *string   `json:"path_geom,omitempty"`
	Bbox               *string   `json:"bbox,omitempty"`
	FirstPoint         *string   `json:"first_point,omitempty"`
	LastPoint          *string   `json:"last_point,omitempty"`
	CoursePoints       []CoursePointResponse `json:"course_points,omitempty"`
	Waypoints          []WaypointResponse    `json:"waypoints,omitempty"`
}


type CoursePointResponse struct {
	ID            string   `json:"id"`
	StepOrder     int32    `json:"step_order"`
	SegDistM      *float64 `json:"seg_dist_m,omitempty"`
	CumDistM      *float64 `json:"cum_dist_m,omitempty"`
	Duration      *float64 `json:"duration,omitempty"`
	Instruction   *string  `json:"instruction,omitempty"`
	RoadName      *string  `json:"road_name,omitempty"`
	ManeuverType  *string  `json:"maneuver_type,omitempty"`
	Modifier      *string  `json:"modifier,omitempty"`
	Location      *string  `json:"location,omitempty"`
	BearingBefore *int32   `json:"bearing_before,omitempty"`
	BearingAfter  *int32   `json:"bearing_after,omitempty"`
}

type WaypointResponse struct {
	ID       string `json:"id"`
	Location *string `json:"location"`
}

// geometryToGeoJSON はorb.Geometryをstringに変換する
func geometryToGeoJSON(geom orb.Geometry) *string {
	if geom == nil {
		return nil
	}
	fc := geojson.NewFeatureCollection()
	fc.Append(geojson.NewFeature(geom))
	b, _ := fc.MarshalJSON()
	result := string(b)
	return &result
}

// pointToGeoJSON は*orb.Pointをstringに変換する
func pointToGeoJSON(point *orb.Point) *string {
	if point == nil {
		return nil
	}
	return geometryToGeoJSON(*point)
}

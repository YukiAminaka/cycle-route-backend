package route

// フロントエンドでは、地理空間データは通常GeoJSON形式で扱われます
type CreateRouteRequest struct {
	Name               string               `json:"name" validate:"required,max=255"`
	Description        string               `json:"description" validate:"max=1000"`
	HighlightedPhotoID *int64               `json:"highlighted_photo_id"`
	Distance           float64              `json:"distance" validate:"required,min=0"`
	Duration           float64              `json:"duration" validate:"required,min=0"`
	ElevationGain      float64              `json:"elevation_gain" validate:"min=0"`
	ElevationLoss      float64              `json:"elevation_loss" validate:"min=0"`
	PathGeom           string               `json:"path_geom" validate:"required"`
	FirstPoint         string               `json:"first_point" validate:"required"`
	LastPoint          string               `json:"last_point" validate:"required"`
	Visibility         int16                `json:"visibility" validate:"required,min=0,max=2"`
	CoursePoints       []CoursePointRequest `json:"course_points"`
	Waypoints          []WaypointRequest    `json:"waypoints"`
}

type UpdateRouteRequest struct {
	Name               string               `json:"name" validate:"required,max=255"`
	Description        string               `json:"description" validate:"max=1000"`
	HighlightedPhotoID *int64               `json:"highlighted_photo_id"`
	Distance           float64              `json:"distance" validate:"required,min=0"`
	Duration           float64              `json:"duration" validate:"required,min=0"`
	ElevationGain      float64              `json:"elevation_gain" validate:"min=0"`
	ElevationLoss      float64              `json:"elevation_loss" validate:"min=0"`
	PathGeom           string               `json:"path_geom" validate:"required"`
	FirstPoint         string               `json:"first_point" validate:"required"`
	LastPoint          string               `json:"last_point" validate:"required"`
	Visibility         int16                `json:"visibility" validate:"required,min=0,max=2"`
	CoursePoints       []CoursePointRequest `json:"course_points"`
	Waypoints          []WaypointRequest    `json:"waypoints"`
}

type CoursePointRequest struct {
	SegDistM      *float64 `json:"seg_dist_m"`
	CumDistM      *float64 `json:"cum_dist_m"`
	Duration      *float64 `json:"duration"`
	Instruction   *string  `json:"instruction"`
	RoadName      *string  `json:"road_name"`
	ManeuverType  *string  `json:"maneuver_type"`
	Modifier      *string  `json:"modifier"`
	Location      string   `json:"location" validate:"required"`
	BearingBefore *int32   `json:"bearing_before"`
	BearingAfter  *int32   `json:"bearing_after"`
}

type WaypointRequest struct {
	Location string `json:"location" validate:"required"`
}
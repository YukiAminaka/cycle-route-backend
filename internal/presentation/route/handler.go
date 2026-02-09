package route

import (
	"errors"

	"github.com/YukiAminaka/cycle-route-backend/internal/pkg/geojson"
	"github.com/YukiAminaka/cycle-route-backend/internal/pkg/geometry"
	"github.com/YukiAminaka/cycle-route-backend/internal/presentation/response"
	routeUsecase "github.com/YukiAminaka/cycle-route-backend/internal/usecase/route"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	createRouteUsecase routeUsecase.ICreateRouteUsecase
	getRouteUsecase    routeUsecase.IGetRouteUsecase
	updateRouteUsecase routeUsecase.IUpdateRouteUsecase
	deleteRouteUsecase routeUsecase.IDeleteRouteUsecase
}


func NewHandler(
	createRouteUsecase routeUsecase.ICreateRouteUsecase,
	getRouteUsecase    routeUsecase.IGetRouteUsecase,
	updateRouteUsecase routeUsecase.IUpdateRouteUsecase,
	deleteRouteUsecase routeUsecase.IDeleteRouteUsecase,
) *Handler {
	return &Handler{
		createRouteUsecase: createRouteUsecase,
		getRouteUsecase:    getRouteUsecase,
		updateRouteUsecase: updateRouteUsecase,
		deleteRouteUsecase: deleteRouteUsecase,
	}
}

// CreateRoute godoc
//	@Summary	ルートを作成する
//	@Tags		routes
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		request	body		CreateRouteRequest	true	"Create Route Request"
//	@Success	201		{object}	RouteResponse
//	@Failure	400		{object}	response.ErrorResponse
//	@Failure	401		{object}	response.ErrorResponse
//	@Failure	500		{object}	response.ErrorResponse
//	@Router		/routes [post]
func (h *Handler) CreateRoute(c *gin.Context) {
	// 認証ミドルウェアからKratosIDを取得
	kratosIDValue, exists := c.Get("kratos_id")
	if !exists {
		response.ReturnStatusUnauthorized(c, errors.New("user not authenticated"))
		return
	}
	kratosID, ok := kratosIDValue.(string)
	if !ok {
		response.ReturnStatusInternalServerError(c, errors.New("invalid kratos_id type"))
		return
	}

	var req CreateRouteRequest
	// ShouldBindJSONはginの提供する関数で、内部的にはctx.Request.Bodyのstreamから読み出してJSON化してくれるもの
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnBadRequest(c, err)
		return
	}

	// GeoJSONをorbの型に変換
	pathGeom, err := geojson.ParseToLineString(req.PathGeom)
	if err != nil {
		response.ReturnBadRequest(c, errors.New("invalid path_geom GeoJSON: "+err.Error()))
		return
	}

	firstPoint, err := geojson.ParseToPoint(req.FirstPoint)
	if err != nil {
		response.ReturnBadRequest(c, errors.New("invalid first_point GeoJSON: "+err.Error()))
		return
	}

	lastPoint, err := geojson.ParseToPoint(req.LastPoint)
	if err != nil {
		response.ReturnBadRequest(c, errors.New("invalid last_point GeoJSON: "+err.Error()))
		return
	}

	// CoursePointsの変換
	coursePoints := make([]routeUsecase.CoursePointInput, len(req.CoursePoints))
	for i, cp := range req.CoursePoints {
		location, err := geojson.ParseToPoint(cp.Location)
		if err != nil {
			response.ReturnBadRequest(c, errors.New("invalid course_point location GeoJSON: "+err.Error()))
			return
		}
		coursePoints[i] = routeUsecase.CoursePointInput{
			SegDistM:      cp.SegDistM,
			CumDistM:      cp.CumDistM,
			Duration:      cp.Duration,
			Instruction:   cp.Instruction,
			RoadName:      cp.RoadName,
			ManeuverType:  cp.ManeuverType,
			Modifier:      cp.Modifier,
			Location:      location,
			BearingBefore: cp.BearingBefore,
			BearingAfter:  cp.BearingAfter,
		}
	}

	// Waypointsの変換
	waypoints := make([]routeUsecase.WaypointInput, len(req.Waypoints))
	for i, wp := range req.Waypoints {
		location, err := geojson.ParseToPoint(wp.Location)
		if err != nil {
			response.ReturnBadRequest(c, errors.New("invalid waypoint location GeoJSON: "+err.Error()))
			return
		}
		waypoints[i] = routeUsecase.WaypointInput{
			Location: location,
		}
	}

	input := routeUsecase.CreateRouteUseCaseInputDto{
		KratosID:           kratosID,
		Name:               req.Name,
		Description:        req.Description,
		HighlightedPhotoID: req.HighlightedPhotoID,
		Distance:           req.Distance,
		Duration:           req.Duration,
		ElevationGain:      req.ElevationGain,
		ElevationLoss:      req.ElevationLoss,
		PathGeom:           pathGeom,
		FirstPoint:         firstPoint,
		LastPoint:          lastPoint,
		Visibility:         req.Visibility,
		CoursePoints:       coursePoints,
		Waypoints:          waypoints,
	}

	dto, err := h.createRouteUsecase.CreateRoute(c.Request.Context(), input)
	if err != nil {
		response.ReturnStatusInternalServerError(c, err)
		return
	}

	res := RouteResponse{
		Route: RouteResponseModel{
			ID:                 dto.ID,
			UserID:             dto.UserID,
			Name:               dto.Name,
			Description:        dto.Description,
			HighlightedPhotoID: dto.HighlightedPhotoID,
			Distance:           dto.Distance,
			Duration:           dto.Duration,
			ElevationGain:      dto.ElevationGain,
			ElevationLoss:      dto.ElevationLoss,
			PathGeom:           geometry.GeometryToGeoJSON(dto.PathGeom),
			Bbox:               geometry.GeometryToGeoJSON(dto.Bbox),
			FirstPoint:         geometry.GeometryToGeoJSON(dto.FirstPoint),
			LastPoint:          geometry.GeometryToGeoJSON(dto.LastPoint),
			Visibility:         dto.Visibility,
		},
	}

	response.ReturnStatusCreated(c, res)
}

// GetRouteByID godoc
//	@Summary	ルートを取得する
//	@Tags		routes
//	@Accept		json
//	@Produce	json
//	@Param		route_id	path		string	true	"Route ID"
//	@Success	200			{object}	RouteResponse
//	@Failure	400			{object}	response.ErrorResponse
//	@Failure	404			{object}	response.ErrorResponse
//	@Failure	500			{object}	response.ErrorResponse
//	@Router		/routes/{route_id} [get]
func (h *Handler) GetRouteByID(c *gin.Context) {
	id := c.Param("route_id")

	dto, err := h.getRouteUsecase.GetRouteByID(c.Request.Context(), id)
	if err != nil {
		response.ReturnStatusInternalServerError(c, err)
		return
	}

	// CoursePointsの変換
	coursePoints := make([]CoursePointResponse, len(dto.CoursePoints))
	for i, cp := range dto.CoursePoints {
		coursePoints[i] = CoursePointResponse{
			ID:            cp.ID,
			StepOrder:     cp.StepOrder,
			SegDistM:      cp.SegDistM,
			CumDistM:      cp.CumDistM,
			Duration:      cp.Duration,
			Instruction:   cp.Instruction,
			RoadName:      cp.RoadName,
			ManeuverType:  cp.ManeuverType,
			Modifier:      cp.Modifier,
			Location:      geometry.GeometryToGeoJSON(cp.Location),
			BearingBefore: cp.BearingBefore,
			BearingAfter:  cp.BearingAfter,
		}
	}

	// Waypointsの変換
	waypoints := make([]WaypointResponse, len(dto.Waypoints))
	for i, wp := range dto.Waypoints {
		waypoints[i] = WaypointResponse{
			ID:       wp.ID,
			Location: geometry.GeometryToGeoJSON(wp.Location),
		}
	}

	res := RouteResponse{
		Route: RouteResponseModel{
			ID:                 dto.ID,
			Name:               dto.Name,
			UserID:             dto.UserID,
			UserName:           dto.UserName,
			Description:        dto.Description,
			HighlightedPhotoID: dto.HighlightedPhotoID,
			Distance:           dto.Distance,
			Duration:           dto.Duration,
			ElevationGain:      dto.ElevationGain,
			ElevationLoss:      dto.ElevationLoss,
			PathGeom:           geometry.GeometryToGeoJSON(dto.PathGeom),
			Bbox:               geometry.GeometryToGeoJSON(dto.Bbox),
			FirstPoint:         geometry.GeometryToGeoJSON(dto.FirstPoint),
			LastPoint:          geometry.GeometryToGeoJSON(dto.LastPoint),
			Visibility:         dto.Visibility,
			CreatedAt:          dto.CreatedAt,
			UpdatedAt:          dto.UpdatedAt,
			CoursePoints:       coursePoints,
			Waypoints:          waypoints,
		},
	}

	response.ReturnStatusOK(c, res)
}

// UpdateRoute godoc
//	@Summary	ルートを更新する
//	@Tags		routes
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		route_id	path	string				true	"Route ID"
//	@Param		request		body	UpdateRouteRequest	true	"Update Route Request"
//	@Success	204
//	@Failure	400	{object}	response.ErrorResponse
//	@Failure	401	{object}	response.ErrorResponse
//	@Failure	404	{object}	response.ErrorResponse
//	@Failure	500	{object}	response.ErrorResponse
//	@Router		/routes/{route_id} [put]
func (h *Handler) UpdateRoute(c *gin.Context) {
	routeID := c.Param("route_id")

	// 認証ミドルウェアからKratosIDを取得
	kratosIDValue, exists := c.Get("kratos_id")
	if !exists {
		response.ReturnStatusUnauthorized(c, errors.New("user not authenticated"))
		return
	}
	kratosID, ok := kratosIDValue.(string)
	if !ok {
		response.ReturnStatusInternalServerError(c, errors.New("invalid kratos_id type"))
		return
	}

	var req UpdateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnBadRequest(c, err)
		return
	}

	// GeoJSONをorbの型に変換
	pathGeom, err := geojson.ParseToLineString(req.PathGeom)
	if err != nil {
		response.ReturnBadRequest(c, errors.New("invalid path_geom GeoJSON: "+err.Error()))
		return
	}

	firstPoint, err := geojson.ParseToPoint(req.FirstPoint)
	if err != nil {
		response.ReturnBadRequest(c, errors.New("invalid first_point GeoJSON: "+err.Error()))
		return
	}

	lastPoint, err := geojson.ParseToPoint(req.LastPoint)
	if err != nil {
		response.ReturnBadRequest(c, errors.New("invalid last_point GeoJSON: "+err.Error()))
		return
	}

	// CoursePointsの変換
	coursePoints := make([]routeUsecase.UpdatedCoursePointInput, len(req.CoursePoints))
	for i, cp := range req.CoursePoints {
		location, err := geojson.ParseToPoint(cp.Location)
		if err != nil {
			response.ReturnBadRequest(c, errors.New("invalid course_point location GeoJSON: "+err.Error()))
			return
		}
		coursePoints[i] = routeUsecase.UpdatedCoursePointInput{
			SegDistM:      cp.SegDistM,
			CumDistM:      cp.CumDistM,
			Duration:      cp.Duration,
			Instruction:   cp.Instruction,
			RoadName:      cp.RoadName,
			ManeuverType:  cp.ManeuverType,
			Modifier:      cp.Modifier,
			Location:      location,
			BearingBefore: cp.BearingBefore,
			BearingAfter:  cp.BearingAfter,
		}
	}

	// Waypointsの変換
	waypoints := make([]routeUsecase.UpdatedWaypointInput, len(req.Waypoints))
	for i, wp := range req.Waypoints {
		location, err := geojson.ParseToPoint(wp.Location)
		if err != nil {
			response.ReturnBadRequest(c, errors.New("invalid waypoint location GeoJSON: "+err.Error()))
			return
		}
		waypoints[i] = routeUsecase.UpdatedWaypointInput{
			Location: location,
		}
	}

	input := routeUsecase.UpdateRouteUseCaseInputDto{
		ID:                 routeID,
		KratosID:           kratosID,
		Name:               req.Name,
		Description:        req.Description,
		HighlightedPhotoID: req.HighlightedPhotoID,
		Distance:           req.Distance,
		Duration:           req.Duration,
		ElevationGain:      req.ElevationGain,
		ElevationLoss:      req.ElevationLoss,
		PathGeom:           pathGeom,
		FirstPoint:         firstPoint,
		LastPoint:          lastPoint,
		Visibility:         req.Visibility,
		CoursePoints:       coursePoints,
		Waypoints:          waypoints,
	}

	if err := h.updateRouteUsecase.UpdateRoute(c.Request.Context(), input); err != nil {
		response.ReturnStatusInternalServerError(c, err)
		return
	}

	response.ReturnStatusNoContent(c)
}

// DeleteRoute godoc
//	@Summary	ルートを削除する
//	@Tags		routes
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		route_id	path	string	true	"Route ID"
//	@Success	204
//	@Failure	400	{object}	response.ErrorResponse
//	@Failure	401	{object}	response.ErrorResponse
//	@Failure	404	{object}	response.ErrorResponse
//	@Failure	500	{object}	response.ErrorResponse
//	@Router		/routes/{route_id} [delete]
func (h *Handler) DeleteRoute(c *gin.Context) {
	routeID := c.Param("route_id")

	// 認証ミドルウェアからKratosIDを取得
	kratosIDValue, exists := c.Get("kratos_id")
	if !exists {
		response.ReturnStatusUnauthorized(c, errors.New("user not authenticated"))
		return
	}
	kratosID, ok := kratosIDValue.(string)
	if !ok {
		response.ReturnStatusInternalServerError(c, errors.New("invalid kratos_id type"))
		return
	}

	if err := h.deleteRouteUsecase.DeleteRoute(c.Request.Context(), routeID, kratosID); err != nil {
		response.ReturnStatusInternalServerError(c, err)
		return
	}

	response.ReturnStatusNoContent(c)
}

// GetRoutesByUserID godoc
//	@Summary	ユーザーのルート一覧を取得する
//	@Tags		routes
//	@Accept		json
//	@Produce	json
//	@Param		user_id			path		string	true	"User ID"
//	@Param		keyword			query		string	false	"Keyword to search in route names"
//	@Param		min_distance	query		string	false	"Minimum distance filter"
//	@Param		max_distance	query		string	false	"Maximum distance filter"
//	@Param		min_elevation	query		string	false	"Minimum elevation gain filter"
//	@Param		max_elevation	query		string	false	"Maximum elevation gain filter"
//	@Param		visibility		query		string	false	"Visibility filter"
//	@Param		author			query		string	false	"Author filter"
//	@Success	200				{object}	RouteListResponse
//	@Failure	400				{object}	response.ErrorResponse
//	@Failure	500				{object}	response.ErrorResponse
//	@Router		/users/{user_id}/routes [get]
func (h *Handler) GetRoutesByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	// keyword := c.Query("keyword")
	// min_distance := c.Query("min_distance")
	// max_distance := c.Query("max_distance")
	// min_elevation := c.Query("min_elevation")
	// max_elevation := c.Query("max_elevation")
	// visibility := c.Query("visibility")
	// author := c.Query("author")


	dtos, err := h.getRouteUsecase.GetRoutesByUserID(c.Request.Context(), userID)
	if err != nil {
		response.ReturnStatusInternalServerError(c, err)
		return
	}

	routes := make([]RouteResponseModel, len(dtos))
	for i, dto := range dtos {
		routes[i] = RouteResponseModel{
			ID:                 dto.ID,
			UserID:             dto.UserID,
			UserName:           dto.UserName,
			Name:               dto.Name,
			Description:        dto.Description,
			HighlightedPhotoID: dto.HighlightedPhotoID,
			Distance:           dto.Distance,
			Duration:           dto.Duration,
			ElevationGain:      dto.ElevationGain,
			ElevationLoss:      dto.ElevationLoss,
			Visibility:         dto.Visibility,
			CreatedAt:          dto.CreatedAt,
			UpdatedAt:          dto.UpdatedAt,
		}
	}

	res := RouteListResponse{
		Routes: routes,
	}

	response.ReturnStatusOK(c, res)
}

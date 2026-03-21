package route

import (
	"context"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	gpxpkg "github.com/YukiAminaka/cycle-route-backend/internal/pkg/gpx"
	"github.com/tkrajina/gpxgo/gpx"
)

type IExportGPXUsecase interface {
	ExportGPX(vtx context.Context, routeID string) ([]byte, error)
}

type exportGPXUsecase struct {
	routeRepo routeDomain.IRouteRepository
}

func NewExportGPXUsecase(routeRepo routeDomain.IRouteRepository) IExportGPXUsecase {
	return &exportGPXUsecase{
		routeRepo: routeRepo,
	}
}

func (u *exportGPXUsecase) ExportGPX(ctx context.Context, routeID string) ([]byte, error) {
	route, err := u.routeRepo.GetRouteByID(ctx, routeID)
	if err != nil {
		return nil, err
	}

	gpxData, err := gpxpkg.RouteToGPX(route)
	if err != nil {
		return nil, err
	}

	return gpxData.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
}

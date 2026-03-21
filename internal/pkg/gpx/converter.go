package gpx

import (
	"errors"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	"github.com/paulmach/orb"

	"github.com/tkrajina/gpxgo/gpx"
)

func RouteToGPX(r *route.Route) (*gpx.GPX, error) {
    g := gpx.GPX{
        Version: "1.1",
        Creator: "rideline",
    }

    rte := gpx.GPXRoute{
        Name:        r.Name(),
        Description: r.Description(),
    }

    // path_geomの各座標をrteptに変換
    lineString, ok := r.PathGeom().Geometry.(orb.LineString)
    if !ok {
        return nil, errors.New("invalid path geometry")
    }
    for _, pr := range lineString {
        rte.Points = append(rte.Points, gpx.GPXPoint{
            Point: gpx.Point{Latitude: pr[1], Longitude: pr[0]},
        })
    }

    g.Routes = append(g.Routes, rte)
    return &g, nil
}

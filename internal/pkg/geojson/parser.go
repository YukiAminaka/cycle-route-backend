package geojson

import (
	"encoding/json"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// ParseToLineString はGeoJSON文字列をorb.LineStringに変換する
func ParseToLineString(geoJSON string) (orb.LineString, error) {
	var g *geojson.Geometry
	if err := json.Unmarshal([]byte(geoJSON), &g); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GeoJSON: %w", err)
	}

	geom := g.Geometry()
	lineString, ok := geom.(orb.LineString)
	if !ok {
		return nil, fmt.Errorf("geometry is not a LineString")
	}

	return lineString, nil
}

// ParseToPoint はGeoJSON文字列をorb.Pointに変換する
func ParseToPoint(geoJSON string) (orb.Point, error) {
	var g *geojson.Geometry
	if err := json.Unmarshal([]byte(geoJSON), &g); err != nil {
		return orb.Point{}, fmt.Errorf("failed to unmarshal GeoJSON: %w", err)
	}

	geom := g.Geometry()
	point, ok := geom.(orb.Point)
	if !ok {
		return orb.Point{}, fmt.Errorf("geometry is not a Point")
	}

	return point, nil
}

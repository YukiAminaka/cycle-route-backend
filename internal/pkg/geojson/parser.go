package geojson

import (
	"encoding/json"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// ParseToLineString はGeoJSON文字列をorb.LineStringに変換する
func ParseToLineString(geoJSON string) (orb.LineString, error) {
	var feature *geojson.Feature
	if err := json.Unmarshal([]byte(geoJSON), &feature); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GeoJSON: %w", err)
	}

	lineString, ok := feature.Geometry.(orb.LineString)
	if !ok {
		return nil, fmt.Errorf("geometry is not a LineString")
	}

	return lineString, nil
}

// ParseToPoint はGeoJSON文字列をorb.Pointに変換する
func ParseToPoint(geoJSON string) (orb.Point, error) {
	var feature *geojson.Feature
	if err := json.Unmarshal([]byte(geoJSON), &feature); err != nil {
		return orb.Point{}, fmt.Errorf("failed to unmarshal GeoJSON: %w", err)
	}

	point, ok := feature.Geometry.(orb.Point)
	if !ok {
		return orb.Point{}, fmt.Errorf("geometry is not a Point")
	}

	return point, nil
}

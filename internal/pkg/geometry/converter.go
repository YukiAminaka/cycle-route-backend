package geometry

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// geometryToGeoJSON はorb.Geometryをstringに変換する
func GeometryToGeoJSON(geom orb.Geometry) *string {
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
func PointToGeoJSON(point *orb.Point) *string {
	if point == nil {
		return nil
	}
	return GeometryToGeoJSON(*point)
}
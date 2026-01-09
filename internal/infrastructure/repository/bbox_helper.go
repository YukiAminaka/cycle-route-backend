package repository

import (
	"fmt"
	"strings"

	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
)

// CalculateBbox は LineString から Bounding Box (Polygon) を計算します
// path_geom から bbox を生成する際に使用します
func CalculateBbox(pathGeom orb.Geometry) dbgen.OrbGeometry {
	// orb.Boundを使って境界ボックスを取得
	bound := pathGeom.Bound()

	// PostGIS用にPolygonを手動で構築
	// Polygonは外部リング(Ring)を持ち、Ringは最初と最後の点が同じである必要がある(closed)
	minX, minY := bound.Min.X(), bound.Min.Y()
	maxX, maxY := bound.Max.X(), bound.Max.Y()

	polygon := orb.Polygon{
		orb.Ring{
			{minX, minY}, // 左下
			{maxX, minY}, // 右下
			{maxX, maxY}, // 右上
			{minX, maxY}, // 左上
			{minX, minY}, // 閉じる(最初の点に戻る)
		},
	}

	return dbgen.OrbGeometry{Geometry: polygon}
}

// CalculateBboxWithPadding はパディング付きの Bounding Box を計算します
// padding は緯度経度での追加マージン(度単位)
// 例: 0.001 = 約111m (緯度1度 ≈ 111km)
func CalculateBboxWithPadding(pathGeom orb.Geometry, padding float64) dbgen.OrbGeometry {
	bound := pathGeom.Bound()

	// パディングを追加
	paddedBound := bound.Pad(padding)

	// PostGIS用にPolygonを手動で構築
	minX, minY := paddedBound.Min.X(), paddedBound.Min.Y()
	maxX, maxY := paddedBound.Max.X(), paddedBound.Max.Y()

	polygon := orb.Polygon{
		orb.Ring{
			{minX, minY},
			{maxX, minY},
			{maxX, maxY},
			{minX, maxY},
			{minX, minY}, // 閉じる
		},
	}

	return dbgen.OrbGeometry{Geometry: polygon}
}

// ParseEWKT はEWKT文字列をorb.Geometryに変換します
// "SRID=4326;POLYGON(...)" -> orb.Polygon
func ParseEWKT(ewktString string) (orb.Geometry, error) {
	// "SRID=4326;" プレフィックスを削除
	wktString := strings.TrimPrefix(ewktString, "SRID=4326;")

	// WKT文字列をパース (wkt.Unmarshalは(orb.Geometry, error)を返す)
	geom, err := wkt.Unmarshal(wktString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WKT: %w", err)
	}

	return geom, nil
}

package dbgen

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
)

// OrbGeometry は任意の PostGIS geometry (Point, LineString, Polygon etc.) を格納できる汎用的な型です。
type OrbGeometry struct {
	orb.Geometry
}

// Scan は DB からの値を orb.Geometry に変換します (sql.Scanner)
func (g *OrbGeometry) Scan(value interface{}) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data, err = hex.DecodeString(v)
		if err != nil {
			return fmt.Errorf("OrbGeometry scan error: failed to decode hex string: %w", err)
		}
	case nil:
		g.Geometry = nil
		return nil
	default:
		return fmt.Errorf("OrbGeometry scan error: expected []byte or string, got %T", value)
	}

	// ewkb.Unmarshal はバイト列を解析し、適切な具体的な型 (orb.Point, orb.LineStringなど)
	// を判定して g.Geometry に格納してくれます。
	geom, _, err := ewkb.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("OrbGeometry scan error: %w", err)
	}

	g.Geometry = geom
	return nil
}

// Value は Go の値を DB 用のバイナリに変換します (driver.Valuer)
func (g OrbGeometry) Value() (driver.Value, error) {
	if g.Geometry == nil {
		return nil, nil
	}

	// // すべてのGeometry型をEWKB形式で送信
	// data, err := ewkb.Marshal(g.Geometry, 4326)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal geometry to EWKB: %w", err)
	// }

	// // Polygon型の場合はデバッグログを出力
	// if _, ok := g.Geometry.(orb.Polygon); ok {
	// 	log.Printf("[DEBUG] Polygon EWKB hex: %x", data)
	// 	wktString := wkt.MarshalString(g.Geometry)
	// 	log.Printf("[DEBUG] Polygon WKT: %s", wktString)
	// }

	return ewkb.Value(g.Geometry, 4326), nil
}

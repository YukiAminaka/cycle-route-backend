package sqlc

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
	// PostGIS に書き込む際は、SRID情報を含めることができる EWKB を推奨。
	// デフォルトのSRID (0) でよい場合も問題なく動作します。
	return ewkb.Value(g.Geometry, 4326).Value() // 第2引数はSRID (例: 4326なら指定可能)
}


type OrbPoint struct {
	orb.Point
} 

// sql.Scannerの実装 (DB -> Go)
func (p *OrbPoint) Scan(value interface{}) error {
    // PostGISなどは通常、EWKBまたはWKB (byte slice) で返します
    _, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("OrbPoint scan error: expected []byte, got %T", value)
    }

    // EWKBまたはWKBとしてスキャン
    // ※ PostGISの場合はEWKBが多いため、ewkbスキャナを使うのが安全です
    s := ewkb.Scanner(&p.Point)
    return s.Scan(value)
}

// driver.Valuerの実装 (Go -> DB)
func (p OrbPoint) Value() (driver.Value, error) {
    return ewkb.Value(p.Point,4326).Value()
}
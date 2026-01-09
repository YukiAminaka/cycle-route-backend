package repository

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestCalculateBbox(t *testing.T) {
	// テスト用のLineStringを作成
	lineString := orb.LineString{
		{139.7000, 35.6800},
		{139.7100, 35.6850},
		{139.7200, 35.6900},
	}

	// bboxを計算
	bbox := CalculateBbox(lineString)
	t.Logf("Generated bbox: %+v", bbox)

	// Polygon型であることを確認
	polygon, ok := bbox.Geometry.(orb.Polygon)
	if !ok {
		t.Fatalf("expected orb.Polygon, got %T", bbox.Geometry)
	}

	// Ringが1つであることを確認
	if len(polygon) != 1 {
		t.Fatalf("expected 1 ring, got %d", len(polygon))
	}

	// Ringが5点(closed)であることを確認
	ring := polygon[0]
	if len(ring) != 5 {
		t.Fatalf("expected 5 points in ring (closed), got %d", len(ring))
	}

	// 最初と最後の点が同じであることを確認
	if !ring[0].Equal(ring[4]) {
		t.Fatalf("ring is not closed: first=%v, last=%v", ring[0], ring[4])
	}

	// 境界値の確認
	expectedMinX, expectedMinY := 139.7000, 35.6800
	expectedMaxX, expectedMaxY := 139.7200, 35.6900

	if ring[0].X() != expectedMinX || ring[0].Y() != expectedMinY {
		t.Errorf("expected min point (%f, %f), got (%f, %f)",
			expectedMinX, expectedMinY, ring[0].X(), ring[0].Y())
	}

	if ring[2].X() != expectedMaxX || ring[2].Y() != expectedMaxY {
		t.Errorf("expected max point (%f, %f), got (%f, %f)",
			expectedMaxX, expectedMaxY, ring[2].X(), ring[2].Y())
	}

	t.Logf("bbox test passed")
}

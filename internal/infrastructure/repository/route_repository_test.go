package repository

import (
	"context"
	"testing"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	"github.com/paulmach/orb"
)

// ヘルパー関数
func ptrFloat64(f float64) *float64 {
	return &f
}

func ptrString(s string) *string {
	return &s
}

func ptrInt32(i int32) *int32 {
	return &i
}

func ptrRouteGeom(g orb.Geometry) *routeDomain.Geometry {
	return &routeDomain.Geometry{Geometry: g}
}

func TestRouteRepository_GetRouteByID(t *testing.T) {
	q := GetTestQueries()
	routeRepository := NewRouteRepository(q)
	ctx := context.Background()
	resetTestData(t)

	tests := []struct {
		name           string
		routeID        string
		wantID         string
		wantUserID     string
		wantName       string
		wantDistance   float64
		wantDuration   float64
		wantCPCount    int // コースポイント数
		wantWPCount    int // ウェイポイント数
		wantErr        bool
	}{
		{
			name:         "IDによって皇居一周ルートが取得できること",
			routeID:      "019b5a50-0000-7000-8000-000000000001",
			wantID:       "019b5a50-0000-7000-8000-000000000001",
			wantUserID:   "70d6037a-b67b-4aa8-b5a3-da393b514f24",
			wantName:     "皇居一周ルート",
			wantDistance: 5000.0,
			wantDuration: 900.0,
			wantCPCount:  4,
			wantWPCount:  3,
			wantErr:      false,
		},
		{
			name:         "IDによって多摩川サイクリングロードが取得できること",
			routeID:      "019b5a50-0000-7000-8000-000000000002",
			wantID:       "019b5a50-0000-7000-8000-000000000002",
			wantUserID:   "70d6037a-b67b-4aa8-b5a3-da393b514f24",
			wantName:     "多摩川サイクリングロード",
			wantDistance: 15000.0,
			wantDuration: 2700.0,
			wantCPCount:  4,
			wantWPCount:  3,
			wantErr:      false,
		},
		{
			name:    "存在しないIDの場合はエラー",
			routeID: "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeRepository.GetRouteByID(ctx, tt.routeID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID mismatch: want %s, got %s", tt.wantID, got.ID())
			}
			if got.UserID() != tt.wantUserID {
				t.Errorf("UserID mismatch: want %s, got %s", tt.wantUserID, got.UserID())
			}
			if got.Name() != tt.wantName {
				t.Errorf("Name mismatch: want %s, got %s", tt.wantName, got.Name())
			}
			if got.Distance() != tt.wantDistance {
				t.Errorf("Distance mismatch: want %f, got %f", tt.wantDistance, got.Distance())
			}
			if got.Duration() != tt.wantDuration {
				t.Errorf("Duration mismatch: want %f, got %f", tt.wantDuration, got.Duration())
			}
			if len(got.CoursePoints()) != tt.wantCPCount {
				t.Errorf("CoursePoints count mismatch: want %d, got %d", tt.wantCPCount, len(got.CoursePoints()))
			}
			if len(got.Waypoints()) != tt.wantWPCount {
				t.Errorf("Waypoints count mismatch: want %d, got %d", tt.wantWPCount, len(got.Waypoints()))
			}
		})
	}
}

func TestRouteRepository_GetRoutesByUserID(t *testing.T) {
	q := GetTestQueries()
	routeRepository := NewRouteRepository(q)
	ctx := context.Background()
	resetTestData(t)

	tests := []struct {
		name      string
		userID    string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "ユーザーIDで複数のルートが取得できること",
			userID:    "70d6037a-b67b-4aa8-b5a3-da393b514f24",
			wantCount: 3, // 皇居、多摩川、削除済み（削除済みは除外される想定）
			wantErr:   false,
		},
		{
			name:      "別のユーザーのルートが取得できること",
			userID:    "019b5a46-1e77-7b9d-ac62-b438a0fc89cb",
			wantCount: 1, // 湘南海岸
			wantErr:   false,
		},
		{
			name:      "ルートがないユーザーの場合は空配列",
			userID:    "019b5a46-71c6-702e-a8a0-c0cf5a2fe757", // newbie_cyclist
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeRepository.GetRoutesByUserID(ctx, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != tt.wantCount {
				t.Errorf("count mismatch: want %d, got %d", tt.wantCount, len(got))
			}
		})
	}
}

func TestRouteRepository_CountRoutesByUserID(t *testing.T) {
	q := GetTestQueries()
	routeRepository := NewRouteRepository(q)
	ctx := context.Background()
	resetTestData(t)

	tests := []struct {
		name      string
		userID    string
		wantCount int64
		wantErr   bool
	}{
		{
			name:      "ユーザーのルート数をカウントできること",
			userID:    "70d6037a-b67b-4aa8-b5a3-da393b514f24",
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "別のユーザーのルート数をカウントできること",
			userID:    "019b5a46-1e77-7b9d-ac62-b438a0fc89cb",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "ルートがないユーザーの場合は0",
			userID:    "019b5a46-71c6-702e-a8a0-c0cf5a2fe757",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeRepository.CountRoutesByUserID(ctx, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.wantCount {
				t.Errorf("count mismatch: want %d, got %d", tt.wantCount, got)
			}
		})
	}
}

func TestRouteRepository_SaveRoute(t *testing.T) {
	q := GetTestQueries()
	routeRepository := NewRouteRepository(q)
	ctx := context.Background()
	resetTestData(t)

	// テスト用のルートを作成
	pathGeom := routeDomain.Geometry{
		Geometry: orb.LineString{
			{139.7000, 35.6800},
			{139.7100, 35.6850},
			{139.7200, 35.6900},
		},
	}
	firstPoint := routeDomain.Geometry{Geometry: orb.Point{139.7000, 35.6800}}
	lastPoint := routeDomain.Geometry{Geometry: orb.Point{139.7200, 35.6900}}

	newRoute, err := routeDomain.NewRoute(
		"70d6037a-b67b-4aa8-b5a3-da393b514f24",
		"新規テストルート",
		"テスト用の説明",
		nil,
		3000.0,
		600.0,
		10.0,
		5.0,
		pathGeom,
		firstPoint,
		lastPoint,
		2,
	)
	if err != nil {
		t.Fatalf("failed to create new route: %v", err)
	}

	// ウェイポイントを追加
	err = newRoute.AddWaypoint(routeDomain.Geometry{Geometry: orb.Point{139.7000, 35.6800}})
	if err != nil {
		t.Fatalf("failed to add waypoint: %v", err)
	}

	// コースポイントを追加
	err = newRoute.AddCoursePoint(
		ptrFloat64(0.0),
		ptrFloat64(0.0),
		ptrFloat64(0.0),
		ptrString("スタート"),
		ptrString("テスト道路"),
		ptrString("depart"),
		nil,
		ptrRouteGeom(orb.Point{139.7000, 35.6800}),
		nil,
		ptrInt32(90),
	)
	if err != nil {
		t.Fatalf("failed to add course point: %v", err)
	}

	tests := []struct {
		name    string
		route   *routeDomain.Route
		wantErr bool
	}{
		{
			name:    "新しいルートを保存できること",
			route:   newRoute,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := routeRepository.SaveRoute(ctx, tt.route)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// 保存後に取得して確認
			saved, err := routeRepository.GetRouteByID(ctx, tt.route.ID())
			if err != nil {
				t.Fatalf("failed to get saved route: %v", err)
			}

			if saved.ID() != tt.route.ID() {
				t.Errorf("ID mismatch: want %s, got %s", tt.route.ID(), saved.ID())
			}
			if saved.Name() != tt.route.Name() {
				t.Errorf("Name mismatch: want %s, got %s", tt.route.Name(), saved.Name())
			}
			if len(saved.Waypoints()) != len(tt.route.Waypoints()) {
				t.Errorf("Waypoints count mismatch: want %d, got %d", len(tt.route.Waypoints()), len(saved.Waypoints()))
			}
			if len(saved.CoursePoints()) != len(tt.route.CoursePoints()) {
				t.Errorf("CoursePoints count mismatch: want %d, got %d", len(tt.route.CoursePoints()), len(saved.CoursePoints()))
			}
		})
	}
}

func TestRouteRepository_DeleteRoute(t *testing.T) {
	q := GetTestQueries()
	routeRepository := NewRouteRepository(q)
	ctx := context.Background()
	resetTestData(t)

	tests := []struct {
		name    string
		routeID string
		wantErr bool
	}{
		{
			name:    "既存のルートを削除できること",
			routeID: "019b5a50-0000-7000-8000-000000000001",
			wantErr: false,
		},
		{
			name:    "存在しないルートの削除はエラー",
			routeID: "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 削除前に存在確認（エラーを期待しない場合のみ）
			if !tt.wantErr {
				_, err := routeRepository.GetRouteByID(ctx, tt.routeID)
				if err != nil {
					t.Fatalf("route should exist before delete: %v", err)
				}
			}

			err := routeRepository.DeleteRoute(ctx, tt.routeID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// 削除後に取得してエラーになることを確認
			_, err = routeRepository.GetRouteByID(ctx, tt.routeID)
			if err == nil {
				t.Error("route should not exist after delete")
			}
		})
	}
}

func TestRouteRepository_UpdateRoute(t *testing.T) {
	q := GetTestQueries()
	routeRepository := NewRouteRepository(q)
	ctx := context.Background()
	resetTestData(t)

	// 既存のルートを取得
	existingRoute, err := routeRepository.GetRouteByID(ctx, "019b5a50-0000-7000-8000-000000000001")
	if err != nil {
		t.Fatalf("failed to get existing route: %v", err)
	}

	// 更新用のルートを作成（ReconstructRouteを使用）
	pathGeom := routeDomain.Geometry{
		Geometry: orb.LineString{
			{139.7528, 35.6850},
			{139.7580, 35.6820},
			{139.7600, 35.6780},
		},
	}
	firstPoint := routeDomain.Geometry{Geometry: orb.Point{139.7528, 35.6850}}
	lastPoint := routeDomain.Geometry{Geometry: orb.Point{139.7600, 35.6780}}

	// bboxは既存のルートから取得（データベースで自動生成されたもの）
	updatedRoute, err := routeDomain.ReconstructRoute(
		existingRoute.ID(),
		existingRoute.UserID(),
		"更新された皇居一周ルート",
		"更新された説明文",
		nil,
		6000.0,  // 距離を変更
		1200.0,  // 時間を変更
		30.0,
		30.0,
		pathGeom,
		existingRoute.Bbox(), // 既存のbboxを使用
		firstPoint,
		lastPoint,
		2,
		existingRoute.CreatedAt(), // 既存の作成日時を使用
		existingRoute.UpdatedAt(), // 既存の更新日時を使用
	)
	if err != nil {
		t.Fatalf("failed to reconstruct route: %v", err)
	}

	// 新しいウェイポイントを追加
	err = updatedRoute.AddWaypoint(routeDomain.Geometry{Geometry: orb.Point{139.7528, 35.6850}})
	if err != nil {
		t.Fatalf("failed to add waypoint: %v", err)
	}

	// 新しいコースポイントを追加
	err = updatedRoute.AddCoursePoint(
		ptrFloat64(0.0),
		ptrFloat64(0.0),
		ptrFloat64(0.0),
		ptrString("更新されたスタート"),
		ptrString("更新された道路"),
		ptrString("depart"),
		nil,
		ptrRouteGeom(orb.Point{139.7528, 35.6850}),
		nil,
		ptrInt32(90),
	)
	if err != nil {
		t.Fatalf("failed to add course point: %v", err)
	}

	tests := []struct {
		name    string
		route   *routeDomain.Route
		wantErr bool
	}{
		{
			name:    "既存のルートを更新できること",
			route:   updatedRoute,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := routeRepository.UpdateRoute(ctx, tt.route)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// 更新後に取得して確認
			updated, err := routeRepository.GetRouteByID(ctx, tt.route.ID())
			if err != nil {
				t.Fatalf("failed to get updated route: %v", err)
			}

			if updated.Name() != tt.route.Name() {
				t.Errorf("Name mismatch: want %s, got %s", tt.route.Name(), updated.Name())
			}
			if updated.Distance() != tt.route.Distance() {
				t.Errorf("Distance mismatch: want %f, got %f", tt.route.Distance(), updated.Distance())
			}
			if updated.Duration() != tt.route.Duration() {
				t.Errorf("Duration mismatch: want %f, got %f", tt.route.Duration(), updated.Duration())
			}
			if len(updated.Waypoints()) != len(tt.route.Waypoints()) {
				t.Errorf("Waypoints count mismatch: want %d, got %d", len(tt.route.Waypoints()), len(updated.Waypoints()))
			}
			if len(updated.CoursePoints()) != len(tt.route.CoursePoints()) {
				t.Errorf("CoursePoints count mismatch: want %d, got %d", len(tt.route.CoursePoints()), len(updated.CoursePoints()))
			}
		})
	}
}
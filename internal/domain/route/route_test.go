package route

import (
	"testing"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/paulmach/orb"
)

//テストケースの設計:
//集約ルートパターンのテスト: CoursePointとWaypointはRouteを通してのみ作成されることを確認
//バリデーションの網羅: 各エンティティの必須フィールドと型チェックを確認
//ビジネスロジックのテスト:
//stepOrderの自動採番
//メトリクス(distance/duration)の再計算
//不変条件の確認: RouteIDが正しく設定されているか

func TestNewRoute(t *testing.T) {
	userID := user.NewUserID()
	type args struct {
		userID             string
		name               string
		description        string
		highlightedPhotoID *int64
		distance           float64
		duration           float64
		elevationGain      float64
		elevationLoss      float64
		pathGeom           Geometry
		firstPoint         Geometry
		lastPoint          Geometry
		visibility         int16
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系: 全てのフィールドが正しい",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: false,
		},
		{
			name: "異常系: userIDが空",
			args: args{
				userID:             "",
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "userID is required",
		},
		{
			name: "異常系: nameが空",
			args: args{
				userID:             string(userID),
				name:               "",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "異常系: pathGeomがnil",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "pathGeom is required",
		},
		{
			name: "異常系: pathGeomがLineStringではない",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.Point{139.6917, 35.6895}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "pathGeom must be a LineString",
		},
		{
			name: "異常系: firstPointがnil",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "firstPoint is required",
		},
		{
			name: "異常系: firstPointがPointではない",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "firstPoint must be a Point",
		},
		{
			name: "異常系: lastPointがnil",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600.0,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "lastPoint is required",
		},
		{
			name: "異常系: lastPointがPointではない",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           600,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "lastPoint must be a Point",
		},
		{
			name: "異常系: distanceが負の値",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           -100.0,
				duration:           600,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "distance must be non-negative",
		},
		{
			name: "異常系: durationが負の値",
			args: args{
				userID:             string(userID),
				name:               "Test Route",
				description:        "This is a test route",
				highlightedPhotoID: nil,
				distance:           100.0,
				duration:           -600,
				elevationGain:      10.0,
				elevationLoss:      5.0,
				pathGeom:           Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				firstPoint:         Geometry{orb.Point{139.6917, 35.6895}},
				lastPoint:          Geometry{orb.Point{139.7000, 35.6900}},
				visibility:         1,
			},
			wantErr: true,
			errMsg:  "duration must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRoute(
				tt.args.userID,
				tt.args.name,
				tt.args.description,
				tt.args.highlightedPhotoID,
				tt.args.distance,
				tt.args.duration,
				tt.args.elevationGain,
				tt.args.elevationLoss,
				tt.args.pathGeom,
				tt.args.firstPoint,
				tt.args.lastPoint,
				tt.args.visibility,
			)

			// エラーチェック
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewRoute() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("NewRoute() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			// エラーが期待されていない場合
			if err != nil {
				t.Errorf("NewRoute() unexpected error = %v", err)
				return
			}

			// 正常系の場合の検証
			if got == nil {
				t.Error("NewRoute() returned nil")
				return
			}

			// フィールドの検証（idとcreatedAtは自動生成されるので除外）
			if got.userID != tt.args.userID {
				t.Errorf("userID = %v, want %v", got.userID, tt.args.userID)
			}
			if got.name != tt.args.name {
				t.Errorf("name = %v, want %v", got.name, tt.args.name)
			}
			if got.description != tt.args.description {
				t.Errorf("description = %v, want %v", got.description, tt.args.description)
			}
			if got.distance != tt.args.distance {
				t.Errorf("distance = %v, want %v", got.distance, tt.args.distance)
			}
			if got.duration != tt.args.duration {
				t.Errorf("duration = %v, want %v", got.duration, tt.args.duration)
			} 
			if got.polyline != "" {
				t.Errorf("polyline = %v, want empty string", got.polyline)
			}
			if got.elevationGain != tt.args.elevationGain {
				t.Errorf("elevationGain = %v, want %v", got.elevationGain, tt.args.elevationGain)
			}
			if got.elevationLoss != tt.args.elevationLoss {
				t.Errorf("elevationLoss = %v, want %v", got.elevationLoss, tt.args.elevationLoss)
			}
			if got.visibility != tt.args.visibility {
				t.Errorf("visibility = %v, want %v", got.visibility, tt.args.visibility)
			}

			// coursePointsとwaypointsが空のスライスで初期化されていることを確認
			if len(got.coursePoints) != 0 {
				t.Errorf("coursePoints length = %v, want 0", len(got.coursePoints))
			}
			if len(got.waypoints) != 0 {
				t.Errorf("waypoints length = %v, want 0", len(got.waypoints))
			}
		})
	}
}

func TestAddCoursePoint(t *testing.T) {
	userID := user.NewUserID()

	type args struct {
		segDistM      *float64
		cumDistM      *float64
		duration      *float64
		instruction   *string
		roadName      *string
		maneuverType  *string
		modifier      *string
		location      *Geometry
		bearingBefore *int32
		bearingAfter  *int32
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系: 全てのフィールドが正しい",
			args: args{
				segDistM:      new(100.0),
				cumDistM:      new(100.0),
				duration:      new(600.0),
				instruction:   new("Turn right"),
				roadName:      new("Main St"),
				maneuverType:  new("turn"),
				modifier:      new("right"),
				location:      &Geometry{orb.Point{139.6917, 35.6895}},
				bearingBefore: new(int32(90)),
				bearingAfter:  new(int32(180)),
			},
			wantErr: false,
		},
		{
			name: "正常系: 必須フィールドのみ",
			args: args{
				segDistM:      nil,
				cumDistM:      nil,
				duration:      nil,
				instruction:   nil,
				roadName:      nil,
				maneuverType:  nil,
				modifier:      nil,
				location:      &Geometry{orb.Point{139.6950, 35.6897}},
				bearingBefore: nil,
				bearingAfter:  nil,
			},
			wantErr: false,
		},
		{
			name: "異常系: locationがnil",
			args: args{
				segDistM:      new(100.0),
				cumDistM:      new(100.0),
				duration:      new(600.0),
				instruction:   new("Turn right"),
				roadName:      new("Main St"),
				maneuverType:  new("turn"),
				modifier:      new("right"),
				location:      nil,
				bearingBefore: new(int32(90)),
				bearingAfter:  new(int32(180)),
			},
			wantErr: true,
			errMsg:  "location is required",
		},
		{
			name: "異常系: locationがPointではない",
			args: args{
				segDistM:      new(100.0),
				cumDistM:      new(100.0),
				duration:      new(600.0),
				instruction:   new("Turn right"),
				roadName:      new("Main St"),
				maneuverType:  new("turn"),
				modifier:      new("right"),
				location:      &Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				bearingBefore: new(int32(90)),
				bearingAfter:  new(int32(180)),
			},
			wantErr: true,
			errMsg:  "location must be a Point geometry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用に新しいルートを作成
			testRoute, _ := NewRoute(
				string(userID),
				"Test Route",
				"This is a test route",
				nil,
				0.0,
				0,
				10.0,
				5.0,
				Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				Geometry{orb.Point{139.6917, 35.6895}},
				Geometry{orb.Point{139.7000, 35.6900}},
				1,
			)

			err := testRoute.AddCoursePoint(
				tt.args.segDistM,
				tt.args.cumDistM,
				tt.args.duration,
				tt.args.instruction,
				tt.args.roadName,
				tt.args.maneuverType,
				tt.args.modifier,
				tt.args.location,
				tt.args.bearingBefore,
				tt.args.bearingAfter,
			)

			// エラーチェック
			if tt.wantErr {
				if err == nil {
					t.Errorf("AddCoursePoint() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("AddCoursePoint() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			// エラーが期待されていない場合
			if err != nil {
				t.Errorf("AddCoursePoint() unexpected error = %v", err)
				return
			}

			// 正常系の場合の検証
			coursePoints := testRoute.CoursePoints()
			if len(coursePoints) != 1 {
				t.Errorf("coursePoints length = %v, want 1", len(coursePoints))
				return
			}

			// 追加されたCoursePointの検証
			cp := coursePoints[0]
			if cp.stepOrder != 0 {
				t.Errorf("stepOrder = %v, want 0", cp.stepOrder)
			}
			if cp.routeID != testRoute.id {
				t.Errorf("routeID = %v, want %v", cp.routeID, testRoute.id)
			}

			// メトリクスの再計算が正しく動作しているか確認
			if tt.args.segDistM != nil {
				if testRoute.distance != *tt.args.segDistM {
					t.Errorf("route distance = %v, want %v", testRoute.distance, *tt.args.segDistM)
				}
			}
			if tt.args.duration != nil {
				if testRoute.duration != *tt.args.duration {
					t.Errorf("route duration = %v, want %v", testRoute.duration, *tt.args.duration)
				}
			}
		})
	}
}

func TestAddCoursePoint_MultiplePoints(t *testing.T) {
	userID := user.NewUserID()

	// テスト用のルートを作成
	route, err := NewRoute(
		string(userID),
		"Test Route",
		"This is a test route",
		nil,
		0.0,
		0,
		10.0,
		5.0,
		Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
		Geometry{orb.Point{139.6917, 35.6895}},
		Geometry{orb.Point{139.7000, 35.6900}},
		1,
	)
	if err != nil {
		t.Fatalf("Failed to create route: %v", err)
	}

	// 複数のCoursePointを追加
	err = route.AddCoursePoint(
		new(100.0),
		new(100.0),
		new(300.0),
		new("Turn right"),
		new("Main St"),
		new("turn"),
		new("right"),
		&Geometry{orb.Point{139.6917, 35.6895}},
		new(int32(90)),
		new(int32(180)),
	)
	if err != nil {
		t.Fatalf("Failed to add first course point: %v", err)
	}

	err = route.AddCoursePoint(
		new(150.0),
		new(250.0),
		new(400.0),
		new("Turn left"),
		new("Second St"),
		new("turn"),
		new("left"),
		&Geometry{orb.Point{139.6950, 35.6897}},
		new(int32(180)),
		new(int32(270)),
	)
	if err != nil {
		t.Fatalf("Failed to add second course point: %v", err)
	}

	err = route.AddCoursePoint(
		new(200.0),
		new(450.0),
		new(500.0),
		new("Go straight"),
		new("Third Ave"),
		new("straight"),
		nil,
		&Geometry{orb.Point{139.7000, 35.6900}},
		new(int32(270)),
		new(int32(270)),
	)
	if err != nil {
		t.Fatalf("Failed to add third course point: %v", err)
	}

	// CoursePointsの数を確認
	coursePoints := route.CoursePoints()
	if len(coursePoints) != 3 {
		t.Errorf("coursePoints length = %v, want 3", len(coursePoints))
	}

	// stepOrderが正しく採番されているか確認
	for i, cp := range coursePoints {
		if cp.stepOrder != int32(i) {
			t.Errorf("coursePoint[%d].stepOrder = %v, want %v", i, cp.stepOrder, i)
		}
	}

	// メトリクスが再計算されているか確認（全ての距離と時間の合計）
	expectedDistance := 100.0 + 150.0 + 200.0
	expectedDuration := 300.0 + 400.0 + 500.0

	if route.distance != expectedDistance {
		t.Errorf("route.distance = %v, want %v", route.distance, expectedDistance)
	}
	if route.duration != expectedDuration {
		t.Errorf("route.duration = %v, want %v", route.duration, expectedDuration)
	}
}

func TestAddWaypoint(t *testing.T) {
	userID := user.NewUserID()

	type args struct {
		location Geometry
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系: Pointを追加",
			args: args{
				location: Geometry{orb.Point{139.6917, 35.6895}},
			},
			wantErr: false,
		},
		{
			name: "異常系: locationがnil",
			args: args{
				location: Geometry{},
			},
			wantErr: true,
			errMsg:  "location is required",
		},
		{
			name: "異常系: locationがPointではない",
			args: args{
				location: Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
			},
			wantErr: true,
			errMsg:  "location must be a Point geometry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用に新しいルートを作成
			testRoute, _ := NewRoute(
				string(userID),
				"Test Route",
				"This is a test route",
				nil,
				100.0,
				600,
				10.0,
				5.0,
				Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
				Geometry{orb.Point{139.6917, 35.6895}},
				Geometry{orb.Point{139.7000, 35.6900}},
				1,
			)

			err := testRoute.AddWaypoint(tt.args.location)

			// エラーチェック
			if tt.wantErr {
				if err == nil {
					t.Errorf("AddWaypoint() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("AddWaypoint() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			// エラーが期待されていない場合
			if err != nil {
				t.Errorf("AddWaypoint() unexpected error = %v", err)
				return
			}

			// 正常系の場合の検証
			waypoints := testRoute.Waypoints()
			if len(waypoints) != 1 {
				t.Errorf("waypoints length = %v, want 1", len(waypoints))
				return
			}

			// 追加されたWaypointの検証
			wp := waypoints[0]
			if wp.routeID != testRoute.id {
				t.Errorf("routeID = %v, want %v", wp.routeID, testRoute.id)
			}
		})
	}
}

func TestAddWaypoint_MultiplePoints(t *testing.T) {
	userID := user.NewUserID()

	// テスト用のルートを作成
	route, err := NewRoute(
		string(userID),
		"Test Route",
		"This is a test route",
		nil,
		100.0,
		600,
		10.0,
		5.0,
		Geometry{orb.LineString{{139.6917, 35.6895}, {139.7000, 35.6900}}},
		Geometry{orb.Point{139.6917, 35.6895}},
		Geometry{orb.Point{139.7000, 35.6900}},
		1,
	)
	if err != nil {
		t.Fatalf("Failed to create route: %v", err)
	}

	// 複数のWaypointを追加
	locations := []Geometry{
		{orb.Point{139.6917, 35.6895}},
		{orb.Point{139.6950, 35.6897}},
		{orb.Point{139.7000, 35.6900}},
	}

	for _, loc := range locations {
		err := route.AddWaypoint(loc)
		if err != nil {
			t.Fatalf("Failed to add waypoint: %v", err)
		}
	}

	// Waypointsの数を確認
	waypoints := route.Waypoints()
	if len(waypoints) != 3 {
		t.Errorf("waypoints length = %v, want 3", len(waypoints))
	}

	// 各Waypointが正しいRouteIDを持っているか確認
	for i, wp := range waypoints {
		if wp.routeID != route.id {
			t.Errorf("waypoint[%d].routeID = %v, want %v", i, wp.routeID, route.id)
		}
	}
}

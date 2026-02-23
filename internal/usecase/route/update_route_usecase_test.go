package route

import (
	"context"
	"errors"
	"strings"
	"testing"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	transactionApp "github.com/YukiAminaka/cycle-route-backend/internal/usecase/transaction"
	"github.com/paulmach/orb"
	"go.uber.org/mock/gomock"
)

// テスト用定数
const (
	testRouteID   = "019b5a50-0000-7000-8000-000000000001"
	testKratosID  = "2eb50f70-3a23-4067-99f6-9fd645686880"
	testUserID    = "019b5a8d-16a7-700a-be92-9ae11e7e5b9a"
	testUserName  = "Test User"
	testRouteName = "Test Route"
	testRouteDesc = "Test Description"

	// ルート属性
	testDistance      = 1000.0
	testDuration      = 3600.0
	testElevationGain = 100.0
	testElevationLoss = 50.0
	testVisibility    = 1
	testPolyline      = "aqtxEshssY]q@"
)

// ジェネリクスを使用したポインタヘルパー
func ptr[T any](v T) *T {
	return &v
}

// テスト用のデフォルトDTO作成ヘルパー
func createDefaultUpdateDTO() UpdateRouteUseCaseInputDto {
	return UpdateRouteUseCaseInputDto{
		ID:                 testRouteID,
		KratosID:           testKratosID,
		Name:               "Updated Route Name",
		Description:        "Updated Route Description",
		HighlightedPhotoID: nil,
		Distance:           200,
		Duration:           1200,
		ElevationGain:      150,
		ElevationLoss:      70,
		PathGeom: orb.LineString{
			{139.713592, 35.670692},
			{139.712618, 35.672179},
		},
		FirstPoint: orb.Point{139.713592, 35.670692},
		LastPoint:  orb.Point{139.712618, 35.672179},
		Visibility: 1,
		CoursePoints: []UpdatedCoursePointInput{
			{
				SegDistM:      ptr(100.0),
				CumDistM:      ptr(100.0),
				Duration:      ptr(600.0),
				Instruction:   ptr("Turn right"),
				RoadName:      ptr("Main St"),
				ManeuverType:  ptr("turn"),
				Modifier:      ptr("right"),
				Location:      orb.Point{139.713592, 35.670692},
				BearingBefore: ptr(int32(0)),
				BearingAfter:  ptr(int32(326)),
			},
			{
				SegDistM:      ptr(100.0),
				CumDistM:      ptr(200.0),
				Duration:      ptr(600.0),
				Instruction:   ptr("Turn left"),
				RoadName:      ptr("Second St"),
				ManeuverType:  ptr("turn"),
				Modifier:      ptr("left"),
				Location:      orb.Point{139.6917, 35.6895},
				BearingBefore: ptr(int32(90)),
				BearingAfter:  ptr(int32(180)),
			},
		},
		Waypoints: []UpdatedWaypointInput{
			{Location: orb.Point{139.713592, 35.670692}},
			{Location: orb.Point{139.712618, 35.672179}},
		},
	}
}

// テスト用のユーザー作成ヘルパー
func createTestUser() *userDomain.User {
	user, _ := userDomain.ReconstructUser(
		userDomain.UserID(testUserID),
		testKratosID,
		testUserName,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, false,
	)
	return user
}

// テスト用のルート作成ヘルパー
func createTestRoute(userID string) *routeDomain.Route {
	route, _ := routeDomain.ReconstructRoute(
		testRouteID,
		userID,
		testRouteName,
		testRouteDesc,
		nil,
		testDistance,
		testDuration,
		testElevationGain,
		testElevationLoss,
		routeDomain.Geometry{},
		routeDomain.Geometry{},
		routeDomain.Geometry{},
		routeDomain.Geometry{},
		testPolyline, testVisibility, "", "",
	)
	return route
}

// テスト用のモックセットアップ構造体
type updateRouteTestMocks struct {
	ctrl               *gomock.Controller
	mockRouteRepo      *routeDomain.MockIRouteRepository
	mockUserRepo       *userDomain.MockIUserRepository
	mockTxManager      *transactionApp.MockTransactionManager
	usecase            IUpdateRouteUsecase
}

// 各サブテスト用にモックを作成
func setupUpdateRouteMocks(t *testing.T) *updateRouteTestMocks {
	ctrl := gomock.NewController(t)
	mockRouteRepo := routeDomain.NewMockIRouteRepository(ctrl)
	mockUserRepo := userDomain.NewMockIUserRepository(ctrl)
	mockTxManager := transactionApp.NewMockTransactionManager(ctrl)
	uc := NewUpdateRouteUsecase(mockUserRepo, mockTxManager, mockRouteRepo)

	return &updateRouteTestMocks{
		ctrl:          ctrl,
		mockRouteRepo: mockRouteRepo,
		mockUserRepo:  mockUserRepo,
		mockTxManager: mockTxManager,
		usecase:       uc,
	}
}

func Test_updateRouteUsecase_UpdateRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		dto            UpdateRouteUseCaseInputDto
		setupMocks     func(m *updateRouteTestMocks)
		wantErr        bool
		wantErrContain string // エラーメッセージに含まれるべき文字列
	}{
		{
			name: "正常系: ルート更新に成功する",
			dto:  createDefaultUpdateDTO(),
			setupMocks: func(m *updateRouteTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), testKratosID).
					Return(createTestUser(), nil)

				m.mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), testRouteID).
					Return(createTestRoute(testUserID), nil)

				m.mockTxManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "正常系: 空のCoursePointsとWaypointsで更新に成功する",
			dto: func() UpdateRouteUseCaseInputDto {
				dto := createDefaultUpdateDTO()
				dto.CoursePoints = []UpdatedCoursePointInput{}
				dto.Waypoints = []UpdatedWaypointInput{}
				return dto
			}(),
			setupMocks: func(m *updateRouteTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), testKratosID).
					Return(createTestUser(), nil)

				m.mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), testRouteID).
					Return(createTestRoute(testUserID), nil)

				m.mockTxManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系: ユーザーが見つからない",
			dto: UpdateRouteUseCaseInputDto{
				ID:       testRouteID,
				KratosID: "invalid-kratos-id",
			},
			setupMocks: func(m *updateRouteTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), "invalid-kratos-id").
					Return(nil, errors.New("user not found"))
			},
			wantErr:        true,
			wantErrContain: "user not found",
		},
		{
			name: "異常系: ルートが見つからない",
			dto: UpdateRouteUseCaseInputDto{
				ID:       "invalid-route-id",
				KratosID: testKratosID,
			},
			setupMocks: func(m *updateRouteTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), testKratosID).
					Return(createTestUser(), nil)

				m.mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), "invalid-route-id").
					Return(nil, errors.New("route not found"))
			},
			wantErr:        true,
			wantErrContain: "route not found",
		},
		{
			name: "異常系: ルートの所有者ではない（権限エラー）",
			dto: UpdateRouteUseCaseInputDto{
				ID:       testRouteID,
				KratosID: testKratosID,
			},
			setupMocks: func(m *updateRouteTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), testKratosID).
					Return(createTestUser(), nil)

				// 異なるユーザーIDを持つルート
				m.mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), testRouteID).
					Return(createTestRoute("different-user-id"), nil)
			},
			wantErr:        true,
			wantErrContain: "unauthorized",
		},
		{
			name: "異常系: トランザクション内での更新に失敗",
			dto:  createDefaultUpdateDTO(),
			setupMocks: func(m *updateRouteTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), testKratosID).
					Return(createTestUser(), nil)

				m.mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), testRouteID).
					Return(createTestRoute(testUserID), nil)

				m.mockTxManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr:        true,
			wantErrContain: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 各サブテストで独立したモックを作成
			mocks := setupUpdateRouteMocks(t)
			tt.setupMocks(mocks)

			gotErr := mocks.usecase.UpdateRoute(context.Background(), tt.dto)

			// エラー検証
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateRoute() unexpected error: %v", gotErr)
					return
				}
				// エラーメッセージの内容を検証
				if tt.wantErrContain != "" && !strings.Contains(gotErr.Error(), tt.wantErrContain) {
					t.Errorf("UpdateRoute() error = %q, want error containing %q", gotErr.Error(), tt.wantErrContain)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateRoute() succeeded unexpectedly, want error")
			}
		})
	}
}

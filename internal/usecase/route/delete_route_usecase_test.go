package route

import (
	"context"
	"errors"
	"testing"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	transactionApp "github.com/YukiAminaka/cycle-route-backend/internal/usecase/transaction"
	"go.uber.org/mock/gomock"
)

func Test_deleteRouteUsecase_DeleteRoute(t *testing.T) {
	// gomockのコントローラーを作成
	ctrl := gomock.NewController(t)
	// IRouteRepositoryのモックを作成
	mockRouteRepo := routeDomain.NewMockIRouteRepository(ctrl)
	// IUserRepositoryのモックを作成
	mockUserRepo := userDomain.NewMockIUserRepository(ctrl)
	// TransactionManagerのモックを作成
	mockTransactionManager := transactionApp.NewMockTransactionManager(ctrl)
	// テスト対象のUsecaseを作成
	uc := NewDeleteRouteUsecase(mockUserRepo, mockTransactionManager, mockRouteRepo)

	tests := []struct {
		name     string // description of this test case
		routeID  string
		kratosID string
		mockFunc func()
		wantErr  bool
	}{
		{
			name:     "正常系: ルート削除に成功する",
			routeID:  "019b5a50-0000-7000-8000-000000000001",
			kratosID: "2eb50f70-3a23-4067-99f6-9fd645686880",
			mockFunc: func() {
				// ユーザー取得のモック
				user, _ := userDomain.ReconstructUser(
					userDomain.UserID("019b5a8d-16a7-700a-be92-9ae11e7e5b9a"),
					"2eb50f70-3a23-4067-99f6-9fd645686880",
					"Test User",
					nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, false,
				)
				mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), "2eb50f70-3a23-4067-99f6-9fd645686880").
					Return(user, nil)

				// ルート取得のモック（所有者が一致）
				route, _ := routeDomain.ReconstructRoute(
					"019b5a50-0000-7000-8000-000000000001",
					"019b5a8d-16a7-700a-be92-9ae11e7e5b9a", // userIDが一致
					"Test Route",
					"Test Description",
					nil, 1000, 3600, 100, 50,
					routeDomain.Geometry{}, routeDomain.Geometry{},
					routeDomain.Geometry{}, routeDomain.Geometry{},
					"", 0, "", "",
				)
				mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), "019b5a50-0000-7000-8000-000000000001").
					Return(route, nil)

				// トランザクションのモック
				mockTransactionManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "異常系: ユーザーが見つからない",
			routeID:  "019b5a50-0000-7000-8000-000000000001",
			kratosID: "invalid-kratos-id",
			mockFunc: func() {
				mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), "invalid-kratos-id").
					Return(nil, errors.New("user not found"))
			},
			wantErr: true,
		},
		{
			name:     "異常系: ルートが見つからない",
			routeID:  "invalid-route-id",
			kratosID: "2eb50f70-3a23-4067-99f6-9fd645686880",
			mockFunc: func() {
				user, _ := userDomain.ReconstructUser(
					userDomain.UserID("019b5a8d-16a7-700a-be92-9ae11e7e5b9a"),
					"2eb50f70-3a23-4067-99f6-9fd645686880",
					"Test User",
					nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, false,
				)
				mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), "2eb50f70-3a23-4067-99f6-9fd645686880").
					Return(user, nil)

				mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), "invalid-route-id").
					Return(nil, errors.New("route not found"))
			},
			wantErr: true,
		},
		{
			name:     "異常系: ルートの所有者ではない（権限エラー）",
			routeID:  "019b5a50-0000-7000-8000-000000000001",
			kratosID: "2eb50f70-3a23-4067-99f6-9fd645686880",
			mockFunc: func() {
				user, _ := userDomain.ReconstructUser(
					userDomain.UserID("019b5a8d-16a7-700a-be92-9ae11e7e5b9a"),
					"2eb50f70-3a23-4067-99f6-9fd645686880",
					"Test User",
					nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, false,
				)
				mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), "2eb50f70-3a23-4067-99f6-9fd645686880").
					Return(user, nil)

				// ルートの所有者が異なる
				route, _ := routeDomain.ReconstructRoute(
					"019b5a50-0000-7000-8000-000000000001",
					"different-user-id", // 異なるユーザーID
					"Test Route",
					"Test Description",
					nil, 1000, 3600, 100, 50,
					routeDomain.Geometry{}, routeDomain.Geometry{},
					routeDomain.Geometry{}, routeDomain.Geometry{},
					"", 0, "", "",
				)
				mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), "019b5a50-0000-7000-8000-000000000001").
					Return(route, nil)
			},
			wantErr: true,
		},
		{
			name:     "異常系: トランザクション内での削除に失敗",
			routeID:  "019b5a50-0000-7000-8000-000000000001",
			kratosID: "2eb50f70-3a23-4067-99f6-9fd645686880",
			mockFunc: func() {
				user, _ := userDomain.ReconstructUser(
					userDomain.UserID("019b5a8d-16a7-700a-be92-9ae11e7e5b9a"),
					"2eb50f70-3a23-4067-99f6-9fd645686880",
					"Test User",
					nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, false,
				)
				mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), "2eb50f70-3a23-4067-99f6-9fd645686880").
					Return(user, nil)

				route, _ := routeDomain.ReconstructRoute(
					"019b5a50-0000-7000-8000-000000000001",
					"019b5a8d-16a7-700a-be92-9ae11e7e5b9a",
					"Test Route",
					"Test Description",
					nil, 1000, 3600, 100, 50,
					routeDomain.Geometry{}, routeDomain.Geometry{},
					routeDomain.Geometry{}, routeDomain.Geometry{},
					"", 0, "", "",
				)
				mockRouteRepo.EXPECT().
					GetRouteByID(gomock.Any(), "019b5a50-0000-7000-8000-000000000001").
					Return(route, nil)

				mockTransactionManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt      // ループ変数をキャプチャ（並列実行時の変数共有を防ぐ）
			t.Parallel()  // テストを並列実行
			tt.mockFunc() // モックの振る舞いを設定

			gotErr := uc.DeleteRoute(context.Background(), tt.routeID, tt.kratosID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DeleteRoute() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DeleteRoute() succeeded unexpectedly")
			}
		})
	}
}

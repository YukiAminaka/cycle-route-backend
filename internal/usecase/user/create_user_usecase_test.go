package usecase

import (
	"context"
	"testing"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"
)

// ptr は文字列のポインタを返すヘルパー関数
// テストデータでポインタ型のフィールドに値を設定する際に使用
func ptr(s string) *string {
	return &s
}

func Test_createUserUsecase_CreateUser(t *testing.T) {
	// gomockのコントローラーを作成
	ctrl := gomock.NewController(t)
	// IUserRepositoryのモックを作成
	mockUserRepo := userDomain.NewMockIUserRepository(ctrl)
	// テスト対象のUsecaseを作成
	uc := NewCreateUserUsecase(mockUserRepo)

	// テストケースの定義
	tests := []struct {
		name     string                      // テストケースの説明
		input    CreateUserUseCaseInputDto   // Usecaseへの入力データ
		mockFunc func()                      // モックの振る舞いを定義する関数
		want     *CreateUserUseCaseOutputDto // 期待される出力
		wantErr  bool                        // エラーが期待されるかどうか
	}{
		{
			name: "正常系: ユーザー作成に成功する",
			input: CreateUserUseCaseInputDto{
				KratosID: "2eb50f70-3a23-4067-99f6-9fd645686880",
				Name:      "Test User",
				FirstName: ptr("Test"),
				LastName:  ptr("User"),
				Email:     ptr("test@example.com"),
			},
			mockFunc: func() {
				// モックの振る舞いを定義
				mockUserRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()). // 期待される引数の型を指定
					DoAndReturn(func(ctx context.Context, user *userDomain.User) (*userDomain.User, error) {
						// 渡されたuserをそのまま返す（実際のリポジトリの挙動を模倣）
						// ReconstructUserを使用してデフォルト値を含む完全なユーザーを返す
						return userDomain.ReconstructUser(
							userDomain.UserID("019b5a8d-16a7-700a-be92-9ae11e7e5b9a"),
							"2eb50f70-3a23-4067-99f6-9fd645686880",
							"Test User",
							nil,                     // highlightedPhotoID (schema.sql: DEFAULT 0、nilで表現)
							ptr("ja"),               // locale (schema.sql: DEFAULT 'ja')
							nil,                     // description
							nil,                     // locality
							nil,                     // administrativeArea
							ptr("JP"),               // countryCode (schema.sql: DEFAULT 'JP')
							nil,                     // postalCode
							nil,                     // geom
							ptr("Test"),             // firstName
							ptr("User"),             // lastName
							ptr("test@example.com"), // email
							false,                   // hasSetLocation (schema.sql: DEFAULT FALSE)
						)
					})
			},
			want: &CreateUserUseCaseOutputDto{
				ID:                 "019b5a8d-16a7-700a-be92-9ae11e7e5b9a",
				Name:               "Test User",
				FirstName:          ptr("Test"),
				LastName:           ptr("User"),
				Email:              ptr("test@example.com"), // 修正: 正しいメールアドレスに
				HighlightedPhotoID: nil,
				Locale:             ptr("ja"),  // schema.sql: DEFAULT 'ja'
				Description:        nil,
				Locality:           nil,
				AdministrativeArea: nil,
				CountryCode:        ptr("JP"),  // schema.sql: DEFAULT 'JP'
				PostalCode:         nil,
				Geom:               nil,
				HasSetLocation:     false,      // schema.sql: DEFAULT FALSE
			},
			wantErr: false,
		},
	}
	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt        // ループ変数をキャプチャ（並列実行時の変数共有を防ぐ）
			t.Parallel()    // テストを並列実行
			tt.mockFunc()   // モックの振る舞いを設定

			// テスト対象のメソッドを実行
			got, err := uc.CreateUser(context.Background(), tt.input)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 結果の検証（IDフィールドは自動生成されるため比較から除外）
			// cmp.Diff(a, b, ...) は a と b の差分を文字列で返す関数
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(CreateUserUseCaseOutputDto{}, "ID")); diff != "" {
				t.Errorf("Run() diff = %v", diff)
			}
		})
	}
}

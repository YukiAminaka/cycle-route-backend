package user

import (
	"context"
	"errors"
	"testing"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
)

const (
	testUserID             = "019b5a8d-16a7-700a-be92-9ae11e7e5b9a"
	testKratosID           = "2eb50f70-3a23-4067-99f6-9fd645686880"
	testUserName           = "Test User"
	testHighlightedPhotoID = int64(12345)
	testLocale             = "ja"
	testDescription        = "test description"
	testLocality           = "shibuya-ku"
	testAdministrativeArea = "tokyo"
	testCountryCode        = "JP"
	testPostalCode         = "150-0002"
	testFirstName          = "Test"
	testLastName           = "User"
	testEmail              = "test@example.com"
	testHasSetLocation     = true
)

func createTestUserByID(userID string) *userDomain.User {
	user, _ := userDomain.ReconstructUser(
		userDomain.UserID(userID),
		testKratosID,
		testUserName,
		new(testHighlightedPhotoID),
		new(testLocale),
		new(testDescription),
		new(testLocality),
		new(testAdministrativeArea),
		new(testCountryCode),
		new(testPostalCode),
		nil,
		new(testFirstName),
		new(testLastName),
		new(testEmail),
		testHasSetLocation,
	)
	return user
}

func createTestUserByKratosID(kratosID string) *userDomain.User {
	user, _ := userDomain.ReconstructUser(
		userDomain.UserID(testUserID),
		kratosID,
		testUserName,
		new(testHighlightedPhotoID),
		new(testLocale),
		new(testDescription),
		new(testLocality),
		new(testAdministrativeArea),
		new(testCountryCode),
		new(testPostalCode),
		nil,
		new(testFirstName),
		new(testLastName),
		new(testEmail),
		testHasSetLocation,
	)
	return user
}

type getUserTestMocks struct {
	mockUserRepo *userDomain.MockIUserRepository
	usecase      IGetUserByIDUsecase
}

func setupGetUserByIDMocks(t *testing.T) *getUserTestMocks {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockUserRepo := userDomain.NewMockIUserRepository(ctrl)
	uc := NewGetUserByIDUsecase(mockUserRepo)
	return &getUserTestMocks{
		mockUserRepo: mockUserRepo,
		usecase:      uc,
	}
}

func Test_getUserByIDUsecase_GetUserByID(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		setupMocks func(m *getUserTestMocks)
		want       *GetUserByIDUseCaseDto
		wantErr    bool
	}{
		{
			name:   "正常系: 有効なIDでユーザー情報が取得できること",
			userID: testUserID,
			setupMocks: func(m *getUserTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByID(gomock.Any(), testUserID).
					Return(createTestUserByID(testUserID), nil)
			},
			want: &GetUserByIDUseCaseDto{
				ID:                 testUserID,
				Name:               testUserName,
				HighlightedPhotoID: new(testHighlightedPhotoID),
				Description:        new(testDescription),
				Locality:           new(testLocality),
				AdministrativeArea: new(testAdministrativeArea),
				CountryCode:        new(testCountryCode),
			},
			wantErr: false,
		},
		{
			name:   "異常系: 存在しないIDでエラーが返されること",
			userID: "non-existent-id",
			setupMocks: func(m *getUserTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByID(gomock.Any(), "non-existent-id").
					Return(nil, errors.New("user not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := setupGetUserByIDMocks(t)
			tt.setupMocks(mocks)

			got, err := mocks.usecase.GetUserByID(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUserByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getUserByIDUsecase_GetUserByKratosID(t *testing.T) {
	tests := []struct {
		name       string
		kratosID   string
		setupMocks func(m *getUserTestMocks)
		want       *GetUserByKratosIDUsecaseDto
		wantErr    bool
	}{
		{
			name:     "正常系: 有効なKratosIDでユーザー情報が取得できること",
			kratosID: testKratosID,
			setupMocks: func(m *getUserTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), testKratosID).
					Return(createTestUserByKratosID(testKratosID), nil)
			},
			want: &GetUserByKratosIDUsecaseDto{
				ID:                 testUserID,
				Name:               testUserName,
				HighlightedPhotoID: new(testHighlightedPhotoID),
				Locale:             new(testLocale),
				Description:        new(testDescription),
				Locality:           new(testLocality),
				AdministrativeArea: new(testAdministrativeArea),
				CountryCode:        new(testCountryCode),
				PostalCode:         new(testPostalCode),
				Geom:               nil,
				FirstName:          new(testFirstName),
				LastName:           new(testLastName),
				Email:              new(testEmail),
				HasSetLocation:     testHasSetLocation,
			},
			wantErr: false,
		},
		{
			name:     "異常系: 存在しないKratosIDでエラーが返されること",
			kratosID: "non-existent-kratos-id",
			setupMocks: func(m *getUserTestMocks) {
				m.mockUserRepo.EXPECT().
					GetUserByKratosID(gomock.Any(), "non-existent-kratos-id").
					Return(nil, errors.New("user not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := setupGetUserByIDMocks(t)
			tt.setupMocks(mocks)

			got, err := mocks.usecase.GetUserByKratosID(context.Background(), tt.kratosID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByKratosID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUserByKratosID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

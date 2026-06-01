package repository

import (
	"context"
	"fmt"
	"testing"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/google/go-cmp/cmp"
	"github.com/paulmach/orb"
)

func TestUserRepository_GetUserByID(t *testing.T) {
	user, err := userDomain.ReconstructUser(
		"70d6037a-b67b-4aa8-b5a3-da393b514f24",
		"019b5a3b-9854-787d-8877-e1732595d5b8",
		"testuser",
		nil, // highlightedPhotoID
		new("ja"),
		new("東京を中心にサイクリングを楽しんでいます。週末ライダーです。"),
		new("渋谷区"),
		new("東京都"),
		new("JP"),
		new("150-0002"),
		&userDomain.Geometry{Geometry: orb.Point{139.7024,35.6598}}, 
		new("太郎"),
		new("山田"),
		new("test@example.com"),
		true,
	)
	userID := user.ID().String()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		want *userDomain.User
	}{
		{
			name: "IDによってユーザーが取得ができること",
			want: user,
		},
	}
	q := GetTestQueries()
	userRepository := NewUserRepository(q)
	ctx := context.Background()
	resetTestData(t)
	for _, tt := range tests {
		t.Run(fmt.Sprintf(": %s", tt.name), func(t *testing.T) {
			got, err := userRepository.GetUserByID(ctx, userID)
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(got.KratosID(), tt.want.KratosID()); diff != "" {
				t.Errorf("FindById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_GetUserByKratosID(t *testing.T) {
	user, err := userDomain.ReconstructUser(
		"70d6037a-b67b-4aa8-b5a3-da393b514f24",
		"019b5a3b-9854-787d-8877-e1732595d5b8",
		"testuser",
		nil, // highlightedPhotoID
		new("ja"),
		new("東京を中心にサイクリングを楽しんでいます。週末ライダーです。"),
		new("渋谷区"),
		new("東京都"),
		new("JP"),
		new("150-0002"),
		&userDomain.Geometry{Geometry: orb.Point{139.7024,35.6598}}, 
		new("太郎"),
		new("山田"),
		new("test@example.com"),
		true,
	)
	kratosID := user.KratosID()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		want *userDomain.User
	}{
		{
			name: "IDによってユーザーが取得ができること",
			want: user,
		},
	}
	q := GetTestQueries()
	userRepository := NewUserRepository(q)
	ctx := context.Background()
	resetTestData(t)
	for _, tt := range tests {
		t.Run(fmt.Sprintf(": %s", tt.name), func(t *testing.T) {
			got, err := userRepository.GetUserByKratosID(ctx, kratosID)
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(got.ID(), tt.want.ID()); diff != "" {
				t.Errorf("FindById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_UpdateUserProfile(t *testing.T) {
	const userID = "70d6037a-b67b-4aa8-b5a3-da393b514f24"
	const kratosID = "019b5a3b-9854-787d-8877-e1732595d5b8"

	updatedUser, err := userDomain.ReconstructUser(
		userDomain.UserID(userID),
		kratosID,
		"updateduser",
		nil,
		new("ja"),
		new("更新後の説明文"),
		nil, nil, nil, nil, nil,
		new("更新太郎"),
		new("更新山田"),
		new("test@example.com"),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		input *userDomain.User
	}{
		{
			name:  "プロフィールが更新されること",
			input: updatedUser,
		},
	}

	q := GetTestQueries()
	userRepository := NewUserRepository(q)
	ctx := context.Background()
	resetTestData(t)

	for _, tt := range tests {
		t.Run(fmt.Sprintf(": %s", tt.name), func(t *testing.T) {
			if err := userRepository.UpdateUserProfile(ctx, tt.input); err != nil {
				t.Fatalf("UpdateUserProfile failed: %v", err)
			}

			got, err := userRepository.GetUserByID(ctx, userID)
			if err != nil {
				t.Fatalf("GetUserByID failed: %v", err)
			}

			if diff := cmp.Diff(got.Name(), tt.input.Name()); diff != "" {
				t.Errorf("Name mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(got.Description(), tt.input.Description()); diff != "" {
				t.Errorf("Description mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(got.FirstName(), tt.input.FirstName()); diff != "" {
				t.Errorf("FirstName mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(got.LastName(), tt.input.LastName()); diff != "" {
				t.Errorf("LastName mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_UpdateUserLocation(t *testing.T) {
	const userID = "70d6037a-b67b-4aa8-b5a3-da393b514f24"
	const kratosID = "019b5a3b-9854-787d-8877-e1732595d5b8"

	updatedUser, err := userDomain.ReconstructUser(
		userDomain.UserID(userID),
		kratosID,
		"testuser",
		nil,
		new("ja"),
		new("東京を中心にサイクリングを楽しんでいます。週末ライダーです。"),
		new("新宿区"),
		new("東京都"),
		new("JP"),
		new("160-0001"),
		&userDomain.Geometry{Geometry: orb.Point{139.7031, 35.6938}},
		new("太郎"),
		new("山田"),
		new("test@example.com"),
		true,
	)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		input *userDomain.User
	}{
		{
			name:  "位置情報が更新されること",
			input: updatedUser,
		},
	}

	q := GetTestQueries()
	userRepository := NewUserRepository(q)
	ctx := context.Background()
	resetTestData(t)

	for _, tt := range tests {
		t.Run(fmt.Sprintf(": %s", tt.name), func(t *testing.T) {
			if err := userRepository.UpdateUserLocation(ctx, tt.input); err != nil {
				t.Fatalf("UpdateUserLocation failed: %v", err)
			}

			got, err := userRepository.GetUserByID(ctx, userID)
			if err != nil {
				t.Fatalf("GetUserByID failed: %v", err)
			}

			if diff := cmp.Diff(got.Locality(), tt.input.Locality()); diff != "" {
				t.Errorf("Locality mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(got.AdministrativeArea(), tt.input.AdministrativeArea()); diff != "" {
				t.Errorf("AdministrativeArea mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(got.CountryCode(), tt.input.CountryCode()); diff != "" {
				t.Errorf("CountryCode mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(got.PostalCode(), tt.input.PostalCode()); diff != "" {
				t.Errorf("PostalCode mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	user, err := userDomain.NewUser("019b5a40-2c63-7c96-a2d2-a8f1ed21ecbd", "newuser", new("newuser@example.com"), new("Test"), new("User"))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name  string
		input *userDomain.User
		want  *userDomain.User
	}{
		{
			name:  "保存かつ取得ができること",
			input: user,
			want:  user,
		},
	}
	q := GetTestQueries()
	userRepository := NewUserRepository(q)
	ctx := context.Background()
	resetTestData(t)
	for _, tt := range tests {
		t.Run(fmt.Sprintf(": %s", tt.name), func(t *testing.T) {
			got, err := userRepository.CreateUser(ctx, tt.input)
			if err != nil {
				t.Fatalf("CreateUser failed: %v", err)
			}

			if got == nil {
				t.Fatal("CreateUser returned nil user")
			}

			if diff := cmp.Diff(got.KratosID(), tt.want.KratosID()); diff != "" {
				t.Errorf("FindById() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

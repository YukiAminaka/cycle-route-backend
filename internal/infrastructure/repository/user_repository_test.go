package repository

import (
	"context"
	"fmt"
	"testing"

	userDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/google/go-cmp/cmp"
	"github.com/paulmach/orb"
)

func ptr(s string) *string {
	return &s
}

func ptrGeom(g orb.Geometry) *userDomain.Geometry {
	return &userDomain.Geometry{Geometry: g}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	user, err := userDomain.ReconstructUser(
		"70d6037a-b67b-4aa8-b5a3-da393b514f24",
		"019b5a3b-9854-787d-8877-e1732595d5b8",
		"testuser",
		nil, // highlightedPhotoID
		ptr("ja"),
		ptr("東京を中心にサイクリングを楽しんでいます。週末ライダーです。"),
		ptr("渋谷区"),
		ptr("東京都"),
		ptr("JP"),
		ptr("150-0002"),
		ptrGeom(orb.Point{139.7024,35.6598}), 
		ptr("太郎"),
		ptr("山田"),
		ptr("test@example.com"),
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

func TestUserRepository_CreateUser(t *testing.T) {
	user, err := userDomain.NewUser("019b5a40-2c63-7c96-a2d2-a8f1ed21ecbd", "newuser",ptr("newuser@example.com"),ptr("Test"),ptr("User"))
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
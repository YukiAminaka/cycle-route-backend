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

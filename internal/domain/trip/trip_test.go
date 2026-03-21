package trip

import (
	"testing"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
)


func Test_newTrip(t *testing.T) {
	userID := user.NewUserID()
	type args struct {
		userID         string
		name           string
		description    string
		visibility     int16
		activityTypeID int32
	}
	tests := []struct {
		name string 
		args args
		wantErr        bool
	}{
		{
			name: "正常系すべての引数が有効",
			args: args{
				userID:         userID.String(),
				name:           "Trip to the park",
				description:    "A nice trip to the park",
				visibility:     1,
				activityTypeID: 2,
			},
			wantErr: false,
		},
		{
			name: "異常系 userIDが空",
			args: args{
				userID:         "",
				name:           "Trip to the park",
				description:    "A nice trip to the park",
				visibility:     1,
				activityTypeID: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTrip(tt.args.userID, tt.args.name, tt.args.description, tt.args.visibility, tt.args.activityTypeID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("newTrip() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newTrip() succeeded unexpectedly")
			}
			
			if got.userID != tt.args.userID {
				t.Errorf("userID = %v, want %v", got.userID, tt.args.userID)
			}
			if got.name != tt.args.name {
				t.Errorf("name = %v, want %v", got.name, tt.args.name)
			}
			if got.description != tt.args.description {
				t.Errorf("description = %v, want %v", got.description, tt.args.description)
			}
			if got.visibility != tt.args.visibility {
				t.Errorf("visibility = %v, want %v", got.visibility, tt.args.visibility)
			}
			if got.activityTypeID != tt.args.activityTypeID {
				t.Errorf("activityTypeID = %v, want %v", got.activityTypeID, tt.args.activityTypeID)
			}
		})
	}
}

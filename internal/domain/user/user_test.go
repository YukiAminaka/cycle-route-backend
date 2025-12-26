package user

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// 文字列のポインタを返すヘルパー関数
func ptr(s string) *string {
	return &s
}

func TestNewUser(t *testing.T) {
	type args struct {
		kratosID string
		name      string
		email     *string
		firstName *string
		lastName  *string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "正常系",
			args: args{
				kratosID: "2eb50f70-3a23-4067-99f6-9fd645686880",
				name:      "John Doe",
				email:     ptr("test@example.com"),  // ポインタ型
				firstName: ptr("John"),             
				lastName:  ptr("Doe"),               
			},
			want: &User{
				kratosID:      "2eb50f70-3a23-4067-99f6-9fd645686880",
				name:           "John Doe",
				highlightedPhotoID: nil,
				locale:         nil,
				description:    nil,
				locality:       nil,
				administrativeArea: nil,
				countryCode:    nil,
				postalCode:     nil,
				geom:           nil,
				email:          ptr("test@example.com"),  
				firstName:      ptr("John"),
				lastName:       ptr("Doe"),
				hasSetLocation: false,
			}, 
			wantErr: false,
		}, 
		{name: "異常系",
			args: args{
				kratosID: "",
				name:      "", 
				email:     nil, 
				firstName: nil, 
				lastName:  nil,
			},
			want:    nil, 
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.kratosID, tt.args.name, tt.args.email, tt.args.firstName, tt.args.lastName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(
				got, tt.want,
				cmp.AllowUnexported(User{}), // 非公開フィールドを比較可能にする
				cmpopts.IgnoreFields(User{}, "id"), // idフィールドは無視する
			)
			if diff != "" {
				t.Errorf("NewUser() = %v, want %v. error is %s", got, tt.want, err)
			}
		})
	}
}

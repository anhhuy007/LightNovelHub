package route

import (
	"Lightnovel/model"
	"testing"
)

func TestIsPasswordValid(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Ok password", args{"12345678"}, true},
		{"Short password", args{"123"}, false},
		{
			"Long password",
			args{"1234567890123456789012345678901234567890123456789012345678901234567890abc"},
			false,
		},
		{
			"Ok password",
			args{"123456789012345678901234567890123456789012345678901234567890123456789"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPasswordValid(tt.args.password); got != tt.want {
				t.Errorf("IsPasswordValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUsernameValid(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Good username", args{"thong"}, true},
		{"Contain strange character", args{"@#$%^&******"}, false},
		{"Short username", args{"th"}, false},
		{"Long username", args{"Repeat3timesRepeat3timesRepeat3times"}, false},
		{"Contain space", args{"thong nguyen"}, false},
		{"Contain utf8 character", args{"thông"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUsernameValid(tt.args.username); got != tt.want {
				t.Errorf("IsUsernameValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkUserMetadata(t *testing.T) {
	tests := []struct {
		name    string
		args    model.UserMetadata
		want    bool
		wantErr ErrorCode
	}{
		{
			"Good input",
			model.UserMetadata{"Thong", "Thong", "thong@mail.com", "12345678"},
			true,
			BadInput,
		},
		{
			"Bad username",
			model.UserMetadata{"thong nguyen", "Thong", "thong@mail.com", "12345678"},
			false,
			BadUsername,
		},
		{
			"Bad displayname",
			model.UserMetadata{"Thong", "T", "thong@mail.com", "12345678"},
			false,
			BadDisplayname,
		},
		{
			"Bad email",
			model.UserMetadata{"Thong", "Thong", "thongmail.com", "12345678"},
			false,
			BadEmail,
		},
		{
			"Display name containt unprintable character",
			model.UserMetadata{"Thong", "Thong\n", "thong@thong.com", "12345678"},
			false,
			BadDisplayname,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkUserMetadata(tt.args)
			if got != tt.want {
				t.Errorf("checkUserMetadata() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantErr {
				t.Errorf("checkUserMetadata() got1 = %v, want %v", got1, tt.wantErr)
			}
		})
	}
}

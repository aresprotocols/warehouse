package util

import "testing"

func TestMd5Str(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{str: "password"},
			want: "5f4dcc3b5aa765d61d8327deb882cf99",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5Str(tt.args.str); got != tt.want {
				t.Errorf("Md5Str() = %v, want %v", got, tt.want)
			}
		})
	}
}

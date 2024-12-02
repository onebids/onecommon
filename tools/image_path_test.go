package tools

import "testing"

func TestConvImagePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "tst", args: args{
			path: "abcc.jpg",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvImagePath(tt.args.path); got != tt.want {
				t.Errorf("ConvImagePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

package random

import "testing"

func TestNewRandomString(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"6 sym",
			args{
				6,
			},
			6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomString(tt.args.size)
			if len(got) != tt.want {
				t.Errorf("NewRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}

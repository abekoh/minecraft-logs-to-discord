package main

import "testing"

func Test_removePrefix(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "remove prefix",
			args: args{
				message: "[15:27:33] [Server thread/INFO]: abekoh lost connection: Disconnected",
			},
			want: "abekoh lost connection: Disconnected",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removePrefix(tt.args.message); got != tt.want {
				t.Errorf("removePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

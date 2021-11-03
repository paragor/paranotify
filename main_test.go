package main

import (
	"reflect"
	"testing"
)

func Test_splitStringByChunks(t *testing.T) {
	type args struct {
		body string
		size int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "short",
			args: args{
				body: "asdf",
				size: 10,
			},
			want: []string{"asdf"},
		},
		{
			name: "long",
			args: args{
				body: "asdfg",
				size: 2,
			},
			want: []string{"as", "df", "g"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitStringByChunks(tt.args.body, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitStringByChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}

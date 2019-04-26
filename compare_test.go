package describe

import (
	"fmt"
	"io"
	"testing"
)

func TestDiffFunc(t *testing.T) {
	type args struct {
		f func(out io.Writer, a, b string)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "basic test",
			args: args{
				f: func(o io.Writer, a, b string) {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DiffFunc(tt.args.f)
		})
	}
}

func TestCompare(t *testing.T) {
	type args struct {
		a interface{}
		b interface{}
	}
	tests := []struct {
		name     string
		args     args
		diffFunc func(out io.Writer, a, b string)
		want     bool
	}{
		{
			name: "basic test",
			args: args{
				a: "matching value",
				b: "matching value",
			},
			want: true,
		},
		{
			name: "basic test, different values",
			args: args{
				a: "differing value 1",
				b: "differing value 2",
			},
			want: false,
		},
		{
			name: "diff func",
			args: args{
				a: "differing value 1",
				b: "differing value 2",
			},
			diffFunc: func(out io.Writer, a, b string) { fmt.Fprintf(out, "%s != %s\n", a, b) },
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DiffFunc(tt.diffFunc)
			if got := Compare(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

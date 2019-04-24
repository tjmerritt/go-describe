package diff

import (
	"bytes"
	"testing"
)

func Test_diffFunc(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name  string
		args  args
		wantF string
	}{
		{
			name: "basic test",
			args: args{
				a: "1",
				b: "2",
			},
			wantF: `*** Got
--- Want
***************
*** 1 ****
! 1
--- 1 ----
! 2
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &bytes.Buffer{}
			diffFunc(f, tt.args.a, tt.args.b)
			if gotF := f.String(); gotF != tt.wantF {
				t.Errorf("diffFunc() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

package describe

import (
	"io"
	"os"
)

var diffFunc func(out io.Writer, a, b string)

// SetDiffFunc set a function that Compare will use to compute and report differences between two strings
func SetDiffFunc(f func(out io.Writer, a, b string)) {
	diffFunc = f
}

// Compare converts two values to there initialization format and then optionally output a diff betwen the two
// representations if the they are different.  Returns true if the representations of the two values are the same.
func Compare(a, b interface{}) bool {
	astr := Value(a)
	bstr := Value(b)
	if astr == bstr {
		return true
	}

	if diffFunc != nil {
		diffFunc(os.Stderr, astr, bstr)
	}

	return false
}

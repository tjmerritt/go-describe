package describe

import (
	"io"
	"os"
)

var diffFunc func(out io.Writer, a, b string)

func SetDiffFunc(f func(out io.Writer, a, b string)) {
	diffFunc = f
}

func Compare(a, b interface{}) bool {
	astr := DescribeValue(a)
	bstr := DescribeValue(b)
	if astr == bstr {
		return true
	}

	if diffFunc != nil {
		diffFunc(os.Stderr, astr, bstr)
	}

	return false
}

package diff

import (
	"fmt"
	"io"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/tjmerritt/go-describe"
)

func init() {
	describe.SetDiffFunc(diffFunc)

}

func diffFunc(f io.Writer, a, b string) {
	diff := difflib.ContextDiff{
		A:        difflib.SplitLines(a),
		B:        difflib.SplitLines(b),
		FromFile: "Got",
		ToFile:   "Want",
		Context:  3,
		Eol:      "\n",
	}
	result, _ := difflib.GetContextDiffString(diff)
	fmt.Fprintf(f, strings.Replace(result, "\t", " ", -1))
}

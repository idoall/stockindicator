package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

// go test -v ./utils/commonutils -run ^TestSkip$
func TestSkip(t *testing.T) {
	input := helpertools.Slice2Chan([]int{2, 4, 6, 8})
	expected := helpertools.Slice2Chan([]int{6, 8})

	result := helpertools.Skip(input, 2)

	err := helpertools.CheckEquals(result, expected)
	if err != nil {
		t.Fatal(err)
	}
}

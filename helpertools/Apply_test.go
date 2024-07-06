package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

// go test -v ./utils/commonutils -run ^TestApply$
func TestApply(t *testing.T) {
	input := helpertools.Slice2Chan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	expected := helpertools.Slice2Chan([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20})

	actual := helpertools.Apply(input, func(n int) int {
		return n * 2
	})

	err := helpertools.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

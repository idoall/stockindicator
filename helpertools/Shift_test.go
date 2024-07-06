package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

// go test -v ./utils/commonutils -run ^TestShift$
func TestShift(t *testing.T) {

	input := helpertools.Slice2Chan([]int{2, 4, 6, 8})
	expected := helpertools.Slice2Chan([]int{0, 0, 0, 0, 2, 4, 6, 8})

	actual := helpertools.Shift(input, 4, 0)

	err := helpertools.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

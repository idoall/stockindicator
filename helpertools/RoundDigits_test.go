package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

func TestRoundDigits(t *testing.T) {
	input := helpertools.Slice2Chan([]float64{10.1234, 5.678, 6.78, 8.91011})
	expected := helpertools.Slice2Chan([]float64{10.12, 5.68, 6.78, 8.91})

	actual := helpertools.RoundDigits(input, 2)

	err := helpertools.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

func TestRoundDigit(t *testing.T) {
	input := 10.1234
	expected := 10.12

	actual := helpertools.RoundDigit(input, 2)

	if actual != expected {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

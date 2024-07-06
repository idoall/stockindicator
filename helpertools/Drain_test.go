package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

// go test -v ./utils/commonutils -run ^TestDrain$
func TestDrain(_ *testing.T) {
	input := helpertools.Slice2Chan([]int{2, 4, 6, 8})
	helpertools.Drain(input)
}

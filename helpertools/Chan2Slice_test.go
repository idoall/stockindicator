package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

func TestChanToSlice(t *testing.T) {
	values := []int{2, 4, 6, 8}
	result := helpertools.Slice2Chan(values)

	chanResult := make(chan int, len(values))
	for _, n := range values {
		chanResult <- n
	}
	close(chanResult)

	err := helpertools.CheckEquals(chanResult, result)
	if err != nil {
		t.Fatal(err)
	}
}

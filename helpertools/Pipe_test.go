package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

func TestPipe(t *testing.T) {
	data := []int{2, 4, 6, 8}
	expected := helpertools.Slice2Chan(data)

	input := helpertools.Slice2Chan(data)
	actual := make(chan int)

	go helpertools.Pipe(input, actual)

	err := helpertools.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

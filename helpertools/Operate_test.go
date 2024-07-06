package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

// go test -v ./utils/commonutils -run ^TestOperate$
func TestOperate(t *testing.T) {
	ac := helpertools.Slice2Chan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helpertools.Slice2Chan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// {
	// 	expected := Slice2Chan([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20})

	// 	actual := Operate(ac, bc, func(a, b int) int {
	// 		return a + b
	// 	})

	// 	fmt.Println(Chan2Slice(expected))
	// 	fmt.Println(Chan2Slice(actual))
	// }

	expected := helpertools.Slice2Chan([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20})

	actual := helpertools.Operate(ac, bc, func(a, b int) int {
		return a + b
	})

	err := helpertools.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// go test -v ./utils/commonutils -run ^TestOperateFirstEnds$
func TestOperateFirstEnds(t *testing.T) {
	ac := helpertools.Slice2Chan([]int{1, 2, 3, 4, 5, 6, 7, 8})
	bc := helpertools.Slice2Chan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	expected := helpertools.Slice2Chan([]int{2, 4, 6, 8, 10, 12, 14, 16})

	actual := helpertools.Operate(ac, bc, func(a, b int) int {
		return a + b
	})

	err := helpertools.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// go test -v ./utils/commonutils -run ^TestOperateSecondEnds$
func TestOperateSecondEnds(t *testing.T) {
	ac := helpertools.Slice2Chan([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	bc := helpertools.Slice2Chan([]int{1, 2, 3, 4, 5, 6, 7, 8})

	expected := helpertools.Slice2Chan([]int{2, 4, 6, 8, 10, 12, 14, 16})

	actual := helpertools.Operate(ac, bc, func(a, b int) int {
		return a + b
	})

	err := helpertools.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

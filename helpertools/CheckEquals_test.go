package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

func TestCheckEqualsNotPairs(t *testing.T) {
	c := helpertools.Slice2Chan([]int{1, 2, 3, 4})

	err := helpertools.CheckEquals(c)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckEqualsNotEndedTheSame(t *testing.T) {
	a := helpertools.Slice2Chan([]int{1, 2, 3, 4})
	b := helpertools.Slice2Chan([]int{1, 2})

	err := helpertools.CheckEquals(a, b)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckEqualsNotEquals(t *testing.T) {
	a := helpertools.Slice2Chan([]int{1, 2, 3, 4})
	b := helpertools.Slice2Chan([]int{1, 3, 3, 4})

	err := helpertools.CheckEquals(a, b)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckEquals(t *testing.T) {
	a := helpertools.Slice2Chan([]int{1, 2, 3, 4})
	b := helpertools.Slice2Chan([]int{1, 2, 3, 4})

	err := helpertools.CheckEquals(a, b)
	if err != nil {
		t.Fatal(err)
	}
}

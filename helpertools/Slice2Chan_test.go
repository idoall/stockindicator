package helpertools_test

import (
	"reflect"
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

func TestSliceToChan(t *testing.T) {
	expected := []int{2, 4, 6, 8}
	actual := helpertools.Chan2Slice(helpertools.Slice2Chan(expected))

	// 判断两个长度是否一致
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}

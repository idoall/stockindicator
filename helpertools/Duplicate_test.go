package helpertools_test

import (
	"testing"

	"github.com/idoall/stockindicator/helpertools"
)

// go test -v ./utils/commonutils -run ^TestDuplicate$
func TestDuplicate(t *testing.T) {

	// temp1 := []float64{-10, 20, -4, -5}
	// temp2 := helpertools.Duplicate[float64](helpertools.Slice2Chan(temp1), 2)

	// for i, source := range temp1 {
	// 	for chanItemIndex, chanItem := range temp2 {
	// 		chainItemValue := <-chanItem
	// 		fmt.Printf("[%d][chanItemIndex:%d]source:%+v\tchainItemValue:%+v\n", i, chanItemIndex, source, chainItemValue)

	// 	}
	// }

	expecteds := []float64{-10, 20, -4, -5}

	outputs := helpertools.Duplicate[float64](helpertools.Slice2Chan(expecteds), 4)

	for i, expected := range expecteds {
		for _, output := range outputs {
			actual := <-output
			if actual != expected {
				t.Fatalf("index %d actual %v expected %v", i, actual, expected)
			}
		}
	}
}

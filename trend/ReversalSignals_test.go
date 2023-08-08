package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./trend -run TestReversalSignals
func TestReversalSignals(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKline()

	stock := NewReversalSignals(list)

	var dataList = stock.GetData()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 10; i < len(dataList)-1; i++ {

		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\tSide:%+v\n", i, v.Time.Format("2006-01-02 15:04:05"), v.Side.String())
	}

}

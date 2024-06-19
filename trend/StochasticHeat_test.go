package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// Run:
// go test -v ./trend -run TestStochasticHeat
func TestStochasticHeat(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	stock := NewDefaultStochasticHeat(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 100; i < len(dataList)-1; i++ {

		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\t Fast:%f\tSlow:%f\tSide:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), v.Fast, v.Slow, side.Data[i].String())
	}
}

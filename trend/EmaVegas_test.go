package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./trend -run TestEmaVegas
func TestEmaVegas(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKline()

	stock := NewDefaultEMAVegas(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-10 {
			break
		}
		var v = dataList[i]
		var klineItem = list[i]
		fmt.Printf("\t[%d]Time:%s\tClose:%f\tOpen:%f\tShort1Value:%f\tShort2Value:%f\tLong1Value:%f\tLong2Value:%f\tSide:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), klineItem.Close, klineItem.Open, v.Short1Value, v.Short2Value, v.Long1Value, v.Long2Value, side.Data[i].String())
	}

}

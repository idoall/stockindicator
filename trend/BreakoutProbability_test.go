package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./trend -run TestBreakoutProbability
func TestBreakoutProbability(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	stock := NewDefaultBreakoutProbability(list)

	var dataList = stock.GetData()

	// var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-10 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Price:%f\tTime:%s\tWin:%.2f\tLoss:%.2f\tProfitability:%.2f%%\n", i, list.Candles[i].Close, v.Time.Format("2006-01-02 15:04:05"), v.Win, v.Loss, v.Win/(v.Win+v.Loss)*100)
		fmt.Printf("\t\tBUY:\n\t\t\t")
		for _, vv := range v.List {
			if vv.Up {
				fmt.Printf("[%.2f/%.2f%%]\t", vv.Price, vv.Per)
			}
		}
		fmt.Printf("\n")
		fmt.Printf("\t\tSELL:\n\t\t\t")
		for _, vv := range v.List {
			if !vv.Up {
				fmt.Printf("[%.2f/%.2f%%]\t", vv.Price, vv.Per)
			}
		}
		fmt.Printf("\n")
	}

}

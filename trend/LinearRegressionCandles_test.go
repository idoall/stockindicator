package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// Run
// go test -v ./trend -test.run TestLinearRegressionCandles
func TestLinearRegressionCandles(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	//计算新的LinearRegressionCandles
	stock := NewDefaultLinearRegressionCandles(list)

	var dataList = stock.GetData()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-5 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d][%s]Open:%.2f\tClose:%.2f\tHigh:%.2f\tLow:%.2f\tSignal:%.2f\n", i, v.Time.Format("2006-01-02 15:04:05"), v.Open, v.Close, v.High, v.Low, v.Signal)
	}

}

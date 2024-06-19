package oscillator

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./oscillator -run TestStochasticOscillator
func TestStochasticOscillator(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	stock := NewDefaultStochasticOscillator(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", side.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-10 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\t K:%f\tD:%f\tSide:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), v.K, v.D, side.Data[i].String())
	}
}

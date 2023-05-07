package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// Run:
// go test -v ./trend -run TestUTBot
func TestUTBot(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKline()

	stock := NewDefaultUTBot(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s %d --\n", stock.Name, stock.Period)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-10 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\t Value:%f\tSide:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), v.Value, side.Data[i])
	}
}

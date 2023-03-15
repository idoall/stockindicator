package channel

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./channel -run TestKeltnerChannel
func TestKeltnerChannel(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKline()

	stock := NewDefaultKeltnerChannel(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s %d --\n", stock.Name, stock.Period)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-5 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\tUpper:%f\tMiddle:%f\tLower:%f\tSide:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), v.Upper, v.Middle, v.Lower, side.Data[i].String())
	}

}

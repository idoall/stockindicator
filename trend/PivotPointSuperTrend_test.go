package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./trend -run TestPivotPointSuperTrend
func TestPivotPointSuperTrend(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKline()

	stock := NewDefaultPivotPointSuperTrend(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 10; i < len(dataList)-1; i++ {

		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\tUp:%f\tUpTrendBegin:%f\tDown:%f\tDownTrendBegin:%f\tSide:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), v.UpTrend, v.UpTrendBegin, v.DownTrend, v.DownTrendBegin, side.Data[i].String())
	}

}

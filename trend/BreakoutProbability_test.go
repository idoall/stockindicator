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
	list := utils.GetTestKline()

	stock := NewDefaultBreakoutProbability(list)

	var dataList = stock.GetData()

	// var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-10 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\t%+v\n", i, v.Time.Format("2006-01-02 15:04:05"), v)
	}

}

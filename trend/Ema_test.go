package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./trend -run TestEma
func TestEma(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()
	//计算新的EMA
	// stock := NewDefaultEma(list)
	stock := NewEma(list, 200)

	dataList := stock.GetData()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-100 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\tEma%d:%f\n", i, v.Time.Format("2006-01-02 15:04:05"), stock.Period, v.Value)
	}
}

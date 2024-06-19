package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// Run
// go test -v ./trend -test.run TestKdj
func TestKdj(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	//计算新的Kdj
	stock := NewDefaultKdj(list)

	var dataList = stock.GetData()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-5 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d][%s]RSV:%.2f\tK:%.2f\tD:%.2f\tJ:%.2f\n", i, v.Time.Format("2006-01-02 15:04:05"), v.RSV, v.K, v.D, v.J)
	}

}

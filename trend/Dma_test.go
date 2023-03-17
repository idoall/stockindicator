package trend

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// Run:
// go test -v ./trend -run TestDma
func TestDma(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKline()

	//计算新的dma
	stock := NewDefaultDma(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-10 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\t %s:Dma:%f Ama:%f Side:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), stock.Name, v.Dma, v.Ama, side.Data[i])
	}

}

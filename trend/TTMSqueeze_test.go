package trend

import (
	"fmt"
	"testing"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Run:
// go test -v ./trend -run TestTTmSqueeze
func TestTTmSqueeze(t *testing.T) {
	loc, _ := time.LoadLocation("UTC-8")
	// handle err
	time.Local = loc // -> this is setting the global timezone

	t.Parallel()
	list := utils.GetTestKline()

	//计算新的Bbi
	stock := NewDefaultTTMSqueeze(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		//if i < len(dataList)-10 {
		//	break
		//}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\t %s:Mom:%f MomColor:%s SqzColor:%s Side:%s\n", i, v.Time.Format("2006-01-02 15:04:05"),
			stock.Name, v.Mom, v.MomColor, v.SqzColor, side.Data[i].PrintColorSide())
	}

}

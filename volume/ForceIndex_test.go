package volume

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./volume -run TestForceIndex
func TestForceIndex(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	stock := NewDefaultForceIndex(list)

	var dataList = stock.GetData()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-5 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\tPrice:%f\tValue:%f\n",
			i,
			v.Time.Format("2006-01-02 15:04:05"),
			list.Candles[i].Close,
			v.Value,
		)
	}
}

package channel

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./channel -run TestUlcerIndex
func TestUlcerIndex(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKline()

	stock := NewDefaultUlcerIndex(list)

	var dataList = stock.GetData()

	fmt.Printf("-- %s %d --\n", stock.Name, stock.Period)
	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-5 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Time:%s\tPrice:%f\tSMA:%f\n",
			i,
			v.Time.Format("2006-01-02 15:04:05"),
			list[i].Close,
			v.Value,
		)
	}
}

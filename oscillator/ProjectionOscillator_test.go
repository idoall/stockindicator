package oscillator

import (
	"fmt"
	"testing"

	"github.com/idoall/stockindicator/utils"
)

// RUN
// go test -v ./oscillator -run TestProjectionOscillator
func TestProjectionOscillator(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	stock := NewDefaultProjectionOscillator(list)

	var dataList = stock.GetData()

	var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", side.Name)
	for i, v := range dataList {
		if i > 5 {
			break
		}
		fmt.Printf("\t[%d]Time:%s\t Side:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), side.Data[i].String())
	}
}

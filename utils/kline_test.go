package utils

import (
	"fmt"
	"testing"
)

// RUN
// go test -v ./utils -run TestToHeikinAshi
func TestToHeikinAshi(t *testing.T) {
	t.Parallel()
	list := GetTestKline()

	var dataList = list.ToHeikinAshi()

	for i := len(dataList) - 1; i > 0; i-- {
		if i < len(dataList)-50 {
			break
		}
		var v = dataList[i]
		fmt.Printf("\t[%d]Kline:%+v\n",
			i,
			list[i],
		)

		fmt.Printf("\t[%d]HeiKinAshi:%+v\n",
			i,
			v,
		)
	}
}

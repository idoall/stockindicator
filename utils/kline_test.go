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

// RUN
// go test -v ./utils -run TestSortCandlesByTimestamp
func TestSortCandlesByTimestamp(t *testing.T) {
	t.Parallel()
	list := GetTestKline()

	var formatStr = "2006-01-02 15:04:05"
	list = list.SortCandlesByTimestamp(false)
	fmt.Printf("ASC\t%s\t%s\n", list[0].Time.Format(formatStr), list[len(list)-1].Time.Format(formatStr))

	list = list.SortCandlesByTimestamp(true)
	fmt.Printf("DESC\t%s\t%s\n", list[0].Time.Format(formatStr), list[len(list)-1].Time.Format(formatStr))
}

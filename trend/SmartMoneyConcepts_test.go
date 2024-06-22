package trend

import (
	"fmt"
	"testing"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Run:
// go test -v ./trend -run TestSmartMoneyConcepts
func TestSmartMoneyConcepts(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	loc, _ := time.LoadLocation("Local")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2023-05-01 00:00:00", loc)
	if err != nil {
		panic(err)
	}
	list.RemoveOutsideRange(startTime, time.Now())

	stock := NewDefaultSmartMoneyConcepts(list)

	// stock.EQHEQL_Enable = false
	// stock.Calculation()

	var dataList = stock.GetData()

	// var side = stock.AnalysisSide()

	fmt.Printf("-- %s --\n", stock.Name)
	for i := len(dataList) - 100; i < len(dataList)-1; i++ {

		var v = dataList[i]
		fmt.Printf("[%d][%s]H_BOS_S[%.2f]\tH_CHoCH_S[%.2f]\tL_BOS_S[%.2f]\tL_ChoCH_S[%.2f]H_BOS_L[%.2f]\tH_CHoCH_L[%.2f]\tL_BOS_L[%.2f]\tL_ChoCH_L[%.2f]\tEQH[%.2f]\tEQL[%.2f]\n",
			i,
			v.Time.Format("2006-01-02 15:04:05"),
			v.HighBOSShort,
			v.HighCHoCHShort,
			v.LowBOSShort,
			v.LowChoCHShort,
			v.HighBOSLong,
			v.HighCHoCHLong,
			v.LowBOSLong,
			v.LowChoCHLong,
			v.EQH,
			v.EQL,
		)
	}
	fmt.Printf("Strong High[%s]:%.2f\tWeak High[%s]:%.2f\tStrong Low[%s]:%.2f\tWeak Low[%s]:%.2f\tLastClose:%.2f\n",
		stock.StrongHigh.Time.Format("2006-01-02 15:04:05"),
		stock.StrongHigh.Value,
		stock.WeakHigh.Time.Format("2006-01-02 15:04:05"),
		stock.WeakHigh.Value,
		stock.StrongLow.Time.Format("2006-01-02 15:04:05"),
		stock.StrongLow.Value,
		stock.WeakLow.Time.Format("2006-01-02 15:04:05"),
		stock.WeakLow.Value,
		list.Candles[len(list.Candles)-1].Close,
	)
	fmt.Printf("OrderBlockBullish:%d\n", len(stock.OrderBlockBullish))
	for i, v := range stock.OrderBlockBullish {
		if i > 20 {
			break
		}
		fmt.Printf("[%d][%s]\tTop:%+v\tHigh:%.2f\tLow:%.2f\n", i, v.Time.Format("2006-01-02 15:04:05"), v.IsTop, v.High, v.Low)
	}
	fmt.Printf("OrderBlockBearish:%d\n", len(stock.OrderBlockBearish))
	for i, v := range stock.OrderBlockBearish {
		if i > 20 {
			break
		}
		fmt.Printf("[%d][%s]\tTop:%+v\tHigh:%.2f\tLow:%.2f\n", i, v.Time.Format("2006-01-02 15:04:05"), v.IsTop, v.High, v.Low)
	}
}

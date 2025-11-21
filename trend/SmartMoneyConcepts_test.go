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

	// 启用 FVG 检测
	stock.FVG_Enable = true
	stock.FVG_AutoThreshold = true
	stock.FVG_KeepHistory = true

	var dataList = stock.GetData()

	fmt.Printf("========== %s ==========\n", stock.Name)
	fmt.Println()

	// 输出 BOS/CHoCH/EQH/EQL 指标信息
	fmt.Printf("========== 市场结构指标 (最近100根K线) ==========\n")
	bosChochCount := 0
	for i := len(dataList) - 100; i < len(dataList)-1; i++ {
		var v = dataList[i]

		// 只显示有信号的K线
		hasSignal := false
		signalInfo := fmt.Sprintf("[%d][%s] ", i, v.Time.Format("01-02 15:04"))

		if v.HighBOSShort > 0 {
			signalInfo += fmt.Sprintf("顶部短线BOS:%.2f ", v.HighBOSShort)
			hasSignal = true
			bosChochCount++
		}
		if v.HighCHoCHShort > 0 {
			signalInfo += fmt.Sprintf("顶部短线CHoCH:%.2f ", v.HighCHoCHShort)
			hasSignal = true
			bosChochCount++
		}
		if v.LowBOSShort > 0 {
			signalInfo += fmt.Sprintf("底部短线BOS:%.2f ", v.LowBOSShort)
			hasSignal = true
			bosChochCount++
		}
		if v.LowChoCHShort > 0 {
			signalInfo += fmt.Sprintf("底部短线CHoCH:%.2f ", v.LowChoCHShort)
			hasSignal = true
			bosChochCount++
		}
		if v.HighBOSLong > 0 {
			signalInfo += fmt.Sprintf("顶部长线BOS:%.2f ", v.HighBOSLong)
			hasSignal = true
			bosChochCount++
		}
		if v.HighCHoCHLong > 0 {
			signalInfo += fmt.Sprintf("顶部长线CHoCH:%.2f ", v.HighCHoCHLong)
			hasSignal = true
			bosChochCount++
		}
		if v.LowBOSLong > 0 {
			signalInfo += fmt.Sprintf("底部长线BOS:%.2f ", v.LowBOSLong)
			hasSignal = true
			bosChochCount++
		}
		if v.LowChoCHLong > 0 {
			signalInfo += fmt.Sprintf("底部长线CHoCH:%.2f ", v.LowChoCHLong)
			hasSignal = true
			bosChochCount++
		}
		if v.EQH > 0 {
			signalInfo += fmt.Sprintf("相等高点:%.2f ", v.EQH)
			hasSignal = true
			bosChochCount++
		}
		if v.EQL > 0 {
			signalInfo += fmt.Sprintf("相等低点:%.2f ", v.EQL)
			hasSignal = true
			bosChochCount++
		}

		if hasSignal {
			fmt.Println(signalInfo)
		}
	}
	fmt.Printf("市场结构信号总数: %d\n", bosChochCount)
	fmt.Println()

	// 统计 FVG 产生和失效事件
	fmt.Printf("========== FVG 公允价值缺口 (最近100根K线) ==========\n")
	fvgNewCount := 0
	fvgFilledCount := 0

	for i := len(dataList) - 100; i < len(dataList)-1; i++ {
		var v = dataList[i]

		// 检测新 FVG 产生
		if v.NewBullishFVG.IsValid() {
			fvgNewCount++
			fmt.Printf("【新看涨FVG】[%d][%s] 区间: %.2f - %.2f\n",
				i,
				v.Time.Format("01-02 15:04"),
				v.NewBullishFVG.Bottom,
				v.NewBullishFVG.Top,
			)
		}

		if v.NewBearishFVG.IsValid() {
			fvgNewCount++
			fmt.Printf("【新看跌FVG】[%d][%s] 区间: %.2f - %.2f\n",
				i,
				v.Time.Format("01-02 15:04"),
				v.NewBearishFVG.Bottom,
				v.NewBearishFVG.Top,
			)
		}

		// 检测 FVG 被填补
		if v.FilledBullishFVG.IsValid() {
			fvgFilledCount++
			fmt.Printf("【看涨FVG失效】[%d][%s] 填补价格: %.2f, 持续: %d根K线\n",
				i,
				v.Time.Format("01-02 15:04"),
				v.FilledBullishFVG.FilledPrice,
				v.FilledBullishFVG.Duration,
			)
		}

		if v.FilledBearishFVG.IsValid() {
			fvgFilledCount++
			fmt.Printf("【看跌FVG失效】[%d][%s] 填补价格: %.2f, 持续: %d根K线\n",
				i,
				v.Time.Format("01-02 15:04"),
				v.FilledBearishFVG.FilledPrice,
				v.FilledBearishFVG.Duration,
			)
		}

		// 显示当前活跃的 FVG
		if v.BullishFVG.IsValid() && !v.NewBullishFVG.IsValid() {
			fmt.Printf("[%d][%s] 当前看涨FVG: %.2f - %.2f (产生于 %s)\n",
				i,
				v.Time.Format("01-02 15:04"),
				v.BullishFVG.Bottom,
				v.BullishFVG.Top,
				v.BullishFVG.Time.Format("01-02 15:04"),
			)
		}

		if v.BearishFVG.IsValid() && !v.NewBearishFVG.IsValid() {
			fmt.Printf("[%d][%s] 当前看跌FVG: %.2f - %.2f (产生于 %s)\n",
				i,
				v.Time.Format("01-02 15:04"),
				v.BearishFVG.Bottom,
				v.BearishFVG.Top,
				v.BearishFVG.Time.Format("01-02 15:04"),
			)
		}
	}

	fmt.Println()
	fmt.Printf("========== FVG 统计 ==========\n")
	fmt.Printf("新产生FVG数量: %d\n", fvgNewCount)
	fmt.Printf("失效FVG数量: %d\n", fvgFilledCount)
	fmt.Printf("历史FVG总数: %d\n", len(stock.FVG_History))

	// 统计 FVG 平均持续时间
	if len(stock.FVG_History) > 0 {
		totalDuration := 0
		bullishCount := 0
		bearishCount := 0
		for _, fvg := range stock.FVG_History {
			totalDuration += fvg.Duration
			if fvg.IsBullish {
				bullishCount++
			} else {
				bearishCount++
			}
		}
		avgDuration := float64(totalDuration) / float64(len(stock.FVG_History))
		fmt.Printf("FVG平均持续时间: %.1f根K线\n", avgDuration)
		fmt.Printf("看涨FVG: %d个, 看跌FVG: %d个\n", bullishCount, bearishCount)

		// 找出持续时间最长的 FVG
		var longestFVG SmartMoneyConceptsDataFairValueGap
		for _, fvg := range stock.FVG_History {
			if fvg.Duration > longestFVG.Duration {
				longestFVG = fvg
			}
		}
		fvgType := "看跌"
		if longestFVG.IsBullish {
			fvgType = "看涨"
		}
		fmt.Printf("最长持续FVG: %s, %d根K线, 区间: %.2f - %.2f\n",
			fvgType, longestFVG.Duration, longestFVG.Bottom, longestFVG.Top)
	}

	fmt.Println()
	fmt.Printf("========== 强弱高低点与订单区块 ==========\n")
	fmt.Printf("强高点[%s]:%.2f\t弱高点[%s]:%.2f\t强低点[%s]:%.2f\t弱低点[%s]:%.2f\t最新收盘:%.2f\n",
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
	fmt.Printf("看涨订单区块数量:%d\n", len(stock.OrderBlockBullish))
	for i, v := range stock.OrderBlockBullish {
		if i > 20 {
			break
		}
		topType := "顶部"
		if !v.IsTop {
			topType = "底部"
		}
		fmt.Printf("[%d][%s]\t类型:%s\t高点:%.2f\t低点:%.2f\n", i, v.Time.Format("2006-01-02 15:04:05"), topType, v.High, v.Low)
	}
	fmt.Printf("看跌订单区块数量:%d\n", len(stock.OrderBlockBearish))
	for i, v := range stock.OrderBlockBearish {
		if i > 20 {
			break
		}
		topType := "顶部"
		if !v.IsTop {
			topType = "底部"
		}
		fmt.Printf("[%d][%s]\t类型:%s\t高点:%.2f\t低点:%.2f\n", i, v.Time.Format("2006-01-02 15:04:05"), topType, v.High, v.Low)
	}
}

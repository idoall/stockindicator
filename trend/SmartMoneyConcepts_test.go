package trend

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/shopspring/decimal"
)

// ========== 测试层回测结构体 ==========

// Run:
// go test -v ./trend -run TestSmartMoneyConcepts
func TestSmartMoneyConcepts(t *testing.T) {
	t.Parallel()
	list := utils.GetTestKlineItem()

	loc, _ := time.LoadLocation("Local")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2025-01-01 00:00:00", loc)
	if err != nil {
		panic(err)
	}
	// end, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-11-22 04:00:00", loc)
	// list.RemoveOutsideRange(startTime, end)

	list.RemoveOutsideRange(startTime, time.Now())

	stock := NewDefaultSmartMoneyConcepts(list)
	output := ""

	// 启用 FVG 检测
	stock.FVG_Enable = true
	stock.FVG_AutoThreshold = true
	stock.FVG_KeepHistory = true

	// 强弱高代点检测
	stock.StrongWeak_Enable = true

	// 启用 CHoCH/BOS 历史记录
	stock.StructureBreak_Enable = true

	// 启用 EQH/EQL 历史记录（EQHEQL_Enable 已在构造函数中默认启用）
	stock.EQHEQL_Enable = true
	stock.EQHEQL_KeepHistory = false

	stock.Calculation()

	output += fmt.Sprintf("========== %s(K线数量:%d) ==========\n", stock.Name, len(list.Candles))
	output += "\n"

	allSignals := make([]SmartMoneyConceptsDataStructureBreak, 0)
	allSignals = append(allSignals, stock.HighBOSShortList...)
	allSignals = append(allSignals, stock.HighBOSLongList...)
	allSignals = append(allSignals, stock.HighCHoCHShortList...)
	allSignals = append(allSignals, stock.HighCHoCHLongList...)
	allSignals = append(allSignals, stock.LowBOSShortList...)
	allSignals = append(allSignals, stock.LowBOSLongList...)
	allSignals = append(allSignals, stock.LowCHoCHShortList...)
	allSignals = append(allSignals, stock.LowCHoCHLongList...)

	// ========== 步骤2: 信号优先级过滤 ==========
	filteredSignals := allSignals

	// 按出现时间倒排序
	sort.Slice(filteredSignals, func(i, j int) bool {
		return filteredSignals[i].Time.After(filteredSignals[j].Time)
	})

	// 遍历所有过滤后的信号
	minCount := len(filteredSignals)
	if minCount > 10 {
		minCount = 10
	}
	// 输出 BOS/CHoCH 指标信息
	output += fmt.Sprintf("========== BOS/CHoCH (最近%d条数据，总数:%d) ==========\n", 10, len(filteredSignals))

	for i := 0; i < minCount && i < len(filteredSignals); i++ {

		if i > minCount {
			output += "...\n"
		}

		signal := filteredSignals[i]
		signalName := fmt.Sprintf("%s%s%s", signal.Position, signal.Type, signal.Period)

		position := ""
		switch signal.Position {
		case PositionTypeHigh:
			position = "做多📈"
		case PositionTypeLow:
			position = "做空📉"
		}

		output += fmt.Sprintf("  %d. %s\t出现时间:[%s]\t价格:%s\t突破时间:[%s]\t突破时收盘价:%s\t持续: %d根K线\t%s\n",
			i+1,
			signalName,
			signal.Time.Format("2006-01-02 15:04:05"),
			signal.BreakPrice,
			signal.BreakTime.Format("2006-01-02 15:04:05"),
			signal.ClosePrice,
			signal.Duration,
			position,
			// signal,
		)
	}

	output += "\n"

	// 显示当前活跃的 FVG 列表
	output += "========== 当前活跃 FVG 列表 ==========\n"
	output += "\n"

	// 显示看涨 FVG 列表（最多显示最后 10 个）
	output += fmt.Sprintf("【看涨 FVG 列表】（共 %d 个）\n", len(stock.BullishFVGList))
	bullishDisplayCount := len(stock.BullishFVGList)
	if bullishDisplayCount > 10 {
		bullishDisplayCount = 10
	}
	for i := 0; i < bullishDisplayCount; i++ {
		fvg := stock.BullishFVGList[i]
		status := "活跃中"
		filledInfo := ""
		if fvg.FilledTime.Unix() > 0 {
			status = "已失效"
			filledInfo = fmt.Sprintf("| 填补时间: %s | 填补价格: %s | 持续: %d根K线",
				fvg.FilledTime.Format("2006-01-02 15:04:05"),
				fvg.FilledPrice,
				fvg.Duration,
			)
		}
		output += fmt.Sprintf("[%d] 产生时间: %s | 区间: %s - %s | 状态: %s %s\n",
			i,
			fvg.Time.Format("2006-01-02 15:04:05"),
			fvg.Bottom,
			fvg.Top,
			status,
			filledInfo,
		)
	}
	output += "\n"

	// 显示看跌 FVG 列表（最多显示最后 10 个）
	output += fmt.Sprintf("【看跌 FVG 列表】（共 %d 个）\n", len(stock.BearishFVGList))
	bearishDisplayCount := len(stock.BearishFVGList)
	if bearishDisplayCount > 10 {
		bearishDisplayCount = 10
	}
	for i := 0; i < bearishDisplayCount; i++ {
		fvg := stock.BearishFVGList[i]
		status := "活跃中"
		filledInfo := ""
		if fvg.FilledTime.Unix() > 0 {
			status = "已失效"
			filledInfo = fmt.Sprintf("| 填补时间: %s | 填补价格: %s | 持续: %d根K线",
				fvg.FilledTime.Format("2006-01-02 15:04:05"),
				fvg.FilledPrice,
				fvg.Duration,
			)
		}
		output += fmt.Sprintf("[%d] 产生时间: %s | 区间: %s - %s | 状态: %s %s\n",
			i,
			fvg.Time.Format("2006-01-02 15:04:05"),
			fvg.Bottom,
			fvg.Top,
			status,
			filledInfo,
		)
	}
	output += "\n"

	output += "========== FVG 统计 ==========\n"
	output += fmt.Sprintf("当前活跃看涨FVG数量: %d\n", len(stock.BullishFVGList))
	output += fmt.Sprintf("当前活跃看跌FVG数量: %d\n", len(stock.BearishFVGList))
	output += fmt.Sprintf("历史FVG总数: %d\n", len(stock.FVG_History))

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
		output += fmt.Sprintf("FVG平均持续时间: %.1f根K线\n", avgDuration)
		output += fmt.Sprintf("看涨FVG: %d个, 看跌FVG: %d个\n", bullishCount, bearishCount)

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
		output += fmt.Sprintf("最长持续FVG: %s, %d根K线, 区间: %s - %s\n",
			fvgType, longestFVG.Duration, longestFVG.Bottom, longestFVG.Top)
	}

	// 输出历史 FVG 详情（最多显示 20 个）
	if len(stock.FVG_History) > 0 {
		output += "\n"
		output += fmt.Sprintf("========== 历史 FVG 详情 (最近%d/%d个) ==========\n", min(20, len(stock.FVG_History)), len(stock.FVG_History))
		displayCount := min(20, len(stock.FVG_History))
		for i := 0; i < displayCount; i++ {
			// 倒序遍历：从最新的 FVG 开始显示
			fvg := stock.FVG_History[len(stock.FVG_History)-1-i]
			fvgType := "看跌"
			if fvg.IsBullish {
				fvgType = "看涨"
			}
			output += fmt.Sprintf("[%d] %s | 产生: %s | 失效: %s | 持续: %d根K线 | 区间: %s - %s | 填补价格: %s\n",
				i,
				fvgType,
				fvg.Time.Format("2006-01-02 15:04:05"),
				fvg.FilledTime.Format("2006-01-02 15:04:05"),
				fvg.Duration,
				fvg.Bottom.String(),
				fvg.Top.String(),
				fvg.FilledPrice.String(),
			)
		}
	}

	output += "\n"
	output += "========== 强弱高低点 ==========\n"

	if stock.StrongWeakDetail != nil {

		if stock.StrongWeakDetail.StrongHigh.Value.GreaterThan(decimal.Zero) {
			output += fmt.Sprintf("\t- 强高点(Strong High): %s (形成于 %s) [%s]\n", stock.StrongWeakDetail.StrongHigh.Value, stock.StrongWeakDetail.StrongHigh.Time.Format("2006-01-02 15:04:05"), stock.StrongWeakDetail.StrongHigh.Role)
		}
		if stock.StrongWeakDetail.WeakHigh.Value.GreaterThan(decimal.Zero) {
			output += fmt.Sprintf("\t- 弱高点(Weak High): %s (形成于 %s) [%s]\n", stock.StrongWeakDetail.WeakHigh.Value, stock.StrongWeakDetail.WeakHigh.Time.Format("2006-01-02 15:04:05"), stock.StrongWeakDetail.WeakHigh.Role)
		}
		if stock.StrongWeakDetail.StrongLow.Value.GreaterThan(decimal.Zero) {
			output += fmt.Sprintf("\t- 弱低点(Strong Low): %s (形成于 %s) [%s]\n", stock.StrongWeakDetail.StrongLow.Value, stock.StrongWeakDetail.StrongLow.Time.Format("2006-01-02 15:04:05"), stock.StrongWeakDetail.StrongLow.Role)
		}
		if stock.StrongWeakDetail.WeakLow.Value.GreaterThan(decimal.Zero) {
			output += fmt.Sprintf("\t- 弱低点(Weak Low): %s (形成于 %s) [%s]\n", stock.StrongWeakDetail.WeakLow.Value, stock.StrongWeakDetail.WeakLow.Time.Format("2006-01-02 15:04:05"), stock.StrongWeakDetail.WeakLow.Role)
		}
	} else {
		output += "○ 暂无强弱高低点数据\n"
	}

	output += "\n"
	output += "========== 订单区块 ==========\n"
	output += fmt.Sprintf("看涨订单区块数量:%d\n", len(stock.OrderBlockBullish))
	for i, v := range stock.OrderBlockBullish {
		if i > 20 {
			break
		}
		topType := "顶部"
		if !v.IsTop {
			topType = "底部"
		}
		output += fmt.Sprintf("[%d][%s]\t类型:%s\t高点:%s\t低点:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), topType, v.High.String(), v.Low.String())
	}
	output += fmt.Sprintf("看跌订单区块数量:%d\n", len(stock.OrderBlockBearish))
	for i, v := range stock.OrderBlockBearish {
		if i > 20 {
			break
		}
		topType := "顶部"
		if !v.IsTop {
			topType = "底部"
		}
		output += fmt.Sprintf("[%d][%s]\t类型:%s\t高点:%s\t低点:%s\n", i, v.Time.Format("2006-01-02 15:04:05"), topType, v.High.String(), v.Low.String())
	}

	output += "\n"
	output += "========== 流动性扫单(EQH/EQL) ==========\n"

	// 检查上方的 EQH（可能被扫单的流动性）
	for _, v := range stock.EQHList {
		output += fmt.Sprintf("EQH:%+v\n", v)
	}
	for _, v := range stock.EQLList {
		output += fmt.Sprintf("EQL:%+v\n", v)
	}

	// 输出历史 FVG 详情（最多显示 20 个）
	if len(stock.EQHEQL_History) > 0 {
		output += "\n"
		output += fmt.Sprintf("========== 历史 EQHEQL 详情 (最近%d/%d个) ==========\n", min(20, len(stock.EQHEQL_History)), len(stock.EQHEQL_History))
		displayCount := min(20, len(stock.EQHEQL_History))
		for i := 0; i < displayCount; i++ {
			// 倒序遍历：从最新的 FVG 开始显示
			eqheql := stock.EQHEQL_History[len(stock.EQHEQL_History)-1-i]

			output += fmt.Sprintf("[%d] %+v\n",
				i,
				eqheql,
			)
		}
	}

	fmt.Println(output)

}

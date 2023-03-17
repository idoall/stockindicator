package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// EMAVegas struct
type EMAVegas struct {
	Name         string
	PeriodShort1 int
	PeriodShort2 int
	PeriodLong1  int
	PeriodLong2  int
	data         []EMAVegasData
	kline        utils.Klines
}

// EMAVegasData EMAVegas函数计算给定期间的双指数移动平均线 (EMAVegas)。

// 双指数移动平均线 (EMAVegas) 是由 Patrick Mulloy 引入的技术指标。目的是减少技术交易者使用的价格图表中存在的噪音量。EMAVegas 使用两个指数移动平均线 (EMA) 来消除滞后。当价格高于平均水平时，它有助于确认上升趋势，当价格低于平均水平时，它有助于确认下降趋势。当价格超过平均线时，可能表示趋势发生变化。
type EMAVegasData struct {
	Short1Value float64
	Short2Value float64
	Long1Value  float64
	Long2Value  float64
	Time        time.Time
}

// NewEMAVegas new Func
func NewEMAVegas(list utils.Klines, periodShort1, periodShort2, periodLong1, periodLong2 int) *EMAVegas {
	m := &EMAVegas{
		Name:         fmt.Sprintf("EMAVegas%d-%d-%d-%d", periodShort1, periodShort2, periodLong1, periodLong2),
		kline:        list,
		PeriodShort1: periodShort1,
		PeriodShort2: periodShort2,
		PeriodLong1:  periodLong1,
		PeriodLong2:  periodLong2,
	}
	return m
}

// NewDefaultEMAVegas new Func
func NewDefaultEMAVegas(list utils.Klines) *EMAVegas {
	return NewEMAVegas(list, 144, 169, 575, 676)
}

// Calculation Func
func (e *EMAVegas) Calculation() *EMAVegas {

	e.data = make([]EMAVegasData, len(e.kline))

	// 计算 EMA1 值
	emaShort1 := NewEma(e.kline, e.PeriodShort1).GetValues()
	emaShort2 := NewEma(e.kline, e.PeriodShort2).GetValues()
	emaLong1 := NewEma(e.kline, e.PeriodLong1).GetValues()
	emaLong2 := NewEma(e.kline, e.PeriodLong2).GetValues()

	for i := 0; i < len(emaShort1); i++ {
		e.data[i] = EMAVegasData{
			Time:        e.kline[i].Time,
			Short1Value: emaShort1[i],
			Short2Value: emaShort2[i],
			Long1Value:  emaLong1[i],
			Long2Value:  emaLong2[i],
		}
	}
	return e
}

// AnalysisSide Func
func (e *EMAVegas) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}
		var prevKlineItem = e.kline[i-1]
		var klineItem = e.kline[i]
		prevItem := e.data[i-1]

		if klineItem.Close > klineItem.Open && ((klineItem.Close > v.Short2Value && prevKlineItem.Close < prevItem.Short2Value) || (klineItem.Close > v.Long2Value && prevKlineItem.Close < prevItem.Short2Value)) {
			sides[i] = utils.Buy
		} else if klineItem.Close < klineItem.Open && ((klineItem.Close < v.Short2Value && prevKlineItem.Close > prevItem.Short2Value) || (klineItem.Close < v.Long2Value && prevKlineItem.Close > prevItem.Short2Value)) {
			sides[i] = utils.Sell
		} else {
			sides[i] = utils.Hold
		}
	}
	return utils.SideData{
		Name: e.Name,
		Data: sides,
	}
}

// GetData return Point
func (e *EMAVegas) GetData() []EMAVegasData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

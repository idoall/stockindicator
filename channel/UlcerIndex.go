package channel

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Ulcer Index (UI) 溃疡指数是一种技术指标，可根据价格下跌的深度和持续时间来衡量下行风险。
// 该指数随着价格远离近期高点而增加，并随着价格升至新高而下跌。
// 该指标通常在 14 天内计算，溃疡指数显示交易者可以预期从该时期的高点回撤的百分比。
//
// High Closings = Max(period, Closings)
// Percentage Drawdown = 100 * ((Closings - High Closings) / High Closings)
// Squared Average = Sma(period, Percent Drawdown * Percent Drawdown)
// Ulcer Index = Sqrt(Squared Average)
type UlcerIndex struct {
	Name   string
	Period int
	data   []UlcerIndexData
	kline  *klines.Item
}

// UlcerIndexData
type UlcerIndexData struct {
	Time  time.Time
	Value float64
}

// NewUlcerIndex new Func
func NewUlcerIndex(klineItem *klines.Item, period int) *UlcerIndex {
	m := &UlcerIndex{
		Name:   fmt.Sprintf("UlcerIndex%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultUlcerIndex new Func
func NewDefaultUlcerIndex(klineItem *klines.Item) *UlcerIndex {
	return NewUlcerIndex(klineItem, 14)
}

// Calculation Func
func (e *UlcerIndex) Calculation() *UlcerIndex {

	var period = e.Period

	var closing = e.kline.GetOHLC().Close

	highClosing := ta.Max(period, closing)
	percentageDrawdown := ta.MultiplyBy(ta.Divide(ta.Subtract(closing, highClosing), highClosing), 100)
	squaredAverage := ta.Ema(period, ta.Multiply(percentageDrawdown, percentageDrawdown))

	ui := ta.Sqrt(squaredAverage)

	for i := 0; i < len(ui); i++ {
		e.data = append(e.data, UlcerIndexData{
			Time:  e.kline.Candles[i].Time,
			Value: ui[i],
		})
	}
	return e
}

// GetData Func
func (e *UlcerIndex) GetData() []UlcerIndexData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

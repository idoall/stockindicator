package channel

import (
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

// Keltner Channels是一个波动性指标，由一位名叫 Chester Keltner 的交易商在他 1960 年的著作《如何在商品中赚钱》中引入。
//
// Linda 版本的 Keltner Channel 使用更广泛，它与布林带非常相似，因为它也由三条线组成。
// 由于该通道源自ATR，而ATR本身就是一个波动率指标，因此Keltner 通道也会随着波动率收缩和扩张，但不像布林带那样波动。
//
// Middle Line = EMA(period, closings)
// Upper Band = EMA(period, closings) + 2 * ATR(period, highs, lows, closings)
// Lower Band = EMA(period, closings) - 2 * ATR(period, highs, lows, closings)
type KeltnerChannel struct {
	Name   string
	Period int
	data   []KeltnerChannelData
	kline  utils.Klines
}

// KeltnerChannelData
type KeltnerChannelData struct {
	Time   time.Time
	Upper  float64
	Middle float64
	Lower  float64
}

// NewKeltnerChannel new Func
func NewKeltnerChannel(list utils.Klines, period int) *KeltnerChannel {
	m := &KeltnerChannel{
		Name:   "KeltnerChannel",
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultKeltnerChannel new Func
func NewDefaultKeltnerChannel(list utils.Klines) *KeltnerChannel {
	return NewKeltnerChannel(list, 20)
}

// Calculation Func
func (e *KeltnerChannel) Calculation() *KeltnerChannel {

	var period = e.Period

	_, atr := trend.NewAtr(e.kline, period).GetValues()
	// _, atr := trend.Atr(period, high, low, closing)
	atr2 := utils.MultiplyBy(atr, 2)

	middleLine := trend.NewEma(e.kline, period).GetValues()
	upperBand := utils.Add(middleLine, atr2)
	lowerBand := utils.Subtract(middleLine, atr2)

	for i := 0; i < len(middleLine); i++ {
		e.data = append(e.data, KeltnerChannelData{
			Time:   e.kline[i].Time,
			Upper:  upperBand[i],
			Middle: middleLine[i],
			Lower:  lowerBand[i],
		})
	}
	return e
}

// AnalysisSide Func
// 当价格冲破上下轨道时，冲入上轨是就是可能的买的信号;
// 反之，冲破下轨时就是可能的卖的信号
//
// 更多使用方法：https://www.oanda.com/bvi-ft/lab-education/technical_analysis/average-candlestick-chart-keltner-channel/
func (e *KeltnerChannel) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]
		var prevPrice = e.kline[i-1].Close
		var price = e.kline[i].Close

		if price > v.Middle && prevPrice < prevItem.Middle {
			sides[i] = utils.Buy
		} else if price < v.Middle && prevPrice > prevItem.Middle {
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

// GetData Func
func (e *KeltnerChannel) GetData() []KeltnerChannelData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

package channel

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
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
	kline  *klines.Item
}

// KeltnerChannelData
type KeltnerChannelData struct {
	Time   time.Time
	Upper  float64
	Middle float64
	Lower  float64
}

// NewKeltnerChannel new Func
func NewKeltnerChannel(klineItem *klines.Item, period int) *KeltnerChannel {
	m := &KeltnerChannel{
		Name:   fmt.Sprintf("KeltnerChannel%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultKeltnerChannel new Func
func NewDefaultKeltnerChannel(klineItem *klines.Item) *KeltnerChannel {
	return NewKeltnerChannel(klineItem, 20)
}

// Calculation Func
func (e *KeltnerChannel) Calculation() *KeltnerChannel {

	var period = e.Period

	_, atr := trend.NewAtr(e.kline, period).GetValues()
	// _, atr := trend.Atr(period, high, low, closing)
	atr2 := ta.MultiplyBy(atr, 2)

	middleLine := trend.NewEma(e.kline, period).GetValues()
	upperBand := ta.Add(middleLine, atr2)
	lowerBand := ta.Subtract(middleLine, atr2)

	for i := 0; i < len(middleLine); i++ {
		e.data = append(e.data, KeltnerChannelData{
			Time:   time.Unix(e.kline.Candles[i].TimeUnix, 0),
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
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var ohlc = e.kline.GetOHLC()
	var closes = ohlc.Close

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]
		var price = closes[i]
		var prevPrice = closes[i-1]

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

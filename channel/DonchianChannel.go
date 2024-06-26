package channel

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Donchian 唐奇安通道指标，由“潮流之父”Richard Donchian 在二十世纪中叶开发。
// 它们是由移动平均线计算生成的三条线，包括由中间范围附近的上限和下限形成的指标。
// 但是，它们之间的区域是为该周期选择的通道。 这是为了帮助他识别市场交易的趋势。
// 因此，它在大多数交易平台上都很常见。
//
// Upper Channel = Mmax(period, closings)
// Lower Channel = Mmin(period, closings)
// Middle Channel = (Upper Channel + Lower Channel) / 2
type DonchianChannel struct {
	Name   string
	Period int
	data   []DonchianChannelData
	kline  *klines.Item
}

// DonchianChannelData
type DonchianChannelData struct {
	Time   time.Time
	Upper  float64
	Middle float64
	Lower  float64
}

// NewDonchianChannel new Func
func NewDonchianChannel(klineItem *klines.Item, period int) *DonchianChannel {
	m := &DonchianChannel{
		Name:   fmt.Sprintf("DonchianChannel%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultDonchianChannel new Func
func NewDefaultDonchianChannel(klineItem *klines.Item) *DonchianChannel {
	return NewDonchianChannel(klineItem, 20)
}

// Calculation Func
func (e *DonchianChannel) Calculation() *DonchianChannel {

	var period = e.Period

	var closing = e.kline.GetOHLC().Close

	upperChannel := ta.Max(period, closing)
	lowerChannel := ta.Min(period, closing)
	middleChannel := ta.DivideBy(ta.Add(upperChannel, lowerChannel), 2)

	for i := 0; i < len(middleChannel); i++ {
		e.data = append(e.data, DonchianChannelData{
			Time:   time.Unix(e.kline.Candles[i].TimeUnix, 0),
			Upper:  upperChannel[i],
			Middle: middleChannel[i],
			Lower:  lowerChannel[i],
		})
	}
	return e
}

// AnalysisSide Func
// 当价格冲破上下轨道时，冲入上轨是就是可能的买的信号;
// 反之，冲破下轨时就是可能的卖的信号
func (e *DonchianChannel) AnalysisSide() utils.SideData {
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
func (e *DonchianChannel) GetData() []DonchianChannelData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

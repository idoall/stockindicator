package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Williams R. 是由美国作家、股市投资家拉里·威廉斯（Larry R. Williams）在1973年出版的《我如何赚得一百万》一书中首先发表，这个指标是一个振荡指标，是依股价的摆动点来度量股票／指数是否处于超买或超卖的现象。
// 它衡量多空双方创出的峰值（最高价）距每天收市价的距离与一定时间内（如7天、14天、28天等）的股价波动范围的比例，以提供出股市趋势反转的讯号。
// 根据最低价、最高价和收盘价计算 Williams R。它是一种动量指标，在 0 和 -100 之间移动，衡量超买和超卖水平。
//
// WR = (Highest High - Closing) / (Highest High - Lowest Low) * -100.
//
// Buy when -80 and below. Sell when -20 and above.
type WilliamsR struct {
	Name   string
	Period int
	data   []WilliamsRData
	kline  utils.Klines
}

// WilliamsRData
type WilliamsRData struct {
	Time  time.Time
	Value float64
}

// NewWilliamsR new Func
func NewWilliamsR(list utils.Klines, period int) *WilliamsR {
	m := &WilliamsR{
		Name:   fmt.Sprintf("WilliamsR%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultWilliamsR new Func
func NewDefaultWilliamsR(list utils.Klines) *WilliamsR {
	return NewWilliamsR(list, 14)
}

// Calculation Func
func (e *WilliamsR) Calculation() *WilliamsR {

	var period = e.Period

	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low
	var closing = ohlc.Close

	highestHigh := utils.Max(period, high)
	lowestLow := utils.Min(period, low)

	result := make([]float64, len(closing))

	for i := 0; i < len(closing); i++ {
		result[i] = (highestHigh[i] - closing[i]) / (highestHigh[i] - lowestLow[i]) * float64(-100)
	}

	for i := 0; i < len(result); i++ {
		e.data = append(e.data, WilliamsRData{
			Time:  e.kline[i].Time,
			Value: result[i],
		})
	}
	return e
}

// AnalysisSide Func
// 1、当%R指标指标小于-80时，说明市场进入了超卖区域，如果此后%R指标从-80下方回头向上穿越-80时，可以参考买入信号。
// 2、当%R指标指标大于-20时，说明市场进入了超买区域，如果此后%R指标从-20上方回头向下穿过-20时，可以参考卖出信号。
func (e *WilliamsR) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]

		if v.Value < -80 && prevItem.Value > -80 {
			sides[i] = utils.Buy
		} else if v.Value > -20 && prevItem.Value < -20 {
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
func (e *WilliamsR) GetData() []WilliamsRData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	return e.data
}

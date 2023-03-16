package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

// Stochastic Oscillator. 显示了一定时期内收盘价相对于高低区间的位置。
// 随机震荡指标不跟随价格，也不跟随交易量或类似的东西。它跟随价格的速度或动量。通常，动量在价格之前改变方向。
// 因此，随机震荡指标的看涨和看跌背离可用于预示反转。
// 还可以使用这个振荡器来识别牛市和熊市的设置，以预测未来的逆转。
// 由于随机震荡指标是区间震荡指标（range-bound），因此它对于识别超买和超卖水平也很有用。
//
// K = (Closing - Lowest Low) / (Highest High - Lowest Low) * 100
// D = 3-Period SMA of K
type StochasticOscillator struct {
	Name   string
	Period int
	data   []StochasticOscillatorData
	kline  utils.Klines
}

// StochasticOscillatorData
type StochasticOscillatorData struct {
	Time time.Time
	K    float64
	D    float64
}

// NewStochasticOscillator new Func
func NewStochasticOscillator(list utils.Klines, period int) *StochasticOscillator {
	m := &StochasticOscillator{
		Name:   fmt.Sprintf("StochasticOscillator%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultStochasticOscillator new Func
func NewDefaultStochasticOscillator(list utils.Klines) *StochasticOscillator {
	return NewStochasticOscillator(list, 14)
}

// Calculation Func
func (e *StochasticOscillator) Calculation() *StochasticOscillator {

	var period = e.Period

	var high, low, closing []float64
	for _, v := range e.kline {
		high = append(high, v.High)
		low = append(low, v.Low)
		closing = append(closing, v.Close)
	}

	highestHigh14 := utils.Max(period, high)
	lowestLow14 := utils.Min(period, low)

	k := utils.MultiplyBy(utils.Divide(utils.Subtract(closing, lowestLow14), utils.Subtract(highestHigh14, lowestLow14)), float64(100))
	d := trend.NewSma(utils.CloseArrayToKline(k), 3).GetValues()

	for i := 0; i < len(k); i++ {
		e.data = append(e.data, StochasticOscillatorData{
			Time: e.kline[i].Time,
			K:    k[i],
			D:    d[i],
		})
	}
	return e
}

// AnalysisSide Func
// 当%D线超过了80 水平线时，就进去了超买区域，如果这时又从上下穿过80水平线时，这是一个可能的做空（卖出）信号。
// 当%D线低于了20水平线 时，就进入了超卖区域，如果这时又从下上穿过20水平线时，这是一个可能的做多（买进）信号。
//
// 当%K线从下向上穿过%D线时，预示着新的上升趋势。
// 当%K线从上向下穿过%D线时，预示着新的下降趋势。
//
// 当指标处在高位（超买区域），并形成依次向下的波峰，而此时价格形成依次向上的波峰，这叫顶背离，是很好的做空（卖出）信号。
// 当指标处在低位（超卖区域），并形成依次向上的波谷，而此时价格形成依次向下的波谷，这叫底背离，是很好的做多（买入）信号。
func (e *StochasticOscillator) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]

		if v.K > v.D && prevItem.K < prevItem.D {
			sides[i] = utils.Buy
		} else if v.K < v.D && prevItem.K > prevItem.D {
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
func (e *StochasticOscillator) GetData() []StochasticOscillatorData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

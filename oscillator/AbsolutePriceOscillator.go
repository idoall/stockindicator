package oscillator

import (
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

// AbsolutePriceOscillator struct
type AbsolutePriceOscillator struct {
	Name       string
	FastPeriod int // 默认一般是14
	SlowPeriod int // 默认一般是30
	data       []AbsolutePriceOscillatorData
	kline      utils.Klines
}

// AbsolutePriceOscillatorPoint 绝对价格震荡指标 (APO)
// AbsolutePriceOscillator函数计算用于跟踪趋势的技术指标。APO 上穿零表示看涨，而下穿零表示看跌。正值表示上升趋势，负值表示下降趋势。
type AbsolutePriceOscillatorData struct {
	Time  time.Time
	Value float64
}

// NewAbsolutePriceOscillator new Func
func NewAbsolutePriceOscillator(list utils.Klines) *AbsolutePriceOscillator {
	m := &AbsolutePriceOscillator{Name: "AbsolutePriceOscillator", kline: list, FastPeriod: 14, SlowPeriod: 30}
	return m
}

// Calculation Func
func (e *AbsolutePriceOscillator) Calculation() *AbsolutePriceOscillator {

	var closing []float64
	for _, v := range e.kline {
		closing = append(closing, v.Close)
	}

	fast := trend.NewEma(utils.CloseArrayToKline(closing), e.FastPeriod).GetValues()
	slow := trend.NewEma(utils.CloseArrayToKline(closing), e.SlowPeriod).GetValues()
	apo := utils.Subtract(fast, slow)

	for i := 0; i < len(apo); i++ {
		e.data = append(e.data, AbsolutePriceOscillatorData{
			Time:  e.kline[i].Time,
			Value: apo[i],
		})
	}
	return e
}

// AnalysisSide Func
func (e *AbsolutePriceOscillator) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		prevItem := e.data[i-1]
		// APO 上穿零表示看涨，而下穿零表示看跌。
		if v.Value > 0 && prevItem.Value < 0 {
			sides[i] = utils.Buy
		} else if v.Value < 0 && prevItem.Value > 0 {
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
func (e *AbsolutePriceOscillator) GetData() []AbsolutePriceOscillatorData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	return e.data
}

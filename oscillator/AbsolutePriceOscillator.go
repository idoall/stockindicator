package oscillator

import (
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// AbsolutePriceOscillator struct
type AbsolutePriceOscillator struct {
	Name       string
	FastPeriod int // 默认一般是14
	SlowPeriod int // 默认一般是30
	data       []AbsolutePriceOscillatorData
	kline      *klines.Item
}

// AbsolutePriceOscillatorPoint 绝对价格震荡指标 (APO)
// AbsolutePriceOscillator函数计算用于跟踪趋势的技术指标。APO 上穿零表示看涨，而下穿零表示看跌。正值表示上升趋势，负值表示下降趋势。
type AbsolutePriceOscillatorData struct {
	Time  time.Time
	Value float64
}

// NewAbsolutePriceOscillator new Func
func NewAbsolutePriceOscillator(klineItem *klines.Item) *AbsolutePriceOscillator {
	m := &AbsolutePriceOscillator{Name: "AbsolutePriceOscillator", kline: klineItem, FastPeriod: 14, SlowPeriod: 30}
	return m
}

// Calculation Func
func (e *AbsolutePriceOscillator) Calculation() *AbsolutePriceOscillator {

	var closing = e.kline.GetOHLC().Close

	fast := ta.Ema(e.FastPeriod, closing)
	slow := ta.Ema(e.SlowPeriod, closing)
	apo := ta.Subtract(fast, slow)

	var data = make([]AbsolutePriceOscillatorData, len(apo))
	for i := 0; i < len(apo); i++ {
		data[i] = AbsolutePriceOscillatorData{
			Time:  time.Unix(e.kline.Candles[i].TimeUnix, 0),
			Value: apo[i],
		}
	}
	e.data = data
	return e
}

// AnalysisSide Func
func (e *AbsolutePriceOscillator) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

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

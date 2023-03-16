package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/volume"
)

// ChaikinOscillator 动量震荡（Awesome Oscillator）是一个用于很累市场动量的指标。
// AO（Awesome Oscillator）计算34个周期和5个周期简单移动平均的差。
// 使用的简单移动平均不是使用收盘价计算的，而是每个柱的中点价格。AO通常被用来确认趋势或预期可能的逆转。
//
// Median Price = ((Low + High) / 2).
// AO = 5-Period SMA - 34-Period SMA.
type ChaikinOscillator struct {
	Name       string
	FastPeriod int
	SlowPeriod int
	data       []ChaikinOscillatorData
	kline      utils.Klines
}

// ChaikinOscillatorData.
type ChaikinOscillatorData struct {
	Time  time.Time
	Value float64
	AD    float64
}

// NewChaikinOscillator new Func
func NewChaikinOscillator(list utils.Klines, fastPeriod, slowPeriod int) *ChaikinOscillator {
	m := &ChaikinOscillator{
		Name:       fmt.Sprintf("ChaikinOscillator%d-%d", fastPeriod, slowPeriod),
		kline:      list,
		FastPeriod: fastPeriod,
		SlowPeriod: slowPeriod,
	}
	return m
}

// NewDefaultChaikinOscillator new Func
func NewDefaultChaikinOscillator(list utils.Klines) *ChaikinOscillator {
	return NewChaikinOscillator(list, 3, 10)
}

// Calculation Func
func (e *ChaikinOscillator) Calculation() *ChaikinOscillator {

	ad := volume.NewAccumulationDistribution(e.kline).GetValues()
	fastEMA := trend.NewEma(utils.CloseArrayToKline(ad), e.FastPeriod).GetValues()
	slowEMA := trend.NewEma(utils.CloseArrayToKline(ad), e.SlowPeriod).GetValues()
	co := utils.Subtract(fastEMA, slowEMA)

	// var co, ad = trend.ChaikinOscillator(e.FastPeriod, e.SlowPeriod, low, high, closing, volume)

	for i := 0; i < len(co); i++ {
		e.data = append(e.data, ChaikinOscillatorData{
			Time:  e.kline[i].Time,
			Value: co[i],
			AD:    ad[i],
		})
	}
	return e
}

// AnalysisSide Func
func (e *ChaikinOscillator) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		prevItem := e.data[i-1]
		// 穿越到零線上方時，會出現看漲金叉。
		if v.Value > 0 && prevItem.Value < 0 {
			sides[i] = utils.Buy
		} else if v.Value < 0 && prevItem.Value > 0 {
			sides[i] = utils.Sell
		} else {
			sides[i] = utils.Hold
		}

		// prevItem := e.data[i-1]
		// // 越过 A/D 线表示看涨。
		// if v.Value > v.AD && prevItem.Value < v.AD {
		// 	sides[i] = utils.Buy
		// } else if v.Value < v.AD && prevItem.Value > v.AD {
		// 	sides[i] = utils.Sell
		// } else {
		// 	sides[i] = utils.Hold
		// }
	}
	return utils.SideData{
		Name: e.Name,
		Data: sides,
	}
}

// GetData Func
func (e *ChaikinOscillator) GetData() []ChaikinOscillatorData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

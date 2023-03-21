package oscillator

import (
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

// AwesomeOscillator struct
type AwesomeOscillator struct {
	Name  string
	data  []AwesomeOscillatorData
	kline utils.Klines
}

// Awesome Oscillator.动量震荡（Awesome Oscillator）是一个用于很累市场动量的指标。
// AO（Awesome Oscillator）计算34个周期和5个周期简单移动平均的差。
// 使用的简单移动平均不是使用收盘价计算的，而是每个柱的中点价格。AO通常被用来确认趋势或预期可能的逆转。
//
// Median Price = ((Low + High) / 2).
// AO = 5-Period SMA - 34-Period SMA.
type AwesomeOscillatorData struct {
	Time  time.Time
	Value float64
}

// NewAwesomeOscillator new Func
func NewAwesomeOscillator(list utils.Klines) *AwesomeOscillator {
	m := &AwesomeOscillator{Name: "AwesomeOscillator", kline: list}
	return m
}

// Calculation Func
func (e *AwesomeOscillator) Calculation() *AwesomeOscillator {

	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low

	medianPrice := utils.DivideBy(utils.Add(low, high), float64(2))
	sma5 := trend.NewEma(utils.CloseArrayToKline(medianPrice), 5).GetValues()
	sma34 := trend.NewEma(utils.CloseArrayToKline(medianPrice), 34).GetValues()
	ao := utils.Subtract(sma5, sma34)

	for i := 0; i < len(ao); i++ {
		e.data = append(e.data, AwesomeOscillatorData{
			Time:  e.kline[i].Time,
			Value: ao[i],
		})
	}
	return e
}

// AnalysisSide Func
func (e *AwesomeOscillator) AnalysisSide() utils.SideData {
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

// GetPoints Func
func (e *AwesomeOscillator) GetData() []AwesomeOscillatorData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

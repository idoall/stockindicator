package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

// ProjectionOscillator Percentage Price Oscillator (PPO). 由Dr. Mel Widner 研仓。
// 与其他不同的指标一样，传统的用法也是超买/超卖，背驰，突破等，不少人利用此指标来交易外汇。
// 传统的参数是12及5，但若应用在港股上，将参数设定为50及10会更好。
// 分析股票时，初步看，每当由50以下重回至50以上有机会是股价重拾升势的时间，值得留意，不过有关的方法仍有待详细测试。
type ProjectionOscillator struct {
	Name   string
	Period int //默认一般是50
	Smooth int // 默认一般是1
	data   []ProjectionOscillatorData
	kline  utils.Klines
}

// ProjectionOscillatorData 投影振荡器策略
// 在po高于spo时提供买入操作，在po低于spo时提供卖出操作。
type ProjectionOscillatorData struct {
	Time time.Time
	Po   float64
	Spo  float64
}

// NewProjectionOscillator new Func
func NewProjectionOscillator(list utils.Klines, period, smooth int) *ProjectionOscillator {
	m := &ProjectionOscillator{
		Name:   fmt.Sprintf("ProjectionOscillator%d-%d", period, smooth),
		kline:  list,
		Period: period,
		Smooth: smooth,
	}
	return m
}

// NewProjectionOscillator new Func
func NewDefaultProjectionOscillator(list utils.Klines) *ProjectionOscillator {
	return NewProjectionOscillator(list, 13, 3)
}

// Calculation Func
func (e *ProjectionOscillator) Calculation() *ProjectionOscillator {

	period := e.Period
	smooth := e.Smooth
	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low
	var closing = ohlc.Close

	x := utils.GenerateNumbers(0, float64(len(closing)), 1)
	mHigh, _ := utils.MovingLeastSquare(period, x, high)
	mLow, _ := utils.MovingLeastSquare(period, x, low)

	vHigh := utils.Add(high, utils.Multiply(mHigh, x))
	vLow := utils.Add(low, utils.Multiply(mLow, x))

	pu := utils.Max(period, vHigh)
	pl := utils.Min(period, vLow)

	po := utils.Divide(utils.MultiplyBy(utils.Subtract(closing, pl), 100), utils.Subtract(pu, pl))
	spo := trend.NewEma(utils.CloseArrayToKline(po), smooth).GetValues()
	for i := 0; i < len(po); i++ {
		e.data = append(e.data, ProjectionOscillatorData{
			Time: e.kline[i].Time,
			Po:   po[i],
			Spo:  spo[i],
		})
	}
	return e
}

// AnalysisSide Func
func (e *ProjectionOscillator) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		// APO 上穿零表示看涨，而下穿零表示看跌。
		if v.Po > v.Spo {
			sides[i] = utils.Buy
		} else if v.Po < v.Spo {
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
func (e *ProjectionOscillator) GetData() []ProjectionOscillatorData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

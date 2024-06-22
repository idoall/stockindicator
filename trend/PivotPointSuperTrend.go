package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/commonutils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// PivotPointSuperTrend struct
type PivotPointSuperTrend struct {
	// Atr Period
	AtrPeriod int
	// Atr Factor
	AtrFactor float64
	// Pivot Point Period
	Period int
	Name   string
	data   []PivotPointSuperTrendData
	ohlc   *klines.OHLC
}

type PivotPointSuperTrendData struct {
	Time           time.Time
	UpTrend        float64
	UpTrendBegin   float64
	DownTrend      float64
	DownTrendBegin float64
}

// NewPivotPointSuperTrend new Func
func NewPivotPointSuperTrend(klineItem *klines.Item, period, atrPeriod int, atrFactor float64) *PivotPointSuperTrend {
	m := &PivotPointSuperTrend{
		Name:      fmt.Sprintf("PivotPointSuperTrend%d-%d-%f", period, atrPeriod, atrFactor),
		AtrPeriod: atrPeriod,
		AtrFactor: atrFactor,
		Period:    period,
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

// NewPivotPointSuperTrend new Func
func NewPivotPointSuperTrendOHLC(ohlc *klines.OHLC, period, atrPeriod int, atrFactor float64) *PivotPointSuperTrend {
	m := &PivotPointSuperTrend{
		Name:      fmt.Sprintf("PivotPointSuperTrend%d-%d-%f", period, atrPeriod, atrFactor),
		ohlc:      ohlc,
		AtrPeriod: atrPeriod,
		AtrFactor: atrFactor,
		Period:    period,
	}
	return m
}

// NewPivotPointSuperTrend new Func
func NewDefaultPivotPointSuperTrend(klineItem *klines.Item) *PivotPointSuperTrend {
	return NewPivotPointSuperTrend(klineItem, 3, 8, 2.8)
}

// Calculation Func
func (e *PivotPointSuperTrend) Calculation() *PivotPointSuperTrend {

	var closes = e.ohlc.Close
	var highs = e.ohlc.High
	var lows = e.ohlc.Low
	var times = e.ohlc.TimeUnix

	var ph = ta.PivotHigh(highs, e.Period, e.Period)
	var pl = ta.PivotLow(lows, e.Period, e.Period)
	var atr = ta.Atr(highs, lows, closes, e.AtrPeriod)

	var tUP = make([]float64, len(closes))
	var tDown = make([]float64, len(closes))
	var tCenter = make([]float64, len(closes))
	var trend = make([]int, len(closes))

	defer func() {
		tUP = nil
		tDown = nil
		ph = nil
		pl = nil
		atr = nil
		trend = nil
		tCenter = nil
	}()

	e.data = make([]PivotPointSuperTrendData, len(closes))

	for i := 1; i < len(closes); i++ {

		var c = closes[i]

		// var center float64
		var lastpp = commonutils.If(ph[i] > 0, ph[i], commonutils.If(pl[i] > 0, pl[i], 0.0).(float64)).(float64)

		if lastpp != 0.0 {
			if tCenter[i-1] == 0.0 {
				tCenter[i] = lastpp
			} else {
				tCenter[i] = (tCenter[i-1]*2 + lastpp) / 3
			}
		} else {
			tCenter[i] = tCenter[i-1]
		}

		var closePrev = closes[i-1]
		var Up = tCenter[i] - (e.AtrFactor * atr[i])
		var Dn = tCenter[i] + (e.AtrFactor * atr[i])

		tUP[i] = commonutils.If(closePrev > tUP[i-1], math.Max(Up, tUP[i-1]), Up).(float64)
		tDown[i] = commonutils.If(closePrev < tDown[i-1], math.Min(Dn, tDown[i-1]), Dn).(float64)
		trend[i] = commonutils.If(c > tDown[i-1], 1, commonutils.If(c < tUP[i-1], -1, trend[i-1]).(int)).(int)

		var trailingsl = commonutils.If(trend[i] == 1, tUP[i], tDown[i]).(float64)

		e.data[i] = PivotPointSuperTrendData{
			Time:           time.Unix(times[i], 0),
			UpTrend:        commonutils.If(trend[i] == 1 && trend[i-1] == 1, trailingsl, 0.0).(float64),
			UpTrendBegin:   commonutils.If(trend[i] == 1 && trend[i-1] == -1, trailingsl, 0.0).(float64),
			DownTrend:      commonutils.If(trend[i] == -1 && trend[i-1] == -1, trailingsl, 0.0).(float64),
			DownTrendBegin: commonutils.If(trend[i] == -1 && trend[i-1] == 1, trailingsl, 0.0).(float64),
		}
	}

	return e
}

// GetValues Func
func (e *PivotPointSuperTrend) GetValues() (up []float64, down []float64) {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	for _, v := range e.data {
		up = append(up, v.UpTrend)
		down = append(down, v.DownTrend)
	}
	return
}

// GetData Func
func (e *PivotPointSuperTrend) GetData() []PivotPointSuperTrendData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// AnalysisSide Func
func (e *PivotPointSuperTrend) AnalysisSide() utils.SideData {

	if len(e.data) == 0 {
		e = e.Calculation()
	}
	sides := make([]utils.Side, len(e.data))

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		if v.UpTrendBegin != 0 {
			sides[i] = utils.Buy
		} else if v.DownTrendBegin != 0 {
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

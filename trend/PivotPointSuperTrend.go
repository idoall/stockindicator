package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/commonutils"
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
	kline  utils.Klines
}

type PivotPointSuperTrendData struct {
	Time           time.Time
	UpTrend        float64
	UpTrendBegin   float64
	DownTrend      float64
	DownTrendBegin float64
}

// NewPivotPointSuperTrend new Func
func NewPivotPointSuperTrend(list utils.Klines, period, atrPeriod int, atrFactor float64) *PivotPointSuperTrend {
	m := &PivotPointSuperTrend{
		Name:      fmt.Sprintf("PivotPointSuperTrend%d-%d-%f", period, atrPeriod, atrFactor),
		kline:     list,
		AtrPeriod: atrPeriod,
		AtrFactor: atrFactor,
		Period:    period,
	}
	return m
}

// NewPivotPointSuperTrend new Func
func NewDefaultPivotPointSuperTrend(list utils.Klines) *PivotPointSuperTrend {
	return NewPivotPointSuperTrend(list, 3, 8, 2.8)
}

// Calculation Func
func (e *PivotPointSuperTrend) Calculation() *PivotPointSuperTrend {

	var close = e.kline.GetOHLC().Close

	var ph = utils.PivotHigh(close, e.Period, e.Period)
	var pl = utils.PivotLow(close, e.Period, e.Period)
	var _, atr = NewAtr(e.kline, e.AtrPeriod).GetValues()

	var tUP = make([]float64, len(e.kline))
	var tDown = make([]float64, len(e.kline))
	var tCenter = make([]float64, len(e.kline))
	var trend = make([]int, len(e.kline))

	defer func() {
		tUP = nil
		tDown = nil
		ph = nil
		pl = nil
		atr = nil
		trend = nil
		tCenter = nil
		close = nil
	}()

	e.data = make([]PivotPointSuperTrendData, len(e.kline))

	for i := 1; i < len(e.kline); i++ {

		var c = close[i]
		var time = e.kline[i].Time

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

		var closePrev = close[i-1]
		var Up = tCenter[i] - (e.AtrFactor * atr[i])
		var Dn = tCenter[i] + (e.AtrFactor * atr[i])

		tUP[i] = commonutils.If(closePrev > tUP[i], math.Max(Up, tUP[1]), Up).(float64)
		tDown[i] = commonutils.If(closePrev < tDown[i], math.Max(Dn, tDown[1]), Dn).(float64)
		trend[i] = commonutils.If(c > tDown[i], 1, commonutils.If(c < tUP[i], -1, trend[i-1]).(int)).(int)

		var trailingsl = commonutils.If(trend[i] == 1, tUP[i], tDown[i]).(float64)

		e.data[i] = PivotPointSuperTrendData{
			Time:           time,
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
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

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

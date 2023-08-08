package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/ta"
)

// SuperTrend 超级趋势 struct
type SuperTrend struct {
	// ATR参数
	AtrPeriod int
	// ATR倍数
	AtrMultiplier int
	// 是否使用ATR算法
	ChangeAtr bool
	Name      string
	data      []SuperTrendData
	kline     utils.Klines
}

type SuperTrendData struct {
	Time           time.Time
	UpTrend        float64
	UpTrendBegin   float64
	DownTrend      float64
	DownTrendBegin float64
}

// NewSuperTrend new Func
func NewSuperTrend(list utils.Klines, atrPeriod, atrMultiplier int, changeAtr bool) *SuperTrend {
	m := &SuperTrend{
		Name:          fmt.Sprintf("SuperTrend%d-%d", atrPeriod, atrMultiplier),
		kline:         list,
		AtrPeriod:     atrPeriod,
		AtrMultiplier: atrMultiplier,
		ChangeAtr:     changeAtr,
	}
	return m
}

// NewSuperTrend new Func
func NewDefaultSuperTrend(list utils.Klines) *SuperTrend {
	return NewSuperTrend(list, 10, 3, true)
}

// Calculation Func
func (e *SuperTrend) Calculation() *SuperTrend {

	var up = make([]float64, len(e.kline))
	var upb = make([]float64, len(e.kline))
	var up1 = make([]float64, len(e.kline))
	var dn = make([]float64, len(e.kline))
	var dnb = make([]float64, len(e.kline))
	var dn1 = make([]float64, len(e.kline))
	var trend = make([]float64, len(e.kline))
	var atr = make([]float64, len(e.kline))

	var src = e.kline.HL2()
	var close = e.kline.GetOHLC().Close

	defer func() {
		up = nil
		upb = nil
		up1 = nil
		dn = nil
		dnb = nil
		dn1 = nil
		trend = nil
		atr = nil
		src = nil
		close = nil
	}()
	// 是否使用原ATR算法
	if e.ChangeAtr {
		_, atr = NewAtr(e.kline, e.AtrPeriod).GetValues()

	} else {
		var tr, _ = NewAtr(e.kline, e.AtrPeriod).GetValues()
		atr = ta.Sma(e.AtrPeriod, tr)
	}

	for i := 0; i < len(src); i++ {

		if i < 1 {
			trend[i] = 1
			continue
		}

		up[i] = src[i] - (float64(e.AtrMultiplier) * atr[i])
		up1[i] = ta.Nz(up[i-1], up[i])
		if close[i-1] > up1[i] {
			up[i] = math.Max(up[i], up1[i])
		}

		dn[i] = src[i] + (float64(e.AtrMultiplier) * atr[i])
		dn1[i] = ta.Nz(dn[i-1], dn[i])
		if close[i-1] < dn1[i] {
			dn[i] = math.Min(dn[i], dn1[i])
		}

		trend[i] = 1
		trend[i] = ta.Nz(trend[i-1], trend[i])

		if trend[i] == -1 && close[i] > dn1[i] {
			trend[i] = 1
		} else if trend[i] == 1 && close[i] < up1[i] {
			trend[i] = -1
		}

		if trend[i] == 1 && trend[i-1] == -1 {
			upb[i] = up[i]
		} else if trend[i] == -1 && trend[i-1] == 1 {
			dnb[i] = dn[i]
		}
	}

	e.data = make([]SuperTrendData, len(e.kline))
	for i, v := range e.kline {
		e.data[i] = SuperTrendData{
			Time:           v.Time,
			UpTrend:        up[i],
			UpTrendBegin:   upb[i],
			DownTrend:      dn[i],
			DownTrendBegin: dnb[i],
		}
	}

	return e
}

// GetValues Func
func (e *SuperTrend) GetValues() (up []float64, down []float64) {
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
func (e *SuperTrend) GetData() []SuperTrendData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// AnalysisSide Func
func (e *SuperTrend) AnalysisSide() utils.SideData {
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

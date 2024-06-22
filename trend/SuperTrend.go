package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
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
	ohlc      *klines.OHLC
	hl2       []float64
}

type SuperTrendData struct {
	Time           time.Time
	UpTrend        float64
	UpTrendBegin   float64
	DownTrend      float64
	DownTrendBegin float64
}

// NewSuperTrend new Func
func NewSuperTrend(klineItem *klines.Item, atrPeriod, atrMultiplier int, changeAtr bool) *SuperTrend {
	m := &SuperTrend{
		Name:          fmt.Sprintf("SuperTrend%d-%d", atrPeriod, atrMultiplier),
		AtrPeriod:     atrPeriod,
		AtrMultiplier: atrMultiplier,
		ChangeAtr:     changeAtr,
	}
	m.ohlc = klineItem.GetOHLC()
	m.hl2 = klineItem.HL2()
	return m
}

func NewSuperTrendOHLC(ohlc *klines.OHLC, hl2 []float64, atrPeriod, atrMultiplier int, changeAtr bool) *SuperTrend {
	m := &SuperTrend{
		Name:          fmt.Sprintf("SuperTrend%d-%d", atrPeriod, atrMultiplier),
		ohlc:          ohlc,
		hl2:           hl2,
		AtrPeriod:     atrPeriod,
		AtrMultiplier: atrMultiplier,
		ChangeAtr:     changeAtr,
	}
	return m
}

// NewSuperTrend new Func
func NewDefaultSuperTrend(klineItem *klines.Item) *SuperTrend {
	return NewSuperTrend(klineItem, 10, 3, true)
}

// Calculation Func
func (e *SuperTrend) Calculation() *SuperTrend {

	var closes = e.ohlc.Close
	var highs = e.ohlc.High
	var lows = e.ohlc.Low
	var times = e.ohlc.Time

	var up = make([]float64, len(closes))
	var upb = make([]float64, len(closes))
	var up1 = make([]float64, len(closes))
	var dn = make([]float64, len(closes))
	var dnb = make([]float64, len(closes))
	var dn1 = make([]float64, len(closes))
	var trend = make([]float64, len(closes))
	var atr = make([]float64, len(closes))

	var src = e.hl2

	defer func() {
		up = nil
		upb = nil
		up1 = nil
		dn = nil
		dnb = nil
		dn1 = nil
		trend = nil
		atr = nil
	}()
	// 是否使用原ATR算法
	if e.ChangeAtr {
		atr = ta.Atr(highs, lows, closes, e.AtrPeriod)
	} else {
		natr := ta.Natr(highs, lows, closes, e.AtrPeriod)
		atr = ta.Sma(e.AtrPeriod, natr)
	}

	for i := 0; i < len(src); i++ {

		if i < 1 {
			trend[i] = 1
			continue
		}

		up[i] = src[i] - (float64(e.AtrMultiplier) * atr[i])
		up1[i] = ta.Nz(up[i-1], up[i])
		if closes[i-1] > up1[i] {
			up[i] = math.Max(up[i], up1[i])
		}

		dn[i] = src[i] + (float64(e.AtrMultiplier) * atr[i])
		dn1[i] = ta.Nz(dn[i-1], dn[i])
		if closes[i-1] < dn1[i] {
			dn[i] = math.Min(dn[i], dn1[i])
		}

		trend[i] = 1
		trend[i] = ta.Nz(trend[i-1], trend[i])

		if trend[i] == -1 && closes[i] > dn1[i] {
			trend[i] = 1
		} else if trend[i] == 1 && closes[i] < up1[i] {
			trend[i] = -1
		}

		if trend[i] == 1 && trend[i-1] == -1 {
			upb[i] = up[i]
		} else if trend[i] == -1 && trend[i-1] == 1 {
			dnb[i] = dn[i]
		}
	}

	e.data = make([]SuperTrendData, len(closes))
	for i := 0; i < len(closes); i++ {
		e.data[i] = SuperTrendData{
			Time:           times[i],
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

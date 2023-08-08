package trend

import (
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/types"
)

// ReversalSignals is the main object
type ReversalSignals struct {
	Name string
	// 增量
	Increment  int
	SmoothFast int
	SmoothSlow int
	PlotNum    int
	MATypes    types.MATypes
	data       []ReversalSignalsData
	kline      utils.Klines
}

type ReversalSignalsData struct {
	Side utils.Side
	Time time.Time
}

// NewReversalSignals new Func
//
//	Args:
//		list K线
func NewReversalSignals(list utils.Klines) *ReversalSignals {
	m := &ReversalSignals{
		Name:  "ReversalSignals",
		kline: list,
	}
	return m
}

// Calculation Func
func (e *ReversalSignals) Calculation() *ReversalSignals {

	var ohlc = e.kline.GetOHLC()
	var lows = ohlc.Low
	var highs = ohlc.High
	var closes = ohlc.Close

	defer func() {
		ohlc = nil
		lows = nil
		closes = nil
		highs = nil
	}()

	var bullishCount, bearishCount int

	e.data = make([]ReversalSignalsData, len(e.kline))

	for i := range e.kline {
		var close = closes[i]
		var high = highs[i]
		var low = lows[i]
		var time = e.kline[i].Time

		e.data[i].Side = utils.Hold
		e.data[i].Time = time

		if i < 4 {
			continue
		}

		var con = close < closes[i-4]

		if con {
			bullishCount = bullishCount + 1
			bearishCount = 0
		} else {
			bearishCount = bearishCount + 1
			bullishCount = 0
		}

		var pbs bool
		if (low <= lows[i-3] && low <= lows[i-2]) || (lows[i-1] < lows[i-3] && lows[i-1] <= lows[i-2]) {
			pbs = true
		}
		if bullishCount == 9 && pbs {
			e.data[i].Side = utils.Buy
		}

		var pss bool
		if (high >= highs[i-3] && high >= highs[i-2]) || (highs[i-1] >= highs[i-3] && highs[i-1] >= highs[i-2]) {
			pss = true
		}
		if bearishCount == 9 && pss {
			e.data[i].Side = utils.Sell
		}

	}

	return e
}

// GetPoints return Point
func (e *ReversalSignals) GetData() []ReversalSignalsData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

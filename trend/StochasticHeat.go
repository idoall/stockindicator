package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/types"
)

// StochasticHeat is the main object
type StochasticHeat struct {
	Name string
	// 增量
	Increment  int
	SmoothFast int
	SmoothSlow int
	PlotNum    int
	MATypes    types.MATypes
	data       []StochasticHeatData
	kline      utils.Klines
}

type StochasticHeatData struct {
	Fast float64
	Slow float64
	Time time.Time
}

// NewStochasticHeat new Func
//
//	Args:
//		list K线
//		inc	增量
//		smoothFast Smooth Fast
//		smoothSlow Smooth Slow
//		plotNum 1-28
//		maTypes Only "SMA","EMA","WMA"
func NewStochasticHeat(list utils.Klines, inc, smoothFast, smoothSlow, plotNum int, maTypes types.MATypes) *StochasticHeat {
	m := &StochasticHeat{
		Name:       fmt.Sprintf("StochasticHeat%d-%d-%d-%d-%d", inc, smoothFast, smoothSlow, plotNum, maTypes),
		kline:      list,
		Increment:  inc,
		SmoothFast: smoothFast,
		SmoothSlow: smoothSlow,
		PlotNum:    plotNum,
		MATypes:    maTypes,
	}
	return m
}

func NewDefaultStochasticHeat(list utils.Klines) *StochasticHeat {
	return NewStochasticHeat(list, 8, 7, 26, 25, types.WMA)
}

// Calculation Func
func (e *StochasticHeat) Calculation() *StochasticHeat {

	var ohlc = e.kline.GetOHLC()
	var highs = ohlc.High
	var lows = ohlc.Low
	var closing = ohlc.Close

	defer func() {
		ohlc = nil
		highs = nil
		lows = nil
		closing = nil
	}()

	var stoch1 = e.getStoch(1, closing, highs, lows)
	var stoch2 = e.getStoch(2, closing, highs, lows)
	var stoch3 = e.getStoch(3, closing, highs, lows)
	var stoch4 = e.getStoch(4, closing, highs, lows)
	var stoch5 = e.getStoch(5, closing, highs, lows)
	var stoch6 = e.getStoch(6, closing, highs, lows)
	var stoch7 = e.getStoch(7, closing, highs, lows)
	var stoch8 = e.getStoch(8, closing, highs, lows)
	var stoch9 = e.getStoch(9, closing, highs, lows)
	var stoch10 = e.getStoch(10, closing, highs, lows)
	var stoch11 = e.getStoch(11, closing, highs, lows)
	var stoch12 = e.getStoch(12, closing, highs, lows)
	var stoch13 = e.getStoch(13, closing, highs, lows)
	var stoch14 = e.getStoch(14, closing, highs, lows)
	var stoch15 = e.getStoch(15, closing, highs, lows)
	var stoch16 = e.getStoch(16, closing, highs, lows)
	var stoch17 = e.getStoch(17, closing, highs, lows)
	var stoch18 = e.getStoch(18, closing, highs, lows)
	var stoch19 = e.getStoch(19, closing, highs, lows)
	var stoch20 = e.getStoch(20, closing, highs, lows)
	var stoch21 = e.getStoch(21, closing, highs, lows)
	var stoch22 = e.getStoch(22, closing, highs, lows)
	var stoch23 = e.getStoch(23, closing, highs, lows)
	var stoch24 = e.getStoch(24, closing, highs, lows)
	var stoch25 = e.getStoch(25, closing, highs, lows)
	var stoch26 = e.getStoch(26, closing, highs, lows)
	var stoch27 = e.getStoch(27, closing, highs, lows)
	var stoch28 = e.getStoch(28, closing, highs, lows)

	var fast = make([]float64, len(e.kline))
	var slow = make([]float64, len(e.kline))

	for i := range e.kline {
		var getAverage = (stoch1[i] + stoch2[i] + stoch3[i] + stoch4[i] + stoch5[i] + stoch6[i] + stoch7[i] + stoch8[i] + stoch9[i] + stoch10[i] + stoch11[i] + stoch12[i] + stoch13[i] + stoch14[i] + stoch15[i] + stoch16[i] + stoch17[i] + stoch18[i] + stoch19[i] + stoch20[i] + stoch21[i] + stoch22[i] + stoch23[i] + stoch24[i] + stoch25[i] + stoch26[i] + stoch27[i] + stoch28[i]) / float64(e.PlotNum)
		fast[i] = ((getAverage / 100) * float64(e.PlotNum))

		if e.MATypes == types.EMA {
			slow = utils.Ema(e.SmoothSlow, fast)
		} else if e.MATypes == types.SMA {
			slow = utils.Sma(e.SmoothSlow, fast)
		} else {
			slow = utils.Wma(e.SmoothSlow, fast)
		}
	}
	for i := 0; i < len(e.kline); i++ {
		e.data = append(e.data, StochasticHeatData{
			Time: e.kline[i].Time,
			Fast: fast[i],
			Slow: slow[i],
		})
	}

	defer func() {
		ohlc = nil
		highs = nil
		lows = nil
		closing = nil
		fast = nil
		slow = nil
		stoch1 = nil
		stoch2 = nil
		stoch3 = nil
		stoch4 = nil
		stoch5 = nil
		stoch6 = nil
		stoch7 = nil
		stoch8 = nil
		stoch9 = nil
		stoch10 = nil
		stoch11 = nil
		stoch12 = nil
		stoch13 = nil
		stoch14 = nil
		stoch15 = nil
		stoch16 = nil
		stoch17 = nil
		stoch18 = nil
		stoch19 = nil
		stoch20 = nil
		stoch21 = nil
		stoch22 = nil
		stoch23 = nil
		stoch24 = nil
		stoch25 = nil
		stoch26 = nil
		stoch27 = nil
		stoch28 = nil
	}()

	return e
}

// getStoch Func
func (e *StochasticHeat) getStoch(i int, closing, highs, lows []float64) []float64 {

	var c = i * e.Increment
	var k, _ = utils.Stochastic(closing, highs, lows, c)
	switch e.MATypes {
	case types.EMA:
		return utils.Ema(e.SmoothFast, k)
	case types.SMA:
		return utils.Sma(e.SmoothFast, k)
	case types.WMA:
		return utils.Wma(e.SmoothFast, k)
	default:
		return k
	}
}

// AnalysisSide Func
func (e *StochasticHeat) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]

		if prevItem.Slow < prevItem.Fast && v.Slow > v.Fast {
			sides[i] = utils.Buy
		} else if prevItem.Slow > prevItem.Fast && v.Slow < v.Fast {
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

// GetPoints return Point
func (e *StochasticHeat) GetData() []StochasticHeatData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

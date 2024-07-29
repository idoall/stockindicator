package klines

import (
	"math"
)

// OHLC is a connector for technical analysis usage
type OHLC struct {
	Open       []float64
	High       []float64
	Low        []float64
	Close      []float64
	Volume     []float64
	BullMarket []bool
	TimeUnix   []int64
}

// GetOHLC returns the entire subset of candles as a friendly type for gct
// technical analysis usage.
func (e *Item) GetOHLC() *OHLC {
	ohlc := &OHLC{
		Open:       make([]float64, len(e.Candles)),
		High:       make([]float64, len(e.Candles)),
		Low:        make([]float64, len(e.Candles)),
		Close:      make([]float64, len(e.Candles)),
		Volume:     make([]float64, len(e.Candles)),
		BullMarket: make([]bool, len(e.Candles)),
		TimeUnix:   make([]int64, len(e.Candles)),
	}
	for x := range e.Candles {
		ohlc.Open[x] = e.Candles[x].Open
		ohlc.High[x] = e.Candles[x].High
		ohlc.Low[x] = e.Candles[x].Low
		ohlc.Close[x] = e.Candles[x].Close
		ohlc.Volume[x] = e.Candles[x].Volume
		ohlc.BullMarket[x] = e.Candles[x].IsBullMarket
		ohlc.TimeUnix[x] = e.Candles[x].TimeUnix
	}
	return ohlc
}

// ToHeikinAshi 转换成平均K线
func (e *Item) ToHeikinAshi() *Item {
	var result = &Item{
		Exchange: e.Exchange,
		Interval: e.Interval,
	}

	result.Candles = make([]*Candle, len(e.Candles))
	for index, candle := range e.Candles {

		var open = (candle.Open + candle.Close) / 2
		if index > 1 {
			var prev = e.Candles[index-1]
			open = (prev.Open + prev.Close) / 2
		}

		var heikinAshiCandle = &Candle{
			Open:     open,
			Close:    (candle.Open + candle.High + candle.Low + candle.Close) / 4,
			High:     math.Max(candle.High, math.Max(open, candle.Close)),
			Low:      math.Max(candle.Low, math.Min(open, candle.Close)),
			TimeUnix: candle.TimeUnix,
		}

		heikinAshiCandle.ChangePercent = (candle.Close - candle.Open) / candle.Open
		if candle.Close > candle.Open {
			heikinAshiCandle.IsBullMarket = true
		}
		result.Candles[index] = heikinAshiCandle
	}

	return result
}

// HL2 (最高价+最低价)/2
func (e *Item) HL2() []float64 {
	result := make([]float64, len(e.Candles))
	for i, candle := range e.Candles {
		result[i] = (candle.High + candle.Low) / 2
	}
	return result
}

// HLC3 (最高价+最低价+收盘价)/3
func (e *Item) HLC3() []float64 {

	result := make([]float64, len(e.Candles))
	for i, candle := range e.Candles {
		result[i] = (candle.High + candle.Low + candle.Close) / 3
	}
	return result
}

// OHLC4 (开盘价 + 最高价 + 最低价 + 收盘价)/4
func (e *Item) OHLC4() []float64 {

	result := make([]float64, len(e.Candles))
	for i, candle := range e.Candles {
		result[i] = (candle.Open + candle.High + candle.Low + candle.Close) / 4
	}
	return result
}

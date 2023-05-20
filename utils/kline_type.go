package utils

import (
	"sort"
	"time"

	"github.com/idoall/stockindicator/container/bst"
)

// Kline struct
type Kline struct {
	Amount float64 // 成交额
	Count  int64   // 成交笔数
	Open   float64 // 开盘价
	Close  float64 // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 // 最低价
	High   float64 // 最高价
	Volume float64 // 成交量
	// 涨跌幅
	ChangePercent float64
	// 是否阳线上涨
	IsBullMarket bool
	Time         time.Time // k线时间
}

type Klines []Kline

// OHLC is a connector for technical analysis usage
type OHLC struct {
	Open          []float64
	High          []float64
	Low           []float64
	Close         []float64
	Volume        []float64
	ChangePercent []float64
	IsBullMarket  []bool
}

// GetOHLC returns the entire subset of candles as a friendly type for gct
// technical analysis usage.
func (e Klines) GetOHLC() *OHLC {
	ohlc := &OHLC{
		Open:          make([]float64, len(e)),
		High:          make([]float64, len(e)),
		Low:           make([]float64, len(e)),
		Close:         make([]float64, len(e)),
		Volume:        make([]float64, len(e)),
		ChangePercent: make([]float64, len(e)),
		IsBullMarket:  make([]bool, len(e)),
	}
	for x := range e {
		ohlc.Open[x] = e[x].Open
		ohlc.High[x] = e[x].High
		ohlc.Low[x] = e[x].Low
		ohlc.Close[x] = e[x].Close
		ohlc.Volume[x] = e[x].Volume
		ohlc.ChangePercent[x] = e[x].ChangePercent
		ohlc.IsBullMarket[x] = e[x].IsBullMarket
	}
	return ohlc
}

// RemoveDuplicates 删除任何重复的蜡烛
func (e Klines) RemoveDuplicates() Klines {
	var result Klines
	lookup := make(map[int64]bool)
	target := 0
	for _, keep := range e {
		if key := keep.Time.Unix(); !lookup[key] {
			lookup[key] = true
			e[target] = keep
			target++
		}
	}
	result = e[:target]
	return result
}

// RemoveOutsideRange 删除开始和结束日期之外的所有蜡烛图。
func (e Klines) RemoveOutsideRange(start, end time.Time) Klines {
	var result Klines
	target := 0
	for _, keep := range e {
		if keep.Time.Equal(start) || (keep.Time.After(start) && keep.Time.Before(end)) {
			e[target] = keep
			target++
		}
	}
	result = e[:target]
	return result
}

// SortCandlesByTimestamp 排序
func (e Klines) SortCandlesByTimestamp(desc bool) Klines {
	if desc {
		sort.Slice(e, func(i, j int) bool { return e[i].Time.After(e[j].Time) })
		return e
	}
	sort.Slice(e, func(i, j int) bool { return e[i].Time.Before(e[j].Time) })
	return e
}

// ToHeikinAshi 转换成平均K线
func (e Klines) ToHeikinAshi() Klines {
	var result Klines

	for _, candle := range e {

		var open = candle.Open
		if len(result) > 1 {
			var prev = e[len(result)-1]
			open = (prev.Open + prev.Close) / 2
		}

		var k = Kline{
			Open:  open,
			Close: (candle.Open + candle.High + candle.Low + candle.Close) / 4,
			High:  bst.New().Inserts([]float64{candle.High, candle.Open, candle.Close}).Max().(float64),
			Low:   bst.New().Inserts([]float64{candle.Close, candle.Open, candle.Close}).Min().(float64),
			Time:  candle.Time,
		}
		k.ChangePercent = (candle.Close - candle.Open) / candle.Open
		if candle.Close > candle.Open {
			k.IsBullMarket = true
		}
		result = append(result, k)
	}

	return result
}

// HL2 (最高价+最低价)/2
func (e Klines) HL2() []float64 {
	var result []float64

	for _, candle := range e {
		result = append(result, (candle.High+candle.Low)/2)
	}
	return result
}

// HLC3 (最高价+最低价+收盘价)/3
func (e Klines) HLC3() []float64 {
	var result []float64

	for _, candle := range e {
		result = append(result, (candle.High+candle.Low+candle.Close)/3)
	}
	return result
}

// OHLC4 (开盘价 + 最高价 + 最低价 + 收盘价)/4
func (e Klines) OHLC4() []float64 {
	var result []float64

	for _, candle := range e {
		result = append(result, (candle.Open+candle.High+candle.Low+candle.Close)/4)
	}
	return result
}

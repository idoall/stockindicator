package utils

import (
	"time"
)

// Kline struct
type Kline struct {
	Amount float64   // 成交额
	Count  int64     // 成交笔数
	Open   float64   // 开盘价
	Close  float64   // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64   // 最低价
	High   float64   // 最高价
	Volume float64   // 成交量
	Time   time.Time // k线时间
}

type Klines []Kline

// OHLC is a connector for technical analysis usage
type OHLC struct {
	Open   []float64
	High   []float64
	Low    []float64
	Close  []float64
	Volume []float64
}

// GetOHLC returns the entire subset of candles as a friendly type for gct
// technical analysis usage.
func (k Klines) GetOHLC() *OHLC {
	ohlc := &OHLC{
		Open:   make([]float64, len(k)),
		High:   make([]float64, len(k)),
		Low:    make([]float64, len(k)),
		Close:  make([]float64, len(k)),
		Volume: make([]float64, len(k)),
	}
	for x := range k {
		ohlc.Open[x] = k[x].Open
		ohlc.High[x] = k[x].High
		ohlc.Low[x] = k[x].Low
		ohlc.Close[x] = k[x].Close
		ohlc.Volume[x] = k[x].Volume
	}
	return ohlc
}

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

func (e Klines) GetCloses() []float64 {
	var vals = make([]float64, len(e))
	for i, v := range e {
		vals[i] = v.Close
	}
	return vals
}

func (e Klines) GetHighs() []float64 {
	var vals = make([]float64, len(e))
	for i, v := range e {
		vals[i] = v.High
	}
	return vals
}

func (e Klines) GetLows() []float64 {
	var vals = make([]float64, len(e))
	for i, v := range e {
		vals[i] = v.Low
	}
	return vals
}

func (e Klines) GetOpen() []float64 {
	var vals = make([]float64, len(e))
	for i, v := range e {
		vals[i] = v.Open
	}
	return vals
}

func (e Klines) GetVolumes() []float64 {
	var vals = make([]float64, len(e))
	for i, v := range e {
		vals[i] = v.Volume
	}
	return vals
}

package trend

import (
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Atr struct
type Atr struct {
	Period int //默认一般是13
	Name   string
	data   []AtrData
	kline  utils.Klines
}

type AtrData struct {
	Time time.Time
	TR   float64
	Atr  float64
}

// NewAtr new Func
func NewAtr(list utils.Klines, period int) *Atr {
	m := &Atr{Name: "Atr", kline: list, Period: period}
	return m
}

// NewAtr new Func
func NewDefaultAtr(list utils.Klines) *Atr {
	return NewAtr(list, 14)
}

// Calculation Func
func (e *Atr) Calculation() *Atr {

	var tr []float64
	for i := 0; i < len(e.kline); i++ {
		klineItem := e.kline[i]
		var AtrPointStruct AtrData
		// TR= | 最高价 - 最低价 | 和 | 最高价 - 昨日收盘价 | 和 | 昨日收盘价 - 最低价 | 的最大值
		var prevClose float64
		if i != 0 {
			prevClose = e.kline[i-1].Close
		}
		AtrPointStruct.TR = math.Max(klineItem.High-klineItem.Low, math.Max(klineItem.High-prevClose, klineItem.Low-prevClose))
		AtrPointStruct.Time = e.kline[i].Time
		e.data = append(e.data, AtrPointStruct)

		tr = append(tr, AtrPointStruct.TR)
	}

	for i, v := range e.data {
		if i == 0 {
			continue
		}
		e.data[i].Atr = (v.TR + float64(e.Period-1)*e.data[i-1].TR) / float64(e.Period)
	}
	// var atr = NewEma(utils.CloseArrayToKline(tr), e.Period).Calculation().GetData()
	// for i := 0; i < len(atr); i++ {
	// 	e.data[i].Atr = atr[i].Value
	// }
	return e
}

// Chandelier Exit. 根据平均真实值 (Atr) 设置追踪止损。
//
// Chandelier Exit Long = 22-Period SMA High - Atr(22) * 3
// Chandelier Exit Short = 22-Period SMA Low + Atr(22) * 3
//
// Returns chandelierExitLong, chandelierExitShort
func (e *Atr) ChandelierExit(period int) ([]float64, []float64) {

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var high, low []float64
	for _, v := range e.kline {
		high = append(high, v.High)
		low = append(low, v.Low)
	}

	highestHigh22 := utils.Max(period, high)
	lowestLow22 := utils.Min(period, low)

	chandelierExitLong := make([]float64, len(e.data))
	chandelierExitShort := make([]float64, len(e.data))

	for i := 0; i < len(chandelierExitLong); i++ {
		chandelierExitLong[i] = highestHigh22[i] - (e.data[i].Atr * float64(3))
		chandelierExitShort[i] = lowestLow22[i] + (e.data[i].Atr * float64(3))
	}

	return chandelierExitLong, chandelierExitShort
}

// GetData Func
func (e *Atr) GetValues() (tr []float64, atr []float64) {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	for _, v := range e.data {
		tr = append(tr, v.TR)
		atr = append(atr, v.Atr)
	}
	return
}

// GetData Func
func (e *Atr) GetData() []AtrData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

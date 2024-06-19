package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Atr struct
type Atr struct {
	Period int //默认一般是13
	Name   string
	data   []AtrData
	kline  *klines.Item
}

type AtrData struct {
	Time time.Time
	TR   float64
	Atr  float64
}

// NewAtr new Func
func NewAtr(klineItem *klines.Item, period int) *Atr {
	m := &Atr{
		Name:   fmt.Sprintf("Atr%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewAtr new Func
func NewDefaultAtr(klineItem *klines.Item) *Atr {
	return NewAtr(klineItem, 14)
}

// Calculation Func
func (e *Atr) Calculation() *Atr {

	var tr = make([]float64, len(e.kline.Candles))

	var ohlc = e.kline.GetOHLC()
	var closes = ohlc.Close
	var highs = ohlc.High
	var lows = ohlc.Low

	for i := 0; i < len(e.kline.Candles); i++ {
		// klineItem := e.kline[i]
		var AtrPointStruct AtrData
		// TR= | 最高价 - 最低价 | 和 | 最高价 - 昨日收盘价 | 和 | 昨日收盘价 - 最低价 | 的最大值
		var prevClose float64
		if i != 0 {
			prevClose = closes[i-1]
		}
		tr[i] = math.Max(highs[i]-lows[i], math.Max(highs[i]-prevClose, lows[i]-prevClose))
		AtrPointStruct.Time = e.kline.Candles[i].Time
	}
	var atr = ta.Rma(e.Period, tr)

	e.data = make([]AtrData, len(e.kline.Candles))
	for i, v := range atr {
		e.data[i] = AtrData{
			Time: e.kline.Candles[i].Time,
			TR:   tr[i],
			Atr:  v,
		}
	}

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

	var ohlc = e.kline.GetOHLC()
	highestHigh22 := ta.Max(period, ohlc.High)
	lowestLow22 := ta.Min(period, ohlc.Low)

	chandelierExitLong := make([]float64, len(e.data))
	chandelierExitShort := make([]float64, len(e.data))

	for i := 0; i < len(chandelierExitLong); i++ {
		chandelierExitLong[i] = highestHigh22[i] - (e.data[i].Atr * float64(3))
		chandelierExitShort[i] = lowestLow22[i] + (e.data[i].Atr * float64(3))
	}

	return chandelierExitLong, chandelierExitShort
}

// GetValues Func
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

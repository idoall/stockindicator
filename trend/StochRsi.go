package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// StochRsi is the main object
type StochRsi struct {
	Name    string
	SmoothK int
	SmoothD int
	// RSI 长度
	RsiLength int
	// 随机指标长度
	StochLength int
	data        []StochRsiData
	kline       *klines.Item
}

type StochRsiData struct {
	K    float64
	D    float64
	Time time.Time
}

// NewStochRsi new Func
func NewStochRsi(klineItem *klines.Item, smoothK, smoothD, rsiLength, stochLength int) *StochRsi {
	m := &StochRsi{
		Name:        fmt.Sprintf("StochRsi%d-%d-%d-%d", smoothK, smoothD, rsiLength, stochLength),
		kline:       klineItem,
		SmoothK:     smoothK,
		SmoothD:     smoothD,
		RsiLength:   rsiLength,
		StochLength: stochLength,
	}
	return m
}

func NewDefaultStochRsi(klineItem *klines.Item) *StochRsi {
	return NewStochRsi(klineItem, 3, 3, 14, 14)
}

// Calculation Func
func (e *StochRsi) Calculation() *StochRsi {

	rsi := NewRsi(e.kline, e.RsiLength).GetValue()

	highestHigh := ta.Max(e.StochLength, rsi)
	lowestLow := ta.Min(e.StochLength, rsi)

	k := ta.Ema(e.SmoothK, ta.MultiplyBy(ta.Divide(ta.Subtract(rsi, lowestLow), ta.Subtract(highestHigh, lowestLow)), float64(100)))
	d := ta.Ema(e.SmoothD, k)

	for i := 0; i < len(k); i++ {
		e.data = append(e.data, StochRsiData{
			Time: time.Unix(e.kline.Candles[i].TimeUnix, 0),
			K:    k[i],
			D:    d[i],
		})
	}
	return e
}

// AnalysisSide Func
func (e *StochRsi) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]

		if v.K > v.D && prevItem.K < prevItem.D {
			sides[i] = utils.Buy
		} else if v.K < v.D && prevItem.K > prevItem.D {
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
func (e *StochRsi) GetData() []StochRsiData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
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
	kline       utils.Klines
}

type StochRsiData struct {
	K    float64
	D    float64
	Time time.Time
}

// NewStochRsi new Func
func NewStochRsi(list utils.Klines, smoothK, smoothD, rsiLength, stochLength int) *StochRsi {
	m := &StochRsi{
		Name:        fmt.Sprintf("StochRsi%d-%d-%d-%d", smoothK, smoothD, rsiLength, stochLength),
		kline:       list,
		SmoothK:     smoothK,
		SmoothD:     smoothD,
		RsiLength:   rsiLength,
		StochLength: stochLength,
	}
	return m
}

func NewDefaultStochRsi(list utils.Klines) *StochRsi {
	return NewStochRsi(list, 3, 3, 14, 14)
}

// Calculation Func
func (e *StochRsi) Calculation() *StochRsi {

	rsi := NewRsi(e.kline, e.RsiLength).GetValue()

	highestHigh := ta.Max(e.StochLength, rsi)
	lowestLow := ta.Min(e.StochLength, rsi)

	k := NewSma(utils.CloseArrayToKline(ta.MultiplyBy(ta.Divide(ta.Subtract(rsi, lowestLow), ta.Subtract(highestHigh, lowestLow)), float64(100))), e.SmoothK).GetValues()
	d := NewSma(utils.CloseArrayToKline(k), e.SmoothD).GetValues()

	for i := 0; i < len(k); i++ {
		e.data = append(e.data, StochRsiData{
			Time: e.kline[i].Time,
			K:    k[i],
			D:    d[i],
		})
	}
	return e
}

// AnalysisSide Func
func (e *StochRsi) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

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

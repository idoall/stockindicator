package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
)

// Rsi is the main object
type Rsi struct {
	Name   string
	Period int //默认计算几天的
	data   []RsiData
	kline  *klines.Item
}

type RsiData struct {
	Value float64
	Time  time.Time
}

// NewRsi new Func
// 使用方法，先添加最早日期的数据,最后一条应该是当前日期的数据，结果与 AICoin 对比完全一致
func NewRsi(klineItem *klines.Item, period int) *Rsi {
	m := &Rsi{
		Name:   fmt.Sprintf("Rsi%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

func NewDefaultRsi(klineItem *klines.Item) *Rsi {
	return NewRsi(klineItem, 14)
}

// Calculation Func
func (e *Rsi) Calculation() *Rsi {

	rsiArray := e.rsi(e.kline.GetOHLC().Close, e.Period)

	rsiArrayLen := len(rsiArray)
	for i := 0; i <= (rsiArrayLen - 1); i++ {
		var p RsiData
		p.Time = e.kline.Candles[i].Time
		p.Value = rsiArray[i]
		e.data = append(e.data, p)
	}
	return e
}

// AnalysisSide Func
func (e *Rsi) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		// 当值低于20认为是超卖，为买入信号
		if v.Value < 20 && e.data[i-1].Value > 20 {
			sides[i] = utils.Buy
		} else if v.Value < 80 && e.data[i-1].Value > 80 {
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
func (e *Rsi) GetData() []RsiData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// GetValue return Value
func (e *Rsi) GetValue() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var result []float64
	for _, v := range e.data {
		result = append(result, v.Value)
	}
	return result
}

func (e *Rsi) rsi(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if len(inReal) < inTimePeriod {
		return outReal
	}

	if inTimePeriod < 2 {
		return outReal
	}

	// variable declarations
	tempValue1 := 0.0
	tempValue2 := 0.0
	outIdx := inTimePeriod
	today := 0
	prevValue := inReal[today]
	prevGain := 0.0
	prevLoss := 0.0
	today++

	for i := inTimePeriod; i > 0; i-- {
		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
	}

	prevLoss /= float64(inTimePeriod)
	prevGain /= float64(inTimePeriod)

	if today > 0 {

		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++

	} else {

		for today < 0 {
			tempValue1 = inReal[today]
			tempValue2 = tempValue1 - prevValue
			prevValue = tempValue1
			prevLoss *= float64(inTimePeriod - 1)
			prevGain *= float64(inTimePeriod - 1)
			if tempValue2 < 0 {
				prevLoss -= tempValue2
			} else {
				prevGain += tempValue2
			}
			prevLoss /= float64(inTimePeriod)
			prevGain /= float64(inTimePeriod)
			today++
		}
	}

	for today < len(inReal) {

		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		prevLoss *= float64(inTimePeriod - 1)
		prevGain *= float64(inTimePeriod - 1)
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
		prevLoss /= float64(inTimePeriod)
		prevGain /= float64(inTimePeriod)
		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}

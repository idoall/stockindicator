package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Sma struct
type Sma struct {
	Name   string
	Period int //默认计算几天的MA,KDJ一般是9，OBV是10、20、30
	data   []SmaData
	kline  *klines.Item
}

type SmaData struct {
	Value float64
	Time  time.Time
}

// NewSma new Func
func NewSma(klineItem *klines.Item, period int) *Sma {
	m := &Sma{
		Name:   fmt.Sprintf("Sma%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultSma new Func
func NewDefaultSma(klineItem *klines.Item) *Sma {
	return NewSma(klineItem, 9)
}

// Calculation Func
func (e *Sma) Calculation() *Sma {

	var period = e.Period

	smas := ta.Sma(period, e.kline.GetOHLC().Close)
	e.data = make([]SmaData, len(e.kline.Candles))

	for i, sma := range smas {
		e.data[i] = SmaData{
			Time:  e.kline.Candles[i].Time,
			Value: sma,
		}
	}

	return e
}

// GetPoints Func
func (e *Sma) GetData() []SmaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// GetValues return Values
func (e *Sma) GetValues() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	val := make([]float64, len(e.data))
	for i, v := range e.data {
		val[i] = v.Value
	}
	// fmt.Println(val)
	return val
}

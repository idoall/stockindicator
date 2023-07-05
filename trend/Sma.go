package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Sma struct
type Sma struct {
	Name   string
	Period int //默认计算几天的MA,KDJ一般是9，OBV是10、20、30
	data   []SmaData
	kline  utils.Klines
}

type SmaData struct {
	Value float64
	Time  time.Time
}

// NewSma new Func
func NewSma(list utils.Klines, period int) *Sma {
	m := &Sma{
		Name:   fmt.Sprintf("Sma%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultSma new Func
func NewDefaultSma(list utils.Klines) *Sma {
	return NewSma(list, 9)
}

// Calculation Func
func (e *Sma) Calculation() *Sma {

	var period = e.Period

	smas := utils.Sma(period, e.kline.GetOHLC().Close)
	e.data = make([]SmaData, len(e.kline))

	for i, sma := range smas {
		e.data[i] = SmaData{
			Time:  e.kline[i].Time,
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

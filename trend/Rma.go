package trend

import (
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Rolling Moving Average (Rma).
type Rma struct {
	Period int //默认13
	data   []RmaData
	kline  utils.Klines
}

type RmaData struct {
	Value float64
	Time  time.Time
}

// NewRma new Func
func NewRma(list utils.Klines, period int) *Rma {
	m := &Rma{kline: list, Period: period}
	return m
}

// Calculation Func
func (e *Rma) Calculation() *Rma {

	result := make([]float64, len(e.kline))

	closeing := make([]float64, len(e.kline))
	for _, v := range e.kline {
		closeing = append(closeing, v.Close)
	}

	period := e.Period
	sum := float64(0)

	for i, value := range closeing {
		count := i + 1

		if i < period {
			sum += value
		} else {
			sum = (result[i-1] * float64(period-1)) + value
			count = period
		}

		e.data = append(e.data, RmaData{
			Time:  e.kline[i].Time,
			Value: sum / float64(count),
		})
	}

	return e
}

// GetData Func
func (e *Rma) GetData() []RmaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

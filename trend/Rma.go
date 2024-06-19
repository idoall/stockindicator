package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Rolling Moving Average (Rma).
type Rma struct {
	Name   string
	Period int //默认13
	data   []RmaData
	kline  *klines.Item
}

type RmaData struct {
	Value float64
	Time  time.Time
}

// NewRma new Func
func NewRma(klineItem *klines.Item, period int) *Rma {
	m := &Rma{
		Name:   fmt.Sprintf("Rma%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

func NewDefaultRma(klineItem *klines.Item) *Rma {
	return NewRma(klineItem, 13)
}

// Calculation Func
func (e *Rma) Calculation() *Rma {

	var rma = ta.Rma(e.Period, e.kline.GetOHLC().Close)

	e.data = make([]RmaData, len(rma))
	for i, value := range rma {

		e.data = append(e.data, RmaData{
			Time:  e.kline.Candles[i].Time,
			Value: value,
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

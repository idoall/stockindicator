package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Ma struct
type Ma struct {
	Name   string
	Period int //默认计算几天的Ma,KDJ一般是9，OBV是10、20、30
	data   []MaData
	kline  *klines.Item
}

type MaData struct {
	Value float64
	Time  time.Time
}

// NewMa new Func
func NewMa(klineItem *klines.Item, period int) *Ma {
	m := &Ma{
		Name:   fmt.Sprintf("Ma%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultMa new Func
func NewDefaultMa(klineItem *klines.Item) *Ma {
	return NewMa(klineItem, 20)
}

// GetPoints return Point
func (e *Ma) GetData() []MaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	return e.data
}

// Calculation Func
func (e *Ma) Calculation() *Ma {

	closes := e.kline.GetOHLC().Close

	// get maData
	maData := ta.Sma(e.Period, closes)

	for i := 0; i < len(maData); i++ {
		p := MaData{}
		p.Time = e.kline.Candles[i].Time
		p.Value = maData[i]
		e.data = append(e.data, p)
	}
	return e
}

package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/ta"
)

// EMAVegas struct
type EMAVegas struct {
	Name         string
	Period       int
	PeriodShort1 int
	PeriodShort2 int
	PeriodLong1  int
	PeriodLong2  int
	data         []EMAVegasData
	kline        utils.Klines
}

// EMAVegasData EMAVegas
type EMAVegasData struct {
	Value       float64
	Short1Value float64
	Short2Value float64
	Long1Value  float64
	Long2Value  float64
	Time        time.Time
}

// NewEMAVegas new Func
func NewEMAVegas(list utils.Klines, period, periodShort1, periodShort2, periodLong1, periodLong2 int) *EMAVegas {
	m := &EMAVegas{
		Name:         fmt.Sprintf("EMAVegas%d-%d-%d-%d", periodShort1, periodShort2, periodLong1, periodLong2),
		kline:        list,
		Period:       period,
		PeriodShort1: periodShort1,
		PeriodShort2: periodShort2,
		PeriodLong1:  periodLong1,
		PeriodLong2:  periodLong2,
	}
	return m
}

// NewDefaultEMAVegas new Func
func NewDefaultEMAVegas(list utils.Klines) *EMAVegas {
	return NewEMAVegas(list, 12, 144, 169, 575, 676)
}

// Calculation Func
func (e *EMAVegas) Calculation() *EMAVegas {

	e.data = make([]EMAVegasData, len(e.kline))

	var closeing = e.kline.GetOHLC().Close

	ema := ta.Ema(e.Period, closeing)
	emaShort1 := ta.Ema(e.PeriodShort1, closeing)
	emaShort2 := ta.Ema(e.PeriodShort2, closeing)
	emaLong1 := ta.Ema(e.PeriodLong1, closeing)
	emaLong2 := ta.Ema(e.PeriodLong2, closeing)

	for i := 0; i < len(emaShort1); i++ {
		e.data[i] = EMAVegasData{
			Time:        e.kline[i].Time,
			Value:       ema[i],
			Short1Value: emaShort1[i],
			Short2Value: emaShort2[i],
			Long1Value:  emaLong1[i],
			Long2Value:  emaLong2[i],
		}
	}
	return e
}

// AnalysisSide Func
func (e *EMAVegas) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}
		var prevKlineItem = e.kline[i-1]
		var klineItem = e.kline[i]
		prevItem := e.data[i-1]
		// var item = e.data[i]

		// if klineItem.Close > klineItem.Open && ((klineItem.Close > v.Short2Value && prevKlineItem.Close < prevItem.Short2Value) || (klineItem.Close > v.Long2Value && prevKlineItem.Close < prevItem.Short2Value)) {
		// 	sides[i] = utils.Buy
		// } else if klineItem.Close < klineItem.Open && ((klineItem.Close < v.Short2Value && prevKlineItem.Close > prevItem.Short2Value) || (klineItem.Close < v.Long2Value && prevKlineItem.Close > prevItem.Short2Value)) {
		// 	sides[i] = utils.Sell
		// } else {
		// 	sides[i] = utils.Hold
		// }

		// 当前是阳线 并且 收盘价大于 Long2 并且 上一根收盘价小于 Long2
		if klineItem.Close > klineItem.Open && ((klineItem.Close > v.Long2Value && prevKlineItem.Close < prevItem.Long2Value) || (klineItem.Close > v.Short2Value && prevKlineItem.Close < v.Short2Value && v.Short1Value > v.Long2Value && v.Short2Value > v.Long2Value)) {
			sides[i] = utils.Buy
		} else if klineItem.Close < klineItem.Open && klineItem.Close < v.Short2Value && prevKlineItem.Close > prevItem.Short2Value {
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

// GetData return Point
func (e *EMAVegas) GetData() []EMAVegasData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

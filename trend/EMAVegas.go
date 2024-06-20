package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
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
	kline        *klines.Item
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
func NewEMAVegas(klineItem *klines.Item, period, periodShort1, periodShort2, periodLong1, periodLong2 int) *EMAVegas {
	m := &EMAVegas{
		Name:         fmt.Sprintf("EMAVegas%d-%d-%d-%d", periodShort1, periodShort2, periodLong1, periodLong2),
		kline:        klineItem,
		Period:       period,
		PeriodShort1: periodShort1,
		PeriodShort2: periodShort2,
		PeriodLong1:  periodLong1,
		PeriodLong2:  periodLong2,
	}
	return m
}

// NewDefaultEMAVegas new Func
func NewDefaultEMAVegas(klineItem *klines.Item) *EMAVegas {
	return NewEMAVegas(klineItem, 12, 144, 169, 575, 676)
}

// Calculation Func
func (e *EMAVegas) Calculation() *EMAVegas {

	e.data = make([]EMAVegasData, len(e.kline.Candles))

	var closeing = e.kline.GetOHLC().Close

	ema := ta.Ema(e.Period, closeing)
	emaShort1 := ta.Ema(e.PeriodShort1, closeing)
	emaShort2 := ta.Ema(e.PeriodShort2, closeing)
	emaLong1 := ta.Ema(e.PeriodLong1, closeing)
	emaLong2 := ta.Ema(e.PeriodLong2, closeing)

	for i := 0; i < len(emaShort1); i++ {
		e.data[i] = EMAVegasData{
			Time:        e.kline.Candles[i].Time,
			Value:       ema[i],
			Short1Value: emaShort1[i],
			Short2Value: emaShort2[i],
			Long1Value:  emaLong1[i],
			Long2Value:  emaLong2[i],
		}
	}
	clear(closeing)
	clear(ema)
	clear(emaShort1)
	clear(emaShort2)
	clear(emaLong1)
	clear(emaLong2)
	return e
}

// AnalysisSide Func
func (e *EMAVegas) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i := range e.data {
		if i < 1 {
			continue
		}
		var klineItem = e.kline.Candles[i]
		var item = e.data[i]

		var t1, d1 bool

		if (i - 10 - 1) < 0 {
			continue
		}

		for k := i - 10; k < i; k++ {
			data := e.data[k]
			dataPrev := e.data[k-1]
			if data.Short2Value > data.Long2Value && dataPrev.Short2Value < dataPrev.Long2Value {
				t1 = true // 近 10 根有金叉
			}
			if data.Short2Value < data.Long2Value && dataPrev.Short2Value > dataPrev.Long2Value {
				d1 = true // 近 10 根有死叉
			}
		}

		if item.Value > item.Short1Value && // ema12>ema144
			item.Short1Value > item.Long1Value && item.Short1Value > item.Long2Value && // 144<575 && 144 > 676
			item.Short2Value > item.Long1Value && item.Short2Value > item.Long2Value && //169>575 && 169 > 676
			klineItem.IsBullMarket && t1 {
			sides[i] = utils.Buy
		} else if item.Value < item.Short1Value &&
			item.Short1Value < item.Long1Value && item.Short1Value < item.Long2Value &&
			item.Short2Value < item.Long1Value && item.Short2Value < item.Long2Value &&
			!klineItem.IsBullMarket && d1 {
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

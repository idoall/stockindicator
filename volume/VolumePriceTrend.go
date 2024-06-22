package volume

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// The Volume Price Trend (VPT) provides a correlation between the
// volume and the price.
//
// VPT = Previous VPT + (Volume * (Current Closing - Previous Closing) / Previous Closing)
type VolumePriceTrend struct {
	Name   string
	Period int
	data   []VolumePriceTrendData
	kline  *klines.Item
}

// VolumePriceTrendData
type VolumePriceTrendData struct {
	Time  time.Time
	Value float64
}

// NewVolumePriceTrend new Func
func NewVolumePriceTrend(klineItem *klines.Item, period int) *VolumePriceTrend {
	m := &VolumePriceTrend{
		Name:   fmt.Sprintf("VolumePriceTrend%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultVolumePriceTrend new Func
func NewDefaultVolumePriceTrend(klineItem *klines.Item) *VolumePriceTrend {
	return NewVolumePriceTrend(klineItem, 14)
}

// Calculation Func
func (e *VolumePriceTrend) Calculation() *VolumePriceTrend {

	period := e.Period
	var ohlc = e.kline.GetOHLC()
	var closing = ohlc.Close
	var volume = ohlc.Volume

	previousClosing := ta.ShiftRightAndFillBy(period, closing[0], closing)
	vpt := ta.Multiply(volume, ta.Divide(ta.Subtract(closing, previousClosing), previousClosing))
	vals := ta.Sum(len(vpt), vpt)

	for i := 0; i < len(vals); i++ {
		e.data = append(e.data, VolumePriceTrendData{
			Time:  time.Unix(e.kline.Candles[i].TimeUnix, 0),
			Value: vals[i],
		})
	}

	return e
}

// AnalysisSide Func
// func (e *VolumePriceTrend) AnalysisSide() utils.SideData {
// 	sides := make([]utils.Side, len(e.kline.Candles))

// 	if len(e.data) == 0 {
// 		e = e.Calculation()
// 	}

// 	for i, v := range e.data {
// 		if i < 1 {
// 			continue
// 		}

// 		var prevItem = e.data[i-1]

// 		if v.Value < 10 && prevItem.Value > 10 {
// 			sides[i] = utils.Buy
// 		} else if v.Value > 90 && prevItem.Value < 90 {
// 			sides[i] = utils.Sell
// 		} else {
// 			sides[i] = utils.Hold
// 		}
// 	}
// 	return utils.SideData{
// 		Name: e.Name,
// 		Data: sides,
// 	}
// }

// GetData Func
func (e *VolumePriceTrend) GetData() []VolumePriceTrendData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

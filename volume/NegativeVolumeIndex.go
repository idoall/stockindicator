package volume

import (
	"time"

	"github.com/idoall/stockindicator/utils/klines"
)

// The Negative Volume Index (NVI) is a cumulative indicator using
// the change in volume to decide when the smart money is active.
//
// If Volume is greather than Previous Volume:
//
//	NVI = Previous NVI
//
// Otherwise:
//
//	NVI = Previous NVI + (((Closing - Previous Closing) / Previous Closing) * Previous NVI)
type NegativeVolumeIndex struct {
	Name  string
	data  []NegativeVolumeIndexData
	kline *klines.Item
}

// NegativeVolumeIndexData
type NegativeVolumeIndexData struct {
	Time  time.Time
	Value float64
}

// NewNegativeVolumeIndex new Func
func NewNegativeVolumeIndex(klineItem *klines.Item) *NegativeVolumeIndex {
	m := &NegativeVolumeIndex{
		Name:  "NegativeVolumeIndex",
		kline: klineItem,
	}
	return m
}

// NewDefaultNegativeVolumeIndex new Func
func NewDefaultNegativeVolumeIndex(klineItem *klines.Item) *NegativeVolumeIndex {
	return NewNegativeVolumeIndex(klineItem)
}

// Calculation Func
func (e *NegativeVolumeIndex) Calculation() *NegativeVolumeIndex {

	var ohlc = e.kline.GetOHLC()
	var closing = ohlc.Close
	var volume = ohlc.Volume

	nvi := make([]float64, len(closing))

	for i := 0; i < len(nvi); i++ {
		if i == 0 {
			nvi[i] = 1000
		} else if volume[i-1] < volume[i] {
			nvi[i] = nvi[i-1]
		} else {
			nvi[i] = nvi[i-1] + (((closing[i] - closing[i-1]) / closing[i-1]) * nvi[i-1])
		}
	}

	for i := 0; i < len(nvi); i++ {
		e.data = append(e.data, NegativeVolumeIndexData{
			Time:  time.Unix(e.kline.Candles[i].TimeUnix, 0),
			Value: nvi[i],
		})
	}

	return e
}

// AnalysisSide Func
// func (e *NegativeVolumeIndex) AnalysisSide() utils.SideData {
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
func (e *NegativeVolumeIndex) GetData() []NegativeVolumeIndexData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

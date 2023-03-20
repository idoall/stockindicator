package volume

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

// The Ease of Movement (EMV) 简易波动指标（Ease of Movement Value）又称EMV指标.
// 它是由RichardW．ArmJr．根据等量图和压缩图的原理设计而成,目的是将价格与成交量的变化结合成一个波动指标来反映股价或指数的变动状况。
// 由于股价的变化和成交量的变化都可以引发该指标数值的变动,因此,EMV实际上也是一个量价合成指标。
//
// Distance Moved = ((High + Low) / 2) - ((Priod High + Prior Low) /2)
// Box Ratio = ((Volume / 100000000) / (High - Low))
// EMV(1) = Distance Moved / Box Ratio
// EMV(14) = SMA(14, EMV(1))
type EaseOfMovement struct {
	Name   string
	Period int
	data   []EaseOfMovementData
	kline  utils.Klines
}

// EaseOfMovementData
type EaseOfMovementData struct {
	Time  time.Time
	Value float64
}

// NewEaseOfMovement new Func
func NewEaseOfMovement(list utils.Klines, period int) *EaseOfMovement {
	m := &EaseOfMovement{
		Name:   fmt.Sprintf("EaseOfMovement%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultEaseOfMovement new Func
func NewDefaultEaseOfMovement(list utils.Klines) *EaseOfMovement {
	return NewEaseOfMovement(list, 14)
}

// Calculation Func
func (e *EaseOfMovement) Calculation() *EaseOfMovement {

	period := e.Period
	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low
	var volume = ohlc.Volume

	distanceMoved := utils.Diff(utils.DivideBy(utils.Add(high, low), 2), 1)
	boxRatio := utils.Divide(utils.DivideBy(volume, float64(100000000)), utils.Subtract(high, low))

	emv := trend.NewSma(utils.CloseArrayToKline(utils.Divide(distanceMoved, boxRatio)), period).GetValues()

	for i := 0; i < len(emv); i++ {
		e.data = append(e.data, EaseOfMovementData{
			Time:  e.kline[i].Time,
			Value: emv[i],
		})
	}

	return e
}

// AnalysisSide Func
//
// 1、当EMV由下往上穿越0轴时，买进。
// 2、当EMV由上往下穿越0轴时，卖出。
// func (e *EaseOfMovement) AnalysisSide() utils.SideData {
// 	sides := make([]utils.Side, len(e.kline))

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
func (e *EaseOfMovement) GetData() []EaseOfMovementData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

package volume

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/ta"
)

// The Force Index (FI) 强力指数指标，简称FI指标,它是由Alexander elder发明的.
// FI指标是用来指示上升或下降趋势的力量大小，它是通过在零线上下移动来表示趋势的强弱。
//
// Force Index = EMA(period, (Current - Previous) * Volume)
type ForceIndex struct {
	Name   string
	Period int
	data   []ForceIndexData
	kline  utils.Klines
}

// ForceIndexData
type ForceIndexData struct {
	Time  time.Time
	Value float64
}

// NewForceIndex new Func
func NewForceIndex(list utils.Klines, period int) *ForceIndex {
	m := &ForceIndex{
		Name:   fmt.Sprintf("ForceIndex%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultForceIndex new Func
func NewDefaultForceIndex(list utils.Klines) *ForceIndex {
	return NewForceIndex(list, 13)
}

// Calculation Func
func (e *ForceIndex) Calculation() *ForceIndex {

	period := e.Period
	var ohlc = e.kline.GetOHLC()
	var closing = ohlc.Close
	var volume = ohlc.Volume

	var vals = trend.NewEma(utils.CloseArrayToKline(ta.Multiply(ta.Diff(closing, 1), volume)), period).GetValues()

	for i := 0; i < len(vals); i++ {
		e.data = append(e.data, ForceIndexData{
			Time:  e.kline[i].Time,
			Value: vals[i],
		})
	}

	return e
}

// AnalysisSide Func
// 　　价格上升，而FI在零线以上，呈上升，则表示价格上升趋势会继续
// 　　价格上升，而FI在零或者趋向于零线时，表示价格上升趋势将要结束
// 　　价格下降，而FI在零线以下，呈下降，则表示价格下降趋势会继续
// 　　价格下降，而FI在零或者趋向于零线时，表示价格下降趋势将要结束
// 　　当强力指数指标在增加趋势的时间段内是负值时，是买入的信号
// 　　当强力指数指标在下降趋势的时间段内变为正值时，是卖出的信号
// 　　如果价格的变化不与相对应的成交量变化相关联，强力指数指标仍旧保持在一个水平上，意味着趋势很快发生变化。
// func (e *ForceIndex) AnalysisSide() utils.SideData {
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
func (e *ForceIndex) GetData() []ForceIndexData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

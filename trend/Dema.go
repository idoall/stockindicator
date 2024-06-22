package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Dema struct
type Dema struct {
	Name   string
	Period int //默认计算几天的Dema
	data   []*DemaData
	ohlc   *klines.OHLC
}

// DemaData Dema函数计算给定期间的双指数移动平均线 (Dema)。

// 双指数移动平均线 (Dema) 是由 Patrick Mulloy 引入的技术指标。目的是减少技术交易者使用的价格图表中存在的噪音量。Dema 使用两个指数移动平均线 (EMA) 来消除滞后。当价格高于平均水平时，它有助于确认上升趋势，当价格低于平均水平时，它有助于确认下降趋势。当价格超过平均线时，可能表示趋势发生变化。
type DemaData struct {
	Value float64
	Time  time.Time
}

// NewDema new Func
func NewDema(klineItem *klines.Item, period int) *Dema {
	m := &Dema{
		Name:   fmt.Sprintf("Dema%d", period),
		Period: period,
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

// NewDema new Func
func NewDemaOHLC(ohlc *klines.OHLC, period int) *Dema {
	m := &Dema{
		Name:   fmt.Sprintf("Dema%d", period),
		ohlc:   ohlc,
		Period: period,
	}
	return m
}

// NewDefaultDema new Func
func NewDefaultDema(klineItem *klines.Item) *Dema {
	return NewDema(klineItem, 20)
}

// Calculation Func
func (e *Dema) Calculation() *Dema {

	period := e.Period
	var close = e.ohlc.Close
	var times = e.ohlc.Time

	e.data = make([]*DemaData, len(close))

	demas := ta.Dema(period, close)

	for i := 0; i < len(demas); i++ {
		e.data[i] = &DemaData{
			Time:  times[i],
			Value: demas[i],
		}
	}
	return e
}

// GetData return Point
func (e *Dema) GetData() []*DemaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// GetValues return Values
func (e *Dema) GetValues() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	val := make([]float64, len(e.data))
	for i, v := range e.data {
		val[i] = v.Value
	}
	return val
}

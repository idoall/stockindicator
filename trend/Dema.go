package trend

import (
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Dema struct
type Dema struct {
	Name   string
	Period int //默认计算几天的Dema
	data   []DemaData
	kline  utils.Klines
}

// DemaData Dema函数计算给定期间的双指数移动平均线 (Dema)。

// 双指数移动平均线 (Dema) 是由 Patrick Mulloy 引入的技术指标。目的是减少技术交易者使用的价格图表中存在的噪音量。Dema 使用两个指数移动平均线 (EMA) 来消除滞后。当价格高于平均水平时，它有助于确认上升趋势，当价格低于平均水平时，它有助于确认下降趋势。当价格超过平均线时，可能表示趋势发生变化。
type DemaData struct {
	Value float64
	Time  time.Time
}

// NewDema new Func
func NewDema(list utils.Klines, period int) *Dema {
	m := &Dema{Name: "Dema", kline: list, Period: period}
	return m
}

// NewDefaultDema new Func
func NewDefaultDema(list utils.Klines) *Dema {
	return NewDema(list, 20)
}

// Calculation Func
func (e *Dema) Calculation() *Dema {

	period := e.Period

	e.data = make([]DemaData, len(e.kline))

	// 计算 EMA1 值
	ema1 := NewEma(e.kline, period).GetValues()

	// 计算 EMA2 值
	ema2 := NewEma(utils.CloseArrayToKline(ema1), period).GetValues()

	// 2 * N日EMA － N日EMA的EMA
	demas := utils.Subtract(utils.MultiplyBy(ema1, 2), ema2)

	for i := 0; i < len(demas); i++ {
		e.data[i] = DemaData{
			Time:  e.kline[i].Time,
			Value: demas[i],
		}
	}
	return e
}

// GetData return Point
func (e *Dema) GetData() []DemaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

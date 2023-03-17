package trend

import (
	"time"

	"github.com/idoall/stockindicator/utils"
)

/*
BBI（Bull and Bear lndex）是一种将不同周期移动平均线加权平均之后的综合指标，属于均线型指标，
一般选用3、6、12、24等4个参数。它是针对普通移动平均线MA指标的一种改进，既有短期移动平均线的灵敏，又有明显的中期趋势特征。

应用法则
价格位于BBI曲线上方，视为多头市场。
价格位于BBI曲线下方，视为空头市场。
高价区收盘价跌破BBI曲线为卖出信号。
低价区收盘价突破BBI曲线为买入信号。
上升回档时，BBI曲线为支持线，可以发挥支撑作用。
下跌反弹时，BBI曲线为压力线，可以发挥阻力作用。
线上阴线买;
线下阳线卖。

BBI=(3日均价+6日均价+12日均价+24日均价)/4
3日均价 = ma(3)
*/

// Bbi struct
type Bbi struct {
	Name                                       string
	periodMa1, periodMa2, periodMa3, periodMa4 int
	data                                       []BbiData
	maPrice                                    []float64
	kline                                      utils.Klines
}

type BbiData struct {
	Value float64
	Time  time.Time
}

// NewBbi new Func
func NewBbi(list utils.Klines, periodMa1, periodMa2, periodMa3, periodMa4 int) *Bbi {
	m := &Bbi{Name: "Bbi", kline: list,
		periodMa1: periodMa1,
		periodMa2: periodMa2,
		periodMa3: periodMa3,
		periodMa4: periodMa4}
	return m
}

// NewDefaultBbi new Func
func NewDefaultBbi(list utils.Klines) *Bbi {
	return NewBbi(list, 3, 6, 12, 24)
}

// Calculation Func
func (e *Bbi) Calculation() *Bbi {
	e.data = make([]BbiData, len(e.kline))

	ma1 := NewMa(e.kline, e.periodMa1).GetData()
	ma2 := NewMa(e.kline, e.periodMa2).GetData()
	ma3 := NewMa(e.kline, e.periodMa3).GetData()
	ma4 := NewMa(e.kline, e.periodMa4).GetData()

	for i := 0; i < len(ma1); i++ {
		e.data[i] = BbiData{
			Time:  e.kline[i].Time,
			Value: (ma1[i].Value + ma2[i].Value + ma3[i].Value + ma4[i].Value) / 4,
		}
	}
	return e
}

// AnalysisSide Func
func (e *Bbi) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]
		var prevPrice = e.kline[i-1].Close
		var price = e.kline[i].Close

		if price > v.Value && prevPrice < prevItem.Value {
			sides[i] = utils.Buy
		} else if price < v.Value && prevPrice > prevItem.Value {
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
func (e *Bbi) GetData() []BbiData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	return e.data
}

// GetValue return Value
func (e *Bbi) GetValue() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var result []float64
	for _, v := range e.data {
		result = append(result, v.Value)
	}
	return result
}

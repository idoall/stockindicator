package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Ema struct
type Ema struct {
	Name   string
	Period int //默认计算几天的Ema
	data   []EmaData
	kline  *klines.Item
}

type EmaData struct {
	Value float64
	Time  time.Time
}

// NewEma new Func
func NewEma(klineItem *klines.Item, period int) *Ema {
	m := &Ema{
		Name:   fmt.Sprintf("Ema%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewEma new Func
func NewDefaultEma(klineItem *klines.Item) *Ema {
	return NewEma(klineItem, 5)
}

// Calculation Func
func (e *Ema) Calculation() *Ema {
	var vals = ta.Ema(e.Period, e.kline.GetOHLC().Close)

	defer func() {
		vals = nil
	}()

	e.data = make([]EmaData, len(vals))
	for i, v := range e.kline.Candles {
		e.data[i] = EmaData{
			Time:  v.Time,
			Value: vals[i],
		}
	}
	return e
}

// GetPoints return Point
func (e *Ema) GetData() []EmaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// GetValues return Values
func (e *Ema) GetValues() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	val := make([]float64, len(e.data))
	for i, v := range e.data {
		val[i] = v.Value
	}
	return val
}

// Add adds a new Value to Ema
// 使用方法，先添加最早日期的数据,最后一条应该是当前日期的数据，结果与 AICoin 对比完全一致
// func (e *Ema) add(timestamp time.Time, value float64) {
// 	p := EmaData{}
// 	p.Time = timestamp

// 	//平滑指数，一般取作2/(N+1)
// 	alpha := 2.0 / (float64(e.Period) + 1.0)

// 	// fmt.Println(alpha)

// 	emaTminusOne := value
// 	if len(e.data) > 0 {
// 		emaTminusOne = e.data[len(e.data)-1].Value
// 	}

// 	// 计算 Ema指数
// 	emaT := alpha*value + (1-alpha)*emaTminusOne
// 	p.Value = emaT
// 	e.data = append(e.data, p)
// }

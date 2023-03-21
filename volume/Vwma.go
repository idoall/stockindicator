package volume

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Vwma struct
type Vwma struct {
	Name   string
	Period int //默认计算几天的Vwma
	data   []VwmaData
	kline  utils.Klines
}

// VwmaPoint 计算成交量加权移动平均线 (Vwma)
// 平均价格数据，重点是成交量，即面积
// 体积越大，权重越大。
//
// Vwma = Sum(Price * Volume) / Sum(Volume) for a given Period.
//
// 返回 vwma
type VwmaData struct {
	Value float64
	Time  time.Time
}

// NewVwma new Func
func NewVwma(list utils.Klines, period int) *Vwma {
	m := &Vwma{
		Name:   fmt.Sprintf("Vwma%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultVwma new Func
func NewDefaultVwma(list utils.Klines) *Vwma {
	return NewVwma(list, 5)
}

// Calculation Func
func (e *Vwma) Calculation() *Vwma {

	var ohlc = e.kline.GetOHLC()
	var closing = ohlc.Close
	var volume = ohlc.Volume

	vwmas := utils.Divide(utils.Sum(e.Period, utils.Multiply(closing, volume)), utils.Sum(e.Period, volume))

	for i := 0; i < len(vwmas); i++ {
		e.data = append(e.data, VwmaData{
			Time:  e.kline[i].Time,
			Value: vwmas[i],
		})
	}
	return e
}

// GetData return Point
func (e *Vwma) GetData() []VwmaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

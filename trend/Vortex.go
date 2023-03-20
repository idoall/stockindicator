package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Vortex struct
type Vortex struct {
	Name   string
	Period int //默认一般是13
	data   []VortexData
	kline  utils.Klines
}

type VortexData struct {
	Time    time.Time
	PlusVi  float64
	MinusVi float64
}

// NewVortex new Func
func NewVortex(list utils.Klines, period int) *Vortex {
	m := &Vortex{
		Name:   fmt.Sprintf("Vortex%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultVortex new Func
func NewDefaultVortex(list utils.Klines) *Vortex {
	return NewVortex(list, 9)
}

// Calculation Func
func (e *Vortex) Calculation() *Vortex {

	period := e.Period

	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low
	var closing = ohlc.Close

	plusVi := make([]float64, len(high))
	minusVi := make([]float64, len(high))

	plusVm := make([]float64, period)
	minusVm := make([]float64, period)
	tr := make([]float64, period)

	var plusVmSum, minusVmSum, trSum float64

	for i := 1; i < len(high); i++ {
		j := i % period

		plusVmSum -= plusVm[j]
		plusVm[j] = math.Abs(high[i] - low[i-1])
		plusVmSum += plusVm[j]

		minusVmSum -= minusVm[j]
		minusVm[j] = math.Abs(low[i] - high[i-1])
		minusVmSum += minusVm[j]

		highLow := high[i] - low[i]
		highPrevClosing := math.Abs(high[i] - closing[i-1])
		lowPrevClosing := math.Abs(low[i] - closing[i-1])

		trSum -= tr[j]
		tr[j] = math.Max(highLow, math.Max(highPrevClosing, lowPrevClosing))
		trSum += tr[j]

		plusVi[i] = plusVmSum / trSum
		minusVi[i] = minusVmSum / trSum
	}

	for i := 0; i < len(plusVi); i++ {
		e.data = append(e.data, VortexData{
			Time:    e.kline[i].Time,
			MinusVi: minusVi[i],
			PlusVi:  plusVi[i],
		})
	}
	return e
}

// GetData Func
func (e *Vortex) GetData() []VortexData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

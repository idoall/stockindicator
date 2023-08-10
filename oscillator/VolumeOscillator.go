package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/ta"
)

// Volume Oscillator. 又名移动平均成交量指标，但是，它并非仅仅计算成交量的移动平均线，
// 而是通过对成交量的长期移动平均线和短期移动平均线之间的比较
// 分析成交量的运行趋势和及时研判趋势转变方向
type VolumeOscillator struct {
	Name        string
	ShortLength int
	LongLength  int
	data        []VolumeOscillatorData
	kline       utils.Klines
}

// VolumeOscillatorData
type VolumeOscillatorData struct {
	Time  time.Time
	Value float64
}

// NewVolumeOscillator new Func
func NewVolumeOscillator(list utils.Klines, shortLength, longLength int) *VolumeOscillator {
	m := &VolumeOscillator{
		Name:        fmt.Sprintf("VolumeOscillator%d-%d", shortLength, longLength),
		kline:       list,
		ShortLength: shortLength,
		LongLength:  longLength,
	}
	return m
}

// NewDefaultVolumeOscillator new Func
func NewDefaultVolumeOscillator(list utils.Klines) *VolumeOscillator {
	return NewVolumeOscillator(list, 5, 10)
}

func (e *VolumeOscillator) Clear() {
	e.data = nil
	e.kline = nil
}

// Calculation Func
func (e *VolumeOscillator) Calculation() *VolumeOscillator {

	var volumes = ta.Nzs(e.kline.GetOHLC().Volume, 0)

	short := trend.NewEma(utils.CloseArrayToKline(volumes), e.ShortLength).GetValues()
	long := trend.NewEma(utils.CloseArrayToKline(volumes), e.LongLength).GetValues()

	oscs := ta.MultiplyBy(ta.Divide(ta.Subtract(short, long), long), 100.0)

	e.data = make([]VolumeOscillatorData, len(oscs))
	for i := 0; i < len(oscs); i++ {
		e.data[i] = VolumeOscillatorData{
			Time:  e.kline[i].Time,
			Value: oscs[i],
		}
	}
	return e
}

// GetData Func
func (e *VolumeOscillator) GetData() []VolumeOscillatorData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// GetValues Func
func (e *VolumeOscillator) GetValues() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	result := make([]float64, len(e.data))
	for i, v := range e.data {
		result[i] = v.Value
	}
	return result
}

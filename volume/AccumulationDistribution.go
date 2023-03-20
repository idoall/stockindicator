package volume

import (
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Accumulation/Distribution Indicator (A/D). 累积/分配指标 (A/D)。 累计指标
// 使用数量和价格来评估股票是否正在累积或分发。
//
// 1、A/D测量资金流向，向上的A/D表明买方占优势，而向下的A/D表明卖方占优势。
// 2、A/D与价格的背离可视为买卖信号，即底背离考虑买入，顶背离考虑卖出
// 3、应当注意A/D忽略了缺口的影响，事实上，跳空缺口的意义是不能轻易忽略的。
// A/D指标无需设置参数。但在应用时，可结合指标的均线进行分析
// MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
// MFV = MFM * Period Volume
// AD = Previous AD + CMFV
type AccumulationDistribution struct {
	Name  string
	data  []AccumulationDistributionData
	kline utils.Klines
}

// AccumulationDistributionData
type AccumulationDistributionData struct {
	Time  time.Time
	Value float64
}

// NewAccumulationDistribution new Func
func NewAccumulationDistribution(list utils.Klines) *AccumulationDistribution {
	m := &AccumulationDistribution{
		Name:  "AccumulationDistribution",
		kline: list,
	}
	return m
}

// NewDefaultAccumulationDistribution new Func
func NewDefaultAccumulationDistribution(list utils.Klines) *AccumulationDistribution {
	return NewAccumulationDistribution(list)
}

// Calculation Func
func (e *AccumulationDistribution) Calculation() *AccumulationDistribution {

	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low
	var closing = ohlc.Close
	var volume = ohlc.Volume

	ad := make([]float64, len(closing))

	for i := 0; i < len(ad); i++ {
		if i > 0 {
			ad[i] = ad[i-1]
		}

		var val = volume[i] * (((closing[i] - low[i]) - (high[i] - closing[i])) / (high[i] - low[i]))

		if math.IsNaN(val) {
			ad[i] += 0
		} else {
			ad[i] += val
		}

		e.data = append(e.data, AccumulationDistributionData{
			Time:  e.kline[i].Time,
			Value: ad[i],
		})
	}

	return e
}

// GetValues return Values
func (e *AccumulationDistribution) GetValues() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	val := make([]float64, len(e.data))
	for i, v := range e.data {
		val[i] = v.Value
	}
	return val
}

// GetData Func
func (e *AccumulationDistribution) GetData() []AccumulationDistributionData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

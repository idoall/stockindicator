package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// AverageDirectionalIndex struct
// ADX（Average Directional Index，平均方向指数）是由J. Welles Wilder在1978年开发的技术分析指标，用于衡量市场趋势的强度。ADX本身不显示趋势的方向，只显示趋势的强弱。ADX通常结合两个方向性指标（+DI 和 -DI）一起使用，以确定市场的上升或下降趋势。
// 2024-07.10 与 TradingView对比，算法、结果一致。
type AverageDirectionalIndex struct {
	Name   string
	Period int
	Len    int
	data   []AverageDirectionalIndexData
	ohlc   *klines.OHLC
}

// AverageDirectionalIndexData AverageDirectionalIndex
type AverageDirectionalIndexData struct {
	DIPlus float64
	DIMins float64
	ADX    float64
	Time   time.Time
}

func (e AverageDirectionalIndexData) String() string {
	return fmt.Sprintf("[%s]DI+:%f DI-:%f ADX:%f", e.Time.Format("2006-01-02 15:04:05"), e.DIPlus, e.DIMins, e.ADX)
}

// NewAverageDirectionalIndex new Func
func NewAverageDirectionalIndex(klineItem *klines.Item, period, length int) *AverageDirectionalIndex {
	m := &AverageDirectionalIndex{
		Name:   fmt.Sprintf("AverageDirectionalIndex%d-%d", period, length),
		Period: period,
		Len:    length,
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

// NewAverageDirectionalIndex new Func
func NewAverageDirectionalIndexOHLC(ohlc *klines.OHLC, period, length int) *AverageDirectionalIndex {
	m := &AverageDirectionalIndex{
		Name:   fmt.Sprintf("AverageDirectionalIndex%d-%d", period, length),
		ohlc:   ohlc,
		Period: period,
		Len:    length,
	}
	return m
}

// NewDefaultAverageDirectionalIndex new Func
func NewDefaultAverageDirectionalIndex(klineItem *klines.Item) *AverageDirectionalIndex {
	return NewAverageDirectionalIndex(klineItem, 30, 14)
}

// calculateTR 计算当前时间段的 True Range (TR)
// TR 是当前高点和当前低点之差、前一收盘价和当前高点之差、前一收盘价和当前低点之差中的最大值
func (e *AverageDirectionalIndex) calculateTR(currentHigh, currentLow, previousClose float64) float64 {
	return math.Max(currentHigh-currentLow, math.Max(math.Abs(currentHigh-previousClose), math.Abs(currentLow-previousClose)))
}

// calculateDM 计算当前时间段的 +DM 和 -DM
// +DM 是当前高点与前一高点之差，如果高点差值大于低点差值并且大于零，则为 +DM，否则为零
// -DM 是当前低点与前一低点之差，如果低点差值大于高点差值并且大于零，则为 -DM，否则为零
func (e *AverageDirectionalIndex) calculateDM(currentHigh, previousHigh, currentLow, previousLow float64) (float64, float64) {
	upMove := currentHigh - previousHigh
	downMove := previousLow - currentLow
	var up, down float64
	if upMove > downMove && upMove > 0 {
		up = upMove
	}
	if downMove > upMove && downMove > 0 {
		down = downMove
	}
	return up, down
}

// smoothValues 对值进行平滑处理
// 使用加权移动平均法对值进行平滑处理，以减少短期波动的影响
func (e *AverageDirectionalIndex) smoothValues(values []float64, period int) []float64 {
	smoothed := make([]float64, len(values))
	for i := 1; i < len(values); i++ {
		smoothed[i] = smoothed[i-1] - smoothed[i-1]/float64(period) + values[i]
	}
	return smoothed
}

// Calculation Func
func (e *AverageDirectionalIndex) Calculation() *AverageDirectionalIndex {

	var high = e.ohlc.High
	var low = e.ohlc.Low
	var close = e.ohlc.Close
	var times = e.ohlc.TimeUnix

	e.data = make([]AverageDirectionalIndexData, len(close))

	tr := make([]float64, len(close))
	plusDM := make([]float64, len(close))
	minusDM := make([]float64, len(close))

	// 计算 TR、+DM 和 -DM
	for i := 1; i < len(close); i++ {
		tr[i] = e.calculateTR(high[i], low[i], close[i-1])
		plusDM[i], minusDM[i] = e.calculateDM(high[i], high[i-1], low[i], low[i-1])
	}

	// 对 TR、+DM 和 -DM 进行平滑处理
	smoothedTR := e.smoothValues(tr, e.Len)
	smoothedPlusDM := e.smoothValues(plusDM, e.Len)
	smoothedMinusDM := e.smoothValues(minusDM, e.Len)

	// 计算 +DI 和 -DI
	plusDI := make([]float64, len(close))
	minusDI := make([]float64, len(close))
	dx := make([]float64, len(close))

	for i := e.Period - 1; i < len(close); i++ {
		plusDI[i] = (smoothedPlusDM[i] / smoothedTR[i]) * 100
		minusDI[i] = (smoothedMinusDM[i] / smoothedTR[i]) * 100
		dx[i] = 100 * (math.Abs(plusDI[i]-minusDI[i]) / (plusDI[i] + minusDI[i]))
	}

	// 对 DX 进行平滑处理
	smoothedDX := ta.Sma(e.Len, dx)
	e.data = make([]AverageDirectionalIndexData, len(close))
	// 计算 ADX
	for i := 0; i < len(close); i++ {
		e.data[i] = AverageDirectionalIndexData{
			Time:   time.Unix(times[i], 0),
			DIPlus: plusDI[i],
			DIMins: minusDI[i],
			ADX:    smoothedDX[i],
		}
	}

	return e
}

// GetData return Point
func (e *AverageDirectionalIndex) GetData() []AverageDirectionalIndexData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// LinearRegressionCandles struct
type LinearRegressionCandles struct {
	SMASignal              bool
	SignalSmoothing        int
	LinearRegressionLength int
	Name                   string
	data                   []LinearRegressionCandlesData
	ohlc                   *klines.OHLC
}

type LinearRegressionCandlesData struct {
	Time   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Signal float64
}

// NewLinearRegressionCandles new Func
func NewLinearRegressionCandles(klineItem *klines.Item, signalSmoothing, linearRegressionLength int, smaSignal bool) *LinearRegressionCandles {
	m := &LinearRegressionCandles{
		Name:                   fmt.Sprintf("LinearRegressionCandles%d-%d", signalSmoothing, linearRegressionLength),
		SignalSmoothing:        signalSmoothing,
		LinearRegressionLength: linearRegressionLength,
		SMASignal:              smaSignal,
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

func NewLinearRegressionCandlesOHLC(ohlc *klines.OHLC, signalSmoothing, linearRegressionLength int, smaSignal bool) *LinearRegressionCandles {
	m := &LinearRegressionCandles{
		Name:                   fmt.Sprintf("LinearRegressionCandles%d-%d", signalSmoothing, linearRegressionLength),
		ohlc:                   ohlc,
		SignalSmoothing:        signalSmoothing,
		LinearRegressionLength: linearRegressionLength,
		SMASignal:              smaSignal,
	}
	return m
}

// NewLinearRegressionCandles new Func
func NewDefaultLinearRegressionCandles(klineItem *klines.Item) *LinearRegressionCandles {
	return NewLinearRegressionCandles(klineItem, 7, 11, true)
}

// Calculation Func
func (e *LinearRegressionCandles) Calculation() *LinearRegressionCandles {

	var ohlc = e.ohlc
	var times = ohlc.Time
	var closes = ohlc.Close
	var bopen = ta.LinearReg(ohlc.Open, e.LinearRegressionLength)
	var bhigh = ta.LinearReg(ohlc.High, e.LinearRegressionLength)
	var blow = ta.LinearReg(ohlc.Low, e.LinearRegressionLength)
	var bclose = ta.LinearReg(closes, e.LinearRegressionLength)

	var signal []float64
	if e.SMASignal {
		signal = ta.Sma(e.SignalSmoothing, bclose)
	} else {
		signal = ta.Ema(e.SignalSmoothing, bclose)
	}
	// ? sma(bclose, signal_length) : ema(bclose, signal_length)

	e.data = make([]LinearRegressionCandlesData, len(closes))

	for i := 0; i < len(closes); i++ {
		e.data[i] = LinearRegressionCandlesData{
			Open:   bopen[i],
			Close:  bclose[i],
			High:   bhigh[i],
			Low:    blow[i],
			Signal: signal[i],
			Time:   times[i],
		}

	}

	defer func() {
		bopen = nil
		bhigh = nil
		blow = nil
		bclose = nil
		signal = nil
	}()
	return e
}

// GetValues Func
func (e *LinearRegressionCandles) GetValues() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var signal = make([]float64, len(e.data))
	for i, v := range e.data {
		signal[i] = v.Signal
	}
	return signal
}

// GetData Func
func (e *LinearRegressionCandles) GetData() []LinearRegressionCandlesData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

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
	kline                  *klines.Item
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
		kline:                  klineItem,
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

	var ohlc = e.kline.GetOHLC()
	var times = ohlc.Time
	var bopen = ta.LinearReg(ohlc.Open, e.LinearRegressionLength)
	var bhigh = ta.LinearReg(ohlc.High, e.LinearRegressionLength)
	var blow = ta.LinearReg(ohlc.Low, e.LinearRegressionLength)
	var bclose = ta.LinearReg(ohlc.Close, e.LinearRegressionLength)

	var signal []float64
	if e.SMASignal {
		signal = ta.Sma(e.SignalSmoothing, bclose)
	} else {
		signal = ta.Ema(e.SignalSmoothing, bclose)
	}
	// ? sma(bclose, signal_length) : ema(bclose, signal_length)

	e.data = make([]LinearRegressionCandlesData, len(e.kline.Candles))

	for i := 0; i < len(e.kline.Candles); i++ {
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
		ohlc = nil
		bopen = nil
		bhigh = nil
		blow = nil
		bclose = nil
		signal = nil
	}()
	return e
}

// GetValues Func
func (e *LinearRegressionCandles) GetValues() (signal []float64) {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	for _, v := range e.data {
		signal = append(signal, v.Signal)
	}
	return
}

// GetData Func
func (e *LinearRegressionCandles) GetData() []LinearRegressionCandlesData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

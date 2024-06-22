package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// UTBot struct
type UTBot struct {
	// Key Vaule. 'This changes the sensitivity
	Period    int
	AtrPeriod int
	Name      string
	ohlc      *klines.OHLC
	data      []UTBotData
}

type UTBotData struct {
	Time  time.Time
	Close float64
	Value float64
}

// NewUTBot new Func
func NewUTBot(klineItem *klines.Item, period, atrPeriod int) *UTBot {
	m := &UTBot{
		Name:      fmt.Sprintf("UTBot%d", period),
		Period:    period,
		AtrPeriod: atrPeriod,
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

// NewUTBot new Func
func NewUTBotOHLC(ohlc *klines.OHLC, period, atrPeriod int) *UTBot {
	m := &UTBot{
		Name:      fmt.Sprintf("UTBot%d", period),
		ohlc:      ohlc,
		Period:    period,
		AtrPeriod: atrPeriod,
	}
	return m
}

// NewUTBot new Func
func NewDefaultUTBot(klineItem *klines.Item) *UTBot {
	return NewUTBot(klineItem, 2, 1)
}

// Calculation Func
func (e *UTBot) Calculation() *UTBot {

	var ohlc = e.ohlc

	var closes = ohlc.Close
	var times = ohlc.Time

	var atr = ta.Atr(ohlc.High, ohlc.Low, closes, e.AtrPeriod)

	var xATRTrailingStop = make([]float64, len(closes))

	defer func() {
		xATRTrailingStop = nil
		atr = nil
	}()

	e.data = make([]UTBotData, len(closes))
	for i, close := range closes {
		if i == 0 {
			continue
		}
		// var prevClose = closeing[i-1]
		var nLoss = float64(e.Period) * atr[i]

		// 计算xATRTrailingStop
		if close > xATRTrailingStop[i-1] && closes[i-1] > xATRTrailingStop[i-1] {
			xATRTrailingStop[i] = math.Max(xATRTrailingStop[i-1], close-nLoss)
		} else if close < xATRTrailingStop[i-1] && closes[i-1] < xATRTrailingStop[i-1] {
			xATRTrailingStop[i] = math.Min(xATRTrailingStop[i-1], close+nLoss)
		} else if close > xATRTrailingStop[i-1] {
			xATRTrailingStop[i] = close - nLoss
		} else {
			xATRTrailingStop[i] = close + nLoss
		}

		// fmt.Printf("[%s]%f\t%f\tatr:%f\txATRTrailingStop:%f\tBuy:%+v\tSell:%+v\n", e.kline.Candles[i].Time.Format("2006-01-02 15:04:05"), close-nLoss, close+nLoss, atr[i], xATRTrailingStop[i], buy, sell)
		e.data[i] = UTBotData{
			Time:  times[i],
			Close: close,
			Value: xATRTrailingStop[i],
		}

	}
	return e
}

// AnalysisSide Func
func (e *UTBot) AnalysisSide() utils.SideData {

	if len(e.data) == 0 {
		e = e.Calculation()
	}
	sides := make([]utils.Side, len(e.data))

	for i, v := range e.data {
		if i < 1 {
			continue
		}
		var close = v.Close
		var closePrev = e.data[i-1].Close
		var val = v.Value
		var valPrev = e.data[i-1].Value

		var above = close > val && closePrev < valPrev
		var below = close < val && closePrev > valPrev

		// 当 DIF、DEA为正，且DIF大于DEA，且DIF向上突破DEA，为买入信号
		if close > val && above {
			sides[i] = utils.Buy
		} else if close < val && below {
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

// GetValues Func
func (e *UTBot) GetValues() (val []float64) {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	for _, v := range e.data {
		val = append(val, v.Value)
	}
	return
}

// GetData Func
func (e *UTBot) GetData() []UTBotData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

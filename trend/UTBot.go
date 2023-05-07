package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// UTBot struct
type UTBot struct {
	// Key Vaule. 'This changes the sensitivity
	Period    int
	AtrPeriod int
	Name      string
	data      []UTBotData
	kline     utils.Klines
}

type UTBotData struct {
	Time  time.Time
	Close float64
	Value float64
}

// NewUTBot new Func
func NewUTBot(list utils.Klines, period, atrPeriod int) *UTBot {
	m := &UTBot{
		Name:      fmt.Sprintf("UTBot%d", period),
		kline:     list,
		Period:    period,
		AtrPeriod: atrPeriod,
	}
	return m
}

// NewUTBot new Func
func NewDefaultUTBot(list utils.Klines) *UTBot {
	return NewUTBot(list, 2, 1)
}

// Calculation Func
func (e *UTBot) Calculation() *UTBot {

	var xATRTrailingStop = make([]float64, len(e.kline))

	var ohlc = e.kline.GetOHLC()

	var atr = (&Atr{}).Atr(ohlc.High, ohlc.Low, ohlc.Close, e.AtrPeriod)
	var closeing = ohlc.Close

	e.data = make([]UTBotData, len(e.kline))
	for i, close := range closeing {
		if i == 0 {
			continue
		}
		// var prevClose = closeing[i-1]
		var nLoss = float64(e.Period) * atr[i]

		// 计算xATRTrailingStop
		if close > xATRTrailingStop[i-1] && closeing[i-1] > xATRTrailingStop[i-1] {
			xATRTrailingStop[i] = math.Max(xATRTrailingStop[i-1], close-nLoss)
		} else if close < xATRTrailingStop[i-1] && closeing[i-1] < xATRTrailingStop[i-1] {
			xATRTrailingStop[i] = math.Min(xATRTrailingStop[i-1], close+nLoss)
		} else if close > xATRTrailingStop[i-1] {
			xATRTrailingStop[i] = close - nLoss
		} else {
			xATRTrailingStop[i] = close + nLoss
		}

		// fmt.Printf("[%s]%f\t%f\tatr:%f\txATRTrailingStop:%f\tBuy:%+v\tSell:%+v\n", e.kline[i].Time.Format("2006-01-02 15:04:05"), close-nLoss, close+nLoss, atr[i], xATRTrailingStop[i], buy, sell)
		e.data[i] = UTBotData{
			Time:  e.kline[i].Time,
			Close: close,
			Value: xATRTrailingStop[i],
		}

	}
	return e
}

// AnalysisSide Func
func (e *UTBot) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

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

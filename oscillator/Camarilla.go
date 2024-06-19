package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Camarilla 卡玛里拉轨道
type Camarilla struct {
	Name      string
	EMAPeriod int
	data      []*CamarillaData
	kline     *klines.Item
}

// CamarillaData.
type CamarillaData struct {
	Time  time.Time
	EMA   float64
	Pivot float64
	High1 float64
	High2 float64
	High3 float64
	High4 float64
	High5 float64
	High6 float64
	Low1  float64
	Low2  float64
	Low3  float64
	Low4  float64
	Low5  float64
	Low6  float64
}

// NewCamarilla new Func
func NewCamarilla(klineItem *klines.Item, emaPeriod int) *Camarilla {
	m := &Camarilla{
		Name:      fmt.Sprintf("Camarilla-%d", emaPeriod),
		kline:     klineItem,
		EMAPeriod: emaPeriod,
	}
	return m
}

// NewDefaultCamarilla new Func
func NewDefaultCamarilla(klineItem *klines.Item) *Camarilla {
	return NewCamarilla(klineItem, 8)
}

// Calculation Func
func (e *Camarilla) Calculation() *Camarilla {

	var ohlc = e.kline.GetOHLC()
	var hlc3 = e.kline.HLC3()
	var closes = ohlc.Close
	var highs = ohlc.High
	var lows = ohlc.Low

	var emas = ta.Ema(e.EMAPeriod, e.kline.GetOHLC().Close)

	var data = make([]*CamarillaData, len(e.kline.Candles))

	for i := 0; i < len(e.kline.Candles); i++ {
		var high = highs[i]
		var low = lows[i]
		var close = closes[i]

		camarillaData := e.getCamarilla(high, low, close)
		camarillaData.Time = e.kline.Candles[i].Time
		camarillaData.EMA = emas[i]
		camarillaData.Pivot = hlc3[i]

		data[i] = camarillaData
	}
	e.data = data
	return e
}

func (e *Camarilla) getCamarilla(high, low, close float64) *CamarillaData {

	var h5 = (high / low) * close
	var h4 = close + (high-low)*1.1/2.0
	var h3 = close + (high-low)*1.1/4.0
	var h2 = close + (high-low)*1.1/6.0
	var h1 = close + (high-low)*1.1/12.0
	var l1 = close - (high-low)*1.1/12.0
	var l2 = close - (high-low)*1.1/6.0
	var l3 = close - (high-low)*1.1/4.0
	var l4 = close - (high-low)*1.1/2.0
	var h6 = h5 + 1.168*(h5-h4)
	var l5 = close - (h5 - close)
	var l6 = close - (h6 - close)

	return &CamarillaData{
		High1: h1,
		High2: h2,
		High3: h3,
		High4: h4,
		High5: h5,
		High6: h6,
		Low1:  l1,
		Low2:  l2,
		Low3:  l3,
		Low4:  l4,
		Low5:  l5,
		Low6:  l6,
	}
}

// AnalysisSide Func
func (e *Camarilla) AnalysisSide(klineItem *klines.Item) utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	// 计算另一个时间维度的K线
	camarillaNewData := NewCamarilla(klineItem, e.EMAPeriod).GetData()

	var securityCamarillaData *CamarillaData

	var closes = e.kline.GetOHLC().Close
	var opens = e.kline.GetOHLC().Open

	for i, v := range e.data {
		for _, camarillaDayDataItem := range camarillaNewData {
			if camarillaDayDataItem.Time.Year() == v.Time.Year() && camarillaDayDataItem.Time.Month() == v.Time.Month() && camarillaDayDataItem.Time.Day() == v.Time.Day() {
				securityCamarillaData = camarillaDayDataItem
				break
			}
		}

		var close = closes[i]
		var open = opens[i]

		// 穿越到零線上方時，會出現看漲金叉。
		if securityCamarillaData != nil && close < v.EMA && close > securityCamarillaData.High4 && open < securityCamarillaData.High4 {
			sides[i] = utils.Buy
		} else if securityCamarillaData != nil && close > v.EMA && close < securityCamarillaData.High4 && open > securityCamarillaData.High4 {
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

// GetData Func
func (e *Camarilla) GetData() []*CamarillaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

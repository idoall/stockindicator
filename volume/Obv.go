package volume

import (
	"time"

	"github.com/idoall/stockindicator/utils/klines"
)

// Obv计算方法：
// 主公式：当日Obv=前一日Obv+今日成交量
// 1.基期Obv值为0，即该股上市的第一天，Obv值为0
// 2.若当日收盘价＞上日收盘价，则当日Obv=前一日Obv＋今日成交量
// 3.若当日收盘价＜上日收盘价，则当日Obv=前一日Obv－今日成交量
// 4.若当日收盘价＝上日收盘价，则当日Obv=前一日Obv

// Obv struct
type Obv struct {
	Name  string
	data  []ObvData
	kline *klines.Item
}

type ObvData struct {
	Value float64
	Time  time.Time
}

// NewObv new Obv
func NewObv(klineItem *klines.Item) *Obv {
	m := &Obv{Name: "Obv", kline: klineItem}
	return m
}

// Calculation Func
func (e *Obv) Calculation() *Obv {

	var ohlc = e.kline.GetOHLC()
	var closes = ohlc.Close
	var volumes = ohlc.Volume
	var times = ohlc.Time

	for i := 0; i < len(e.kline.Candles); i++ {

		var value float64

		//由于Obv的计算方法过于简单化，所以容易受到偶然因素的影响，为了提高Obv的准确性，可以采取多空比率净额法对其进行修正。
		//多空比率净额= [（收盘价－最低价）－（最高价-收盘价）] ÷（ 最高价－最低价）×V
		// value = ((item.Close - item.Low) - (item.High - item.Close)) / (item.High - item.Close) * item.Vol

		if i-1 == -1 {
			value = 0
		} else if closes[i] > closes[i-1] {
			value = e.data[i-1].Value + volumes[i]
		} else if closes[i] < closes[i-1] {
			value = e.data[i-1].Value - volumes[i]
		} else if closes[i] == closes[i-1] {
			value = e.data[i-1].Value
		}
		var p ObvData
		p.Value = value
		p.Time = times[i]
		e.data = append(e.data, p)
	}
	return e
}

// GetPoints Func
func (e *Obv) GetData() []ObvData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

/*
1、计算移动平均值（EMA）
12日EMA的算式为
EMA（12）=前一日EMA（12）×11/13+今日收盘价×2/13
26日EMA的算式为
EMA（26）=前一日EMA（26）×25/27+今日收盘价×2/27
2、计算离差值（DIF）
DIF=今日EMA（12）－今日EMA（26）
3、计算DIF的9日EMA
根据离差值计算其9日的EMA，即离差平均值，是所求的Macd值。为了不与指标原名相混淆，此值又名
DEA或DEM。
今日DEA（Macd）=前一日DEA×8/10+今日DIF×2/10计算出的DIF和DEA的数值均为正值或负值。
用（DIF-DEA）×2即为Macd柱状图。
*/

// Macd is the main object
type Macd struct {
	Name         string
	PeriodShort  int //默认12
	PeriodSignal int //信号长度默认9
	PeriodLong   int //默认26
	data         []MacdData
	kline        *klines.Item
}

type MacdData struct {
	Time time.Time
	DIF  float64
	DEA  float64
	Macd float64
	// Hist float64
}

// NewMacd new Func
// 使用方法，先添加最早日期的数据,最后一条应该是当前日期的数据，结果与 AICoin 对比完全一致
func NewMacd(klineItem *klines.Item, short, signal, long int) *Macd {
	m := &Macd{
		Name:         fmt.Sprintf("Macd%d-%d-%d", short, signal, long),
		PeriodShort:  short,
		PeriodSignal: signal,
		PeriodLong:   long,
		kline:        klineItem,
	}
	return m
}

// NewDefaultMacd new Func
func NewDefaultMacd(klineItem *klines.Item) *Macd {
	return NewMacd(klineItem, 12, 9, 26)
}

// Calculation Func
func (e *Macd) Calculation() *Macd {

	closes := e.kline.GetOHLC().Close
	// 计算DIF
	difs := ta.Subtract(
		ta.Ema(e.PeriodShort, closes),
		ta.Ema(e.PeriodLong, closes),
	)
	// 计算DEA
	deas := ta.Ema(e.PeriodSignal, difs)

	e.data = make([]MacdData, len(e.kline.Candles))
	for i, dif := range difs {
		e.data[i] = MacdData{
			Time: time.Unix(e.kline.Candles[i].TimeUnix, 0),
			DIF:  dif,
			DEA:  deas[i],
			Macd: (dif - deas[i]) * 2,
			// Hist: dif - (dif-deas[i])*2,
		}
	}

	return e
}

// AnalysisSide Func
func (e *Macd) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		// 当 DIF、DEA为正，且DIF大于DEA，且DIF向上突破DEA，为买入信号
		if v.DIF > 0 && v.DEA > 0 && v.DIF > v.DEA && e.data[i-1].DIF < e.data[i-1].DEA {
			sides[i] = utils.Buy
		} else if v.DIF < 0 && v.DEA < 0 && v.DIF < v.DEA && e.data[i-1].DIF > e.data[i-1].DEA {
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

// GetPoints return Data
func (e *Macd) GetData() []MacdData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// GetDifs return DIFS
func (e *Macd) GetDifs() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	var difs = make([]float64, len(e.data))
	for i, v := range e.data {
		difs[i] = v.DIF
	}
	return difs
}

// GetMACDs return MACDS
func (e *Macd) GetMACDs() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	var val = make([]float64, len(e.data))
	for i, v := range e.data {
		val[i] = v.Macd
	}
	return val
}

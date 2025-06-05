package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

/*
使用AiCoin的CCI指标计算商品通道指数。
*/

// Cci struct
type AiCoinCCI struct {
	Name   string
	Period int     //默认计算几天的
	factor float64 //计算系数
	data   []AiCoinCCIData
	kline  *klines.Item
}

type AiCoinCCIData struct {
	Value float64
	Time  time.Time
}

// NewCci new Func
func NewAiCoinCCI(klineItem *klines.Item, period int) *AiCoinCCI {
	m := &AiCoinCCI{
		Name:   fmt.Sprintf("AiCoinCCI-%d", period),
		kline:  klineItem,
		Period: period,
		factor: 0.015,
	}
	return m
}

// NewDefaultAiCoinCCI new Func
func NewDefaultAiCoinCCI(klineItem *klines.Item) *AiCoinCCI {
	return NewAiCoinCCI(klineItem, 20)
}

// ComputeCCI 计算 CCI 序列
func ComputeCCI(close, matp, md []float64) []float64 {
	n := len(close)
	cci := make([]float64, n)

	for i := 0; i < n; i++ {
		if md[i] == 0 {
			cci[i] = 0 // 或者 math.NaN()
			continue
		}
		cci[i] = (close[i] - matp[i]) / (0.015 * md[i])
	}
	return cci
}

// Calculation Func
func (e *AiCoinCCI) Calculation() *AiCoinCCI {

	close := e.kline.GetOHLC().Close

	ccis := ta.AiCoinCCI(close, e.Period, e.factor)

	//计算 Cci
	for i := 0; i < len(ccis); i++ {
		e.data = append(e.data, AiCoinCCIData{
			Time:  time.Unix(e.kline.Candles[i].TimeUnix, 0),
			Value: ccis[i],
		})
	}
	return e
}

// AnalysisSide Func
func (e *AiCoinCCI) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		prevItem := e.data[i-1]
		// APO 上穿零表示看涨，而下穿零表示看跌。
		if v.Value < -100.0 && prevItem.Value > -100 {
			sides[i] = utils.Buy
		} else if v.Value > 100 && prevItem.Value < 100 {
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

// GetPoints return Point
func (e *AiCoinCCI) GetData() []AiCoinCCIData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	return e.data
}

// GetValue return Value
func (e *AiCoinCCI) GetValue() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var result []float64
	for _, v := range e.data {
		result = append(result, v.Value)
	}
	return result
}

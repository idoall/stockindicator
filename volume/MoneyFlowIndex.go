package volume

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// The Money Flow Index (MFI) 资金流量指标,1989年3月由JWellesWilder's首次发表MFI指标的用法。
// MFI指标实际是将RSI加以修改后，演变而来。
// RSI以成交价为计算基础；MFI指标则结合价和量，将其列入综合考虑的范围。可以说，MFI指标是成交量的RSI指标。
//
// 资金流量指数是一种使用交易量和交易价格来判断超买(overbought)和超卖(oversold)的技术指标。
// 该指标计算的是过去一段时间里(一般为14个周期)，资金流入和资金流出的比值并将其归一化到0-100之间。
//
// Raw Money Flow = Typical Price * Volume
// Money Ratio = Positive Money Flow / Negative Money Flow
// Money Flow Index = 100 - (100 / (1 + Money Ratio))
type MoneyFlowIndex struct {
	Name   string
	Period int
	data   []MoneyFlowIndexData
	kline  *klines.Item
}

// MoneyFlowIndexData
type MoneyFlowIndexData struct {
	Time  time.Time
	Value float64
}

// NewMoneyFlowIndex new Func
func NewMoneyFlowIndex(klineItem *klines.Item, period int) *MoneyFlowIndex {
	m := &MoneyFlowIndex{
		Name:   fmt.Sprintf("MoneyFlowIndex%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultMoneyFlowIndex new Func
func NewDefaultMoneyFlowIndex(klineItem *klines.Item) *MoneyFlowIndex {
	return NewMoneyFlowIndex(klineItem, 14)
}

// Calculation Func
func (e *MoneyFlowIndex) Calculation() *MoneyFlowIndex {

	period := e.Period
	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low
	var closing = ohlc.Close
	var volume = ohlc.Volume

	typicalPrice := make([]float64, len(closing))
	for i := 0; i < len(typicalPrice); i++ {
		typicalPrice[i] = (high[i] + low[i] + closing[i]) / float64(3)
	}
	rawMoneyFlow := ta.Multiply(typicalPrice, volume)

	signs := ta.ExtractSign(ta.Diff(rawMoneyFlow, 1))
	moneyFlow := ta.Multiply(signs, rawMoneyFlow)

	positiveMoneyFlow := ta.KeepPositives(moneyFlow)
	negativeMoneyFlow := ta.KeepNegatives(moneyFlow)

	moneyRatio := ta.Divide(
		ta.Sum(period, positiveMoneyFlow),
		ta.Sum(period, ta.MultiplyBy(negativeMoneyFlow, -1)))

	moneyFlowIndex := ta.AddBy(ta.MultiplyBy(ta.Pow(ta.AddBy(moneyRatio, 1), -1), -100), 100)

	for i := 0; i < len(moneyFlowIndex); i++ {
		e.data = append(e.data, MoneyFlowIndexData{
			Time:  time.Unix(e.kline.Candles[i].TimeUnix, 0),
			Value: moneyFlowIndex[i],
		})
	}

	return e
}

// AnalysisSide Func
// 超买/超卖：资金流量指标很少会出现小于10或大于90。当指标小于10时，说明出现了超卖，股价后续会回升，此时可以考虑买入。当指标大于90时，说明出现了超买，股价承压，需要考虑止盈卖出。
// 趋势转向：股票价格仍在上升，但是资金流量指标却向下穿过80线，说明股价向上的趋势快要结束，需要注意股价下跌风险。或者股票价格仍在下跌，但是资金流量指标向上穿过20线，说明股价向下的趋势快要结束，可以考虑买入做多。
// 底/顶背离：股价创新高(当前价格高于近期最高点)，但是资金流量指标却没有创新高，说明股价的新高缺乏支撑，存在下跌风险。
func (e *MoneyFlowIndex) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]

		if v.Value < 10 && prevItem.Value > 10 {
			sides[i] = utils.Buy
		} else if v.Value > 90 && prevItem.Value < 90 {
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
func (e *MoneyFlowIndex) GetData() []MoneyFlowIndexData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

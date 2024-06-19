package oscillator

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// Ichimoku Cloud. 也称为 Ichimoku Kinko Hyo，计算一个多功能指标，定义支撑和阻力，识别趋势方向，衡量动量，并提供交易信号。
// 经过三十多年的研究和测试，Goichi Hosada认为（9,26,52）的周期设置效果最佳。当时，日本的商业时间表中包括星期六，所以9代表一周半（6 + 3天）的时间。数字26和52分别代表一个月和两个月的时间。
// 在加密货币市场中，许多交易者通常将Ichimoku的周期范围设为（9,26,52）到（10,30,60），以此适应7*24小时的市场。甚至还可以直接将周期设置为（20,60,120），以减少错误信号的产生。
//
// Tenkan-sen (Conversion Line) = (9-Period High + 9-Period Low) / 2
// Kijun-sen (Base Line) = (26-Period High + 26-Period Low) / 2
// Senkou Span A (Leading Span A) = (Conversion Line + Base Line) / 2
// Senkou Span B (Leading Span B) = (52-Period High + 52-Period Low) / 2
// Chikou Span (Lagging Span) = Closing plotted 26 days in the past.
//
// Returns conversionLine, baseLine, leadingSpanA, leadingSpanB, laggingSpan
type IchimokuCloud struct {
	Name               string
	ConversionPeriod   int
	LeadingSpanBPeriod int
	LaggingLinePeriod  int
	Smooth             int // 默认一般是1
	data               []IchimokuCloudData
	kline              *klines.Item
}

// IchimokuCloudData
type IchimokuCloudData struct {
	Time           time.Time
	ConversionLine float64
	BaseLine       float64
	LeadingSpanA   float64
	LeadingSpanB   float64
	LaggingLine    float64
}

// NewIchimokuCloud new Func
func NewIchimokuCloud(klineItem *klines.Item, conversionPeriod, leadingSpanBPeriod, laggingLinePeriod int) *IchimokuCloud {
	m := &IchimokuCloud{
		Name:               fmt.Sprintf("IchimokuCloud%d-%d-%d", conversionPeriod, leadingSpanBPeriod, laggingLinePeriod),
		kline:              klineItem,
		ConversionPeriod:   conversionPeriod,
		LeadingSpanBPeriod: leadingSpanBPeriod,
		LaggingLinePeriod:  laggingLinePeriod,
	}
	return m
}

// NewDefaultIchimokuCloud new Func
func NewDefaultIchimokuCloud(klineItem *klines.Item) *IchimokuCloud {
	return NewIchimokuCloud(klineItem, 20, 60, 120)
}

// Calculation Func
func (e *IchimokuCloud) Calculation() *IchimokuCloud {

	conversionPeriod := e.ConversionPeriod
	leadingSpanBPeriod := e.LeadingSpanBPeriod
	laggingLinePeriod := e.LaggingLinePeriod
	var ohlc = e.kline.GetOHLC()
	var high = ohlc.High
	var low = ohlc.Low
	var closing = ohlc.Close

	// 转换线
	conversionLine := ta.DivideBy(ta.Add(ta.Max(9, high), ta.Min(conversionPeriod, low)), float64(2))
	// 基线
	baseLine := ta.DivideBy(ta.Add(ta.Max(26, high), ta.Min(26, low)), float64(2))
	// 先行带A（Senkou Span A）：通过转换线和基线的移动平均值预计未来26日内趋势
	leadingSpanA := ta.DivideBy(ta.Add(conversionLine, baseLine), float64(2))
	// 先行带B（Senkou Span B）：通过52日移动平均值预计未来26日内趋势。
	leadingSpanB := ta.DivideBy(ta.Add(ta.Max(leadingSpanBPeriod, high), ta.Min(leadingSpanBPeriod, low)), float64(2))
	// 迟行带（Chikou Span）：今日收盘价与过去26日中线的差值。
	laggingLine := ta.ShiftRight(laggingLinePeriod, closing)

	for i := 0; i < len(conversionLine); i++ {
		e.data = append(e.data, IchimokuCloudData{
			Time:           e.kline.Candles[i].Time,
			ConversionLine: conversionLine[i],
			BaseLine:       baseLine[i],
			LeadingSpanA:   leadingSpanA[i],
			LeadingSpanB:   leadingSpanB[i],
			LaggingLine:    laggingLine[i],
		})
	}
	return e
}

// AnalysisSide Func
// 先行带A（3）和先行带B（4）之间的空间称为云带（Kumo），该参数是Ichimoku系统中最值得注意的元素。
// 两条先行带能够预计26日内的市场趋势，因此被视为先行指标。
// 另一方面，迟行带（5）是一个滞后指标，体现了过去26日的趋势。
//
// 动量信号
//
//	市场价格高于基线（看涨信号），低于基线（看跌）。
//	TK cross：转换线在基线上方移动（看涨），在基线下方移动（看跌）。
//
// 趋势跟踪信号
//
//	市场价格高于云带（看涨），低于云带（看跌）。
//	云带颜色从红色变为绿色（看涨），从绿色变为红色（看跌）。
//	迟行带高于市场价格（看涨），低于市场价格（看跌）。
func (e *IchimokuCloud) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var ohlc = e.kline.GetOHLC()
	var closes = ohlc.Close

	for i, v := range e.data {
		if i < e.ConversionPeriod {
			continue
		}

		var close = closes[i]
		var prevItem = e.data[i-1]
		var prevClose = closes[i-1]
		// 市场价格高于云带（看涨），低于云带（看跌）。
		if close > v.LeadingSpanA && close > v.LeadingSpanB && prevClose < prevItem.LeadingSpanA && prevClose < prevItem.LeadingSpanB {
			sides[i] = utils.Buy
		} else if close < v.LeadingSpanA && close < v.LeadingSpanB && prevClose > prevItem.LeadingSpanA && prevClose > prevItem.LeadingSpanB {
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
func (e *IchimokuCloud) GetData() []IchimokuCloudData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

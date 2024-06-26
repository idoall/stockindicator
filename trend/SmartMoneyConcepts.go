package trend

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/utils/commonutils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// SmartMoneyConcepts struct
// 简称SMC交易策略
// 特征变化（CHOCH）： 特征变化意味着市场在一段时间内改变了其趋势或订单流。
// 结构突破 (BOS)：  “结构突破”用于描述超过图表上关键支撑位或阻力位的重大价格变动。该水平通常由机构交易员或其他重要投资者预先确定，他们将其设定为市场变动的关键阈值。
type SmartMoneyConcepts struct {

	// Real Time Swing Structure
	SwingLenght int

	// Bars Confirmation，最小是1
	//
	// 用于确认相等高点和相等低点的条数
	EQHEQL_BarsConfirmation int

	// Threshold，1-5之间
	//
	// 用于检测相等的高低的范围（0，1）内的灵敏度阈值\n\n较低的值将返回较少但更相关的结果
	EQHEQL_Threshold int

	EQHEQL_Enable bool

	Name string
	data []SmartMoneyConceptsData
	ohlc *klines.OHLC
	// Order Blocks Number
	OrderBlockNumber  int
	OrderBlockBullish SmartMoneyConceptsDataOrderBlockList
	OrderBlockBearish SmartMoneyConceptsDataOrderBlockList

	// 内部变量
	os1Prev        int
	os2Prev        int
	idxPrevBullish int
	idxPrevBearish int

	// K线到最后时显示压力、阻力位
	StrongHigh SmartMoneyConceptsDataStrongWeak
	WeakHigh   SmartMoneyConceptsDataStrongWeak
	StrongLow  SmartMoneyConceptsDataStrongWeak
	WeakLow    SmartMoneyConceptsDataStrongWeak
}

type SmartMoneyConceptsData struct {
	Time time.Time
	// 顶部短线结构突破
	HighBOSShort float64
	// 顶部的短线特征变化
	HighCHoCHShort float64
	// 顶部长线结构突破
	HighBOSLong float64
	// 顶部的长线特征变化
	HighCHoCHLong float64
	// 底部短线结构突破
	LowBOSShort float64
	// 底部的短线特征变化
	LowChoCHShort float64
	// 底部长线结构突破
	LowBOSLong float64
	// 底部的长线特征变化
	LowChoCHLong float64

	// 双顶、双底区域
	EQH float64
	EQL float64
}

// NewSmartMoneyConcepts new Func
func NewSmartMoneyConcepts(klineItem *klines.Item, swingLenght int, eqheql_BarsConfirmation, eqheql_Threshold int) *SmartMoneyConcepts {
	m := &SmartMoneyConcepts{
		Name:                    fmt.Sprintf("SmartMoneyConcepts%d-%d-%d", swingLenght, eqheql_BarsConfirmation, eqheql_Threshold),
		SwingLenght:             swingLenght,
		EQHEQL_BarsConfirmation: eqheql_BarsConfirmation,
		EQHEQL_Threshold:        eqheql_Threshold,
		// EQHEQL_Enable:           eqheql_enable,
		OrderBlockNumber: 5,
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

// NewSmartMoneyConcepts new Func
func NewSmartMoneyConceptsOHLC(opens, highs, lows, closes []float64, times []int64, swingLenght int, eqheql_BarsConfirmation, eqheql_Threshold int) *SmartMoneyConcepts {
	m := &SmartMoneyConcepts{
		Name:                    fmt.Sprintf("SmartMoneyConcepts%d-%d-%d", swingLenght, eqheql_BarsConfirmation, eqheql_Threshold),
		SwingLenght:             swingLenght,
		EQHEQL_BarsConfirmation: eqheql_BarsConfirmation,
		EQHEQL_Threshold:        eqheql_Threshold,
		// EQHEQL_Enable:           eqheql_enable,
		OrderBlockNumber: 5,
	}
	m.ohlc = &klines.OHLC{
		Open:     opens,
		High:     highs,
		Low:      lows,
		Close:    closes,
		TimeUnix: times,
	}

	return m
}

// NewSmartMoneyConcepts new Func
func NewDefaultSmartMoneyConcepts(klineItem *klines.Item) *SmartMoneyConcepts {
	return NewSmartMoneyConcepts(klineItem, 50, 3, 1)
}

func (e *SmartMoneyConcepts) Clear() {
	e.data = nil
	e.OrderBlockBearish = nil
	e.OrderBlockBullish = nil
}

// Calculation Func
func (e *SmartMoneyConcepts) Calculation() *SmartMoneyConcepts {

	var ohlc = e.ohlc
	var closes = ohlc.Close
	var highs = ohlc.High
	var lows = ohlc.Low
	var times = ohlc.TimeUnix

	var atr = ta.Atr(highs, lows, closes, 200)

	e.data = make([]SmartMoneyConceptsData, len(closes))

	var trend = 0
	var itrend = 0

	var top_y = 0.0
	var top_x = 0
	var btm_y = 0.0
	var btm_x = 0

	var itop_y = 0.0
	var itop_x = 0
	var ibtm_y = 0.0
	var ibtm_x = 0

	var trail_up = 0.0
	var trail_dn = 0.0
	var trail_up_x = 0
	var trail_dn_x = 0

	var top_cross = true
	var btm_cross = true
	var itop_cross = true
	var ibtm_cross = true

	var eq_prev_top = 0.0
	var eq_top_x = 0

	var eq_prev_btm = 0.0
	var eq_btm_x = 0

	defer func() {
		closes = nil
		highs = nil
		lows = nil
	}()

	for i := 1; i < len(closes); i++ {

		var close = closes[i]
		var high = highs[i]
		var low = lows[i]
		var timeUnix = times[i]

		if trail_up == 0.0 {
			trail_up = high
		}

		if trail_dn == 0.0 {
			trail_dn = low
		}

		if i < e.SwingLenght {
			continue
		}

		var top, btm = e.swings(highs[:i+1], lows[:i+1], e.SwingLenght, i, 1)
		var itop, ibtm = e.swings(highs[:i+1], lows[:i+1], 5, i, 2)

		if top != 0.0 {
			top_cross = true
			// var txt_top =
			top_y = top
			top_x = i - e.SwingLenght

			trail_up = top
			trail_up_x = i - e.SwingLenght
		}

		if itop != 0.0 {
			itop_cross = true
			itop_y = itop
			itop_x = i - 5
		}

		trail_up = math.Max(high, trail_up)
		if trail_up == high {
			trail_up_x = i
		}

		if btm != 0.0 {
			btm_cross = true

			btm_y = btm
			btm_x = i - e.SwingLenght

			trail_dn = btm
			trail_dn_x = i - e.SwingLenght
		}

		if ibtm != 0.0 {
			ibtm_cross = true

			ibtm_y = ibtm
			ibtm_x = i - 5
		}

		trail_dn = math.Min(low, trail_dn)
		if trail_dn == low {
			trail_dn_x = i
		}

		//
		// Pivot High BOS/CHoCH
		//
		if ta.CrossOver(closes[:i+1], itop_y) && itop_cross && top_y != itop_y {
			var choch = false

			if itrend < 0 {
				choch = true
			}

			if choch {
				for candleIndex := itop_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].HighCHoCHShort = itop_y
				}
			} else {
				for candleIndex := itop_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].HighBOSShort = itop_y
				}
			}

			itop_cross = false
			itrend = 1

			e.OrderBlockBullish = e.obCoord(false, i, itop_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBullish, 1)

		}

		if ta.CrossOver(closes[:i+1], top_y) && top_cross {
			var choch = false

			if trend < 0 {
				choch = true
			}
			if choch {
				for candleIndex := top_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].HighCHoCHLong = top_y
				}
			} else {
				for candleIndex := top_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].HighBOSLong = top_y
				}
			}

			top_cross = false
			trend = 1

			e.OrderBlockBearish = e.obCoord(false, i, top_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBearish, 2)
		}

		//
		// Pivot LOW BOS/CHoCH
		//
		if ta.CrossUnder(closes[:i+1], ibtm_y) && ibtm_cross && btm_y != ibtm_y {

			var choch = false

			if itrend > 0 {
				choch = true
			}
			if choch {
				for candleIndex := ibtm_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].LowChoCHShort = ibtm_y
				}
			} else {
				for candleIndex := ibtm_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].LowBOSShort = ibtm_y
				}
			}

			ibtm_cross = false
			itrend = -1

			e.OrderBlockBullish = e.obCoord(true, i, ibtm_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBullish, 1)
		}

		if ta.CrossUnder(closes[:i+1], btm_y) && btm_cross {

			var choch = false
			if itrend > 0 {
				choch = true
			}

			if choch {
				for candleIndex := btm_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].LowChoCHLong = btm_y
				}
			} else {
				for candleIndex := btm_x; candleIndex <= i; candleIndex++ {
					e.data[candleIndex].LowBOSLong = btm_y
				}
			}

			btm_cross = false
			trend = -1

			e.OrderBlockBearish = e.obCoord(true, i, btm_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBearish, 2)
		}

		//
		// Order Blocks
		//
		var orderBlockBullishList = e.OrderBlockBullish
		var orderBlockBullishRemovelist SmartMoneyConceptsDataOrderBlockList
		for j, v := range orderBlockBullishList {
			if close < v.Close && !v.IsTop {
				orderBlockBullishRemovelist = append(orderBlockBullishRemovelist, v)
			} else if close > v.High && v.IsTop {
				orderBlockBullishRemovelist = append(orderBlockBullishRemovelist, v)
			}
			if j > e.OrderBlockNumber && !orderBlockBullishRemovelist.Contains(v) {
				orderBlockBullishRemovelist = append(orderBlockBullishRemovelist, v)
			}
		}
		for _, v := range orderBlockBullishRemovelist {
			e.OrderBlockBullish = e.OrderBlockBullish.Remove(v)
		}

		var orderBlockBearishList = e.OrderBlockBearish
		var orderBlockBearishRemovelist SmartMoneyConceptsDataOrderBlockList
		for j, v := range orderBlockBearishList {
			if close < v.Close && !v.IsTop {
				orderBlockBearishRemovelist = append(orderBlockBearishRemovelist, v)
			} else if close > v.High && v.IsTop {
				orderBlockBearishRemovelist = append(orderBlockBearishRemovelist, v)
			}
			if j > e.OrderBlockNumber && !orderBlockBearishRemovelist.Contains(v) {
				orderBlockBearishRemovelist = append(orderBlockBearishRemovelist, v)
			}
		}
		for _, v := range orderBlockBearishRemovelist {
			e.OrderBlockBearish = e.OrderBlockBearish.Remove(v)
		}

		//
		// EQH/EQL
		//
		if e.EQHEQL_Enable {
			var eq_topArr = ta.PivotHigh(highs[:i+1], e.EQHEQL_BarsConfirmation, e.EQHEQL_BarsConfirmation)
			var eq_top = eq_topArr[len(eq_topArr)-1]
			var eq_btmArr = ta.PivotLow(lows[:i+1], e.EQHEQL_BarsConfirmation, e.EQHEQL_BarsConfirmation)
			var eq_btm = eq_btmArr[len(eq_btmArr)-1]

			if eq_top != 0.0 {
				var max = math.Max(eq_top, eq_prev_top)
				var min = math.Min(eq_top, eq_prev_top)
				if max < (min + atr[i]*float64(e.EQHEQL_Threshold)/10.0) {
					// 划EQH线
					for candleIndex := eq_top_x; candleIndex <= i-e.EQHEQL_BarsConfirmation; candleIndex++ {
						e.data[candleIndex].EQH = eq_prev_top
					}
				}
				eq_prev_top = eq_top
				eq_top_x = i - e.EQHEQL_BarsConfirmation
			}

			if eq_btm != 0.0 {
				var max = math.Max(eq_btm, eq_prev_btm)
				var min = math.Min(eq_btm, eq_prev_btm)
				if min > (max - atr[i]*float64(e.EQHEQL_Threshold)/10.0) {
					// 划EQL线
					for candleIndex := eq_btm_x; candleIndex <= i-e.EQHEQL_BarsConfirmation; candleIndex++ {
						e.data[candleIndex].EQL = eq_prev_btm
					}
				}
				eq_prev_btm = eq_btm
				eq_btm_x = i - e.EQHEQL_BarsConfirmation
			}
		}

		//
		// Strong High/Weak High/Strong Low/Weak Low
		//
		if trend < 0 {
			e.StrongHigh = SmartMoneyConceptsDataStrongWeak{
				Time:  time.Unix(times[trail_up_x], 0),
				Value: trail_up,
			}
		} else {
			e.WeakHigh = SmartMoneyConceptsDataStrongWeak{
				Time:  time.Unix(times[trail_up_x], 0),
				Value: trail_up,
			}
		}

		if trend > 0 {
			e.StrongLow = SmartMoneyConceptsDataStrongWeak{
				Time:  time.Unix(times[trail_dn_x], 0),
				Value: trail_dn,
			}
		} else {
			e.WeakLow = SmartMoneyConceptsDataStrongWeak{
				Time:  time.Unix(times[trail_dn_x], 0),
				Value: trail_dn,
			}
		}

		e.data[i].Time = time.Unix(timeUnix, 0)
	}

	return e
}

func (e *SmartMoneyConcepts) obCoord(useMax bool, index, loc int, highs, lows, atr []float64, list SmartMoneyConceptsDataOrderBlockList, obType int) SmartMoneyConceptsDataOrderBlockList {

	var min = 99999999.0
	var max = 0.0
	var idx = 1
	var ob_threshold = atr

	if useMax {
		for i := index; i > (loc - 1); i-- {
			var h = highs[i]
			var l = lows[i]
			if (h - l) < ob_threshold[i]*2 {
				var idxPrev = e.idxPrevBullish
				if obType == 2 {
					idxPrev = e.idxPrevBearish
				}
				max = math.Max(h, max)
				min = commonutils.If(max == h, l, min).(float64)
				idx = commonutils.If(max == h, i, idxPrev).(int)
				if obType == 1 {
					e.idxPrevBullish = idx
				} else if obType == 2 {
					e.idxPrevBearish = idx
				}
			}

		}
	} else {
		for i := index; i > (loc - 1); i-- {
			var h = highs[i]
			var l = lows[i]
			if (h - l) < ob_threshold[i]*2 {
				var idxPrev = e.idxPrevBullish
				if obType == 2 {
					idxPrev = e.idxPrevBearish
				}
				min = math.Min(l, min)
				max = commonutils.If(min == l, h, max).(float64)
				idx = commonutils.If(min == l, i, idxPrev).(int)
				if obType == 1 {
					e.idxPrevBullish = idx
				} else if obType == 2 {
					e.idxPrevBearish = idx
				}
			}
		}
	}

	var times = e.ohlc.TimeUnix
	var closes = e.ohlc.Close
	var opens = e.ohlc.Open
	list = list.Add(SmartMoneyConceptsDataOrderBlock{
		IsTop: useMax,
		Time:  time.Unix(times[idx], 0),
		Open:  opens[idx],
		Close: closes[idx],
		High:  highs[idx],
		Low:   lows[idx],
	})
	return list

}

// swings 摆动检测，测量
func (e *SmartMoneyConcepts) swings(highs, lows []float64, lenght, index int, osType int) (float64, float64) {
	var upper = ta.Highest(highs, lenght)
	var lower = ta.Lowest(lows, lenght)

	var h = highs[len(highs)-1]
	if len(highs)-lenght-1 >= 0 {
		h = highs[len(highs)-lenght-1]
	}
	var l = lows[len(lows)-1]
	if len(lows)-lenght-1 >= 0 {
		l = lows[len(lows)-lenght-1]
	}

	var os, osPrev int
	if osType == 1 {
		osPrev = e.os1Prev
	} else if osType == 2 {
		osPrev = e.os2Prev
	}
	if h > upper {
		os = 0
	} else if l < lower {
		os = 1
	} else {
		os = osPrev
	}

	var top, btm float64
	if os == 0 && osPrev != 0 {
		top = h
	}

	if os == 1 && osPrev != 1 {
		btm = l
	}

	if osType == 1 {
		e.os1Prev = os
	} else if osType == 2 {
		e.os2Prev = os
	}

	return top, btm
}

// GetData Func
func (e *SmartMoneyConcepts) GetData() []SmartMoneyConceptsData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// AnalysisSide Func
// func (e *SmartMoneyConcepts) AnalysisSide() utils.SideData {
// 	sides := make([]utils.Side, len(e.kline.Candles))

// 	if len(e.data) == 0 {
// 		e = e.Calculation()
// 	}

// 	for i, v := range e.data {
// 		if i < 1 {
// 			continue
// 		}

// 		if v.UpTrendBegin != 0 {
// 			sides[i] = utils.Buy
// 		} else if v.DownTrendBegin != 0 {
// 			sides[i] = utils.Sell
// 		} else {
// 			sides[i] = utils.Hold
// 		}
// 	}
// 	return utils.SideData{
// 		Name: e.Name,
// 		Data: sides,
// 	}
// }

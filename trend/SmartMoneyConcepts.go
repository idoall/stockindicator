package trend

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/idoall/stockindicator/utils/commonutils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
	"github.com/shopspring/decimal"
)

// decimalPool 对象池，用于复用 decimal.Decimal 对象，减少 GC 压力
var decimalPool = sync.Pool{
	New: func() interface{} {
		return new(decimal.Decimal)
	},
}

// SmartMoneyConcepts 相关常量
const (
	// 摆动点检测相关常量
	ShortSwingPeriod       = 5  // 短周期摆动点周期（固定5根K线）
	DefaultLongSwingPeriod = 50 // 默认长周期摆动点周期

	// 订单区块相关常量
	DefaultOrderBlockCount  = 5   // 默认保留的订单区块数量
	OrderBlockATRMultiplier = 2.0 // 订单区块过滤的ATR倍数阈值

	// ATR 相关常量
	DefaultATRPeriod = 200 // 默认ATR计算周期

	// FVG 相关常量
	FVGAutoThresholdMultiplier = 2.0 // FVG 动态阈值倍数

	// EQH/EQL 相关常量
	DefaultEQHEQLBarsConfirmation = 3 // 默认确认K线数量
	DefaultEQHEQLThreshold        = 1 // 默认灵敏度阈值
)

// NewSmartMoneyConcepts 创建一个新的 SmartMoneyConcepts 指标实例
//
// 参数:
//   - klineItem: K线数据项，包含 OHLC 数据
//   - swingLenght: 摆动结构的周期长度，推荐值 50
//   - eqheql_BarsConfirmation: 相等高低点确认所需的K线数量，推荐值 3
//   - eqheql_Threshold: 相等高低点检测阈值（1-5），推荐值 1
//
// 返回:
//   - *SmartMoneyConcepts: SMC 指标实例
func NewSmartMoneyConcepts(klineItem *klines.Item, swingLenght int, eqheql_BarsConfirmation, eqheql_Threshold int) *SmartMoneyConcepts {
	m := &SmartMoneyConcepts{
		Name:                    fmt.Sprintf("SmartMoneyConcepts%d-%d-%d", swingLenght, eqheql_BarsConfirmation, eqheql_Threshold),
		SwingLength:             swingLenght,
		EQHEQL_BarsConfirmation: eqheql_BarsConfirmation,
		EQHEQL_Threshold:        eqheql_Threshold,
		// EQHEQL_Enable:           eqheql_enable,
		OrderBlockNumber:  DefaultOrderBlockCount, // 默认保留5个订单区块
		FVG_Enable:        true,                   // 默认启用 FVG 检测
		FVG_AutoThreshold: true,                   // 默认启用自动阈值过滤
		FVG_KeepHistory:   true,                   // 默认保留历史记录
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

// NewSmartMoneyConceptsOHLC 使用原始 OHLC 数据创建 SmartMoneyConcepts 指标实例
//
// 当你已有分离的 OHLC 数组数据时使用此构造函数
//
// 参数:
//   - opens: 开盘价数组
//   - highs: 最高价数组
//   - lows: 最低价数组
//   - closes: 收盘价数组
//   - times: Unix 时间戳数组
//   - swingLenght: 摆动结构的周期长度
//   - eqheql_BarsConfirmation: 相等高低点确认所需的K线数量
//   - eqheql_Threshold: 相等高低点检测阈值（1-5）
//
// 返回:
//   - *SmartMoneyConcepts: SMC 指标实例
func NewSmartMoneyConceptsOHLC(opens, highs, lows, closes []float64, times []int64, swingLenght int, eqheql_BarsConfirmation, eqheql_Threshold int) *SmartMoneyConcepts {
	m := &SmartMoneyConcepts{
		Name:                    fmt.Sprintf("SmartMoneyConcepts%d-%d-%d", swingLenght, eqheql_BarsConfirmation, eqheql_Threshold),
		SwingLength:             swingLenght,
		EQHEQL_BarsConfirmation: eqheql_BarsConfirmation,
		EQHEQL_Threshold:        eqheql_Threshold,
		// EQHEQL_Enable:           eqheql_enable,
		OrderBlockNumber:  DefaultOrderBlockCount, // 默认保留5个订单区块
		FVG_Enable:        true,                   // 默认启用 FVG 检测
		FVG_AutoThreshold: true,                   // 默认启用自动阈值过滤
		FVG_KeepHistory:   true,                   // 默认保留历史记录
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

// NewDefaultSmartMoneyConcepts 使用默认参数创建 SmartMoneyConcepts 指标实例
//
// 默认参数:
//   - swingLenght: 50（较长的摆动周期，适合日线级别）
//   - eqheql_BarsConfirmation: 3（中等确认强度）
//   - eqheql_Threshold: 1（最严格的相等检测）
//
// 参数:
//   - klineItem: K线数据项
//
// 返回:
//   - *SmartMoneyConcepts: 使用默认参数的 SMC 指标实例
func NewDefaultSmartMoneyConcepts(klineItem *klines.Item) *SmartMoneyConcepts {
	return NewSmartMoneyConcepts(klineItem, 50, 3, 1)
}

// Clear 清空所有计算结果和订单区块数据
// 当需要重新计算或释放内存时调用
func (e *SmartMoneyConcepts) Clear() {
	// e.data = nil
	e.OrderBlockBearish = nil
	e.OrderBlockBullish = nil
	e.cachedATR = nil // 清空 ATR 缓存
}

// Calculation 执行 SMC 指标的完整计算
//
// 该方法会遍历所有K线数据，计算以下指标：
//  1. 摆动高点和低点（基于长周期和短周期）
//  2. BOS（结构突破）和 CHoCH（特征变化）信号
//  3. 订单区块（Order Blocks）
//  4. 相等高点/低点（EQH/EQL）
//  5. 强弱高低点（Strong/Weak High/Low）
//
// 返回:
//   - *SmartMoneyConcepts: 返回自身，支持链式调用
func (e *SmartMoneyConcepts) Calculation() *SmartMoneyConcepts {

	var ohlc = e.ohlc
	var closes = ohlc.Close
	var highs = ohlc.High
	var lows = ohlc.Low
	var times = ohlc.TimeUnix

	// 设置默认 ATR 周期
	if e.atrPeriod == 0 {
		e.atrPeriod = DefaultATRPeriod
	}

	// 计算或复用缓存的 ATR（用于订单区块阈值和相等高低点检测）
	var atr []float64
	if e.cachedATR != nil && len(e.cachedATR) == len(closes) {
		// 复用缓存，避免重复计算
		atr = e.cachedATR
	} else {
		// 重新计算并缓存
		atr = ta.Atr(highs, lows, closes, e.atrPeriod)
		e.cachedATR = atr
	}

	// 初始化结果数据数组
	// e.data = make([]SmartMoneyConceptsData, len(closes))

	// ========== 趋势状态变量 ==========
	// trend: 长周期趋势方向（1=上涨，-1=下跌，0=中性）
	var trend = 0
	// itrend: 短周期（内部）趋势方向（1=上涨，-1=下跌，0=中性）
	var itrend = 0

	// ========== 长周期摆动点变量 ==========
	// top_y: 当前摆动高点的价格
	var top_y = 0.0
	// top_x: 当前摆动高点的索引位置
	var top_x = 0
	// btm_y: 当前摆动低点的价格
	var btm_y = 0.0
	// btm_x: 当前摆动低点的索引位置
	var btm_x = 0

	// ========== 短周期（5根K线）摆动点变量 ==========
	// itop_y: 内部摆动高点的价格
	var itop_y = 0.0
	// itop_x: 内部摆动高点的索引位置
	var itop_x = 0
	// ibtm_y: 内部摆动低点的价格
	var ibtm_y = 0.0
	// ibtm_x: 内部摆动低点的索引位置
	var ibtm_x = 0

	// ========== 跟踪高低点变量（用于强弱高低点计算） ==========
	// trail_up: 跟踪的最高价格
	var trail_up = 0.0
	// trail_dn: 跟踪的最低价格
	var trail_dn = 0.0
	// trail_up_x: 跟踪最高价的索引位置
	var trail_up_x = 0
	// trail_dn_x: 跟踪最低价的索引位置
	var trail_dn_x = 0

	// ========== 交叉标志变量（避免重复触发） ==========
	// top_cross: 长周期高点是否可以触发交叉
	var top_cross = true
	// btm_cross: 长周期低点是否可以触发交叉
	var btm_cross = true
	// itop_cross: 短周期高点是否可以触发交叉
	var itop_cross = true
	// ibtm_cross: 短周期低点是否可以触发交叉
	var ibtm_cross = true

	// ========== 相等高低点变量 ==========
	// eq_prev_top: 前一个枢轴高点的价格
	var eq_prev_top = 0.0
	// eq_top_x: 前一个枢轴高点的索引位置
	var eq_top_x = 0

	// eq_prev_btm: 前一个枢轴低点的价格
	var eq_prev_btm = 0.0
	// eq_btm_x: 前一个枢轴低点的索引位置
	var eq_btm_x = 0

	// FVG 相关变量
	var cumulativeDelta = 0.0 // 累积的 K 线百分比变化（用于动态阈值）

	// 预计算循环常量，避免重复访问结构体字段
	swingLength := e.SwingLength
	orderBlockNumber := e.OrderBlockNumber
	eqheqlEnable := e.EQHEQL_Enable
	eqheqlBarsConfirmation := e.EQHEQL_BarsConfirmation
	eqheqlThreshold := e.EQHEQL_Threshold
	fvgEnable := e.FVG_Enable
	strongWeakEnable := e.StrongWeak_Enable

	// 主循环：遍历所有K线数据进行SMC分析
	for i := 1; i < len(closes); i++ {

		var close = closes[i]
		var high = highs[i]
		var low = lows[i]
		// var timeUnix = times[i]

		// 初始化跟踪高点（第一次遇到时）
		if trail_up == 0.0 {
			trail_up = high
		}

		// 初始化跟踪低点（第一次遇到时）
		if trail_dn == 0.0 {
			trail_dn = low
		}

		// 跳过前SwingLength根K线，因为无法计算摆动点
		if i < swingLength {
			continue
		}

		// 计算长周期摆动点（基于SwingLength参数）
		// top: 检测到的摆动高点价格，btm: 检测到的摆动低点价格
		var top, btm = e.swings(highs[:i+1], lows[:i+1], swingLength, 1)

		// 计算短周期摆动点（固定5根K线）
		// itop: 内部摆动高点，ibtm: 内部摆动低点
		var itop, ibtm = e.swings(highs[:i+1], lows[:i+1], ShortSwingPeriod, 2)

		// ========== 处理长周期摆动高点 ==========
		if top != 0.0 {
			// 允许该高点触发交叉信号
			top_cross = true
			// 更新当前摆动高点价格和位置
			top_y = top
			top_x = i - swingLength

			// 更新跟踪高点（用于强弱高点判断）
			trail_up = top
			trail_up_x = i - swingLength
		}

		// ========== 处理短周期摆动高点 ==========
		if itop != 0.0 {
			// 允许该高点触发交叉信号
			itop_cross = true
			// 更新内部摆动高点价格和位置
			itop_y = itop
			itop_x = i - ShortSwingPeriod
		}

		// 持续更新跟踪高点为最高价
		trail_up = math.Max(high, trail_up)
		if trail_up == high {
			trail_up_x = i
		}

		// ========== 处理长周期摆动低点 ==========
		if btm != 0.0 {
			// 允许该低点触发交叉信号
			btm_cross = true

			// 更新当前摆动低点价格和位置
			btm_y = btm
			btm_x = i - swingLength

			// 更新跟踪低点（用于强弱低点判断）
			trail_dn = btm
			trail_dn_x = i - swingLength
		}

		// ========== 处理短周期摆动低点 ==========
		if ibtm != 0.0 {
			// 允许该低点触发交叉信号
			ibtm_cross = true

			// 更新内部摆动低点价格和位置
			ibtm_y = ibtm
			ibtm_x = i - ShortSwingPeriod
		}

		// 持续更新跟踪低点为最低价
		trail_dn = math.Min(low, trail_dn)
		if trail_dn == low {
			trail_dn_x = i
		}

		// ========== 检测短周期摆动高点的突破（BOS 或 CHoCH） ==========
		// 当收盘价向上穿越内部摆动高点时触发
		if ta.CrossOver(closes[:i+1], itop_y) && itop_cross && top_y != itop_y {
			// 处理结构突破信号
			_, signal, shouldRecord := e.handleStructureBreak(
				i, itop_x, itop_y,
				itrend, itrend < 0, // 如果之前是下跌趋势，则为CHoCH
				PeriodTypeShort, PositionTypeHigh,
				closes, times,
			)
			if shouldRecord {
				e.recordStructureBreak(signal)
			}

			// 重置交叉标志，避免重复触发
			itop_cross = false
			// 更新短周期趋势为上涨
			itrend = 1

			// 记录看涨订单区块（突破前的最后下跌区域）
			e.OrderBlockBullish = e.obCoord(false, i, itop_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBullish, 1)
		}

		// ========== 检测长周期摆动高点的突破（BOS 或 CHoCH） ==========
		// 当收盘价向上穿越长周期摆动高点时触发
		if ta.CrossOver(closes[:i+1], top_y) && top_cross {
			// 处理结构突破信号
			_, signal, shouldRecord := e.handleStructureBreak(
				i, top_x, top_y,
				trend, trend < 0, // 如果之前是下跌趋势，则为CHoCH
				PeriodTypeLong, PositionTypeHigh,
				closes, times,
			)
			if shouldRecord {
				e.recordStructureBreak(signal)
			}

			// 重置交叉标志
			top_cross = false
			// 更新长周期趋势为上涨
			trend = 1

			// 记录看跌订单区块（突破前的最后上涨区域，可能成为未来阻力位）
			e.OrderBlockBearish = e.obCoord(false, i, top_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBearish, 2)
		}

		// ========== 检测短周期摆动低点的下破（BOS 或 CHoCH） ==========
		// 当收盘价向下穿越内部摆动低点时触发
		if ta.CrossUnder(closes[:i+1], ibtm_y) && ibtm_cross && btm_y != ibtm_y {
			// 处理结构突破信号
			_, signal, shouldRecord := e.handleStructureBreak(
				i, ibtm_x, ibtm_y,
				itrend, itrend > 0, // 如果之前是上涨趋势，则为CHoCH
				PeriodTypeShort, PositionTypeLow,
				closes, times,
			)
			if shouldRecord {
				e.recordStructureBreak(signal)
			}

			// 重置交叉标志
			ibtm_cross = false
			// 更新短周期趋势为下跌
			itrend = -1

			// 记录看涨订单区块（下破前的最后上涨区域，可能成为未来支撑位）
			e.OrderBlockBullish = e.obCoord(true, i, ibtm_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBullish, 1)
		}

		// ========== 检测长周期摆动低点的下破（BOS 或 CHoCH） ==========
		// 当收盘价向下穿越长周期摆动低点时触发
		if ta.CrossUnder(closes[:i+1], btm_y) && btm_cross {
			// 处理结构突破信号
			_, signal, shouldRecord := e.handleStructureBreak(
				i, btm_x, btm_y,
				trend, itrend > 0, // 如果之前是上涨趋势，则为CHoCH
				PeriodTypeLong, PositionTypeLow,
				closes, times,
			)
			if shouldRecord {
				e.recordStructureBreak(signal)
			}

			// 重置交叉标志
			btm_cross = false
			// 更新长周期趋势为下跌
			trend = -1

			// 记录看跌订单区块（下破前的最后下跌区域，可能成为未来阻力位）
			e.OrderBlockBearish = e.obCoord(true, i, btm_x, highs[:i+1], lows[:i+1], atr[:i+1], e.OrderBlockBearish, 2)
		}

		// ========== 订单区块管理和清理 ==========
		// 使用单次遍历优化，避免 O(n²) 复杂度

		// 从对象池获取 closeDecimal，仅在需要时才初始化
		var closeDecimal *decimal.Decimal
		var closeDecimalFromPool bool

		// 清理看涨订单区块（单次遍历，O(n)）- 使用原地过滤算法
		if len(e.OrderBlockBullish) > 0 {
			if closeDecimal == nil {
				closeDecimal = decimalPool.Get().(*decimal.Decimal)
				*closeDecimal = decimal.NewFromFloat(close)
				closeDecimalFromPool = true
			}
			validCount := 0
			for j := range e.OrderBlockBullish {
				// 保留条件：
				// 1. 价格未跌破 Low（支撑仍然有效）
				// 2. 在保留数量范围内（保留最新的 OrderBlockNumber 个）
				if !closeDecimal.LessThan(e.OrderBlockBullish[j].Low) && j <= orderBlockNumber {
					if validCount != j {
						e.OrderBlockBullish[validCount] = e.OrderBlockBullish[j]
					}
					validCount++
				}
			}
			e.OrderBlockBullish = e.OrderBlockBullish[:validCount]
		}

		// 清理看跌订单区块（单次遍历，O(n)）- 使用原地过滤算法
		if len(e.OrderBlockBearish) > 0 {
			if closeDecimal == nil {
				closeDecimal = decimalPool.Get().(*decimal.Decimal)
				*closeDecimal = decimal.NewFromFloat(close)
				closeDecimalFromPool = true
			}
			validCount := 0
			for j := range e.OrderBlockBearish {
				// 保留条件：
				// 1. 价格未突破 High（阻力仍然有效）
				// 2. 在保留数量范围内（保留最新的 OrderBlockNumber 个）
				if !closeDecimal.GreaterThan(e.OrderBlockBearish[j].High) && j <= orderBlockNumber {
					if validCount != j {
						e.OrderBlockBearish[validCount] = e.OrderBlockBearish[j]
					}
					validCount++
				}
			}
			e.OrderBlockBearish = e.OrderBlockBearish[:validCount]
		}

		// 归还对象到对象池
		if closeDecimalFromPool {
			decimalPool.Put(closeDecimal)
		}

		// ========== 相等高点/低点（EQH/EQL）检测 ==========
		// 识别价格相近的高点或低点，这些区域通常是流动性池
		if eqheqlEnable {
			// 计算枢轴高点（需要左右各EQHEQL_BarsConfirmation根K线确认）
			var eq_topArr = ta.PivotHigh(highs[:i+1], eqheqlBarsConfirmation, eqheqlBarsConfirmation)
			var eq_top = eq_topArr[len(eq_topArr)-1]

			// 计算枢轴低点（需要左右各EQHEQL_BarsConfirmation根K线确认）
			var eq_btmArr = ta.PivotLow(lows[:i+1], eqheqlBarsConfirmation, eqheqlBarsConfirmation)
			var eq_btm = eq_btmArr[len(eq_btmArr)-1]

			// 检测相等高点
			if eq_top != 0.0 {
				var max = math.Max(eq_top, eq_prev_top)
				var min = math.Min(eq_top, eq_prev_top)

				// 判断两个高点是否足够接近（基于ATR阈值）
				// 如果最大值和最小值之间的差距小于ATR*阈值，则认为是相等高点
				if max < (min + atr[i]*float64(eqheqlThreshold)/10.0) {
					// 标记相等高点线段
					// for candleIndex := eq_top_x; candleIndex <= i-e.EQHEQL_BarsConfirmation; candleIndex++ {
					// 	e.data[candleIndex].EQH = eq_prev_top
					// }

					// 如果启用了 EQH/EQL 历史记录，添加到列表
					if eq_prev_top > 0 { // 确保有前一个高点
						priceDiff := math.Abs(eq_top - eq_prev_top)
						priceDiffPercent := priceDiff / eq_prev_top

						newEQH := SmartMoneyConceptsDataEqualHighLow{
							Time:             time.Unix(times[i-eqheqlBarsConfirmation], 0),
							StartIndex:       i - eqheqlBarsConfirmation,        // 存储索引，用于快速计算 Duration
							Price:            decimal.NewFromFloat(eq_prev_top), // 使用相等的价格
							IsHigh:           true,
							FirstPointTime:   time.Unix(times[eq_top_x], 0),
							FirstPointPrice:  decimal.NewFromFloat(eq_prev_top),
							SecondPointTime:  time.Unix(times[i-eqheqlBarsConfirmation], 0),
							SecondPointPrice: decimal.NewFromFloat(eq_top),
							PointDistance:    (i - eqheqlBarsConfirmation) - eq_top_x,
							PriceDifference:  decimal.NewFromFloat(priceDiff),
							PriceDiffPercent: decimal.NewFromFloat(priceDiffPercent),
						}

						e.EQHList = append(e.EQHList, newEQH)
						// e.data[i-e.EQHEQL_BarsConfirmation].NewEQH = newEQH
					}
				}

				// 更新前一个高点的价格和位置
				eq_prev_top = eq_top
				eq_top_x = i - eqheqlBarsConfirmation
			}

			// 检测相等低点
			if eq_btm != 0.0 {
				var max = math.Max(eq_btm, eq_prev_btm)
				var min = math.Min(eq_btm, eq_prev_btm)

				// 判断两个低点是否足够接近（基于ATR阈值）
				// 如果最小值和最大值之间的差距小于ATR*阈值，则认为是相等低点
				if min > (max - atr[i]*float64(e.EQHEQL_Threshold)/10.0) {
					// 标记相等低点线段
					// for candleIndex := eq_btm_x; candleIndex <= i-e.EQHEQL_BarsConfirmation; candleIndex++ {
					// 	e.data[candleIndex].EQL = eq_prev_btm
					// }

					// 如果启用了 EQH/EQL 历史记录，添加到列表
					if eq_prev_btm > 0 { // 确保有前一个低点
						priceDiff := math.Abs(eq_btm - eq_prev_btm)
						priceDiffPercent := priceDiff / eq_prev_btm

						newEQL := SmartMoneyConceptsDataEqualHighLow{
							Time:             time.Unix(times[i-e.EQHEQL_BarsConfirmation], 0),
							StartIndex:       i - e.EQHEQL_BarsConfirmation,     // 存储索引，用于快速计算 Duration
							Price:            decimal.NewFromFloat(eq_prev_btm), // 使用相等的价格
							IsHigh:           false,
							FirstPointTime:   time.Unix(times[eq_btm_x], 0),
							FirstPointPrice:  decimal.NewFromFloat(eq_prev_btm),
							SecondPointTime:  time.Unix(times[i-e.EQHEQL_BarsConfirmation], 0),
							SecondPointPrice: decimal.NewFromFloat(eq_btm),
							PointDistance:    (i - e.EQHEQL_BarsConfirmation) - eq_btm_x,
							PriceDifference:  decimal.NewFromFloat(priceDiff),
							PriceDiffPercent: decimal.NewFromFloat(priceDiffPercent),
						}

						e.EQLList = append(e.EQLList, newEQL)
						// e.data[i-e.EQHEQL_BarsConfirmation].NewEQL = newEQL
					}
				}

				// 更新前一个低点的价格和位置
				eq_prev_btm = eq_btm
				eq_btm_x = i - e.EQHEQL_BarsConfirmation
			}
		}

		// ========== 更新强弱高低点 ==========
		// 根据当前趋势方向判断高低点的强弱属性

		// 处理高点：下跌趋势中的高点为强高点，上涨趋势中的高点为弱高点
		if strongWeakEnable {
			if e.StrongWeakDetail == nil {
				e.StrongWeakDetail = &StrongWeakDetail{}
			}
			if trend < 0 {
				// 下跌趋势：高点是重要的阻力位（强高点）
				e.StrongWeakDetail.StrongHigh = SmartMoneyConceptsDataStrongWeak{
					Time:  time.Unix(times[trail_up_x], 0),
					Value: decimal.NewFromFloat(trail_up),
					Role:  CoinSMCKeyLevelRoleStrongResistance, // 强阻力位
				}
			} else {
				// 上涨趋势：高点可能会被突破（弱高点）
				e.StrongWeakDetail.WeakHigh = SmartMoneyConceptsDataStrongWeak{
					Time:  time.Unix(times[trail_up_x], 0),
					Value: decimal.NewFromFloat(trail_up),
					Role:  CoinSMCKeyLevelRoleWeakResistance, // 弱阻力位
				}
			}

			// 处理低点：上涨趋势中的低点为强低点，下跌趋势中的低点为弱低点
			if trend > 0 {
				// 上涨趋势：低点是重要的支撑位（强低点）
				e.StrongWeakDetail.StrongLow = SmartMoneyConceptsDataStrongWeak{
					Time:  time.Unix(times[trail_dn_x], 0),
					Value: decimal.NewFromFloat(trail_dn),
					Role:  CoinSMCKeyLevelRoleStrongSupport, // 强支撑位
				}
			} else {
				// 下跌趋势：低点可能会被突破（弱低点）
				e.StrongWeakDetail.WeakLow = SmartMoneyConceptsDataStrongWeak{
					Time:  time.Unix(times[trail_dn_x], 0),
					Value: decimal.NewFromFloat(trail_dn),
					Role:  CoinSMCKeyLevelRoleWeakSupport, // 弱支撑位
				}
			}
		}

		// 设置当前K线的时间戳
		// e.data[i].Time = time.Unix(timeUnix, 0)

		// ========== CHoCH/BOS 结构突破处理 ==========
		// 注意：信号在创建时就已经包含了完整的生命周期信息（Duration、FilledTime等）
		// 因此不需要额外的失效检测逻辑
		// if e.StructureBreak_Enable {
		// 	e.checkStructureBreakInvalidation(i, highs, lows, closes, times)
		// }

		// ========== EQH/EQL 相等高低点处理 ==========
		if e.EQHEQL_Enable {
			e.checkEQHEQLInvalidation(i, highs, lows, times)
		}

		// ========== Fair Value Gap 处理 ==========
		if fvgEnable && i >= 2 { // 至少需要 3 根 K 线

			// 1. 检查现有 FVG 是否失效
			e.checkFVGInvalidation(i, highs, lows, times)

			// 2. 检测新的 FVG
			e.detectNewFVG(i, ohlc.Open, highs, lows, closes, times, cumulativeDelta)

			// 3. 更新累积 delta（用于动态阈值）
			if i > 0 && ohlc.Open[i-1] != 0 {
				barDeltaPercent := (closes[i-1] - ohlc.Open[i-1]) / ohlc.Open[i-1]
				cumulativeDelta += math.Abs(barDeltaPercent)
			}
		}

	}

	return e
}

// handleStructureBreak 处理结构突破信号（BOS/CHoCH）的公共逻辑
//
// 该函数提取了4个重复代码块的公共逻辑，用于检测和记录市场结构突破信号。
//
// 参数:
//   - currentIndex: 当前K线索引
//   - pivotIndex: 摆动点K线索引
//   - pivotPrice: 摆动点价���
//   - currentTrend: 当前趋势（1=上涨，-1=下跌）
//   - trendCheck: CHoCH判断条件（当前趋势 < 0 或 > 0）
//   - period: 周期类型（PeriodTypeShort 或 PeriodTypeLong）
//   - position: 位置类型（PositionTypeHigh 或 PositionTypeLow）
//   - closes: 收盘价数组
//   - times: 时间戳数组
//
// 返回:
//   - signalType: 信号类型（SignalTypeCHoCH 或 SignalTypeBOS）
//   - signal: 生成的结构突破信号（如果启用了记录）
//   - shouldRecord: 是否应该记录该信号
func (e *SmartMoneyConcepts) handleStructureBreak(
	currentIndex, pivotIndex int,
	pivotPrice float64,
	currentTrend int,
	trendCheck bool,
	period PeriodType,
	position PositionType,
	closes []float64,
	times []int64,
) (SignalType, SmartMoneyConceptsDataStructureBreak, bool) {

	// 判断是否为CHoCH（趋势反转）
	var signalType SignalType
	if trendCheck {
		signalType = SignalTypeCHoCH
	} else {
		signalType = SignalTypeBOS
	}

	// 如果未启用结构突破记录，直接返回
	if !e.StructureBreak_Enable {
		return signalType, SmartMoneyConceptsDataStructureBreak{}, false
	}

	// 创建结构突破信号
	pivotTime := time.Unix(times[pivotIndex], 0)
	signal := SmartMoneyConceptsDataStructureBreak{
		Time:       pivotTime, // Time = PivotTime，信号产生时间就是摆动点时间
		Type:       signalType,
		Period:     period,
		Position:   position,
		BreakPrice: decimal.NewFromFloat(pivotPrice),
		ClosePrice: decimal.NewFromFloat(closes[currentIndex]),
		BreakTime:  time.Unix(times[currentIndex], 0), // 实际突破时间
		PivotTime:  pivotTime,
		PivotPrice: decimal.NewFromFloat(pivotPrice),
		Duration:   currentIndex - pivotIndex + 1, // K线标记区间长度
	}

	return signalType, signal, true
}

// recordStructureBreak 将结构突破信号记录到对应的列表中
//
// 根据信号的周期、位置和类型，将信号添加到对应的列表中。
//
// 参数:
//   - signal: 结构突破信号
func (e *SmartMoneyConcepts) recordStructureBreak(signal SmartMoneyConceptsDataStructureBreak) {
	// 根据位置、类型和周期选择对应的列表
	if signal.Position == PositionTypeHigh {
		if signal.Type == SignalTypeCHoCH {
			if signal.Period == PeriodTypeShort {
				e.HighCHoCHShortList = append(e.HighCHoCHShortList, signal)
			} else {
				e.HighCHoCHLongList = append(e.HighCHoCHLongList, signal)
			}
		} else {
			if signal.Period == PeriodTypeShort {
				e.HighBOSShortList = append(e.HighBOSShortList, signal)
			} else {
				e.HighBOSLongList = append(e.HighBOSLongList, signal)
			}
		}
	} else {
		if signal.Type == SignalTypeCHoCH {
			if signal.Period == PeriodTypeShort {
				e.LowCHoCHShortList = append(e.LowCHoCHShortList, signal)
			} else {
				e.LowCHoCHLongList = append(e.LowCHoCHLongList, signal)
			}
		} else {
			if signal.Period == PeriodTypeShort {
				e.LowBOSShortList = append(e.LowBOSShortList, signal)
			} else {
				e.LowBOSLongList = append(e.LowBOSLongList, signal)
			}
		}
	}
}

// obCoord 识别并记录订单区块（Order Block）的坐标和范围
//
// 订单区块是指机构在市场反转前最后一次下单的价格区域。
// 该函数在检测到市场结构突破（BOS/CHoCH）时被调用，用于定位订单区块的位置。
//
// 参数:
//   - useMax: true=寻找最高价区域（看涨OB），false=寻找最低价区域（看跌OB）
//   - index: 当前K线索引位置
//   - loc: 摆动点位置（突破开始的位置）
//   - highs: 最高价数组
//   - lows: 最低价数组
//   - atr: ATR数组，用作波动率阈值
//   - list: 现有的订单区块列表
//   - obType: 订单区块类型（1=看涨，2=看跌）
//
// 返回:
//   - SmartMoneyConceptsDataOrderBlockList: 更新后的订单区块列表
func (e *SmartMoneyConcepts) obCoord(useMax bool, index, loc int, highs, lows, atr []float64, list SmartMoneyConceptsDataOrderBlockList, obType int) SmartMoneyConceptsDataOrderBlockList {

	var min = 99999999.0   // 初始化为极大值
	var max = 0.0          // 初始化为极小值
	var idx = 1            // 订单区块K线的索引位置
	var ob_threshold = atr // ATR阈值，用于过滤波动过大的K线

	// 寻找最高价区域的订单区块（用于看涨突破后识别支撑）
	if useMax {
		// 从当前位置向后扫描到摆动点位置
		for i := index; i > (loc - 1); i-- {
			var h = highs[i]
			var l = lows[i]

			// 仅考虑波动范围小于2倍ATR的K线（排除异常波动）
			if (h - l) < ob_threshold[i]*OrderBlockATRMultiplier {
				var idxPrev = e.idxPrevBullish
				if obType == 2 {
					idxPrev = e.idxPrevBearish
				}

				// 寻找该区间内的最高价
				max = math.Max(h, max)
				// 如果当前K线创造了新高，记录其低点作为OB下边界
				min = commonutils.If(max == h, l, min).(float64)
				// 如果当前K线创造了新高，记录其索引
				idx = commonutils.If(max == h, i, idxPrev).(int)

				// 更新对应类型的OB索引
				switch obType {
				case 1:
					e.idxPrevBullish = idx
				case 2:
					e.idxPrevBearish = idx
				}
			}
		}
	} else {
		// 寻找最低价区域的订单区块（用于看跌突破后识别阻力）
		for i := index; i > (loc - 1); i-- {
			var h = highs[i]
			var l = lows[i]

			// 仅考虑波动范围小于2倍ATR的K线（排除异常波动）
			if (h - l) < ob_threshold[i]*OrderBlockATRMultiplier {
				var idxPrev = e.idxPrevBullish
				if obType == 2 {
					idxPrev = e.idxPrevBearish
				}

				// 寻找该区间内的最低价
				min = math.Min(l, min)
				// 如果当前K线创造了新低，记录其高点作为OB上边界
				max = commonutils.If(min == l, h, max).(float64)
				// 如果当前K线创造了新低，记录其索引
				idx = commonutils.If(min == l, i, idxPrev).(int)

				// 更新对应类型的OB索引
				switch obType {
				case 1:
					e.idxPrevBullish = idx
				case 2:
					e.idxPrevBearish = idx
				}

			}
		}
	}

	// 获取OHLC数据
	var times = e.ohlc.TimeUnix
	var closes = e.ohlc.Close
	var opens = e.ohlc.Open

	// 创建新的订单区块并添加到列表
	list = list.Add(SmartMoneyConceptsDataOrderBlock{
		IsTop: useMax,                            // 是否是顶部区块
		Time:  time.Unix(times[idx], 0),          // 订单区块的时间
		Open:  decimal.NewFromFloat(opens[idx]),  // 开盘价
		Close: decimal.NewFromFloat(closes[idx]), // 收盘价
		High:  decimal.NewFromFloat(highs[idx]),  // 最高价（阻力位）
		Low:   decimal.NewFromFloat(lows[idx]),   // 最低价（支撑位）
	})
	return list
}

// swings 摆动点检测函数
//
// 该函数使用振荡器方法检测价格的摆动高点和摆动低点。
// 摆动高点是指价格在回溯期内创建新高后开始回落的点。
// 摆动低点是指价格在回溯期内创建新低后开始反弹的点。
//
// 工作原理：
// 1. 计算指定周期内的最高价和最低价
// 2. 使用振荡器状态（os）跟踪当前市场方向：
//   - os = 0: 看涨状态（价格创新高）
//   - os = 1: 看跌状态（价格创新低）
//
// 3. 当振荡器状态改变时，记录摆动点
//
// 参数:
//   - highs: 最高价数组
//   - lows: 最低价数组
//   - length: 回溯周期长度
//   - index: 当前K线索引（用于调试，实际未使用）
//   - osType: 振荡器类型（1=长周期，2=短周期）
//
// 返回:
//   - top: 检测到的摆动高点价格（0表示未检测到）
//   - btm: 检测到的摆动低点价格（0表示未检测到）
func (e *SmartMoneyConcepts) swings(highs, lows []float64, length, osType int) (float64, float64) {
	// 计算回溯期内的最高价和最低价
	var upper = ta.Highest(highs, length)
	var lower = ta.Lowest(lows, length)

	// 获取length根K线之前的高点和低点
	// 这是用来判断是否形成摆动点的参考价格
	var h = highs[len(highs)-1]
	if len(highs)-length-1 >= 0 {
		h = highs[len(highs)-length-1]
	}
	var l = lows[len(lows)-1]
	if len(lows)-length-1 >= 0 {
		l = lows[len(lows)-length-1]
	}

	// 获取前一个振荡器状态
	var os, osPrev int
	switch osType {
	case 1:
		osPrev = e.os1Prev // 长周期振荡器状态
	case 2:
		osPrev = e.os2Prev // 短周期振荡器状态
	}

	// 更新当前振荡器状态
	if h > upper {
		os = 0 // 价格创新高，进入看涨状态
	} else if l < lower {
		os = 1 // 价格创新低，进入看跌状态
	} else {
		os = osPrev // 保持之前的状态
	}

	// 检测摆动点：当振荡器状态发生改变时
	var top, btm float64

	// 从非看涨状态转为看涨状态时，前一个高点即为摆动高点
	if os == 0 && osPrev != 0 {
		top = h
	}

	// 从非看跌状态转为看跌状态时，前一个低点即为摆动低点
	if os == 1 && osPrev != 1 {
		btm = l
	}

	// 保存当前振荡器状态供下次使用
	switch osType {
	case 1:
		e.os1Prev = os
	case 2:
		e.os2Prev = os
	}

	return top, btm
}

// checkFVGInvalidation 检查并清除失效的 FVG（遍历所有活跃 FVG）
//
// 失效条件：
//   - 看涨 FVG: 当前 K 线最低价 < FVG 底部（价格向下回填缺口）
//   - 看跌 FVG: 当前 K 线最高价 > FVG 顶部（价格向上回填缺口）
//
// 当 FVG 失效时，会记录失效信息（填补时间、填补价格、持续时间）
// 并将 FVG 添加到历史记录中（如果启用了历史记录功能）
//
// 参数:
//   - index: 当前 K 线索引
//   - highs: 最高价数组
//   - lows: 最低价数组
//   - times: 时间戳数组
func (e *SmartMoneyConcepts) checkFVGInvalidation(index int, highs, lows []float64, times []int64) {
	// 从对象池获取 decimal 对象
	lowsDecimal := decimalPool.Get().(*decimal.Decimal)
	highDecimal := decimalPool.Get().(*decimal.Decimal)
	defer func() {
		decimalPool.Put(lowsDecimal)
		decimalPool.Put(highDecimal)
	}()

	*lowsDecimal = decimal.NewFromFloat(lows[index])
	*highDecimal = decimal.NewFromFloat(highs[index])

	// 检查所有看涨 FVG 是否失效 - 使用原地过滤算法
	validCount := 0
	for i := range e.BullishFVGList {
		fvg := &e.BullishFVGList[i]

		// 更新持续时间（直接使用 StartIndex，避免 O(n) 查找）
		fvg.Duration = index - fvg.StartIndex

		if lowsDecimal.LessThan(fvg.Bottom) {
			// 价格跌破 FVG 底部，FVG 失效

			// 记录失效信息
			fvg.FilledTime = time.Unix(times[index], 0)
			fvg.FilledPrice = *lowsDecimal

			// 如果启用了历史记录，添加到历史列表
			if e.FVG_KeepHistory {
				e.FVG_History = append(e.FVG_History, *fvg)
			}
			// 不保留到有效列表
		} else {
			// FVG 仍然有效，保留在列表中
			if validCount != i {
				e.BullishFVGList[validCount] = *fvg
			}
			validCount++
		}
	}
	e.BullishFVGList = e.BullishFVGList[:validCount]

	// 检查所有看跌 FVG 是否失效 - 使用原地过滤算法
	validCount = 0
	for i := range e.BearishFVGList {
		fvg := &e.BearishFVGList[i]

		// 更新持续时间（直接使用 StartIndex，避免 O(n) 查找）
		fvg.Duration = index - fvg.StartIndex

		if highDecimal.GreaterThan(fvg.Top) {
			// 价格突破 FVG 顶部，FVG 失效

			// 记录失效信息
			fvg.FilledTime = time.Unix(times[index], 0)
			fvg.FilledPrice = *highDecimal

			// 如果启用了历史记录，添加到历史列表
			if e.FVG_KeepHistory {
				e.FVG_History = append(e.FVG_History, *fvg)
			}
			// 不保留到有效列表
		} else {
			// FVG 仍然有效，保留在列表中
			if validCount != i {
				e.BearishFVGList[validCount] = *fvg
			}
			validCount++
		}
	}
	e.BearishFVGList = e.BearishFVGList[:validCount]
}

// detectNewFVG 检测新的 Fair Value Gap 并添加到列表
//
// FVG 形成条件（三根 K 线模式）：
//
//	看涨 FVG: 当前 K 线低点 > 两根前 K 线高点（存在向上缺口）
//	         且上一根 K 线收盘价 > 两根前 K 线高点（确认方向）
//	看跌 FVG: 当前 K 线高点 < 两根前 K 线低点（存在向下缺口）
//	         且上一根 K 线收盘价 < 两根前 K 线低点（确认方向）
//
// 参数:
//   - index: 当前 K 线索引
//   - opens: 开盘价数组
//   - highs: 最高价数组
//   - lows: 最低价数组
//   - closes: 收盘价数组
//   - times: 时间戳数组
//   - cumulativeDelta: 累积的 K 线百分比变化（用于阈值计算）
func (e *SmartMoneyConcepts) detectNewFVG(index int, opens, highs, lows, closes []float64, times []int64, cumulativeDelta float64) {
	// 获取三根 K 线的数据
	last2High := highs[index-2]  // 两根前的高点
	last2Low := lows[index-2]    // 两根前的低点
	lastClose := closes[index-1] // 上一根的收盘价
	lastOpen := opens[index-1]   // 上一根的开盘价
	currentHigh := highs[index]  // 当前高点
	currentLow := lows[index]    // 当前低点

	// 计算上一根 K 线的涨跌幅百分比
	barDeltaPercent := 0.0
	if lastOpen != 0 {
		barDeltaPercent = (lastClose - lastOpen) / lastOpen
	}

	// 计算动态阈值
	threshold := 0.0
	if e.FVG_AutoThreshold && index > 0 {
		// 阈值 = 平均 K 线变化的 2 倍
		threshold = (cumulativeDelta / float64(index)) * FVGAutoThresholdMultiplier
	}

	// 检测看涨 FVG
	if currentLow > last2High && // 存在向上缺口
		lastClose > last2High && // 收盘价确认方向
		barDeltaPercent > threshold { // 超过动态阈值

		// 创建新的看涨 FVG
		// 注意: FVG 由三根 K 线形成,记录中间那根 K 线(index-1)的时间
		newFVG := SmartMoneyConceptsDataFairValueGap{
			Time:       time.Unix(times[index-1], 0),     // 中间 K 线的时间
			StartIndex: index - 1,                        // 存储索引，用于快速计算 Duration
			Top:        decimal.NewFromFloat(currentLow), // FVG 上边界 = 当前低点
			Bottom:     decimal.NewFromFloat(last2High),  // FVG 下边界 = 两根前高点
			IsBullish:  true,
			Duration:   0, // 初始持续时间为0，从index开始计算
		}

		// 添加到看涨 FVG 列表（添加到列表头部，类似 Pine Script 的 unshift）
		e.BullishFVGList = append([]SmartMoneyConceptsDataFairValueGap{newFVG}, e.BullishFVGList...)

		// 在当前 K 线数据中标记新产生的 FVG
		// e.data[index].NewBullishFVG = newFVG
	}

	// 检测看跌 FVG
	if currentHigh < last2Low && // 存在向下缺口
		lastClose < last2Low && // 收盘价确认方向
		-barDeltaPercent > threshold { // 超过动态阈值（注意负号）

		// 创建新的看跌 FVG
		// 注意: FVG 由三根 K 线形成,记录中间那根 K 线(index-1)��时间
		newFVG := SmartMoneyConceptsDataFairValueGap{
			Time:       time.Unix(times[index-1], 0),      // 中间 K 线的时间
			StartIndex: index - 1,                         // 存储索引，用于快速计算 Duration
			Top:        decimal.NewFromFloat(last2Low),    // FVG 上边界 = 两根前低点
			Bottom:     decimal.NewFromFloat(currentHigh), // FVG 下边界 = 当前高点
			IsBullish:  false,
			Duration:   0, // 初始持续时间为0，从index开始计算
		}

		// 添加到看跌 FVG 列表（添加到列表头部，类似 Pine Script 的 unshift）
		e.BearishFVGList = append([]SmartMoneyConceptsDataFairValueGap{newFVG}, e.BearishFVGList...)

		// 在当前 K 线数据中标记新产生的 FVG
		// e.data[index].NewBearishFVG = newFVG
	}
}

// checkEQHEQLInvalidation 检查并更新所有 EQH/EQL 信号的状态
//
// 功能说明：
//  1. 遍历所有 EQH 和 EQL 列表
//  2. 更新每个信号的持续时间
//  3. 检查失效条件并标记失效信号
//  4. 将失效的信号移至历史记录（如果启用）
//
// 失效条件：
//   - EQH（相等高点）: 价格向上突破该水平（high > Price）
//   - EQL（相等低点）: 价格向下突破该水平（low < Price）
//
// 参数:
//   - index: 当前 K 线索引
//   - highs: 最高价数组
//   - lows: 最低价数组
//   - times: 时间戳数组
func (e *SmartMoneyConcepts) checkEQHEQLInvalidation(index int, highs, lows []float64, times []int64) {
	// 从对象池获取 decimal 对象
	currentHigh := decimalPool.Get().(*decimal.Decimal)
	currentLow := decimalPool.Get().(*decimal.Decimal)
	defer func() {
		decimalPool.Put(currentHigh)
		decimalPool.Put(currentLow)
	}()

	*currentHigh = decimal.NewFromFloat(highs[index])
	*currentLow = decimal.NewFromFloat(lows[index])
	currentTime := time.Unix(times[index], 0)

	// 检查所有 EQH（相等高点）是否失效 - 使用原地过滤算法
	validCount := 0
	for i := range e.EQHList {
		eqh := &e.EQHList[i]

		// 更新持续时间（直接使用 StartIndex，避免 O(n) 查找）
		eqh.Duration = index - eqh.StartIndex

		// 检查失效条件：价格向上突破相等高点
		if currentHigh.GreaterThan(eqh.Price) {
			eqh.FilledTime = currentTime
			eqh.FilledPrice = *currentHigh
			eqh.BreakDirection = "Upward"

			// 添加到历史记录
			if e.EQHEQL_KeepHistory {
				e.EQHEQL_History = append(e.EQHEQL_History, *eqh)
			}
			// 不保留到有效列表
		} else {
			// EQH 仍然有效，更新后保留
			if validCount != i {
				e.EQHList[validCount] = *eqh
			}
			validCount++
		}
	}
	e.EQHList = e.EQHList[:validCount]

	// 检查所有 EQL（相等低点）是否失效 - 使用原地过滤算法
	validCount = 0
	for i := range e.EQLList {
		eql := &e.EQLList[i]

		// 更新持续时间（直接使用 StartIndex，避免 O(n) 查找）
		eql.Duration = index - eql.StartIndex

		// 检查失效条件：价格向下突破相等低点
		if currentLow.LessThan(eql.Price) {
			eql.FilledTime = currentTime
			eql.FilledPrice = *currentLow
			eql.BreakDirection = "Downward"

			// 添加到历史记录
			if e.EQHEQL_KeepHistory {
				e.EQHEQL_History = append(e.EQHEQL_History, *eql)
			}
			// 不保留到有效列表
		} else {
			// EQL 仍然有效，更新后保留
			if validCount != i {
				e.EQLList[validCount] = *eql
			}
			validCount++
		}
	}
	e.EQLList = e.EQLList[:validCount]
}

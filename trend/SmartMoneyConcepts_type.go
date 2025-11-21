package trend

import (
	"time"
)

// 强弱高低点
type SmartMoneyConceptsDataStrongWeak struct {
	Time  time.Time
	Value float64
}

// 订单区块管理
type SmartMoneyConceptsDataOrderBlock struct {
	IsTop bool
	Time  time.Time
	Close float64
	High  float64
	Low   float64
	Open  float64
}

// SmartMoneyConceptsDataFairValueGap Fair Value Gap（公允价值缺口）数据结构
//
// FVG 是三根 K 线之间形成的价格缺口，代表市场失衡区域。
// 看涨 FVG：当前 K 线低点 > 两根前 K 线高点（存在向上缺口）
// 看跌 FVG：当前 K 线高点 < 两根前 K 线低点（存在向下缺口）
type SmartMoneyConceptsDataFairValueGap struct {
	// Time FVG 产生的时间
	Time time.Time

	// Top Fair Value Gap 的上边界价格
	Top float64

	// Bottom Fair Value Gap 的下边界价格
	Bottom float64

	// IsBullish 是否为看涨 FVG
	// true: 看涨 FVG（价格可能回调至此获得支撑）
	// false: 看跌 FVG（价格可能反弹至此遇到阻力）
	IsBullish bool

	// FilledTime FVG 被填补的时间
	// 零值表示 FVG 尚未被填补
	FilledTime time.Time

	// FilledPrice FVG 被填补时的价格
	// 看涨 FVG：价格跌破 Bottom 时的 low 值
	// 看跌 FVG：价格突破 Top 时的 high 值
	FilledPrice float64

	// Duration FVG 持续的 K 线数量
	// 从产生到被填补之间经过的 K 线数量
	Duration int
}

// IsValid 判断 FVG 是否有效（非零值）
//
// 返回:
//   - bool: true 表示 FVG 有效，false 表示 FVG 为空
func (f SmartMoneyConceptsDataFairValueGap) IsValid() bool {
	return f.Top != 0 && f.Bottom != 0
}

// Clear 清空 FVG 数据，将所有字段重置为零值
func (f *SmartMoneyConceptsDataFairValueGap) Clear() {
	f.Time = time.Time{}
	f.Top = 0
	f.Bottom = 0
	f.IsBullish = false
	f.FilledTime = time.Time{}
	f.FilledPrice = 0
	f.Duration = 0
}

// SmartMoneyConceptsData 每个K线的 SMC 指标数据
// 存储该K线时刻的市场结构信息，包括 BOS、CHoCH 和 EQH/EQL 等信号
type SmartMoneyConceptsData struct {
	// Time K线时间戳
	Time time.Time

	// ========== 高点相关信号（阻力位突破和特征变化） ==========

	// HighBOSShort 顶部短线结构突破价位
	// 基于短周期（5根K线）计算的高点突破信号
	// 非零值表示在该价位发生了看涨方向的结构突破
	HighBOSShort float64

	// HighCHoCHShort 顶部短线特征变化价位
	// 基于短周期（5根K线）计算的高点特征变化信号
	// 非零值表示市场可能从下跌趋势转为上涨趋势
	HighCHoCHShort float64

	// HighBOSLong 顶部长线结构突破价位
	// 基于长周期（SwingLenght）计算的高点突破信号
	// 表示更强的趋势延续信号
	HighBOSLong float64

	// HighCHoCHLong 顶部长线特征变化价位
	// 基于长周期（SwingLenght）计算的高点特征变化信号
	// 表示更明确的趋势反转信号
	HighCHoCHLong float64

	// ========== 低点相关信号（支撑位突破和特征变化） ==========

	// LowBOSShort 底部短线结构突破价位
	// 基于短周期（5根K线）计算的低点突破信号
	// 非零值表示在该价位发生了看跌方向的结构突破
	LowBOSShort float64

	// LowChoCHShort 底部短线特征变化价位
	// 基于短周期（5根K线）计算的低点特征变化信号
	// 非零值表示市场可能从上涨趋势转为下跌趋势
	LowChoCHShort float64

	// LowBOSLong 底部长线结构突破价位
	// 基于长周期（SwingLenght）计算的低点突破信号
	// 表示更强的趋势延续信号
	LowBOSLong float64

	// LowChoCHLong 底部长线特征变化价位
	// 基于长周期（SwingLenght）计算的低点特征变化信号
	// 表示更明确的趋势反转信号
	LowChoCHLong float64

	// ========== 相等高低点（流动性区域） ==========

	// EQH 相等高点价位（Equal Highs）
	// 非零值表示检测到两个或多个相近的高点
	// 这些区域通常是机构扫单的目标，可能发生快速突破或反转
	EQH float64

	// EQL 相等低点价位（Equal Lows）
	// 非零值表示检测到两个或多个相近的低点
	// 这些区域通常是机构扫单的目标，可能发生快速突破或反转
	EQL float64

	// ========== Fair Value Gaps（公允价值缺口）==========

	// BullishFVG 当前活跃的看涨 FVG
	// 当有看涨 FVG 存在时，该字段持续有值，直到 FVG 被填补后清零
	// 可用于判断当前是否存在潜在的支撑区域
	BullishFVG SmartMoneyConceptsDataFairValueGap

	// BearishFVG 当前活跃的看跌 FVG
	// 当有看跌 FVG 存在时，该字段持续有值，直到 FVG 被填补后清零
	// 可用于判断当前是否存在潜在的阻力区域
	BearishFVG SmartMoneyConceptsDataFairValueGap

	// NewBullishFVG 本根 K 线新产生的看涨 FVG
	// 仅在 FVG 产生的那根 K 线有值，其他时候为零值
	// 可用于捕捉 FVG 产生事件
	NewBullishFVG SmartMoneyConceptsDataFairValueGap

	// NewBearishFVG 本根 K 线新产生的看跌 FVG
	// 仅在 FVG 产生的那根 K 线有值，其他时候为零值
	// 可用于捕捉 FVG 产生事件
	NewBearishFVG SmartMoneyConceptsDataFairValueGap

	// FilledBullishFVG 本根 K 线被填补的看涨 FVG
	// 仅在 FVG 失效的那根 K 线有值，记录被填补的 FVG 完整信息（包括持续时间等）
	// 可用于捕捉 FVG 失效事件和统计分析
	FilledBullishFVG SmartMoneyConceptsDataFairValueGap

	// FilledBearishFVG 本根 K 线被填补的看跌 FVG
	// 仅在 FVG 失效的那根 K 线有值，记录被填补的 FVG 完整信息（包括持续时间等）
	// 可用于捕捉 FVG 失效事件和统计分析
	FilledBearishFVG SmartMoneyConceptsDataFairValueGap
}

func (e SmartMoneyConceptsDataOrderBlock) Equal(k SmartMoneyConceptsDataOrderBlock) bool {
	return e.IsTop == k.IsTop &&
		e.Time.Equal(k.Time)
}

type SmartMoneyConceptsDataOrderBlockList []SmartMoneyConceptsDataOrderBlock

func (e SmartMoneyConceptsDataOrderBlockList) Contains(k SmartMoneyConceptsDataOrderBlock) bool {
	for i := range e {
		if e[i].Equal(k) {
			return true
		}
	}
	return false
}

func (e SmartMoneyConceptsDataOrderBlockList) Add(k SmartMoneyConceptsDataOrderBlock) SmartMoneyConceptsDataOrderBlockList {
	if !e.Contains(k) {
		e = append([]SmartMoneyConceptsDataOrderBlock{k}, e...)
		// e = append(e, k)
	}
	return e
}

func (e SmartMoneyConceptsDataOrderBlockList) Remove(k SmartMoneyConceptsDataOrderBlock) SmartMoneyConceptsDataOrderBlockList {
	list := make(SmartMoneyConceptsDataOrderBlockList, len(e))
	copy(list, e)

	for x := range e {
		if e[x].Equal(k) {
			return append(list[:x], list[x+1:]...)
		}
	}

	return nil
}

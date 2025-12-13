package trend

import (
	"time"

	"github.com/idoall/stockindicator/utils/klines"
	"github.com/shopspring/decimal"
)

// CoinSMCKeyLevelRole 关键价位角色
type CoinSMCKeyLevelRole uint32

const (
	CoinSMCKeyLevelRoleUnknown          CoinSMCKeyLevelRole = 0
	CoinSMCKeyLevelRoleStrongSupport    CoinSMCKeyLevelRole = 1 << iota // 强支撑位
	CoinSMCKeyLevelRoleWeakSupport                                      // 弱支撑位
	CoinSMCKeyLevelRoleStrongResistance                                 // 强阻力位
	CoinSMCKeyLevelRoleWeakResistance                                   // 弱阻力位
)

func (r CoinSMCKeyLevelRole) String() string {
	switch r {
	case CoinSMCKeyLevelRoleStrongSupport:
		return "强支撑位"
	case CoinSMCKeyLevelRoleWeakSupport:
		return "弱支撑位"
	case CoinSMCKeyLevelRoleStrongResistance:
		return "强阻力位"
	case CoinSMCKeyLevelRoleWeakResistance:
		return "弱阻力位"
	default:
		return "未知"
	}
}

// 强弱高低点
type SmartMoneyConceptsDataStrongWeak struct {
	Time  time.Time
	Value decimal.Decimal
	Role  CoinSMCKeyLevelRole // 支撑位或阻力位角色
}

// 订单区块管理
type SmartMoneyConceptsDataOrderBlock struct {
	IsTop bool
	Time  time.Time
	Close decimal.Decimal
	High  decimal.Decimal
	Low   decimal.Decimal
	Open  decimal.Decimal
}

// SmartMoneyConceptsDataEqualHighLow 相等高低点（EQH/EQL）数据结构
//
// 代表两个或多个相近的高点/低点，通常是机构的流动性扫单区域
type SmartMoneyConceptsDataEqualHighLow struct {
	// ========== 基础信息 ==========

	// Time 产生时间（第二个相等点的时间）
	Time time.Time

	// StartIndex 产生时的 K 线索引
	// 用于快速计算 Duration，避免 O(n) 时间戳查找
	StartIndex int

	// Price 相等高/低点的价格
	Price decimal.Decimal

	// IsHigh 是否为相等高点
	// true: EQH(相等高点), false: EQL(相等低点)
	IsHigh bool

	// ========== 关联的两个点位信息 ==========

	// FirstPointTime 第一个点的时间
	FirstPointTime time.Time

	// FirstPointPrice 第一个点的价格
	FirstPointPrice decimal.Decimal

	// SecondPointTime 第二个点的时间
	SecondPointTime time.Time

	// SecondPointPrice 第二个点的价格
	SecondPointPrice decimal.Decimal

	// PointDistance 两个点之间的K线数量
	PointDistance int

	// ========== 失效信息 ==========

	// FilledTime 失效时间（价格突破该水平）
	// 零值表示尚未失效
	FilledTime time.Time

	// FilledPrice 填补价格（突破时的价格）
	// EQH: 向上突破时的 high 值
	// EQL: 向下突破时的 low 值
	FilledPrice decimal.Decimal

	// Duration 持续的K线数量
	// 从产生到失效之间经过的K线数量
	Duration int

	// BreakDirection 突破方向
	// "Upward"(向上突破) 或 "Downward"(向下突破)
	BreakDirection string

	// ========== 价格偏差统计 ==========

	// PriceDifference 两个点之间的价格差值（绝对值）
	PriceDifference decimal.Decimal

	// PriceDiffPercent 价格差值百分比
	PriceDiffPercent decimal.Decimal
}

// IsValid 判断 EQH/EQL 是否有效（非零值）
func (e SmartMoneyConceptsDataEqualHighLow) IsValid() bool {
	return !e.Price.IsZero() && !e.Time.IsZero()
}

// SignalType 信号类型枚举
type SignalType int

const (
	SignalTypeUnknown SignalType = 0
	SignalTypeCHoCH   SignalType = 1 << iota // 趋势反转 (Change of Character)
	SignalTypeBOS                            // 趋势延续 (Break of Structure)
)

func (s SignalType) String() string {
	switch s {
	case SignalTypeCHoCH:
		return "CHoCH"
	case SignalTypeBOS:
		return "BOS"
	default:
		return "Unknown"
	}
}

// PositionType 位置类型枚举
type PositionType int

const (
	PositionTypeUnknown PositionType = 0
	PositionTypeHigh    PositionType = 1 << iota // 顶部突破 (做多信号)
	PositionTypeLow                              // 底部突破 (做空信号)
)

func (p PositionType) String() string {
	switch p {
	case PositionTypeHigh:
		return "High"
	case PositionTypeLow:
		return "Low"
	default:
		return "Unknown"
	}
}

// PeriodType 周期类型枚举
type PeriodType int

const (
	PeriodTypeUnknown PeriodType = 0
	PeriodTypeShort   PeriodType = 1 << iota // 短周期 (5根K线)
	PeriodTypeLong                           // 长周期 (SwingLength根K线)
)

func (p PeriodType) String() string {
	switch p {
	case PeriodTypeShort:
		return "Short"
	case PeriodTypeLong:
		return "Long"
	default:
		return "Unknown"
	}
}

// SwingPriceSet 摆动点价格集合
// 用于止损计算时传递摆动点价格，简化函数参数
type SwingPriceSet struct {
	Short    float64 // 短周期摆动点（5根K线）
	Long     float64 // 长周期摆动点（SwingLength根K线）
	Trailing float64 // 追踪点（持续更新的最高/最低点）
}

// SmartMoneyConceptsDataStructureBreak 市场结构突破信号（CHoCH/BOS）数据结构
//
// 记录市场结构突破（BOS）和特征变化（CHoCH）的完整信息
type SmartMoneyConceptsDataStructureBreak struct {
	// ========== 基础分类信息 ==========

	// Time 摆动点形成时间（等同于 PivotTime）
	// 注意：此字段为了向后兼容保留，实际存储的是摆动点形成的时间
	Time time.Time

	// Type 信号类型
	// SignalTypeCHoCH: Change of Character（特征变化，趋势反转）
	// SignalTypeBOS: Break of Structure（结构突破，趋势延续）
	Type SignalType

	// Period 周期类型
	// PeriodTypeShort: 短线（5根K线）
	// PeriodTypeLong: 长线（SwingLength根K线）
	Period PeriodType

	// Position 位置类型
	// PositionTypeHigh: 顶部突破（做多信号）
	// PositionTypeLow: 底部突破（做空信号）
	Position PositionType

	// ========== 突破点位信息 ==========

	// BreakPrice 突破的价格水平
	BreakPrice decimal.Decimal

	// ClosePrice 突破时的收盘价
	ClosePrice decimal.Decimal

	// BreakTime 实际突破时间
	// 当价格穿越摆动点位时的K线时间（即CrossOver/CrossUnder发生的时刻）
	// 这是突破信号产生的准确时间点
	// 关系：摆动点形成(PivotTime) → 价格突破(BreakTime) → 下单入场(EntryTime)
	BreakTime time.Time

	// ========== 被突破的点位信息 ==========

	// PivotTime 被突破的高/低点产生时间
	PivotTime time.Time

	// PivotPrice 被突破的高/低点价格
	PivotPrice decimal.Decimal

	// Duration 持续的K线数量
	// 从产生到失效之间经过的K线数量
	Duration int
}

// IsValid 判断结构突破信号是否有效（非零值）
func (s SmartMoneyConceptsDataStructureBreak) IsValid() bool {
	return !s.BreakPrice.IsZero() && !s.Time.IsZero()
}

// GetSignalID 获取完整的信号标识
// 返回如 "HighBOSShort", "LowCHoCHLong" 等
func (s SmartMoneyConceptsDataStructureBreak) GetSignalID() string {
	return s.Position.String() + s.Type.String() + s.Period.String()
}

// IsLongSignal 判断是否做多信号
func (s SmartMoneyConceptsDataStructureBreak) IsLongSignal() bool {
	return s.Position == PositionTypeHigh
}

// IsShortSignal 判断是否做空信号
func (s SmartMoneyConceptsDataStructureBreak) IsShortSignal() bool {
	return s.Position == PositionTypeLow
}

// SmartMoneyConceptsDataFairValueGap Fair Value Gap（公允价值缺口）数据结构
//
// FVG 是三根 K 线之间形成的价格缺口，代表市场失衡区域。
// 看涨 FVG：当前 K 线低点 > 两根前 K 线高点（存在向上缺口）
// 看跌 FVG：当前 K 线高点 < 两根前 K 线低点（存在向下缺口）
type SmartMoneyConceptsDataFairValueGap struct {
	// Time FVG 产生的时间
	Time time.Time

	// StartIndex FVG 产生时的 K 线索引
	// 用于快速计算 Duration，避免 O(n) 时间戳查找
	StartIndex int

	// Top Fair Value Gap 的上边界价格
	Top decimal.Decimal

	// Bottom Fair Value Gap 的下边界价格
	Bottom decimal.Decimal

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
	FilledPrice decimal.Decimal

	// Duration FVG 持续的 K 线数量
	// 从产生到被填补之间经过的 K 线数量
	Duration int
}

// IsValid 判断 FVG 是否有效（非零值）
//
// 返回:
//   - bool: true 表示 FVG 有效，false 表示 FVG 为空
func (f SmartMoneyConceptsDataFairValueGap) IsValid() bool {
	return f.Top != decimal.Zero && f.Bottom != decimal.Zero
}

// Clear 清空 FVG 数据，将所有字段重置为零值
func (f *SmartMoneyConceptsDataFairValueGap) Clear() {
	f.Time = time.Time{}
	f.StartIndex = 0
	f.Top = decimal.Zero
	f.Bottom = decimal.Zero
	f.IsBullish = false
	f.FilledTime = time.Time{}
	f.FilledPrice = decimal.Zero
	f.Duration = 0
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

type StrongWeakDetail struct {

	// StrongHigh 强高点（下跌趋势中的高点）
	// 表示可能的阻力位
	StrongHigh SmartMoneyConceptsDataStrongWeak

	// WeakHigh 弱高点（上涨趋势中的高点）
	// 可能会被突破，成为支撑位
	WeakHigh SmartMoneyConceptsDataStrongWeak

	// StrongLow 强低点（上涨趋势中的低点）
	// 表示可能的支撑位
	StrongLow SmartMoneyConceptsDataStrongWeak

	// WeakLow 弱低点（下跌趋势中的低点）
	// 可能会被突破，成为阻力位
	WeakLow SmartMoneyConceptsDataStrongWeak
}

// SmartMoneyConcepts 智能资金概念（SMC）交易策略结构体
//
// SMC 是一种基于机构交易行为的技术分析方法，通过识别市场结构变化来捕捉趋势转折点。
// 该策略主要关注以下几个核心概念：
//
//  1. CHoCH (Change of Character - 特征变化):
//     表示市场在一段时间内改变了其趋势或订单流方向。
//     当价格突破前一个摆动高点/低点，且之前处于相反趋势时，形成 CHoCH。
//     这是趋势反转的早期信号。
//
//  2. BOS (Break of Structure - 结构突破):
//     描述价格突破图表上关键支撑位或阻力位的重大价格变动。
//     当价格突破前一个摆动高点/低点，且保持当前趋势方向时，形成 BOS。
//     这是趋势延续的确认信号。
//
//  3. Order Blocks (订单区块):
//     机构交易者在市场反转前最后一次下单的价格区域。
//     这些区域往往成为未来的支撑或阻力位。
//
//  4. EQH/EQL (Equal Highs/Equal Lows - 相等高点/相等低点):
//     两个或多个相近的高点或低点，通常表示流动性积累区域。
//     机构可能在这些区域进行扫单操作。
//
//  5. Strong/Weak High/Low (强/弱高点/低点):
//     根据趋势背景判断高点和低点的强弱。
//     强高点：下跌趋势中的高点；强低点：上涨趋势中的低点。
//     弱高点：上涨趋势中的高点；弱低点：下跌趋势中的低点。
type SmartMoneyConcepts struct {

	// SwingLength 实时摆动结构的周期长度
	// 用于计算摆动高点和摆动低点的回溯期数
	// 较大的值会产生更平滑但滞后的摆动点
	SwingLength int

	// EQHEQL_BarsConfirmation 相等高点/低点的确认K线数量（最小值为1）
	// 用于确认相等高点和相等低点所需的K线数量
	// 增加此值可以减少虚假信号，但可能错过一些有效信号
	EQHEQL_BarsConfirmation int

	// EQHEQL_Threshold 相等高点/低点检测的灵敏度阈值（取值范围：1-5）
	// 用于检测相等高低点的价格接近程度
	// 较低的值将返回较少但更精确的相等高低点信号
	// 阈值基于 ATR（平均真实波幅）的倍数
	EQHEQL_Threshold int

	// EQHEQL_Enable 是否启用相等高点/低点检测
	EQHEQL_Enable bool

	// Name 指标名称，包含参数信息
	Name string

	// data 存储每个K线的计算结果
	// data []SmartMoneyConceptsData

	// ohlc OHLC 数据（开盘价、最高价、最低价、收盘价）
	ohlc *klines.OHLC

	// OrderBlockNumber 保留的订单区块数量
	// 控制在图表上显示多少个最近的订单区块
	OrderBlockNumber int

	// OrderBlockBullish 看涨订单区块列表
	// 存储可能作为支撑位的价格区域
	OrderBlockBullish SmartMoneyConceptsDataOrderBlockList

	// OrderBlockBearish 看跌订单区块列表
	// 存储可能作为阻力位的价格区域
	OrderBlockBearish SmartMoneyConceptsDataOrderBlockList

	// ========== 内部状态变量 ==========

	// os1Prev 长周期摆动的前一个振荡器状态（0=看涨，1=看跌）
	os1Prev int

	// os2Prev 短周期摆动的前一个振荡器状态（0=看涨，1=看跌）
	os2Prev int

	// idxPrevBullish 前一个看涨订单区块的索引位置
	idxPrevBullish int

	// idxPrevBearish 前一个看跌订单区块的索引位置
	idxPrevBearish int

	// ========== 强弱高低点（仅在K线末尾显示） ==========
	StrongWeak_Enable bool

	StrongWeakDetail *StrongWeakDetail // 强弱高低点分析

	// ========== Fair Value Gap 配置 ==========

	// FVG_Enable 是否启用 Fair Value Gap 检测
	// true: 启用 FVG 检测，false: 禁用 FVG 检测
	FVG_Enable bool

	// FVG_AutoThreshold 是否启用自动阈值过滤
	// true: 使用动态阈值过滤波动性小的 FVG
	// false: 不使用阈值，检测所有 FVG
	FVG_AutoThreshold bool

	// FVG_KeepHistory 是否保留历史 FVG 记录
	// true: 保留所有已失效的 FVG 用于统计分析
	// false: 不保留历史，节省内存
	FVG_KeepHistory bool

	// ========== Fair Value Gap 历史记录 ==========

	// FVG_History 所有已失效的 FVG 历史记录
	// 用于统计分析、回测等场景
	FVG_History []SmartMoneyConceptsDataFairValueGap

	// ========== 活跃 FVG 列表（公开访问）==========

	// BullishFVGList 当前所有活跃的看涨 FVG 列表
	// 可以同时存在多个看涨 FVG，每个都有独立的生命周期
	BullishFVGList []SmartMoneyConceptsDataFairValueGap

	// BearishFVGList 当前所有活跃的看跌 FVG 列表
	// 可以同时存在多个看跌 FVG，每个都有独立的生命周期
	BearishFVGList []SmartMoneyConceptsDataFairValueGap

	// ========== EQH/EQL 配置参数 ==========

	// EQHEQL_KeepHistory 是否保留历史 EQH/EQL 记录
	// true: 保留所有已失效的 EQH/EQL 用于统计分析
	// false: 不保留历史，节省内存
	// 注意: EQHEQL_Enable 参数已在前面定义（第59行），用于控制是否启用 EQH/EQL 检测
	EQHEQL_KeepHistory bool

	// ========== EQH/EQL 列表（公开访问）==========

	// EQHList 当前所有活跃的相等高点列表
	EQHList []SmartMoneyConceptsDataEqualHighLow

	// EQLList 当前所有活跃的相等低点列表
	EQLList []SmartMoneyConceptsDataEqualHighLow

	// EQHEQL_History 所有已失效的 EQH/EQL 历史记录
	EQHEQL_History []SmartMoneyConceptsDataEqualHighLow

	// ========== CHoCH/BOS 配置参数 ==========

	// StructureBreak_Enable 是否启用 CHoCH/BOS 历史记录功能
	// true: 记录 CHoCH/BOS 的生命周期（产生、失效、持续时间等）
	// false: 不记录历史
	StructureBreak_Enable bool

	// ========== CHoCH/BOS 列表（按类型分开存储）==========

	// HighBOSShortList 顶部短线BOS列表
	HighBOSShortList []SmartMoneyConceptsDataStructureBreak

	// HighBOSLongList 顶部长线BOS列表
	HighBOSLongList []SmartMoneyConceptsDataStructureBreak

	// HighCHoCHShortList 顶部短线CHoCH列表
	HighCHoCHShortList []SmartMoneyConceptsDataStructureBreak

	// HighCHoCHLongList 顶部长线CHoCH列表
	HighCHoCHLongList []SmartMoneyConceptsDataStructureBreak

	// LowBOSShortList 底部短线BOS列表
	LowBOSShortList []SmartMoneyConceptsDataStructureBreak

	// LowBOSLongList 底部长线BOS列表
	LowBOSLongList []SmartMoneyConceptsDataStructureBreak

	// LowCHoCHShortList 底部短线CHoCH列表
	LowCHoCHShortList []SmartMoneyConceptsDataStructureBreak

	// LowCHoCHLongList 底部长线CHoCH列表
	LowCHoCHLongList []SmartMoneyConceptsDataStructureBreak

	// ========== 性能优化缓存 ==========

	// cachedATR 缓存的 ATR 值，避免重复计算
	// 在 Calculation() 方法中首次计算后缓存，后续调用直接复用
	cachedATR []float64

	// atrPeriod ATR 计算周期（默认 200）
	// 用于订单区块过滤和相等高低点检测
	atrPeriod int
}

package klines

import (
	"errors"
	"fmt"
	"time"
)

// Consts here define basic time intervals
const (
	HundredMilliseconds  = Interval(100 * time.Millisecond)
	ThousandMilliseconds = 10 * HundredMilliseconds
	TenSecond            = Interval(10 * time.Second)
	FifteenSecond        = Interval(15 * time.Second)
	ThirtySecond         = 2 * FifteenSecond
	OneMin               = Interval(time.Minute)
	ThreeMin             = 3 * OneMin
	FiveMin              = 5 * OneMin
	TenMin               = 10 * OneMin
	FifteenMin           = 15 * OneMin
	ThirtyMin            = 30 * OneMin
	OneHour              = Interval(time.Hour)
	TwoHour              = 2 * OneHour
	ThreeHour            = 3 * OneHour
	FourHour             = 4 * OneHour
	SixHour              = 6 * OneHour
	SevenHour            = 7 * OneHour
	EightHour            = 8 * OneHour
	TwelveHour           = 12 * OneHour
	OneDay               = 24 * OneHour
	TwoDay               = 2 * OneDay
	ThreeDay             = 3 * OneDay
	SevenDay             = 7 * OneDay
	FifteenDay           = 15 * OneDay
	OneWeek              = 7 * OneDay
	TwoWeek              = 2 * OneWeek
	ThreeWeek            = 3 * OneWeek
	OneMonth             = 30 * OneDay
	ThreeMonth           = 90 * OneDay
	SixMonth             = 2 * ThreeMonth
	NineMonth            = 3 * ThreeMonth
	OneYear              = 365 * OneDay
	FiveDay              = 5 * OneDay
)

var (
	errNilKline = errors.New("KlineItem item is nil")
	// ErrCanOnlyUpscaleCandles returns when attempting to upscale candles
	ErrCanOnlyUpscaleCandles = errors.New("interval must be a longer duration to scale")
	// ErrInvalidInterval defines when an interval is invalid e.g. interval <= 0
	ErrInvalidInterval = errors.New("invalid/unset interval")
	// ErrWholeNumberScaling returns when old interval data cannot neatly fit into new interval size
	ErrWholeNumberScaling            = errors.New("old interval must scale properly into new candle")
	errCandleDataNotPadded           = errors.New("candle data not padded")
	ErrInsufficientCandleData        = errors.New("insufficient candle data to generate new candle")
	errCannotEstablishTimeWindow     = errors.New("cannot establish time window")
	errCandleOpenTimeIsNotUTCAligned = errors.New("candle open time is not UTC aligned")
)

// Item holds all the relevant information for internal kline elements
type Item struct {
	Exchange string
	Interval Interval
	Symbol   string
	Code     string
	Candles  []*Candle
}

// Interval type for kline Interval usage
type Interval time.Duration

// Candle struct
type Candle struct {
	Amount float64 `json:"Amount,omitempty"` // 成交额
	Count  int64   `json:"Count,omitempty"`  // 成交笔数
	Open   float64 `json:"Open,omitempty"`   // 开盘价
	Close  float64 `json:"Close,omitempty"`  // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 `json:"Low,omitempty"`    // 最低价
	High   float64 `json:"High,omitempty"`   // 最高价
	Volume float64 `json:"Volume,omitempty"` // 成交量
	// 涨跌幅
	ChangePercent float64 `json:"ChangePercent,omitempty"`
	// 是否阳线上涨
	IsBullMarket bool `json:"IsBullMarket,omitempty"`
	// Time         time.Time `json:"Time,omitempty"`
	TimeUnix int64 `json:"TimeUnix,omitempty"` // k线时间
}

func (e *Candle) String() string {
	return fmt.Sprintf("%s Open:%f High:%f Low:%f Close:%f Volume:%f Bull:%+v ChangePercent:%.4f%%", time.Unix(e.TimeUnix, 0).Format("2006-01-02 15:04:05"), e.Open, e.High, e.Low, e.Close, e.Volume, e.IsBullMarket, e.ChangePercent*100)
}

// IntervalRangeHolder 保存整个间隔范围
// 以及所有内容的开始结束日期
type IntervalRangeHolder struct {
	Start  IntervalTime
	End    IntervalTime
	Ranges []IntervalRange
	Limit  int
}

// IntervalRange 是基于交易所 API 请求限制的蜡烛子集
type IntervalRange struct {
	Start     IntervalTime
	End       IntervalTime
	Intervals []IntervalData
}

// 间隔数据用于监控哪些蜡烛包含数据
// 确定是否缺少任何数据
type IntervalData struct {
	Start   IntervalTime
	End     IntervalTime
	HasData bool
}

// IntervalTime benchmarks demonstrate, see
// BenchmarkJustifyIntervalTimeStoringUnixValues1 &&
// BenchmarkJustifyIntervalTimeStoringUnixValues2
type IntervalTime struct {
	Time  time.Time
	Ticks int64
}

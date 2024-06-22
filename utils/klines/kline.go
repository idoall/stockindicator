package klines

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/idoall/stockindicator/utils/commonutils"
)

// String returns numeric string
func (i Interval) String() string {
	return i.Duration().String()
}

// Duration returns interval casted as time.Duration for compatibility
func (i Interval) Duration() time.Duration {
	return time.Duration(i)
}

// ConvertToNewInterval 允许将蜡烛图缩放为更大的蜡烛图
// 例如，如果蜡烛图足够多，则将一日蜡烛图转换为三日蜡烛图。不完整的蜡烛图不会被转换，例如，4 根一日蜡烛图将
// 转换为一根三日蜡烛图，跳过第四根。
func (e *Item) ConvertToNewInterval(newInterval Interval) (*Item, error) {
	if e == nil {
		return nil, errNilKline
	}
	if e.Interval <= 0 {
		return nil, fmt.Errorf("%w for old candle", ErrInvalidInterval)
	}
	if newInterval <= 0 {
		return nil, fmt.Errorf("%w for new candle", ErrInvalidInterval)
	}
	if newInterval <= e.Interval {
		return nil, fmt.Errorf("%w %s is less than or equal to %s",
			ErrCanOnlyUpscaleCandles,
			newInterval,
			e.Interval)
	}
	if newInterval%e.Interval != 0 {
		return nil, fmt.Errorf("%s %w %s",
			e.Interval,
			ErrWholeNumberScaling,
			newInterval)
	}

	start := time.Unix(e.Candles[0].TimeUnix, 0)
	end := time.Unix(e.Candles[len(e.Candles)-1].TimeUnix, 0).Add(e.Interval.Duration())

	var window time.Duration
	if start.Before(end) {
		window = end.Sub(start)
	} else {
		window = start.Sub(end)
	}

	if expected := int(window / e.Interval.Duration()); expected != len(e.Candles) {
		return nil, fmt.Errorf("%w expected candles %d but have only %d when converting from %s to %s interval",
			errCandleDataNotPadded,
			expected,
			len(e.Candles),
			e.Interval,
			newInterval)
	}

	oldIntervalsPerNewCandle := int(newInterval / e.Interval)
	candles := make([]*Candle, len(e.Candles)/oldIntervalsPerNewCandle)
	if len(candles) == 0 {
		return nil, fmt.Errorf("%w to %v no candle data", ErrInsufficientCandleData, newInterval)
	}

	// 有的时候不取整，重新计算开始、结束时间
	newIntervalRangeHolder, err := CalculateCandleDateRanges(start, end, newInterval, 100000)
	if err != nil {
		panic(err)
	}

	var target int

	// 对要返回的数据时间先赋值，后面好计算
	for _, rangeHolder := range newIntervalRangeHolder.Ranges {
		for _, intervalsData := range rangeHolder.Intervals {
			if target >= len(candles) {
				break
			}
			candles[target] = &Candle{
				TimeUnix: intervalsData.Start.Time.Unix(),
			}
			target++
		}
	}

	target = 0
	for x := range e.Candles {

		if e.Candles[x].TimeUnix == candles[target].TimeUnix && candles[target].Open == 0 {
			candles[target].Open = e.Candles[x].Open
		}

		if e.Candles[x].TimeUnix == candles[target].TimeUnix || e.Candles[x].TimeUnix > candles[target].TimeUnix && e.Candles[x].High > candles[target].High {
			candles[target].High = e.Candles[x].High
		}

		if e.Candles[x].TimeUnix == candles[target].TimeUnix || e.Candles[x].TimeUnix > candles[target].TimeUnix && (candles[target].Low == 0 || e.Candles[x].Low < candles[target].Low) {
			candles[target].Low = e.Candles[x].Low
		}

		if e.Candles[x].TimeUnix == candles[target].TimeUnix || (e.Candles[x].TimeUnix > candles[target].TimeUnix && e.Candles[x].TimeUnix < (candles[target].TimeUnix+int64(newInterval.Duration().Seconds()))) {
			candles[target].Volume += e.Candles[x].Volume
		}
		if (e.Candles[x].TimeUnix + int64(e.Interval.Duration().Seconds())) == (candles[target].TimeUnix + int64(newInterval.Duration().Seconds())) {
			candles[target].Close = e.Candles[x].Close
			candles[target].ChangePercent = (candles[target].Close - candles[target].Open) / candles[target].Open
			if candles[target].Close > candles[target].Open {
				candles[target].IsBullMarket = true
			}
		}

		if x < len(e.Candles)-1 && target < len(candles)-1 && e.Candles[x+1].TimeUnix == candles[target+1].TimeUnix {
			// 注意：下面检查了后续切片的长度，因此如果我们无法制作完整的蜡烛，我们可以
			// 立即中断。例如，一小时蜡烛中有 60 分钟
			// 蜡烛，我们还剩下 59 分钟蜡烛。
			// 整个过程被劈开。
			if len(e.Candles[x:])-1 < oldIntervalsPerNewCandle {
				break
			}

			target++
		}
	}
	return &Item{
		Exchange: e.Exchange,
		Interval: newInterval,
		Candles:  candles,
	}, nil
}

// RemoveOutsideRangeCopy 创建一个新的 Item，删除开始和结束日期之外的所有蜡烛。
func (e *Item) RemoveOutsideRangeCopy(start, end time.Time) *Item {

	candles := make([]*Candle, len(e.Candles))
	for i, candle := range e.Candles {
		copiedCandle := *candle
		candles[i] = &copiedCandle
		i++
	}

	klineItem := &Item{
		Exchange: e.Exchange,
		Interval: e.Interval,
		Candles:  candles,
	}
	klineItem.RemoveOutsideRange(start, end)
	return klineItem
}

// RemoveOutsideRange 删除开始和结束日期之外的所有蜡烛。
func (e *Item) RemoveOutsideRange(start, end time.Time) {
	target := 0
	for _, keep := range e.Candles {
		if keep.TimeUnix == start.Unix() || (keep.TimeUnix > start.Unix() && keep.TimeUnix < end.Unix()) {
			e.Candles[target] = keep
			target++
		}
	}
	e.Candles = e.Candles[:target]
}

// RemoveDuplicates 删除所有重复的蜡烛。注意：此函数中使用了就地过滤以进行优化并保持切片引用指针
func (e *Item) RemoveDuplicates() {
	lookup := make(map[int64]bool)
	target := 0
	for _, keep := range e.Candles {
		// 如果时间不存在
		if key := keep.TimeUnix; !lookup[key] {
			// 添加时间
			lookup[key] = true
			// 重新设置索引
			e.Candles[target] = keep
			target++
		}
	}
	e.Candles = e.Candles[:target]
}

// SortCandlesByTimestamp sorts candles by timestamp
func (e *Item) SortCandlesByTimestamp(desc bool) {
	if desc {
		sort.Slice(e.Candles, func(i, j int) bool { return e.Candles[i].TimeUnix < e.Candles[j].TimeUnix })
		return
	}
	sort.Slice(e.Candles, func(i, j int) bool { return e.Candles[i].TimeUnix > e.Candles[j].TimeUnix })
}

// addPadding 插入填充时间对齐，当交易所不提供所有数据时
// 当在某个时间间隔内没有活动时。
// Start 定义请求开始，由于从此点开始可能没有活动
// 需要指定这一点。 ExclusiveEnd 定义结束日期
// 其中不包括蜡烛，因此从开始开始的所有内容基本上都可以
// 用空格添加。
func (e *Item) addPadding(start, exclusiveEnd time.Time, purgeOnPartial bool) error {
	if e == nil {
		return errNilKline
	}

	if e.Interval <= 0 {
		return ErrInvalidInterval
	}

	window := exclusiveEnd.Sub(start)
	if window <= 0 {
		return errCannotEstablishTimeWindow
	}
	padded := make([]*Candle, int(window/e.Interval.Duration()))

	var target int
	for x := range padded {
		switch {
		case target >= len(e.Candles):
			padded[x] = &Candle{
				TimeUnix: start.Unix(),
			}
		case e.Candles[target].TimeUnix != start.Unix():
			if e.Candles[target].TimeUnix < start.Unix() {
				return fmt.Errorf("%w when it should be %s truncated at a %s interval",
					errCandleOpenTimeIsNotUTCAligned,
					start.Add(e.Interval.Duration()),
					e.Interval)
			}
			padded[x] = &Candle{
				TimeUnix: start.Unix(),
			}
		default:
			padded[x] = e.Candles[target]
			target++
		}
		start = start.Add(e.Interval.Duration())
	}

	// NOTE: This checks if the end time exceeds time.Now() and we are capturing
	// a partially created candle. This will only delete an element if it is
	// empty.
	if purgeOnPartial {
		lastElement := padded[len(padded)-1]
		if lastElement.Volume == 0 &&
			lastElement.Open == 0 &&
			lastElement.High == 0 &&
			lastElement.Low == 0 &&
			lastElement.Close == 0 {
			padded = padded[:len(padded)-1]
		}
	}
	e.Candles = padded
	return nil
}

func (e *Item) Clear() {
	clear(e.Candles)
	e.Candles = nil
}

// TotalCandlesPerInterval returns the total number of candle intervals between the start and end date
func TotalCandlesPerInterval(start, end time.Time, interval Interval) int64 {
	if interval <= 0 {
		return 0
	}
	window := end.Sub(start)
	return int64(window) / int64(interval)
}

// CalculateCandleDateRanges 将计算日期范围内间隔内的预期蜡烛数据
// 如果 API 在请求中可以生成的蜡烛数量有限，它将自动将
// 范围分离到限制中
func CalculateCandleDateRanges(start, end time.Time, interval Interval, limit uint32) (*IntervalRangeHolder, error) {
	if err := commonutils.StartEndTimeCheck(start, end); err != nil && !errors.Is(err, commonutils.ErrStartAfterTimeNow) {
		return nil, err
	}
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}

	start = start.Round(interval.Duration())
	end = end.Round(interval.Duration())
	window := end.Sub(start)
	count := int64(window) / int64(interval)
	requests := float64(count) / float64(limit)

	switch {
	case requests <= 1:
		requests = 1
	case limit == 0:
		requests, limit = 1, uint32(count)
	case requests-float64(int64(requests)) > 0:
		requests++
	}

	potentialRequests := make([]IntervalRange, int(requests))
	requestStart := start
	for x := range potentialRequests {
		potentialRequests[x].Start = CreateIntervalTime(requestStart)

		count -= int64(limit)
		if count < 0 {
			potentialRequests[x].Intervals = make([]IntervalData, count+int64(limit))
		} else {
			potentialRequests[x].Intervals = make([]IntervalData, limit)
		}

		for y := range potentialRequests[x].Intervals {
			potentialRequests[x].Intervals[y].Start = CreateIntervalTime(requestStart)
			requestStart = requestStart.Add(interval.Duration())
			potentialRequests[x].Intervals[y].End = CreateIntervalTime(requestStart)
		}
		potentialRequests[x].End = CreateIntervalTime(requestStart)
	}
	return &IntervalRangeHolder{
		Start:  CreateIntervalTime(start),
		End:    CreateIntervalTime(requestStart),
		Ranges: potentialRequests,
		Limit:  int(limit),
	}, nil
}

// HasDataAtDate determines whether a there is any data at a set
// date inside the existing limits
func (h *IntervalRangeHolder) HasDataAtDate(t time.Time) bool {
	tu := t.Unix()
	if tu < h.Start.Ticks || tu > h.End.Ticks {
		return false
	}
	for i := range h.Ranges {
		if tu < h.Ranges[i].Start.Ticks || tu >= h.Ranges[i].End.Ticks {
			continue
		}

		for j := range h.Ranges[i].Intervals {
			if tu >= h.Ranges[i].Intervals[j].Start.Ticks &&
				tu < h.Ranges[i].Intervals[j].End.Ticks {
				return h.Ranges[i].Intervals[j].HasData
			}
		}
	}
	return false
}

// CreateIntervalTime is a simple helper function to set the time twice
func CreateIntervalTime(tt time.Time) IntervalTime {
	return IntervalTime{Time: tt, Ticks: tt.Unix()}
}

// Equal allows for easier unix comparison
func (i *IntervalTime) Equal(tt time.Time) bool {
	return tt.Unix() == i.Ticks
}

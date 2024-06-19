package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
)

// Kdj is the main object
type Kdj struct {
	Name   string
	Period int //默认计算几天的
	data   []KdjData
	kline  *klines.Item
}

type KdjData struct {
	Time time.Time
	RSV  float64
	K    float64
	D    float64
	J    float64
}

// NewKdj new Func
// 超短线选择：K线：9（部分人选择设置为6）；D线：3；J线：3；
// 短线选择：K线：18；D线：3；J线：3；
// 中线选择：；K线：24；D线：3；J线：3。
func NewKdj(klineItem *klines.Item, period int) *Kdj {
	m := &Kdj{
		Name:   fmt.Sprintf("Kdj%d", period),
		kline:  klineItem,
		Period: period,
	}
	return m
}

// NewDefaultKdj new Func
func NewDefaultKdj(klineItem *klines.Item) *Kdj {
	return NewKdj(klineItem, 9)
}

// Calculation Func
func (e *Kdj) Calculation() *Kdj {

	//计算 rsv , k , d
	rsv, k, d := e.calculationKD(e.kline)
	arrayLen := len(rsv)
	e.data = make([]KdjData, len(rsv))
	for i := 0; i < arrayLen; i++ {
		e.data[i] = KdjData{
			RSV:  rsv[i],
			K:    k[i],
			D:    d[i],
			Time: e.kline.Candles[i].Time,
		}
	}
	j := e.calculationJ()
	for i := 0; i < arrayLen; i++ {
		e.data[i].J = j[i]
	}
	return e
}

// AnalysisCrossSide Func
// 分析金、死叉买卖方向
func (e *Kdj) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		prevItem := e.data[i-1]
		// 当 K、J同时金叉D，为买入信号
		if v.K > v.D && v.J > v.D && prevItem.K < prevItem.D && prevItem.J < prevItem.D {
			sides[i] = utils.Buy
		} else if v.K < v.D && v.J < v.D && prevItem.K > prevItem.D && prevItem.J > prevItem.D {
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

// GetData return Data
func (e *Kdj) GetData() []KdjData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// GetListK Func
func (e *Kdj) GetListK() []float64 {
	var result []float64
	for _, v := range e.data {
		result = append(result, v.K)
	}
	return result
}

// GetListD Func
func (e *Kdj) GetListD() []float64 {
	var result []float64
	for _, v := range e.data {
		result = append(result, v.D)
	}
	return result
}

// GetListJ Func
func (e *Kdj) GetListJ() []float64 {
	var result []float64
	for _, v := range e.data {
		result = append(result, v.J)
	}
	return result
}

// calculationKD 计算出kd值
func (e *Kdj) calculationKD(records *klines.Item) (rsv, k, d []float64) {

	var periodLowArr, periodHighArr []float64

	var ohlc = records.GetOHLC()
	var lows = ohlc.Low
	var highs = ohlc.High
	var closes = ohlc.Close

	length := len(lows)
	rsv = make([]float64, length)
	k = make([]float64, length)
	d = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodLowArr = append(periodLowArr, lows[i])
		periodHighArr = append(periodHighArr, highs[i])

		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.

		if e.Period == len(periodLowArr) {
			lowest := e.arrayLowest(periodLowArr)
			highest := e.arrayHighest(periodHighArr)

			if highest-lowest < 0.000001 {
				rsv[i] = 100
			} else {
				rsv[i] = (closes[i] - lowest) / (highest - lowest) * 100
			}

			k[i] = (2.0/3)*k[i-1] + 1.0/3*rsv[i]
			d[i] = (2.0/3)*d[i-1] + 1.0/3*k[i]
			// remove first value in array.
			periodLowArr = periodLowArr[1:]
			periodHighArr = periodHighArr[1:]
		} else {
			k[i] = 50
			d[i] = 50
			rsv[i] = 0
		}
	}
	return rsv, k, d

}

// calculationJ 计算J值
func (e *Kdj) calculationJ() []float64 {
	length := len(e.data)
	var j []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		item := e.data[i]
		j[i] = 3*item.K - 2*item.D

	}
	return j
}

// func (e *Kdj) highest(priceArray []float64, periods int) []float64 {
// 	var periodArr []float64
// 	length := len(priceArray)
// 	var HighestLine []float64 = make([]float64, length)

// 	// Loop through the entire array.
// 	for i := 0; i < length; i++ {
// 		// add points to the array.
// 		periodArr = append(periodArr, priceArray[i])
// 		// 1: Check if array is "filled" else create null point in line.
// 		// 2: Calculate average.
// 		// 3: Remove first value.
// 		if periods == len(periodArr) {
// 			HighestLine[i] = e.arrayHighest(periodArr)

// 			// remove first value in array.
// 			periodArr = periodArr[1:]
// 		} else {
// 			HighestLine[i] = 0
// 		}
// 	}

// 	return HighestLine
// }

// func (e *Kdj) lowest(priceArray []float64, periods int) []float64 {
// 	var periodArr []float64
// 	length := len(priceArray)
// 	var LowestLine []float64 = make([]float64, length)

// 	// Loop through the entire array.
// 	for i := 0; i < length; i++ {
// 		// add points to the array.
// 		periodArr = append(periodArr, priceArray[i])
// 		// 1: Check if array is "filled" else create null point in line.
// 		// 2: Calculate average.
// 		// 3: Remove first value.
// 		if periods == len(periodArr) {
// 			LowestLine[i] = e.arrayLowest(periodArr)

// 			// remove first value in array.
// 			periodArr = periodArr[1:]
// 		} else {
// 			LowestLine[i] = 0
// 		}
// 	}

// 	return LowestLine
// }

func (e *Kdj) arrayLowest(priceArray []float64) float64 {
	length := len(priceArray)
	var lowest = priceArray[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if priceArray[i] < lowest {
			lowest = priceArray[i]
		}
	}

	return lowest
}

func (e *Kdj) arrayHighest(priceArray []float64) float64 {
	length := len(priceArray)
	var highest = priceArray[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if priceArray[i] > highest {
			highest = priceArray[i]
		}
	}

	return highest
}

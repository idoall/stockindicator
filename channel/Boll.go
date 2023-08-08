// Copyright 2016 mshk.top, lion@mshk.top
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package channel

import (
	"fmt"
	"math"
	"time"

	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/ta"
)

/*
日Boll指标的计算公式
中轨线=N日的移动平均线
上轨线=中轨线+两倍的标准差
下轨线=中轨线－两倍的标准差
日Boll指标的计算过程
1）计算MA
MA=N日内的收盘价之和÷N
2）计算标准差MD
MD=平方根N日的（C－MA）的两次方之和除以N
3）计算MB、UP、DN线
MB=（N－1）日的MA
UP=MB+2×MD
DN=MB－2×MD
*/

// 原始的加速通道有兩種，一種是以 20 週移動平均線來計算的通道
// 另一種則是以 20 日移動平均線來計算的通道；
// 兩者算法相同，只是時間長度不同。
// 前者適合交易e 一個月至三個月左右的中期趨勢，後者適合拿來做短線價差交易，時間長度通常在一個月以內。
type Boll struct {
	Name    string
	PeriodN int //计算周期
	PeriodK int //带宽
	data    []BollData
	kline   utils.Klines
}

type BollData struct {
	Upper  float64
	Middle float64
	Lower  float64
	Time   time.Time
}

// NewBoll Func
// 使用方法，先添加最早日期的数据,最后一条应该是当前日期的数据，结果与 AICoin 对比完全一致
func NewBoll(list utils.Klines, periodN, periodK int) *Boll {
	return &Boll{Name: fmt.Sprintf("Boll%d-%d", periodN, periodK), PeriodN: periodN, PeriodK: periodK, kline: list}
}

// NewDefaultBoll Func
func NewDefaultBoll(list utils.Klines) *Boll {
	return NewBoll(list, 20, 2.0)
}

// dma MD=平方根N日的（C－MA）的两次方之和除以N
func (e *Boll) dma(sma []float64) []float64 {
	result := make([]float64, len(e.kline))
	period := e.PeriodN
	sum2 := 0.0
	for i, v := range e.kline {
		sum2 += v.Close * v.Close
		if i < period-1 {
			result[i] = 0.0
		} else {
			result[i] = math.Sqrt(sum2/float64(period) - sma[i]*sma[i])
			w := e.kline[i-(period-1)].Close
			sum2 -= w * w
		}
	}

	return result
}

// Calculation Func
func (e *Boll) Calculation() *Boll {
	l := len(e.kline)

	e.data = make([]BollData, l)

	var middle = trend.NewSma(e.kline, e.PeriodN).GetValues()
	var md = e.dma(middle)
	var upper = ta.Add(middle, md)
	var lower = ta.Subtract(middle, md)

	for i, _ := range middle {
		e.data[i] = BollData{
			Time:   e.kline[i].Time,
			Middle: middle[i],
			Upper:  upper[i],
			Lower:  lower[i],
		}
	}
	return e
}

// AnalysisSide Func
// 收盘价高于 upperBand 时提供卖出操作
// 收盘价低于 lowerBand 值时提供买入操作。
func (e *Boll) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]
		var price = e.kline[i].Close
		var prevPrice = e.kline[i-1].Close

		if price > v.Upper && prevPrice < prevItem.Upper {
			sides[i] = utils.Buy
		} else if price < v.Lower && prevPrice > prevPrice {
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

// GetPoints Func
func (e *Boll) GetData() []BollData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

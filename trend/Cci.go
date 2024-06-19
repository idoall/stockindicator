package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

/*
TYP:=(HIGH+LOW+CLOSE)/3;
       Cci:(TYP-MA(TYP,N))/(0.015*AVEDEV(TYP,N));
TYP比较容易理解，（最高价+最低价+收盘价）÷3
MA(TYP,N) 也比较简单，就是N天的TYP的平均值
AVEDEV(TYP,N) 比较难理解，是对TYP进行绝对平均偏差的计算。
也就是说N天的TYP减去MA(TYP,N)的绝对值的和的平均值。
表达式：
MA = MA(TYP,N)
AVEDEV(TYP,N) =( | 第N天的TYP - MA |   +  | 第N-1天的TYP - MA | + ...... + | 第1天的TYP - MA | ) ÷ N
Cci = （TYP－MA）÷ AVEDEV(TYP,N)   ÷0.015

计算商品通道指数有几个步骤。
以下示例适用于典型的20周期cci：
cci =（典型价格 - tp的20周期平均值）/（.015 x平均偏差）
典型价格（tp）=（高+低+近）/ 3
常数= 0.015
出于缩放目的，该常数被设置为.015。
通过包含常数，大多数cci值将落入100到-100的范围内。
计算平均偏差有三个步骤。
1.减去最近的20个期间，简单地从该时期的每个典型价格（tp）移动。
2.严格使用绝对值对这些数字进行求和。
3.将步骤3中生成的值除以期间总数
*/

// Cci struct
type Cci struct {
	Name      string
	Period    int     //默认计算几天的
	SMAPeriod int     //默认计算几天的
	factor    float64 //计算系数
	data      []CciData
	kline     *klines.Item
}

type CciData struct {
	Value float64
	Time  time.Time
}

// NewCci new Func
func NewCci(klineItem *klines.Item, period, smaPeriod int) *Cci {
	m := &Cci{
		Name:      fmt.Sprintf("CCI%d-SMA%d", period, smaPeriod),
		kline:     klineItem,
		Period:    period,
		SMAPeriod: smaPeriod,
		factor:    0.015,
	}
	return m
}

// NewDefaultCci new Func
func NewDefaultCci(klineItem *klines.Item) *Cci {
	return NewCci(klineItem, 20, 20)
}

// Calculation Func
func (e *Cci) Calculation() *Cci {

	hlc3 := e.kline.HLC3()

	ccis := ta.CCI(hlc3, e.Period, e.SMAPeriod)

	//计算 Cci
	for i := 0; i < len(ccis); i++ {
		e.data = append(e.data, CciData{
			Time:  e.kline.Candles[i].Time,
			Value: ccis[i],
		})
	}
	return e
}

// func (e *Cci) Calculation() *Cci {

// 	// 计算TYP
// 	// TYP:=(HIGH+LOW+CLOSE)/3;
// 	for i := 0; i < len(e.kline.Candles); i++ {
// 		item := e.kline[i]
// 		typicalPrice := (item.High + item.Low + item.Close) / 3.0
// 		e.typicalPrice = append(e.typicalPrice, typicalPrice)
// 	}

// 	// 计算MA
// 	// MA = MA(TYP,N)
// 	// var closeArray []float64
// 	// for _, v := range e.kline.Candles {
// 	// 	closeArray = append(closeArray, v.Close)
// 	// }
// 	var tempKlineArray *klines.Item
// 	for i := 0; i < len(e.typicalPrice); i++ {
// 		tempKlineArray = append(tempKlineArray, utils.Kline{
// 			Close: e.typicalPrice[i],
// 			Time:  e.kline.Candles[i].Time,
// 		})
// 	}
// 	maPoints := NewMa(tempKlineArray, e.Period).GetData()
// 	for _, v := range maPoints {
// 		e.maPrice = append(e.maPrice, v.Value)
// 	}

// 	//计算平均偏差有三个步骤。
// 	// 1.减去最近的20个期间，简单地从该时期的每个典型价格（tp）移动。
// 	// 2.严格使用绝对值对这些数字进行求和。
// 	// 3.将步骤3中生成的值除以期间总数
// 	for i := 0; i < len(e.maPrice); i++ {
// 		if i < e.Period-1 {
// 			e.avedevPrice = append(e.avedevPrice, 0.0)
// 			continue
// 		}

// 		var avedevSum float64
// 		for j := 0; j < e.Period; j++ {
// 			avedevSum += math.Abs(e.typicalPrice[i-j] - e.maPrice[i])
// 		}
// 		tempAvedevPrice, _ := commonutils.FloatFromString(fmt.Sprintf("%d", e.Period))
// 		e.avedevPrice = append(e.avedevPrice, avedevSum/tempAvedevPrice)
// 	}

// 	//计算 Cci
// 	// cci =（典型价格 - tp的20周期平均值）/（.015 x平均偏差）
// 	for i := 0; i < len(e.maPrice); i++ {
// 		var p CciData
// 		p.Time = e.kline.Candles[i].Time
// 		if i < e.Period-1 {
// 			p.Value = 0
// 			e.data = append(e.data, p)
// 			continue
// 		}

// 		p.Value = (e.typicalPrice[i] - e.maPrice[i]) / (e.avedevPrice[i] * e.factor)
// 		e.data = append(e.data, p)
// 	}
// 	return e
// }

// AnalysisSide Func
func (e *Cci) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline.Candles))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		prevItem := e.data[i-1]
		// APO 上穿零表示看涨，而下穿零表示看跌。
		if v.Value < -100.0 && prevItem.Value > -100 {
			sides[i] = utils.Buy
		} else if v.Value > 100 && prevItem.Value < 100 {
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

// GetPoints return Point
func (e *Cci) GetData() []CciData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	return e.data
}

// GetValue return Value
func (e *Cci) GetValue() []float64 {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	var result []float64
	for _, v := range e.data {
		result = append(result, v.Value)
	}
	return result
}

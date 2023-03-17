package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Trix struct
// TRIX(Triple Exponentially Smoothed Average)根据移动平均线理论，对一条平均线进行三次平滑处理，再根据这条移动平均线的变动情况来预测价格的长期走势。
// 其最大的优点就是可以过滤短期波动的干扰，以避免频繁操作而带来的失误和损失，因此最适合于对行情的中长期走势的研判。
// TRIX由下往上交叉MATRIX时，为长期买进信号。
// TRIX由上往下交叉MATRIX时，为长期卖出信号。
// TRIX与价格产生背离时，应注意随时会反转。
// 盘整行情本指标不适用。
// 计算公式
// 　　1. TR=收盘价的N日指数移动平均；
// 　　2.TRIX=（TR-昨日TR）/昨日TR*100；
// 　　3.MATRIX=TRIX的M日简单移动平均；
// 　　4.参数N设为12，参数M设为20；
// TR：= EMA（EMA（EMA（CLOSE，N），N），N）;
// TRIX = TR-REF（TR，1））/REF（TR，1）*100;
// MATR = MA（TRIX，M）
type Trix struct {
	Name    string
	PeriodN int
	PeriodM int
	data    []TrixData
	kline   utils.Klines
}

type TrixData struct {
	Trix   float64
	MaTrix float64
	Time   time.Time
}

// NewTrix new Func
func NewTrix(list utils.Klines, periodN, periodM int) *Trix {
	m := &Trix{
		Name:    fmt.Sprintf("Trix%d-%d", periodN, periodM),
		kline:   list,
		PeriodN: periodN,
		PeriodM: periodM,
	}
	return m
}

// NewDefaultTrix new Func
func NewDefaultTrix(list utils.Klines) *Trix {
	return NewTrix(list, 12, 9)
}

// Calculation Func
func (e *Trix) Calculation() *Trix {
	e.data = make([]TrixData, len(e.kline))

	periodN := e.PeriodN
	periodM := e.PeriodM

	// 计算 EMA1 值
	ema1 := NewEma(e.kline, periodN).GetValues()
	// 计算 EMA2 值
	ema2 := NewEma(utils.CloseArrayToKline(ema1), periodN).GetValues()
	// 计算 TR 值
	tr := NewEma(utils.CloseArrayToKline(ema2), periodN).GetValues()
	trix := utils.MultiplyBy(utils.PercentDiff(tr, 1), 100)
	maTrix := NewMa(utils.CloseArrayToKline(trix), periodM).GetData()

	for i := 0; i < len(maTrix); i++ {
		e.data[i] = TrixData{
			Time:   e.kline[i].Time,
			Trix:   trix[i],
			MaTrix: maTrix[i].Value,
		}
	}
	return e
}

// AnalysisSide Func
func (e *Trix) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]
		if v.Trix > v.MaTrix && prevItem.Trix < prevItem.MaTrix {
			sides[i] = utils.Buy
		} else if v.Trix < v.MaTrix && prevItem.Trix > prevItem.MaTrix {
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

// GetData return Point
func (e *Trix) GetData() []TrixData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

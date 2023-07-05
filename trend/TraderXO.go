package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// TraderXO struct
type TraderXO struct {
	FastPeriod int
	SlowPeriod int
	Name       string
	data       []TraderXOData
	kline      utils.Klines
}

type TraderXOData struct {
	Time time.Time
	Fast float64
	Slow float64
}

// NewTraderXO new Func
func NewTraderXO(list utils.Klines, fastPeriod, slowPeriod int) *TraderXO {
	m := &TraderXO{
		Name:       fmt.Sprintf("TraderXO%d-%d", fastPeriod, slowPeriod),
		kline:      list,
		FastPeriod: fastPeriod,
		SlowPeriod: slowPeriod,
	}
	return m
}

// NewTraderXO new Func
func NewDefaultTraderXO(list utils.Klines) *TraderXO {
	return NewTraderXO(list, 12, 25)
}

// Calculation Func
func (e *TraderXO) Calculation() *TraderXO {

	// Define EMAs
	v_fastEMAList := NewEma(e.kline, e.FastPeriod).GetValues()
	v_slowEMAList := NewEma(e.kline, e.SlowPeriod).GetValues()

	defer func() {
		v_fastEMAList = nil
		v_slowEMAList = nil
	}()

	e.data = make([]TraderXOData, len(e.kline))
	for i := 0; i < len(e.kline); i++ {

		e.data[i] = TraderXOData{
			Time: e.kline[i].Time,
			Fast: v_fastEMAList[i],
			Slow: v_slowEMAList[i],
		}

	}

	return e
}

// GetData Func
func (e *TraderXO) GetValues() (fast []float64, slow []float64) {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	for _, v := range e.data {
		fast = append(fast, v.Fast)
		slow = append(slow, v.Slow)
	}
	return
}

// GetData Func
func (e *TraderXO) GetData() []TraderXOData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

// AnalysisSide Func
func (e *TraderXO) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	// sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}
	var countBuy = 0
	var countSell = 0
	var buys = make([]bool, len(e.data))
	var sells = make([]bool, len(e.data))

	for i := 0; i < len(e.data); i++ {
		if i < 1 {
			sides[i] = utils.Hold
			continue
		}
		v := e.data[i]

		buys[i] = v.Fast > v.Slow
		sells[i] = v.Fast < v.Slow

		if buys[i] {
			countBuy += 1
		}

		if buys[i] {
			countSell = 0
		}

		if sells[i] {
			countSell += 1
		}

		if sells[i] {
			countBuy = 0
		}
		buysignal := countBuy < 2 && countBuy > 0 && countSell < 1 && buys[i] && !buys[i-1]
		sellsignal := countSell > 0 && countSell < 2 && countBuy < 1 && sells[i] && !sells[i-1]

		// 当值低于20认为是超卖，为买入信号
		if buysignal {
			sides[i] = utils.Buy
		} else if sellsignal {
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

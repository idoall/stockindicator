package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
	"github.com/idoall/stockindicator/utils/klines"
	"github.com/idoall/stockindicator/utils/ta"
)

// TraderXO struct
type TraderXO struct {
	FastPeriod int
	SlowPeriod int
	Name       string
	data       []TraderXOData
	ohlc       *klines.OHLC
}

type TraderXOData struct {
	Time time.Time
	Fast float64
	Slow float64
}

// NewTraderXO new Func
func NewTraderXO(klineItem *klines.Item, fastPeriod, slowPeriod int) *TraderXO {
	m := &TraderXO{
		Name:       fmt.Sprintf("TraderXO%d-%d", fastPeriod, slowPeriod),
		FastPeriod: fastPeriod,
		SlowPeriod: slowPeriod,
	}
	m.ohlc = klineItem.GetOHLC()
	return m
}

// NewTraderXO new Func
func NewTraderXOOHLC(ohlc *klines.OHLC, fastPeriod, slowPeriod int) *TraderXO {
	m := &TraderXO{
		Name:       fmt.Sprintf("TraderXO%d-%d", fastPeriod, slowPeriod),
		ohlc:       ohlc,
		FastPeriod: fastPeriod,
		SlowPeriod: slowPeriod,
	}
	return m
}

// NewTraderXO new Func
func NewDefaultTraderXO(klineItem *klines.Item) *TraderXO {
	return NewTraderXO(klineItem, 12, 25)
}

// Calculation Func
func (e *TraderXO) Calculation() *TraderXO {

	var closes = e.ohlc.Close
	var times = e.ohlc.Time
	// Define EMAs
	v_fastEMAList := ta.Ema(e.FastPeriod, closes)
	v_slowEMAList := ta.Ema(e.SlowPeriod, closes)

	defer func() {
		v_fastEMAList = nil
		v_slowEMAList = nil
	}()

	e.data = make([]TraderXOData, len(closes))
	for i := 0; i < len(closes); i++ {

		e.data[i] = TraderXOData{
			Time: times[i],
			Fast: v_fastEMAList[i],
			Slow: v_slowEMAList[i],
		}

	}

	return e
}

// GetData Func
func (e *TraderXO) GetValues() ([]float64, []float64) {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	var fast = make([]float64, len(e.data))
	var slow = make([]float64, len(e.data))
	for i, v := range e.data {
		fast[i] = v.Fast
		slow[i] = v.Slow
	}
	return fast, slow
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

	if len(e.data) == 0 {
		e = e.Calculation()
	}
	sides := make([]utils.Side, len(e.data))
	// sides := make([]utils.Side, len(e.kline.Candles))

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

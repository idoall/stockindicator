package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Dma struct
// DMA(Different of Moving Average)是利用两条不同期间的平均线，来判断当前买卖能量的大小和未来价格趋势的一种中短期投资指标。
// DMA线向上交叉AMA线，买进；DMA线向下交叉AMA线，卖出。
// 当DMA和AMA均>0(即在图形上表示为它们处于零线以上)并向上移动时，一般表示为市场处于多头行情中，为买入信号，可以买入或持有观望；
// 当DMA和AMA均<0(即在图形上表示为它们处于零线以下)并向下移动时，一般表示为市场处于空头行情中，为卖出信号，可以卖出或观望。
// 当DMA和AMA均<0时，经过一段时间的下跌后，如果两者同时从低位向上移动时，为买进信号；
// 当DMA和AMA均>0，在经过一段时间的上涨后，如果两者同时从高位向下移动时，为卖出信号。
// DMA指标与价格产生背离时的交叉信号，可信度较高。
// DMA指标亦适于结合形态理论进行分析。
// DMA指标、MACD指标、TRIX指标三者构成一组指标群，互相验证。
type Dma struct {
	Name                        string
	PeriodN1, PeriodN2, PeriodM int
	data                        []DmaData
	kline                       utils.Klines
}

type DmaData struct {
	Dma  float64
	Ama  float64
	Time time.Time
}

// NewDma new Func
func NewDma(list utils.Klines, periodN1, periodN2, periodM int) *Dma {
	m := &Dma{
		Name:     fmt.Sprintf("Dma%d-%d-%d", periodN1, periodN2, periodM),
		kline:    list,
		PeriodN1: periodN1,
		PeriodN2: periodN2,
		PeriodM:  periodM,
	}
	return m
}

// NewDefaultDma new Func
func NewDefaultDma(list utils.Klines) *Dma {
	return NewDma(list, 10, 50, 10)
}

// Calculation Func
func (e *Dma) Calculation() *Dma {
	e.data = make([]DmaData, len(e.kline))

	periodN1 := e.PeriodN1
	periodN2 := e.PeriodN2
	periodM := e.PeriodM

	// 计算 MA1 值
	ma1 := NewMa(e.kline, periodN1).GetValues()
	// 计算 MA2 值
	ma2 := NewMa(e.kline, periodN2).GetValues()
	dma := utils.Subtract(ma1, ma2)
	ama := NewMa(utils.CloseArrayToKline(dma), periodM).GetValues()

	for i := 0; i < len(dma); i++ {
		e.data[i] = DmaData{
			Time: e.kline[i].Time,
			Dma:  dma[i],
			Ama:  ama[i],
		}
	}
	return e
}

// AnalysisSide Func
func (e *Dma) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]
		if v.Dma > v.Ama && prevItem.Dma < prevItem.Ama {
			sides[i] = utils.Buy
		} else if v.Dma < v.Ama && prevItem.Dma > prevItem.Ama {
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
func (e *Dma) GetData() []DmaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

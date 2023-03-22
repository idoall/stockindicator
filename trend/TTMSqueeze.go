package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// TTMSqueeze Pro 挤压指标
// Usage: https://cn.tradingview.com/script/0drMdHsO-TTM-Squeeze-Pro/
type TTMSqueeze struct {
	Name                                     string
	Period                                   int
	MultBB, MultLowKC, MultMidKC, MultHighKC float64
	data                                     []TTMSqueezeData
	kline                                    utils.Klines
}

type TTMSqueezeColor int

const (
	SqzGreen  TTMSqueezeColor = 0
	SqzOrange TTMSqueezeColor = 1 << iota
	SqzRed
	SqzBlack
)

// String implements the stringer interface
func (s TTMSqueezeColor) String() string {
	switch s {
	case SqzGreen:
		return "Green"
	case SqzOrange:
		return "Red"
	case SqzRed:
		return "Red"
	case SqzBlack:
		return "Black"
	default:
		return "UNKNOWN"
	}
}

type MomColor int

const (
	MomGreen MomColor = 0
	MomBlue  MomColor = 1 << iota
	MomRed
	MomYellow
	MomUnknown
)

// String implements the stringer interface
func (s MomColor) String() string {
	switch s {
	case MomGreen:
		return "Green"
	case MomBlue:
		return "Blue"
	case MomRed:
		return "Red"
	case MomYellow:
		return "Yellow"
	default:
		return "UNKNOWN"
	}
}

type TTMSqueezeData struct {
	Time     time.Time
	Mom      float64
	MomColor MomColor
	SqzColor TTMSqueezeColor
}

// NewTTMSqueeze new Func
func NewTTMSqueeze(list utils.Klines, period int, multBB, multHighKC, multMidKC, multLowKC float64) *TTMSqueeze {
	m := &TTMSqueeze{
		Name:       fmt.Sprintf("TTMSqueeze%d-%f-%f-%f-%f", period, multBB, multLowKC, multMidKC, multHighKC),
		kline:      list,
		Period:     period,
		MultBB:     multBB,
		MultLowKC:  multLowKC,
		MultMidKC:  multMidKC,
		MultHighKC: multHighKC,
	}
	return m
}

// NewDefaultTTMSqueeze new Func
func NewDefaultTTMSqueeze(list utils.Klines) *TTMSqueeze {
	return NewTTMSqueeze(list, 20, 2.0, 1.0, 1.5, 2.0)
}

// Calculation Func
func (e *TTMSqueeze) Calculation() *TTMSqueeze {
	e.data = make([]TTMSqueezeData, len(e.kline))

	period := e.Period
	multBB := e.MultBB
	multLowKC := e.MultLowKC
	multMidKC := e.MultMidKC
	multHighKC := e.MultHighKC

	highPrice := e.kline.GetOHLC().High
	lowPrice := e.kline.GetOHLC().Low
	closePrice := e.kline.GetOHLC().Close

	// Calculate BB
	basis := NewSma(e.kline, period).GetValues()
	dev := utils.MultiplyBy(utils.StdDev(closePrice, period, 1.0), multBB)
	upperBB := utils.Add(basis, dev)
	lowerBB := utils.Subtract(basis, dev)

	// Calculate KC
	basisKc := NewSma(e.kline, period).GetValues()
	devKC := NewSma(utils.CloseArrayToKline(utils.TrueRange(highPrice, lowPrice, closePrice)), period).GetValues()
	lowDevKC := utils.MultiplyBy(devKC, multLowKC)
	midDevKC := utils.MultiplyBy(devKC, multMidKC)
	highDevKC := utils.MultiplyBy(devKC, multHighKC)

	upperLowKC := utils.Add(basisKc, lowDevKC)
	lowerLowKC := utils.Subtract(basisKc, lowDevKC)
	upperMidKC := utils.Add(basisKc, midDevKC)
	lowerMidKC := utils.Subtract(basisKc, midDevKC)
	upperHighKC := utils.Add(basisKc, highDevKC)
	lowerHighKC := utils.Subtract(basisKc, highDevKC)

	//MOMENTUM OSCILLATOR
	avg := utils.AvgPrice(
		utils.AvgPrice(
			utils.Max(period, highPrice),
			utils.Min(period, lowPrice)),
		NewSma(e.kline, period).GetValues(),
	)
	x := utils.GenerateNumbers(0, float64(len(closePrice)), 1)
	mom := utils.MovingLinearRegressionUsingLeastSquare(period, x, utils.Subtract(closePrice, avg))

	for i := 0; i < len(e.kline); i++ {
		//NoSqz := lowerBB[i] < lowerLowKC[i] || upperBB[i] > upperLowKC[i]       //NO SQUEEZE: GREEN
		LowSqz := lowerBB[i] >= lowerLowKC[i] || upperBB[i] <= upperLowKC[i]    //LOW COMPRESSION: BLACK
		MidSqz := lowerBB[i] >= lowerMidKC[i] || upperBB[i] <= upperMidKC[i]    //MID COMPRESSION: RED
		HighSqz := lowerBB[i] >= lowerHighKC[i] || upperBB[i] <= upperHighKC[i] //HIGH COMPRESSION: ORANGE

		//SQUEEZE DOTS COLOR
		//sq_color := HighSqz ? SqzOrange : MidSqz ? SqzRed : LowSqz ? SqzBlack : SqzGreen
		sqzColor := utils.IfCase(HighSqz, SqzOrange, utils.IfCase(MidSqz, SqzRed, utils.IfCase(LowSqz, SqzBlack, SqzGreen)))

		//MOMENTUM HISTOGRAM COLOR
		momColor := MomUnknown
		if i > 0 {
			up := utils.IfCase(mom[i] > mom[i-1], MomGreen, MomBlue)
			down := utils.IfCase(mom[i] < mom[i-1], MomRed, MomYellow)
			momColor = utils.IfCase(mom[i] > 0, up, down)
		}

		e.data[i] = TTMSqueezeData{
			Time:     e.kline[i].Time,
			Mom:      mom[i],
			MomColor: momColor,
			SqzColor: sqzColor,
		}
	}

	return e
}

// AnalysisSide Func
func (e *TTMSqueeze) AnalysisSide() utils.SideData {
	sides := make([]utils.Side, len(e.kline))

	if len(e.data) == 0 {
		e = e.Calculation()
	}

	for i, v := range e.data {
		if i < 1 {
			continue
		}

		var prevItem = e.data[i-1]
		if v.SqzColor == SqzGreen && prevItem.SqzColor != SqzGreen {
			sides[i] = utils.Buy
		} else if v.SqzColor != SqzGreen && prevItem.SqzColor == SqzGreen {
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
func (e *TTMSqueeze) GetData() []TTMSqueezeData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}
	return e.data
}

package trend

import (
	"fmt"
	"time"

	"github.com/idoall/stockindicator/utils"
)

// Ma struct
type Ma struct {
	Name   string
	Period int //默认计算几天的Ma,KDJ一般是9，OBV是10、20、30
	data   []MaData
	kline  utils.Klines
}

type MaData struct {
	Value float64
	Time  time.Time
}

// NewMa new Func
func NewMa(list utils.Klines, period int) *Ma {
	m := &Ma{
		Name:   fmt.Sprintf("Ma%d", period),
		kline:  list,
		Period: period,
	}
	return m
}

// NewDefaultMa new Func
func NewDefaultMa(list utils.Klines) *Ma {
	return NewMa(list, 20)
}

// GetPoints return Point
func (e *Ma) GetData() []MaData {
	if len(e.data) == 0 {
		e = e.Calculation()
	}

	return e.data
}

// Calculation Func
func (e *Ma) Calculation() *Ma {
	for i := 0; i < len(e.kline); i++ {
		if i < e.Period-1 {
			p := MaData{}
			p.Time = e.kline[i].Time
			p.Value = 0.0
			e.data = append(e.data, p)
			continue
		}
		var sum float64
		for j := 0; j < e.Period; j++ {

			sum += e.kline[i-j].Close
		}

		p := MaData{}
		p.Time = e.kline[i].Time
		p.Value = sum / float64(e.Period)
		e.data = append(e.data, p)
	}
	return e
}

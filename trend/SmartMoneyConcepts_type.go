package trend

import (
	"time"

	"github.com/idoall/stockindicator/utils"
)

type SmartMoneyConceptsDataStrongWeak struct {
	Time  time.Time
	Value float64
}
type SmartMoneyConceptsDataOrderBlock struct {
	IsTop bool
	Kline utils.Kline
}

func (e SmartMoneyConceptsDataOrderBlock) Equal(k SmartMoneyConceptsDataOrderBlock) bool {
	return e.IsTop == k.IsTop &&
		e.Kline.Time.Equal(k.Kline.Time) &&
		e.Kline.High == k.Kline.High &&
		e.Kline.Low == k.Kline.Low &&
		e.Kline.Open == k.Kline.Open &&
		e.Kline.Close == k.Kline.Close
}

type SmartMoneyConceptsDataOrderBlockList []SmartMoneyConceptsDataOrderBlock

func (e SmartMoneyConceptsDataOrderBlockList) Contains(k SmartMoneyConceptsDataOrderBlock) bool {
	for i := range e {
		if e[i].Equal(k) {
			return true
		}
	}
	return false
}

func (e SmartMoneyConceptsDataOrderBlockList) Add(k SmartMoneyConceptsDataOrderBlock) SmartMoneyConceptsDataOrderBlockList {
	if !e.Contains(k) {
		e = append([]SmartMoneyConceptsDataOrderBlock{k}, e...)
		// e = append(e, k)
	}
	return e
}

func (e SmartMoneyConceptsDataOrderBlockList) Remove(k SmartMoneyConceptsDataOrderBlock) SmartMoneyConceptsDataOrderBlockList {
	list := make(SmartMoneyConceptsDataOrderBlockList, len(e))
	copy(list, e)

	for x := range e {
		if e[x].Equal(k) {
			return append(list[:x], list[x+1:]...)
		}
	}

	return nil
}

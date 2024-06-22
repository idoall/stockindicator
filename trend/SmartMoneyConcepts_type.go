package trend

import (
	"time"
)

type SmartMoneyConceptsDataStrongWeak struct {
	Time  time.Time
	Value float64
}
type SmartMoneyConceptsDataOrderBlock struct {
	IsTop bool
	Time  time.Time
	Close float64
	High  float64
	Low   float64
	Open  float64
}

func (e SmartMoneyConceptsDataOrderBlock) Equal(k SmartMoneyConceptsDataOrderBlock) bool {
	return e.IsTop == k.IsTop &&
		e.Time.Equal(k.Time)
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

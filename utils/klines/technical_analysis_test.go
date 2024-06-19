package klines

import (
	"testing"
)

func TestGetOHLC(t *testing.T) {
	t.Parallel()
	if (&Item{Candles: []*Candle{{Open: 1337}}}).GetOHLC() == nil {
		t.Fatal("unexpected value")
	}
}

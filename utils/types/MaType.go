package types

type MATypes uint32

const (
	UnknownMATypes MATypes = 0
	SMA            MATypes = 1 << iota
	EMA
	WMA
	DEMA
	TEMA
	TRIMA
	KAMA
	MAMA
	T3MA
)

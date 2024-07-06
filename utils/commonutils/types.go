package commonutils

type Integer interface {
	int | int8 | int16 | int32 | int64
}

type Float interface {
	float32 | float64
}

type Number interface {
	Integer | Float
}

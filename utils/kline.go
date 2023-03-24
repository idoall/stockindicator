package utils

func CloseArrayToKline(closeing []float64) (result Klines) {
	for _, v := range closeing {
		result = append(result, Kline{
			Close: v,
		})
	}
	return result
}

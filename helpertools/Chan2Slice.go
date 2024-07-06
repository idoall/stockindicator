package helpertools

func Chan2Slice[T any](c <-chan T) []T {
	var result []T

	for n := range c {
		result = append(result, n)
	}

	return result
}

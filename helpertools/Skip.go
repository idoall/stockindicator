package helpertools

// Skip 从 float64 的给定通道中跳过指定数量的元素。
//
// Example:
//
//	c := Slice2Chan([]int{2, 4, 6, 8, 10})
//	actual := Skip(c, 2)
//	fmt.Println(Chan2Slice(actual)) // [6 8 10]
func Skip[T any](c <-chan T, count int) <-chan T {
	result := make(chan T, cap(c))

	go func() {
		for i := 0; i < count; i++ {
			<-c
		}

		Pipe(c, result)
	}()

	return result
}

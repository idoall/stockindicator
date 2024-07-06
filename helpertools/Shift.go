package helpertools

// Shift 接收一个数字通道，按照指定的数量向右移动它们，
// 并使用提供的填充值填充任何缺失的值。
//
// Example:
//
//	ch := Slice2Chan([]int{2, 4, 6, 8})
//	result := Shift(ch, 4, 0)
//	fmt.Println(Chan2Slice(result))	// [0 0 0 0 2 4 6 8]
func Shift[T any](c <-chan T, count int, fill T) <-chan T {
	result := make(chan T, cap(c)+count)

	go func() {
		for i := 0; i < count; i++ {
			result <- fill
		}

		Pipe(c, result)
	}()

	return result
}

package helpertools

// Duplicate 会复制给定的只接收通道，方法是读取从该通道传出的每个值
// 并将它们发送到请求数量的新输出通道上。
//
// Example:
//
//	temp1 := []float64{-10, 20, -4, -5}
//	temp2 := helpertools.Duplicate[float64](helpertools.Slice2Chan(temp1), 2)
//
//	for i, source := range temp1 {
//		for chanItemIndex, chanItem := range temp2 {
//			chainItemValue := <-chanItem
//			fmt.Printf("[%d][chanItemIndex:%d]source:%+v\tchainItemValue:%+v\n", i, chanItemIndex, source, chainItemValue)
//		}
//	}
func Duplicate[T any](input <-chan T, count int) []<-chan T {

	outputs := make([]chan T, count)
	result := make([]<-chan T, count)

	for i := range outputs {
		outputs[i] = make(chan T, cap(input))
		result[i] = outputs[i]
	}

	go func() {
		for _, output := range outputs {
			defer close(output)
		}

		for n := range input {
			for _, output := range outputs {
				output <- n
			}
		}
	}()

	return result
}

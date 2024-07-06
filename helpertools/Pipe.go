package helpertools

// Pipe 函数接受一个输入通道和一个输出通道，并将
// 输入通道中的所有元素复制到输出通道。
//
// Example:
//
//	input := Slice2Chan([]int{2, 4, 6, 8})
//	output := make(chan int)
//	Pipe(input, output)
//	fmt.println(Chan2Slice(output)) // [2, 4, 6, 8]
func Pipe[T any](f <-chan T, t chan<- T) {
	defer close(t)
	for n := range f {
		t <- n
	}
}

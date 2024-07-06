package helpertools

// Drain 排空指定通道。它会阻塞调用者。
func Drain[T any](c <-chan T) {
	for {
		_, ok := <-c
		if !ok {
			break
		}
	}
}

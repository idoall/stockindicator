package helpertools

// 操作将提供的操作函数应用于两个
// 数字输入通道的相应值，并将结果值发送到输出通道。
func Operate[A any, B any, R any](ac <-chan A, bc <-chan B, o func(A, B) R) <-chan R {
	oc := make(chan R)

	go func() {
		defer close(oc)

		for {
			an, ok := <-ac
			if !ok {
				Drain(bc)
				break
			}

			bn, ok := <-bc
			if !ok {
				Drain(ac)
				break
			}

			oc <- o(an, bn)
		}
	}()

	return oc
}

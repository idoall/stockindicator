package helpertools

func Slice2Chan[T any](list []T) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		for _, n := range list {
			c <- n
		}
	}()

	return c
}

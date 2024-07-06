package helpertools

import "github.com/idoall/stockindicator/utils/commonutils"

// Apply 将给定的转换函数应用于输入通道中的每个元素，并返回包含转换后的值的新通道。转换函数将 float64 值作为输入
// 并返回 float64 值作为输出。
//
// Example:
//
//	timesTwo := helpertools.Apply(c, func(n int) int {
//		return n * 2
//	})
func Apply[T commonutils.Number](c <-chan T, f func(T) T) <-chan T {
	ac := make(chan T)

	go func() {
		defer close(ac)

		for n := range c {
			ac <- f(n)
		}
	}()

	return ac
}

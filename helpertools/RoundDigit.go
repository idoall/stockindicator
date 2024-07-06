package helpertools

import (
	"math"

	"github.com/idoall/stockindicator/utils/commonutils"
)

// RoundDigit rounds the given float64 number to d decimal places.
//
// Example:
//
//	n := helper.RoundDigit(10.1234, 2)
//	fmt.Println(n) // 10.12
func RoundDigit[T commonutils.Number](n T, d int) T {
	// m := math.Pow(10, float64(d))
	factor := float64(1)
	for i := 0; i < d; i++ {
		factor *= 10
	}
	return T(math.Round(float64(n)*factor) / factor)
}

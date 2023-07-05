package utils

import (
	"math"
	"time"

	"github.com/idoall/stockindicator/container/bst"
)

// Multiply 对 values 数组元素遍历乘以 multiplier 进行乘法运算.
func MultiplyBy(values []float64, multiplier float64) []float64 {
	result := make([]float64, len(values))

	for i, value := range values {
		result[i] = value * multiplier
	}

	return result
}

// Multiply 将 values1 与 values2 数组元素进行乘法运算.
func Multiply(values1, values2 []float64) []float64 {

	result := make([]float64, len(values1))

	for i := 0; i < len(result); i++ {
		result[i] = values1[i] * values2[i]
	}

	return result
}

// Divide 对 values 数组元素遍历进行 per 百分比运算，并返回结果.
func DivideBy(values []float64, per float64) []float64 {
	return MultiplyBy(values, float64(1)/per)
}

// Divide 将 values1 数组与 values2 数组，遍历每一对进行除法运算.
func Divide(values1, values2 []float64) []float64 {

	result := make([]float64, len(values1))

	for i := 0; i < len(result); i++ {
		val := values1[i] / values2[i]
		if math.IsNaN(val) || math.IsInf(val, 0) {
			val = 0
		}
		result[i] = val
	}

	return result
}

// Add 将 values1 与遍历 values2 数组元素相加.
func Add(values1, values2 []float64) []float64 {

	result := make([]float64, len(values1))
	for i := 0; i < len(result); i++ {
		result[i] = values1[i] + values2[i]
	}

	return result
}

// AddBy 将 values 数组遍历加上 addition.
func AddBy(values []float64, addition float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(result); i++ {
		result[i] = values[i] + addition
	}

	return result
}

// Subtract 遍历 values1 数组与 value2 数组进行减法运算.
func Subtract(values1, values2 []float64) []float64 {
	// 将 value2 进行负值运算
	subtract := MultiplyBy(values2, float64(-1))
	return Add(values1, subtract)
}

// Diff 按 before 右移后，再进行减法运算.
func Diff(values []float64, before int) []float64 {

	return Subtract(values, ShiftRight(before, values))
}

// PercentDiff 将 values 数组 从 before 开始，计算当前值-上一项值/上一项值的百分比
func PercentDiff(values []float64, before int) []float64 {
	result := make([]float64, len(values))

	for i := before; i < len(values); i++ {
		result[i] = (values[i] - values[i-before]) / values[i-before]
	}

	return result
}

// ShiftRightAndFillBy 将 values 数组的值 按 period 右移，填充 fill
func ShiftRightAndFillBy(period int, fill float64, values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(result); i++ {
		if i < period {
			result[i] = fill
		} else {
			result[i] = values[i-period]
		}
	}

	return result
}

// ShiftRight 将 values 数组的值按 period 进行右移.
func ShiftRight(period int, values []float64) []float64 {
	return ShiftRightAndFillBy(period, 0, values)
}

// RoundDigits 将值进行四舍五入.
func RoundDigits(value float64, digits int) float64 {

	// 计算10 的 digits 次幂
	n := math.Pow(10, float64(digits))

	// 返回最接近的整数，从零四舍五入
	return math.Round(value*n) / n
}

// RoundDigitsAll 将 values 数组遍历进行四舍五入.
func RoundDigitsAll(values []float64, digits int) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(result); i++ {
		result[i] = RoundDigits(values[i], digits)
	}

	return result
}

// GenerateNumbers 生成数字数组.
func GenerateNumbers(begin, end, step float64) []float64 {
	n := int(math.Round((end - begin) / step))

	numbers := make([]float64, n)

	for i := 0; i < n; i++ {
		numbers[i] = begin + (step * float64(i))
	}

	return numbers
}

// Pow 遍历 base 数组与 exponent 进行幂运行.
func Pow(base []float64, exponent float64) []float64 {
	result := make([]float64, len(base))

	for i := 0; i < len(result); i++ {
		result[i] = math.Pow(base[i], exponent)
	}

	return result
}

// ExtractSign 判断 values 数组每一项是否是正数.
func ExtractSign(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(result); i++ {
		if values[i] >= 0 {
			result[i] = 1
		} else {
			result[i] = -1
		}
	}

	return result
}

// KeepPositives 遍历 values 数组，如果小于0改为0.
func KeepPositives(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		if values[i] > 0 {
			result[i] = values[i]
		} else {
			result[i] = 0
		}
	}

	return result
}

// KeepNegatives 遍历 values 数组，如果大于0改为0..
func KeepNegatives(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		if values[i] < 0 {
			result[i] = values[i]
		} else {
			result[i] = 0
		}
	}

	return result
}

// Max 遍历数组 values，从右向左在 period 周期内，用最大值替换索引位置的值.
func Max(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)
	bst := bst.New()

	for i := 0; i < len(values); i++ {
		bst.Insert(values[i])

		if i >= period {
			bst.Remove(buffer[i%period])
		}

		buffer[i%period] = values[i]
		result[i] = bst.Max().(float64)
	}

	return result
}

// Min 遍历数组 values，从右向左在 period 周期内，用最小值替换索引位置的值.
func Min(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	buffer := make([]float64, period)
	bst := bst.New()

	for i := 0; i < len(values); i++ {
		bst.Insert(values[i])

		if i >= period {
			bst.Remove(buffer[i%period])
		}

		buffer[i%period] = values[i]
		result[i] = bst.Min().(float64)
	}

	return result
}

// PivotHigh 此函数返回枢轴高点的价格。 如果没有枢轴高点，则返回“0”。
func PivotHigh(values []float64, left, right int) []float64 {
	pivot := make([]float64, len(values))

	for i := left + right; i < len(values); i++ {
		rolling := values[i-right-left : i+1]

		bst := bst.New()
		bst.Inserts(rolling)
		m := bst.Max().(float64)
		if values[i-right] == m {
			pivot[i] = m
		}
	}
	return pivot
}

// PivotLow 此函数返回枢轴低点的价格。 如果没有枢轴低点，则返回“0”。
func PivotLow(values []float64, left, right int) []float64 {
	pivot := make([]float64, len(values))

	for i := left + right; i < len(values); i++ {
		rolling := values[i-right-left : i+1]

		bst := bst.New()
		bst.Inserts(rolling)
		m := bst.Min().(float64)
		if values[i-right] == m {
			pivot[i] = m
		}
	}
	return pivot
}

// PivotMax 以 index 为起始点，向 左、右 开始查找最大值
func PivotMax(values []float64, index, left, right int) (val float64, position int) {
	var leftVal, rightVal float64

	val = values[index]
	position = index

	var leftPositions = index
	var rightPositions = index
	// 向左找数量
	for i := 0; i <= left; i++ {
		newIndex := index - i
		if newIndex < 0 {
			break
		}
		if values[newIndex] > leftVal {
			leftVal = values[newIndex]
			leftPositions = newIndex
		}
	}

	// 向右找数量
	for i := 0; i <= right; i++ {
		newIndex := index + i
		if newIndex > len(values)-1 {
			break
		}
		if values[newIndex] > rightVal {
			rightVal = values[newIndex]
			rightPositions = newIndex
		}
	}

	if leftVal > rightVal {
		return leftVal, leftPositions
	} else if leftVal < rightVal {
		return rightVal, rightPositions
	}
	return values[index], index

}

// PivotMin 以 index 为起始点，向 左、右 开始查找最小值
func PivotMin(values []float64, index, left, right int) (val float64, position int) {

	var leftVal = values[index]
	var rightVal = values[index]

	var leftPositions = index
	var rightPositions = index

	// 向左找数量
	for i := 0; i <= left; i++ {
		newIndex := index - i
		if newIndex < 0 {
			break
		}
		if values[newIndex] < leftVal {
			leftVal = values[newIndex]
			leftPositions = newIndex
		}
	}

	// 向右找数量
	for i := 0; i <= right; i++ {
		newIndex := index + i
		if newIndex > len(values)-1 {
			break
		}
		if values[newIndex] < rightVal {
			rightVal = values[newIndex]
			rightPositions = newIndex
		}
	}

	if leftVal < rightVal {
		return leftVal, leftPositions
	} else if leftVal > rightVal {
		return rightVal, rightPositions
	}
	return values[index], index

}

// Sum 遍历数组 values，从左向右在 period 周期内总和，替换索引位置的值.
func Sum(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := 0.0

	for i := 0; i < len(values); i++ {
		sum += values[i]

		if i >= period {
			sum -= values[i-period]
		}

		result[i] = sum
	}

	return result
}

// Sqrt 依次计算 values 数组的平方根.
func Sqrt(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		result[i] = math.Sqrt(values[i])
	}

	return result
}

// Abs 将 values 数组的所有值转换为绝对值.
func Abs(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		result[i] = math.Abs(values[i])
	}

	return result
}

// mean 返回float64值数组的平均值
func mean(values []float64) float64 {
	var total float64 = 0
	for x := range values {
		total += values[x]
	}
	return total / float64(len(values))
}

// trueRange 返回高低闭合的真实范围
func trueRange(inHigh, inLow, inClose []float64) []float64 {
	outReal := make([]float64, len(inClose))

	startIdx := 1
	outIdx := startIdx
	today := startIdx
	for today < len(inClose) {
		tempLT := inLow[today]
		tempHT := inHigh[today]
		tempCY := inClose[today-1]
		greatest := tempHT - tempLT
		val2 := math.Abs(tempCY - tempHT)
		if val2 > greatest {
			greatest = val2
		}
		val3 := math.Abs(tempCY - tempLT)
		if val3 > greatest {
			greatest = val3
		}
		outReal[outIdx] = greatest
		outIdx++
		today++
	}

	return outReal
}

// variance 返回给定时间段的方差
func variance(inReal []float64, inTimePeriod int) []float64 {
	outReal := make([]float64, len(inReal))

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	periodTotal1 := 0.0
	periodTotal2 := 0.0
	trailingIdx := startIdx - nbInitialElementNeeded
	i := trailingIdx
	if inTimePeriod > 1 {
		for i < startIdx {
			tempReal := inReal[i]
			periodTotal1 += tempReal
			tempReal *= tempReal
			periodTotal2 += tempReal
			i++
		}
	}
	outIdx := startIdx
	for ok := true; ok; {
		tempReal := inReal[i]
		periodTotal1 += tempReal
		tempReal *= tempReal
		periodTotal2 += tempReal
		meanValue1 := periodTotal1 / float64(inTimePeriod)
		meanValue2 := periodTotal2 / float64(inTimePeriod)
		tempReal = inReal[trailingIdx]
		periodTotal1 -= tempReal
		tempReal *= tempReal
		periodTotal2 -= tempReal
		outReal[outIdx] = meanValue2 - meanValue1*meanValue1
		i++
		trailingIdx++
		outIdx++
		ok = i < len(inReal)
	}
	return outReal
}

// stdDev - Standard Deviation
func stdDev(inReal []float64, inTimePeriod int, inNbDev float64) []float64 {
	outReal := variance(inReal, inTimePeriod)

	if inNbDev != 1.0 {
		for i := 0; i < len(inReal); i++ {
			tempReal := outReal[i]
			if !(tempReal < 0.00000000000001) {
				outReal[i] = math.Sqrt(tempReal) * inNbDev
			} else {
				outReal[i] = 0.0
			}
		}
	} else {
		for i := 0; i < len(inReal); i++ {
			tempReal := outReal[i]
			if !(tempReal < 0.00000000000001) {
				outReal[i] = math.Sqrt(tempReal)
			} else {
				outReal[i] = 0.0
			}
		}
	}
	return outReal
}

// Nzs 以系列中的指定数替换NaN值。
func Nzs(vals []float64, replacement float64) []float64 {
	for i := 0; i < len(vals); i++ {
		vals[i] = Nz(vals[i], replacement)
	}
	return vals
}

// Nz 以系列中的指定数替换NaN值。
func Nz(vals float64, replacement float64) float64 {

	if math.IsNaN(vals) {
		return replacement
	}

	return vals
}

// LinearReg - 线性回归
func LinearReg(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		m := (inTimePeriodF*sumXY - sumX*sumY) / divisor
		b := (sumY - m*sumX) / inTimePeriodF
		outReal[outIdx] = b + m*(inTimePeriodF-1)
		outIdx++
		today++
	}
	return outReal
}

// LinearRegAngle - Linear Regression Angle
func LinearRegAngle(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		m := (inTimePeriodF*sumXY - sumX*sumY) / divisor
		outReal[outIdx] = math.Atan(m) * (180.0 / math.Pi)
		outIdx++
		today++
	}
	return outReal
}

// LinearRegIntercept - Linear Regression Intercept
func LinearRegIntercept(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		m := (inTimePeriodF*sumXY - sumX*sumY) / divisor
		outReal[outIdx] = (sumY - m*sumX) / inTimePeriodF
		outIdx++
		today++
	}
	return outReal
}

// LinearRegSlope - 线性回归斜率指标
func LinearRegSlope(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		outReal[outIdx] = (inTimePeriodF*sumXY - sumX*sumY) / divisor
		outIdx++
		today++
	}
	return outReal
}

// TRange - True Range
func TRange(inHigh []float64, inLow []float64, inClose []float64) []float64 {

	outReal := make([]float64, len(inClose))

	startIdx := 1
	outIdx := startIdx
	today := startIdx
	for today < len(inClose) {
		tempLT := inLow[today]
		tempHT := inHigh[today]
		tempCY := inClose[today-1]
		greatest := tempHT - tempLT
		val2 := math.Abs(tempCY - tempHT)
		if val2 > greatest {
			greatest = val2
		}
		val3 := math.Abs(tempCY - tempLT)
		if val3 > greatest {
			greatest = val3
		}
		outReal[outIdx] = greatest
		outIdx++
		today++
	}

	return outReal
}

func Stochastic(closing, highs, lows []float64, period int) (k, d []float64) {
	highestHigh := Max(period, highs)
	lowestLow := Min(period, lows)

	k = MultiplyBy(Divide(Subtract(closing, lowestLow), Subtract(highestHigh, lowestLow)), float64(100))
	d = Sma(3, k)
	return k, d
}

func Sma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1
		sum += value

		if i >= period {
			sum -= values[i-period]
			count = period
		}

		val := sum / float64(count)
		if math.IsNaN(val) || math.IsInf(val, -1) {
			result[i] = 0
		} else {
			result[i] = val
		}
	}

	values = nil
	return result
}

// Rolling Moving Average (RMA).
//
// R[0] to R[p-1] is SMA(values)
// R[p] and after is R[i] = ((R[i-1]*(p-1)) + v[i]) / p
//
// Returns r.
func Rma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1

		if i < period {
			sum += value
		} else {
			sum = (result[i-1] * float64(period-1)) + value
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}

func Ema(period int, values []float64) []float64 {
	result := make([]float64, len(values))

	// k := float64(2) / float64(1+period)

	for i, value := range values {
		if i > 0 {
			// result[i] = (value * k) + (result[i-1] * float64(1-k))
			// result[i] = k*(values[i]-result[i-1]) + result[i-1]
			result[i] = (2*value + float64(period-1)*(result[i-1])) / float64(period+1)
		} else {
			result[i] = value
		}
	}
	values = nil
	return result
}

// Wma - Weighted Moving Average
func Wma(period int, values []float64) []float64 {

	result := make([]float64, len(values))

	lookbackTotal := period - 1
	startIdx := lookbackTotal

	if period == 1 {
		copy(result, values)
		return result
	}
	divider := (period * (period + 1)) >> 1
	outIdx := period - 1
	trailingIdx := startIdx - lookbackTotal
	periodSum, periodSub := 0.0, 0.0
	inIdx := trailingIdx
	i := 1
	for inIdx < startIdx {
		tempReal := values[inIdx]
		periodSub += tempReal
		periodSum += tempReal * float64(i)
		inIdx++
		i++
	}
	trailingValue := 0.0
	for inIdx < len(values) {
		tempReal := values[inIdx]
		periodSub += tempReal
		periodSub -= trailingValue
		periodSum += tempReal * float64(period)
		trailingValue = values[trailingIdx]
		result[outIdx] = periodSum / float64(divider)
		periodSum -= periodSub
		inIdx++
		trailingIdx++
		outIdx++
	}
	return result
}

func GetTestKline() Klines {
	return []Kline{
		Kline{Open: 9986.300000, Close: 9800.010000, Low: 9705.000000, High: 10035.960000, Volume: 100683.796400, Time: time.UnixMicro(1588896000000000), ChangePercent: -0.018655, IsBullMarket: false},
		Kline{Open: 9800.020000, Close: 9539.400000, Low: 9520.000000, High: 9914.250000, Volume: 81950.679567, Time: time.UnixMicro(1588982400000000), ChangePercent: -0.026594, IsBullMarket: false},
		Kline{Open: 9539.100000, Close: 8722.770000, Low: 8117.000000, High: 9574.830000, Volume: 183865.182028, Time: time.UnixMicro(1589068800000000), ChangePercent: -0.085577, IsBullMarket: false},
		Kline{Open: 8722.770000, Close: 8561.520000, Low: 8200.000000, High: 9168.000000, Volume: 168807.251832, Time: time.UnixMicro(1589155200000000), ChangePercent: -0.018486, IsBullMarket: false},
		Kline{Open: 8562.040000, Close: 8810.790000, Low: 8528.780000, High: 8978.260000, Volume: 86522.780066, Time: time.UnixMicro(1589241600000000), ChangePercent: 0.029053, IsBullMarket: true},
		Kline{Open: 8810.990000, Close: 9309.370000, Low: 8792.990000, High: 9398.000000, Volume: 92466.274018, Time: time.UnixMicro(1589328000000000), ChangePercent: 0.056563, IsBullMarket: true},
		Kline{Open: 9309.350000, Close: 9791.980000, Low: 9256.760000, High: 9939.000000, Volume: 129565.377470, Time: time.UnixMicro(1589414400000000), ChangePercent: 0.051844, IsBullMarket: true},
		Kline{Open: 9791.970000, Close: 9316.420000, Low: 9150.000000, High: 9845.620000, Volume: 115890.761516, Time: time.UnixMicro(1589500800000000), ChangePercent: -0.048565, IsBullMarket: false},
		Kline{Open: 9315.960000, Close: 9381.270000, Low: 9220.000000, High: 9588.000000, Volume: 59587.627862, Time: time.UnixMicro(1589587200000000), ChangePercent: 0.007011, IsBullMarket: true},
		Kline{Open: 9380.810000, Close: 9680.040000, Low: 9322.100000, High: 9888.000000, Volume: 68647.764323, Time: time.UnixMicro(1589673600000000), ChangePercent: 0.031898, IsBullMarket: true},
		Kline{Open: 9681.110000, Close: 9733.930000, Low: 9464.230000, High: 9950.000000, Volume: 82006.603583, Time: time.UnixMicro(1589760000000000), ChangePercent: 0.005456, IsBullMarket: true},
		Kline{Open: 9733.930000, Close: 9775.530000, Low: 9474.000000, High: 9897.210000, Volume: 78539.760454, Time: time.UnixMicro(1589846400000000), ChangePercent: 0.004274, IsBullMarket: true},
		Kline{Open: 9775.130000, Close: 9511.430000, Low: 9326.000000, High: 9842.000000, Volume: 74923.738090, Time: time.UnixMicro(1589932800000000), ChangePercent: -0.026977, IsBullMarket: false},
		Kline{Open: 9511.430000, Close: 9068.650000, Low: 8815.000000, High: 9578.470000, Volume: 108928.780969, Time: time.UnixMicro(1590019200000000), ChangePercent: -0.046552, IsBullMarket: false},
		Kline{Open: 9067.510000, Close: 9170.000000, Low: 8933.520000, High: 9271.000000, Volume: 58943.131024, Time: time.UnixMicro(1590105600000000), ChangePercent: 0.011303, IsBullMarket: true},
		Kline{Open: 9170.000000, Close: 9179.150000, Low: 9070.000000, High: 9307.850000, Volume: 43526.296966, Time: time.UnixMicro(1590192000000000), ChangePercent: 0.000998, IsBullMarket: true},
		Kline{Open: 9179.010000, Close: 8720.340000, Low: 8700.000000, High: 9298.000000, Volume: 70379.866450, Time: time.UnixMicro(1590278400000000), ChangePercent: -0.049969, IsBullMarket: false},
		Kline{Open: 8718.140000, Close: 8900.350000, Low: 8642.720000, High: 8979.660000, Volume: 62833.910949, Time: time.UnixMicro(1590364800000000), ChangePercent: 0.020900, IsBullMarket: true},
		Kline{Open: 8900.350000, Close: 8841.180000, Low: 8700.000000, High: 9017.670000, Volume: 58299.770138, Time: time.UnixMicro(1590451200000000), ChangePercent: -0.006648, IsBullMarket: false},
		Kline{Open: 8841.000000, Close: 9204.070000, Low: 8811.730000, High: 9225.000000, Volume: 68910.355514, Time: time.UnixMicro(1590537600000000), ChangePercent: 0.041067, IsBullMarket: true},
		Kline{Open: 9204.070000, Close: 9575.890000, Low: 9110.000000, High: 9625.470000, Volume: 74110.787662, Time: time.UnixMicro(1590624000000000), ChangePercent: 0.040397, IsBullMarket: true},
		Kline{Open: 9575.870000, Close: 9427.070000, Low: 9330.000000, High: 9605.260000, Volume: 57374.362961, Time: time.UnixMicro(1590710400000000), ChangePercent: -0.015539, IsBullMarket: false},
		Kline{Open: 9426.600000, Close: 9697.720000, Low: 9331.230000, High: 9740.000000, Volume: 55665.272540, Time: time.UnixMicro(1590796800000000), ChangePercent: 0.028761, IsBullMarket: true},
		Kline{Open: 9697.720000, Close: 9448.270000, Low: 9381.410000, High: 9700.000000, Volume: 48333.786403, Time: time.UnixMicro(1590883200000000), ChangePercent: -0.025723, IsBullMarket: false},
		Kline{Open: 9448.270000, Close: 10200.770000, Low: 9421.670000, High: 10380.000000, Volume: 76649.126960, Time: time.UnixMicro(1590969600000000), ChangePercent: 0.079644, IsBullMarket: true},
		Kline{Open: 10202.710000, Close: 9518.040000, Low: 9266.000000, High: 10228.990000, Volume: 108970.773151, Time: time.UnixMicro(1591056000000000), ChangePercent: -0.067107, IsBullMarket: false},
		Kline{Open: 9518.020000, Close: 9666.240000, Low: 9365.210000, High: 9690.000000, Volume: 46252.644939, Time: time.UnixMicro(1591142400000000), ChangePercent: 0.015573, IsBullMarket: true},
		Kline{Open: 9666.320000, Close: 9789.060000, Low: 9450.000000, High: 9881.630000, Volume: 57456.100969, Time: time.UnixMicro(1591228800000000), ChangePercent: 0.012698, IsBullMarket: true},
		Kline{Open: 9788.140000, Close: 9621.160000, Low: 9581.000000, High: 9854.750000, Volume: 47788.050050, Time: time.UnixMicro(1591315200000000), ChangePercent: -0.017059, IsBullMarket: false},
		Kline{Open: 9621.170000, Close: 9666.300000, Low: 9531.050000, High: 9735.000000, Volume: 32752.950893, Time: time.UnixMicro(1591401600000000), ChangePercent: 0.004691, IsBullMarket: true},
		Kline{Open: 9666.850000, Close: 9746.990000, Low: 9372.460000, High: 9802.000000, Volume: 57952.848385, Time: time.UnixMicro(1591488000000000), ChangePercent: 0.008290, IsBullMarket: true},
		Kline{Open: 9746.990000, Close: 9782.010000, Low: 9633.000000, High: 9800.000000, Volume: 40664.664125, Time: time.UnixMicro(1591574400000000), ChangePercent: 0.003593, IsBullMarket: true},
		Kline{Open: 9782.000000, Close: 9772.430000, Low: 9570.000000, High: 9877.000000, Volume: 46024.001289, Time: time.UnixMicro(1591660800000000), ChangePercent: -0.000978, IsBullMarket: false},
		Kline{Open: 9772.440000, Close: 9885.000000, Low: 9704.180000, High: 9992.720000, Volume: 47130.762982, Time: time.UnixMicro(1591747200000000), ChangePercent: 0.011518, IsBullMarket: true},
		Kline{Open: 9885.220000, Close: 9280.400000, Low: 9113.000000, High: 9964.000000, Volume: 94418.984730, Time: time.UnixMicro(1591833600000000), ChangePercent: -0.061184, IsBullMarket: false},
		Kline{Open: 9278.880000, Close: 9465.130000, Low: 9232.510000, High: 9557.120000, Volume: 50119.066932, Time: time.UnixMicro(1591920000000000), ChangePercent: 0.020072, IsBullMarket: true},
		Kline{Open: 9464.960000, Close: 9473.340000, Low: 9351.000000, High: 9494.730000, Volume: 27759.784851, Time: time.UnixMicro(1592006400000000), ChangePercent: 0.000885, IsBullMarket: true},
		Kline{Open: 9473.340000, Close: 9342.100000, Low: 9245.000000, High: 9480.990000, Volume: 30055.506608, Time: time.UnixMicro(1592092800000000), ChangePercent: -0.013854, IsBullMarket: false},
		Kline{Open: 9342.100000, Close: 9426.020000, Low: 8910.450000, High: 9495.000000, Volume: 86107.924707, Time: time.UnixMicro(1592179200000000), ChangePercent: 0.008983, IsBullMarket: true},
		Kline{Open: 9426.050000, Close: 9525.590000, Low: 9373.090000, High: 9589.000000, Volume: 52052.446927, Time: time.UnixMicro(1592265600000000), ChangePercent: 0.010560, IsBullMarket: true},
		Kline{Open: 9526.970000, Close: 9465.140000, Low: 9236.610000, High: 9565.000000, Volume: 48046.411152, Time: time.UnixMicro(1592352000000000), ChangePercent: -0.006490, IsBullMarket: false},
		Kline{Open: 9465.130000, Close: 9386.320000, Low: 9280.000000, High: 9489.000000, Volume: 37381.953765, Time: time.UnixMicro(1592438400000000), ChangePercent: -0.008326, IsBullMarket: false},
		Kline{Open: 9386.320000, Close: 9310.230000, Low: 9215.790000, High: 9438.300000, Volume: 45330.983673, Time: time.UnixMicro(1592524800000000), ChangePercent: -0.008106, IsBullMarket: false},
		Kline{Open: 9310.230000, Close: 9358.950000, Low: 9170.950000, High: 9395.000000, Volume: 30329.065384, Time: time.UnixMicro(1592611200000000), ChangePercent: 0.005233, IsBullMarket: true},
		Kline{Open: 9358.950000, Close: 9294.690000, Low: 9281.540000, High: 9422.000000, Volume: 24316.926234, Time: time.UnixMicro(1592697600000000), ChangePercent: -0.006866, IsBullMarket: false},
		Kline{Open: 9294.690000, Close: 9685.690000, Low: 9277.090000, High: 9780.000000, Volume: 57895.468343, Time: time.UnixMicro(1592784000000000), ChangePercent: 0.042067, IsBullMarket: true},
		Kline{Open: 9685.690000, Close: 9624.890000, Low: 9577.030000, High: 9720.000000, Volume: 41031.029380, Time: time.UnixMicro(1592870400000000), ChangePercent: -0.006277, IsBullMarket: false},
		Kline{Open: 9624.330000, Close: 9296.490000, Low: 9208.000000, High: 9670.000000, Volume: 61571.561464, Time: time.UnixMicro(1592956800000000), ChangePercent: -0.034064, IsBullMarket: false},
		Kline{Open: 9298.330000, Close: 9249.490000, Low: 9009.690000, High: 9340.000000, Volume: 55831.619156, Time: time.UnixMicro(1593043200000000), ChangePercent: -0.005253, IsBullMarket: false},
		Kline{Open: 9249.490000, Close: 9162.210000, Low: 9045.450000, High: 9298.000000, Volume: 50292.298277, Time: time.UnixMicro(1593129600000000), ChangePercent: -0.009436, IsBullMarket: false},
		Kline{Open: 9162.210000, Close: 9012.000000, Low: 8833.000000, High: 9196.240000, Volume: 46290.930113, Time: time.UnixMicro(1593216000000000), ChangePercent: -0.016395, IsBullMarket: false},
		Kline{Open: 9012.000000, Close: 9116.350000, Low: 8948.060000, High: 9191.000000, Volume: 30688.176421, Time: time.UnixMicro(1593302400000000), ChangePercent: 0.011579, IsBullMarket: true},
		Kline{Open: 9116.160000, Close: 9192.560000, Low: 9024.670000, High: 9238.000000, Volume: 42120.293261, Time: time.UnixMicro(1593388800000000), ChangePercent: 0.008381, IsBullMarket: true},
		Kline{Open: 9192.930000, Close: 9138.550000, Low: 9064.890000, High: 9205.000000, Volume: 31463.162801, Time: time.UnixMicro(1593475200000000), ChangePercent: -0.005915, IsBullMarket: false},
		Kline{Open: 9138.080000, Close: 9232.000000, Low: 9080.100000, High: 9292.000000, Volume: 38488.528699, Time: time.UnixMicro(1593561600000000), ChangePercent: 0.010278, IsBullMarket: true},
		Kline{Open: 9231.990000, Close: 9086.540000, Low: 8940.000000, High: 9261.960000, Volume: 45725.168076, Time: time.UnixMicro(1593648000000000), ChangePercent: -0.015755, IsBullMarket: false},
		Kline{Open: 9086.540000, Close: 9058.260000, Low: 9037.470000, High: 9125.000000, Volume: 28943.420177, Time: time.UnixMicro(1593734400000000), ChangePercent: -0.003112, IsBullMarket: false},
		Kline{Open: 9057.790000, Close: 9135.460000, Low: 9040.040000, High: 9190.000000, Volume: 26441.968484, Time: time.UnixMicro(1593820800000000), ChangePercent: 0.008575, IsBullMarket: true},
		Kline{Open: 9135.000000, Close: 9069.410000, Low: 8893.030000, High: 9145.240000, Volume: 34073.653627, Time: time.UnixMicro(1593907200000000), ChangePercent: -0.007180, IsBullMarket: false},
		Kline{Open: 9069.410000, Close: 9344.200000, Low: 9055.920000, High: 9375.000000, Volume: 54463.132277, Time: time.UnixMicro(1593993600000000), ChangePercent: 0.030299, IsBullMarket: true},
		Kline{Open: 9342.470000, Close: 9257.390000, Low: 9203.000000, High: 9379.420000, Volume: 34587.336678, Time: time.UnixMicro(1594080000000000), ChangePercent: -0.009107, IsBullMarket: false},
		Kline{Open: 9257.400000, Close: 9436.060000, Low: 9231.000000, High: 9470.000000, Volume: 56140.517781, Time: time.UnixMicro(1594166400000000), ChangePercent: 0.019299, IsBullMarket: true},
		Kline{Open: 9436.060000, Close: 9232.430000, Low: 9160.000000, High: 9440.790000, Volume: 48044.450645, Time: time.UnixMicro(1594252800000000), ChangePercent: -0.021580, IsBullMarket: false},
		Kline{Open: 9232.420000, Close: 9288.340000, Low: 9125.000000, High: 9317.480000, Volume: 38295.494006, Time: time.UnixMicro(1594339200000000), ChangePercent: 0.006057, IsBullMarket: true},
		Kline{Open: 9288.340000, Close: 9234.030000, Low: 9178.250000, High: 9299.280000, Volume: 22561.366000, Time: time.UnixMicro(1594425600000000), ChangePercent: -0.005847, IsBullMarket: false},
		Kline{Open: 9234.020000, Close: 9302.750000, Low: 9157.500000, High: 9345.000000, Volume: 30872.702286, Time: time.UnixMicro(1594512000000000), ChangePercent: 0.007443, IsBullMarket: true},
		Kline{Open: 9303.310000, Close: 9242.620000, Low: 9200.890000, High: 9343.820000, Volume: 42740.069115, Time: time.UnixMicro(1594598400000000), ChangePercent: -0.006523, IsBullMarket: false},
		Kline{Open: 9242.610000, Close: 9255.850000, Low: 9113.000000, High: 9279.540000, Volume: 45772.552509, Time: time.UnixMicro(1594684800000000), ChangePercent: 0.001432, IsBullMarket: true},
		Kline{Open: 9255.850000, Close: 9197.600000, Low: 9160.570000, High: 9276.490000, Volume: 39053.579665, Time: time.UnixMicro(1594771200000000), ChangePercent: -0.006293, IsBullMarket: false},
		Kline{Open: 9197.600000, Close: 9133.720000, Low: 9047.250000, High: 9226.150000, Volume: 43375.571191, Time: time.UnixMicro(1594857600000000), ChangePercent: -0.006945, IsBullMarket: false},
		Kline{Open: 9133.720000, Close: 9154.320000, Low: 9089.810000, High: 9186.830000, Volume: 28054.358741, Time: time.UnixMicro(1594944000000000), ChangePercent: 0.002255, IsBullMarket: true},
		Kline{Open: 9154.310000, Close: 9170.280000, Low: 9121.100000, High: 9219.300000, Volume: 22554.541457, Time: time.UnixMicro(1595030400000000), ChangePercent: 0.001745, IsBullMarket: true},
		Kline{Open: 9170.300000, Close: 9208.990000, Low: 9101.350000, High: 9232.270000, Volume: 26052.019417, Time: time.UnixMicro(1595116800000000), ChangePercent: 0.004219, IsBullMarket: true},
		Kline{Open: 9208.990000, Close: 9160.780000, Low: 9131.000000, High: 9221.520000, Volume: 35458.764082, Time: time.UnixMicro(1595203200000000), ChangePercent: -0.005235, IsBullMarket: false},
		Kline{Open: 9160.780000, Close: 9390.000000, Low: 9152.800000, High: 9437.730000, Volume: 60413.582486, Time: time.UnixMicro(1595289600000000), ChangePercent: 0.025022, IsBullMarket: true},
		Kline{Open: 9390.000000, Close: 9518.160000, Low: 9261.000000, High: 9544.000000, Volume: 48815.004107, Time: time.UnixMicro(1595376000000000), ChangePercent: 0.013649, IsBullMarket: true},
		Kline{Open: 9518.160000, Close: 9603.270000, Low: 9440.330000, High: 9664.000000, Volume: 51856.233500, Time: time.UnixMicro(1595462400000000), ChangePercent: 0.008942, IsBullMarket: true},
		Kline{Open: 9603.270000, Close: 9537.800000, Low: 9463.440000, High: 9637.000000, Volume: 43931.136205, Time: time.UnixMicro(1595548800000000), ChangePercent: -0.006817, IsBullMarket: false},
		Kline{Open: 9538.100000, Close: 9700.420000, Low: 9513.000000, High: 9732.900000, Volume: 40679.545416, Time: time.UnixMicro(1595635200000000), ChangePercent: 0.017018, IsBullMarket: true},
		Kline{Open: 9700.420000, Close: 9931.540000, Low: 9650.000000, High: 10111.000000, Volume: 65279.269319, Time: time.UnixMicro(1595721600000000), ChangePercent: 0.023826, IsBullMarket: true},
		Kline{Open: 9931.540000, Close: 11029.960000, Low: 9917.210000, High: 11394.860000, Volume: 150188.933144, Time: time.UnixMicro(1595808000000000), ChangePercent: 0.110599, IsBullMarket: true},
		Kline{Open: 11029.960000, Close: 10906.270000, Low: 10565.000000, High: 11242.230000, Volume: 97267.734187, Time: time.UnixMicro(1595894400000000), ChangePercent: -0.011214, IsBullMarket: false},
		Kline{Open: 10906.270000, Close: 11100.530000, Low: 10812.000000, High: 11342.820000, Volume: 76838.094233, Time: time.UnixMicro(1595980800000000), ChangePercent: 0.017812, IsBullMarket: true},
		Kline{Open: 11100.520000, Close: 11099.610000, Low: 10831.000000, High: 11170.000000, Volume: 60794.826456, Time: time.UnixMicro(1596067200000000), ChangePercent: -0.000082, IsBullMarket: false},
		Kline{Open: 11099.790000, Close: 11335.460000, Low: 10960.000000, High: 11444.000000, Volume: 70063.660974, Time: time.UnixMicro(1596153600000000), ChangePercent: 0.021232, IsBullMarket: true},
		Kline{Open: 11335.460000, Close: 11801.170000, Low: 11220.000000, High: 11861.000000, Volume: 85087.485126, Time: time.UnixMicro(1596240000000000), ChangePercent: 0.041084, IsBullMarket: true},
		Kline{Open: 11801.170000, Close: 11071.350000, Low: 10518.500000, High: 12123.460000, Volume: 97553.077604, Time: time.UnixMicro(1596326400000000), ChangePercent: -0.061843, IsBullMarket: false},
		Kline{Open: 11071.360000, Close: 11219.810000, Low: 10936.000000, High: 11473.000000, Volume: 56931.841475, Time: time.UnixMicro(1596412800000000), ChangePercent: 0.013408, IsBullMarket: true},
		Kline{Open: 11219.680000, Close: 11191.970000, Low: 11000.000000, High: 11414.980000, Volume: 58629.113709, Time: time.UnixMicro(1596499200000000), ChangePercent: -0.002470, IsBullMarket: false},
		Kline{Open: 11191.990000, Close: 11744.910000, Low: 11093.000000, High: 11780.930000, Volume: 74970.256852, Time: time.UnixMicro(1596585600000000), ChangePercent: 0.049403, IsBullMarket: true},
		Kline{Open: 11744.910000, Close: 11762.460000, Low: 11562.500000, High: 11900.000000, Volume: 63529.085020, Time: time.UnixMicro(1596672000000000), ChangePercent: 0.001494, IsBullMarket: true},
		Kline{Open: 11762.470000, Close: 11594.230000, Low: 11322.000000, High: 11909.940000, Volume: 65755.926022, Time: time.UnixMicro(1596758400000000), ChangePercent: -0.014303, IsBullMarket: false},
		Kline{Open: 11594.360000, Close: 11761.410000, Low: 11512.000000, High: 11808.270000, Volume: 41858.161040, Time: time.UnixMicro(1596844800000000), ChangePercent: 0.014408, IsBullMarket: true},
		Kline{Open: 11761.020000, Close: 11681.680000, Low: 11521.970000, High: 11797.110000, Volume: 41493.067342, Time: time.UnixMicro(1596931200000000), ChangePercent: -0.006746, IsBullMarket: false},
		Kline{Open: 11681.690000, Close: 11892.920000, Low: 11450.000000, High: 12067.350000, Volume: 84952.337887, Time: time.UnixMicro(1597017600000000), ChangePercent: 0.018082, IsBullMarket: true},
		Kline{Open: 11892.900000, Close: 11392.080000, Low: 11125.000000, High: 11935.000000, Volume: 90748.284634, Time: time.UnixMicro(1597104000000000), ChangePercent: -0.042111, IsBullMarket: false},
		Kline{Open: 11392.090000, Close: 11564.330000, Low: 11150.000000, High: 11617.520000, Volume: 64909.613644, Time: time.UnixMicro(1597190400000000), ChangePercent: 0.015119, IsBullMarket: true},
		Kline{Open: 11564.340000, Close: 11780.000000, Low: 11270.360000, High: 11792.960000, Volume: 70132.491502, Time: time.UnixMicro(1597276800000000), ChangePercent: 0.018649, IsBullMarket: true},
		Kline{Open: 11779.770000, Close: 11760.540000, Low: 11634.030000, High: 11850.000000, Volume: 59818.852697, Time: time.UnixMicro(1597363200000000), ChangePercent: -0.001632, IsBullMarket: false},
		Kline{Open: 11760.550000, Close: 11852.400000, Low: 11680.000000, High: 11980.000000, Volume: 56237.905050, Time: time.UnixMicro(1597449600000000), ChangePercent: 0.007810, IsBullMarket: true},
		Kline{Open: 11852.400000, Close: 11911.000000, Low: 11686.000000, High: 11931.720000, Volume: 41368.236906, Time: time.UnixMicro(1597536000000000), ChangePercent: 0.004944, IsBullMarket: true},
		Kline{Open: 11910.990000, Close: 12281.130000, Low: 11769.780000, High: 12468.000000, Volume: 84734.211540, Time: time.UnixMicro(1597622400000000), ChangePercent: 0.031076, IsBullMarket: true},
		Kline{Open: 12281.150000, Close: 11945.010000, Low: 11817.930000, High: 12387.770000, Volume: 75923.835527, Time: time.UnixMicro(1597708800000000), ChangePercent: -0.027370, IsBullMarket: false},
		Kline{Open: 11945.100000, Close: 11754.590000, Low: 11561.000000, High: 12020.080000, Volume: 73940.169606, Time: time.UnixMicro(1597795200000000), ChangePercent: -0.015949, IsBullMarket: false},
		Kline{Open: 11754.380000, Close: 11853.550000, Low: 11668.000000, High: 11888.000000, Volume: 46085.254351, Time: time.UnixMicro(1597881600000000), ChangePercent: 0.008437, IsBullMarket: true},
		Kline{Open: 11853.540000, Close: 11531.340000, Low: 11485.810000, High: 11878.000000, Volume: 64448.306142, Time: time.UnixMicro(1597968000000000), ChangePercent: -0.027182, IsBullMarket: false},
		Kline{Open: 11531.230000, Close: 11662.960000, Low: 11376.810000, High: 11686.000000, Volume: 43678.701646, Time: time.UnixMicro(1598054400000000), ChangePercent: 0.011424, IsBullMarket: true},
		Kline{Open: 11663.510000, Close: 11648.130000, Low: 11514.130000, High: 11718.070000, Volume: 37900.004690, Time: time.UnixMicro(1598140800000000), ChangePercent: -0.001319, IsBullMarket: false},
		Kline{Open: 11648.120000, Close: 11748.200000, Low: 11585.090000, High: 11824.900000, Volume: 46212.391867, Time: time.UnixMicro(1598227200000000), ChangePercent: 0.008592, IsBullMarket: true},
		Kline{Open: 11748.190000, Close: 11318.420000, Low: 11117.640000, High: 11767.850000, Volume: 69590.923272, Time: time.UnixMicro(1598313600000000), ChangePercent: -0.036582, IsBullMarket: false},
		Kline{Open: 11318.420000, Close: 11461.430000, Low: 11244.000000, High: 11539.320000, Volume: 53998.231231, Time: time.UnixMicro(1598400000000000), ChangePercent: 0.012635, IsBullMarket: true},
		Kline{Open: 11461.420000, Close: 11330.380000, Low: 11125.000000, High: 11592.200000, Volume: 63246.036383, Time: time.UnixMicro(1598486400000000), ChangePercent: -0.011433, IsBullMarket: false},
		Kline{Open: 11330.380000, Close: 11526.910000, Low: 11276.890000, High: 11542.650000, Volume: 45953.908365, Time: time.UnixMicro(1598572800000000), ChangePercent: 0.017345, IsBullMarket: true},
		Kline{Open: 11526.900000, Close: 11465.840000, Low: 11417.040000, High: 11580.020000, Volume: 32973.799200, Time: time.UnixMicro(1598659200000000), ChangePercent: -0.005297, IsBullMarket: false},
		Kline{Open: 11465.840000, Close: 11711.160000, Low: 11458.000000, High: 11719.000000, Volume: 43177.879054, Time: time.UnixMicro(1598745600000000), ChangePercent: 0.021396, IsBullMarket: true},
		Kline{Open: 11711.170000, Close: 11649.510000, Low: 11570.000000, High: 11800.770000, Volume: 55353.617744, Time: time.UnixMicro(1598832000000000), ChangePercent: -0.005265, IsBullMarket: false},
		Kline{Open: 11649.510000, Close: 11921.970000, Low: 11515.000000, High: 12050.850000, Volume: 78148.193668, Time: time.UnixMicro(1598918400000000), ChangePercent: 0.023388, IsBullMarket: true},
		Kline{Open: 11921.970000, Close: 11388.540000, Low: 11160.100000, High: 11954.570000, Volume: 87221.845602, Time: time.UnixMicro(1599004800000000), ChangePercent: -0.044743, IsBullMarket: false},
		Kline{Open: 11388.540000, Close: 10140.850000, Low: 9960.800000, High: 11462.600000, Volume: 121950.106015, Time: time.UnixMicro(1599091200000000), ChangePercent: -0.109557, IsBullMarket: false},
		Kline{Open: 10138.290000, Close: 10446.250000, Low: 9875.500000, High: 10627.050000, Volume: 92733.599113, Time: time.UnixMicro(1599177600000000), ChangePercent: 0.030376, IsBullMarket: true},
		Kline{Open: 10446.250000, Close: 10166.690000, Low: 9825.000000, High: 10565.680000, Volume: 90001.605568, Time: time.UnixMicro(1599264000000000), ChangePercent: -0.026762, IsBullMarket: false},
		Kline{Open: 10166.690000, Close: 10256.200000, Low: 9994.860000, High: 10347.140000, Volume: 56368.788815, Time: time.UnixMicro(1599350400000000), ChangePercent: 0.008804, IsBullMarket: true},
		Kline{Open: 10255.890000, Close: 10373.440000, Low: 9875.000000, High: 10410.750000, Volume: 62620.230676, Time: time.UnixMicro(1599436800000000), ChangePercent: 0.011462, IsBullMarket: true},
		Kline{Open: 10373.450000, Close: 10126.650000, Low: 9850.000000, High: 10438.000000, Volume: 73491.878418, Time: time.UnixMicro(1599523200000000), ChangePercent: -0.023792, IsBullMarket: false},
		Kline{Open: 10126.660000, Close: 10219.200000, Low: 9981.010000, High: 10343.000000, Volume: 49347.113776, Time: time.UnixMicro(1599609600000000), ChangePercent: 0.009138, IsBullMarket: true},
		Kline{Open: 10219.290000, Close: 10336.870000, Low: 10070.830000, High: 10483.350000, Volume: 58253.753750, Time: time.UnixMicro(1599696000000000), ChangePercent: 0.011506, IsBullMarket: true},
		Kline{Open: 10336.860000, Close: 10387.890000, Low: 10200.000000, High: 10397.600000, Volume: 43830.254467, Time: time.UnixMicro(1599782400000000), ChangePercent: 0.004937, IsBullMarket: true},
		Kline{Open: 10387.890000, Close: 10440.920000, Low: 10269.250000, High: 10477.970000, Volume: 35379.153096, Time: time.UnixMicro(1599868800000000), ChangePercent: 0.005105, IsBullMarket: true},
		Kline{Open: 10440.670000, Close: 10332.830000, Low: 10200.000000, High: 10580.110000, Volume: 43837.609865, Time: time.UnixMicro(1599955200000000), ChangePercent: -0.010329, IsBullMarket: false},
		Kline{Open: 10332.840000, Close: 10671.770000, Low: 10212.340000, High: 10750.000000, Volume: 67059.291361, Time: time.UnixMicro(1600041600000000), ChangePercent: 0.032801, IsBullMarket: true},
		Kline{Open: 10671.770000, Close: 10785.310000, Low: 10606.480000, High: 10930.040000, Volume: 61822.452786, Time: time.UnixMicro(1600128000000000), ChangePercent: 0.010639, IsBullMarket: true},
		Kline{Open: 10785.230000, Close: 10954.010000, Low: 10661.220000, High: 11093.000000, Volume: 64991.512440, Time: time.UnixMicro(1600214400000000), ChangePercent: 0.015649, IsBullMarket: true},
		Kline{Open: 10954.010000, Close: 10939.990000, Low: 10745.830000, High: 11045.460000, Volume: 55601.614529, Time: time.UnixMicro(1600300800000000), ChangePercent: -0.001280, IsBullMarket: false},
		Kline{Open: 10940.000000, Close: 10933.390000, Low: 10812.840000, High: 11038.030000, Volume: 47266.728275, Time: time.UnixMicro(1600387200000000), ChangePercent: -0.000604, IsBullMarket: false},
		Kline{Open: 10933.400000, Close: 11080.650000, Low: 10887.370000, High: 11179.790000, Volume: 38440.036858, Time: time.UnixMicro(1600473600000000), ChangePercent: 0.013468, IsBullMarket: true},
		Kline{Open: 11080.640000, Close: 10920.280000, Low: 10723.000000, High: 11080.640000, Volume: 39157.922565, Time: time.UnixMicro(1600560000000000), ChangePercent: -0.014472, IsBullMarket: false},
		Kline{Open: 10920.280000, Close: 10417.220000, Low: 10296.350000, High: 10988.860000, Volume: 70683.431179, Time: time.UnixMicro(1600646400000000), ChangePercent: -0.046067, IsBullMarket: false},
		Kline{Open: 10417.220000, Close: 10529.610000, Low: 10353.000000, High: 10572.710000, Volume: 43991.235476, Time: time.UnixMicro(1600732800000000), ChangePercent: 0.010789, IsBullMarket: true},
		Kline{Open: 10529.610000, Close: 10241.460000, Low: 10136.820000, High: 10537.150000, Volume: 51876.568079, Time: time.UnixMicro(1600819200000000), ChangePercent: -0.027366, IsBullMarket: false},
		Kline{Open: 10241.460000, Close: 10736.320000, Low: 10190.930000, High: 10795.240000, Volume: 57676.619427, Time: time.UnixMicro(1600905600000000), ChangePercent: 0.048319, IsBullMarket: true},
		Kline{Open: 10736.330000, Close: 10686.670000, Low: 10556.240000, High: 10760.530000, Volume: 48101.117008, Time: time.UnixMicro(1600992000000000), ChangePercent: -0.004625, IsBullMarket: false},
		Kline{Open: 10686.570000, Close: 10728.600000, Low: 10644.680000, High: 10820.940000, Volume: 28420.836659, Time: time.UnixMicro(1601078400000000), ChangePercent: 0.003933, IsBullMarket: true},
		Kline{Open: 10728.590000, Close: 10774.250000, Low: 10594.820000, High: 10799.000000, Volume: 30549.483253, Time: time.UnixMicro(1601164800000000), ChangePercent: 0.004256, IsBullMarket: true},
		Kline{Open: 10774.260000, Close: 10696.120000, Low: 10626.000000, High: 10950.000000, Volume: 50095.251734, Time: time.UnixMicro(1601251200000000), ChangePercent: -0.007252, IsBullMarket: false},
		Kline{Open: 10696.110000, Close: 10840.480000, Low: 10635.870000, High: 10867.540000, Volume: 41874.898399, Time: time.UnixMicro(1601337600000000), ChangePercent: 0.013497, IsBullMarket: true},
		Kline{Open: 10840.580000, Close: 10776.590000, Low: 10665.130000, High: 10849.340000, Volume: 39596.027322, Time: time.UnixMicro(1601424000000000), ChangePercent: -0.005903, IsBullMarket: false},
		Kline{Open: 10776.590000, Close: 10619.130000, Low: 10437.000000, High: 10920.000000, Volume: 60866.332893, Time: time.UnixMicro(1601510400000000), ChangePercent: -0.014611, IsBullMarket: false},
		Kline{Open: 10619.130000, Close: 10570.400000, Low: 10374.000000, High: 10664.640000, Volume: 50130.393705, Time: time.UnixMicro(1601596800000000), ChangePercent: -0.004589, IsBullMarket: false},
		Kline{Open: 10570.400000, Close: 10542.060000, Low: 10496.460000, High: 10603.560000, Volume: 22298.221341, Time: time.UnixMicro(1601683200000000), ChangePercent: -0.002681, IsBullMarket: false},
		Kline{Open: 10542.070000, Close: 10666.630000, Low: 10517.870000, High: 10696.870000, Volume: 23212.001595, Time: time.UnixMicro(1601769600000000), ChangePercent: 0.011816, IsBullMarket: true},
		Kline{Open: 10666.620000, Close: 10792.210000, Low: 10615.640000, High: 10798.000000, Volume: 34025.761653, Time: time.UnixMicro(1601856000000000), ChangePercent: 0.011774, IsBullMarket: true},
		Kline{Open: 10792.200000, Close: 10599.660000, Low: 10525.000000, High: 10800.000000, Volume: 48674.740471, Time: time.UnixMicro(1601942400000000), ChangePercent: -0.017841, IsBullMarket: false},
		Kline{Open: 10599.650000, Close: 10666.390000, Low: 10546.170000, High: 10681.870000, Volume: 32811.990279, Time: time.UnixMicro(1602028800000000), ChangePercent: 0.006296, IsBullMarket: true},
		Kline{Open: 10666.400000, Close: 10925.570000, Low: 10530.410000, High: 10950.000000, Volume: 51959.691572, Time: time.UnixMicro(1602115200000000), ChangePercent: 0.024298, IsBullMarket: true},
		Kline{Open: 10925.440000, Close: 11050.640000, Low: 10829.000000, High: 11104.640000, Volume: 48240.073237, Time: time.UnixMicro(1602201600000000), ChangePercent: 0.011459, IsBullMarket: true},
		Kline{Open: 11050.640000, Close: 11293.220000, Low: 11050.510000, High: 11491.000000, Volume: 43648.036943, Time: time.UnixMicro(1602288000000000), ChangePercent: 0.021952, IsBullMarket: true},
		Kline{Open: 11293.220000, Close: 11369.020000, Low: 11221.000000, High: 11445.000000, Volume: 29043.851339, Time: time.UnixMicro(1602374400000000), ChangePercent: 0.006712, IsBullMarket: true},
		Kline{Open: 11369.020000, Close: 11528.250000, Low: 11172.000000, High: 11720.010000, Volume: 52825.283710, Time: time.UnixMicro(1602460800000000), ChangePercent: 0.014006, IsBullMarket: true},
		Kline{Open: 11528.240000, Close: 11420.560000, Low: 11300.000000, High: 11557.000000, Volume: 42205.283709, Time: time.UnixMicro(1602547200000000), ChangePercent: -0.009341, IsBullMarket: false},
		Kline{Open: 11420.570000, Close: 11417.890000, Low: 11280.000000, High: 11547.980000, Volume: 41415.106015, Time: time.UnixMicro(1602633600000000), ChangePercent: -0.000235, IsBullMarket: false},
		Kline{Open: 11417.890000, Close: 11505.120000, Low: 11250.830000, High: 11617.340000, Volume: 48760.717679, Time: time.UnixMicro(1602720000000000), ChangePercent: 0.007640, IsBullMarket: true},
		Kline{Open: 11505.130000, Close: 11319.320000, Low: 11200.000000, High: 11541.150000, Volume: 48797.749502, Time: time.UnixMicro(1602806400000000), ChangePercent: -0.016150, IsBullMarket: false},
		Kline{Open: 11319.240000, Close: 11360.200000, Low: 11255.000000, High: 11402.420000, Volume: 22368.915241, Time: time.UnixMicro(1602892800000000), ChangePercent: 0.003619, IsBullMarket: true},
		Kline{Open: 11360.310000, Close: 11503.140000, Low: 11346.220000, High: 11505.000000, Volume: 23284.041191, Time: time.UnixMicro(1602979200000000), ChangePercent: 0.012573, IsBullMarket: true},
		Kline{Open: 11503.140000, Close: 11751.470000, Low: 11407.960000, High: 11823.990000, Volume: 47414.534692, Time: time.UnixMicro(1603065600000000), ChangePercent: 0.021588, IsBullMarket: true},
		Kline{Open: 11751.460000, Close: 11909.990000, Low: 11677.590000, High: 12038.380000, Volume: 62134.750663, Time: time.UnixMicro(1603152000000000), ChangePercent: 0.013490, IsBullMarket: true},
		Kline{Open: 11910.000000, Close: 12780.960000, Low: 11886.950000, High: 13217.680000, Volume: 114584.456767, Time: time.UnixMicro(1603238400000000), ChangePercent: 0.073128, IsBullMarket: true},
		Kline{Open: 12780.750000, Close: 12968.520000, Low: 12678.080000, High: 13185.000000, Volume: 70038.824144, Time: time.UnixMicro(1603324800000000), ChangePercent: 0.014692, IsBullMarket: true},
		Kline{Open: 12968.840000, Close: 12923.070000, Low: 12720.080000, High: 13027.690000, Volume: 50386.999841, Time: time.UnixMicro(1603411200000000), ChangePercent: -0.003529, IsBullMarket: false},
		Kline{Open: 12923.060000, Close: 13111.730000, Low: 12870.000000, High: 13166.730000, Volume: 35952.209070, Time: time.UnixMicro(1603497600000000), ChangePercent: 0.014599, IsBullMarket: true},
		Kline{Open: 13111.730000, Close: 13028.830000, Low: 12888.000000, High: 13350.000000, Volume: 38481.579504, Time: time.UnixMicro(1603584000000000), ChangePercent: -0.006323, IsBullMarket: false},
		Kline{Open: 13029.640000, Close: 13052.190000, Low: 12765.000000, High: 13238.810000, Volume: 60951.672986, Time: time.UnixMicro(1603670400000000), ChangePercent: 0.001731, IsBullMarket: true},
		Kline{Open: 13052.150000, Close: 13636.170000, Low: 13019.870000, High: 13789.290000, Volume: 80811.019450, Time: time.UnixMicro(1603756800000000), ChangePercent: 0.044745, IsBullMarket: true},
		Kline{Open: 13636.160000, Close: 13266.400000, Low: 12888.000000, High: 13859.480000, Volume: 94440.561226, Time: time.UnixMicro(1603843200000000), ChangePercent: -0.027116, IsBullMarket: false},
		Kline{Open: 13266.400000, Close: 13455.700000, Low: 12920.770000, High: 13642.910000, Volume: 74872.602132, Time: time.UnixMicro(1603929600000000), ChangePercent: 0.014269, IsBullMarket: true},
		Kline{Open: 13455.690000, Close: 13560.100000, Low: 13115.000000, High: 13669.980000, Volume: 70657.778881, Time: time.UnixMicro(1604016000000000), ChangePercent: 0.007760, IsBullMarket: true},
		Kline{Open: 13560.100000, Close: 13791.000000, Low: 13411.500000, High: 14100.000000, Volume: 67339.238515, Time: time.UnixMicro(1604102400000000), ChangePercent: 0.017028, IsBullMarket: true},
		Kline{Open: 13791.000000, Close: 13761.500000, Low: 13603.000000, High: 13895.000000, Volume: 36285.648526, Time: time.UnixMicro(1604188800000000), ChangePercent: -0.002139, IsBullMarket: false},
		Kline{Open: 13761.490000, Close: 13549.370000, Low: 13195.050000, High: 13830.000000, Volume: 64566.421908, Time: time.UnixMicro(1604275200000000), ChangePercent: -0.015414, IsBullMarket: false},
		Kline{Open: 13549.630000, Close: 14023.530000, Low: 13284.990000, High: 14066.110000, Volume: 74115.630787, Time: time.UnixMicro(1604361600000000), ChangePercent: 0.034975, IsBullMarket: true},
		Kline{Open: 14023.530000, Close: 14144.010000, Low: 13525.000000, High: 14259.000000, Volume: 93016.988262, Time: time.UnixMicro(1604448000000000), ChangePercent: 0.008591, IsBullMarket: true},
		Kline{Open: 14144.010000, Close: 15590.020000, Low: 14093.560000, High: 15750.000000, Volume: 143741.522673, Time: time.UnixMicro(1604534400000000), ChangePercent: 0.102235, IsBullMarket: true},
		Kline{Open: 15590.020000, Close: 15579.920000, Low: 15166.000000, High: 15960.000000, Volume: 122618.197695, Time: time.UnixMicro(1604620800000000), ChangePercent: -0.000648, IsBullMarket: false},
		Kline{Open: 15579.930000, Close: 14818.300000, Low: 14344.220000, High: 15753.520000, Volume: 101431.206553, Time: time.UnixMicro(1604707200000000), ChangePercent: -0.048885, IsBullMarket: false},
		Kline{Open: 14818.300000, Close: 15475.100000, Low: 14703.880000, High: 15650.000000, Volume: 65547.178574, Time: time.UnixMicro(1604793600000000), ChangePercent: 0.044324, IsBullMarket: true},
		Kline{Open: 15475.100000, Close: 15328.410000, Low: 14805.540000, High: 15840.000000, Volume: 108976.334134, Time: time.UnixMicro(1604880000000000), ChangePercent: -0.009479, IsBullMarket: false},
		Kline{Open: 15328.410000, Close: 15297.210000, Low: 15072.460000, High: 15460.000000, Volume: 61681.919606, Time: time.UnixMicro(1604966400000000), ChangePercent: -0.002035, IsBullMarket: false},
		Kline{Open: 15297.210000, Close: 15684.240000, Low: 15272.680000, High: 15965.000000, Volume: 78469.746458, Time: time.UnixMicro(1605052800000000), ChangePercent: 0.025301, IsBullMarket: true},
		Kline{Open: 15684.250000, Close: 16291.860000, Low: 15440.640000, High: 16340.700000, Volume: 102196.356592, Time: time.UnixMicro(1605139200000000), ChangePercent: 0.038740, IsBullMarket: true},
		Kline{Open: 16291.850000, Close: 16320.700000, Low: 15952.350000, High: 16480.000000, Volume: 75691.881014, Time: time.UnixMicro(1605225600000000), ChangePercent: 0.001771, IsBullMarket: true},
		Kline{Open: 16320.040000, Close: 16070.450000, Low: 15670.000000, High: 16326.990000, Volume: 59116.347179, Time: time.UnixMicro(1605312000000000), ChangePercent: -0.015293, IsBullMarket: false},
		Kline{Open: 16069.560000, Close: 15957.000000, Low: 15774.720000, High: 16180.000000, Volume: 43596.841513, Time: time.UnixMicro(1605398400000000), ChangePercent: -0.007005, IsBullMarket: false},
		Kline{Open: 15957.000000, Close: 16713.570000, Low: 15864.000000, High: 16880.000000, Volume: 81300.675924, Time: time.UnixMicro(1605484800000000), ChangePercent: 0.047413, IsBullMarket: true},
		Kline{Open: 16713.080000, Close: 17659.380000, Low: 16538.000000, High: 17858.820000, Volume: 115221.403102, Time: time.UnixMicro(1605571200000000), ChangePercent: 0.056620, IsBullMarket: true},
		Kline{Open: 17659.380000, Close: 17776.120000, Low: 17214.450000, High: 18476.930000, Volume: 149019.788134, Time: time.UnixMicro(1605657600000000), ChangePercent: 0.006611, IsBullMarket: true},
		Kline{Open: 17777.750000, Close: 17802.820000, Low: 17335.650000, High: 18179.800000, Volume: 93009.561008, Time: time.UnixMicro(1605744000000000), ChangePercent: 0.001410, IsBullMarket: true},
		Kline{Open: 17802.810000, Close: 18655.670000, Low: 17740.040000, High: 18815.220000, Volume: 88423.018489, Time: time.UnixMicro(1605830400000000), ChangePercent: 0.047906, IsBullMarket: true},
		Kline{Open: 18655.660000, Close: 18703.800000, Low: 18308.580000, High: 18965.900000, Volume: 75577.458394, Time: time.UnixMicro(1605916800000000), ChangePercent: 0.002580, IsBullMarket: true},
		Kline{Open: 18703.800000, Close: 18414.430000, Low: 17610.860000, High: 18750.000000, Volume: 81645.737778, Time: time.UnixMicro(1606003200000000), ChangePercent: -0.015471, IsBullMarket: false},
		Kline{Open: 18413.880000, Close: 18368.000000, Low: 18000.000000, High: 18766.000000, Volume: 82961.506093, Time: time.UnixMicro(1606089600000000), ChangePercent: -0.002492, IsBullMarket: false},
		Kline{Open: 18368.010000, Close: 19160.010000, Low: 18018.000000, High: 19418.970000, Volume: 113581.509241, Time: time.UnixMicro(1606176000000000), ChangePercent: 0.043118, IsBullMarket: true},
		Kline{Open: 19160.000000, Close: 18719.110000, Low: 18500.270000, High: 19484.210000, Volume: 93266.576887, Time: time.UnixMicro(1606262400000000), ChangePercent: -0.023011, IsBullMarket: false},
		Kline{Open: 18718.830000, Close: 17149.470000, Low: 16188.000000, High: 18915.030000, Volume: 181005.246693, Time: time.UnixMicro(1606348800000000), ChangePercent: -0.083839, IsBullMarket: false},
		Kline{Open: 17149.470000, Close: 17139.520000, Low: 16438.080000, High: 17457.620000, Volume: 85297.024787, Time: time.UnixMicro(1606435200000000), ChangePercent: -0.000580, IsBullMarket: false},
		Kline{Open: 17139.530000, Close: 17719.850000, Low: 16865.560000, High: 17880.490000, Volume: 64910.699970, Time: time.UnixMicro(1606521600000000), ChangePercent: 0.033859, IsBullMarket: true},
		Kline{Open: 17719.840000, Close: 18184.990000, Low: 17517.000000, High: 18360.050000, Volume: 55329.016303, Time: time.UnixMicro(1606608000000000), ChangePercent: 0.026250, IsBullMarket: true},
		Kline{Open: 18185.000000, Close: 19695.870000, Low: 18184.990000, High: 19863.160000, Volume: 115463.466888, Time: time.UnixMicro(1606694400000000), ChangePercent: 0.083083, IsBullMarket: true},
		Kline{Open: 19695.870000, Close: 18764.960000, Low: 18001.120000, High: 19888.000000, Volume: 127698.762652, Time: time.UnixMicro(1606780800000000), ChangePercent: -0.047264, IsBullMarket: false},
		Kline{Open: 18764.960000, Close: 19204.090000, Low: 18330.000000, High: 19342.000000, Volume: 75911.013478, Time: time.UnixMicro(1606867200000000), ChangePercent: 0.023402, IsBullMarket: true},
		Kline{Open: 19204.080000, Close: 19421.900000, Low: 18867.200000, High: 19598.000000, Volume: 66689.391279, Time: time.UnixMicro(1606953600000000), ChangePercent: 0.011342, IsBullMarket: true},
		Kline{Open: 19422.340000, Close: 18650.520000, Low: 18565.310000, High: 19527.000000, Volume: 71283.668200, Time: time.UnixMicro(1607040000000000), ChangePercent: -0.039739, IsBullMarket: false},
		Kline{Open: 18650.510000, Close: 19147.660000, Low: 18500.000000, High: 19177.000000, Volume: 42922.748573, Time: time.UnixMicro(1607126400000000), ChangePercent: 0.026656, IsBullMarket: true},
		Kline{Open: 19147.660000, Close: 19359.400000, Low: 18857.000000, High: 19420.000000, Volume: 37043.091861, Time: time.UnixMicro(1607212800000000), ChangePercent: 0.011058, IsBullMarket: true},
		Kline{Open: 19358.670000, Close: 19166.900000, Low: 18902.880000, High: 19420.910000, Volume: 41372.296293, Time: time.UnixMicro(1607299200000000), ChangePercent: -0.009906, IsBullMarket: false},
		Kline{Open: 19166.900000, Close: 18324.110000, Low: 18200.000000, High: 19294.840000, Volume: 61626.947614, Time: time.UnixMicro(1607385600000000), ChangePercent: -0.043971, IsBullMarket: false},
		Kline{Open: 18324.110000, Close: 18541.280000, Low: 17650.000000, High: 18639.570000, Volume: 79585.553801, Time: time.UnixMicro(1607472000000000), ChangePercent: 0.011852, IsBullMarket: true},
		Kline{Open: 18541.290000, Close: 18254.630000, Low: 17911.120000, High: 18557.320000, Volume: 52890.675094, Time: time.UnixMicro(1607558400000000), ChangePercent: -0.015461, IsBullMarket: false},
		Kline{Open: 18254.810000, Close: 18036.530000, Low: 17572.330000, High: 18292.730000, Volume: 72610.724259, Time: time.UnixMicro(1607644800000000), ChangePercent: -0.011957, IsBullMarket: false},
		Kline{Open: 18036.530000, Close: 18808.690000, Low: 18020.700000, High: 18948.660000, Volume: 49519.978432, Time: time.UnixMicro(1607731200000000), ChangePercent: 0.042811, IsBullMarket: true},
		Kline{Open: 18808.690000, Close: 19174.990000, Low: 18711.120000, High: 19411.000000, Volume: 56560.821744, Time: time.UnixMicro(1607817600000000), ChangePercent: 0.019475, IsBullMarket: true},
		Kline{Open: 19174.990000, Close: 19273.140000, Low: 19000.000000, High: 19349.000000, Volume: 47257.201294, Time: time.UnixMicro(1607904000000000), ChangePercent: 0.005119, IsBullMarket: true},
		Kline{Open: 19273.690000, Close: 19426.430000, Low: 19050.000000, High: 19570.000000, Volume: 61834.366011, Time: time.UnixMicro(1607990400000000), ChangePercent: 0.007925, IsBullMarket: true},
		Kline{Open: 19426.430000, Close: 21335.520000, Low: 19278.600000, High: 21560.000000, Volume: 114306.335570, Time: time.UnixMicro(1608076800000000), ChangePercent: 0.098273, IsBullMarket: true},
		Kline{Open: 21335.520000, Close: 22797.160000, Low: 21230.000000, High: 23800.000000, Volume: 184882.476748, Time: time.UnixMicro(1608163200000000), ChangePercent: 0.068507, IsBullMarket: true},
		Kline{Open: 22797.150000, Close: 23107.390000, Low: 22350.000000, High: 23285.180000, Volume: 79646.134315, Time: time.UnixMicro(1608249600000000), ChangePercent: 0.013609, IsBullMarket: true},
		Kline{Open: 23107.390000, Close: 23821.610000, Low: 22750.000000, High: 24171.470000, Volume: 86045.064677, Time: time.UnixMicro(1608336000000000), ChangePercent: 0.030909, IsBullMarket: true},
		Kline{Open: 23821.600000, Close: 23455.520000, Low: 23060.000000, High: 24295.000000, Volume: 76690.145685, Time: time.UnixMicro(1608422400000000), ChangePercent: -0.015368, IsBullMarket: false},
		Kline{Open: 23455.540000, Close: 22719.710000, Low: 21815.000000, High: 24102.770000, Volume: 88030.297243, Time: time.UnixMicro(1608508800000000), ChangePercent: -0.031371, IsBullMarket: false},
		Kline{Open: 22719.880000, Close: 23810.790000, Low: 22353.400000, High: 23837.100000, Volume: 87033.126160, Time: time.UnixMicro(1608595200000000), ChangePercent: 0.048016, IsBullMarket: true},
		Kline{Open: 23810.790000, Close: 23232.760000, Low: 22600.000000, High: 24100.000000, Volume: 119047.259733, Time: time.UnixMicro(1608681600000000), ChangePercent: -0.024276, IsBullMarket: false},
		Kline{Open: 23232.390000, Close: 23729.200000, Low: 22703.420000, High: 23794.430000, Volume: 69013.834252, Time: time.UnixMicro(1608768000000000), ChangePercent: 0.021384, IsBullMarket: true},
		Kline{Open: 23728.990000, Close: 24712.470000, Low: 23433.600000, High: 24789.860000, Volume: 79519.943569, Time: time.UnixMicro(1608854400000000), ChangePercent: 0.041446, IsBullMarket: true},
		Kline{Open: 24712.470000, Close: 26493.390000, Low: 24500.000000, High: 26867.030000, Volume: 97806.513386, Time: time.UnixMicro(1608940800000000), ChangePercent: 0.072066, IsBullMarket: true},
		Kline{Open: 26493.400000, Close: 26281.660000, Low: 25700.000000, High: 28422.000000, Volume: 148455.586214, Time: time.UnixMicro(1609027200000000), ChangePercent: -0.007992, IsBullMarket: false},
		Kline{Open: 26281.540000, Close: 27079.410000, Low: 26101.000000, High: 27500.000000, Volume: 79721.742496, Time: time.UnixMicro(1609113600000000), ChangePercent: 0.030359, IsBullMarket: true},
		Kline{Open: 27079.420000, Close: 27385.000000, Low: 25880.000000, High: 27410.000000, Volume: 69411.592606, Time: time.UnixMicro(1609200000000000), ChangePercent: 0.011285, IsBullMarket: true},
		Kline{Open: 27385.000000, Close: 28875.540000, Low: 27320.000000, High: 28996.000000, Volume: 95356.057826, Time: time.UnixMicro(1609286400000000), ChangePercent: 0.054429, IsBullMarket: true},
		Kline{Open: 28875.550000, Close: 28923.630000, Low: 27850.000000, High: 29300.000000, Volume: 75508.505152, Time: time.UnixMicro(1609372800000000), ChangePercent: 0.001665, IsBullMarket: true},
		Kline{Open: 28923.630000, Close: 29331.690000, Low: 28624.570000, High: 29600.000000, Volume: 54182.925011, Time: time.UnixMicro(1609459200000000), ChangePercent: 0.014108, IsBullMarket: true},
		Kline{Open: 29331.700000, Close: 32178.330000, Low: 28946.530000, High: 33300.000000, Volume: 129993.873362, Time: time.UnixMicro(1609545600000000), ChangePercent: 0.097050, IsBullMarket: true},
		Kline{Open: 32176.450000, Close: 33000.050000, Low: 31962.990000, High: 34778.110000, Volume: 120957.566750, Time: time.UnixMicro(1609632000000000), ChangePercent: 0.025596, IsBullMarket: true},
		Kline{Open: 33000.050000, Close: 31988.710000, Low: 28130.000000, High: 33600.000000, Volume: 140899.885690, Time: time.UnixMicro(1609718400000000), ChangePercent: -0.030647, IsBullMarket: false},
		Kline{Open: 31989.750000, Close: 33949.530000, Low: 29900.000000, High: 34360.000000, Volume: 116049.997038, Time: time.UnixMicro(1609804800000000), ChangePercent: 0.061263, IsBullMarket: true},
		Kline{Open: 33949.530000, Close: 36769.360000, Low: 33288.000000, High: 36939.210000, Volume: 127139.201310, Time: time.UnixMicro(1609891200000000), ChangePercent: 0.083059, IsBullMarket: true},
		Kline{Open: 36769.360000, Close: 39432.280000, Low: 36300.000000, High: 40365.000000, Volume: 132825.700437, Time: time.UnixMicro(1609977600000000), ChangePercent: 0.072422, IsBullMarket: true},
		Kline{Open: 39432.480000, Close: 40582.810000, Low: 36500.000000, High: 41950.000000, Volume: 139789.957499, Time: time.UnixMicro(1610064000000000), ChangePercent: 0.029172, IsBullMarket: true},
		Kline{Open: 40586.960000, Close: 40088.220000, Low: 38720.000000, High: 41380.000000, Volume: 75785.979675, Time: time.UnixMicro(1610150400000000), ChangePercent: -0.012288, IsBullMarket: false},
		Kline{Open: 40088.220000, Close: 38150.020000, Low: 35111.110000, High: 41350.000000, Volume: 118209.544503, Time: time.UnixMicro(1610236800000000), ChangePercent: -0.048348, IsBullMarket: false},
		Kline{Open: 38150.020000, Close: 35404.470000, Low: 30420.000000, High: 38264.740000, Volume: 251156.138287, Time: time.UnixMicro(1610323200000000), ChangePercent: -0.071967, IsBullMarket: false},
		Kline{Open: 35410.370000, Close: 34051.240000, Low: 32531.000000, High: 36628.000000, Volume: 133948.151996, Time: time.UnixMicro(1610409600000000), ChangePercent: -0.038382, IsBullMarket: false},
		Kline{Open: 34049.150000, Close: 37371.380000, Low: 32380.000000, High: 37850.000000, Volume: 124477.914938, Time: time.UnixMicro(1610496000000000), ChangePercent: 0.097572, IsBullMarket: true},
		Kline{Open: 37371.380000, Close: 39144.500000, Low: 36701.230000, High: 40100.000000, Volume: 102950.389421, Time: time.UnixMicro(1610582400000000), ChangePercent: 0.047446, IsBullMarket: true},
		Kline{Open: 39145.210000, Close: 36742.220000, Low: 34408.000000, High: 39747.760000, Volume: 118300.920916, Time: time.UnixMicro(1610668800000000), ChangePercent: -0.061387, IsBullMarket: false},
		Kline{Open: 36737.430000, Close: 35994.980000, Low: 35357.800000, High: 37950.000000, Volume: 86348.431508, Time: time.UnixMicro(1610755200000000), ChangePercent: -0.020210, IsBullMarket: false},
		Kline{Open: 35994.980000, Close: 35828.610000, Low: 33850.000000, High: 36852.500000, Volume: 80157.727384, Time: time.UnixMicro(1610841600000000), ChangePercent: -0.004622, IsBullMarket: false},
		Kline{Open: 35824.990000, Close: 36631.270000, Low: 34800.000000, High: 37469.830000, Volume: 70698.118750, Time: time.UnixMicro(1610928000000000), ChangePercent: 0.022506, IsBullMarket: true},
		Kline{Open: 36622.460000, Close: 35891.490000, Low: 35844.060000, High: 37850.000000, Volume: 79611.307769, Time: time.UnixMicro(1611014400000000), ChangePercent: -0.019960, IsBullMarket: false},
		Kline{Open: 35901.940000, Close: 35468.230000, Low: 33400.000000, High: 36415.310000, Volume: 89368.422918, Time: time.UnixMicro(1611100800000000), ChangePercent: -0.012080, IsBullMarket: false},
		Kline{Open: 35468.230000, Close: 30850.130000, Low: 30071.000000, High: 35600.000000, Volume: 135004.076658, Time: time.UnixMicro(1611187200000000), ChangePercent: -0.130204, IsBullMarket: false},
		Kline{Open: 30851.990000, Close: 32945.170000, Low: 28850.000000, High: 33826.530000, Volume: 142971.684049, Time: time.UnixMicro(1611273600000000), ChangePercent: 0.067846, IsBullMarket: true},
		Kline{Open: 32950.000000, Close: 32078.000000, Low: 31390.160000, High: 33456.000000, Volume: 64595.287675, Time: time.UnixMicro(1611360000000000), ChangePercent: -0.026464, IsBullMarket: false},
		Kline{Open: 32078.000000, Close: 32259.900000, Low: 30900.000000, High: 33071.000000, Volume: 57978.037966, Time: time.UnixMicro(1611446400000000), ChangePercent: 0.005671, IsBullMarket: true},
		Kline{Open: 32259.450000, Close: 32254.200000, Low: 31910.000000, High: 34875.000000, Volume: 88499.226921, Time: time.UnixMicro(1611532800000000), ChangePercent: -0.000163, IsBullMarket: false},
		Kline{Open: 32254.190000, Close: 32467.770000, Low: 30837.370000, High: 32921.880000, Volume: 84972.206910, Time: time.UnixMicro(1611619200000000), ChangePercent: 0.006622, IsBullMarket: true},
		Kline{Open: 32464.010000, Close: 30366.150000, Low: 29241.720000, High: 32557.290000, Volume: 95911.961711, Time: time.UnixMicro(1611705600000000), ChangePercent: -0.064621, IsBullMarket: false},
		Kline{Open: 30362.190000, Close: 33364.860000, Low: 29842.100000, High: 33783.980000, Volume: 92621.145617, Time: time.UnixMicro(1611792000000000), ChangePercent: 0.098895, IsBullMarket: true},
		Kline{Open: 33368.180000, Close: 34252.200000, Low: 31915.400000, High: 38531.900000, Volume: 231827.005626, Time: time.UnixMicro(1611878400000000), ChangePercent: 0.026493, IsBullMarket: true},
		Kline{Open: 34246.280000, Close: 34262.880000, Low: 32825.000000, High: 34933.000000, Volume: 84889.681340, Time: time.UnixMicro(1611964800000000), ChangePercent: 0.000485, IsBullMarket: true},
		Kline{Open: 34262.890000, Close: 33092.980000, Low: 32171.670000, High: 34342.690000, Volume: 68742.280384, Time: time.UnixMicro(1612051200000000), ChangePercent: -0.034145, IsBullMarket: false},
		Kline{Open: 33092.970000, Close: 33526.370000, Low: 32296.160000, High: 34717.270000, Volume: 82718.276882, Time: time.UnixMicro(1612137600000000), ChangePercent: 0.013096, IsBullMarket: true},
		Kline{Open: 33517.090000, Close: 35466.240000, Low: 33418.000000, High: 35984.330000, Volume: 78056.659880, Time: time.UnixMicro(1612224000000000), ChangePercent: 0.058154, IsBullMarket: true},
		Kline{Open: 35472.710000, Close: 37618.870000, Low: 35362.380000, High: 37662.630000, Volume: 80784.333663, Time: time.UnixMicro(1612310400000000), ChangePercent: 0.060502, IsBullMarket: true},
		Kline{Open: 37620.260000, Close: 36936.660000, Low: 36161.950000, High: 38708.270000, Volume: 92080.735898, Time: time.UnixMicro(1612396800000000), ChangePercent: -0.018171, IsBullMarket: false},
		Kline{Open: 36936.650000, Close: 38290.240000, Low: 36570.000000, High: 38310.120000, Volume: 66681.334275, Time: time.UnixMicro(1612483200000000), ChangePercent: 0.036646, IsBullMarket: true},
		Kline{Open: 38289.320000, Close: 39186.940000, Low: 38215.940000, High: 40955.510000, Volume: 98757.311183, Time: time.UnixMicro(1612569600000000), ChangePercent: 0.023443, IsBullMarket: true},
		Kline{Open: 39181.010000, Close: 38795.690000, Low: 37351.000000, High: 39700.000000, Volume: 84363.679763, Time: time.UnixMicro(1612656000000000), ChangePercent: -0.009834, IsBullMarket: false},
		Kline{Open: 38795.690000, Close: 46374.870000, Low: 37988.890000, High: 46794.450000, Volume: 138597.536914, Time: time.UnixMicro(1612742400000000), ChangePercent: 0.195361, IsBullMarket: true},
		Kline{Open: 46374.860000, Close: 46420.420000, Low: 44961.090000, High: 48142.190000, Volume: 115499.861712, Time: time.UnixMicro(1612828800000000), ChangePercent: 0.000982, IsBullMarket: true},
		Kline{Open: 46420.420000, Close: 44807.580000, Low: 43727.000000, High: 47310.000000, Volume: 97154.182200, Time: time.UnixMicro(1612915200000000), ChangePercent: -0.034744, IsBullMarket: false},
		Kline{Open: 44807.580000, Close: 47969.510000, Low: 43994.020000, High: 48678.900000, Volume: 89561.081454, Time: time.UnixMicro(1613001600000000), ChangePercent: 0.070567, IsBullMarket: true},
		Kline{Open: 47968.660000, Close: 47287.600000, Low: 46125.000000, High: 48985.800000, Volume: 85870.035697, Time: time.UnixMicro(1613088000000000), ChangePercent: -0.014198, IsBullMarket: false},
		Kline{Open: 47298.150000, Close: 47153.690000, Low: 46202.530000, High: 48150.000000, Volume: 63768.097399, Time: time.UnixMicro(1613174400000000), ChangePercent: -0.003054, IsBullMarket: false},
		Kline{Open: 47156.780000, Close: 48577.790000, Low: 47014.170000, High: 49707.430000, Volume: 73735.475533, Time: time.UnixMicro(1613260800000000), ChangePercent: 0.030134, IsBullMarket: true},
		Kline{Open: 48580.470000, Close: 47911.100000, Low: 45570.790000, High: 49010.920000, Volume: 79398.156784, Time: time.UnixMicro(1613347200000000), ChangePercent: -0.013779, IsBullMarket: false},
		Kline{Open: 47911.100000, Close: 49133.450000, Low: 47003.620000, High: 50689.180000, Volume: 88813.266298, Time: time.UnixMicro(1613433600000000), ChangePercent: 0.025513, IsBullMarket: true},
		Kline{Open: 49133.450000, Close: 52119.710000, Low: 48947.000000, High: 52618.740000, Volume: 85743.637818, Time: time.UnixMicro(1613520000000000), ChangePercent: 0.060779, IsBullMarket: true},
		Kline{Open: 52117.670000, Close: 51552.600000, Low: 50901.900000, High: 52530.000000, Volume: 60758.046954, Time: time.UnixMicro(1613606400000000), ChangePercent: -0.010842, IsBullMarket: false},
		Kline{Open: 51552.610000, Close: 55906.000000, Low: 50710.200000, High: 56368.000000, Volume: 79659.778020, Time: time.UnixMicro(1613692800000000), ChangePercent: 0.084446, IsBullMarket: true},
		Kline{Open: 55906.000000, Close: 55841.190000, Low: 53863.930000, High: 57700.460000, Volume: 80948.205314, Time: time.UnixMicro(1613779200000000), ChangePercent: -0.001159, IsBullMarket: false},
		Kline{Open: 55841.190000, Close: 57408.570000, Low: 55477.590000, High: 58352.800000, Volume: 58166.708511, Time: time.UnixMicro(1613865600000000), ChangePercent: 0.028069, IsBullMarket: true},
		Kline{Open: 57412.350000, Close: 54087.670000, Low: 47622.000000, High: 57508.470000, Volume: 134019.434944, Time: time.UnixMicro(1613952000000000), ChangePercent: -0.057909, IsBullMarket: false},
		Kline{Open: 54087.670000, Close: 48891.000000, Low: 44892.560000, High: 54183.590000, Volume: 169375.025051, Time: time.UnixMicro(1614038400000000), ChangePercent: -0.096079, IsBullMarket: false},
		Kline{Open: 48891.000000, Close: 49676.200000, Low: 46988.690000, High: 51374.990000, Volume: 91881.209252, Time: time.UnixMicro(1614124800000000), ChangePercent: 0.016060, IsBullMarket: true},
		Kline{Open: 49676.210000, Close: 47073.730000, Low: 46674.340000, High: 52041.730000, Volume: 83310.673121, Time: time.UnixMicro(1614211200000000), ChangePercent: -0.052389, IsBullMarket: false},
		Kline{Open: 47073.730000, Close: 46276.870000, Low: 44106.780000, High: 48424.110000, Volume: 109423.200663, Time: time.UnixMicro(1614297600000000), ChangePercent: -0.016928, IsBullMarket: false},
		Kline{Open: 46276.880000, Close: 46106.430000, Low: 45000.000000, High: 48394.000000, Volume: 66060.834292, Time: time.UnixMicro(1614384000000000), ChangePercent: -0.003683, IsBullMarket: false},
		Kline{Open: 46103.670000, Close: 45135.660000, Low: 43000.000000, High: 46638.460000, Volume: 83055.369042, Time: time.UnixMicro(1614470400000000), ChangePercent: -0.020996, IsBullMarket: false},
		Kline{Open: 45134.110000, Close: 49587.030000, Low: 44950.530000, High: 49790.000000, Volume: 85086.111648, Time: time.UnixMicro(1614556800000000), ChangePercent: 0.098660, IsBullMarket: true},
		Kline{Open: 49595.760000, Close: 48440.650000, Low: 47047.600000, High: 50200.000000, Volume: 64221.062140, Time: time.UnixMicro(1614643200000000), ChangePercent: -0.023290, IsBullMarket: false},
		Kline{Open: 48436.610000, Close: 50349.370000, Low: 48100.710000, High: 52640.000000, Volume: 81035.913705, Time: time.UnixMicro(1614729600000000), ChangePercent: 0.039490, IsBullMarket: true},
		Kline{Open: 50349.370000, Close: 48374.090000, Low: 47500.000000, High: 51773.880000, Volume: 82649.716829, Time: time.UnixMicro(1614816000000000), ChangePercent: -0.039231, IsBullMarket: false},
		Kline{Open: 48374.090000, Close: 48751.710000, Low: 46300.000000, High: 49448.930000, Volume: 78192.496372, Time: time.UnixMicro(1614902400000000), ChangePercent: 0.007806, IsBullMarket: true},
		Kline{Open: 48746.810000, Close: 48882.200000, Low: 47070.000000, High: 49200.000000, Volume: 44399.234242, Time: time.UnixMicro(1614988800000000), ChangePercent: 0.002777, IsBullMarket: true},
		Kline{Open: 48882.200000, Close: 50971.750000, Low: 48882.200000, High: 51450.030000, Volume: 55235.028032, Time: time.UnixMicro(1615075200000000), ChangePercent: 0.042747, IsBullMarket: true},
		Kline{Open: 50959.110000, Close: 52375.170000, Low: 49274.670000, High: 52402.780000, Volume: 66987.359664, Time: time.UnixMicro(1615161600000000), ChangePercent: 0.027788, IsBullMarket: true},
		Kline{Open: 52375.180000, Close: 54884.500000, Low: 51789.410000, High: 54895.000000, Volume: 71656.737076, Time: time.UnixMicro(1615248000000000), ChangePercent: 0.047910, IsBullMarket: true},
		Kline{Open: 54874.670000, Close: 55851.590000, Low: 53005.000000, High: 57387.690000, Volume: 84749.238943, Time: time.UnixMicro(1615334400000000), ChangePercent: 0.017803, IsBullMarket: true},
		Kline{Open: 55851.590000, Close: 57773.160000, Low: 54272.820000, High: 58150.000000, Volume: 81914.812859, Time: time.UnixMicro(1615420800000000), ChangePercent: 0.034405, IsBullMarket: true},
		Kline{Open: 57773.150000, Close: 57221.720000, Low: 54962.840000, High: 58081.510000, Volume: 73405.406047, Time: time.UnixMicro(1615507200000000), ChangePercent: -0.009545, IsBullMarket: false},
		Kline{Open: 57221.720000, Close: 61188.390000, Low: 56078.230000, High: 61844.000000, Volume: 83245.091346, Time: time.UnixMicro(1615593600000000), ChangePercent: 0.069321, IsBullMarket: true},
		Kline{Open: 61188.380000, Close: 58968.310000, Low: 58966.780000, High: 61724.790000, Volume: 52601.052750, Time: time.UnixMicro(1615680000000000), ChangePercent: -0.036283, IsBullMarket: false},
		Kline{Open: 58976.080000, Close: 55605.200000, Low: 54600.000000, High: 60633.430000, Volume: 102771.427298, Time: time.UnixMicro(1615766400000000), ChangePercent: -0.057157, IsBullMarket: false},
		Kline{Open: 55605.200000, Close: 56900.750000, Low: 53271.340000, High: 56938.290000, Volume: 77986.694355, Time: time.UnixMicro(1615852800000000), ChangePercent: 0.023299, IsBullMarket: true},
		Kline{Open: 56900.740000, Close: 58912.970000, Low: 54123.690000, High: 58974.730000, Volume: 70421.620841, Time: time.UnixMicro(1615939200000000), ChangePercent: 0.035364, IsBullMarket: true},
		Kline{Open: 58912.970000, Close: 57648.160000, Low: 57023.000000, High: 60129.970000, Volume: 66580.406675, Time: time.UnixMicro(1616025600000000), ChangePercent: -0.021469, IsBullMarket: false},
		Kline{Open: 57641.000000, Close: 58030.010000, Low: 56270.740000, High: 59468.000000, Volume: 52392.652961, Time: time.UnixMicro(1616112000000000), ChangePercent: 0.006749, IsBullMarket: true},
		Kline{Open: 58030.010000, Close: 58102.280000, Low: 57820.170000, High: 59880.000000, Volume: 44476.941776, Time: time.UnixMicro(1616198400000000), ChangePercent: 0.001245, IsBullMarket: true},
		Kline{Open: 58100.020000, Close: 57351.560000, Low: 55450.110000, High: 58589.100000, Volume: 48564.470274, Time: time.UnixMicro(1616284800000000), ChangePercent: -0.012882, IsBullMarket: false},
		Kline{Open: 57351.560000, Close: 54083.250000, Low: 53650.000000, High: 58430.730000, Volume: 62581.626169, Time: time.UnixMicro(1616371200000000), ChangePercent: -0.056987, IsBullMarket: false},
		Kline{Open: 54083.250000, Close: 54340.890000, Low: 53000.000000, High: 55830.900000, Volume: 59789.365427, Time: time.UnixMicro(1616457600000000), ChangePercent: 0.004764, IsBullMarket: true},
		Kline{Open: 54342.800000, Close: 52303.650000, Low: 51700.000000, High: 57200.000000, Volume: 83537.465021, Time: time.UnixMicro(1616544000000000), ChangePercent: -0.037524, IsBullMarket: false},
		Kline{Open: 52303.660000, Close: 51293.780000, Low: 50427.560000, High: 53287.000000, Volume: 87400.534538, Time: time.UnixMicro(1616630400000000), ChangePercent: -0.019308, IsBullMarket: false},
		Kline{Open: 51293.780000, Close: 55025.590000, Low: 51214.600000, High: 55073.460000, Volume: 63813.774692, Time: time.UnixMicro(1616716800000000), ChangePercent: 0.072754, IsBullMarket: true},
		Kline{Open: 55025.590000, Close: 55817.140000, Low: 53950.000000, High: 56700.360000, Volume: 50105.475055, Time: time.UnixMicro(1616803200000000), ChangePercent: 0.014385, IsBullMarket: true},
		Kline{Open: 55817.140000, Close: 55777.630000, Low: 54691.840000, High: 56559.750000, Volume: 39050.387511, Time: time.UnixMicro(1616889600000000), ChangePercent: -0.000708, IsBullMarket: false},
		Kline{Open: 55777.650000, Close: 57635.470000, Low: 54800.010000, High: 58405.820000, Volume: 67857.937398, Time: time.UnixMicro(1616976000000000), ChangePercent: 0.033308, IsBullMarket: true},
		Kline{Open: 57635.460000, Close: 58746.570000, Low: 57071.350000, High: 59368.000000, Volume: 55122.443122, Time: time.UnixMicro(1617062400000000), ChangePercent: 0.019278, IsBullMarket: true},
		Kline{Open: 58746.570000, Close: 58740.550000, Low: 56769.000000, High: 59800.000000, Volume: 60975.542666, Time: time.UnixMicro(1617148800000000), ChangePercent: -0.000102, IsBullMarket: false},
		Kline{Open: 58739.460000, Close: 58720.440000, Low: 57935.450000, High: 59490.000000, Volume: 47415.617220, Time: time.UnixMicro(1617235200000000), ChangePercent: -0.000324, IsBullMarket: false},
		Kline{Open: 58720.450000, Close: 58950.010000, Low: 58428.570000, High: 60200.000000, Volume: 47382.418781, Time: time.UnixMicro(1617321600000000), ChangePercent: 0.003909, IsBullMarket: true},
		Kline{Open: 58950.010000, Close: 57051.940000, Low: 56880.000000, High: 59791.720000, Volume: 47409.852113, Time: time.UnixMicro(1617408000000000), ChangePercent: -0.032198, IsBullMarket: false},
		Kline{Open: 57051.950000, Close: 58202.010000, Low: 56388.000000, High: 58492.850000, Volume: 41314.081973, Time: time.UnixMicro(1617494400000000), ChangePercent: 0.020158, IsBullMarket: true},
		Kline{Open: 58202.010000, Close: 59129.990000, Low: 56777.770000, High: 59272.000000, Volume: 54258.015790, Time: time.UnixMicro(1617580800000000), ChangePercent: 0.015944, IsBullMarket: true},
		Kline{Open: 59129.990000, Close: 57991.150000, Low: 57413.020000, High: 59495.240000, Volume: 54201.000727, Time: time.UnixMicro(1617667200000000), ChangePercent: -0.019260, IsBullMarket: false},
		Kline{Open: 57990.030000, Close: 55953.450000, Low: 55473.000000, High: 58655.000000, Volume: 71228.405659, Time: time.UnixMicro(1617753600000000), ChangePercent: -0.035119, IsBullMarket: false},
		Kline{Open: 55953.440000, Close: 58077.520000, Low: 55700.000000, High: 58153.310000, Volume: 44283.147019, Time: time.UnixMicro(1617840000000000), ChangePercent: 0.037962, IsBullMarket: true},
		Kline{Open: 58077.520000, Close: 58142.540000, Low: 57654.000000, High: 58894.900000, Volume: 40831.884911, Time: time.UnixMicro(1617926400000000), ChangePercent: 0.001120, IsBullMarket: true},
		Kline{Open: 58142.550000, Close: 59769.130000, Low: 57900.010000, High: 61500.000000, Volume: 69906.424117, Time: time.UnixMicro(1618012800000000), ChangePercent: 0.027976, IsBullMarket: true},
		Kline{Open: 59769.130000, Close: 60002.430000, Low: 59232.520000, High: 60699.000000, Volume: 41156.715391, Time: time.UnixMicro(1618099200000000), ChangePercent: 0.003903, IsBullMarket: true},
		Kline{Open: 59998.800000, Close: 59860.000000, Low: 59350.590000, High: 61300.000000, Volume: 56375.037117, Time: time.UnixMicro(1618185600000000), ChangePercent: -0.002313, IsBullMarket: false},
		Kline{Open: 59860.010000, Close: 63575.000000, Low: 59805.150000, High: 63777.770000, Volume: 82848.688746, Time: time.UnixMicro(1618272000000000), ChangePercent: 0.062061, IsBullMarket: true},
		Kline{Open: 63575.010000, Close: 62959.530000, Low: 61301.000000, High: 64854.000000, Volume: 82616.343993, Time: time.UnixMicro(1618358400000000), ChangePercent: -0.009681, IsBullMarket: false},
		Kline{Open: 62959.530000, Close: 63159.980000, Low: 62020.000000, High: 63800.000000, Volume: 51649.700340, Time: time.UnixMicro(1618444800000000), ChangePercent: 0.003184, IsBullMarket: true},
		Kline{Open: 63158.740000, Close: 61334.800000, Low: 60000.000000, High: 63520.610000, Volume: 91764.139884, Time: time.UnixMicro(1618531200000000), ChangePercent: -0.028879, IsBullMarket: false},
		Kline{Open: 61334.810000, Close: 60006.660000, Low: 59580.910000, High: 62506.050000, Volume: 58912.256128, Time: time.UnixMicro(1618617600000000), ChangePercent: -0.021654, IsBullMarket: false},
		Kline{Open: 60006.670000, Close: 56150.010000, Low: 50931.300000, High: 60499.000000, Volume: 124882.131824, Time: time.UnixMicro(1618704000000000), ChangePercent: -0.064271, IsBullMarket: false},
		Kline{Open: 56150.010000, Close: 55633.140000, Low: 54221.580000, High: 57526.810000, Volume: 78229.042267, Time: time.UnixMicro(1618790400000000), ChangePercent: -0.009205, IsBullMarket: false},
		Kline{Open: 55633.140000, Close: 56425.000000, Low: 53329.960000, High: 57076.240000, Volume: 72744.482151, Time: time.UnixMicro(1618876800000000), ChangePercent: 0.014234, IsBullMarket: true},
		Kline{Open: 56425.000000, Close: 53787.630000, Low: 53536.020000, High: 56757.910000, Volume: 66984.756909, Time: time.UnixMicro(1618963200000000), ChangePercent: -0.046741, IsBullMarket: false},
		Kline{Open: 53787.620000, Close: 51690.960000, Low: 50500.000000, High: 55521.480000, Volume: 104656.631337, Time: time.UnixMicro(1619049600000000), ChangePercent: -0.038980, IsBullMarket: false},
		Kline{Open: 51690.950000, Close: 51125.140000, Low: 47500.000000, High: 52131.850000, Volume: 132230.780719, Time: time.UnixMicro(1619136000000000), ChangePercent: -0.010946, IsBullMarket: false},
		Kline{Open: 51110.560000, Close: 50047.840000, Low: 48657.140000, High: 51166.220000, Volume: 55361.512573, Time: time.UnixMicro(1619222400000000), ChangePercent: -0.020793, IsBullMarket: false},
		Kline{Open: 50047.840000, Close: 49066.770000, Low: 46930.000000, High: 50567.910000, Volume: 58255.645004, Time: time.UnixMicro(1619308800000000), ChangePercent: -0.019603, IsBullMarket: false},
		Kline{Open: 49066.760000, Close: 54001.390000, Low: 48753.440000, High: 54356.620000, Volume: 86310.802124, Time: time.UnixMicro(1619395200000000), ChangePercent: 0.100570, IsBullMarket: true},
		Kline{Open: 54001.380000, Close: 55011.970000, Low: 53222.000000, High: 55460.000000, Volume: 54064.034675, Time: time.UnixMicro(1619481600000000), ChangePercent: 0.018714, IsBullMarket: true},
		Kline{Open: 55011.970000, Close: 54846.220000, Low: 53813.160000, High: 56428.000000, Volume: 55130.459015, Time: time.UnixMicro(1619568000000000), ChangePercent: -0.003013, IsBullMarket: false},
		Kline{Open: 54846.230000, Close: 53555.000000, Low: 52330.940000, High: 55195.840000, Volume: 52486.019455, Time: time.UnixMicro(1619654400000000), ChangePercent: -0.023543, IsBullMarket: false},
		Kline{Open: 53555.000000, Close: 57694.270000, Low: 53013.010000, High: 57963.000000, Volume: 68578.910045, Time: time.UnixMicro(1619740800000000), ChangePercent: 0.077290, IsBullMarket: true},
		Kline{Open: 57697.250000, Close: 57800.370000, Low: 56956.140000, High: 58458.070000, Volume: 42600.351836, Time: time.UnixMicro(1619827200000000), ChangePercent: 0.001787, IsBullMarket: true},
		Kline{Open: 57797.350000, Close: 56578.210000, Low: 56035.250000, High: 57911.020000, Volume: 36812.878863, Time: time.UnixMicro(1619913600000000), ChangePercent: -0.021093, IsBullMarket: false},
		Kline{Open: 56578.210000, Close: 57169.390000, Low: 56435.000000, High: 58981.440000, Volume: 57649.931286, Time: time.UnixMicro(1620000000000000), ChangePercent: 0.010449, IsBullMarket: true},
		Kline{Open: 57169.390000, Close: 53200.010000, Low: 53046.690000, High: 57200.000000, Volume: 85324.625903, Time: time.UnixMicro(1620086400000000), ChangePercent: -0.069432, IsBullMarket: false},
		Kline{Open: 53205.050000, Close: 57436.110000, Low: 52900.000000, High: 58069.820000, Volume: 77263.923439, Time: time.UnixMicro(1620172800000000), ChangePercent: 0.079524, IsBullMarket: true},
		Kline{Open: 57436.110000, Close: 56393.680000, Low: 55200.000000, High: 58360.000000, Volume: 70181.671908, Time: time.UnixMicro(1620259200000000), ChangePercent: -0.018149, IsBullMarket: false},
		Kline{Open: 56393.680000, Close: 57314.750000, Low: 55241.630000, High: 58650.000000, Volume: 74542.747829, Time: time.UnixMicro(1620345600000000), ChangePercent: 0.016333, IsBullMarket: true},
		Kline{Open: 57315.490000, Close: 58862.050000, Low: 56900.000000, High: 59500.000000, Volume: 69709.906028, Time: time.UnixMicro(1620432000000000), ChangePercent: 0.026983, IsBullMarket: true},
		Kline{Open: 58866.530000, Close: 58240.840000, Low: 56235.660000, High: 59300.000000, Volume: 69806.119910, Time: time.UnixMicro(1620518400000000), ChangePercent: -0.010629, IsBullMarket: false},
		Kline{Open: 58240.830000, Close: 55816.140000, Low: 53400.000000, High: 59500.000000, Volume: 89586.349250, Time: time.UnixMicro(1620604800000000), ChangePercent: -0.041632, IsBullMarket: false},
		Kline{Open: 55816.140000, Close: 56670.020000, Low: 54370.000000, High: 56862.430000, Volume: 64329.540550, Time: time.UnixMicro(1620691200000000), ChangePercent: 0.015298, IsBullMarket: true},
		Kline{Open: 56670.020000, Close: 49631.320000, Low: 48600.000000, High: 58000.010000, Volume: 99842.789836, Time: time.UnixMicro(1620777600000000), ChangePercent: -0.124205, IsBullMarket: false},
		Kline{Open: 49537.150000, Close: 49670.970000, Low: 46000.000000, High: 51367.190000, Volume: 147332.002121, Time: time.UnixMicro(1620864000000000), ChangePercent: 0.002701, IsBullMarket: true},
		Kline{Open: 49671.920000, Close: 49841.450000, Low: 48799.750000, High: 51483.000000, Volume: 80082.204306, Time: time.UnixMicro(1620950400000000), ChangePercent: 0.003413, IsBullMarket: true},
		Kline{Open: 49844.160000, Close: 46762.990000, Low: 46555.000000, High: 50700.000000, Volume: 89437.449359, Time: time.UnixMicro(1621036800000000), ChangePercent: -0.061816, IsBullMarket: false},
		Kline{Open: 46762.990000, Close: 46431.500000, Low: 43825.390000, High: 49795.890000, Volume: 114269.812775, Time: time.UnixMicro(1621123200000000), ChangePercent: -0.007089, IsBullMarket: false},
		Kline{Open: 46426.830000, Close: 43538.040000, Low: 42001.000000, High: 46686.000000, Volume: 166657.172736, Time: time.UnixMicro(1621209600000000), ChangePercent: -0.062222, IsBullMarket: false},
		Kline{Open: 43538.020000, Close: 42849.780000, Low: 42250.020000, High: 45799.290000, Volume: 116979.860784, Time: time.UnixMicro(1621296000000000), ChangePercent: -0.015808, IsBullMarket: false},
		Kline{Open: 42849.780000, Close: 36690.090000, Low: 30000.000000, High: 43584.900000, Volume: 354347.243161, Time: time.UnixMicro(1621382400000000), ChangePercent: -0.143751, IsBullMarket: false},
		Kline{Open: 36671.230000, Close: 40526.640000, Low: 34850.000000, High: 42451.670000, Volume: 203017.596923, Time: time.UnixMicro(1621468800000000), ChangePercent: 0.105134, IsBullMarket: true},
		Kline{Open: 40525.390000, Close: 37252.010000, Low: 33488.000000, High: 42200.000000, Volume: 202100.888258, Time: time.UnixMicro(1621555200000000), ChangePercent: -0.080774, IsBullMarket: false},
		Kline{Open: 37263.350000, Close: 37449.730000, Low: 35200.620000, High: 38829.000000, Volume: 126542.243689, Time: time.UnixMicro(1621641600000000), ChangePercent: 0.005002, IsBullMarket: true},
		Kline{Open: 37458.510000, Close: 34655.250000, Low: 31111.010000, High: 38270.640000, Volume: 217136.046593, Time: time.UnixMicro(1621728000000000), ChangePercent: -0.074836, IsBullMarket: false},
		Kline{Open: 34681.440000, Close: 38796.290000, Low: 34031.000000, High: 39920.000000, Volume: 161630.893971, Time: time.UnixMicro(1621814400000000), ChangePercent: 0.118647, IsBullMarket: true},
		Kline{Open: 38810.990000, Close: 38324.720000, Low: 36419.620000, High: 39791.770000, Volume: 111996.228404, Time: time.UnixMicro(1621900800000000), ChangePercent: -0.012529, IsBullMarket: false},
		Kline{Open: 38324.720000, Close: 39241.910000, Low: 37800.440000, High: 40841.000000, Volume: 104780.773396, Time: time.UnixMicro(1621987200000000), ChangePercent: 0.023932, IsBullMarket: true},
		Kline{Open: 39241.920000, Close: 38529.980000, Low: 37134.270000, High: 40411.140000, Volume: 86547.158794, Time: time.UnixMicro(1622073600000000), ChangePercent: -0.018142, IsBullMarket: false},
		Kline{Open: 38529.990000, Close: 35663.490000, Low: 34684.000000, High: 38877.830000, Volume: 135377.629720, Time: time.UnixMicro(1622160000000000), ChangePercent: -0.074397, IsBullMarket: false},
		Kline{Open: 35661.790000, Close: 34605.150000, Low: 33632.760000, High: 37338.580000, Volume: 112663.092689, Time: time.UnixMicro(1622246400000000), ChangePercent: -0.029629, IsBullMarket: false},
		Kline{Open: 34605.150000, Close: 35641.270000, Low: 33379.000000, High: 36488.000000, Volume: 73535.386967, Time: time.UnixMicro(1622332800000000), ChangePercent: 0.029941, IsBullMarket: true},
		Kline{Open: 35641.260000, Close: 37253.810000, Low: 34153.840000, High: 37499.000000, Volume: 94160.735289, Time: time.UnixMicro(1622419200000000), ChangePercent: 0.045244, IsBullMarket: true},
		Kline{Open: 37253.820000, Close: 36693.090000, Low: 35666.000000, High: 37894.810000, Volume: 81234.663770, Time: time.UnixMicro(1622505600000000), ChangePercent: -0.015052, IsBullMarket: false},
		Kline{Open: 36694.850000, Close: 37568.680000, Low: 35920.000000, High: 38225.000000, Volume: 67587.372495, Time: time.UnixMicro(1622592000000000), ChangePercent: 0.023813, IsBullMarket: true},
		Kline{Open: 37568.680000, Close: 39246.790000, Low: 37170.000000, High: 39476.000000, Volume: 75889.106011, Time: time.UnixMicro(1622678400000000), ChangePercent: 0.044668, IsBullMarket: true},
		Kline{Open: 39246.780000, Close: 36829.000000, Low: 35555.150000, High: 39289.070000, Volume: 91317.799245, Time: time.UnixMicro(1622764800000000), ChangePercent: -0.061605, IsBullMarket: false},
		Kline{Open: 36829.150000, Close: 35513.200000, Low: 34800.000000, High: 37925.000000, Volume: 70459.621490, Time: time.UnixMicro(1622851200000000), ChangePercent: -0.035731, IsBullMarket: false},
		Kline{Open: 35516.070000, Close: 35796.310000, Low: 35222.000000, High: 36480.000000, Volume: 47650.206637, Time: time.UnixMicro(1622937600000000), ChangePercent: 0.007891, IsBullMarket: true},
		Kline{Open: 35796.310000, Close: 33552.790000, Low: 33300.000000, High: 36900.000000, Volume: 77574.952573, Time: time.UnixMicro(1623024000000000), ChangePercent: -0.062675, IsBullMarket: false},
		Kline{Open: 33556.960000, Close: 33380.810000, Low: 31000.000000, High: 34068.010000, Volume: 123251.189037, Time: time.UnixMicro(1623110400000000), ChangePercent: -0.005249, IsBullMarket: false},
		Kline{Open: 33380.800000, Close: 37388.050000, Low: 32396.820000, High: 37534.790000, Volume: 136607.597517, Time: time.UnixMicro(1623196800000000), ChangePercent: 0.120047, IsBullMarket: true},
		Kline{Open: 37388.050000, Close: 36675.720000, Low: 35782.000000, High: 38491.000000, Volume: 109527.284943, Time: time.UnixMicro(1623283200000000), ChangePercent: -0.019052, IsBullMarket: false},
		Kline{Open: 36677.830000, Close: 37331.980000, Low: 35936.770000, High: 37680.400000, Volume: 78466.005300, Time: time.UnixMicro(1623369600000000), ChangePercent: 0.017835, IsBullMarket: true},
		Kline{Open: 37331.980000, Close: 35546.110000, Low: 34600.360000, High: 37463.630000, Volume: 87717.549990, Time: time.UnixMicro(1623456000000000), ChangePercent: -0.047838, IsBullMarket: false},
		Kline{Open: 35546.120000, Close: 39020.570000, Low: 34757.000000, High: 39380.000000, Volume: 86921.025555, Time: time.UnixMicro(1623542400000000), ChangePercent: 0.097745, IsBullMarket: true},
		Kline{Open: 39020.560000, Close: 40516.290000, Low: 38730.000000, High: 41064.050000, Volume: 108522.391949, Time: time.UnixMicro(1623628800000000), ChangePercent: 0.038332, IsBullMarket: true},
		Kline{Open: 40516.280000, Close: 40144.040000, Low: 39506.400000, High: 41330.000000, Volume: 80679.622838, Time: time.UnixMicro(1623715200000000), ChangePercent: -0.009187, IsBullMarket: false},
		Kline{Open: 40143.800000, Close: 38349.010000, Low: 38116.010000, High: 40527.140000, Volume: 87771.976937, Time: time.UnixMicro(1623801600000000), ChangePercent: -0.044709, IsBullMarket: false},
		Kline{Open: 38349.000000, Close: 38092.970000, Low: 37365.000000, High: 39559.880000, Volume: 79541.307119, Time: time.UnixMicro(1623888000000000), ChangePercent: -0.006676, IsBullMarket: false},
		Kline{Open: 38092.970000, Close: 35819.840000, Low: 35129.290000, High: 38202.840000, Volume: 95228.042935, Time: time.UnixMicro(1623974400000000), ChangePercent: -0.059673, IsBullMarket: false},
		Kline{Open: 35820.480000, Close: 35483.720000, Low: 34803.520000, High: 36457.000000, Volume: 68712.449461, Time: time.UnixMicro(1624060800000000), ChangePercent: -0.009401, IsBullMarket: false},
		Kline{Open: 35483.720000, Close: 35600.160000, Low: 33336.000000, High: 36137.720000, Volume: 89878.170850, Time: time.UnixMicro(1624147200000000), ChangePercent: 0.003282, IsBullMarket: true},
		Kline{Open: 35600.170000, Close: 31608.930000, Low: 31251.230000, High: 35750.000000, Volume: 168778.873159, Time: time.UnixMicro(1624233600000000), ChangePercent: -0.112113, IsBullMarket: false},
		Kline{Open: 31614.120000, Close: 32509.560000, Low: 28805.000000, High: 33298.780000, Volume: 204208.179762, Time: time.UnixMicro(1624320000000000), ChangePercent: 0.028324, IsBullMarket: true},
		Kline{Open: 32509.560000, Close: 33678.070000, Low: 31683.000000, High: 34881.000000, Volume: 126966.100563, Time: time.UnixMicro(1624406400000000), ChangePercent: 0.035944, IsBullMarket: true},
		Kline{Open: 33675.070000, Close: 34663.090000, Low: 32286.570000, High: 35298.000000, Volume: 86625.804260, Time: time.UnixMicro(1624492800000000), ChangePercent: 0.029340, IsBullMarket: true},
		Kline{Open: 34663.080000, Close: 31584.450000, Low: 31275.000000, High: 35500.000000, Volume: 116061.130356, Time: time.UnixMicro(1624579200000000), ChangePercent: -0.088816, IsBullMarket: false},
		Kline{Open: 31576.090000, Close: 32283.650000, Low: 30151.000000, High: 32730.000000, Volume: 107820.375287, Time: time.UnixMicro(1624665600000000), ChangePercent: 0.022408, IsBullMarket: true},
		Kline{Open: 32283.650000, Close: 34700.340000, Low: 31973.450000, High: 34749.000000, Volume: 96613.244211, Time: time.UnixMicro(1624752000000000), ChangePercent: 0.074858, IsBullMarket: true},
		Kline{Open: 34702.490000, Close: 34494.890000, Low: 33862.720000, High: 35297.710000, Volume: 82222.267819, Time: time.UnixMicro(1624838400000000), ChangePercent: -0.005982, IsBullMarket: false},
		Kline{Open: 34494.890000, Close: 35911.730000, Low: 34225.430000, High: 36600.000000, Volume: 90788.796220, Time: time.UnixMicro(1624924800000000), ChangePercent: 0.041074, IsBullMarket: true},
		Kline{Open: 35911.720000, Close: 35045.000000, Low: 34017.550000, High: 36100.000000, Volume: 77152.197634, Time: time.UnixMicro(1625011200000000), ChangePercent: -0.024135, IsBullMarket: false},
		Kline{Open: 35045.000000, Close: 33504.690000, Low: 32711.000000, High: 35057.570000, Volume: 71708.266112, Time: time.UnixMicro(1625097600000000), ChangePercent: -0.043952, IsBullMarket: false},
		Kline{Open: 33502.330000, Close: 33786.550000, Low: 32699.000000, High: 33977.040000, Volume: 56172.181378, Time: time.UnixMicro(1625184000000000), ChangePercent: 0.008484, IsBullMarket: true},
		Kline{Open: 33786.540000, Close: 34669.130000, Low: 33316.730000, High: 34945.610000, Volume: 43044.578641, Time: time.UnixMicro(1625270400000000), ChangePercent: 0.026123, IsBullMarket: true},
		Kline{Open: 34669.120000, Close: 35286.510000, Low: 34357.150000, High: 35967.850000, Volume: 43703.475789, Time: time.UnixMicro(1625356800000000), ChangePercent: 0.017808, IsBullMarket: true},
		Kline{Open: 35288.130000, Close: 33690.140000, Low: 33125.550000, High: 35293.780000, Volume: 64123.874245, Time: time.UnixMicro(1625443200000000), ChangePercent: -0.045284, IsBullMarket: false},
		Kline{Open: 33690.150000, Close: 34220.010000, Low: 33532.000000, High: 35118.880000, Volume: 58210.596349, Time: time.UnixMicro(1625529600000000), ChangePercent: 0.015727, IsBullMarket: true},
		Kline{Open: 34220.020000, Close: 33862.120000, Low: 33777.770000, High: 35059.090000, Volume: 53807.521675, Time: time.UnixMicro(1625616000000000), ChangePercent: -0.010459, IsBullMarket: false},
		Kline{Open: 33862.110000, Close: 32875.710000, Low: 32077.000000, High: 33929.640000, Volume: 70136.480320, Time: time.UnixMicro(1625702400000000), ChangePercent: -0.029130, IsBullMarket: false},
		Kline{Open: 32875.710000, Close: 33815.810000, Low: 32261.070000, High: 34100.000000, Volume: 47153.939899, Time: time.UnixMicro(1625788800000000), ChangePercent: 0.028596, IsBullMarket: true},
		Kline{Open: 33815.810000, Close: 33502.870000, Low: 33004.780000, High: 34262.000000, Volume: 34761.175468, Time: time.UnixMicro(1625875200000000), ChangePercent: -0.009254, IsBullMarket: false},
		Kline{Open: 33502.870000, Close: 34258.990000, Low: 33306.470000, High: 34666.000000, Volume: 31572.647448, Time: time.UnixMicro(1625961600000000), ChangePercent: 0.022569, IsBullMarket: true},
		Kline{Open: 34259.000000, Close: 33086.630000, Low: 32658.340000, High: 34678.430000, Volume: 48181.403762, Time: time.UnixMicro(1626048000000000), ChangePercent: -0.034221, IsBullMarket: false},
		Kline{Open: 33086.940000, Close: 32729.770000, Low: 32202.250000, High: 33340.000000, Volume: 41126.361008, Time: time.UnixMicro(1626134400000000), ChangePercent: -0.010795, IsBullMarket: false},
		Kline{Open: 32729.120000, Close: 32820.020000, Low: 31550.000000, High: 33114.030000, Volume: 46777.823484, Time: time.UnixMicro(1626220800000000), ChangePercent: 0.002777, IsBullMarket: true},
		Kline{Open: 32820.030000, Close: 31880.000000, Low: 31133.000000, High: 33185.250000, Volume: 51639.576353, Time: time.UnixMicro(1626307200000000), ChangePercent: -0.028642, IsBullMarket: false},
		Kline{Open: 31874.490000, Close: 31383.870000, Low: 31020.000000, High: 32249.180000, Volume: 48499.864154, Time: time.UnixMicro(1626393600000000), ChangePercent: -0.015392, IsBullMarket: false},
		Kline{Open: 31383.860000, Close: 31520.070000, Low: 31164.310000, High: 31955.920000, Volume: 34012.242132, Time: time.UnixMicro(1626480000000000), ChangePercent: 0.004340, IsBullMarket: true},
		Kline{Open: 31520.070000, Close: 31778.560000, Low: 31108.970000, High: 32435.000000, Volume: 35923.716186, Time: time.UnixMicro(1626566400000000), ChangePercent: 0.008201, IsBullMarket: true},
		Kline{Open: 31778.570000, Close: 30839.650000, Low: 30407.440000, High: 31899.000000, Volume: 47340.468499, Time: time.UnixMicro(1626652800000000), ChangePercent: -0.029546, IsBullMarket: false},
		Kline{Open: 30839.650000, Close: 29790.350000, Low: 29278.000000, High: 31063.070000, Volume: 61034.049017, Time: time.UnixMicro(1626739200000000), ChangePercent: -0.034024, IsBullMarket: false},
		Kline{Open: 29790.340000, Close: 32144.510000, Low: 29482.610000, High: 32858.000000, Volume: 82796.265128, Time: time.UnixMicro(1626825600000000), ChangePercent: 0.079025, IsBullMarket: true},
		Kline{Open: 32144.510000, Close: 32287.830000, Low: 31708.000000, High: 32591.350000, Volume: 46148.092433, Time: time.UnixMicro(1626912000000000), ChangePercent: 0.004459, IsBullMarket: true},
		Kline{Open: 32287.580000, Close: 33634.090000, Low: 31924.320000, High: 33650.000000, Volume: 50112.863626, Time: time.UnixMicro(1626998400000000), ChangePercent: 0.041704, IsBullMarket: true},
		Kline{Open: 33634.100000, Close: 34258.140000, Low: 33401.140000, High: 34500.000000, Volume: 47977.550138, Time: time.UnixMicro(1627084800000000), ChangePercent: 0.018554, IsBullMarket: true},
		Kline{Open: 34261.510000, Close: 35381.020000, Low: 33851.120000, High: 35398.000000, Volume: 47852.928313, Time: time.UnixMicro(1627171200000000), ChangePercent: 0.032675, IsBullMarket: true},
		Kline{Open: 35381.020000, Close: 37237.600000, Low: 35205.780000, High: 40550.000000, Volume: 152452.512724, Time: time.UnixMicro(1627257600000000), ChangePercent: 0.052474, IsBullMarket: true},
		Kline{Open: 37241.330000, Close: 39457.870000, Low: 36383.000000, High: 39542.610000, Volume: 88397.267015, Time: time.UnixMicro(1627344000000000), ChangePercent: 0.059518, IsBullMarket: true},
		Kline{Open: 39456.610000, Close: 40019.560000, Low: 38772.000000, High: 40900.000000, Volume: 101344.528441, Time: time.UnixMicro(1627430400000000), ChangePercent: 0.014268, IsBullMarket: true},
		Kline{Open: 40019.570000, Close: 40016.480000, Low: 39200.000000, High: 40640.000000, Volume: 53998.439283, Time: time.UnixMicro(1627516800000000), ChangePercent: -0.000077, IsBullMarket: false},
		Kline{Open: 40018.490000, Close: 42206.370000, Low: 38313.230000, High: 42316.710000, Volume: 73602.784805, Time: time.UnixMicro(1627603200000000), ChangePercent: 0.054672, IsBullMarket: true},
		Kline{Open: 42206.360000, Close: 41461.830000, Low: 41000.150000, High: 42448.000000, Volume: 44849.791012, Time: time.UnixMicro(1627689600000000), ChangePercent: -0.017640, IsBullMarket: false},
		Kline{Open: 41461.840000, Close: 39845.440000, Low: 39422.010000, High: 42599.000000, Volume: 53953.186326, Time: time.UnixMicro(1627776000000000), ChangePercent: -0.038985, IsBullMarket: false},
		Kline{Open: 39850.270000, Close: 39147.820000, Low: 38690.000000, High: 40480.010000, Volume: 50837.351954, Time: time.UnixMicro(1627862400000000), ChangePercent: -0.017627, IsBullMarket: false},
		Kline{Open: 39146.860000, Close: 38207.050000, Low: 37642.030000, High: 39780.000000, Volume: 57117.435853, Time: time.UnixMicro(1627948800000000), ChangePercent: -0.024007, IsBullMarket: false},
		Kline{Open: 38207.040000, Close: 39723.180000, Low: 37508.560000, High: 39969.660000, Volume: 52329.352430, Time: time.UnixMicro(1628035200000000), ChangePercent: 0.039682, IsBullMarket: true},
		Kline{Open: 39723.170000, Close: 40862.460000, Low: 37332.700000, High: 41350.000000, Volume: 84343.755621, Time: time.UnixMicro(1628121600000000), ChangePercent: 0.028681, IsBullMarket: true},
		Kline{Open: 40862.460000, Close: 42836.870000, Low: 39853.860000, High: 43392.430000, Volume: 75753.941347, Time: time.UnixMicro(1628208000000000), ChangePercent: 0.048318, IsBullMarket: true},
		Kline{Open: 42836.870000, Close: 44572.540000, Low: 42446.410000, High: 44700.000000, Volume: 73396.740808, Time: time.UnixMicro(1628294400000000), ChangePercent: 0.040518, IsBullMarket: true},
		Kline{Open: 44572.540000, Close: 43794.370000, Low: 43261.000000, High: 45310.000000, Volume: 69329.092698, Time: time.UnixMicro(1628380800000000), ChangePercent: -0.017459, IsBullMarket: false},
		Kline{Open: 43794.360000, Close: 46253.400000, Low: 42779.000000, High: 46454.150000, Volume: 74587.884845, Time: time.UnixMicro(1628467200000000), ChangePercent: 0.056150, IsBullMarket: true},
		Kline{Open: 46248.870000, Close: 45584.990000, Low: 44589.460000, High: 46700.000000, Volume: 53814.643421, Time: time.UnixMicro(1628553600000000), ChangePercent: -0.014355, IsBullMarket: false},
		Kline{Open: 45585.000000, Close: 45511.000000, Low: 45341.140000, High: 46743.470000, Volume: 52734.901977, Time: time.UnixMicro(1628640000000000), ChangePercent: -0.001623, IsBullMarket: false},
		Kline{Open: 45510.670000, Close: 44399.000000, Low: 43770.000000, High: 46218.120000, Volume: 55266.108781, Time: time.UnixMicro(1628726400000000), ChangePercent: -0.024427, IsBullMarket: false},
		Kline{Open: 44400.060000, Close: 47800.000000, Low: 44217.390000, High: 47886.000000, Volume: 48239.370431, Time: time.UnixMicro(1628812800000000), ChangePercent: 0.076575, IsBullMarket: true},
		Kline{Open: 47799.990000, Close: 47068.510000, Low: 45971.030000, High: 48144.000000, Volume: 46114.359022, Time: time.UnixMicro(1628899200000000), ChangePercent: -0.015303, IsBullMarket: false},
		Kline{Open: 47068.500000, Close: 46973.820000, Low: 45500.000000, High: 47372.270000, Volume: 42110.711334, Time: time.UnixMicro(1628985600000000), ChangePercent: -0.002012, IsBullMarket: false},
		Kline{Open: 46973.820000, Close: 45901.290000, Low: 45660.000000, High: 48053.830000, Volume: 52480.574014, Time: time.UnixMicro(1629072000000000), ChangePercent: -0.022833, IsBullMarket: false},
		Kline{Open: 45901.300000, Close: 44695.950000, Low: 44376.000000, High: 47160.000000, Volume: 57039.341629, Time: time.UnixMicro(1629158400000000), ChangePercent: -0.026260, IsBullMarket: false},
		Kline{Open: 44695.950000, Close: 44705.290000, Low: 44203.280000, High: 46000.000000, Volume: 54099.415985, Time: time.UnixMicro(1629244800000000), ChangePercent: 0.000209, IsBullMarket: true},
		Kline{Open: 44699.370000, Close: 46760.620000, Low: 43927.700000, High: 47033.000000, Volume: 53411.753920, Time: time.UnixMicro(1629331200000000), ChangePercent: 0.046114, IsBullMarket: true},
		Kline{Open: 46760.620000, Close: 49322.470000, Low: 46622.990000, High: 49382.990000, Volume: 56850.352228, Time: time.UnixMicro(1629417600000000), ChangePercent: 0.054786, IsBullMarket: true},
		Kline{Open: 49322.470000, Close: 48821.870000, Low: 48222.000000, High: 49757.040000, Volume: 46745.136584, Time: time.UnixMicro(1629504000000000), ChangePercent: -0.010150, IsBullMarket: false},
		Kline{Open: 48821.880000, Close: 49239.220000, Low: 48050.000000, High: 49500.000000, Volume: 37007.887795, Time: time.UnixMicro(1629590400000000), ChangePercent: 0.008548, IsBullMarket: true},
		Kline{Open: 49239.220000, Close: 49488.850000, Low: 49029.000000, High: 50500.000000, Volume: 52462.541954, Time: time.UnixMicro(1629676800000000), ChangePercent: 0.005070, IsBullMarket: true},
		Kline{Open: 49488.850000, Close: 47674.010000, Low: 47600.000000, High: 49860.000000, Volume: 51014.594748, Time: time.UnixMicro(1629763200000000), ChangePercent: -0.036672, IsBullMarket: false},
		Kline{Open: 47674.010000, Close: 48973.320000, Low: 47126.280000, High: 49264.300000, Volume: 44655.830342, Time: time.UnixMicro(1629849600000000), ChangePercent: 0.027254, IsBullMarket: true},
		Kline{Open: 48973.320000, Close: 46843.870000, Low: 46250.000000, High: 49352.840000, Volume: 49371.277774, Time: time.UnixMicro(1629936000000000), ChangePercent: -0.043482, IsBullMarket: false},
		Kline{Open: 46843.860000, Close: 49069.900000, Low: 46348.000000, High: 49149.930000, Volume: 42068.104965, Time: time.UnixMicro(1630022400000000), ChangePercent: 0.047520, IsBullMarket: true},
		Kline{Open: 49069.900000, Close: 48895.350000, Low: 48346.880000, High: 49299.000000, Volume: 26681.063786, Time: time.UnixMicro(1630108800000000), ChangePercent: -0.003557, IsBullMarket: false},
		Kline{Open: 48895.350000, Close: 48767.830000, Low: 47762.540000, High: 49632.270000, Volume: 32652.283473, Time: time.UnixMicro(1630195200000000), ChangePercent: -0.002608, IsBullMarket: false},
		Kline{Open: 48767.840000, Close: 46982.910000, Low: 46853.000000, High: 48888.610000, Volume: 40288.350830, Time: time.UnixMicro(1630281600000000), ChangePercent: -0.036601, IsBullMarket: false},
		Kline{Open: 46982.910000, Close: 47100.890000, Low: 46700.000000, High: 48246.110000, Volume: 48645.527370, Time: time.UnixMicro(1630368000000000), ChangePercent: 0.002511, IsBullMarket: true},
		Kline{Open: 47100.890000, Close: 48810.520000, Low: 46512.000000, High: 49156.000000, Volume: 49904.655280, Time: time.UnixMicro(1630454400000000), ChangePercent: 0.036297, IsBullMarket: true},
		Kline{Open: 48810.510000, Close: 49246.640000, Low: 48584.060000, High: 50450.130000, Volume: 54410.770538, Time: time.UnixMicro(1630540800000000), ChangePercent: 0.008935, IsBullMarket: true},
		Kline{Open: 49246.630000, Close: 49999.140000, Low: 48316.840000, High: 51000.000000, Volume: 59025.644157, Time: time.UnixMicro(1630627200000000), ChangePercent: 0.015280, IsBullMarket: true},
		Kline{Open: 49998.000000, Close: 49915.640000, Low: 49370.000000, High: 50535.690000, Volume: 34664.659590, Time: time.UnixMicro(1630713600000000), ChangePercent: -0.001647, IsBullMarket: false},
		Kline{Open: 49917.540000, Close: 51756.880000, Low: 49450.000000, High: 51900.000000, Volume: 40544.835873, Time: time.UnixMicro(1630800000000000), ChangePercent: 0.036848, IsBullMarket: true},
		Kline{Open: 51756.880000, Close: 52663.900000, Low: 50969.330000, High: 52780.000000, Volume: 49249.667081, Time: time.UnixMicro(1630886400000000), ChangePercent: 0.017525, IsBullMarket: true},
		Kline{Open: 52666.200000, Close: 46863.730000, Low: 42843.050000, High: 52920.000000, Volume: 123048.802719, Time: time.UnixMicro(1630972800000000), ChangePercent: -0.110174, IsBullMarket: false},
		Kline{Open: 46868.570000, Close: 46048.310000, Low: 44412.020000, High: 47340.990000, Volume: 65069.315200, Time: time.UnixMicro(1631059200000000), ChangePercent: -0.017501, IsBullMarket: false},
		Kline{Open: 46048.310000, Close: 46395.140000, Low: 45513.080000, High: 47399.970000, Volume: 50651.660020, Time: time.UnixMicro(1631145600000000), ChangePercent: 0.007532, IsBullMarket: true},
		Kline{Open: 46395.140000, Close: 44850.910000, Low: 44132.290000, High: 47033.000000, Volume: 49048.266180, Time: time.UnixMicro(1631232000000000), ChangePercent: -0.033284, IsBullMarket: false},
		Kline{Open: 44842.200000, Close: 45173.690000, Low: 44722.220000, High: 45987.930000, Volume: 30440.408100, Time: time.UnixMicro(1631318400000000), ChangePercent: 0.007392, IsBullMarket: true},
		Kline{Open: 45173.680000, Close: 46025.240000, Low: 44742.060000, High: 46460.000000, Volume: 32094.280520, Time: time.UnixMicro(1631404800000000), ChangePercent: 0.018851, IsBullMarket: true},
		Kline{Open: 46025.230000, Close: 44940.730000, Low: 43370.000000, High: 46880.000000, Volume: 65429.150560, Time: time.UnixMicro(1631491200000000), ChangePercent: -0.023563, IsBullMarket: false},
		Kline{Open: 44940.720000, Close: 47111.520000, Low: 44594.440000, High: 47250.000000, Volume: 44855.850990, Time: time.UnixMicro(1631577600000000), ChangePercent: 0.048304, IsBullMarket: true},
		Kline{Open: 47103.280000, Close: 48121.410000, Low: 46682.320000, High: 48500.000000, Volume: 43204.711740, Time: time.UnixMicro(1631664000000000), ChangePercent: 0.021615, IsBullMarket: true},
		Kline{Open: 48121.400000, Close: 47737.820000, Low: 47021.100000, High: 48557.000000, Volume: 40725.088950, Time: time.UnixMicro(1631750400000000), ChangePercent: -0.007971, IsBullMarket: false},
		Kline{Open: 47737.810000, Close: 47299.980000, Low: 46699.560000, High: 48150.000000, Volume: 34461.927760, Time: time.UnixMicro(1631836800000000), ChangePercent: -0.009172, IsBullMarket: false},
		Kline{Open: 47299.980000, Close: 48292.740000, Low: 47035.560000, High: 48843.200000, Volume: 30906.470380, Time: time.UnixMicro(1631923200000000), ChangePercent: 0.020989, IsBullMarket: true},
		Kline{Open: 48292.750000, Close: 47241.750000, Low: 46829.180000, High: 48372.830000, Volume: 29847.243490, Time: time.UnixMicro(1632009600000000), ChangePercent: -0.021763, IsBullMarket: false},
		Kline{Open: 47241.750000, Close: 43015.620000, Low: 42500.000000, High: 47347.250000, Volume: 78003.524443, Time: time.UnixMicro(1632096000000000), ChangePercent: -0.089458, IsBullMarket: false},
		Kline{Open: 43016.640000, Close: 40734.380000, Low: 39600.000000, High: 43639.000000, Volume: 84534.080485, Time: time.UnixMicro(1632182400000000), ChangePercent: -0.053055, IsBullMarket: false},
		Kline{Open: 40734.090000, Close: 43543.610000, Low: 40565.390000, High: 44000.550000, Volume: 58349.055420, Time: time.UnixMicro(1632268800000000), ChangePercent: 0.068972, IsBullMarket: true},
		Kline{Open: 43546.370000, Close: 44865.260000, Low: 43069.090000, High: 44978.000000, Volume: 48699.576550, Time: time.UnixMicro(1632355200000000), ChangePercent: 0.030287, IsBullMarket: true},
		Kline{Open: 44865.260000, Close: 42810.570000, Low: 40675.000000, High: 45200.000000, Volume: 84113.426292, Time: time.UnixMicro(1632441600000000), ChangePercent: -0.045797, IsBullMarket: false},
		Kline{Open: 42810.580000, Close: 42670.640000, Low: 41646.280000, High: 42966.840000, Volume: 33594.571890, Time: time.UnixMicro(1632528000000000), ChangePercent: -0.003269, IsBullMarket: false},
		Kline{Open: 42670.630000, Close: 43160.900000, Low: 40750.000000, High: 43950.000000, Volume: 49879.997650, Time: time.UnixMicro(1632614400000000), ChangePercent: 0.011490, IsBullMarket: true},
		Kline{Open: 43160.900000, Close: 42147.350000, Low: 42098.000000, High: 44350.000000, Volume: 39776.843830, Time: time.UnixMicro(1632700800000000), ChangePercent: -0.023483, IsBullMarket: false},
		Kline{Open: 42147.350000, Close: 41026.540000, Low: 40888.000000, High: 42787.380000, Volume: 43372.262400, Time: time.UnixMicro(1632787200000000), ChangePercent: -0.026593, IsBullMarket: false},
		Kline{Open: 41025.010000, Close: 41524.280000, Low: 40753.880000, High: 42590.000000, Volume: 33511.534870, Time: time.UnixMicro(1632873600000000), ChangePercent: 0.012170, IsBullMarket: true},
		Kline{Open: 41524.290000, Close: 43824.100000, Low: 41410.170000, High: 44141.370000, Volume: 46381.227810, Time: time.UnixMicro(1632960000000000), ChangePercent: 0.055385, IsBullMarket: true},
		Kline{Open: 43820.010000, Close: 48141.610000, Low: 43283.030000, High: 48495.000000, Volume: 66244.874920, Time: time.UnixMicro(1633046400000000), ChangePercent: 0.098622, IsBullMarket: true},
		Kline{Open: 48141.600000, Close: 47634.900000, Low: 47430.180000, High: 48336.590000, Volume: 30508.981310, Time: time.UnixMicro(1633132800000000), ChangePercent: -0.010525, IsBullMarket: false},
		Kline{Open: 47634.890000, Close: 48200.010000, Low: 47088.000000, High: 49228.080000, Volume: 30825.056010, Time: time.UnixMicro(1633219200000000), ChangePercent: 0.011864, IsBullMarket: true},
		Kline{Open: 48200.010000, Close: 49224.940000, Low: 46891.000000, High: 49536.120000, Volume: 46796.493720, Time: time.UnixMicro(1633305600000000), ChangePercent: 0.021264, IsBullMarket: true},
		Kline{Open: 49224.930000, Close: 51471.990000, Low: 49022.400000, High: 51886.300000, Volume: 52125.667930, Time: time.UnixMicro(1633392000000000), ChangePercent: 0.045649, IsBullMarket: true},
		Kline{Open: 51471.990000, Close: 55315.000000, Low: 50382.410000, High: 55750.000000, Volume: 79877.545181, Time: time.UnixMicro(1633478400000000), ChangePercent: 0.074662, IsBullMarket: true},
		Kline{Open: 55315.000000, Close: 53785.220000, Low: 53357.000000, High: 55332.310000, Volume: 54917.377660, Time: time.UnixMicro(1633564800000000), ChangePercent: -0.027656, IsBullMarket: false},
		Kline{Open: 53785.220000, Close: 53951.430000, Low: 53617.610000, High: 56100.000000, Volume: 46160.257850, Time: time.UnixMicro(1633651200000000), ChangePercent: 0.003090, IsBullMarket: true},
		Kline{Open: 53955.670000, Close: 54949.720000, Low: 53661.670000, High: 55489.000000, Volume: 55177.080130, Time: time.UnixMicro(1633737600000000), ChangePercent: 0.018423, IsBullMarket: true},
		Kline{Open: 54949.720000, Close: 54659.000000, Low: 54080.000000, High: 56561.310000, Volume: 89237.836128, Time: time.UnixMicro(1633824000000000), ChangePercent: -0.005291, IsBullMarket: false},
		Kline{Open: 54659.010000, Close: 57471.350000, Low: 54415.060000, High: 57839.040000, Volume: 52933.165751, Time: time.UnixMicro(1633910400000000), ChangePercent: 0.051452, IsBullMarket: true},
		Kline{Open: 57471.350000, Close: 55996.930000, Low: 53879.000000, High: 57680.000000, Volume: 53471.285500, Time: time.UnixMicro(1633996800000000), ChangePercent: -0.025655, IsBullMarket: false},
		Kline{Open: 55996.910000, Close: 57367.000000, Low: 54167.190000, High: 57777.000000, Volume: 55808.444920, Time: time.UnixMicro(1634083200000000), ChangePercent: 0.024467, IsBullMarket: true},
		Kline{Open: 57370.830000, Close: 57347.940000, Low: 56818.050000, High: 58532.540000, Volume: 43053.336781, Time: time.UnixMicro(1634169600000000), ChangePercent: -0.000399, IsBullMarket: false},
		Kline{Open: 57347.940000, Close: 61672.420000, Low: 56850.000000, High: 62933.000000, Volume: 82512.908022, Time: time.UnixMicro(1634256000000000), ChangePercent: 0.075408, IsBullMarket: true},
		Kline{Open: 61672.420000, Close: 60875.570000, Low: 60150.000000, High: 62378.420000, Volume: 35467.880960, Time: time.UnixMicro(1634342400000000), ChangePercent: -0.012921, IsBullMarket: false},
		Kline{Open: 60875.570000, Close: 61528.330000, Low: 58963.000000, High: 61718.390000, Volume: 39099.241240, Time: time.UnixMicro(1634428800000000), ChangePercent: 0.010723, IsBullMarket: true},
		Kline{Open: 61528.320000, Close: 62009.840000, Low: 59844.450000, High: 62695.780000, Volume: 51798.448440, Time: time.UnixMicro(1634515200000000), ChangePercent: 0.007826, IsBullMarket: true},
		Kline{Open: 62005.600000, Close: 64280.590000, Low: 61322.220000, High: 64486.000000, Volume: 53628.107744, Time: time.UnixMicro(1634601600000000), ChangePercent: 0.036690, IsBullMarket: true},
		Kline{Open: 64280.590000, Close: 66001.410000, Low: 63481.400000, High: 67000.000000, Volume: 51428.934856, Time: time.UnixMicro(1634688000000000), ChangePercent: 0.026770, IsBullMarket: true},
		Kline{Open: 66001.400000, Close: 62193.150000, Low: 62000.000000, High: 66639.740000, Volume: 68538.645370, Time: time.UnixMicro(1634774400000000), ChangePercent: -0.057700, IsBullMarket: false},
		Kline{Open: 62193.150000, Close: 60688.220000, Low: 60000.000000, High: 63732.390000, Volume: 52119.358860, Time: time.UnixMicro(1634860800000000), ChangePercent: -0.024198, IsBullMarket: false},
		Kline{Open: 60688.230000, Close: 61286.750000, Low: 59562.150000, High: 61747.640000, Volume: 27626.936780, Time: time.UnixMicro(1634947200000000), ChangePercent: 0.009862, IsBullMarket: true},
		Kline{Open: 61286.750000, Close: 60852.220000, Low: 59510.630000, High: 61500.000000, Volume: 31226.576760, Time: time.UnixMicro(1635033600000000), ChangePercent: -0.007090, IsBullMarket: false},
		Kline{Open: 60852.220000, Close: 63078.780000, Low: 60650.000000, High: 63710.630000, Volume: 36853.838060, Time: time.UnixMicro(1635120000000000), ChangePercent: 0.036590, IsBullMarket: true},
		Kline{Open: 63078.780000, Close: 60328.810000, Low: 59817.550000, High: 63293.480000, Volume: 40217.500830, Time: time.UnixMicro(1635206400000000), ChangePercent: -0.043596, IsBullMarket: false},
		Kline{Open: 60328.810000, Close: 58413.440000, Low: 58000.000000, High: 61496.000000, Volume: 62124.490160, Time: time.UnixMicro(1635292800000000), ChangePercent: -0.031749, IsBullMarket: false},
		Kline{Open: 58413.440000, Close: 60575.890000, Low: 57820.000000, High: 62499.000000, Volume: 61056.353010, Time: time.UnixMicro(1635379200000000), ChangePercent: 0.037020, IsBullMarket: true},
		Kline{Open: 60575.900000, Close: 62253.710000, Low: 60174.810000, High: 62980.000000, Volume: 43973.904140, Time: time.UnixMicro(1635465600000000), ChangePercent: 0.027698, IsBullMarket: true},
		Kline{Open: 62253.700000, Close: 61859.190000, Low: 60673.000000, High: 62359.250000, Volume: 31478.125660, Time: time.UnixMicro(1635552000000000), ChangePercent: -0.006337, IsBullMarket: false},
		Kline{Open: 61859.190000, Close: 61299.800000, Low: 59945.360000, High: 62405.300000, Volume: 39267.637940, Time: time.UnixMicro(1635638400000000), ChangePercent: -0.009043, IsBullMarket: false},
		Kline{Open: 61299.810000, Close: 60911.110000, Low: 59405.000000, High: 62437.740000, Volume: 44687.666720, Time: time.UnixMicro(1635724800000000), ChangePercent: -0.006341, IsBullMarket: false},
		Kline{Open: 60911.120000, Close: 63219.990000, Low: 60624.680000, High: 64270.000000, Volume: 46368.284100, Time: time.UnixMicro(1635811200000000), ChangePercent: 0.037906, IsBullMarket: true},
		Kline{Open: 63220.570000, Close: 62896.480000, Low: 60382.760000, High: 63500.000000, Volume: 43336.090490, Time: time.UnixMicro(1635897600000000), ChangePercent: -0.005126, IsBullMarket: false},
		Kline{Open: 62896.490000, Close: 61395.010000, Low: 60677.010000, High: 63086.310000, Volume: 35930.933140, Time: time.UnixMicro(1635984000000000), ChangePercent: -0.023872, IsBullMarket: false},
		Kline{Open: 61395.010000, Close: 60937.120000, Low: 60721.000000, High: 62595.720000, Volume: 31604.487490, Time: time.UnixMicro(1636070400000000), ChangePercent: -0.007458, IsBullMarket: false},
		Kline{Open: 60940.180000, Close: 61470.610000, Low: 60050.000000, High: 61560.490000, Volume: 25590.574080, Time: time.UnixMicro(1636156800000000), ChangePercent: 0.008704, IsBullMarket: true},
		Kline{Open: 61470.620000, Close: 63273.590000, Low: 61322.780000, High: 63286.350000, Volume: 25515.688300, Time: time.UnixMicro(1636243200000000), ChangePercent: 0.029331, IsBullMarket: true},
		Kline{Open: 63273.580000, Close: 67525.830000, Low: 63273.580000, High: 67789.000000, Volume: 54442.094554, Time: time.UnixMicro(1636329600000000), ChangePercent: 0.067204, IsBullMarket: true},
		Kline{Open: 67525.820000, Close: 66947.660000, Low: 66222.400000, High: 68524.250000, Volume: 44661.378068, Time: time.UnixMicro(1636416000000000), ChangePercent: -0.008562, IsBullMarket: false},
		Kline{Open: 66947.670000, Close: 64882.430000, Low: 62822.900000, High: 69000.000000, Volume: 65171.504046, Time: time.UnixMicro(1636502400000000), ChangePercent: -0.030849, IsBullMarket: false},
		Kline{Open: 64882.420000, Close: 64774.260000, Low: 64100.000000, High: 65600.070000, Volume: 37237.980580, Time: time.UnixMicro(1636588800000000), ChangePercent: -0.001667, IsBullMarket: false},
		Kline{Open: 64774.250000, Close: 64122.230000, Low: 62278.000000, High: 65450.700000, Volume: 44490.108160, Time: time.UnixMicro(1636675200000000), ChangePercent: -0.010066, IsBullMarket: false},
		Kline{Open: 64122.220000, Close: 64380.000000, Low: 63360.220000, High: 65000.000000, Volume: 22504.973830, Time: time.UnixMicro(1636761600000000), ChangePercent: 0.004020, IsBullMarket: true},
		Kline{Open: 64380.010000, Close: 65519.100000, Low: 63576.270000, High: 65550.510000, Volume: 25705.073470, Time: time.UnixMicro(1636848000000000), ChangePercent: 0.017693, IsBullMarket: true},
		Kline{Open: 65519.110000, Close: 63606.740000, Low: 63400.000000, High: 66401.820000, Volume: 37829.371240, Time: time.UnixMicro(1636934400000000), ChangePercent: -0.029188, IsBullMarket: false},
		Kline{Open: 63606.730000, Close: 60058.870000, Low: 58574.070000, High: 63617.310000, Volume: 77455.156090, Time: time.UnixMicro(1637020800000000), ChangePercent: -0.055778, IsBullMarket: false},
		Kline{Open: 60058.870000, Close: 60344.870000, Low: 58373.000000, High: 60840.230000, Volume: 46289.384910, Time: time.UnixMicro(1637107200000000), ChangePercent: 0.004762, IsBullMarket: true},
		Kline{Open: 60344.860000, Close: 56891.620000, Low: 56474.260000, High: 60976.000000, Volume: 62146.999310, Time: time.UnixMicro(1637193600000000), ChangePercent: -0.057225, IsBullMarket: false},
		Kline{Open: 56891.620000, Close: 58052.240000, Low: 55600.000000, High: 58320.000000, Volume: 50715.887260, Time: time.UnixMicro(1637280000000000), ChangePercent: 0.020401, IsBullMarket: true},
		Kline{Open: 58057.100000, Close: 59707.510000, Low: 57353.000000, High: 59845.000000, Volume: 33811.590100, Time: time.UnixMicro(1637366400000000), ChangePercent: 0.028427, IsBullMarket: true},
		Kline{Open: 59707.520000, Close: 58622.020000, Low: 58486.650000, High: 60029.760000, Volume: 31902.227850, Time: time.UnixMicro(1637452800000000), ChangePercent: -0.018180, IsBullMarket: false},
		Kline{Open: 58617.700000, Close: 56247.180000, Low: 55610.000000, High: 59444.000000, Volume: 51724.320470, Time: time.UnixMicro(1637539200000000), ChangePercent: -0.040440, IsBullMarket: false},
		Kline{Open: 56243.830000, Close: 57541.270000, Low: 55317.000000, High: 58009.990000, Volume: 49917.850170, Time: time.UnixMicro(1637625600000000), ChangePercent: 0.023068, IsBullMarket: true},
		Kline{Open: 57541.260000, Close: 57138.290000, Low: 55837.000000, High: 57735.000000, Volume: 39612.049640, Time: time.UnixMicro(1637712000000000), ChangePercent: -0.007003, IsBullMarket: false},
		Kline{Open: 57138.290000, Close: 58960.360000, Low: 57000.000000, High: 59398.900000, Volume: 42153.515220, Time: time.UnixMicro(1637798400000000), ChangePercent: 0.031889, IsBullMarket: true},
		Kline{Open: 58960.370000, Close: 53726.530000, Low: 53500.000000, High: 59150.000000, Volume: 65927.870660, Time: time.UnixMicro(1637884800000000), ChangePercent: -0.088769, IsBullMarket: false},
		Kline{Open: 53723.720000, Close: 54721.030000, Low: 53610.000000, High: 55280.000000, Volume: 29716.999570, Time: time.UnixMicro(1637971200000000), ChangePercent: 0.018564, IsBullMarket: true},
		Kline{Open: 54716.470000, Close: 57274.880000, Low: 53256.640000, High: 57445.050000, Volume: 36163.713700, Time: time.UnixMicro(1638057600000000), ChangePercent: 0.046758, IsBullMarket: true},
		Kline{Open: 57274.890000, Close: 57776.250000, Low: 56666.670000, High: 58865.970000, Volume: 40125.280090, Time: time.UnixMicro(1638144000000000), ChangePercent: 0.008754, IsBullMarket: true},
		Kline{Open: 57776.250000, Close: 56950.560000, Low: 55875.550000, High: 59176.990000, Volume: 49161.051940, Time: time.UnixMicro(1638230400000000), ChangePercent: -0.014291, IsBullMarket: false},
		Kline{Open: 56950.560000, Close: 57184.070000, Low: 56458.010000, High: 59053.550000, Volume: 44956.636560, Time: time.UnixMicro(1638316800000000), ChangePercent: 0.004100, IsBullMarket: true},
		Kline{Open: 57184.070000, Close: 56480.340000, Low: 55777.770000, High: 57375.470000, Volume: 37574.059760, Time: time.UnixMicro(1638403200000000), ChangePercent: -0.012306, IsBullMarket: false},
		Kline{Open: 56484.260000, Close: 53601.050000, Low: 51680.000000, High: 57600.000000, Volume: 58927.690270, Time: time.UnixMicro(1638489600000000), ChangePercent: -0.051044, IsBullMarket: false},
		Kline{Open: 53601.050000, Close: 49152.470000, Low: 42000.300000, High: 53859.100000, Volume: 114203.373748, Time: time.UnixMicro(1638576000000000), ChangePercent: -0.082994, IsBullMarket: false},
		Kline{Open: 49152.460000, Close: 49396.330000, Low: 47727.210000, High: 49699.050000, Volume: 45580.820120, Time: time.UnixMicro(1638662400000000), ChangePercent: 0.004962, IsBullMarket: true},
		Kline{Open: 49396.320000, Close: 50441.920000, Low: 47100.000000, High: 50891.110000, Volume: 58571.215750, Time: time.UnixMicro(1638748800000000), ChangePercent: 0.021168, IsBullMarket: true},
		Kline{Open: 50441.910000, Close: 50588.950000, Low: 50039.740000, High: 51936.330000, Volume: 38253.468770, Time: time.UnixMicro(1638835200000000), ChangePercent: 0.002915, IsBullMarket: true},
		Kline{Open: 50588.950000, Close: 50471.190000, Low: 48600.000000, High: 51200.000000, Volume: 38425.924660, Time: time.UnixMicro(1638921600000000), ChangePercent: -0.002328, IsBullMarket: false},
		Kline{Open: 50471.190000, Close: 47545.590000, Low: 47320.000000, High: 50797.760000, Volume: 37692.686650, Time: time.UnixMicro(1639008000000000), ChangePercent: -0.057966, IsBullMarket: false},
		Kline{Open: 47535.900000, Close: 47140.540000, Low: 46852.000000, High: 50125.000000, Volume: 44233.573910, Time: time.UnixMicro(1639094400000000), ChangePercent: -0.008317, IsBullMarket: false},
		Kline{Open: 47140.540000, Close: 49389.990000, Low: 46751.000000, High: 49485.710000, Volume: 28889.193580, Time: time.UnixMicro(1639180800000000), ChangePercent: 0.047718, IsBullMarket: true},
		Kline{Open: 49389.990000, Close: 50053.900000, Low: 48638.000000, High: 50777.000000, Volume: 26017.934210, Time: time.UnixMicro(1639267200000000), ChangePercent: 0.013442, IsBullMarket: true},
		Kline{Open: 50053.900000, Close: 46702.750000, Low: 45672.750000, High: 50189.970000, Volume: 50869.520930, Time: time.UnixMicro(1639353600000000), ChangePercent: -0.066951, IsBullMarket: false},
		Kline{Open: 46702.760000, Close: 48343.280000, Low: 46290.000000, High: 48700.410000, Volume: 39955.984450, Time: time.UnixMicro(1639440000000000), ChangePercent: 0.035127, IsBullMarket: true},
		Kline{Open: 48336.950000, Close: 48864.980000, Low: 46547.000000, High: 49500.000000, Volume: 51629.181000, Time: time.UnixMicro(1639526400000000), ChangePercent: 0.010924, IsBullMarket: true},
		Kline{Open: 48864.980000, Close: 47632.380000, Low: 47511.000000, High: 49436.430000, Volume: 31949.867390, Time: time.UnixMicro(1639612800000000), ChangePercent: -0.025225, IsBullMarket: false},
		Kline{Open: 47632.380000, Close: 46131.200000, Low: 45456.000000, High: 47995.960000, Volume: 43104.488700, Time: time.UnixMicro(1639699200000000), ChangePercent: -0.031516, IsBullMarket: false},
		Kline{Open: 46133.830000, Close: 46834.480000, Low: 45500.000000, High: 47392.370000, Volume: 25020.052710, Time: time.UnixMicro(1639785600000000), ChangePercent: 0.015187, IsBullMarket: true},
		Kline{Open: 46834.470000, Close: 46681.230000, Low: 46406.910000, High: 48300.010000, Volume: 29305.706650, Time: time.UnixMicro(1639872000000000), ChangePercent: -0.003272, IsBullMarket: false},
		Kline{Open: 46681.240000, Close: 46914.160000, Low: 45558.850000, High: 47537.570000, Volume: 35848.506090, Time: time.UnixMicro(1639958400000000), ChangePercent: 0.004990, IsBullMarket: true},
		Kline{Open: 46914.170000, Close: 48889.880000, Low: 46630.000000, High: 49328.960000, Volume: 37713.929240, Time: time.UnixMicro(1640044800000000), ChangePercent: 0.042113, IsBullMarket: true},
		Kline{Open: 48887.590000, Close: 48588.160000, Low: 48421.870000, High: 49576.130000, Volume: 27004.202200, Time: time.UnixMicro(1640131200000000), ChangePercent: -0.006125, IsBullMarket: false},
		Kline{Open: 48588.170000, Close: 50838.810000, Low: 47920.420000, High: 51375.000000, Volume: 35192.540460, Time: time.UnixMicro(1640217600000000), ChangePercent: 0.046321, IsBullMarket: true},
		Kline{Open: 50838.820000, Close: 50820.000000, Low: 50384.430000, High: 51810.000000, Volume: 31684.842690, Time: time.UnixMicro(1640304000000000), ChangePercent: -0.000370, IsBullMarket: false},
		Kline{Open: 50819.990000, Close: 50399.660000, Low: 50142.320000, High: 51156.230000, Volume: 19135.516130, Time: time.UnixMicro(1640390400000000), ChangePercent: -0.008271, IsBullMarket: false},
		Kline{Open: 50399.670000, Close: 50775.490000, Low: 49412.000000, High: 51280.000000, Volume: 22569.889140, Time: time.UnixMicro(1640476800000000), ChangePercent: 0.007457, IsBullMarket: true},
		Kline{Open: 50775.480000, Close: 50701.440000, Low: 50449.000000, High: 52088.000000, Volume: 28792.215660, Time: time.UnixMicro(1640563200000000), ChangePercent: -0.001458, IsBullMarket: false},
		Kline{Open: 50701.440000, Close: 47543.740000, Low: 47313.010000, High: 50704.050000, Volume: 45853.339240, Time: time.UnixMicro(1640649600000000), ChangePercent: -0.062280, IsBullMarket: false},
		Kline{Open: 47543.740000, Close: 46464.660000, Low: 46096.990000, High: 48139.080000, Volume: 39498.870000, Time: time.UnixMicro(1640736000000000), ChangePercent: -0.022697, IsBullMarket: false},
		Kline{Open: 46464.660000, Close: 47120.870000, Low: 45900.000000, High: 47900.000000, Volume: 30352.295690, Time: time.UnixMicro(1640822400000000), ChangePercent: 0.014123, IsBullMarket: true},
		Kline{Open: 47120.880000, Close: 46216.930000, Low: 45678.000000, High: 48548.260000, Volume: 34937.997960, Time: time.UnixMicro(1640908800000000), ChangePercent: -0.019184, IsBullMarket: false},
		Kline{Open: 46216.930000, Close: 47722.650000, Low: 46208.370000, High: 47954.630000, Volume: 19604.463250, Time: time.UnixMicro(1640995200000000), ChangePercent: 0.032579, IsBullMarket: true},
		Kline{Open: 47722.660000, Close: 47286.180000, Low: 46654.000000, High: 47990.000000, Volume: 18340.460400, Time: time.UnixMicro(1641081600000000), ChangePercent: -0.009146, IsBullMarket: false},
		Kline{Open: 47286.180000, Close: 46446.100000, Low: 45696.000000, High: 47570.000000, Volume: 27662.077100, Time: time.UnixMicro(1641168000000000), ChangePercent: -0.017766, IsBullMarket: false},
		Kline{Open: 46446.100000, Close: 45832.010000, Low: 45500.000000, High: 47557.540000, Volume: 35491.413600, Time: time.UnixMicro(1641254400000000), ChangePercent: -0.013222, IsBullMarket: false},
		Kline{Open: 45832.010000, Close: 43451.130000, Low: 42500.000000, High: 47070.000000, Volume: 51784.118570, Time: time.UnixMicro(1641340800000000), ChangePercent: -0.051948, IsBullMarket: false},
		Kline{Open: 43451.140000, Close: 43082.310000, Low: 42430.580000, High: 43816.000000, Volume: 38880.373050, Time: time.UnixMicro(1641427200000000), ChangePercent: -0.008488, IsBullMarket: false},
		Kline{Open: 43082.300000, Close: 41566.480000, Low: 40610.000000, High: 43145.830000, Volume: 54836.508180, Time: time.UnixMicro(1641513600000000), ChangePercent: -0.035184, IsBullMarket: false},
		Kline{Open: 41566.480000, Close: 41679.740000, Low: 40501.000000, High: 42300.000000, Volume: 32952.731110, Time: time.UnixMicro(1641600000000000), ChangePercent: 0.002725, IsBullMarket: true},
		Kline{Open: 41679.740000, Close: 41864.620000, Low: 41200.020000, High: 42786.700000, Volume: 22724.394260, Time: time.UnixMicro(1641686400000000), ChangePercent: 0.004436, IsBullMarket: true},
		Kline{Open: 41864.620000, Close: 41822.490000, Low: 39650.000000, High: 42248.500000, Volume: 50729.170190, Time: time.UnixMicro(1641772800000000), ChangePercent: -0.001006, IsBullMarket: false},
		Kline{Open: 41822.490000, Close: 42729.290000, Low: 41268.930000, High: 43100.000000, Volume: 37296.437290, Time: time.UnixMicro(1641859200000000), ChangePercent: 0.021682, IsBullMarket: true},
		Kline{Open: 42729.290000, Close: 43902.660000, Low: 42450.000000, High: 44322.000000, Volume: 33943.292800, Time: time.UnixMicro(1641945600000000), ChangePercent: 0.027461, IsBullMarket: true},
		Kline{Open: 43902.650000, Close: 42560.110000, Low: 42311.220000, High: 44500.000000, Volume: 34910.877620, Time: time.UnixMicro(1642032000000000), ChangePercent: -0.030580, IsBullMarket: false},
		Kline{Open: 42558.350000, Close: 43059.960000, Low: 41725.950000, High: 43448.780000, Volume: 32640.882920, Time: time.UnixMicro(1642118400000000), ChangePercent: 0.011786, IsBullMarket: true},
		Kline{Open: 43059.960000, Close: 43084.290000, Low: 42555.000000, High: 43800.000000, Volume: 21936.056160, Time: time.UnixMicro(1642204800000000), ChangePercent: 0.000565, IsBullMarket: true},
		Kline{Open: 43084.290000, Close: 43071.660000, Low: 42581.790000, High: 43475.000000, Volume: 20602.352710, Time: time.UnixMicro(1642291200000000), ChangePercent: -0.000293, IsBullMarket: false},
		Kline{Open: 43071.660000, Close: 42201.620000, Low: 41540.420000, High: 43176.180000, Volume: 27562.086130, Time: time.UnixMicro(1642377600000000), ChangePercent: -0.020200, IsBullMarket: false},
		Kline{Open: 42201.630000, Close: 42352.120000, Low: 41250.000000, High: 42691.000000, Volume: 29324.082570, Time: time.UnixMicro(1642464000000000), ChangePercent: 0.003566, IsBullMarket: true},
		Kline{Open: 42352.120000, Close: 41660.010000, Low: 41138.560000, High: 42559.130000, Volume: 31685.721590, Time: time.UnixMicro(1642550400000000), ChangePercent: -0.016342, IsBullMarket: false},
		Kline{Open: 41660.000000, Close: 40680.910000, Low: 40553.310000, High: 43505.000000, Volume: 42330.339530, Time: time.UnixMicro(1642636800000000), ChangePercent: -0.023502, IsBullMarket: false},
		Kline{Open: 40680.920000, Close: 36445.310000, Low: 35440.450000, High: 41100.000000, Volume: 88860.891999, Time: time.UnixMicro(1642723200000000), ChangePercent: -0.104118, IsBullMarket: false},
		Kline{Open: 36445.310000, Close: 35071.420000, Low: 34008.000000, High: 36835.220000, Volume: 90471.338961, Time: time.UnixMicro(1642809600000000), ChangePercent: -0.037697, IsBullMarket: false},
		Kline{Open: 35071.420000, Close: 36244.550000, Low: 34601.010000, High: 36499.000000, Volume: 44279.523540, Time: time.UnixMicro(1642896000000000), ChangePercent: 0.033450, IsBullMarket: true},
		Kline{Open: 36244.550000, Close: 36660.350000, Low: 32917.170000, High: 37550.000000, Volume: 91904.753211, Time: time.UnixMicro(1642982400000000), ChangePercent: 0.011472, IsBullMarket: true},
		Kline{Open: 36660.350000, Close: 36958.320000, Low: 35701.000000, High: 37545.140000, Volume: 49232.401830, Time: time.UnixMicro(1643068800000000), ChangePercent: 0.008128, IsBullMarket: true},
		Kline{Open: 36958.320000, Close: 36809.340000, Low: 36234.630000, High: 38919.980000, Volume: 69830.160360, Time: time.UnixMicro(1643155200000000), ChangePercent: -0.004031, IsBullMarket: false},
		Kline{Open: 36807.240000, Close: 37160.100000, Low: 35507.010000, High: 37234.470000, Volume: 53020.879340, Time: time.UnixMicro(1643241600000000), ChangePercent: 0.009587, IsBullMarket: true},
		Kline{Open: 37160.110000, Close: 37716.560000, Low: 36155.010000, High: 38000.000000, Volume: 42154.269560, Time: time.UnixMicro(1643328000000000), ChangePercent: 0.014974, IsBullMarket: true},
		Kline{Open: 37716.570000, Close: 38166.840000, Low: 37268.440000, High: 38720.740000, Volume: 26129.496820, Time: time.UnixMicro(1643414400000000), ChangePercent: 0.011938, IsBullMarket: true},
		Kline{Open: 38166.830000, Close: 37881.760000, Low: 37351.630000, High: 38359.260000, Volume: 21430.665270, Time: time.UnixMicro(1643500800000000), ChangePercent: -0.007469, IsBullMarket: false},
		Kline{Open: 37881.750000, Close: 38466.900000, Low: 36632.610000, High: 38744.000000, Volume: 36855.245800, Time: time.UnixMicro(1643587200000000), ChangePercent: 0.015447, IsBullMarket: true},
		Kline{Open: 38466.900000, Close: 38694.590000, Low: 38000.000000, High: 39265.200000, Volume: 34574.446630, Time: time.UnixMicro(1643673600000000), ChangePercent: 0.005919, IsBullMarket: true},
		Kline{Open: 38694.590000, Close: 36896.360000, Low: 36586.950000, High: 38855.920000, Volume: 35794.681300, Time: time.UnixMicro(1643760000000000), ChangePercent: -0.046472, IsBullMarket: false},
		Kline{Open: 36896.370000, Close: 37311.610000, Low: 36250.000000, High: 37387.000000, Volume: 32081.109990, Time: time.UnixMicro(1643846400000000), ChangePercent: 0.011254, IsBullMarket: true},
		Kline{Open: 37311.980000, Close: 41574.250000, Low: 37026.730000, High: 41772.330000, Volume: 64703.958740, Time: time.UnixMicro(1643932800000000), ChangePercent: 0.114233, IsBullMarket: true},
		Kline{Open: 41571.700000, Close: 41382.590000, Low: 40843.010000, High: 41913.690000, Volume: 32532.343720, Time: time.UnixMicro(1644019200000000), ChangePercent: -0.004549, IsBullMarket: false},
		Kline{Open: 41382.600000, Close: 42380.870000, Low: 41116.560000, High: 42656.000000, Volume: 22405.167040, Time: time.UnixMicro(1644105600000000), ChangePercent: 0.024123, IsBullMarket: true},
		Kline{Open: 42380.870000, Close: 43839.990000, Low: 41645.850000, High: 44500.500000, Volume: 51060.620060, Time: time.UnixMicro(1644192000000000), ChangePercent: 0.034429, IsBullMarket: true},
		Kline{Open: 43839.990000, Close: 44042.990000, Low: 42666.000000, High: 45492.000000, Volume: 64880.293870, Time: time.UnixMicro(1644278400000000), ChangePercent: 0.004630, IsBullMarket: true},
		Kline{Open: 44043.000000, Close: 44372.720000, Low: 43117.920000, High: 44799.000000, Volume: 34428.167290, Time: time.UnixMicro(1644364800000000), ChangePercent: 0.007486, IsBullMarket: true},
		Kline{Open: 44372.710000, Close: 43495.440000, Low: 43174.010000, High: 45821.000000, Volume: 62357.290910, Time: time.UnixMicro(1644451200000000), ChangePercent: -0.019770, IsBullMarket: false},
		Kline{Open: 43495.440000, Close: 42373.730000, Low: 41938.510000, High: 43920.000000, Volume: 44975.168700, Time: time.UnixMicro(1644537600000000), ChangePercent: -0.025789, IsBullMarket: false},
		Kline{Open: 42373.730000, Close: 42217.870000, Low: 41688.880000, High: 43079.490000, Volume: 26556.856810, Time: time.UnixMicro(1644624000000000), ChangePercent: -0.003678, IsBullMarket: false},
		Kline{Open: 42217.870000, Close: 42053.660000, Low: 41870.000000, High: 42760.000000, Volume: 17732.081130, Time: time.UnixMicro(1644710400000000), ChangePercent: -0.003890, IsBullMarket: false},
		Kline{Open: 42053.650000, Close: 42535.940000, Low: 41550.560000, High: 42842.400000, Volume: 34010.130600, Time: time.UnixMicro(1644796800000000), ChangePercent: 0.011468, IsBullMarket: true},
		Kline{Open: 42535.940000, Close: 44544.860000, Low: 42427.030000, High: 44751.400000, Volume: 38095.195760, Time: time.UnixMicro(1644883200000000), ChangePercent: 0.047229, IsBullMarket: true},
		Kline{Open: 44544.850000, Close: 43873.560000, Low: 43307.000000, High: 44549.970000, Volume: 28471.872700, Time: time.UnixMicro(1644969600000000), ChangePercent: -0.015070, IsBullMarket: false},
		Kline{Open: 43873.560000, Close: 40515.700000, Low: 40073.210000, High: 44164.710000, Volume: 47245.994940, Time: time.UnixMicro(1645056000000000), ChangePercent: -0.076535, IsBullMarket: false},
		Kline{Open: 40515.710000, Close: 39974.440000, Low: 39450.000000, High: 40959.880000, Volume: 43845.922410, Time: time.UnixMicro(1645142400000000), ChangePercent: -0.013360, IsBullMarket: false},
		Kline{Open: 39974.450000, Close: 40079.170000, Low: 39639.030000, High: 40444.320000, Volume: 18042.055100, Time: time.UnixMicro(1645228800000000), ChangePercent: 0.002620, IsBullMarket: true},
		Kline{Open: 40079.170000, Close: 38386.890000, Low: 38000.000000, High: 40125.440000, Volume: 33439.290110, Time: time.UnixMicro(1645315200000000), ChangePercent: -0.042223, IsBullMarket: false},
		Kline{Open: 38386.890000, Close: 37008.160000, Low: 36800.000000, High: 39494.350000, Volume: 62347.684960, Time: time.UnixMicro(1645401600000000), ChangePercent: -0.035917, IsBullMarket: false},
		Kline{Open: 37008.160000, Close: 38230.330000, Low: 36350.000000, High: 38429.000000, Volume: 53785.945890, Time: time.UnixMicro(1645488000000000), ChangePercent: 0.033024, IsBullMarket: true},
		Kline{Open: 38230.330000, Close: 37250.010000, Low: 37036.790000, High: 39249.930000, Volume: 43560.732000, Time: time.UnixMicro(1645574400000000), ChangePercent: -0.025642, IsBullMarket: false},
		Kline{Open: 37250.020000, Close: 38327.210000, Low: 34322.280000, High: 39843.000000, Volume: 120476.294580, Time: time.UnixMicro(1645660800000000), ChangePercent: 0.028918, IsBullMarket: true},
		Kline{Open: 38328.680000, Close: 39219.170000, Low: 38014.370000, High: 39683.530000, Volume: 56574.571250, Time: time.UnixMicro(1645747200000000), ChangePercent: 0.023233, IsBullMarket: true},
		Kline{Open: 39219.160000, Close: 39116.720000, Low: 38573.180000, High: 40348.450000, Volume: 29361.256800, Time: time.UnixMicro(1645833600000000), ChangePercent: -0.002612, IsBullMarket: false},
		Kline{Open: 39116.730000, Close: 37699.070000, Low: 37000.000000, High: 39855.700000, Volume: 46229.447190, Time: time.UnixMicro(1645920000000000), ChangePercent: -0.036242, IsBullMarket: false},
		Kline{Open: 37699.080000, Close: 43160.000000, Low: 37450.170000, High: 44225.840000, Volume: 73945.638580, Time: time.UnixMicro(1646006400000000), ChangePercent: 0.144856, IsBullMarket: true},
		Kline{Open: 43160.000000, Close: 44421.200000, Low: 42809.980000, High: 44949.000000, Volume: 61743.098730, Time: time.UnixMicro(1646092800000000), ChangePercent: 0.029222, IsBullMarket: true},
		Kline{Open: 44421.200000, Close: 43892.980000, Low: 43334.090000, High: 45400.000000, Volume: 57782.650810, Time: time.UnixMicro(1646179200000000), ChangePercent: -0.011891, IsBullMarket: false},
		Kline{Open: 43892.990000, Close: 42454.000000, Low: 41832.280000, High: 44101.120000, Volume: 50940.610210, Time: time.UnixMicro(1646265600000000), ChangePercent: -0.032784, IsBullMarket: false},
		Kline{Open: 42454.000000, Close: 39148.660000, Low: 38550.000000, High: 42527.300000, Volume: 61964.684980, Time: time.UnixMicro(1646352000000000), ChangePercent: -0.077857, IsBullMarket: false},
		Kline{Open: 39148.650000, Close: 39397.960000, Low: 38407.590000, High: 39613.240000, Volume: 30363.133410, Time: time.UnixMicro(1646438400000000), ChangePercent: 0.006368, IsBullMarket: true},
		Kline{Open: 39397.970000, Close: 38420.810000, Low: 38088.570000, High: 39693.870000, Volume: 39677.261580, Time: time.UnixMicro(1646524800000000), ChangePercent: -0.024802, IsBullMarket: false},
		Kline{Open: 38420.800000, Close: 37988.000000, Low: 37155.000000, High: 39547.570000, Volume: 63994.115590, Time: time.UnixMicro(1646611200000000), ChangePercent: -0.011265, IsBullMarket: false},
		Kline{Open: 37988.010000, Close: 38730.630000, Low: 37867.650000, High: 39362.080000, Volume: 55583.066380, Time: time.UnixMicro(1646697600000000), ChangePercent: 0.019549, IsBullMarket: true},
		Kline{Open: 38730.630000, Close: 41941.710000, Low: 38656.450000, High: 42594.060000, Volume: 67392.587990, Time: time.UnixMicro(1646784000000000), ChangePercent: 0.082908, IsBullMarket: true},
		Kline{Open: 41941.700000, Close: 39422.000000, Low: 38539.730000, High: 42039.630000, Volume: 71962.931540, Time: time.UnixMicro(1646870400000000), ChangePercent: -0.060076, IsBullMarket: false},
		Kline{Open: 39422.010000, Close: 38729.570000, Low: 38223.600000, High: 40236.260000, Volume: 59018.764200, Time: time.UnixMicro(1646956800000000), ChangePercent: -0.017565, IsBullMarket: false},
		Kline{Open: 38729.570000, Close: 38807.360000, Low: 38660.520000, High: 39486.710000, Volume: 24034.364320, Time: time.UnixMicro(1647043200000000), ChangePercent: 0.002009, IsBullMarket: true},
		Kline{Open: 38807.350000, Close: 37777.340000, Low: 37578.510000, High: 39310.000000, Volume: 32791.823590, Time: time.UnixMicro(1647129600000000), ChangePercent: -0.026542, IsBullMarket: false},
		Kline{Open: 37777.350000, Close: 39671.370000, Low: 37555.000000, High: 39947.120000, Volume: 46945.453750, Time: time.UnixMicro(1647216000000000), ChangePercent: 0.050136, IsBullMarket: true},
		Kline{Open: 39671.370000, Close: 39280.330000, Low: 38098.330000, High: 39887.610000, Volume: 46015.549260, Time: time.UnixMicro(1647302400000000), ChangePercent: -0.009857, IsBullMarket: false},
		Kline{Open: 39280.330000, Close: 41114.000000, Low: 38828.480000, High: 41718.000000, Volume: 88120.761670, Time: time.UnixMicro(1647388800000000), ChangePercent: 0.046682, IsBullMarket: true},
		Kline{Open: 41114.010000, Close: 40917.900000, Low: 40500.000000, High: 41478.820000, Volume: 37189.380870, Time: time.UnixMicro(1647475200000000), ChangePercent: -0.004770, IsBullMarket: false},
		Kline{Open: 40917.890000, Close: 41757.510000, Low: 40135.040000, High: 42325.020000, Volume: 45408.009690, Time: time.UnixMicro(1647561600000000), ChangePercent: 0.020520, IsBullMarket: true},
		Kline{Open: 41757.510000, Close: 42201.130000, Low: 41499.290000, High: 42400.000000, Volume: 29067.181080, Time: time.UnixMicro(1647648000000000), ChangePercent: 0.010624, IsBullMarket: true},
		Kline{Open: 42201.130000, Close: 41262.110000, Low: 40911.000000, High: 42296.260000, Volume: 30653.334680, Time: time.UnixMicro(1647734400000000), ChangePercent: -0.022251, IsBullMarket: false},
		Kline{Open: 41262.110000, Close: 41002.250000, Low: 40467.940000, High: 41544.220000, Volume: 39426.248770, Time: time.UnixMicro(1647820800000000), ChangePercent: -0.006298, IsBullMarket: false},
		Kline{Open: 41002.260000, Close: 42364.130000, Low: 40875.510000, High: 43361.000000, Volume: 59454.942940, Time: time.UnixMicro(1647907200000000), ChangePercent: 0.033215, IsBullMarket: true},
		Kline{Open: 42364.130000, Close: 42882.760000, Low: 41751.470000, High: 43025.960000, Volume: 40828.870390, Time: time.UnixMicro(1647993600000000), ChangePercent: 0.012242, IsBullMarket: true},
		Kline{Open: 42882.760000, Close: 43991.460000, Low: 42560.460000, High: 44220.890000, Volume: 56195.123740, Time: time.UnixMicro(1648080000000000), ChangePercent: 0.025854, IsBullMarket: true},
		Kline{Open: 43991.460000, Close: 44313.160000, Low: 43579.000000, High: 45094.140000, Volume: 54614.436480, Time: time.UnixMicro(1648166400000000), ChangePercent: 0.007313, IsBullMarket: true},
		Kline{Open: 44313.160000, Close: 44511.270000, Low: 44071.970000, High: 44792.990000, Volume: 23041.617410, Time: time.UnixMicro(1648252800000000), ChangePercent: 0.004471, IsBullMarket: true},
		Kline{Open: 44511.270000, Close: 46827.760000, Low: 44421.460000, High: 46999.000000, Volume: 41874.910710, Time: time.UnixMicro(1648339200000000), ChangePercent: 0.052043, IsBullMarket: true},
		Kline{Open: 46827.760000, Close: 47122.210000, Low: 46663.560000, High: 48189.840000, Volume: 58949.261400, Time: time.UnixMicro(1648425600000000), ChangePercent: 0.006288, IsBullMarket: true},
		Kline{Open: 47122.210000, Close: 47434.800000, Low: 46950.850000, High: 48096.470000, Volume: 36772.284570, Time: time.UnixMicro(1648512000000000), ChangePercent: 0.006634, IsBullMarket: true},
		Kline{Open: 47434.790000, Close: 47067.990000, Low: 46445.420000, High: 47700.220000, Volume: 40947.208500, Time: time.UnixMicro(1648598400000000), ChangePercent: -0.007733, IsBullMarket: false},
		Kline{Open: 47067.990000, Close: 45510.340000, Low: 45200.000000, High: 47600.000000, Volume: 48645.126670, Time: time.UnixMicro(1648684800000000), ChangePercent: -0.033094, IsBullMarket: false},
		Kline{Open: 45510.350000, Close: 46283.490000, Low: 44200.000000, High: 46720.090000, Volume: 56271.064740, Time: time.UnixMicro(1648771200000000), ChangePercent: 0.016988, IsBullMarket: true},
		Kline{Open: 46283.490000, Close: 45811.000000, Low: 45620.000000, High: 47213.000000, Volume: 37073.535820, Time: time.UnixMicro(1648857600000000), ChangePercent: -0.010209, IsBullMarket: false},
		Kline{Open: 45810.990000, Close: 46407.350000, Low: 45530.920000, High: 47444.110000, Volume: 33394.677940, Time: time.UnixMicro(1648944000000000), ChangePercent: 0.013018, IsBullMarket: true},
		Kline{Open: 46407.360000, Close: 46580.510000, Low: 45118.000000, High: 46890.710000, Volume: 44641.875140, Time: time.UnixMicro(1649030400000000), ChangePercent: 0.003731, IsBullMarket: true},
		Kline{Open: 46580.500000, Close: 45497.550000, Low: 45353.810000, High: 47200.000000, Volume: 42192.748520, Time: time.UnixMicro(1649116800000000), ChangePercent: -0.023249, IsBullMarket: false},
		Kline{Open: 45497.540000, Close: 43170.470000, Low: 43121.000000, High: 45507.140000, Volume: 60849.329360, Time: time.UnixMicro(1649203200000000), ChangePercent: -0.051147, IsBullMarket: false},
		Kline{Open: 43170.470000, Close: 43444.190000, Low: 42727.350000, High: 43900.990000, Volume: 37396.541560, Time: time.UnixMicro(1649289600000000), ChangePercent: 0.006340, IsBullMarket: true},
		Kline{Open: 43444.200000, Close: 42252.010000, Low: 42107.140000, High: 43970.620000, Volume: 42375.042030, Time: time.UnixMicro(1649376000000000), ChangePercent: -0.027442, IsBullMarket: false},
		Kline{Open: 42252.020000, Close: 42753.970000, Low: 42125.480000, High: 42800.000000, Volume: 17891.660470, Time: time.UnixMicro(1649462400000000), ChangePercent: 0.011880, IsBullMarket: true},
		Kline{Open: 42753.960000, Close: 42158.850000, Low: 41868.000000, High: 43410.300000, Volume: 22771.094030, Time: time.UnixMicro(1649548800000000), ChangePercent: -0.013919, IsBullMarket: false},
		Kline{Open: 42158.850000, Close: 39530.450000, Low: 39200.000000, High: 42414.710000, Volume: 63560.447210, Time: time.UnixMicro(1649635200000000), ChangePercent: -0.062345, IsBullMarket: false},
		Kline{Open: 39530.450000, Close: 40074.940000, Low: 39254.630000, High: 40699.000000, Volume: 57751.017780, Time: time.UnixMicro(1649721600000000), ChangePercent: 0.013774, IsBullMarket: true},
		Kline{Open: 40074.950000, Close: 41147.790000, Low: 39588.540000, High: 41561.310000, Volume: 41342.272540, Time: time.UnixMicro(1649808000000000), ChangePercent: 0.026771, IsBullMarket: true},
		Kline{Open: 41147.780000, Close: 39942.380000, Low: 39551.940000, High: 41500.000000, Volume: 36807.014010, Time: time.UnixMicro(1649894400000000), ChangePercent: -0.029294, IsBullMarket: false},
		Kline{Open: 39942.370000, Close: 40551.900000, Low: 39766.400000, High: 40870.360000, Volume: 24026.357390, Time: time.UnixMicro(1649980800000000), ChangePercent: 0.015260, IsBullMarket: true},
		Kline{Open: 40551.900000, Close: 40378.710000, Low: 39991.550000, High: 40709.350000, Volume: 15805.447180, Time: time.UnixMicro(1650067200000000), ChangePercent: -0.004271, IsBullMarket: false},
		Kline{Open: 40378.700000, Close: 39678.120000, Low: 39546.170000, High: 40595.670000, Volume: 19988.492590, Time: time.UnixMicro(1650153600000000), ChangePercent: -0.017350, IsBullMarket: false},
		Kline{Open: 39678.110000, Close: 40801.130000, Low: 38536.510000, High: 41116.730000, Volume: 54243.495750, Time: time.UnixMicro(1650240000000000), ChangePercent: 0.028303, IsBullMarket: true},
		Kline{Open: 40801.130000, Close: 41493.180000, Low: 40571.000000, High: 41760.000000, Volume: 35788.858430, Time: time.UnixMicro(1650326400000000), ChangePercent: 0.016962, IsBullMarket: true},
		Kline{Open: 41493.190000, Close: 41358.190000, Low: 40820.000000, High: 42199.000000, Volume: 40877.350410, Time: time.UnixMicro(1650412800000000), ChangePercent: -0.003254, IsBullMarket: false},
		Kline{Open: 41358.190000, Close: 40480.010000, Low: 39751.000000, High: 42976.000000, Volume: 59316.276570, Time: time.UnixMicro(1650499200000000), ChangePercent: -0.021234, IsBullMarket: false},
		Kline{Open: 40480.010000, Close: 39709.180000, Low: 39177.000000, High: 40795.060000, Volume: 46664.019600, Time: time.UnixMicro(1650585600000000), ChangePercent: -0.019042, IsBullMarket: false},
		Kline{Open: 39709.190000, Close: 39441.600000, Low: 39285.000000, High: 39980.000000, Volume: 20291.423750, Time: time.UnixMicro(1650672000000000), ChangePercent: -0.006739, IsBullMarket: false},
		Kline{Open: 39441.610000, Close: 39450.130000, Low: 38929.620000, High: 39940.000000, Volume: 26703.611860, Time: time.UnixMicro(1650758400000000), ChangePercent: 0.000216, IsBullMarket: true},
		Kline{Open: 39450.120000, Close: 40426.080000, Low: 38200.000000, High: 40616.000000, Volume: 63037.127840, Time: time.UnixMicro(1650844800000000), ChangePercent: 0.024739, IsBullMarket: true},
		Kline{Open: 40426.080000, Close: 38112.650000, Low: 37702.260000, High: 40797.310000, Volume: 66650.258000, Time: time.UnixMicro(1650931200000000), ChangePercent: -0.057226, IsBullMarket: false},
		Kline{Open: 38112.640000, Close: 39235.720000, Low: 37881.310000, High: 39474.720000, Volume: 57083.122720, Time: time.UnixMicro(1651017600000000), ChangePercent: 0.029467, IsBullMarket: true},
		Kline{Open: 39235.720000, Close: 39742.070000, Low: 38881.430000, High: 40372.630000, Volume: 56086.671500, Time: time.UnixMicro(1651104000000000), ChangePercent: 0.012905, IsBullMarket: true},
		Kline{Open: 39742.060000, Close: 38596.110000, Low: 38175.000000, High: 39925.250000, Volume: 51453.657150, Time: time.UnixMicro(1651190400000000), ChangePercent: -0.028835, IsBullMarket: false},
		Kline{Open: 38596.110000, Close: 37630.800000, Low: 37578.200000, High: 38795.380000, Volume: 35321.189890, Time: time.UnixMicro(1651276800000000), ChangePercent: -0.025011, IsBullMarket: false},
		Kline{Open: 37630.800000, Close: 38468.350000, Low: 37386.380000, High: 38675.000000, Volume: 38812.241040, Time: time.UnixMicro(1651363200000000), ChangePercent: 0.022257, IsBullMarket: true},
		Kline{Open: 38468.350000, Close: 38525.160000, Low: 38052.000000, High: 39167.340000, Volume: 53200.926280, Time: time.UnixMicro(1651449600000000), ChangePercent: 0.001477, IsBullMarket: true},
		Kline{Open: 38525.160000, Close: 37728.950000, Low: 37517.800000, High: 38651.510000, Volume: 40316.453580, Time: time.UnixMicro(1651536000000000), ChangePercent: -0.020667, IsBullMarket: false},
		Kline{Open: 37728.950000, Close: 39690.000000, Low: 37670.000000, High: 40023.770000, Volume: 62574.617360, Time: time.UnixMicro(1651622400000000), ChangePercent: 0.051977, IsBullMarket: true},
		Kline{Open: 39690.000000, Close: 36552.970000, Low: 35571.900000, High: 39845.510000, Volume: 88722.433550, Time: time.UnixMicro(1651708800000000), ChangePercent: -0.079038, IsBullMarket: false},
		Kline{Open: 36552.970000, Close: 36013.770000, Low: 35258.000000, High: 36675.630000, Volume: 68437.801870, Time: time.UnixMicro(1651795200000000), ChangePercent: -0.014751, IsBullMarket: false},
		Kline{Open: 36013.770000, Close: 35472.390000, Low: 34785.000000, High: 36146.300000, Volume: 34281.706820, Time: time.UnixMicro(1651881600000000), ChangePercent: -0.015033, IsBullMarket: false},
		Kline{Open: 35472.400000, Close: 34038.400000, Low: 33713.950000, High: 35514.220000, Volume: 72445.643440, Time: time.UnixMicro(1651968000000000), ChangePercent: -0.040426, IsBullMarket: false},
		Kline{Open: 34038.390000, Close: 30076.310000, Low: 30033.330000, High: 34243.150000, Volume: 191876.926428, Time: time.UnixMicro(1652054400000000), ChangePercent: -0.116400, IsBullMarket: false},
		Kline{Open: 30074.230000, Close: 31017.100000, Low: 29730.400000, High: 32658.990000, Volume: 165532.003110, Time: time.UnixMicro(1652140800000000), ChangePercent: 0.031351, IsBullMarket: true},
		Kline{Open: 31017.110000, Close: 29103.940000, Low: 27785.000000, High: 32162.590000, Volume: 207063.739278, Time: time.UnixMicro(1652227200000000), ChangePercent: -0.061681, IsBullMarket: false},
		Kline{Open: 29103.940000, Close: 29029.750000, Low: 26700.000000, High: 30243.000000, Volume: 204507.263138, Time: time.UnixMicro(1652313600000000), ChangePercent: -0.002549, IsBullMarket: false},
		Kline{Open: 29029.740000, Close: 29287.050000, Low: 28751.670000, High: 31083.370000, Volume: 97872.369570, Time: time.UnixMicro(1652400000000000), ChangePercent: 0.008864, IsBullMarket: true},
		Kline{Open: 29287.050000, Close: 30086.740000, Low: 28630.000000, High: 30343.270000, Volume: 51095.878630, Time: time.UnixMicro(1652486400000000), ChangePercent: 0.027305, IsBullMarket: true},
		Kline{Open: 30086.740000, Close: 31328.890000, Low: 29480.000000, High: 31460.000000, Volume: 46275.669120, Time: time.UnixMicro(1652572800000000), ChangePercent: 0.041286, IsBullMarket: true},
		Kline{Open: 31328.890000, Close: 29874.010000, Low: 29087.040000, High: 31328.900000, Volume: 73082.196580, Time: time.UnixMicro(1652659200000000), ChangePercent: -0.046439, IsBullMarket: false},
		Kline{Open: 29874.010000, Close: 30444.930000, Low: 29450.380000, High: 30788.370000, Volume: 56724.133070, Time: time.UnixMicro(1652745600000000), ChangePercent: 0.019111, IsBullMarket: true},
		Kline{Open: 30444.930000, Close: 28715.320000, Low: 28654.470000, High: 30709.990000, Volume: 59749.157990, Time: time.UnixMicro(1652832000000000), ChangePercent: -0.056811, IsBullMarket: false},
		Kline{Open: 28715.330000, Close: 30319.230000, Low: 28691.380000, High: 30545.180000, Volume: 67877.364150, Time: time.UnixMicro(1652918400000000), ChangePercent: 0.055855, IsBullMarket: true},
		Kline{Open: 30319.220000, Close: 29201.010000, Low: 28730.000000, High: 30777.330000, Volume: 60517.253250, Time: time.UnixMicro(1653004800000000), ChangePercent: -0.036881, IsBullMarket: false},
		Kline{Open: 29201.010000, Close: 29445.060000, Low: 28947.280000, High: 29656.180000, Volume: 20987.131240, Time: time.UnixMicro(1653091200000000), ChangePercent: 0.008358, IsBullMarket: true},
		Kline{Open: 29445.070000, Close: 30293.940000, Low: 29255.110000, High: 30487.990000, Volume: 36158.987480, Time: time.UnixMicro(1653177600000000), ChangePercent: 0.028829, IsBullMarket: true},
		Kline{Open: 30293.930000, Close: 29109.150000, Low: 28866.350000, High: 30670.510000, Volume: 63901.499320, Time: time.UnixMicro(1653264000000000), ChangePercent: -0.039109, IsBullMarket: false},
		Kline{Open: 29109.140000, Close: 29654.580000, Low: 28669.000000, High: 29845.860000, Volume: 59442.960360, Time: time.UnixMicro(1653350400000000), ChangePercent: 0.018738, IsBullMarket: true},
		Kline{Open: 29654.580000, Close: 29542.150000, Low: 29294.210000, High: 30223.740000, Volume: 59537.386590, Time: time.UnixMicro(1653436800000000), ChangePercent: -0.003791, IsBullMarket: false},
		Kline{Open: 29542.140000, Close: 29201.350000, Low: 28019.560000, High: 29886.640000, Volume: 94581.654630, Time: time.UnixMicro(1653523200000000), ChangePercent: -0.011536, IsBullMarket: false},
		Kline{Open: 29201.350000, Close: 28629.800000, Low: 28282.900000, High: 29397.660000, Volume: 90998.520100, Time: time.UnixMicro(1653609600000000), ChangePercent: -0.019573, IsBullMarket: false},
		Kline{Open: 28629.810000, Close: 29031.330000, Low: 28450.000000, High: 29266.000000, Volume: 34479.351270, Time: time.UnixMicro(1653696000000000), ChangePercent: 0.014025, IsBullMarket: true},
		Kline{Open: 29031.330000, Close: 29468.100000, Low: 28839.210000, High: 29587.780000, Volume: 27567.347640, Time: time.UnixMicro(1653782400000000), ChangePercent: 0.015045, IsBullMarket: true},
		Kline{Open: 29468.100000, Close: 31734.220000, Low: 29299.620000, High: 32222.000000, Volume: 96785.947600, Time: time.UnixMicro(1653868800000000), ChangePercent: 0.076901, IsBullMarket: true},
		Kline{Open: 31734.230000, Close: 31801.040000, Low: 31200.010000, High: 32399.000000, Volume: 62433.116320, Time: time.UnixMicro(1653955200000000), ChangePercent: 0.002105, IsBullMarket: true},
		Kline{Open: 31801.050000, Close: 29805.830000, Low: 29301.000000, High: 31982.970000, Volume: 103395.633820, Time: time.UnixMicro(1654041600000000), ChangePercent: -0.062741, IsBullMarket: false},
		Kline{Open: 29805.840000, Close: 30452.620000, Low: 29594.550000, High: 30689.000000, Volume: 56961.429280, Time: time.UnixMicro(1654128000000000), ChangePercent: 0.021700, IsBullMarket: true},
		Kline{Open: 30452.630000, Close: 29700.210000, Low: 29282.360000, High: 30699.000000, Volume: 54067.447270, Time: time.UnixMicro(1654214400000000), ChangePercent: -0.024708, IsBullMarket: false},
		Kline{Open: 29700.210000, Close: 29864.040000, Low: 29485.000000, High: 29988.880000, Volume: 25617.901130, Time: time.UnixMicro(1654300800000000), ChangePercent: 0.005516, IsBullMarket: true},
		Kline{Open: 29864.030000, Close: 29919.210000, Low: 29531.420000, High: 30189.000000, Volume: 23139.928100, Time: time.UnixMicro(1654387200000000), ChangePercent: 0.001848, IsBullMarket: true},
		Kline{Open: 29919.200000, Close: 31373.100000, Low: 29890.230000, High: 31765.640000, Volume: 68836.924560, Time: time.UnixMicro(1654473600000000), ChangePercent: 0.048594, IsBullMarket: true},
		Kline{Open: 31373.100000, Close: 31125.330000, Low: 29218.960000, High: 31589.600000, Volume: 110674.516580, Time: time.UnixMicro(1654560000000000), ChangePercent: -0.007898, IsBullMarket: false},
		Kline{Open: 31125.320000, Close: 30204.770000, Low: 29843.880000, High: 31327.220000, Volume: 68542.612760, Time: time.UnixMicro(1654646400000000), ChangePercent: -0.029576, IsBullMarket: false},
		Kline{Open: 30204.770000, Close: 30109.930000, Low: 29944.100000, High: 30700.000000, Volume: 46291.186500, Time: time.UnixMicro(1654732800000000), ChangePercent: -0.003140, IsBullMarket: false},
		Kline{Open: 30109.930000, Close: 29091.880000, Low: 28850.000000, High: 30382.800000, Volume: 76204.249800, Time: time.UnixMicro(1654819200000000), ChangePercent: -0.033811, IsBullMarket: false},
		Kline{Open: 29091.870000, Close: 28424.700000, Low: 28099.990000, High: 29440.410000, Volume: 65901.398950, Time: time.UnixMicro(1654905600000000), ChangePercent: -0.022933, IsBullMarket: false},
		Kline{Open: 28424.710000, Close: 26574.530000, Low: 26560.000000, High: 28544.960000, Volume: 92474.598809, Time: time.UnixMicro(1654992000000000), ChangePercent: -0.065091, IsBullMarket: false},
		Kline{Open: 26574.530000, Close: 22487.410000, Low: 21925.770000, High: 26895.840000, Volume: 254611.034966, Time: time.UnixMicro(1655078400000000), ChangePercent: -0.153798, IsBullMarket: false},
		Kline{Open: 22485.270000, Close: 22136.410000, Low: 20846.000000, High: 23362.880000, Volume: 187201.646710, Time: time.UnixMicro(1655164800000000), ChangePercent: -0.015515, IsBullMarket: false},
		Kline{Open: 22136.420000, Close: 22583.720000, Low: 20111.620000, High: 22800.000000, Volume: 200774.493467, Time: time.UnixMicro(1655251200000000), ChangePercent: 0.020207, IsBullMarket: true},
		Kline{Open: 22583.720000, Close: 20401.310000, Low: 20232.000000, High: 22995.730000, Volume: 99673.594290, Time: time.UnixMicro(1655337600000000), ChangePercent: -0.096636, IsBullMarket: false},
		Kline{Open: 20400.600000, Close: 20468.810000, Low: 20246.660000, High: 21365.430000, Volume: 86694.336630, Time: time.UnixMicro(1655424000000000), ChangePercent: 0.003344, IsBullMarket: true},
		Kline{Open: 20468.810000, Close: 18970.790000, Low: 17622.000000, High: 20792.060000, Volume: 196441.655524, Time: time.UnixMicro(1655510400000000), ChangePercent: -0.073185, IsBullMarket: false},
		Kline{Open: 18970.790000, Close: 20574.000000, Low: 17960.410000, High: 20815.950000, Volume: 128320.875950, Time: time.UnixMicro(1655596800000000), ChangePercent: 0.084509, IsBullMarket: true},
		Kline{Open: 20574.000000, Close: 20573.890000, Low: 19637.030000, High: 21090.000000, Volume: 109028.941540, Time: time.UnixMicro(1655683200000000), ChangePercent: -0.000005, IsBullMarket: false},
		Kline{Open: 20573.900000, Close: 20723.520000, Low: 20348.400000, High: 21723.000000, Volume: 104371.074900, Time: time.UnixMicro(1655769600000000), ChangePercent: 0.007272, IsBullMarket: true},
		Kline{Open: 20723.510000, Close: 19987.990000, Low: 19770.510000, High: 20900.000000, Volume: 92133.979380, Time: time.UnixMicro(1655856000000000), ChangePercent: -0.035492, IsBullMarket: false},
		Kline{Open: 19988.000000, Close: 21110.130000, Low: 19890.070000, High: 21233.000000, Volume: 83127.087160, Time: time.UnixMicro(1655942400000000), ChangePercent: 0.056140, IsBullMarket: true},
		Kline{Open: 21110.120000, Close: 21237.690000, Low: 20736.720000, High: 21558.410000, Volume: 77430.366220, Time: time.UnixMicro(1656028800000000), ChangePercent: 0.006043, IsBullMarket: true},
		Kline{Open: 21237.680000, Close: 21491.190000, Low: 20906.620000, High: 21614.500000, Volume: 51431.677940, Time: time.UnixMicro(1656115200000000), ChangePercent: 0.011937, IsBullMarket: true},
		Kline{Open: 21491.180000, Close: 21038.070000, Low: 20964.730000, High: 21888.000000, Volume: 53278.104640, Time: time.UnixMicro(1656201600000000), ChangePercent: -0.021084, IsBullMarket: false},
		Kline{Open: 21038.080000, Close: 20742.560000, Low: 20510.000000, High: 21539.850000, Volume: 64475.001300, Time: time.UnixMicro(1656288000000000), ChangePercent: -0.014047, IsBullMarket: false},
		Kline{Open: 20742.570000, Close: 20281.290000, Low: 20202.010000, High: 21212.100000, Volume: 63801.083200, Time: time.UnixMicro(1656374400000000), ChangePercent: -0.022238, IsBullMarket: false},
		Kline{Open: 20281.280000, Close: 20123.010000, Low: 19854.920000, High: 20432.310000, Volume: 77309.043790, Time: time.UnixMicro(1656460800000000), ChangePercent: -0.007804, IsBullMarket: false},
		Kline{Open: 20123.000000, Close: 19942.210000, Low: 18626.000000, High: 20179.080000, Volume: 93846.648060, Time: time.UnixMicro(1656547200000000), ChangePercent: -0.008984, IsBullMarket: false},
		Kline{Open: 19942.210000, Close: 19279.800000, Low: 18975.000000, High: 20918.350000, Volume: 111844.594940, Time: time.UnixMicro(1656633600000000), ChangePercent: -0.033216, IsBullMarket: false},
		Kline{Open: 19279.800000, Close: 19252.810000, Low: 18977.010000, High: 19467.390000, Volume: 46180.302100, Time: time.UnixMicro(1656720000000000), ChangePercent: -0.001400, IsBullMarket: false},
		Kline{Open: 19252.820000, Close: 19315.830000, Low: 18781.000000, High: 19647.630000, Volume: 51087.466310, Time: time.UnixMicro(1656806400000000), ChangePercent: 0.003273, IsBullMarket: true},
		Kline{Open: 19315.830000, Close: 20236.710000, Low: 19055.310000, High: 20354.010000, Volume: 74814.046010, Time: time.UnixMicro(1656892800000000), ChangePercent: 0.047675, IsBullMarket: true},
		Kline{Open: 20236.710000, Close: 20175.830000, Low: 19304.400000, High: 20750.000000, Volume: 96041.137560, Time: time.UnixMicro(1656979200000000), ChangePercent: -0.003008, IsBullMarket: false},
		Kline{Open: 20175.840000, Close: 20564.510000, Low: 19761.250000, High: 20675.220000, Volume: 82439.580800, Time: time.UnixMicro(1657065600000000), ChangePercent: 0.019264, IsBullMarket: true},
		Kline{Open: 20564.510000, Close: 21624.980000, Low: 20251.680000, High: 21838.100000, Volume: 85014.582610, Time: time.UnixMicro(1657152000000000), ChangePercent: 0.051568, IsBullMarket: true},
		Kline{Open: 21624.990000, Close: 21594.750000, Low: 21189.260000, High: 22527.370000, Volume: 403081.573490, Time: time.UnixMicro(1657238400000000), ChangePercent: -0.001398, IsBullMarket: false},
		Kline{Open: 21594.750000, Close: 21591.830000, Low: 21322.120000, High: 21980.000000, Volume: 178417.844680, Time: time.UnixMicro(1657324800000000), ChangePercent: -0.000135, IsBullMarket: false},
		Kline{Open: 21592.150000, Close: 20862.470000, Low: 20655.000000, High: 21607.650000, Volume: 192188.215560, Time: time.UnixMicro(1657411200000000), ChangePercent: -0.033794, IsBullMarket: false},
		Kline{Open: 20861.110000, Close: 19963.610000, Low: 19875.230000, High: 20868.480000, Volume: 137535.407240, Time: time.UnixMicro(1657497600000000), ChangePercent: -0.043023, IsBullMarket: false},
		Kline{Open: 19963.610000, Close: 19328.750000, Low: 19240.000000, High: 20059.420000, Volume: 139506.458620, Time: time.UnixMicro(1657584000000000), ChangePercent: -0.031801, IsBullMarket: false},
		Kline{Open: 19331.280000, Close: 20234.870000, Low: 18910.940000, High: 20366.610000, Volume: 209250.248880, Time: time.UnixMicro(1657670400000000), ChangePercent: 0.046742, IsBullMarket: true},
		Kline{Open: 20234.870000, Close: 20588.840000, Low: 19616.070000, High: 20900.000000, Volume: 174809.216960, Time: time.UnixMicro(1657756800000000), ChangePercent: 0.017493, IsBullMarket: true},
		Kline{Open: 20588.840000, Close: 20830.040000, Low: 20382.290000, High: 21200.000000, Volume: 143343.304900, Time: time.UnixMicro(1657843200000000), ChangePercent: 0.011715, IsBullMarket: true},
		Kline{Open: 20830.040000, Close: 21195.600000, Low: 20478.610000, High: 21588.940000, Volume: 121011.673930, Time: time.UnixMicro(1657929600000000), ChangePercent: 0.017550, IsBullMarket: true},
		Kline{Open: 21195.600000, Close: 20798.160000, Low: 20750.010000, High: 21684.540000, Volume: 118229.452500, Time: time.UnixMicro(1658016000000000), ChangePercent: -0.018751, IsBullMarket: false},
		Kline{Open: 20799.580000, Close: 22432.580000, Low: 20762.450000, High: 22777.630000, Volume: 239942.731320, Time: time.UnixMicro(1658102400000000), ChangePercent: 0.078511, IsBullMarket: true},
		Kline{Open: 22432.580000, Close: 23396.620000, Low: 21579.540000, High: 23800.000000, Volume: 263770.765740, Time: time.UnixMicro(1658188800000000), ChangePercent: 0.042975, IsBullMarket: true},
		Kline{Open: 23398.480000, Close: 23223.300000, Low: 22906.190000, High: 24276.740000, Volume: 238762.170940, Time: time.UnixMicro(1658275200000000), ChangePercent: -0.007487, IsBullMarket: false},
		Kline{Open: 23223.300000, Close: 23152.190000, Low: 22341.460000, High: 23442.770000, Volume: 184817.681910, Time: time.UnixMicro(1658361600000000), ChangePercent: -0.003062, IsBullMarket: false},
		Kline{Open: 23152.190000, Close: 22684.830000, Low: 22500.000000, High: 23756.490000, Volume: 171598.439660, Time: time.UnixMicro(1658448000000000), ChangePercent: -0.020186, IsBullMarket: false},
		Kline{Open: 22684.830000, Close: 22451.070000, Low: 21934.570000, High: 23000.770000, Volume: 122137.773750, Time: time.UnixMicro(1658534400000000), ChangePercent: -0.010305, IsBullMarket: false},
		Kline{Open: 22448.580000, Close: 22579.680000, Low: 22257.150000, High: 23014.640000, Volume: 115189.672770, Time: time.UnixMicro(1658620800000000), ChangePercent: 0.005840, IsBullMarket: true},
		Kline{Open: 22577.130000, Close: 21310.900000, Low: 21250.000000, High: 22666.000000, Volume: 180344.766430, Time: time.UnixMicro(1658707200000000), ChangePercent: -0.056085, IsBullMarket: false},
		Kline{Open: 21310.900000, Close: 21254.670000, Low: 20706.500000, High: 21347.820000, Volume: 177817.243260, Time: time.UnixMicro(1658793600000000), ChangePercent: -0.002639, IsBullMarket: false},
		Kline{Open: 21254.670000, Close: 22952.450000, Low: 21042.530000, High: 23112.630000, Volume: 210971.197960, Time: time.UnixMicro(1658880000000000), ChangePercent: 0.079878, IsBullMarket: true},
		Kline{Open: 22954.310000, Close: 23842.930000, Low: 22582.130000, High: 24199.720000, Volume: 236029.074100, Time: time.UnixMicro(1658966400000000), ChangePercent: 0.038713, IsBullMarket: true},
		Kline{Open: 23845.250000, Close: 23773.750000, Low: 23414.030000, High: 24442.660000, Volume: 198298.506230, Time: time.UnixMicro(1659052800000000), ChangePercent: -0.002999, IsBullMarket: false},
		Kline{Open: 23777.280000, Close: 23643.510000, Low: 23502.250000, High: 24668.000000, Volume: 151060.132110, Time: time.UnixMicro(1659139200000000), ChangePercent: -0.005626, IsBullMarket: false},
		Kline{Open: 23644.640000, Close: 23293.320000, Low: 23227.310000, High: 24194.820000, Volume: 127743.324830, Time: time.UnixMicro(1659225600000000), ChangePercent: -0.014858, IsBullMarket: false},
		Kline{Open: 23296.360000, Close: 23268.010000, Low: 22850.000000, High: 23509.680000, Volume: 144210.162190, Time: time.UnixMicro(1659312000000000), ChangePercent: -0.001217, IsBullMarket: false},
		Kline{Open: 23266.900000, Close: 22987.790000, Low: 22654.370000, High: 23459.890000, Volume: 158073.282250, Time: time.UnixMicro(1659398400000000), ChangePercent: -0.011996, IsBullMarket: false},
		Kline{Open: 22985.930000, Close: 22818.370000, Low: 22681.220000, High: 23647.680000, Volume: 145948.809950, Time: time.UnixMicro(1659484800000000), ChangePercent: -0.007290, IsBullMarket: false},
		Kline{Open: 22816.910000, Close: 22622.980000, Low: 22400.000000, High: 23223.320000, Volume: 154854.670160, Time: time.UnixMicro(1659571200000000), ChangePercent: -0.008499, IsBullMarket: false},
		Kline{Open: 22622.410000, Close: 23312.420000, Low: 22586.950000, High: 23472.860000, Volume: 175251.697490, Time: time.UnixMicro(1659657600000000), ChangePercent: 0.030501, IsBullMarket: true},
		Kline{Open: 23313.560000, Close: 22954.210000, Low: 22909.520000, High: 23354.360000, Volume: 83911.803070, Time: time.UnixMicro(1659744000000000), ChangePercent: -0.015414, IsBullMarket: false},
		Kline{Open: 22954.210000, Close: 23174.390000, Low: 22844.620000, High: 23402.000000, Volume: 88890.008770, Time: time.UnixMicro(1659830400000000), ChangePercent: 0.009592, IsBullMarket: true},
		Kline{Open: 23174.390000, Close: 23810.000000, Low: 23154.250000, High: 24245.000000, Volume: 170958.441520, Time: time.UnixMicro(1659916800000000), ChangePercent: 0.027427, IsBullMarket: true},
		Kline{Open: 23810.980000, Close: 23149.950000, Low: 22865.000000, High: 23933.250000, Volume: 143182.508580, Time: time.UnixMicro(1660003200000000), ChangePercent: -0.027762, IsBullMarket: false},
		Kline{Open: 23151.320000, Close: 23954.050000, Low: 22664.690000, High: 24226.000000, Volume: 208916.549530, Time: time.UnixMicro(1660089600000000), ChangePercent: 0.034673, IsBullMarket: true},
		Kline{Open: 23954.050000, Close: 23934.390000, Low: 23852.130000, High: 24918.540000, Volume: 249759.795570, Time: time.UnixMicro(1660176000000000), ChangePercent: -0.000821, IsBullMarket: false},
		Kline{Open: 23933.090000, Close: 24403.680000, Low: 23583.000000, High: 24456.500000, Volume: 174207.570400, Time: time.UnixMicro(1660262400000000), ChangePercent: 0.019663, IsBullMarket: true},
		Kline{Open: 24401.700000, Close: 24441.380000, Low: 24291.220000, High: 24888.000000, Volume: 152852.254350, Time: time.UnixMicro(1660348800000000), ChangePercent: 0.001626, IsBullMarket: true},
		Kline{Open: 24443.060000, Close: 24305.240000, Low: 24144.000000, High: 25047.560000, Volume: 151206.144730, Time: time.UnixMicro(1660435200000000), ChangePercent: -0.005638, IsBullMarket: false},
		Kline{Open: 24305.250000, Close: 24094.820000, Low: 23773.220000, High: 25211.320000, Volume: 242539.547580, Time: time.UnixMicro(1660521600000000), ChangePercent: -0.008658, IsBullMarket: false},
		Kline{Open: 24093.040000, Close: 23854.740000, Low: 23671.220000, High: 24247.490000, Volume: 179324.948210, Time: time.UnixMicro(1660608000000000), ChangePercent: -0.009891, IsBullMarket: false},
		Kline{Open: 23856.150000, Close: 23342.660000, Low: 23180.400000, High: 24446.710000, Volume: 210668.687660, Time: time.UnixMicro(1660694400000000), ChangePercent: -0.021524, IsBullMarket: false},
		Kline{Open: 23342.660000, Close: 23191.200000, Low: 23111.040000, High: 23600.000000, Volume: 144185.970110, Time: time.UnixMicro(1660780800000000), ChangePercent: -0.006489, IsBullMarket: false},
		Kline{Open: 23191.450000, Close: 20834.390000, Low: 20783.570000, High: 23208.670000, Volume: 283995.877470, Time: time.UnixMicro(1660867200000000), ChangePercent: -0.101635, IsBullMarket: false},
		Kline{Open: 20834.390000, Close: 21140.070000, Low: 20761.900000, High: 21382.850000, Volume: 183041.683630, Time: time.UnixMicro(1660953600000000), ChangePercent: 0.014672, IsBullMarket: true},
		Kline{Open: 21140.070000, Close: 21515.610000, Low: 21069.110000, High: 21800.000000, Volume: 159200.684100, Time: time.UnixMicro(1661040000000000), ChangePercent: 0.017764, IsBullMarket: true},
		Kline{Open: 21516.700000, Close: 21399.830000, Low: 20890.180000, High: 21548.710000, Volume: 222222.045260, Time: time.UnixMicro(1661126400000000), ChangePercent: -0.005432, IsBullMarket: false},
		Kline{Open: 21400.750000, Close: 21529.120000, Low: 20890.140000, High: 21684.870000, Volume: 200967.771640, Time: time.UnixMicro(1661212800000000), ChangePercent: 0.005998, IsBullMarket: true},
		Kline{Open: 21529.110000, Close: 21368.080000, Low: 21145.000000, High: 21900.000000, Volume: 174383.220460, Time: time.UnixMicro(1661299200000000), ChangePercent: -0.007480, IsBullMarket: false},
		Kline{Open: 21368.050000, Close: 21559.040000, Low: 21310.150000, High: 21819.880000, Volume: 169915.783010, Time: time.UnixMicro(1661385600000000), ChangePercent: 0.008938, IsBullMarket: true},
		Kline{Open: 21559.040000, Close: 20241.050000, Low: 20107.900000, High: 21886.770000, Volume: 273811.619550, Time: time.UnixMicro(1661472000000000), ChangePercent: -0.061134, IsBullMarket: false},
		Kline{Open: 20239.140000, Close: 20037.600000, Low: 19800.000000, High: 20402.930000, Volume: 162582.460320, Time: time.UnixMicro(1661558400000000), ChangePercent: -0.009958, IsBullMarket: false},
		Kline{Open: 20037.600000, Close: 19555.610000, Low: 19520.000000, High: 20171.180000, Volume: 139307.959760, Time: time.UnixMicro(1661644800000000), ChangePercent: -0.024054, IsBullMarket: false},
		Kline{Open: 19555.610000, Close: 20285.730000, Low: 19550.790000, High: 20433.620000, Volume: 210509.495450, Time: time.UnixMicro(1661731200000000), ChangePercent: 0.037336, IsBullMarket: true},
		Kline{Open: 20287.200000, Close: 19811.660000, Low: 19540.000000, High: 20576.250000, Volume: 256634.355290, Time: time.UnixMicro(1661817600000000), ChangePercent: -0.023440, IsBullMarket: false},
		Kline{Open: 19813.030000, Close: 20050.020000, Low: 19797.940000, High: 20490.000000, Volume: 276946.607650, Time: time.UnixMicro(1661904000000000), ChangePercent: 0.011961, IsBullMarket: true},
		Kline{Open: 20048.440000, Close: 20131.460000, Low: 19565.660000, High: 20208.370000, Volume: 245289.972630, Time: time.UnixMicro(1661990400000000), ChangePercent: 0.004141, IsBullMarket: true},
		Kline{Open: 20132.640000, Close: 19951.860000, Low: 19755.290000, High: 20441.260000, Volume: 245986.603300, Time: time.UnixMicro(1662076800000000), ChangePercent: -0.008979, IsBullMarket: false},
		Kline{Open: 19950.980000, Close: 19831.900000, Low: 19652.720000, High: 20055.930000, Volume: 146639.032040, Time: time.UnixMicro(1662163200000000), ChangePercent: -0.005969, IsBullMarket: false},
		Kline{Open: 19832.450000, Close: 20000.300000, Low: 19583.100000, High: 20029.230000, Volume: 145588.778930, Time: time.UnixMicro(1662249600000000), ChangePercent: 0.008463, IsBullMarket: true},
		Kline{Open: 20000.300000, Close: 19796.840000, Low: 19633.830000, High: 20057.270000, Volume: 222543.010570, Time: time.UnixMicro(1662336000000000), ChangePercent: -0.010173, IsBullMarket: false},
		Kline{Open: 19795.340000, Close: 18790.610000, Low: 18649.510000, High: 20180.000000, Volume: 356315.057180, Time: time.UnixMicro(1662422400000000), ChangePercent: -0.050756, IsBullMarket: false},
		Kline{Open: 18790.610000, Close: 19292.840000, Low: 18510.770000, High: 19464.060000, Volume: 287394.778800, Time: time.UnixMicro(1662508800000000), ChangePercent: 0.026728, IsBullMarket: true},
		Kline{Open: 19292.850000, Close: 19319.770000, Low: 19012.000000, High: 19458.250000, Volume: 262813.282730, Time: time.UnixMicro(1662595200000000), ChangePercent: 0.001395, IsBullMarket: true},
		Kline{Open: 19320.540000, Close: 21360.110000, Low: 19291.750000, High: 21597.220000, Volume: 428919.746520, Time: time.UnixMicro(1662681600000000), ChangePercent: 0.105565, IsBullMarket: true},
		Kline{Open: 21361.620000, Close: 21648.340000, Low: 21111.130000, High: 21810.800000, Volume: 307997.335040, Time: time.UnixMicro(1662768000000000), ChangePercent: 0.013422, IsBullMarket: true},
		Kline{Open: 21647.210000, Close: 21826.870000, Low: 21350.000000, High: 21860.000000, Volume: 280702.551490, Time: time.UnixMicro(1662854400000000), ChangePercent: 0.008299, IsBullMarket: true},
		Kline{Open: 21826.870000, Close: 22395.740000, Low: 21538.510000, High: 22488.000000, Volume: 395395.618280, Time: time.UnixMicro(1662940800000000), ChangePercent: 0.026063, IsBullMarket: true},
		Kline{Open: 22395.440000, Close: 20173.570000, Low: 19860.000000, High: 22799.000000, Volume: 431915.033330, Time: time.UnixMicro(1663027200000000), ChangePercent: -0.099211, IsBullMarket: false},
		Kline{Open: 20173.620000, Close: 20226.710000, Low: 19617.620000, High: 20541.480000, Volume: 340826.401510, Time: time.UnixMicro(1663113600000000), ChangePercent: 0.002632, IsBullMarket: true},
		Kline{Open: 20227.170000, Close: 19701.880000, Low: 19497.000000, High: 20330.240000, Volume: 333069.760760, Time: time.UnixMicro(1663200000000000), ChangePercent: -0.025970, IsBullMarket: false},
		Kline{Open: 19701.880000, Close: 19803.300000, Low: 19320.010000, High: 19890.000000, Volume: 283791.070640, Time: time.UnixMicro(1663286400000000), ChangePercent: 0.005148, IsBullMarket: true},
		Kline{Open: 19803.300000, Close: 20113.620000, Low: 19748.080000, High: 20189.000000, Volume: 179350.243380, Time: time.UnixMicro(1663372800000000), ChangePercent: 0.015670, IsBullMarket: true},
		Kline{Open: 20112.610000, Close: 19416.180000, Low: 19335.620000, High: 20117.260000, Volume: 254217.469040, Time: time.UnixMicro(1663459200000000), ChangePercent: -0.034627, IsBullMarket: false},
		Kline{Open: 19417.450000, Close: 19537.020000, Low: 18232.560000, High: 19686.200000, Volume: 380512.403060, Time: time.UnixMicro(1663545600000000), ChangePercent: 0.006158, IsBullMarket: true},
		Kline{Open: 19537.020000, Close: 18875.000000, Low: 18711.870000, High: 19634.620000, Volume: 324098.328600, Time: time.UnixMicro(1663632000000000), ChangePercent: -0.033885, IsBullMarket: false},
		Kline{Open: 18874.310000, Close: 18461.360000, Low: 18125.980000, High: 19956.000000, Volume: 385034.100210, Time: time.UnixMicro(1663718400000000), ChangePercent: -0.021879, IsBullMarket: false},
		Kline{Open: 18461.360000, Close: 19401.630000, Low: 18356.390000, High: 19550.170000, Volume: 379321.721110, Time: time.UnixMicro(1663804800000000), ChangePercent: 0.050932, IsBullMarket: true},
		Kline{Open: 19401.630000, Close: 19289.910000, Low: 18531.420000, High: 19500.000000, Volume: 385886.918290, Time: time.UnixMicro(1663891200000000), ChangePercent: -0.005758, IsBullMarket: false},
		Kline{Open: 19288.570000, Close: 18920.500000, Low: 18805.340000, High: 19316.140000, Volume: 239496.567460, Time: time.UnixMicro(1663977600000000), ChangePercent: -0.019082, IsBullMarket: false},
		Kline{Open: 18921.990000, Close: 18807.380000, Low: 18629.200000, High: 19180.210000, Volume: 191191.449200, Time: time.UnixMicro(1664064000000000), ChangePercent: -0.006057, IsBullMarket: false},
		Kline{Open: 18809.130000, Close: 19227.820000, Low: 18680.720000, High: 19318.960000, Volume: 439239.219430, Time: time.UnixMicro(1664150400000000), ChangePercent: 0.022260, IsBullMarket: true},
		Kline{Open: 19226.680000, Close: 19079.130000, Low: 18816.320000, High: 20385.860000, Volume: 593260.741610, Time: time.UnixMicro(1664236800000000), ChangePercent: -0.007674, IsBullMarket: false},
		Kline{Open: 19078.100000, Close: 19412.820000, Low: 18471.280000, High: 19790.000000, Volume: 521385.455470, Time: time.UnixMicro(1664323200000000), ChangePercent: 0.017545, IsBullMarket: true},
		Kline{Open: 19412.820000, Close: 19591.510000, Low: 18843.010000, High: 19645.520000, Volume: 406424.932560, Time: time.UnixMicro(1664409600000000), ChangePercent: 0.009205, IsBullMarket: true},
		Kline{Open: 19590.540000, Close: 19422.610000, Low: 19155.360000, High: 20185.000000, Volume: 444322.953400, Time: time.UnixMicro(1664496000000000), ChangePercent: -0.008572, IsBullMarket: false},
		Kline{Open: 19422.610000, Close: 19310.950000, Low: 19159.420000, High: 19484.000000, Volume: 165625.139590, Time: time.UnixMicro(1664582400000000), ChangePercent: -0.005749, IsBullMarket: false},
		Kline{Open: 19312.240000, Close: 19056.800000, Low: 18920.350000, High: 19395.910000, Volume: 206812.470320, Time: time.UnixMicro(1664668800000000), ChangePercent: -0.013227, IsBullMarket: false},
		Kline{Open: 19057.740000, Close: 19629.080000, Low: 18959.680000, High: 19719.100000, Volume: 293585.752120, Time: time.UnixMicro(1664755200000000), ChangePercent: 0.029979, IsBullMarket: true},
		Kline{Open: 19629.080000, Close: 20337.820000, Low: 19490.600000, High: 20475.000000, Volume: 327012.001270, Time: time.UnixMicro(1664841600000000), ChangePercent: 0.036107, IsBullMarket: true},
		Kline{Open: 20337.820000, Close: 20158.260000, Low: 19730.000000, High: 20365.600000, Volume: 312239.752240, Time: time.UnixMicro(1664928000000000), ChangePercent: -0.008829, IsBullMarket: false},
		Kline{Open: 20158.260000, Close: 19960.670000, Low: 19853.000000, High: 20456.600000, Volume: 320122.170200, Time: time.UnixMicro(1665014400000000), ChangePercent: -0.009802, IsBullMarket: false},
		Kline{Open: 19960.670000, Close: 19530.090000, Low: 19320.000000, High: 20068.820000, Volume: 220874.839130, Time: time.UnixMicro(1665100800000000), ChangePercent: -0.021571, IsBullMarket: false},
		Kline{Open: 19530.090000, Close: 19417.960000, Low: 19237.140000, High: 19627.380000, Volume: 102480.098420, Time: time.UnixMicro(1665187200000000), ChangePercent: -0.005741, IsBullMarket: false},
		Kline{Open: 19416.520000, Close: 19439.020000, Low: 19316.040000, High: 19558.000000, Volume: 113900.826810, Time: time.UnixMicro(1665273600000000), ChangePercent: 0.001159, IsBullMarket: true},
		Kline{Open: 19439.960000, Close: 19131.870000, Low: 19020.250000, High: 19525.000000, Volume: 212509.098490, Time: time.UnixMicro(1665360000000000), ChangePercent: -0.015848, IsBullMarket: false},
		Kline{Open: 19131.870000, Close: 19060.000000, Low: 18860.000000, High: 19268.090000, Volume: 243473.842860, Time: time.UnixMicro(1665446400000000), ChangePercent: -0.003757, IsBullMarket: false},
		Kline{Open: 19060.000000, Close: 19155.530000, Low: 18965.880000, High: 19238.310000, Volume: 213826.267310, Time: time.UnixMicro(1665532800000000), ChangePercent: 0.005012, IsBullMarket: true},
		Kline{Open: 19155.100000, Close: 19375.130000, Low: 18190.000000, High: 19513.790000, Volume: 399756.683370, Time: time.UnixMicro(1665619200000000), ChangePercent: 0.011487, IsBullMarket: true},
		Kline{Open: 19375.580000, Close: 19176.930000, Low: 19070.370000, High: 19951.870000, Volume: 351634.326010, Time: time.UnixMicro(1665705600000000), ChangePercent: -0.010253, IsBullMarket: false},
		Kline{Open: 19176.930000, Close: 19069.390000, Low: 18975.180000, High: 19227.680000, Volume: 113847.642320, Time: time.UnixMicro(1665792000000000), ChangePercent: -0.005608, IsBullMarket: false},
		Kline{Open: 19068.400000, Close: 19262.980000, Low: 19063.740000, High: 19425.840000, Volume: 131894.618850, Time: time.UnixMicro(1665878400000000), ChangePercent: 0.010204, IsBullMarket: true},
		Kline{Open: 19262.980000, Close: 19549.860000, Low: 19152.030000, High: 19676.960000, Volume: 222813.876340, Time: time.UnixMicro(1665964800000000), ChangePercent: 0.014893, IsBullMarket: true},
		Kline{Open: 19548.480000, Close: 19327.440000, Low: 19091.000000, High: 19706.660000, Volume: 260313.078480, Time: time.UnixMicro(1666051200000000), ChangePercent: -0.011307, IsBullMarket: false},
		Kline{Open: 19327.440000, Close: 19123.970000, Low: 19065.970000, High: 19360.160000, Volume: 186137.295380, Time: time.UnixMicro(1666137600000000), ChangePercent: -0.010528, IsBullMarket: false},
		Kline{Open: 19123.350000, Close: 19041.920000, Low: 18900.000000, High: 19347.820000, Volume: 223530.130680, Time: time.UnixMicro(1666224000000000), ChangePercent: -0.004258, IsBullMarket: false},
		Kline{Open: 19041.920000, Close: 19164.370000, Low: 18650.000000, High: 19250.000000, Volume: 269310.757690, Time: time.UnixMicro(1666310400000000), ChangePercent: 0.006431, IsBullMarket: true},
		Kline{Open: 19164.370000, Close: 19204.350000, Low: 19112.720000, High: 19257.000000, Volume: 110403.908370, Time: time.UnixMicro(1666396800000000), ChangePercent: 0.002086, IsBullMarket: true},
		Kline{Open: 19204.290000, Close: 19570.400000, Low: 19070.110000, High: 19695.000000, Volume: 167057.201840, Time: time.UnixMicro(1666483200000000), ChangePercent: 0.019064, IsBullMarket: true},
		Kline{Open: 19570.400000, Close: 19329.720000, Low: 19157.000000, High: 19601.150000, Volume: 256168.144670, Time: time.UnixMicro(1666569600000000), ChangePercent: -0.012298, IsBullMarket: false},
		Kline{Open: 19330.600000, Close: 20080.070000, Low: 19237.000000, High: 20415.870000, Volume: 326370.670610, Time: time.UnixMicro(1666656000000000), ChangePercent: 0.038771, IsBullMarket: true},
		Kline{Open: 20079.020000, Close: 20771.590000, Low: 20050.410000, High: 21020.000000, Volume: 380492.695760, Time: time.UnixMicro(1666742400000000), ChangePercent: 0.034492, IsBullMarket: true},
		Kline{Open: 20771.610000, Close: 20295.110000, Low: 20200.000000, High: 20872.210000, Volume: 328643.577910, Time: time.UnixMicro(1666828800000000), ChangePercent: -0.022940, IsBullMarket: false},
		Kline{Open: 20295.110000, Close: 20591.840000, Low: 20000.090000, High: 20750.000000, Volume: 287039.945690, Time: time.UnixMicro(1666915200000000), ChangePercent: 0.014621, IsBullMarket: true},
		Kline{Open: 20591.840000, Close: 20809.670000, Low: 20554.010000, High: 21085.000000, Volume: 254881.777550, Time: time.UnixMicro(1667001600000000), ChangePercent: 0.010578, IsBullMarket: true},
		Kline{Open: 20809.680000, Close: 20627.480000, Low: 20515.000000, High: 20931.210000, Volume: 192795.608860, Time: time.UnixMicro(1667088000000000), ChangePercent: -0.008756, IsBullMarket: false},
		Kline{Open: 20627.480000, Close: 20490.740000, Low: 20237.950000, High: 20845.920000, Volume: 303567.616280, Time: time.UnixMicro(1667174400000000), ChangePercent: -0.006629, IsBullMarket: false},
		Kline{Open: 20490.740000, Close: 20483.620000, Low: 20330.740000, High: 20700.000000, Volume: 279932.437710, Time: time.UnixMicro(1667260800000000), ChangePercent: -0.000347, IsBullMarket: false},
		Kline{Open: 20482.810000, Close: 20151.840000, Low: 20048.040000, High: 20800.000000, Volume: 373716.272990, Time: time.UnixMicro(1667347200000000), ChangePercent: -0.016158, IsBullMarket: false},
		Kline{Open: 20151.840000, Close: 20207.820000, Low: 20031.240000, High: 20393.320000, Volume: 319185.154400, Time: time.UnixMicro(1667433600000000), ChangePercent: 0.002778, IsBullMarket: true},
		Kline{Open: 20207.120000, Close: 21148.520000, Low: 20180.960000, High: 21302.050000, Volume: 453694.391650, Time: time.UnixMicro(1667520000000000), ChangePercent: 0.046588, IsBullMarket: true},
		Kline{Open: 21148.520000, Close: 21299.370000, Low: 21080.650000, High: 21480.650000, Volume: 245621.985250, Time: time.UnixMicro(1667606400000000), ChangePercent: 0.007133, IsBullMarket: true},
		Kline{Open: 21299.370000, Close: 20905.580000, Low: 20886.130000, High: 21365.270000, Volume: 230036.972170, Time: time.UnixMicro(1667692800000000), ChangePercent: -0.018488, IsBullMarket: false},
		Kline{Open: 20905.580000, Close: 20591.130000, Low: 20384.890000, High: 21069.770000, Volume: 386977.603370, Time: time.UnixMicro(1667779200000000), ChangePercent: -0.015041, IsBullMarket: false},
		Kline{Open: 20590.670000, Close: 18547.230000, Low: 17166.830000, High: 20700.880000, Volume: 760705.362783, Time: time.UnixMicro(1667865600000000), ChangePercent: -0.099241, IsBullMarket: false},
		Kline{Open: 18545.380000, Close: 15922.810000, Low: 15588.000000, High: 18587.760000, Volume: 731926.929729, Time: time.UnixMicro(1667952000000000), ChangePercent: -0.141414, IsBullMarket: false},
		Kline{Open: 15922.680000, Close: 17601.150000, Low: 15754.260000, High: 18199.000000, Volume: 608448.364320, Time: time.UnixMicro(1668038400000000), ChangePercent: 0.105414, IsBullMarket: true},
		Kline{Open: 17602.450000, Close: 17070.310000, Low: 16361.600000, High: 17695.000000, Volume: 393552.864920, Time: time.UnixMicro(1668124800000000), ChangePercent: -0.030231, IsBullMarket: false},
		Kline{Open: 17069.980000, Close: 16812.080000, Low: 16631.390000, High: 17119.100000, Volume: 167819.960350, Time: time.UnixMicro(1668211200000000), ChangePercent: -0.015108, IsBullMarket: false},
		Kline{Open: 16813.160000, Close: 16329.850000, Low: 16229.000000, High: 16954.280000, Volume: 184960.788460, Time: time.UnixMicro(1668297600000000), ChangePercent: -0.028746, IsBullMarket: false},
		Kline{Open: 16331.780000, Close: 16619.460000, Low: 15815.210000, High: 17190.000000, Volume: 380210.777500, Time: time.UnixMicro(1668384000000000), ChangePercent: 0.017615, IsBullMarket: true},
		Kline{Open: 16617.720000, Close: 16900.570000, Low: 16527.720000, High: 17134.690000, Volume: 282461.843910, Time: time.UnixMicro(1668470400000000), ChangePercent: 0.017021, IsBullMarket: true},
		Kline{Open: 16900.570000, Close: 16662.760000, Low: 16378.610000, High: 17015.920000, Volume: 261493.408090, Time: time.UnixMicro(1668556800000000), ChangePercent: -0.014071, IsBullMarket: false},
		Kline{Open: 16661.610000, Close: 16692.560000, Low: 16410.740000, High: 16751.000000, Volume: 228038.978730, Time: time.UnixMicro(1668643200000000), ChangePercent: 0.001858, IsBullMarket: true},
		Kline{Open: 16692.560000, Close: 16700.450000, Low: 16546.040000, High: 17011.000000, Volume: 214224.181840, Time: time.UnixMicro(1668729600000000), ChangePercent: 0.000473, IsBullMarket: true},
		Kline{Open: 16699.430000, Close: 16700.680000, Low: 16553.530000, High: 16822.410000, Volume: 104963.155580, Time: time.UnixMicro(1668816000000000), ChangePercent: 0.000075, IsBullMarket: true},
		Kline{Open: 16700.680000, Close: 16280.230000, Low: 16180.000000, High: 16753.330000, Volume: 154842.134780, Time: time.UnixMicro(1668902400000000), ChangePercent: -0.025176, IsBullMarket: false},
		Kline{Open: 16279.500000, Close: 15781.290000, Low: 15476.000000, High: 16319.000000, Volume: 324096.997753, Time: time.UnixMicro(1668988800000000), ChangePercent: -0.030604, IsBullMarket: false},
		Kline{Open: 15781.290000, Close: 16226.940000, Low: 15616.630000, High: 16315.000000, Volume: 239548.066230, Time: time.UnixMicro(1669075200000000), ChangePercent: 0.028239, IsBullMarket: true},
		Kline{Open: 16227.960000, Close: 16603.110000, Low: 16160.200000, High: 16706.000000, Volume: 264927.704080, Time: time.UnixMicro(1669161600000000), ChangePercent: 0.023118, IsBullMarket: true},
		Kline{Open: 16603.110000, Close: 16598.950000, Low: 16458.050000, High: 16812.630000, Volume: 206565.923460, Time: time.UnixMicro(1669248000000000), ChangePercent: -0.000251, IsBullMarket: false},
		Kline{Open: 16599.550000, Close: 16522.140000, Low: 16342.810000, High: 16666.000000, Volume: 182089.495330, Time: time.UnixMicro(1669334400000000), ChangePercent: -0.004663, IsBullMarket: false},
		Kline{Open: 16521.350000, Close: 16458.570000, Low: 16385.000000, High: 16701.990000, Volume: 181804.816660, Time: time.UnixMicro(1669420800000000), ChangePercent: -0.003800, IsBullMarket: false},
		Kline{Open: 16457.610000, Close: 16428.780000, Low: 16401.000000, High: 16600.000000, Volume: 162025.476070, Time: time.UnixMicro(1669507200000000), ChangePercent: -0.001752, IsBullMarket: false},
		Kline{Open: 16428.770000, Close: 16212.910000, Low: 15995.270000, High: 16487.040000, Volume: 252695.403670, Time: time.UnixMicro(1669593600000000), ChangePercent: -0.013139, IsBullMarket: false},
		Kline{Open: 16212.180000, Close: 16442.530000, Low: 16100.000000, High: 16548.710000, Volume: 248106.250090, Time: time.UnixMicro(1669680000000000), ChangePercent: 0.014208, IsBullMarket: true},
		Kline{Open: 16442.910000, Close: 17163.640000, Low: 16428.300000, High: 17249.000000, Volume: 303019.807190, Time: time.UnixMicro(1669766400000000), ChangePercent: 0.043832, IsBullMarket: true},
		Kline{Open: 17165.530000, Close: 16977.370000, Low: 16855.010000, High: 17324.000000, Volume: 232818.182180, Time: time.UnixMicro(1669852800000000), ChangePercent: -0.010962, IsBullMarket: false},
		Kline{Open: 16978.000000, Close: 17092.740000, Low: 16787.850000, High: 17105.730000, Volume: 202372.206200, Time: time.UnixMicro(1669939200000000), ChangePercent: 0.006758, IsBullMarket: true},
		Kline{Open: 17092.130000, Close: 16885.200000, Low: 16858.740000, High: 17188.980000, Volume: 154542.573060, Time: time.UnixMicro(1670025600000000), ChangePercent: -0.012107, IsBullMarket: false},
		Kline{Open: 16885.200000, Close: 17105.700000, Low: 16878.250000, High: 17202.840000, Volume: 178619.133870, Time: time.UnixMicro(1670112000000000), ChangePercent: 0.013059, IsBullMarket: true},
		Kline{Open: 17106.650000, Close: 16966.350000, Low: 16867.000000, High: 17424.250000, Volume: 233703.292250, Time: time.UnixMicro(1670198400000000), ChangePercent: -0.008201, IsBullMarket: false},
		Kline{Open: 16966.350000, Close: 17088.960000, Low: 16906.370000, High: 17107.010000, Volume: 218730.768830, Time: time.UnixMicro(1670284800000000), ChangePercent: 0.007227, IsBullMarket: true},
		Kline{Open: 17088.960000, Close: 16836.640000, Low: 16678.830000, High: 17142.210000, Volume: 220657.413340, Time: time.UnixMicro(1670371200000000), ChangePercent: -0.014765, IsBullMarket: false},
		Kline{Open: 16836.640000, Close: 17224.100000, Low: 16733.490000, High: 17299.000000, Volume: 236417.978040, Time: time.UnixMicro(1670457600000000), ChangePercent: 0.023013, IsBullMarket: true},
		Kline{Open: 17224.100000, Close: 17128.560000, Low: 17058.210000, High: 17360.000000, Volume: 238422.064650, Time: time.UnixMicro(1670544000000000), ChangePercent: -0.005547, IsBullMarket: false},
		Kline{Open: 17128.560000, Close: 17127.490000, Low: 17092.000000, High: 17227.720000, Volume: 140573.979370, Time: time.UnixMicro(1670630400000000), ChangePercent: -0.000062, IsBullMarket: false},
		Kline{Open: 17127.490000, Close: 17085.050000, Low: 17071.000000, High: 17270.990000, Volume: 155286.478710, Time: time.UnixMicro(1670716800000000), ChangePercent: -0.002478, IsBullMarket: false},
		Kline{Open: 17085.050000, Close: 17209.830000, Low: 16871.850000, High: 17241.890000, Volume: 226227.496940, Time: time.UnixMicro(1670803200000000), ChangePercent: 0.007303, IsBullMarket: true},
		Kline{Open: 17208.930000, Close: 17774.700000, Low: 17080.140000, High: 18000.000000, Volume: 284462.913900, Time: time.UnixMicro(1670889600000000), ChangePercent: 0.032877, IsBullMarket: true},
		Kline{Open: 17775.820000, Close: 17803.150000, Low: 17660.940000, High: 18387.950000, Volume: 266681.222090, Time: time.UnixMicro(1670976000000000), ChangePercent: 0.001537, IsBullMarket: true},
		Kline{Open: 17804.010000, Close: 17356.340000, Low: 17275.510000, High: 17854.820000, Volume: 223701.968820, Time: time.UnixMicro(1671062400000000), ChangePercent: -0.025144, IsBullMarket: false},
		Kline{Open: 17356.960000, Close: 16632.120000, Low: 16527.320000, High: 17531.730000, Volume: 253379.041160, Time: time.UnixMicro(1671148800000000), ChangePercent: -0.041761, IsBullMarket: false},
		Kline{Open: 16631.500000, Close: 16776.520000, Low: 16579.850000, High: 16796.820000, Volume: 144825.666840, Time: time.UnixMicro(1671235200000000), ChangePercent: 0.008720, IsBullMarket: true},
		Kline{Open: 16777.540000, Close: 16738.210000, Low: 16663.070000, High: 16863.260000, Volume: 112619.316460, Time: time.UnixMicro(1671321600000000), ChangePercent: -0.002344, IsBullMarket: false},
		Kline{Open: 16739.000000, Close: 16438.880000, Low: 16256.300000, High: 16815.990000, Volume: 179094.283050, Time: time.UnixMicro(1671408000000000), ChangePercent: -0.017929, IsBullMarket: false},
		Kline{Open: 16438.880000, Close: 16895.560000, Low: 16397.200000, High: 17061.270000, Volume: 248808.923240, Time: time.UnixMicro(1671494400000000), ChangePercent: 0.027780, IsBullMarket: true},
		Kline{Open: 16896.150000, Close: 16824.670000, Low: 16723.000000, High: 16925.000000, Volume: 156810.963620, Time: time.UnixMicro(1671580800000000), ChangePercent: -0.004231, IsBullMarket: false},
		Kline{Open: 16824.680000, Close: 16821.430000, Low: 16559.850000, High: 16868.520000, Volume: 176044.272350, Time: time.UnixMicro(1671667200000000), ChangePercent: -0.000193, IsBullMarket: false},
		Kline{Open: 16821.900000, Close: 16778.500000, Low: 16731.130000, High: 16955.140000, Volume: 161612.009470, Time: time.UnixMicro(1671753600000000), ChangePercent: -0.002580, IsBullMarket: false},
		Kline{Open: 16778.520000, Close: 16836.120000, Low: 16776.620000, High: 16869.990000, Volume: 100224.290960, Time: time.UnixMicro(1671840000000000), ChangePercent: 0.003433, IsBullMarket: true},
		Kline{Open: 16835.730000, Close: 16832.110000, Low: 16721.000000, High: 16857.960000, Volume: 125441.072020, Time: time.UnixMicro(1671926400000000), ChangePercent: -0.000215, IsBullMarket: false},
		Kline{Open: 16832.110000, Close: 16919.390000, Low: 16791.000000, High: 16944.520000, Volume: 124564.006560, Time: time.UnixMicro(1672012800000000), ChangePercent: 0.005185, IsBullMarket: true},
		Kline{Open: 16919.390000, Close: 16706.360000, Low: 16592.370000, High: 16972.830000, Volume: 173749.586160, Time: time.UnixMicro(1672099200000000), ChangePercent: -0.012591, IsBullMarket: false},
		Kline{Open: 16706.060000, Close: 16547.310000, Low: 16465.330000, High: 16785.190000, Volume: 193037.565770, Time: time.UnixMicro(1672185600000000), ChangePercent: -0.009503, IsBullMarket: false},
		Kline{Open: 16547.320000, Close: 16633.470000, Low: 16488.910000, High: 16664.410000, Volume: 160998.471580, Time: time.UnixMicro(1672272000000000), ChangePercent: 0.005206, IsBullMarket: true},
		Kline{Open: 16633.470000, Close: 16607.480000, Low: 16333.000000, High: 16677.350000, Volume: 164916.311740, Time: time.UnixMicro(1672358400000000), ChangePercent: -0.001563, IsBullMarket: false},
		Kline{Open: 16607.480000, Close: 16542.400000, Low: 16470.000000, High: 16644.090000, Volume: 114490.428640, Time: time.UnixMicro(1672444800000000), ChangePercent: -0.003919, IsBullMarket: false},
		Kline{Open: 16541.770000, Close: 16616.750000, Low: 16499.010000, High: 16628.000000, Volume: 96925.413740, Time: time.UnixMicro(1672531200000000), ChangePercent: 0.004533, IsBullMarket: true},
		Kline{Open: 16617.170000, Close: 16672.870000, Low: 16548.700000, High: 16799.230000, Volume: 121888.571910, Time: time.UnixMicro(1672617600000000), ChangePercent: 0.003352, IsBullMarket: true},
		Kline{Open: 16672.780000, Close: 16675.180000, Low: 16605.280000, High: 16778.400000, Volume: 159541.537330, Time: time.UnixMicro(1672704000000000), ChangePercent: 0.000144, IsBullMarket: true},
		Kline{Open: 16675.650000, Close: 16850.360000, Low: 16652.660000, High: 16991.870000, Volume: 220362.188620, Time: time.UnixMicro(1672790400000000), ChangePercent: 0.010477, IsBullMarket: true},
		Kline{Open: 16850.360000, Close: 16831.850000, Low: 16753.000000, High: 16879.820000, Volume: 163473.566410, Time: time.UnixMicro(1672876800000000), ChangePercent: -0.001098, IsBullMarket: false},
		Kline{Open: 16831.850000, Close: 16950.650000, Low: 16679.000000, High: 17041.000000, Volume: 207401.284150, Time: time.UnixMicro(1672963200000000), ChangePercent: 0.007058, IsBullMarket: true},
		Kline{Open: 16950.310000, Close: 16943.570000, Low: 16908.000000, High: 16981.910000, Volume: 104526.568800, Time: time.UnixMicro(1673049600000000), ChangePercent: -0.000398, IsBullMarket: false},
		Kline{Open: 16943.830000, Close: 17127.830000, Low: 16911.000000, High: 17176.990000, Volume: 135155.896950, Time: time.UnixMicro(1673136000000000), ChangePercent: 0.010859, IsBullMarket: true},
		Kline{Open: 17127.830000, Close: 17178.260000, Low: 17104.660000, High: 17398.800000, Volume: 266211.527230, Time: time.UnixMicro(1673222400000000), ChangePercent: 0.002944, IsBullMarket: true},
		Kline{Open: 17179.040000, Close: 17440.660000, Low: 17146.340000, High: 17499.000000, Volume: 221382.425810, Time: time.UnixMicro(1673308800000000), ChangePercent: 0.015229, IsBullMarket: true},
		Kline{Open: 17440.640000, Close: 17943.260000, Low: 17315.600000, High: 18000.000000, Volume: 262221.606530, Time: time.UnixMicro(1673395200000000), ChangePercent: 0.028819, IsBullMarket: true},
		Kline{Open: 17943.260000, Close: 18846.620000, Low: 17892.050000, High: 19117.040000, Volume: 454568.321780, Time: time.UnixMicro(1673481600000000), ChangePercent: 0.050345, IsBullMarket: true},
		Kline{Open: 18846.620000, Close: 19930.010000, Low: 18714.120000, High: 20000.000000, Volume: 368615.878230, Time: time.UnixMicro(1673568000000000), ChangePercent: 0.057485, IsBullMarket: true},
		Kline{Open: 19930.010000, Close: 20954.920000, Low: 19888.050000, High: 21258.000000, Volume: 393913.749510, Time: time.UnixMicro(1673654400000000), ChangePercent: 0.051425, IsBullMarket: true},
		Kline{Open: 20952.760000, Close: 20871.500000, Low: 20551.010000, High: 21050.740000, Volume: 178542.225490, Time: time.UnixMicro(1673740800000000), ChangePercent: -0.003878, IsBullMarket: false},
		Kline{Open: 20872.990000, Close: 21185.650000, Low: 20611.480000, High: 21474.050000, Volume: 293078.082620, Time: time.UnixMicro(1673827200000000), ChangePercent: 0.014979, IsBullMarket: true},
		Kline{Open: 21185.650000, Close: 21134.810000, Low: 20841.310000, High: 21647.450000, Volume: 275407.744090, Time: time.UnixMicro(1673913600000000), ChangePercent: -0.002400, IsBullMarket: false},
		Kline{Open: 21132.290000, Close: 20677.470000, Low: 20407.150000, High: 21650.000000, Volume: 350916.019490, Time: time.UnixMicro(1674000000000000), ChangePercent: -0.021523, IsBullMarket: false},
		Kline{Open: 20677.470000, Close: 21071.590000, Low: 20659.190000, High: 21192.000000, Volume: 251385.849250, Time: time.UnixMicro(1674086400000000), ChangePercent: 0.019060, IsBullMarket: true},
		Kline{Open: 21071.590000, Close: 22667.210000, Low: 20861.280000, High: 22755.930000, Volume: 338079.136590, Time: time.UnixMicro(1674172800000000), ChangePercent: 0.075724, IsBullMarket: true},
		Kline{Open: 22666.000000, Close: 22783.550000, Low: 22422.000000, High: 23371.800000, Volume: 346445.484320, Time: time.UnixMicro(1674259200000000), ChangePercent: 0.005186, IsBullMarket: true},
		Kline{Open: 22783.350000, Close: 22707.880000, Low: 22292.370000, High: 23078.710000, Volume: 253577.752860, Time: time.UnixMicro(1674345600000000), ChangePercent: -0.003313, IsBullMarket: false},
		Kline{Open: 22706.020000, Close: 22916.450000, Low: 22500.000000, High: 23180.000000, Volume: 293588.379380, Time: time.UnixMicro(1674432000000000), ChangePercent: 0.009268, IsBullMarket: true},
		Kline{Open: 22917.810000, Close: 22632.890000, Low: 22462.930000, High: 23162.200000, Volume: 293158.782540, Time: time.UnixMicro(1674518400000000), ChangePercent: -0.012432, IsBullMarket: false},
		Kline{Open: 22631.940000, Close: 23060.940000, Low: 22300.000000, High: 23816.730000, Volume: 346042.832230, Time: time.UnixMicro(1674604800000000), ChangePercent: 0.018956, IsBullMarket: true},
		Kline{Open: 23060.420000, Close: 23009.650000, Low: 22850.010000, High: 23282.470000, Volume: 288924.435810, Time: time.UnixMicro(1674691200000000), ChangePercent: -0.002202, IsBullMarket: false},
		Kline{Open: 23009.650000, Close: 23074.160000, Low: 22534.880000, High: 23500.000000, Volume: 280833.863150, Time: time.UnixMicro(1674777600000000), ChangePercent: 0.002804, IsBullMarket: true},
		Kline{Open: 23074.160000, Close: 23022.600000, Low: 22878.460000, High: 23189.000000, Volume: 148115.710850, Time: time.UnixMicro(1674864000000000), ChangePercent: -0.002235, IsBullMarket: false},
		Kline{Open: 23021.400000, Close: 23742.300000, Low: 22967.760000, High: 23960.540000, Volume: 295688.792040, Time: time.UnixMicro(1674950400000000), ChangePercent: 0.031314, IsBullMarket: true},
		Kline{Open: 23743.370000, Close: 22826.150000, Low: 22500.000000, High: 23800.510000, Volume: 302405.901210, Time: time.UnixMicro(1675036800000000), ChangePercent: -0.038631, IsBullMarket: false},
		Kline{Open: 22827.380000, Close: 23125.130000, Low: 22714.770000, High: 23320.000000, Volume: 264649.349090, Time: time.UnixMicro(1675123200000000), ChangePercent: 0.013044, IsBullMarket: true},
		Kline{Open: 23125.130000, Close: 23732.660000, Low: 22760.230000, High: 23812.660000, Volume: 310790.422710, Time: time.UnixMicro(1675209600000000), ChangePercent: 0.026271, IsBullMarket: true},
		Kline{Open: 23731.410000, Close: 23488.940000, Low: 23363.270000, High: 24255.000000, Volume: 364177.207510, Time: time.UnixMicro(1675296000000000), ChangePercent: -0.010217, IsBullMarket: false},
		Kline{Open: 23489.330000, Close: 23431.900000, Low: 23204.620000, High: 23715.700000, Volume: 332571.029040, Time: time.UnixMicro(1675382400000000), ChangePercent: -0.002445, IsBullMarket: false},
		Kline{Open: 23431.900000, Close: 23326.840000, Low: 23253.960000, High: 23587.780000, Volume: 166126.472950, Time: time.UnixMicro(1675468800000000), ChangePercent: -0.004484, IsBullMarket: false},
		Kline{Open: 23327.660000, Close: 22932.910000, Low: 22743.000000, High: 23433.330000, Volume: 209251.339170, Time: time.UnixMicro(1675555200000000), ChangePercent: -0.016922, IsBullMarket: false},
		Kline{Open: 22932.910000, Close: 22762.520000, Low: 22628.130000, High: 23158.250000, Volume: 265371.606900, Time: time.UnixMicro(1675641600000000), ChangePercent: -0.007430, IsBullMarket: false},
		Kline{Open: 22762.520000, Close: 23240.460000, Low: 22745.780000, High: 23350.250000, Volume: 308006.724820, Time: time.UnixMicro(1675728000000000), ChangePercent: 0.020997, IsBullMarket: true},
		Kline{Open: 23242.420000, Close: 22963.000000, Low: 22665.850000, High: 23452.000000, Volume: 280056.307170, Time: time.UnixMicro(1675814400000000), ChangePercent: -0.012022, IsBullMarket: false},
		Kline{Open: 22961.850000, Close: 21796.350000, Low: 21688.000000, High: 23011.390000, Volume: 402894.695500, Time: time.UnixMicro(1675900800000000), ChangePercent: -0.050758, IsBullMarket: false},
		Kline{Open: 21797.830000, Close: 21625.190000, Low: 21451.000000, High: 21938.160000, Volume: 338591.942470, Time: time.UnixMicro(1675987200000000), ChangePercent: -0.007920, IsBullMarket: false},
		Kline{Open: 21625.190000, Close: 21862.550000, Low: 21599.780000, High: 21906.320000, Volume: 177021.584330, Time: time.UnixMicro(1676073600000000), ChangePercent: 0.010976, IsBullMarket: true},
		Kline{Open: 21862.020000, Close: 21783.540000, Low: 21630.000000, High: 22090.000000, Volume: 204435.651630, Time: time.UnixMicro(1676160000000000), ChangePercent: -0.003590, IsBullMarket: false},
		Kline{Open: 21782.370000, Close: 21773.970000, Low: 21351.070000, High: 21894.990000, Volume: 295730.767910, Time: time.UnixMicro(1676246400000000), ChangePercent: -0.000386, IsBullMarket: false},
		Kline{Open: 21774.630000, Close: 22199.840000, Low: 21532.770000, High: 22319.080000, Volume: 361958.401090, Time: time.UnixMicro(1676332800000000), ChangePercent: 0.019528, IsBullMarket: true},
		Kline{Open: 22199.840000, Close: 24324.050000, Low: 22047.280000, High: 24380.000000, Volume: 375669.161110, Time: time.UnixMicro(1676419200000000), ChangePercent: 0.095686, IsBullMarket: true},
		Kline{Open: 24322.870000, Close: 23517.720000, Low: 23505.250000, High: 25250.000000, Volume: 450080.683660, Time: time.UnixMicro(1676505600000000), ChangePercent: -0.033103, IsBullMarket: false},
		Kline{Open: 23517.720000, Close: 24569.970000, Low: 23339.370000, High: 25021.110000, Volume: 496813.213760, Time: time.UnixMicro(1676592000000000), ChangePercent: 0.044743, IsBullMarket: true},
		Kline{Open: 24568.240000, Close: 24631.950000, Low: 24430.000000, High: 24877.000000, Volume: 216917.252130, Time: time.UnixMicro(1676678400000000), ChangePercent: 0.002593, IsBullMarket: true},
		Kline{Open: 24632.050000, Close: 24271.760000, Low: 24192.570000, High: 25192.000000, Volume: 300395.995420, Time: time.UnixMicro(1676764800000000), ChangePercent: -0.014627, IsBullMarket: false},
		Kline{Open: 24272.510000, Close: 24842.200000, Low: 23840.830000, High: 25121.230000, Volume: 346938.569970, Time: time.UnixMicro(1676851200000000), ChangePercent: 0.023471, IsBullMarket: true},
		Kline{Open: 24843.890000, Close: 24452.160000, Low: 24148.340000, High: 25250.000000, Volume: 376000.828680, Time: time.UnixMicro(1676937600000000), ChangePercent: -0.015768, IsBullMarket: false},
		Kline{Open: 24450.670000, Close: 24182.210000, Low: 23574.690000, High: 24476.050000, Volume: 379425.753650, Time: time.UnixMicro(1677024000000000), ChangePercent: -0.010980, IsBullMarket: false},
		Kline{Open: 24182.210000, Close: 23940.200000, Low: 23608.000000, High: 24599.590000, Volume: 398400.454370, Time: time.UnixMicro(1677110400000000), ChangePercent: -0.010008, IsBullMarket: false},
		Kline{Open: 23940.200000, Close: 23185.290000, Low: 22841.190000, High: 24132.350000, Volume: 343582.574530, Time: time.UnixMicro(1677196800000000), ChangePercent: -0.031533, IsBullMarket: false},
		Kline{Open: 23184.040000, Close: 23157.070000, Low: 22722.000000, High: 23219.130000, Volume: 191311.810100, Time: time.UnixMicro(1677283200000000), ChangePercent: -0.001163, IsBullMarket: false},
		Kline{Open: 23157.070000, Close: 23554.850000, Low: 23059.180000, High: 23689.990000, Volume: 202323.736230, Time: time.UnixMicro(1677369600000000), ChangePercent: 0.017177, IsBullMarket: true},
		Kline{Open: 23554.850000, Close: 23492.090000, Low: 23106.770000, High: 23897.990000, Volume: 283706.085900, Time: time.UnixMicro(1677456000000000), ChangePercent: -0.002664, IsBullMarket: false},
		Kline{Open: 23492.090000, Close: 23141.570000, Low: 23020.970000, High: 23600.000000, Volume: 264140.998940, Time: time.UnixMicro(1677542400000000), ChangePercent: -0.014921, IsBullMarket: false},
		Kline{Open: 23141.570000, Close: 23628.970000, Low: 23020.030000, High: 24000.000000, Volume: 315287.417370, Time: time.UnixMicro(1677628800000000), ChangePercent: 0.021062, IsBullMarket: true},
		Kline{Open: 23629.760000, Close: 23465.320000, Low: 23195.900000, High: 23796.930000, Volume: 239315.452190, Time: time.UnixMicro(1677715200000000), ChangePercent: -0.006959, IsBullMarket: false},
		Kline{Open: 23465.320000, Close: 22354.340000, Low: 21971.130000, High: 23476.950000, Volume: 319954.197850, Time: time.UnixMicro(1677801600000000), ChangePercent: -0.047346, IsBullMarket: false},
		Kline{Open: 22354.340000, Close: 22346.570000, Low: 22157.080000, High: 22410.000000, Volume: 121257.381320, Time: time.UnixMicro(1677888000000000), ChangePercent: -0.000348, IsBullMarket: false},
		Kline{Open: 22346.570000, Close: 22430.240000, Low: 22189.220000, High: 22662.090000, Volume: 154841.757860, Time: time.UnixMicro(1677974400000000), ChangePercent: 0.003744, IsBullMarket: true},
		Kline{Open: 22430.240000, Close: 22410.000000, Low: 22258.000000, High: 22602.190000, Volume: 203751.829570, Time: time.UnixMicro(1678060800000000), ChangePercent: -0.000902, IsBullMarket: false},
		Kline{Open: 22409.410000, Close: 22197.960000, Low: 21927.000000, High: 22557.910000, Volume: 292519.809120, Time: time.UnixMicro(1678147200000000), ChangePercent: -0.009436, IsBullMarket: false},
		Kline{Open: 22198.560000, Close: 21705.440000, Low: 21580.000000, High: 22287.000000, Volume: 301460.572720, Time: time.UnixMicro(1678233600000000), ChangePercent: -0.022214, IsBullMarket: false},
		Kline{Open: 21704.370000, Close: 20362.220000, Low: 20042.720000, High: 21834.990000, Volume: 443658.285840, Time: time.UnixMicro(1678320000000000), ChangePercent: -0.061838, IsBullMarket: false},
		Kline{Open: 20362.210000, Close: 20150.690000, Low: 19549.090000, High: 20367.780000, Volume: 618456.467100, Time: time.UnixMicro(1678406400000000), ChangePercent: -0.010388, IsBullMarket: false},
		Kline{Open: 20150.690000, Close: 20455.730000, Low: 19765.030000, High: 20686.510000, Volume: 427831.821330, Time: time.UnixMicro(1678492800000000), ChangePercent: 0.015138, IsBullMarket: true},
		Kline{Open: 20455.730000, Close: 21997.110000, Low: 20270.600000, High: 22150.000000, Volume: 430944.942880, Time: time.UnixMicro(1678579200000000), ChangePercent: 0.075352, IsBullMarket: true},
		Kline{Open: 21998.050000, Close: 24113.480000, Low: 21813.880000, High: 24500.000000, Volume: 687889.312590, Time: time.UnixMicro(1678665600000000), ChangePercent: 0.096164, IsBullMarket: true},
		Kline{Open: 24112.270000, Close: 24670.410000, Low: 23976.420000, High: 26386.870000, Volume: 699360.934230, Time: time.UnixMicro(1678752000000000), ChangePercent: 0.023148, IsBullMarket: true},
		Kline{Open: 24670.410000, Close: 24285.660000, Low: 23896.950000, High: 25196.970000, Volume: 581450.729840, Time: time.UnixMicro(1678838400000000), ChangePercent: -0.015596, IsBullMarket: false},
		Kline{Open: 24285.660000, Close: 24998.780000, Low: 24123.000000, High: 25167.400000, Volume: 439421.329980, Time: time.UnixMicro(1678924800000000), ChangePercent: 0.029364, IsBullMarket: true},
		Kline{Open: 24998.780000, Close: 27395.130000, Low: 24890.000000, High: 27756.840000, Volume: 624460.680910, Time: time.UnixMicro(1679011200000000), ChangePercent: 0.095859, IsBullMarket: true},
		Kline{Open: 27395.130000, Close: 26907.490000, Low: 26578.000000, High: 27724.850000, Volume: 371238.971740, Time: time.UnixMicro(1679097600000000), ChangePercent: -0.017800, IsBullMarket: false},
		Kline{Open: 26907.490000, Close: 27972.870000, Low: 26827.220000, High: 28390.100000, Volume: 372066.990540, Time: time.UnixMicro(1679184000000000), ChangePercent: 0.039594, IsBullMarket: true},
		Kline{Open: 27972.870000, Close: 27717.010000, Low: 27124.470000, High: 28472.000000, Volume: 477378.233730, Time: time.UnixMicro(1679270400000000), ChangePercent: -0.009147, IsBullMarket: false},
		Kline{Open: 27717.010000, Close: 28105.470000, Low: 27303.100000, High: 28438.550000, Volume: 420929.742200, Time: time.UnixMicro(1679356800000000), ChangePercent: 0.014015, IsBullMarket: true},
		Kline{Open: 28107.810000, Close: 27250.970000, Low: 26601.800000, High: 28868.050000, Volume: 224113.412960, Time: time.UnixMicro(1679443200000000), ChangePercent: -0.030484, IsBullMarket: false},
		Kline{Open: 27250.970000, Close: 28295.410000, Low: 27105.000000, High: 28750.000000, Volume: 128649.608180, Time: time.UnixMicro(1679529600000000), ChangePercent: 0.038327, IsBullMarket: true},
		Kline{Open: 28295.420000, Close: 27454.470000, Low: 27000.000000, High: 28374.300000, Volume: 86242.065440, Time: time.UnixMicro(1679616000000000), ChangePercent: -0.029720, IsBullMarket: false},
		Kline{Open: 27454.460000, Close: 27462.950000, Low: 27156.090000, High: 27787.330000, Volume: 50844.081020, Time: time.UnixMicro(1679702400000000), ChangePercent: 0.000309, IsBullMarket: true},
		Kline{Open: 27462.960000, Close: 27968.050000, Low: 27417.760000, High: 28194.400000, Volume: 49671.703530, Time: time.UnixMicro(1679788800000000), ChangePercent: 0.018392, IsBullMarket: true},
		Kline{Open: 27968.050000, Close: 26822.480000, Low: 26656.000000, High: 28023.860000, Volume: 50348.512590, Time: time.UnixMicro(1679875200000000), ChangePercent: -0.040960, IsBullMarket: false},
	}
}

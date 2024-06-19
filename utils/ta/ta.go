package ta

import (
	"math"

	"github.com/idoall/stockindicator/container/bst"
)

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

// 简单移动均线简写为SMA，有时候也直接记为MA。 移动平均线，SMA(N)它将指定周期内的收盘价格之和除以周期N得到的一个指标
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
		if i < 1 {
			continue
		}
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

// 计算平均偏差
func MeanDeviation(prices []float64, period int, ma []float64) []float64 {
	md := make([]float64, len(prices))
	for i := period - 1; i < len(prices); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += math.Abs(prices[j] - ma[i])
		}
		md[i] = sum / float64(period)
	}
	return md
}

// 计算CCI
func CCI(hlc3 []float64, period, smaPeriod int) []float64 {
	ccis := make([]float64, len(hlc3))

	// 计算典型价格、移动平均和平均偏差
	typicalPrices := make([]float64, len(hlc3))
	for i := 0; i < len(hlc3); i++ {
		// 计算典型价格
		typicalPrices[i] = hlc3[i]
	}
	ma := Sma(smaPeriod, typicalPrices)
	md := MeanDeviation(typicalPrices, period, ma)

	// 计算CCI值
	for i := period - 1; i < len(hlc3); i++ {
		typicalPrice := typicalPrices[i]
		cci := (typicalPrice - ma[i]) / (0.015 * md[i])
		ccis[i] = cci
	}

	return ccis
}

// Atr（真实波动幅度均值）返回真实范围的RMA。
//
//	真实波动幅度是max(high - low, abs(high - close[1]), abs(low - close[1]))。
func Atr(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	inTimePeriodF := float64(inTimePeriod)

	if inTimePeriod < 1 {
		return outReal
	}

	if inTimePeriod <= 1 {
		return TRange(inHigh, inLow, inClose)
	}

	outIdx := inTimePeriod
	today := inTimePeriod + 1

	tr := TRange(inHigh, inLow, inClose)
	prevATRTemp := Rma(inTimePeriod, tr)
	prevATR := prevATRTemp[inTimePeriod]
	outReal[inTimePeriod] = prevATR

	for outIdx = inTimePeriod + 1; outIdx < len(inClose); outIdx++ {
		prevATR *= inTimePeriodF - 1.0
		prevATR += tr[today]
		prevATR /= inTimePeriodF
		outReal[outIdx] = prevATR
		today++
	}

	return outReal
}

// Natr - Normalized Average True Range
func Natr(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	if inTimePeriod < 1 {
		return outReal
	}

	if inTimePeriod <= 1 {
		return TRange(inHigh, inLow, inClose)
	}

	inTimePeriodF := float64(inTimePeriod)
	outIdx := inTimePeriod
	today := inTimePeriod

	tr := TRange(inHigh, inLow, inClose)
	prevATRTemp := Sma(inTimePeriod, tr)
	prevATR := prevATRTemp[inTimePeriod]

	tempValue := inClose[today]
	if tempValue != 0.0 {
		outReal[outIdx] = (prevATR / tempValue) * 100.0
	} else {
		outReal[outIdx] = 0.0
	}

	for outIdx = inTimePeriod + 1; outIdx < len(inClose); outIdx++ {
		today++
		prevATR *= inTimePeriodF - 1.0
		prevATR += tr[today]
		prevATR /= inTimePeriodF
		tempValue = inClose[today]
		if tempValue != 0.0 {
			outReal[outIdx] = (prevATR / tempValue) * 100.0
		} else {
			outReal[0] = 0.0
		}
	}

	return outReal
}

// Cum 计算souce的累计总和,返回`source`所有元素的总和
func Cum(source []float64) float64 {
	var val float64
	for _, v := range source {
		val += v
	}
	return val
}

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

// Highest 返回给定数目的最高值。
//
//	Args:
//		values 要计算的数组
//		lenght 计算的长度
func Highest(values []float64, lenght int) float64 {
	var result = values[len(values)-1]
	for i := len(values) - 1; i >= len(values)-lenght; i-- {
		if i < 0 {
			break
		}
		if values[i] > result {
			result = values[i]
		}
	}

	return result
}

// Lowest 返回给定数目的最低值。
//
//	Args:
//		values 要计算的数组
//		lenght 计算的长度
func Lowest(values []float64, lenght int) float64 {
	var result = values[len(values)-1]
	for i := len(values) - 1; i >= len(values)-lenght; i-- {
		if i < 0 {
			break
		}
		if values[i] < result {
			result = values[i]
		}
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

	// val = values[index]
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

// Mean 返回float64值数组的平均值
func Mean(values []float64) float64 {
	var total float64 = 0
	for x := range values {
		total += values[x]
	}
	return total / float64(len(values))
}

// trueRange 返回高低闭合的真实范围
func TrueRange(inHigh, inLow, inClose []float64) []float64 {
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

// Variance 返回给定时间段的方差
func Variance(inReal []float64, inTimePeriod int) []float64 {
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

// stdDev - 标准差
func StdDev(inReal []float64, inTimePeriod int, inNbDev float64) []float64 {
	outReal := Variance(inReal, inTimePeriod)

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

// CrossOver 判断 values 最后的值是否向上交叉 source
//
//	当前K线上values的值大于source，并且前一根K线小于source，返回 true，否则返回 false
//
//	Arg:
//		values	要检查的数组
//		source	要判断的值
func CrossOver(values []float64, source float64) bool {
	return values[len(values)-1] > source && values[len(values)-2] < source
}

// CrossUnder 判断 values 最后的值是否向下交叉 source
//
//	当前K线上values的值小于source，并且前一根K线大于source，返回 true，否则返回 false
//
//	Arg:
//		values	要检查的数组
//		source	要判断的值
func CrossUnder(values []float64, source float64) bool {
	return values[len(values)-1] < source && values[len(values)-2] > source
}

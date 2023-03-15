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
		if math.IsNaN(val) {
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

func GetTestKline() Klines {
	return []Kline{
		Kline{Amount: 0, Count: 0, Open: 19193.700000, Close: 19004.500000, Low: 18600.000000, High: 19995.000000, Volume: 14700.175002, Time: time.UnixMicro(1663776000000000)},
		Kline{Amount: 0, Count: 0, Open: 19004.500000, Close: 18462.000000, Low: 18144.500000, High: 19068.500000, Volume: 5412.460913, Time: time.UnixMicro(1663790400000000)},
		Kline{Amount: 0, Count: 0, Open: 18461.100000, Close: 18688.800000, Low: 18360.000000, High: 18787.000000, Volume: 2244.625410, Time: time.UnixMicro(1663804800000000)},
		Kline{Amount: 0, Count: 0, Open: 18689.200000, Close: 18933.500000, Low: 18620.300000, High: 18968.600000, Volume: 3041.719755, Time: time.UnixMicro(1663819200000000)},
		Kline{Amount: 0, Count: 0, Open: 18933.500000, Close: 19231.100000, Low: 18904.100000, High: 19282.200000, Volume: 2480.367947, Time: time.UnixMicro(1663833600000000)},
		Kline{Amount: 0, Count: 0, Open: 19232.300000, Close: 18986.200000, Low: 18798.300000, High: 19321.800000, Volume: 2877.679940, Time: time.UnixMicro(1663848000000000)},
		Kline{Amount: 0, Count: 0, Open: 18987.000000, Close: 19307.700000, Low: 18781.100000, High: 19491.400000, Volume: 1768.031474, Time: time.UnixMicro(1663862400000000)},
		Kline{Amount: 0, Count: 0, Open: 19308.400000, Close: 19399.400000, Low: 19161.200000, High: 19526.800000, Volume: 1493.340755, Time: time.UnixMicro(1663876800000000)},
		Kline{Amount: 0, Count: 0, Open: 19399.600000, Close: 19391.500000, Low: 19227.400000, High: 19466.200000, Volume: 843.218430, Time: time.UnixMicro(1663891200000000)},
		Kline{Amount: 0, Count: 0, Open: 19391.600000, Close: 19224.000000, Low: 19202.500000, High: 19496.000000, Volume: 1372.686588, Time: time.UnixMicro(1663905600000000)},
		Kline{Amount: 0, Count: 0, Open: 19224.000000, Close: 18860.000000, Low: 18837.900000, High: 19268.000000, Volume: 1713.984983, Time: time.UnixMicro(1663920000000000)},
		Kline{Amount: 0, Count: 0, Open: 18859.900000, Close: 18716.300000, Low: 18533.400000, High: 19079.500000, Volume: 2404.982478, Time: time.UnixMicro(1663934400000000)},
		Kline{Amount: 0, Count: 0, Open: 18714.900000, Close: 18812.500000, Low: 18571.100000, High: 18849.300000, Volume: 1368.216583, Time: time.UnixMicro(1663948800000000)},
		Kline{Amount: 0, Count: 0, Open: 18812.500000, Close: 19289.100000, Low: 18800.900000, High: 19404.300000, Volume: 1826.775536, Time: time.UnixMicro(1663963200000000)},
		Kline{Amount: 0, Count: 0, Open: 19289.100000, Close: 19122.500000, Low: 19020.500000, High: 19308.400000, Volume: 993.543529, Time: time.UnixMicro(1663977600000000)},
		Kline{Amount: 0, Count: 0, Open: 19122.600000, Close: 18995.200000, Low: 18926.500000, High: 19186.600000, Volume: 813.994491, Time: time.UnixMicro(1663992000000000)},
		Kline{Amount: 0, Count: 0, Open: 18995.000000, Close: 19057.000000, Low: 18943.000000, High: 19114.500000, Volume: 668.120321, Time: time.UnixMicro(1664006400000000)},
		Kline{Amount: 0, Count: 0, Open: 19055.300000, Close: 19121.800000, Low: 19015.400000, High: 19228.900000, Volume: 767.535762, Time: time.UnixMicro(1664020800000000)},
		Kline{Amount: 0, Count: 0, Open: 19122.700000, Close: 19105.900000, Low: 19056.000000, High: 19192.700000, Volume: 430.112413, Time: time.UnixMicro(1664035200000000)},
		Kline{Amount: 0, Count: 0, Open: 19105.700000, Close: 18921.600000, Low: 18810.200000, High: 19150.800000, Volume: 715.042848, Time: time.UnixMicro(1664049600000000)},
		Kline{Amount: 0, Count: 0, Open: 18922.200000, Close: 19052.200000, Low: 18915.400000, High: 19059.800000, Volume: 547.856820, Time: time.UnixMicro(1664064000000000)},
		Kline{Amount: 0, Count: 0, Open: 19052.300000, Close: 19099.300000, Low: 18942.200000, High: 19132.300000, Volume: 679.342013, Time: time.UnixMicro(1664078400000000)},
		Kline{Amount: 0, Count: 0, Open: 19099.300000, Close: 19118.300000, Low: 18980.900000, High: 19188.000000, Volume: 485.577039, Time: time.UnixMicro(1664092800000000)},
		Kline{Amount: 0, Count: 0, Open: 19118.300000, Close: 18990.200000, Low: 18846.400000, High: 19144.000000, Volume: 1048.177322, Time: time.UnixMicro(1664107200000000)},
		Kline{Amount: 0, Count: 0, Open: 18990.100000, Close: 18916.500000, Low: 18866.300000, High: 19127.600000, Volume: 901.846720, Time: time.UnixMicro(1664121600000000)},
		Kline{Amount: 0, Count: 0, Open: 18916.300000, Close: 18809.200000, Low: 18635.000000, High: 18950.800000, Volume: 1004.584468, Time: time.UnixMicro(1664136000000000)},
		Kline{Amount: 0, Count: 0, Open: 18809.200000, Close: 18871.000000, Low: 18727.400000, High: 18946.200000, Volume: 1297.420333, Time: time.UnixMicro(1664150400000000)},
		Kline{Amount: 0, Count: 0, Open: 18869.700000, Close: 18911.200000, Low: 18687.600000, High: 18932.600000, Volume: 993.288604, Time: time.UnixMicro(1664164800000000)},
		Kline{Amount: 0, Count: 0, Open: 18911.200000, Close: 18875.200000, Low: 18805.100000, High: 19318.900000, Volume: 2426.968516, Time: time.UnixMicro(1664179200000000)},
		Kline{Amount: 0, Count: 0, Open: 18876.300000, Close: 19114.300000, Low: 18823.200000, High: 19276.400000, Volume: 1821.406396, Time: time.UnixMicro(1664193600000000)},
		Kline{Amount: 0, Count: 0, Open: 19114.400000, Close: 19203.900000, Low: 18971.700000, High: 19266.100000, Volume: 874.252442, Time: time.UnixMicro(1664208000000000)},
		Kline{Amount: 0, Count: 0, Open: 19202.500000, Close: 19228.400000, Low: 19073.800000, High: 19244.000000, Volume: 580.843239, Time: time.UnixMicro(1664222400000000)},
		Kline{Amount: 0, Count: 0, Open: 19228.400000, Close: 20101.400000, Low: 19190.900000, High: 20298.600000, Volume: 5539.753458, Time: time.UnixMicro(1664236800000000)},
		Kline{Amount: 0, Count: 0, Open: 20101.500000, Close: 20179.800000, Low: 19958.900000, High: 20345.300000, Volume: 2002.657728, Time: time.UnixMicro(1664251200000000)},
		Kline{Amount: 0, Count: 0, Open: 20179.700000, Close: 20230.300000, Low: 20096.100000, High: 20333.500000, Volume: 1185.306650, Time: time.UnixMicro(1664265600000000)},
		Kline{Amount: 0, Count: 0, Open: 20229.500000, Close: 19886.600000, Low: 19871.300000, High: 20380.300000, Volume: 2056.153879, Time: time.UnixMicro(1664280000000000)},
		Kline{Amount: 0, Count: 0, Open: 19887.700000, Close: 19080.900000, Low: 18820.000000, High: 19950.000000, Volume: 4573.678088, Time: time.UnixMicro(1664294400000000)},
		Kline{Amount: 0, Count: 0, Open: 19080.100000, Close: 19080.000000, Low: 18930.000000, High: 19168.400000, Volume: 1233.651180, Time: time.UnixMicro(1664308800000000)},
		Kline{Amount: 0, Count: 0, Open: 19080.100000, Close: 18745.600000, Low: 18461.900000, High: 19235.700000, Volume: 3002.183035, Time: time.UnixMicro(1664323200000000)},
		Kline{Amount: 0, Count: 0, Open: 18745.500000, Close: 18783.300000, Low: 18624.900000, High: 18881.700000, Volume: 1051.116858, Time: time.UnixMicro(1664337600000000)},
		Kline{Amount: 0, Count: 0, Open: 18783.200000, Close: 18954.500000, Low: 18566.700000, High: 19150.000000, Volume: 3051.645783, Time: time.UnixMicro(1664352000000000)},
		Kline{Amount: 0, Count: 0, Open: 18954.600000, Close: 19488.900000, Low: 18947.100000, High: 19498.000000, Volume: 2473.153404, Time: time.UnixMicro(1664366400000000)},
		Kline{Amount: 0, Count: 0, Open: 19488.900000, Close: 19544.000000, Low: 19261.400000, High: 19665.400000, Volume: 2637.829807, Time: time.UnixMicro(1664380800000000)},
		Kline{Amount: 0, Count: 0, Open: 19544.400000, Close: 19414.100000, Low: 19388.700000, High: 19777.000000, Volume: 1285.110606, Time: time.UnixMicro(1664395200000000)},
		Kline{Amount: 0, Count: 0, Open: 19414.100000, Close: 19509.200000, Low: 19282.000000, High: 19604.400000, Volume: 1241.916267, Time: time.UnixMicro(1664409600000000)},
		Kline{Amount: 0, Count: 0, Open: 19509.100000, Close: 19361.400000, Low: 19270.000000, High: 19626.700000, Volume: 1467.068458, Time: time.UnixMicro(1664424000000000)},
		Kline{Amount: 0, Count: 0, Open: 19361.300000, Close: 19466.200000, Low: 19300.000000, High: 19547.400000, Volume: 1116.957912, Time: time.UnixMicro(1664438400000000)},
		Kline{Amount: 0, Count: 0, Open: 19465.500000, Close: 19291.100000, Low: 18842.400000, High: 19475.000000, Volume: 3747.120426, Time: time.UnixMicro(1664452800000000)},
		Kline{Amount: 0, Count: 0, Open: 19291.000000, Close: 19404.200000, Low: 19122.000000, High: 19647.800000, Volume: 2042.654595, Time: time.UnixMicro(1664467200000000)},
		Kline{Amount: 0, Count: 0, Open: 19405.000000, Close: 19589.400000, Low: 19380.100000, High: 19615.600000, Volume: 629.114507, Time: time.UnixMicro(1664481600000000)},
		Kline{Amount: 0, Count: 0, Open: 19589.600000, Close: 19399.500000, Low: 19314.000000, High: 19699.000000, Volume: 1190.674498, Time: time.UnixMicro(1664496000000000)},
		Kline{Amount: 0, Count: 0, Open: 19398.700000, Close: 19568.800000, Low: 19390.600000, High: 19633.700000, Volume: 1252.705561, Time: time.UnixMicro(1664510400000000)},
		Kline{Amount: 0, Count: 0, Open: 19568.700000, Close: 19435.200000, Low: 19403.600000, High: 19630.000000, Volume: 832.224350, Time: time.UnixMicro(1664524800000000)},
		Kline{Amount: 0, Count: 0, Open: 19435.100000, Close: 19750.600000, Low: 19155.400000, High: 20176.800000, Volume: 6023.910184, Time: time.UnixMicro(1664539200000000)},
		Kline{Amount: 0, Count: 0, Open: 19751.700000, Close: 19477.600000, Low: 19467.800000, High: 19869.800000, Volume: 1944.841445, Time: time.UnixMicro(1664553600000000)},
		Kline{Amount: 0, Count: 0, Open: 19477.100000, Close: 19421.700000, Low: 19245.500000, High: 19539.200000, Volume: 1166.136573, Time: time.UnixMicro(1664568000000000)},
		Kline{Amount: 0, Count: 0, Open: 19422.600000, Close: 19418.700000, Low: 19353.800000, High: 19481.600000, Volume: 508.899184, Time: time.UnixMicro(1664582400000000)},
		Kline{Amount: 0, Count: 0, Open: 19418.700000, Close: 19311.900000, Low: 19258.300000, High: 19422.000000, Volume: 683.383719, Time: time.UnixMicro(1664596800000000)},
		Kline{Amount: 0, Count: 0, Open: 19311.500000, Close: 19310.900000, Low: 19160.100000, High: 19368.300000, Volume: 720.348612, Time: time.UnixMicro(1664611200000000)},
		Kline{Amount: 0, Count: 0, Open: 19311.200000, Close: 19333.500000, Low: 19244.800000, High: 19399.300000, Volume: 585.470904, Time: time.UnixMicro(1664625600000000)},
		Kline{Amount: 0, Count: 0, Open: 19333.100000, Close: 19268.600000, Low: 19207.500000, High: 19352.900000, Volume: 434.628335, Time: time.UnixMicro(1664640000000000)},
		Kline{Amount: 0, Count: 0, Open: 19268.600000, Close: 19311.300000, Low: 19201.800000, High: 19339.800000, Volume: 357.235239, Time: time.UnixMicro(1664654400000000)},
		Kline{Amount: 0, Count: 0, Open: 19311.400000, Close: 19309.800000, Low: 19239.300000, High: 19345.400000, Volume: 327.117804, Time: time.UnixMicro(1664668800000000)},
		Kline{Amount: 0, Count: 0, Open: 19309.800000, Close: 19285.700000, Low: 19234.200000, High: 19396.700000, Volume: 587.460169, Time: time.UnixMicro(1664683200000000)},
		Kline{Amount: 0, Count: 0, Open: 19285.800000, Close: 19196.200000, Low: 19041.500000, High: 19298.700000, Volume: 760.333267, Time: time.UnixMicro(1664697600000000)},
		Kline{Amount: 0, Count: 0, Open: 19195.800000, Close: 19143.900000, Low: 19069.200000, High: 19249.200000, Volume: 540.392447, Time: time.UnixMicro(1664712000000000)},
		Kline{Amount: 0, Count: 0, Open: 19144.700000, Close: 19258.500000, Low: 19133.400000, High: 19345.600000, Volume: 550.411887, Time: time.UnixMicro(1664726400000000)},
		Kline{Amount: 0, Count: 0, Open: 19259.000000, Close: 19055.100000, Low: 18920.000000, High: 19298.400000, Volume: 1400.987720, Time: time.UnixMicro(1664740800000000)},
		Kline{Amount: 0, Count: 0, Open: 19055.200000, Close: 19170.200000, Low: 18963.700000, High: 19255.100000, Volume: 710.248987, Time: time.UnixMicro(1664755200000000)},
		Kline{Amount: 0, Count: 0, Open: 19170.100000, Close: 19197.200000, Low: 19057.600000, High: 19310.300000, Volume: 1282.536932, Time: time.UnixMicro(1664769600000000)},
		Kline{Amount: 0, Count: 0, Open: 19197.300000, Close: 19250.500000, Low: 19124.000000, High: 19253.400000, Volume: 435.129845, Time: time.UnixMicro(1664784000000000)},
		Kline{Amount: 0, Count: 0, Open: 19249.600000, Close: 19335.200000, Low: 19135.900000, High: 19492.700000, Volume: 1789.245012, Time: time.UnixMicro(1664798400000000)},
		Kline{Amount: 0, Count: 0, Open: 19335.300000, Close: 19543.200000, Low: 19312.200000, High: 19646.200000, Volume: 2020.160675, Time: time.UnixMicro(1664812800000000)},
		Kline{Amount: 0, Count: 0, Open: 19541.900000, Close: 19629.600000, Low: 19511.200000, High: 19725.000000, Volume: 810.530419, Time: time.UnixMicro(1664827200000000)},
		Kline{Amount: 0, Count: 0, Open: 19628.600000, Close: 19597.000000, Low: 19492.900000, High: 19721.700000, Volume: 856.812015, Time: time.UnixMicro(1664841600000000)},
		Kline{Amount: 0, Count: 0, Open: 19597.100000, Close: 19967.900000, Low: 19558.100000, High: 19974.500000, Volume: 4962.463383, Time: time.UnixMicro(1664856000000000)},
		Kline{Amount: 0, Count: 0, Open: 19966.300000, Close: 19936.300000, Low: 19814.300000, High: 20167.800000, Volume: 2650.151003, Time: time.UnixMicro(1664870400000000)},
		Kline{Amount: 0, Count: 0, Open: 19936.400000, Close: 20097.800000, Low: 19907.000000, High: 20270.000000, Volume: 4831.201542, Time: time.UnixMicro(1664884800000000)},
		Kline{Amount: 0, Count: 0, Open: 20098.000000, Close: 20217.800000, Low: 19873.600000, High: 20268.000000, Volume: 1462.772171, Time: time.UnixMicro(1664899200000000)},
		Kline{Amount: 0, Count: 0, Open: 20218.300000, Close: 20335.100000, Low: 20149.600000, High: 20462.400000, Volume: 2836.349808, Time: time.UnixMicro(1664913600000000)},
		Kline{Amount: 0, Count: 0, Open: 20336.600000, Close: 20189.600000, Low: 20150.100000, High: 20361.300000, Volume: 786.559439, Time: time.UnixMicro(1664928000000000)},
		Kline{Amount: 0, Count: 0, Open: 20189.600000, Close: 20224.000000, Low: 20074.500000, High: 20297.300000, Volume: 872.580852, Time: time.UnixMicro(1664942400000000)},
		Kline{Amount: 0, Count: 0, Open: 20224.100000, Close: 20025.900000, Low: 19959.200000, High: 20240.600000, Volume: 1117.087967, Time: time.UnixMicro(1664956800000000)},
		Kline{Amount: 0, Count: 0, Open: 20027.100000, Close: 19999.100000, Low: 19731.200000, High: 20068.200000, Volume: 1974.721185, Time: time.UnixMicro(1664971200000000)},
		Kline{Amount: 0, Count: 0, Open: 20000.100000, Close: 20118.400000, Low: 19952.900000, High: 20346.700000, Volume: 1556.534339, Time: time.UnixMicro(1664985600000000)},
		Kline{Amount: 0, Count: 0, Open: 20119.900000, Close: 20158.600000, Low: 19958.600000, High: 20212.900000, Volume: 676.902474, Time: time.UnixMicro(1665000000000000)},
		Kline{Amount: 0, Count: 0, Open: 20158.600000, Close: 20358.600000, Low: 20146.600000, High: 20439.000000, Volume: 1833.259426, Time: time.UnixMicro(1665014400000000)},
		Kline{Amount: 0, Count: 0, Open: 20358.600000, Close: 20223.300000, Low: 20162.200000, High: 20447.900000, Volume: 1206.242685, Time: time.UnixMicro(1665028800000000)},
		Kline{Amount: 0, Count: 0, Open: 20222.400000, Close: 20243.000000, Low: 20086.300000, High: 20258.100000, Volume: 805.711294, Time: time.UnixMicro(1665043200000000)},
		Kline{Amount: 0, Count: 0, Open: 20243.100000, Close: 20016.100000, Low: 19881.000000, High: 20321.700000, Volume: 1937.053942, Time: time.UnixMicro(1665057600000000)},
		Kline{Amount: 0, Count: 0, Open: 20014.100000, Close: 20010.100000, Low: 19953.200000, High: 20256.500000, Volume: 3253.079265, Time: time.UnixMicro(1665072000000000)},
		Kline{Amount: 0, Count: 0, Open: 20010.100000, Close: 19961.100000, Low: 19855.700000, High: 20070.200000, Volume: 594.524712, Time: time.UnixMicro(1665086400000000)},
		Kline{Amount: 0, Count: 0, Open: 19962.100000, Close: 19984.500000, Low: 19945.300000, High: 20052.800000, Volume: 385.163183, Time: time.UnixMicro(1665100800000000)},
		Kline{Amount: 0, Count: 0, Open: 19984.400000, Close: 19928.200000, Low: 19811.100000, High: 20010.700000, Volume: 744.375470, Time: time.UnixMicro(1665115200000000)},
		Kline{Amount: 0, Count: 0, Open: 19928.200000, Close: 20000.100000, Low: 19920.300000, High: 20050.100000, Volume: 1620.015171, Time: time.UnixMicro(1665129600000000)},
		Kline{Amount: 0, Count: 0, Open: 20000.100000, Close: 19579.000000, Low: 19473.600000, High: 20050.000000, Volume: 5296.773411, Time: time.UnixMicro(1665144000000000)},
		Kline{Amount: 0, Count: 0, Open: 19579.000000, Close: 19462.500000, Low: 19320.500000, High: 19607.400000, Volume: 2920.228321, Time: time.UnixMicro(1665158400000000)},
		Kline{Amount: 0, Count: 0, Open: 19463.000000, Close: 19529.600000, Low: 19446.900000, High: 19640.400000, Volume: 1580.766122, Time: time.UnixMicro(1665172800000000)},
		Kline{Amount: 0, Count: 0, Open: 19529.600000, Close: 19503.200000, Low: 19495.400000, High: 19626.200000, Volume: 760.562322, Time: time.UnixMicro(1665187200000000)},
		Kline{Amount: 0, Count: 0, Open: 19503.300000, Close: 19498.900000, Low: 19427.100000, High: 19544.600000, Volume: 538.918488, Time: time.UnixMicro(1665201600000000)},
		Kline{Amount: 0, Count: 0, Open: 19499.000000, Close: 19536.700000, Low: 19452.800000, High: 19544.400000, Volume: 574.069598, Time: time.UnixMicro(1665216000000000)},
		Kline{Amount: 0, Count: 0, Open: 19536.600000, Close: 19488.800000, Low: 19416.400000, High: 19538.600000, Volume: 350.010787, Time: time.UnixMicro(1665230400000000)},
		Kline{Amount: 0, Count: 0, Open: 19488.900000, Close: 19469.000000, Low: 19439.300000, High: 19531.300000, Volume: 226.350888, Time: time.UnixMicro(1665244800000000)},
		Kline{Amount: 0, Count: 0, Open: 19469.100000, Close: 19413.200000, Low: 19215.900000, High: 19481.000000, Volume: 965.257761, Time: time.UnixMicro(1665259200000000)},
		Kline{Amount: 0, Count: 0, Open: 19413.300000, Close: 19369.300000, Low: 19332.100000, High: 19449.500000, Volume: 309.347973, Time: time.UnixMicro(1665273600000000)},
		Kline{Amount: 0, Count: 0, Open: 19368.300000, Close: 19407.700000, Low: 19323.000000, High: 19435.500000, Volume: 364.535269, Time: time.UnixMicro(1665288000000000)},
		Kline{Amount: 0, Count: 0, Open: 19407.700000, Close: 19468.000000, Low: 19375.000000, High: 19549.500000, Volume: 541.414789, Time: time.UnixMicro(1665302400000000)},
		Kline{Amount: 0, Count: 0, Open: 19468.000000, Close: 19536.200000, Low: 19463.800000, High: 19547.400000, Volume: 606.447587, Time: time.UnixMicro(1665316800000000)},
		Kline{Amount: 0, Count: 0, Open: 19536.200000, Close: 19467.100000, Low: 19340.000000, High: 19558.900000, Volume: 1870.601165, Time: time.UnixMicro(1665331200000000)},
		Kline{Amount: 0, Count: 0, Open: 19467.200000, Close: 19441.300000, Low: 19369.800000, High: 19496.400000, Volume: 503.178741, Time: time.UnixMicro(1665345600000000)},
		Kline{Amount: 0, Count: 0, Open: 19441.300000, Close: 19451.500000, Low: 19405.300000, High: 19522.700000, Volume: 550.517157, Time: time.UnixMicro(1665360000000000)},
		Kline{Amount: 0, Count: 0, Open: 19451.500000, Close: 19410.300000, Low: 19371.400000, High: 19480.100000, Volume: 529.123096, Time: time.UnixMicro(1665374400000000)},
		Kline{Amount: 0, Count: 0, Open: 19410.300000, Close: 19333.200000, Low: 19127.300000, High: 19416.400000, Volume: 1502.313042, Time: time.UnixMicro(1665388800000000)},
		Kline{Amount: 0, Count: 0, Open: 19333.400000, Close: 19271.400000, Low: 19225.100000, High: 19437.500000, Volume: 1201.174534, Time: time.UnixMicro(1665403200000000)},
		Kline{Amount: 0, Count: 0, Open: 19270.100000, Close: 19212.100000, Low: 19103.000000, High: 19356.600000, Volume: 1230.159020, Time: time.UnixMicro(1665417600000000)},
		Kline{Amount: 0, Count: 0, Open: 19212.300000, Close: 19129.100000, Low: 19020.000000, High: 19278.700000, Volume: 1170.321959, Time: time.UnixMicro(1665432000000000)},
		Kline{Amount: 0, Count: 0, Open: 19130.100000, Close: 19043.000000, Low: 18954.000000, High: 19130.100000, Volume: 1507.606619, Time: time.UnixMicro(1665446400000000)},
		Kline{Amount: 0, Count: 0, Open: 19043.400000, Close: 19083.100000, Low: 19010.700000, High: 19115.000000, Volume: 628.315084, Time: time.UnixMicro(1665460800000000)},
		Kline{Amount: 0, Count: 0, Open: 19083.000000, Close: 19150.500000, Low: 19032.000000, High: 19159.900000, Volume: 649.633456, Time: time.UnixMicro(1665475200000000)},
		Kline{Amount: 0, Count: 0, Open: 19150.500000, Close: 19145.200000, Low: 18849.500000, High: 19263.900000, Volume: 6631.531799, Time: time.UnixMicro(1665489600000000)},
		Kline{Amount: 0, Count: 0, Open: 19146.900000, Close: 18996.200000, Low: 18921.700000, High: 19181.300000, Volume: 694.687155, Time: time.UnixMicro(1665504000000000)},
		Kline{Amount: 0, Count: 0, Open: 18996.100000, Close: 19061.100000, Low: 18989.300000, High: 19062.900000, Volume: 324.615290, Time: time.UnixMicro(1665518400000000)},
		Kline{Amount: 0, Count: 0, Open: 19060.500000, Close: 19048.800000, Low: 19026.100000, High: 19127.900000, Volume: 457.589689, Time: time.UnixMicro(1665532800000000)},
		Kline{Amount: 0, Count: 0, Open: 19048.800000, Close: 19114.700000, Low: 19044.900000, High: 19198.600000, Volume: 664.186650, Time: time.UnixMicro(1665547200000000)},
		Kline{Amount: 0, Count: 0, Open: 19114.600000, Close: 19117.000000, Low: 19092.600000, High: 19199.900000, Volume: 458.018513, Time: time.UnixMicro(1665561600000000)},
		Kline{Amount: 0, Count: 0, Open: 19117.100000, Close: 19140.100000, Low: 18957.100000, High: 19183.100000, Volume: 2072.786340, Time: time.UnixMicro(1665576000000000)},
		Kline{Amount: 0, Count: 0, Open: 19140.100000, Close: 19137.300000, Low: 19041.200000, High: 19173.500000, Volume: 565.031325, Time: time.UnixMicro(1665590400000000)},
		Kline{Amount: 0, Count: 0, Open: 19137.300000, Close: 19155.200000, Low: 19133.200000, High: 19233.700000, Volume: 366.883112, Time: time.UnixMicro(1665604800000000)},
		Kline{Amount: 0, Count: 0, Open: 19155.200000, Close: 19089.000000, Low: 19019.000000, High: 19172.700000, Volume: 636.330764, Time: time.UnixMicro(1665619200000000)},
		Kline{Amount: 0, Count: 0, Open: 19089.100000, Close: 19009.200000, Low: 18952.500000, High: 19113.000000, Volume: 897.683850, Time: time.UnixMicro(1665633600000000)},
		Kline{Amount: 0, Count: 0, Open: 19009.200000, Close: 18753.400000, Low: 18563.100000, High: 19039.500000, Volume: 3037.908144, Time: time.UnixMicro(1665648000000000)},
		Kline{Amount: 0, Count: 0, Open: 18753.400000, Close: 18955.100000, Low: 18155.000000, High: 19082.500000, Volume: 15032.477726, Time: time.UnixMicro(1665662400000000)},
		Kline{Amount: 0, Count: 0, Open: 18957.800000, Close: 19370.700000, Low: 18904.400000, High: 19507.000000, Volume: 3742.458485, Time: time.UnixMicro(1665676800000000)},
		Kline{Amount: 0, Count: 0, Open: 19370.600000, Close: 19377.800000, Low: 19327.500000, High: 19496.000000, Volume: 928.413697, Time: time.UnixMicro(1665691200000000)},
		Kline{Amount: 0, Count: 0, Open: 19377.800000, Close: 19812.700000, Low: 19339.200000, High: 19954.800000, Volume: 2700.817100, Time: time.UnixMicro(1665705600000000)},
		Kline{Amount: 0, Count: 0, Open: 19812.800000, Close: 19642.100000, Low: 19570.000000, High: 19845.600000, Volume: 1171.842873, Time: time.UnixMicro(1665720000000000)},
		Kline{Amount: 0, Count: 0, Open: 19641.600000, Close: 19597.800000, Low: 19530.000000, High: 19742.500000, Volume: 1303.880419, Time: time.UnixMicro(1665734400000000)},
		Kline{Amount: 0, Count: 0, Open: 19597.900000, Close: 19357.600000, Low: 19308.200000, High: 19849.900000, Volume: 2393.769182, Time: time.UnixMicro(1665748800000000)},
		Kline{Amount: 0, Count: 0, Open: 19357.100000, Close: 19167.600000, Low: 19111.400000, High: 19403.900000, Volume: 1341.406345, Time: time.UnixMicro(1665763200000000)},
		Kline{Amount: 0, Count: 0, Open: 19168.500000, Close: 19178.800000, Low: 19075.800000, High: 19213.100000, Volume: 513.988061, Time: time.UnixMicro(1665777600000000)},
		Kline{Amount: 0, Count: 0, Open: 19178.800000, Close: 19204.500000, Low: 19136.400000, High: 19227.400000, Volume: 547.182479, Time: time.UnixMicro(1665792000000000)},
		Kline{Amount: 0, Count: 0, Open: 19203.900000, Close: 19168.100000, Low: 19138.300000, High: 19209.700000, Volume: 248.490567, Time: time.UnixMicro(1665806400000000)},
		Kline{Amount: 0, Count: 0, Open: 19168.200000, Close: 19165.900000, Low: 19035.500000, High: 19179.900000, Volume: 858.921368, Time: time.UnixMicro(1665820800000000)},
		Kline{Amount: 0, Count: 0, Open: 19165.800000, Close: 19129.300000, Low: 19101.000000, High: 19196.400000, Volume: 449.415895, Time: time.UnixMicro(1665835200000000)},
		Kline{Amount: 0, Count: 0, Open: 19129.300000, Close: 19098.300000, Low: 19056.700000, High: 19162.700000, Volume: 597.323189, Time: time.UnixMicro(1665849600000000)},
		Kline{Amount: 0, Count: 0, Open: 19098.400000, Close: 19070.600000, Low: 18975.700000, High: 19130.500000, Volume: 507.802831, Time: time.UnixMicro(1665864000000000)},
		Kline{Amount: 0, Count: 0, Open: 19070.500000, Close: 19121.700000, Low: 19065.700000, High: 19174.900000, Volume: 263.180618, Time: time.UnixMicro(1665878400000000)},
		Kline{Amount: 0, Count: 0, Open: 19121.700000, Close: 19136.100000, Low: 19113.300000, High: 19160.600000, Volume: 218.461256, Time: time.UnixMicro(1665892800000000)},
		Kline{Amount: 0, Count: 0, Open: 19136.100000, Close: 19145.400000, Low: 19090.000000, High: 19193.300000, Volume: 334.388044, Time: time.UnixMicro(1665907200000000)},
		Kline{Amount: 0, Count: 0, Open: 19145.400000, Close: 19158.000000, Low: 19100.900000, High: 19170.000000, Volume: 399.443184, Time: time.UnixMicro(1665921600000000)},
		Kline{Amount: 0, Count: 0, Open: 19158.100000, Close: 19377.000000, Low: 19121.200000, High: 19380.900000, Volume: 902.094595, Time: time.UnixMicro(1665936000000000)},
		Kline{Amount: 0, Count: 0, Open: 19377.000000, Close: 19263.400000, Low: 19105.000000, High: 19421.700000, Volume: 1944.596676, Time: time.UnixMicro(1665950400000000)},
		Kline{Amount: 0, Count: 0, Open: 19263.400000, Close: 19188.500000, Low: 19160.200000, High: 19305.100000, Volume: 486.963135, Time: time.UnixMicro(1665964800000000)},
		Kline{Amount: 0, Count: 0, Open: 19188.600000, Close: 19295.400000, Low: 19165.000000, High: 19316.400000, Volume: 401.988358, Time: time.UnixMicro(1665979200000000)},
		Kline{Amount: 0, Count: 0, Open: 19295.400000, Close: 19458.800000, Low: 19244.900000, High: 19543.800000, Volume: 1440.103553, Time: time.UnixMicro(1665993600000000)},
		Kline{Amount: 0, Count: 0, Open: 19458.800000, Close: 19532.700000, Low: 19421.900000, High: 19675.200000, Volume: 3248.066976, Time: time.UnixMicro(1666008000000000)},
		Kline{Amount: 0, Count: 0, Open: 19532.600000, Close: 19526.000000, Low: 19441.200000, High: 19589.400000, Volume: 783.838639, Time: time.UnixMicro(1666022400000000)},
		Kline{Amount: 0, Count: 0, Open: 19525.900000, Close: 19547.600000, Low: 19471.500000, High: 19621.800000, Volume: 560.296996, Time: time.UnixMicro(1666036800000000)},
		Kline{Amount: 0, Count: 0, Open: 19547.700000, Close: 19547.400000, Low: 19481.300000, High: 19620.100000, Volume: 506.173678, Time: time.UnixMicro(1666051200000000)},
		Kline{Amount: 0, Count: 0, Open: 19547.500000, Close: 19667.400000, Low: 19534.600000, High: 19709.200000, Volume: 2155.151423, Time: time.UnixMicro(1666065600000000)},
		Kline{Amount: 0, Count: 0, Open: 19667.400000, Close: 19563.500000, Low: 19476.200000, High: 19671.400000, Volume: 913.787746, Time: time.UnixMicro(1666080000000000)},
		Kline{Amount: 0, Count: 0, Open: 19564.500000, Close: 19372.100000, Low: 19278.000000, High: 19677.500000, Volume: 2555.545423, Time: time.UnixMicro(1666094400000000)},
		Kline{Amount: 0, Count: 0, Open: 19372.100000, Close: 19217.000000, Low: 19090.100000, High: 19438.000000, Volume: 3342.645019, Time: time.UnixMicro(1666108800000000)},
		Kline{Amount: 0, Count: 0, Open: 19217.000000, Close: 19330.100000, Low: 19198.200000, High: 19427.500000, Volume: 572.835455, Time: time.UnixMicro(1666123200000000)},
		Kline{Amount: 0, Count: 0, Open: 19330.000000, Close: 19285.700000, Low: 19216.000000, High: 19358.600000, Volume: 467.995744, Time: time.UnixMicro(1666137600000000)},
		Kline{Amount: 0, Count: 0, Open: 19285.700000, Close: 19219.600000, Low: 19149.600000, High: 19304.400000, Volume: 443.410209, Time: time.UnixMicro(1666152000000000)},
		Kline{Amount: 0, Count: 0, Open: 19219.600000, Close: 19201.500000, Low: 19147.100000, High: 19277.100000, Volume: 482.231301, Time: time.UnixMicro(1666166400000000)},
		Kline{Amount: 0, Count: 0, Open: 19201.500000, Close: 19267.300000, Low: 19069.100000, High: 19302.300000, Volume: 1078.693254, Time: time.UnixMicro(1666180800000000)},
		Kline{Amount: 0, Count: 0, Open: 19267.200000, Close: 19228.900000, Low: 19096.100000, High: 19276.800000, Volume: 850.685844, Time: time.UnixMicro(1666195200000000)},
		Kline{Amount: 0, Count: 0, Open: 19228.800000, Close: 19126.900000, Low: 19088.300000, High: 19246.400000, Volume: 543.375981, Time: time.UnixMicro(1666209600000000)},
		Kline{Amount: 0, Count: 0, Open: 19126.600000, Close: 19041.300000, Low: 18902.100000, High: 19161.100000, Volume: 1153.408526, Time: time.UnixMicro(1666224000000000)},
		Kline{Amount: 0, Count: 0, Open: 19043.300000, Close: 19144.100000, Low: 19042.000000, High: 19206.900000, Volume: 680.264321, Time: time.UnixMicro(1666238400000000)},
		Kline{Amount: 0, Count: 0, Open: 19144.100000, Close: 19209.300000, Low: 19108.000000, High: 19252.300000, Volume: 676.235283, Time: time.UnixMicro(1666252800000000)},
		Kline{Amount: 0, Count: 0, Open: 19209.100000, Close: 19255.900000, Low: 19131.200000, High: 19346.200000, Volume: 1026.796662, Time: time.UnixMicro(1666267200000000)},
		Kline{Amount: 0, Count: 0, Open: 19256.000000, Close: 19063.900000, Low: 18967.800000, High: 19256.000000, Volume: 1005.730481, Time: time.UnixMicro(1666281600000000)},
		Kline{Amount: 0, Count: 0, Open: 19063.300000, Close: 19043.900000, Low: 18932.800000, High: 19096.700000, Volume: 828.457655, Time: time.UnixMicro(1666296000000000)},
		Kline{Amount: 0, Count: 0, Open: 19043.900000, Close: 19050.900000, Low: 19001.700000, High: 19132.500000, Volume: 518.485967, Time: time.UnixMicro(1666310400000000)},
		Kline{Amount: 0, Count: 0, Open: 19051.000000, Close: 19030.200000, Low: 18990.000000, High: 19081.300000, Volume: 468.048181, Time: time.UnixMicro(1666324800000000)},
		Kline{Amount: 0, Count: 0, Open: 19030.300000, Close: 18944.300000, Low: 18900.300000, High: 19040.000000, Volume: 792.380985, Time: time.UnixMicro(1666339200000000)},
		Kline{Amount: 0, Count: 0, Open: 18944.300000, Close: 19165.000000, Low: 18648.900000, High: 19168.800000, Volume: 2523.458944, Time: time.UnixMicro(1666353600000000)},
		Kline{Amount: 0, Count: 0, Open: 19165.100000, Close: 19189.600000, Low: 19093.300000, High: 19250.000000, Volume: 889.608527, Time: time.UnixMicro(1666368000000000)},
		Kline{Amount: 0, Count: 0, Open: 19189.500000, Close: 19161.900000, Low: 19139.200000, High: 19248.400000, Volume: 355.494661, Time: time.UnixMicro(1666382400000000)},
		Kline{Amount: 0, Count: 0, Open: 19162.100000, Close: 19152.100000, Low: 19115.000000, High: 19181.700000, Volume: 308.026773, Time: time.UnixMicro(1666396800000000)},
		Kline{Amount: 0, Count: 0, Open: 19152.100000, Close: 19163.100000, Low: 19135.000000, High: 19173.900000, Volume: 357.246693, Time: time.UnixMicro(1666411200000000)},
		Kline{Amount: 0, Count: 0, Open: 19163.500000, Close: 19180.000000, Low: 19140.900000, High: 19185.700000, Volume: 209.846787, Time: time.UnixMicro(1666425600000000)},
		Kline{Amount: 0, Count: 0, Open: 19180.100000, Close: 19235.300000, Low: 19157.700000, High: 19252.000000, Volume: 454.779440, Time: time.UnixMicro(1666440000000000)},
		Kline{Amount: 0, Count: 0, Open: 19235.300000, Close: 19184.800000, Low: 19137.000000, High: 19247.300000, Volume: 201.238876, Time: time.UnixMicro(1666454400000000)},
		Kline{Amount: 0, Count: 0, Open: 19184.900000, Close: 19202.700000, Low: 19175.800000, High: 19222.300000, Volume: 113.112977, Time: time.UnixMicro(1666468800000000)},
		Kline{Amount: 0, Count: 0, Open: 19202.800000, Close: 19197.100000, Low: 19167.600000, High: 19217.900000, Volume: 170.900489, Time: time.UnixMicro(1666483200000000)},
		Kline{Amount: 0, Count: 0, Open: 19197.100000, Close: 19163.300000, Low: 19151.000000, High: 19218.900000, Volume: 209.479647, Time: time.UnixMicro(1666497600000000)},
		Kline{Amount: 0, Count: 0, Open: 19163.200000, Close: 19156.200000, Low: 19130.700000, High: 19187.300000, Volume: 208.021636, Time: time.UnixMicro(1666512000000000)},
		Kline{Amount: 0, Count: 0, Open: 19156.200000, Close: 19193.000000, Low: 19070.100000, High: 19209.300000, Volume: 471.992813, Time: time.UnixMicro(1666526400000000)},
		Kline{Amount: 0, Count: 0, Open: 19193.100000, Close: 19499.800000, Low: 19171.400000, High: 19527.400000, Volume: 1378.809753, Time: time.UnixMicro(1666540800000000)},
		Kline{Amount: 0, Count: 0, Open: 19499.800000, Close: 19570.500000, Low: 19447.100000, High: 19691.900000, Volume: 1203.846178, Time: time.UnixMicro(1666555200000000)},
		Kline{Amount: 0, Count: 0, Open: 19570.700000, Close: 19391.200000, Low: 19341.300000, High: 19598.800000, Volume: 1124.612564, Time: time.UnixMicro(1666569600000000)},
		Kline{Amount: 0, Count: 0, Open: 19391.200000, Close: 19321.900000, Low: 19250.300000, High: 19391.200000, Volume: 888.546083, Time: time.UnixMicro(1666584000000000)},
		Kline{Amount: 0, Count: 0, Open: 19322.000000, Close: 19423.100000, Low: 19293.800000, High: 19444.200000, Volume: 459.585669, Time: time.UnixMicro(1666598400000000)},
		Kline{Amount: 0, Count: 0, Open: 19423.200000, Close: 19303.600000, Low: 19161.000000, High: 19452.000000, Volume: 2520.936859, Time: time.UnixMicro(1666612800000000)},
		Kline{Amount: 0, Count: 0, Open: 19303.700000, Close: 19342.000000, Low: 19224.200000, High: 19422.000000, Volume: 487.760595, Time: time.UnixMicro(1666627200000000)},
		Kline{Amount: 0, Count: 0, Open: 19342.100000, Close: 19332.300000, Low: 19306.600000, High: 19415.200000, Volume: 282.017738, Time: time.UnixMicro(1666641600000000)},
		Kline{Amount: 0, Count: 0, Open: 19331.400000, Close: 19341.900000, Low: 19246.800000, High: 19361.800000, Volume: 371.812323, Time: time.UnixMicro(1666656000000000)},
		Kline{Amount: 0, Count: 0, Open: 19342.000000, Close: 19308.500000, Low: 19274.100000, High: 19374.000000, Volume: 262.389397, Time: time.UnixMicro(1666670400000000)},
		Kline{Amount: 0, Count: 0, Open: 19307.600000, Close: 19299.900000, Low: 19242.600000, High: 19323.800000, Volume: 304.219431, Time: time.UnixMicro(1666684800000000)},
		Kline{Amount: 0, Count: 0, Open: 19300.000000, Close: 19747.500000, Low: 19264.000000, High: 19807.900000, Volume: 2127.983023, Time: time.UnixMicro(1666699200000000)},
		Kline{Amount: 0, Count: 0, Open: 19747.600000, Close: 20271.400000, Low: 19721.100000, High: 20412.300000, Volume: 7316.699989, Time: time.UnixMicro(1666713600000000)},
		Kline{Amount: 0, Count: 0, Open: 20271.800000, Close: 20082.500000, Low: 19983.600000, High: 20293.200000, Volume: 1935.774257, Time: time.UnixMicro(1666728000000000)},
		Kline{Amount: 0, Count: 0, Open: 20082.500000, Close: 20251.600000, Low: 20053.200000, High: 20307.500000, Volume: 1260.922046, Time: time.UnixMicro(1666742400000000)},
		Kline{Amount: 0, Count: 0, Open: 20251.600000, Close: 20345.500000, Low: 20136.100000, High: 20377.600000, Volume: 1515.907170, Time: time.UnixMicro(1666756800000000)},
		Kline{Amount: 0, Count: 0, Open: 20345.400000, Close: 20609.400000, Low: 20283.500000, High: 20790.200000, Volume: 4242.780621, Time: time.UnixMicro(1666771200000000)},
		Kline{Amount: 0, Count: 0, Open: 20609.300000, Close: 20845.000000, Low: 20344.200000, High: 21020.500000, Volume: 4824.220769, Time: time.UnixMicro(1666785600000000)},
		Kline{Amount: 0, Count: 0, Open: 20845.100000, Close: 20761.000000, Low: 20639.500000, High: 20989.900000, Volume: 2665.284469, Time: time.UnixMicro(1666800000000000)},
		Kline{Amount: 0, Count: 0, Open: 20761.500000, Close: 20770.800000, Low: 20694.400000, High: 20906.000000, Volume: 1630.863048, Time: time.UnixMicro(1666814400000000)},
		Kline{Amount: 0, Count: 0, Open: 20769.200000, Close: 20752.400000, Low: 20590.900000, High: 20847.500000, Volume: 960.287549, Time: time.UnixMicro(1666828800000000)},
		Kline{Amount: 0, Count: 0, Open: 20752.400000, Close: 20713.200000, Low: 20682.200000, High: 20870.300000, Volume: 616.242816, Time: time.UnixMicro(1666843200000000)},
		Kline{Amount: 0, Count: 0, Open: 20713.100000, Close: 20623.300000, Low: 20457.600000, High: 20750.000000, Volume: 1285.171310, Time: time.UnixMicro(1666857600000000)},
		Kline{Amount: 0, Count: 0, Open: 20623.200000, Close: 20562.200000, Low: 20435.300000, High: 20768.700000, Volume: 5149.982500, Time: time.UnixMicro(1666872000000000)},
		Kline{Amount: 0, Count: 0, Open: 20562.100000, Close: 20628.700000, Low: 20495.100000, High: 20675.000000, Volume: 709.612769, Time: time.UnixMicro(1666886400000000)},
		Kline{Amount: 0, Count: 0, Open: 20629.300000, Close: 20290.400000, Low: 20200.000000, High: 20633.300000, Volume: 2943.995720, Time: time.UnixMicro(1666900800000000)},
		Kline{Amount: 0, Count: 0, Open: 20291.500000, Close: 20295.900000, Low: 20164.000000, High: 20334.300000, Volume: 1499.510333, Time: time.UnixMicro(1666915200000000)},
		Kline{Amount: 0, Count: 0, Open: 20295.900000, Close: 20078.900000, Low: 20049.900000, High: 20306.200000, Volume: 3728.774840, Time: time.UnixMicro(1666929600000000)},
		Kline{Amount: 0, Count: 0, Open: 20079.000000, Close: 20176.100000, Low: 20027.600000, High: 20246.200000, Volume: 971.073586, Time: time.UnixMicro(1666944000000000)},
		Kline{Amount: 0, Count: 0, Open: 20175.300000, Close: 20481.500000, Low: 20000.000000, High: 20587.100000, Volume: 3170.418750, Time: time.UnixMicro(1666958400000000)},
		Kline{Amount: 0, Count: 0, Open: 20482.200000, Close: 20609.000000, Low: 20450.800000, High: 20746.800000, Volume: 1445.614426, Time: time.UnixMicro(1666972800000000)},
		Kline{Amount: 0, Count: 0, Open: 20610.000000, Close: 20592.500000, Low: 20569.000000, High: 20735.300000, Volume: 400.146006, Time: time.UnixMicro(1666987200000000)},
		Kline{Amount: 0, Count: 0, Open: 20592.500000, Close: 20752.800000, Low: 20556.900000, High: 20800.000000, Volume: 1277.871191, Time: time.UnixMicro(1667001600000000)},
		Kline{Amount: 0, Count: 0, Open: 20752.900000, Close: 20772.400000, Low: 20618.700000, High: 20777.700000, Volume: 759.072096, Time: time.UnixMicro(1667016000000000)},
		Kline{Amount: 0, Count: 0, Open: 20772.000000, Close: 20697.200000, Low: 20653.900000, High: 21080.000000, Volume: 3895.348517, Time: time.UnixMicro(1667030400000000)},
		Kline{Amount: 0, Count: 0, Open: 20697.200000, Close: 20923.100000, Low: 20679.100000, High: 20943.100000, Volume: 2865.002634, Time: time.UnixMicro(1667044800000000)},
		Kline{Amount: 0, Count: 0, Open: 20923.100000, Close: 20838.200000, Low: 20760.600000, High: 20979.500000, Volume: 1953.904766, Time: time.UnixMicro(1667059200000000)},
		Kline{Amount: 0, Count: 0, Open: 20838.200000, Close: 20808.100000, Low: 20705.100000, High: 20899.000000, Volume: 776.123779, Time: time.UnixMicro(1667073600000000)},
		Kline{Amount: 0, Count: 0, Open: 20808.200000, Close: 20773.000000, Low: 20702.600000, High: 20835.100000, Volume: 641.278816, Time: time.UnixMicro(1667088000000000)},
		Kline{Amount: 0, Count: 0, Open: 20773.100000, Close: 20850.600000, Low: 20741.300000, High: 20927.700000, Volume: 1860.407062, Time: time.UnixMicro(1667102400000000)},
		Kline{Amount: 0, Count: 0, Open: 20850.700000, Close: 20774.000000, Low: 20634.300000, High: 20885.200000, Volume: 1045.179807, Time: time.UnixMicro(1667116800000000)},
		Kline{Amount: 0, Count: 0, Open: 20773.800000, Close: 20679.000000, Low: 20608.800000, High: 20777.800000, Volume: 849.653032, Time: time.UnixMicro(1667131200000000)},
		Kline{Amount: 0, Count: 0, Open: 20679.100000, Close: 20644.900000, Low: 20536.100000, High: 20696.500000, Volume: 1057.320167, Time: time.UnixMicro(1667145600000000)},
		Kline{Amount: 0, Count: 0, Open: 20645.000000, Close: 20626.700000, Low: 20509.700000, High: 20737.400000, Volume: 444.287511, Time: time.UnixMicro(1667160000000000)},
		Kline{Amount: 0, Count: 0, Open: 20627.000000, Close: 20520.600000, Low: 20430.000000, High: 20673.200000, Volume: 700.258152, Time: time.UnixMicro(1667174400000000)},
		Kline{Amount: 0, Count: 0, Open: 20520.600000, Close: 20547.300000, Low: 20457.200000, High: 20590.200000, Volume: 399.596702, Time: time.UnixMicro(1667188800000000)},
		Kline{Amount: 0, Count: 0, Open: 20547.400000, Close: 20733.100000, Low: 20452.800000, High: 20846.000000, Volume: 970.104313, Time: time.UnixMicro(1667203200000000)},
		Kline{Amount: 0, Count: 0, Open: 20733.200000, Close: 20400.000000, Low: 20241.100000, High: 20752.200000, Volume: 2439.557291, Time: time.UnixMicro(1667217600000000)},
		Kline{Amount: 0, Count: 0, Open: 20400.000000, Close: 20369.900000, Low: 20331.000000, High: 20527.400000, Volume: 851.480687, Time: time.UnixMicro(1667232000000000)},
		Kline{Amount: 0, Count: 0, Open: 20370.600000, Close: 20491.000000, Low: 20356.800000, High: 20520.500000, Volume: 364.634775, Time: time.UnixMicro(1667246400000000)},
		Kline{Amount: 0, Count: 0, Open: 20491.100000, Close: 20483.800000, Low: 20439.300000, High: 20585.200000, Volume: 475.252761, Time: time.UnixMicro(1667260800000000)},
		Kline{Amount: 0, Count: 0, Open: 20483.800000, Close: 20598.600000, Low: 20466.400000, High: 20661.300000, Volume: 600.576383, Time: time.UnixMicro(1667275200000000)},
		Kline{Amount: 0, Count: 0, Open: 20598.600000, Close: 20525.200000, Low: 20492.300000, High: 20694.600000, Volume: 1092.423997, Time: time.UnixMicro(1667289600000000)},
		Kline{Amount: 0, Count: 0, Open: 20525.700000, Close: 20447.500000, Low: 20332.200000, High: 20569.300000, Volume: 1051.511789, Time: time.UnixMicro(1667304000000000)},
		Kline{Amount: 0, Count: 0, Open: 20447.400000, Close: 20450.000000, Low: 20375.800000, High: 20515.900000, Volume: 534.523657, Time: time.UnixMicro(1667318400000000)},
		Kline{Amount: 0, Count: 0, Open: 20450.000000, Close: 20483.200000, Low: 20433.300000, High: 20500.300000, Volume: 262.014572, Time: time.UnixMicro(1667332800000000)},
		Kline{Amount: 0, Count: 0, Open: 20484.000000, Close: 20520.400000, Low: 20406.600000, High: 20555.500000, Volume: 378.723540, Time: time.UnixMicro(1667347200000000)},
		Kline{Amount: 0, Count: 0, Open: 20521.000000, Close: 20426.300000, Low: 20357.900000, High: 20559.700000, Volume: 595.936439, Time: time.UnixMicro(1667361600000000)},
		Kline{Amount: 0, Count: 0, Open: 20426.400000, Close: 20434.200000, Low: 20334.600000, High: 20536.700000, Volume: 658.885718, Time: time.UnixMicro(1667376000000000)},
		Kline{Amount: 0, Count: 0, Open: 20434.200000, Close: 20407.600000, Low: 20361.600000, High: 20455.000000, Volume: 590.254266, Time: time.UnixMicro(1667390400000000)},
		Kline{Amount: 0, Count: 0, Open: 20407.500000, Close: 20253.800000, Low: 20100.000000, High: 20800.000000, Volume: 5653.562103, Time: time.UnixMicro(1667404800000000)},
		Kline{Amount: 0, Count: 0, Open: 20252.800000, Close: 20151.900000, Low: 20056.500000, High: 20282.000000, Volume: 1043.037970, Time: time.UnixMicro(1667419200000000)},
		Kline{Amount: 0, Count: 0, Open: 20152.000000, Close: 20310.600000, Low: 20124.800000, High: 20337.400000, Volume: 975.385743, Time: time.UnixMicro(1667433600000000)},
		Kline{Amount: 0, Count: 0, Open: 20310.900000, Close: 20302.300000, Low: 20279.500000, High: 20393.500000, Volume: 597.726186, Time: time.UnixMicro(1667448000000000)},
		Kline{Amount: 0, Count: 0, Open: 20302.300000, Close: 20130.400000, Low: 20055.000000, High: 20337.300000, Volume: 804.704558, Time: time.UnixMicro(1667462400000000)},
		Kline{Amount: 0, Count: 0, Open: 20130.200000, Close: 20252.400000, Low: 20032.200000, High: 20334.500000, Volume: 1173.279596, Time: time.UnixMicro(1667476800000000)},
		Kline{Amount: 0, Count: 0, Open: 20252.400000, Close: 20261.000000, Low: 20181.800000, High: 20313.600000, Volume: 340.670399, Time: time.UnixMicro(1667491200000000)},
		Kline{Amount: 0, Count: 0, Open: 20261.000000, Close: 20206.300000, Low: 20162.800000, High: 20267.800000, Volume: 366.070080, Time: time.UnixMicro(1667505600000000)},
		Kline{Amount: 0, Count: 0, Open: 20206.200000, Close: 20324.400000, Low: 20184.100000, High: 20328.000000, Volume: 377.205291, Time: time.UnixMicro(1667520000000000)},
		Kline{Amount: 0, Count: 0, Open: 20324.300000, Close: 20609.900000, Low: 20309.400000, High: 20686.400000, Volume: 1371.312420, Time: time.UnixMicro(1667534400000000)},
		Kline{Amount: 0, Count: 0, Open: 20609.900000, Close: 20555.200000, Low: 20547.800000, High: 20665.000000, Volume: 891.537735, Time: time.UnixMicro(1667548800000000)},
		Kline{Amount: 0, Count: 0, Open: 20555.200000, Close: 20836.200000, Low: 20367.600000, High: 21300.000000, Volume: 6172.497424, Time: time.UnixMicro(1667563200000000)},
		Kline{Amount: 0, Count: 0, Open: 20836.300000, Close: 21106.900000, Low: 20675.200000, High: 21106.900000, Volume: 1605.953619, Time: time.UnixMicro(1667577600000000)},
		Kline{Amount: 0, Count: 0, Open: 21106.900000, Close: 21145.100000, Low: 21055.100000, High: 21205.000000, Volume: 583.322754, Time: time.UnixMicro(1667592000000000)},
		Kline{Amount: 0, Count: 0, Open: 21145.200000, Close: 21447.700000, Low: 21084.200000, High: 21475.900000, Volume: 1444.578684, Time: time.UnixMicro(1667606400000000)},
		Kline{Amount: 0, Count: 0, Open: 21447.700000, Close: 21384.600000, Low: 21365.000000, High: 21464.800000, Volume: 485.016230, Time: time.UnixMicro(1667620800000000)},
		Kline{Amount: 0, Count: 0, Open: 21384.600000, Close: 21402.900000, Low: 21238.900000, High: 21422.900000, Volume: 708.703780, Time: time.UnixMicro(1667635200000000)},
		Kline{Amount: 0, Count: 0, Open: 21402.900000, Close: 21301.900000, Low: 21211.200000, High: 21417.600000, Volume: 625.835716, Time: time.UnixMicro(1667649600000000)},
		Kline{Amount: 0, Count: 0, Open: 21302.800000, Close: 21333.200000, Low: 21263.300000, High: 21344.500000, Volume: 318.112918, Time: time.UnixMicro(1667664000000000)},
		Kline{Amount: 0, Count: 0, Open: 21333.200000, Close: 21300.500000, Low: 21238.700000, High: 21379.000000, Volume: 398.040994, Time: time.UnixMicro(1667678400000000)},
		Kline{Amount: 0, Count: 0, Open: 21300.400000, Close: 21233.800000, Low: 21173.800000, High: 21362.600000, Volume: 642.084419, Time: time.UnixMicro(1667692800000000)},
		Kline{Amount: 0, Count: 0, Open: 21233.800000, Close: 21194.200000, Low: 21148.700000, High: 21250.400000, Volume: 318.752207, Time: time.UnixMicro(1667707200000000)},
		Kline{Amount: 0, Count: 0, Open: 21194.200000, Close: 21256.700000, Low: 21178.900000, High: 21309.900000, Volume: 840.546431, Time: time.UnixMicro(1667721600000000)},
		Kline{Amount: 0, Count: 0, Open: 21256.700000, Close: 21244.200000, Low: 21191.000000, High: 21289.000000, Volume: 263.830830, Time: time.UnixMicro(1667736000000000)},
		Kline{Amount: 0, Count: 0, Open: 21244.300000, Close: 21197.800000, Low: 21154.500000, High: 21280.000000, Volume: 624.712931, Time: time.UnixMicro(1667750400000000)},
		Kline{Amount: 0, Count: 0, Open: 21197.900000, Close: 20908.900000, Low: 20880.000000, High: 21212.000000, Volume: 1082.070764, Time: time.UnixMicro(1667764800000000)},
		Kline{Amount: 0, Count: 0, Open: 20908.900000, Close: 20913.900000, Low: 20821.300000, High: 21067.300000, Volume: 839.352440, Time: time.UnixMicro(1667779200000000)},
		Kline{Amount: 0, Count: 0, Open: 20913.900000, Close: 20632.200000, Low: 20615.200000, High: 20919.700000, Volume: 1528.336238, Time: time.UnixMicro(1667793600000000)},
		Kline{Amount: 0, Count: 0, Open: 20632.100000, Close: 20731.600000, Low: 20548.200000, High: 20808.000000, Volume: 1065.185910, Time: time.UnixMicro(1667808000000000)},
		Kline{Amount: 0, Count: 0, Open: 20731.600000, Close: 20767.600000, Low: 20625.800000, High: 20799.300000, Volume: 703.330041, Time: time.UnixMicro(1667822400000000)},
		Kline{Amount: 0, Count: 0, Open: 20767.500000, Close: 20846.500000, Low: 20656.100000, High: 20897.000000, Volume: 564.132616, Time: time.UnixMicro(1667836800000000)},
		Kline{Amount: 0, Count: 0, Open: 20846.600000, Close: 20587.300000, Low: 20385.900000, High: 20888.400000, Volume: 1034.674810, Time: time.UnixMicro(1667851200000000)},
		Kline{Amount: 0, Count: 0, Open: 20587.600000, Close: 20151.200000, Low: 20122.000000, High: 20668.200000, Volume: 1622.514138, Time: time.UnixMicro(1667865600000000)},
		Kline{Amount: 0, Count: 0, Open: 20151.300000, Close: 19798.500000, Low: 19338.400000, High: 20241.500000, Volume: 5498.271334, Time: time.UnixMicro(1667880000000000)},
		Kline{Amount: 0, Count: 0, Open: 19798.400000, Close: 19701.400000, Low: 19586.500000, High: 19799.400000, Volume: 891.335100, Time: time.UnixMicro(1667894400000000)},
		Kline{Amount: 0, Count: 0, Open: 19701.400000, Close: 19526.000000, Low: 19220.000000, High: 19745.700000, Volume: 4962.958047, Time: time.UnixMicro(1667908800000000)},
		Kline{Amount: 0, Count: 0, Open: 19526.000000, Close: 18255.900000, Low: 16888.000000, High: 20690.800000, Volume: 42346.918779, Time: time.UnixMicro(1667923200000000)},
		Kline{Amount: 0, Count: 0, Open: 18255.900000, Close: 18545.000000, Low: 17800.000000, High: 18721.600000, Volume: 6630.088346, Time: time.UnixMicro(1667937600000000)},
		Kline{Amount: 0, Count: 0, Open: 18545.100000, Close: 18310.400000, Low: 18000.000000, High: 18585.700000, Volume: 2967.417601, Time: time.UnixMicro(1667952000000000)},
		Kline{Amount: 0, Count: 0, Open: 18311.200000, Close: 18236.900000, Low: 18170.900000, High: 18506.600000, Volume: 2248.983090, Time: time.UnixMicro(1667966400000000)},
		Kline{Amount: 0, Count: 0, Open: 18236.900000, Close: 17815.200000, Low: 17260.000000, High: 18253.200000, Volume: 8051.076276, Time: time.UnixMicro(1667980800000000)},
		Kline{Amount: 0, Count: 0, Open: 17816.200000, Close: 17124.000000, Low: 16920.000000, High: 18025.000000, Volume: 9476.377376, Time: time.UnixMicro(1667995200000000)},
		Kline{Amount: 0, Count: 0, Open: 17123.900000, Close: 16779.500000, Low: 16431.200000, High: 17327.900000, Volume: 5486.162540, Time: time.UnixMicro(1668009600000000)},
		Kline{Amount: 0, Count: 0, Open: 16779.600000, Close: 15921.000000, Low: 15511.100000, High: 16900.000000, Volume: 15188.292272, Time: time.UnixMicro(1668024000000000)},
		Kline{Amount: 0, Count: 0, Open: 15920.900000, Close: 16372.100000, Low: 15753.500000, High: 16577.000000, Volume: 4414.484737, Time: time.UnixMicro(1668038400000000)},
		Kline{Amount: 0, Count: 0, Open: 16372.000000, Close: 16760.800000, Low: 16207.600000, High: 16922.700000, Volume: 4268.204840, Time: time.UnixMicro(1668052800000000)},
		Kline{Amount: 0, Count: 0, Open: 16760.800000, Close: 16402.300000, Low: 16269.700000, High: 16963.600000, Volume: 3595.513306, Time: time.UnixMicro(1668067200000000)},
		Kline{Amount: 0, Count: 0, Open: 16403.800000, Close: 17553.300000, Low: 16392.400000, High: 17965.800000, Volume: 10248.153289, Time: time.UnixMicro(1668081600000000)},
		Kline{Amount: 0, Count: 0, Open: 17557.200000, Close: 17457.000000, Low: 17125.600000, High: 17993.500000, Volume: 3273.236589, Time: time.UnixMicro(1668096000000000)},
		Kline{Amount: 0, Count: 0, Open: 17456.900000, Close: 17601.100000, Low: 17373.100000, High: 18193.400000, Volume: 3480.604159, Time: time.UnixMicro(1668110400000000)},
		Kline{Amount: 0, Count: 0, Open: 17601.600000, Close: 17062.700000, Low: 16883.200000, High: 17696.400000, Volume: 2907.350354, Time: time.UnixMicro(1668124800000000)},
		Kline{Amount: 0, Count: 0, Open: 17064.000000, Close: 17411.100000, Low: 17003.300000, High: 17499.000000, Volume: 2905.943035, Time: time.UnixMicro(1668139200000000)},
		Kline{Amount: 0, Count: 0, Open: 17411.000000, Close: 17351.100000, Low: 17226.600000, High: 17544.200000, Volume: 1891.073895, Time: time.UnixMicro(1668153600000000)},
		Kline{Amount: 0, Count: 0, Open: 17351.400000, Close: 16873.500000, Low: 16366.000000, High: 17544.400000, Volume: 6199.117967, Time: time.UnixMicro(1668168000000000)},
		Kline{Amount: 0, Count: 0, Open: 16873.600000, Close: 16855.400000, Low: 16540.000000, High: 17042.100000, Volume: 1752.178549, Time: time.UnixMicro(1668182400000000)},
		Kline{Amount: 0, Count: 0, Open: 16851.600000, Close: 17063.800000, Low: 16560.000000, High: 17144.000000, Volume: 1464.707595, Time: time.UnixMicro(1668196800000000)},
		Kline{Amount: 0, Count: 0, Open: 17070.000000, Close: 16907.700000, Low: 16779.500000, High: 17120.500000, Volume: 915.210741, Time: time.UnixMicro(1668211200000000)},
		Kline{Amount: 0, Count: 0, Open: 16909.300000, Close: 16823.100000, Low: 16631.200000, High: 16937.200000, Volume: 1491.017641, Time: time.UnixMicro(1668225600000000)},
		Kline{Amount: 0, Count: 0, Open: 16823.100000, Close: 16871.500000, Low: 16723.000000, High: 16974.500000, Volume: 885.737960, Time: time.UnixMicro(1668240000000000)},
		Kline{Amount: 0, Count: 0, Open: 16871.800000, Close: 16931.600000, Low: 16758.000000, High: 16995.500000, Volume: 722.884126, Time: time.UnixMicro(1668254400000000)},
		Kline{Amount: 0, Count: 0, Open: 16932.200000, Close: 16897.200000, Low: 16829.200000, High: 16963.200000, Volume: 368.638442, Time: time.UnixMicro(1668268800000000)},
		Kline{Amount: 0, Count: 0, Open: 16897.100000, Close: 16813.100000, Low: 16757.900000, High: 16897.100000, Volume: 411.112302, Time: time.UnixMicro(1668283200000000)},
		Kline{Amount: 0, Count: 0, Open: 16813.200000, Close: 16918.100000, Low: 16786.100000, High: 16955.000000, Volume: 447.609368, Time: time.UnixMicro(1668297600000000)},
		Kline{Amount: 0, Count: 0, Open: 16918.000000, Close: 16774.400000, Low: 16562.300000, High: 16927.000000, Volume: 1035.724386, Time: time.UnixMicro(1668312000000000)},
		Kline{Amount: 0, Count: 0, Open: 16774.500000, Close: 16660.700000, Low: 16437.700000, High: 16793.800000, Volume: 1284.185813, Time: time.UnixMicro(1668326400000000)},
		Kline{Amount: 0, Count: 0, Open: 16659.500000, Close: 16575.700000, Low: 16523.300000, High: 16710.000000, Volume: 542.955823, Time: time.UnixMicro(1668340800000000)},
		Kline{Amount: 0, Count: 0, Open: 16575.800000, Close: 16561.500000, Low: 16467.200000, High: 16621.300000, Volume: 817.484178, Time: time.UnixMicro(1668355200000000)},
		Kline{Amount: 0, Count: 0, Open: 16560.500000, Close: 16336.500000, Low: 16241.400000, High: 16588.700000, Volume: 1027.456248, Time: time.UnixMicro(1668369600000000)},
		Kline{Amount: 0, Count: 0, Open: 16333.800000, Close: 16139.800000, Low: 15915.000000, High: 16410.000000, Volume: 2308.258404, Time: time.UnixMicro(1668384000000000)},
		Kline{Amount: 0, Count: 0, Open: 16139.800000, Close: 16823.200000, Low: 15820.000000, High: 16885.000000, Volume: 3483.376184, Time: time.UnixMicro(1668398400000000)},
		Kline{Amount: 0, Count: 0, Open: 16821.700000, Close: 16753.500000, Low: 16635.400000, High: 16916.000000, Volume: 1622.510019, Time: time.UnixMicro(1668412800000000)},
		Kline{Amount: 0, Count: 0, Open: 16753.400000, Close: 16592.100000, Low: 16425.200000, High: 17198.700000, Volume: 3457.325225, Time: time.UnixMicro(1668427200000000)},
		Kline{Amount: 0, Count: 0, Open: 16592.100000, Close: 16263.000000, Low: 16189.600000, High: 16676.400000, Volume: 1260.784278, Time: time.UnixMicro(1668441600000000)},
		Kline{Amount: 0, Count: 0, Open: 16261.300000, Close: 16623.500000, Low: 16201.000000, High: 16767.100000, Volume: 892.286542, Time: time.UnixMicro(1668456000000000)},
		Kline{Amount: 0, Count: 0, Open: 16621.700000, Close: 16831.000000, Low: 16533.000000, High: 16899.100000, Volume: 1319.565567, Time: time.UnixMicro(1668470400000000)},
		Kline{Amount: 0, Count: 0, Open: 16830.800000, Close: 16848.600000, Low: 16644.000000, High: 16940.100000, Volume: 903.066180, Time: time.UnixMicro(1668484800000000)},
		Kline{Amount: 0, Count: 0, Open: 16848.700000, Close: 16779.200000, Low: 16701.400000, High: 17000.000000, Volume: 1311.912002, Time: time.UnixMicro(1668499200000000)},
		Kline{Amount: 0, Count: 0, Open: 16780.200000, Close: 16984.300000, Low: 16727.900000, High: 17125.500000, Volume: 3119.719834, Time: time.UnixMicro(1668513600000000)},
		Kline{Amount: 0, Count: 0, Open: 16984.200000, Close: 16854.600000, Low: 16624.000000, High: 17068.000000, Volume: 2554.558868, Time: time.UnixMicro(1668528000000000)},
		Kline{Amount: 0, Count: 0, Open: 16858.400000, Close: 16896.800000, Low: 16747.000000, High: 16937.300000, Volume: 705.758683, Time: time.UnixMicro(1668542400000000)},
		Kline{Amount: 0, Count: 0, Open: 16895.600000, Close: 16970.400000, Low: 16766.200000, High: 17000.000000, Volume: 1037.466625, Time: time.UnixMicro(1668556800000000)},
		Kline{Amount: 0, Count: 0, Open: 16970.500000, Close: 16825.200000, Low: 16803.200000, High: 16976.000000, Volume: 1068.821162, Time: time.UnixMicro(1668571200000000)},
		Kline{Amount: 0, Count: 0, Open: 16826.100000, Close: 16707.700000, Low: 16655.000000, High: 16851.500000, Volume: 1156.101905, Time: time.UnixMicro(1668585600000000)},
		Kline{Amount: 0, Count: 0, Open: 16708.100000, Close: 16436.300000, Low: 16381.300000, High: 16722.000000, Volume: 2141.674702, Time: time.UnixMicro(1668600000000000)},
		Kline{Amount: 0, Count: 0, Open: 16438.600000, Close: 16562.700000, Low: 16408.200000, High: 16648.000000, Volume: 1079.448198, Time: time.UnixMicro(1668614400000000)},
		Kline{Amount: 0, Count: 0, Open: 16560.500000, Close: 16660.000000, Low: 16492.500000, High: 16778.000000, Volume: 739.311256, Time: time.UnixMicro(1668628800000000)},
		Kline{Amount: 0, Count: 0, Open: 16660.100000, Close: 16525.900000, Low: 16513.700000, High: 16750.000000, Volume: 588.425876, Time: time.UnixMicro(1668643200000000)},
		Kline{Amount: 0, Count: 0, Open: 16525.800000, Close: 16609.500000, Low: 16420.000000, High: 16622.200000, Volume: 601.141539, Time: time.UnixMicro(1668657600000000)},
		Kline{Amount: 0, Count: 0, Open: 16610.100000, Close: 16597.600000, Low: 16472.700000, High: 16635.800000, Volume: 763.518570, Time: time.UnixMicro(1668672000000000)},
		Kline{Amount: 0, Count: 0, Open: 16597.600000, Close: 16536.600000, Low: 16435.000000, High: 16720.000000, Volume: 1158.881478, Time: time.UnixMicro(1668686400000000)},
		Kline{Amount: 0, Count: 0, Open: 16538.700000, Close: 16620.800000, Low: 16520.700000, High: 16755.000000, Volume: 607.692264, Time: time.UnixMicro(1668700800000000)},
		Kline{Amount: 0, Count: 0, Open: 16620.800000, Close: 16692.400000, Low: 16611.200000, High: 16744.200000, Volume: 426.843995, Time: time.UnixMicro(1668715200000000)},
		Kline{Amount: 0, Count: 0, Open: 16692.500000, Close: 16854.000000, Low: 16679.700000, High: 16988.800000, Volume: 1968.209724, Time: time.UnixMicro(1668729600000000)},
		Kline{Amount: 0, Count: 0, Open: 16853.900000, Close: 16749.700000, Low: 16730.800000, High: 16863.300000, Volume: 1096.512883, Time: time.UnixMicro(1668744000000000)},
		Kline{Amount: 0, Count: 0, Open: 16749.800000, Close: 16750.800000, Low: 16685.700000, High: 16848.900000, Volume: 1986.822820, Time: time.UnixMicro(1668758400000000)},
		Kline{Amount: 0, Count: 0, Open: 16751.600000, Close: 16694.700000, Low: 16610.700000, High: 16818.400000, Volume: 1063.360414, Time: time.UnixMicro(1668772800000000)},
		Kline{Amount: 0, Count: 0, Open: 16695.500000, Close: 16658.100000, Low: 16555.100000, High: 16733.900000, Volume: 909.905261, Time: time.UnixMicro(1668787200000000)},
		Kline{Amount: 0, Count: 0, Open: 16657.600000, Close: 16699.100000, Low: 16621.400000, High: 16729.400000, Volume: 394.860919, Time: time.UnixMicro(1668801600000000)},
		Kline{Amount: 0, Count: 0, Open: 16700.900000, Close: 16629.400000, Low: 16575.400000, High: 16707.500000, Volume: 674.179758, Time: time.UnixMicro(1668816000000000)},
		Kline{Amount: 0, Count: 0, Open: 16629.500000, Close: 16603.500000, Low: 16557.500000, High: 16650.000000, Volume: 530.222049, Time: time.UnixMicro(1668830400000000)},
		Kline{Amount: 0, Count: 0, Open: 16603.500000, Close: 16673.700000, Low: 16585.100000, High: 16678.000000, Volume: 377.130740, Time: time.UnixMicro(1668844800000000)},
		Kline{Amount: 0, Count: 0, Open: 16675.500000, Close: 16658.200000, Low: 16635.000000, High: 16684.000000, Volume: 249.874302, Time: time.UnixMicro(1668859200000000)},
		Kline{Amount: 0, Count: 0, Open: 16657.800000, Close: 16626.500000, Low: 16609.300000, High: 16665.300000, Volume: 238.823435, Time: time.UnixMicro(1668873600000000)},
		Kline{Amount: 0, Count: 0, Open: 16626.600000, Close: 16707.700000, Low: 16626.500000, High: 16818.600000, Volume: 569.724728, Time: time.UnixMicro(1668888000000000)},
		Kline{Amount: 0, Count: 0, Open: 16702.800000, Close: 16682.700000, Low: 16670.000000, High: 16750.000000, Volume: 313.675363, Time: time.UnixMicro(1668902400000000)},
		Kline{Amount: 0, Count: 0, Open: 16682.700000, Close: 16721.900000, Low: 16673.000000, High: 16742.000000, Volume: 224.059077, Time: time.UnixMicro(1668916800000000)},
		Kline{Amount: 0, Count: 0, Open: 16722.800000, Close: 16523.000000, Low: 16461.700000, High: 16734.600000, Volume: 1084.395644, Time: time.UnixMicro(1668931200000000)},
		Kline{Amount: 0, Count: 0, Open: 16522.200000, Close: 16581.800000, Low: 16458.000000, High: 16622.100000, Volume: 617.125217, Time: time.UnixMicro(1668945600000000)},
		Kline{Amount: 0, Count: 0, Open: 16580.900000, Close: 16594.000000, Low: 16531.400000, High: 16610.800000, Volume: 229.078710, Time: time.UnixMicro(1668960000000000)},
		Kline{Amount: 0, Count: 0, Open: 16594.000000, Close: 16282.700000, Low: 16178.000000, High: 16598.000000, Volume: 1140.512764, Time: time.UnixMicro(1668974400000000)},
		Kline{Amount: 0, Count: 0, Open: 16282.100000, Close: 16076.400000, Low: 15901.500000, High: 16293.700000, Volume: 2325.740903, Time: time.UnixMicro(1668988800000000)},
		Kline{Amount: 0, Count: 0, Open: 16075.900000, Close: 16040.200000, Low: 15962.400000, High: 16331.500000, Volume: 1740.337593, Time: time.UnixMicro(1669003200000000)},
		Kline{Amount: 0, Count: 0, Open: 16041.500000, Close: 16087.000000, Low: 15972.900000, High: 16173.300000, Volume: 948.667687, Time: time.UnixMicro(1669017600000000)},
		Kline{Amount: 0, Count: 0, Open: 16087.100000, Close: 16133.800000, Low: 16032.000000, High: 16300.100000, Volume: 2522.916982, Time: time.UnixMicro(1669032000000000)},
		Kline{Amount: 0, Count: 0, Open: 16133.800000, Close: 15728.700000, Low: 15588.900000, High: 16150.100000, Volume: 3697.803913, Time: time.UnixMicro(1669046400000000)},
		Kline{Amount: 0, Count: 0, Open: 15724.300000, Close: 15783.500000, Low: 15450.000000, High: 15892.900000, Volume: 1800.600248, Time: time.UnixMicro(1669060800000000)},
		Kline{Amount: 0, Count: 0, Open: 15783.600000, Close: 15846.700000, Low: 15749.700000, High: 15990.000000, Volume: 1306.850700, Time: time.UnixMicro(1669075200000000)},
		Kline{Amount: 0, Count: 0, Open: 15846.800000, Close: 15763.300000, Low: 15624.200000, High: 15883.200000, Volume: 1508.913022, Time: time.UnixMicro(1669089600000000)},
		Kline{Amount: 0, Count: 0, Open: 15763.900000, Close: 15752.300000, Low: 15665.000000, High: 15824.800000, Volume: 1203.338962, Time: time.UnixMicro(1669104000000000)},
		Kline{Amount: 0, Count: 0, Open: 15751.900000, Close: 16237.400000, Low: 15691.600000, High: 16315.300000, Volume: 2474.328178, Time: time.UnixMicro(1669118400000000)},
		Kline{Amount: 0, Count: 0, Open: 16240.400000, Close: 16168.100000, Low: 16052.500000, High: 16287.400000, Volume: 1000.483898, Time: time.UnixMicro(1669132800000000)},
		Kline{Amount: 0, Count: 0, Open: 16166.800000, Close: 16227.100000, Low: 16086.900000, High: 16268.400000, Volume: 755.652433, Time: time.UnixMicro(1669147200000000)},
		Kline{Amount: 0, Count: 0, Open: 16227.500000, Close: 16586.000000, Low: 16165.400000, High: 16590.000000, Volume: 2689.725605, Time: time.UnixMicro(1669161600000000)},
		Kline{Amount: 0, Count: 0, Open: 16588.800000, Close: 16578.200000, Low: 16440.000000, High: 16661.000000, Volume: 1624.748208, Time: time.UnixMicro(1669176000000000)},
		Kline{Amount: 0, Count: 0, Open: 16579.900000, Close: 16586.000000, Low: 16500.000000, High: 16673.300000, Volume: 1080.768540, Time: time.UnixMicro(1669190400000000)},
		Kline{Amount: 0, Count: 0, Open: 16583.600000, Close: 16426.100000, Low: 16346.700000, High: 16602.100000, Volume: 1131.696068, Time: time.UnixMicro(1669204800000000)},
		Kline{Amount: 0, Count: 0, Open: 16426.200000, Close: 16475.800000, Low: 16325.000000, High: 16720.400000, Volume: 1035.734146, Time: time.UnixMicro(1669219200000000)},
		Kline{Amount: 0, Count: 0, Open: 16472.900000, Close: 16606.800000, Low: 16429.100000, High: 16669.600000, Volume: 477.314278, Time: time.UnixMicro(1669233600000000)},
		Kline{Amount: 0, Count: 0, Open: 16604.900000, Close: 16712.100000, Low: 16536.400000, High: 16810.800000, Volume: 1280.846749, Time: time.UnixMicro(1669248000000000)},
		Kline{Amount: 0, Count: 0, Open: 16711.600000, Close: 16622.200000, Low: 16610.000000, High: 16736.800000, Volume: 995.128393, Time: time.UnixMicro(1669262400000000)},
		Kline{Amount: 0, Count: 0, Open: 16622.600000, Close: 16579.000000, Low: 16497.000000, High: 16635.100000, Volume: 1457.946548, Time: time.UnixMicro(1669276800000000)},
		Kline{Amount: 0, Count: 0, Open: 16578.900000, Close: 16558.100000, Low: 16463.700000, High: 16605.000000, Volume: 560.191892, Time: time.UnixMicro(1669291200000000)},
		Kline{Amount: 0, Count: 0, Open: 16558.100000, Close: 16576.700000, Low: 16520.600000, High: 16674.400000, Volume: 822.803728, Time: time.UnixMicro(1669305600000000)},
		Kline{Amount: 0, Count: 0, Open: 16576.800000, Close: 16602.800000, Low: 16538.300000, High: 16637.600000, Volume: 285.516867, Time: time.UnixMicro(1669320000000000)},
		Kline{Amount: 0, Count: 0, Open: 16602.800000, Close: 16492.100000, Low: 16472.000000, High: 16620.500000, Volume: 443.893125, Time: time.UnixMicro(1669334400000000)},
		Kline{Amount: 0, Count: 0, Open: 16492.100000, Close: 16467.300000, Low: 16350.100000, High: 16502.000000, Volume: 962.822792, Time: time.UnixMicro(1669348800000000)},
		Kline{Amount: 0, Count: 0, Open: 16466.800000, Close: 16532.600000, Low: 16414.200000, High: 16593.600000, Volume: 756.556106, Time: time.UnixMicro(1669363200000000)},
		Kline{Amount: 0, Count: 0, Open: 16534.000000, Close: 16511.100000, Low: 16443.500000, High: 16582.000000, Volume: 1023.350330, Time: time.UnixMicro(1669377600000000)},
		Kline{Amount: 0, Count: 0, Open: 16513.200000, Close: 16525.100000, Low: 16461.900000, High: 16552.400000, Volume: 542.691262, Time: time.UnixMicro(1669392000000000)},
		Kline{Amount: 0, Count: 0, Open: 16525.200000, Close: 16522.400000, Low: 16484.400000, High: 16630.300000, Volume: 380.106774, Time: time.UnixMicro(1669406400000000)},
		Kline{Amount: 0, Count: 0, Open: 16520.200000, Close: 16623.100000, Low: 16512.300000, High: 16700.000000, Volume: 745.492817, Time: time.UnixMicro(1669420800000000)},
		Kline{Amount: 0, Count: 0, Open: 16623.100000, Close: 16587.800000, Low: 16533.600000, High: 16655.200000, Volume: 577.300620, Time: time.UnixMicro(1669435200000000)},
		Kline{Amount: 0, Count: 0, Open: 16587.800000, Close: 16581.800000, Low: 16559.500000, High: 16637.000000, Volume: 469.790189, Time: time.UnixMicro(1669449600000000)},
		Kline{Amount: 0, Count: 0, Open: 16581.800000, Close: 16513.200000, Low: 16481.700000, High: 16681.300000, Volume: 889.344757, Time: time.UnixMicro(1669464000000000)},
		Kline{Amount: 0, Count: 0, Open: 16513.600000, Close: 16488.500000, Low: 16422.500000, High: 16533.300000, Volume: 575.840448, Time: time.UnixMicro(1669478400000000)},
		Kline{Amount: 0, Count: 0, Open: 16488.600000, Close: 16458.400000, Low: 16382.900000, High: 16546.100000, Volume: 353.314450, Time: time.UnixMicro(1669492800000000)},
		Kline{Amount: 0, Count: 0, Open: 16459.700000, Close: 16529.600000, Low: 16451.600000, High: 16570.200000, Volume: 368.203993, Time: time.UnixMicro(1669507200000000)},
		Kline{Amount: 0, Count: 0, Open: 16529.600000, Close: 16563.700000, Low: 16528.300000, High: 16598.300000, Volume: 344.567721, Time: time.UnixMicro(1669521600000000)},
		Kline{Amount: 0, Count: 0, Open: 16563.700000, Close: 16553.500000, Low: 16513.200000, High: 16596.300000, Volume: 508.572979, Time: time.UnixMicro(1669536000000000)},
		Kline{Amount: 0, Count: 0, Open: 16553.400000, Close: 16575.000000, Low: 16520.400000, High: 16584.900000, Volume: 483.139863, Time: time.UnixMicro(1669550400000000)},
		Kline{Amount: 0, Count: 0, Open: 16575.000000, Close: 16555.400000, Low: 16523.100000, High: 16579.000000, Volume: 179.796892, Time: time.UnixMicro(1669564800000000)},
		Kline{Amount: 0, Count: 0, Open: 16555.400000, Close: 16427.400000, Low: 16400.000000, High: 16597.200000, Volume: 559.462550, Time: time.UnixMicro(1669579200000000)},
		Kline{Amount: 0, Count: 0, Open: 16429.800000, Close: 16181.200000, Low: 16055.100000, High: 16484.300000, Volume: 2230.714160, Time: time.UnixMicro(1669593600000000)},
		Kline{Amount: 0, Count: 0, Open: 16181.200000, Close: 16227.700000, Low: 16140.000000, High: 16249.900000, Volume: 427.115696, Time: time.UnixMicro(1669608000000000)},
		Kline{Amount: 0, Count: 0, Open: 16227.300000, Close: 16217.700000, Low: 16172.000000, High: 16270.800000, Volume: 689.580774, Time: time.UnixMicro(1669622400000000)},
		Kline{Amount: 0, Count: 0, Open: 16215.300000, Close: 16145.800000, Low: 16123.500000, High: 16314.200000, Volume: 903.059053, Time: time.UnixMicro(1669636800000000)},
		Kline{Amount: 0, Count: 0, Open: 16146.900000, Close: 16223.400000, Low: 16000.000000, High: 16393.500000, Volume: 1334.834872, Time: time.UnixMicro(1669651200000000)},
		Kline{Amount: 0, Count: 0, Open: 16225.000000, Close: 16215.000000, Low: 16193.200000, High: 16269.600000, Volume: 275.589372, Time: time.UnixMicro(1669665600000000)},
		Kline{Amount: 0, Count: 0, Open: 16215.100000, Close: 16286.900000, Low: 16101.400000, High: 16310.600000, Volume: 861.220285, Time: time.UnixMicro(1669680000000000)},
		Kline{Amount: 0, Count: 0, Open: 16287.000000, Close: 16463.600000, Low: 16282.600000, High: 16537.200000, Volume: 1562.731382, Time: time.UnixMicro(1669694400000000)},
		Kline{Amount: 0, Count: 0, Open: 16463.700000, Close: 16498.100000, Low: 16432.500000, High: 16546.200000, Volume: 1011.617474, Time: time.UnixMicro(1669708800000000)},
		Kline{Amount: 0, Count: 0, Open: 16497.900000, Close: 16392.700000, Low: 16333.800000, High: 16507.000000, Volume: 947.371280, Time: time.UnixMicro(1669723200000000)},
		Kline{Amount: 0, Count: 0, Open: 16392.600000, Close: 16426.100000, Low: 16347.300000, High: 16478.000000, Volume: 726.392665, Time: time.UnixMicro(1669737600000000)},
		Kline{Amount: 0, Count: 0, Open: 16426.100000, Close: 16442.800000, Low: 16424.500000, High: 16525.200000, Volume: 411.178442, Time: time.UnixMicro(1669752000000000)},
		Kline{Amount: 0, Count: 0, Open: 16442.800000, Close: 16847.700000, Low: 16430.700000, High: 17091.600000, Volume: 4118.076522, Time: time.UnixMicro(1669766400000000)},
		Kline{Amount: 0, Count: 0, Open: 16845.700000, Close: 16884.600000, Low: 16822.900000, High: 16918.900000, Volume: 963.320898, Time: time.UnixMicro(1669780800000000)},
		Kline{Amount: 0, Count: 0, Open: 16884.700000, Close: 16875.800000, Low: 16843.600000, High: 16928.400000, Volume: 1755.775472, Time: time.UnixMicro(1669795200000000)},
		Kline{Amount: 0, Count: 0, Open: 16875.800000, Close: 16865.400000, Low: 16768.800000, High: 16891.000000, Volume: 2958.305896, Time: time.UnixMicro(1669809600000000)},
		Kline{Amount: 0, Count: 0, Open: 16865.700000, Close: 17059.800000, Low: 16712.700000, High: 17146.600000, Volume: 1763.021267, Time: time.UnixMicro(1669824000000000)},
		Kline{Amount: 0, Count: 0, Open: 17059.700000, Close: 17164.800000, Low: 17028.800000, High: 17255.000000, Volume: 1398.252809, Time: time.UnixMicro(1669838400000000)},
		Kline{Amount: 0, Count: 0, Open: 17164.700000, Close: 17146.400000, Low: 17089.500000, High: 17230.600000, Volume: 1288.963165, Time: time.UnixMicro(1669852800000000)},
		Kline{Amount: 0, Count: 0, Open: 17146.300000, Close: 17067.000000, Low: 17051.300000, High: 17164.500000, Volume: 953.708594, Time: time.UnixMicro(1669867200000000)},
		Kline{Amount: 0, Count: 0, Open: 17067.000000, Close: 17103.300000, Low: 17042.300000, High: 17133.900000, Volume: 728.907845, Time: time.UnixMicro(1669881600000000)},
		Kline{Amount: 0, Count: 0, Open: 17103.300000, Close: 16978.200000, Low: 16900.800000, High: 17335.500000, Volume: 2303.627512, Time: time.UnixMicro(1669896000000000)},
		Kline{Amount: 0, Count: 0, Open: 16980.600000, Close: 16956.700000, Low: 16905.100000, High: 16998.300000, Volume: 589.807742, Time: time.UnixMicro(1669910400000000)},
		Kline{Amount: 0, Count: 0, Open: 16956.700000, Close: 16980.100000, Low: 16865.200000, High: 16984.800000, Volume: 565.394361, Time: time.UnixMicro(1669924800000000)},
		Kline{Amount: 0, Count: 0, Open: 16979.900000, Close: 16903.200000, Low: 16863.500000, High: 17045.800000, Volume: 793.934198, Time: time.UnixMicro(1669939200000000)},
		Kline{Amount: 0, Count: 0, Open: 16903.200000, Close: 16955.000000, Low: 16885.000000, High: 16974.300000, Volume: 721.212755, Time: time.UnixMicro(1669953600000000)},
		Kline{Amount: 0, Count: 0, Open: 16955.000000, Close: 16993.800000, Low: 16938.800000, High: 17051.700000, Volume: 938.857102, Time: time.UnixMicro(1669968000000000)},
		Kline{Amount: 0, Count: 0, Open: 16993.800000, Close: 16942.000000, Low: 16776.800000, High: 17104.100000, Volume: 3421.761599, Time: time.UnixMicro(1669982400000000)},
		Kline{Amount: 0, Count: 0, Open: 16940.900000, Close: 16989.500000, Low: 16907.000000, High: 16999.100000, Volume: 344.488335, Time: time.UnixMicro(1669996800000000)},
		Kline{Amount: 0, Count: 0, Open: 16989.600000, Close: 17091.300000, Low: 16982.800000, High: 17094.000000, Volume: 1391.291917, Time: time.UnixMicro(1670011200000000)},
		Kline{Amount: 0, Count: 0, Open: 17091.400000, Close: 17019.800000, Low: 17010.000000, High: 17169.200000, Volume: 610.377689, Time: time.UnixMicro(1670025600000000)},
		Kline{Amount: 0, Count: 0, Open: 17019.800000, Close: 16963.500000, Low: 16928.500000, High: 17030.400000, Volume: 522.903205, Time: time.UnixMicro(1670040000000000)},
		Kline{Amount: 0, Count: 0, Open: 16963.500000, Close: 16945.600000, Low: 16891.500000, High: 16984.000000, Volume: 591.802154, Time: time.UnixMicro(1670054400000000)},
		Kline{Amount: 0, Count: 0, Open: 16945.600000, Close: 16962.600000, Low: 16923.700000, High: 16980.900000, Volume: 303.248443, Time: time.UnixMicro(1670068800000000)},
		Kline{Amount: 0, Count: 0, Open: 16962.600000, Close: 16955.000000, Low: 16939.000000, High: 16981.400000, Volume: 191.205065, Time: time.UnixMicro(1670083200000000)},
		Kline{Amount: 0, Count: 0, Open: 16955.000000, Close: 16887.100000, Low: 16865.000000, High: 16957.600000, Volume: 436.953958, Time: time.UnixMicro(1670097600000000)},
		Kline{Amount: 0, Count: 0, Open: 16886.800000, Close: 16962.200000, Low: 16883.200000, High: 17049.200000, Volume: 384.037081, Time: time.UnixMicro(1670112000000000)},
		Kline{Amount: 0, Count: 0, Open: 16962.200000, Close: 16998.900000, Low: 16962.200000, High: 17087.200000, Volume: 536.613479, Time: time.UnixMicro(1670126400000000)},
		Kline{Amount: 0, Count: 0, Open: 16998.300000, Close: 16948.600000, Low: 16906.100000, High: 17048.800000, Volume: 338.702781, Time: time.UnixMicro(1670140800000000)},
		Kline{Amount: 0, Count: 0, Open: 16949.000000, Close: 17030.500000, Low: 16915.100000, High: 17066.300000, Volume: 2034.063670, Time: time.UnixMicro(1670155200000000)},
		Kline{Amount: 0, Count: 0, Open: 17030.600000, Close: 17091.600000, Low: 16972.900000, High: 17155.000000, Volume: 1852.573965, Time: time.UnixMicro(1670169600000000)},
		Kline{Amount: 0, Count: 0, Open: 17091.600000, Close: 17109.000000, Low: 17077.300000, High: 17200.000000, Volume: 1562.641204, Time: time.UnixMicro(1670184000000000)},
		Kline{Amount: 0, Count: 0, Open: 17109.000000, Close: 17194.900000, Low: 17080.900000, High: 17341.600000, Volume: 2235.732511, Time: time.UnixMicro(1670198400000000)},
		Kline{Amount: 0, Count: 0, Open: 17195.000000, Close: 17312.300000, Low: 17192.700000, High: 17420.100000, Volume: 1456.323056, Time: time.UnixMicro(1670212800000000)},
		Kline{Amount: 0, Count: 0, Open: 17312.200000, Close: 17311.700000, Low: 17284.000000, High: 17386.600000, Volume: 1295.492494, Time: time.UnixMicro(1670227200000000)},
		Kline{Amount: 0, Count: 0, Open: 17311.600000, Close: 17084.700000, Low: 17020.700000, High: 17313.300000, Volume: 1872.975327, Time: time.UnixMicro(1670241600000000)},
		Kline{Amount: 0, Count: 0, Open: 17084.700000, Close: 16928.900000, Low: 16870.700000, High: 17128.600000, Volume: 1590.520312, Time: time.UnixMicro(1670256000000000)},
		Kline{Amount: 0, Count: 0, Open: 16926.100000, Close: 16966.600000, Low: 16902.500000, High: 16989.600000, Volume: 310.246502, Time: time.UnixMicro(1670270400000000)},
		Kline{Amount: 0, Count: 0, Open: 16967.400000, Close: 17001.200000, Low: 16967.400000, High: 17100.300000, Volume: 2485.088141, Time: time.UnixMicro(1670284800000000)},
		Kline{Amount: 0, Count: 0, Open: 17001.300000, Close: 17005.100000, Low: 16960.400000, High: 17023.000000, Volume: 275.538301, Time: time.UnixMicro(1670299200000000)},
		Kline{Amount: 0, Count: 0, Open: 17005.100000, Close: 16983.900000, Low: 16910.200000, High: 17041.700000, Volume: 699.557302, Time: time.UnixMicro(1670313600000000)},
		Kline{Amount: 0, Count: 0, Open: 16983.900000, Close: 16984.100000, Low: 16932.000000, High: 17030.400000, Volume: 550.933501, Time: time.UnixMicro(1670328000000000)},
		Kline{Amount: 0, Count: 0, Open: 16984.200000, Close: 16971.400000, Low: 16915.900000, High: 17017.800000, Volume: 478.158443, Time: time.UnixMicro(1670342400000000)},
		Kline{Amount: 0, Count: 0, Open: 16971.500000, Close: 17087.900000, Low: 16956.000000, High: 17110.800000, Volume: 690.298428, Time: time.UnixMicro(1670356800000000)},
		Kline{Amount: 0, Count: 0, Open: 17089.000000, Close: 17027.200000, Low: 17022.400000, High: 17141.400000, Volume: 678.713771, Time: time.UnixMicro(1670371200000000)},
		Kline{Amount: 0, Count: 0, Open: 17027.300000, Close: 16765.900000, Low: 16730.800000, High: 17044.800000, Volume: 1987.342859, Time: time.UnixMicro(1670385600000000)},
		Kline{Amount: 0, Count: 0, Open: 16765.900000, Close: 16796.200000, Low: 16686.000000, High: 16835.300000, Volume: 993.064941, Time: time.UnixMicro(1670400000000000)},
		Kline{Amount: 0, Count: 0, Open: 16795.100000, Close: 16838.200000, Low: 16768.600000, High: 16897.800000, Volume: 736.904312, Time: time.UnixMicro(1670414400000000)},
		Kline{Amount: 0, Count: 0, Open: 16838.300000, Close: 16814.800000, Low: 16782.200000, High: 16848.400000, Volume: 287.115546, Time: time.UnixMicro(1670428800000000)},
		Kline{Amount: 0, Count: 0, Open: 16816.200000, Close: 16837.000000, Low: 16798.000000, High: 16865.900000, Volume: 159.140661, Time: time.UnixMicro(1670443200000000)},
		Kline{Amount: 0, Count: 0, Open: 16837.100000, Close: 16832.400000, Low: 16781.200000, High: 16887.700000, Volume: 279.805949, Time: time.UnixMicro(1670457600000000)},
		Kline{Amount: 0, Count: 0, Open: 16832.400000, Close: 16826.700000, Low: 16788.000000, High: 16846.400000, Volume: 370.453607, Time: time.UnixMicro(1670472000000000)},
		Kline{Amount: 0, Count: 0, Open: 16826.600000, Close: 16849.700000, Low: 16733.300000, High: 16865.200000, Volume: 709.949347, Time: time.UnixMicro(1670486400000000)},
		Kline{Amount: 0, Count: 0, Open: 16850.500000, Close: 16918.600000, Low: 16803.200000, High: 16957.000000, Volume: 1217.774467, Time: time.UnixMicro(1670500800000000)},
		Kline{Amount: 0, Count: 0, Open: 16918.700000, Close: 17252.800000, Low: 16912.400000, High: 17296.900000, Volume: 1899.191079, Time: time.UnixMicro(1670515200000000)},
		Kline{Amount: 0, Count: 0, Open: 17252.900000, Close: 17226.000000, Low: 17159.400000, High: 17293.000000, Volume: 563.260303, Time: time.UnixMicro(1670529600000000)},
		Kline{Amount: 0, Count: 0, Open: 17226.000000, Close: 17204.000000, Low: 17191.400000, High: 17292.800000, Volume: 443.936628, Time: time.UnixMicro(1670544000000000)},
		Kline{Amount: 0, Count: 0, Open: 17204.100000, Close: 17211.800000, Low: 17193.500000, High: 17232.200000, Volume: 385.491467, Time: time.UnixMicro(1670558400000000)},
		Kline{Amount: 0, Count: 0, Open: 17211.800000, Close: 17244.400000, Low: 17200.300000, High: 17266.000000, Volume: 547.931502, Time: time.UnixMicro(1670572800000000)},
		Kline{Amount: 0, Count: 0, Open: 17244.400000, Close: 17179.200000, Low: 17067.200000, High: 17353.900000, Volume: 1993.369447, Time: time.UnixMicro(1670587200000000)},
		Kline{Amount: 0, Count: 0, Open: 17179.300000, Close: 17144.000000, Low: 17123.500000, High: 17189.800000, Volume: 517.511740, Time: time.UnixMicro(1670601600000000)},
		Kline{Amount: 0, Count: 0, Open: 17144.000000, Close: 17127.200000, Low: 17063.800000, High: 17154.500000, Volume: 313.221867, Time: time.UnixMicro(1670616000000000)},
		Kline{Amount: 0, Count: 0, Open: 17127.300000, Close: 17146.300000, Low: 17119.400000, High: 17167.600000, Volume: 206.792413, Time: time.UnixMicro(1670630400000000)},
		Kline{Amount: 0, Count: 0, Open: 17146.300000, Close: 17154.900000, Low: 17131.000000, High: 17174.400000, Volume: 237.002853, Time: time.UnixMicro(1670644800000000)},
		Kline{Amount: 0, Count: 0, Open: 17154.900000, Close: 17168.100000, Low: 17113.700000, High: 17170.800000, Volume: 231.014284, Time: time.UnixMicro(1670659200000000)},
		Kline{Amount: 0, Count: 0, Open: 17168.100000, Close: 17212.300000, Low: 17140.200000, High: 17229.400000, Volume: 560.343965, Time: time.UnixMicro(1670673600000000)},
		Kline{Amount: 0, Count: 0, Open: 17212.300000, Close: 17177.200000, Low: 17161.300000, High: 17230.000000, Volume: 453.088351, Time: time.UnixMicro(1670688000000000)},
		Kline{Amount: 0, Count: 0, Open: 17176.500000, Close: 17127.900000, Low: 17092.000000, High: 17182.300000, Volume: 218.164685, Time: time.UnixMicro(1670702400000000)},
		Kline{Amount: 0, Count: 0, Open: 17128.000000, Close: 17156.700000, Low: 17125.900000, High: 17173.500000, Volume: 215.104606, Time: time.UnixMicro(1670716800000000)},
		Kline{Amount: 0, Count: 0, Open: 17156.700000, Close: 17181.900000, Low: 17154.900000, High: 17199.900000, Volume: 220.249401, Time: time.UnixMicro(1670731200000000)},
		Kline{Amount: 0, Count: 0, Open: 17181.800000, Close: 17167.300000, Low: 17139.400000, High: 17183.500000, Volume: 187.896617, Time: time.UnixMicro(1670745600000000)},
		Kline{Amount: 0, Count: 0, Open: 17167.400000, Close: 17154.900000, Low: 17125.300000, High: 17167.900000, Volume: 488.960117, Time: time.UnixMicro(1670760000000000)},
		Kline{Amount: 0, Count: 0, Open: 17155.000000, Close: 17180.300000, Low: 17140.000000, High: 17269.300000, Volume: 597.406042, Time: time.UnixMicro(1670774400000000)},
		Kline{Amount: 0, Count: 0, Open: 17180.300000, Close: 17082.300000, Low: 17072.100000, High: 17180.300000, Volume: 681.497960, Time: time.UnixMicro(1670788800000000)},
		Kline{Amount: 0, Count: 0, Open: 17083.100000, Close: 16916.600000, Low: 16870.600000, High: 17087.200000, Volume: 1341.026711, Time: time.UnixMicro(1670803200000000)},
		Kline{Amount: 0, Count: 0, Open: 16916.600000, Close: 16925.000000, Low: 16903.900000, High: 16946.200000, Volume: 396.381670, Time: time.UnixMicro(1670817600000000)},
		Kline{Amount: 0, Count: 0, Open: 16925.000000, Close: 16988.400000, Low: 16903.600000, High: 17006.300000, Volume: 547.678505, Time: time.UnixMicro(1670832000000000)},
		Kline{Amount: 0, Count: 0, Open: 16988.500000, Close: 17008.000000, Low: 16919.600000, High: 17052.700000, Volume: 904.884566, Time: time.UnixMicro(1670846400000000)},
		Kline{Amount: 0, Count: 0, Open: 17008.000000, Close: 17044.900000, Low: 16982.300000, High: 17046.700000, Volume: 462.942724, Time: time.UnixMicro(1670860800000000)},
		Kline{Amount: 0, Count: 0, Open: 17044.900000, Close: 17207.100000, Low: 17044.900000, High: 17238.000000, Volume: 819.086040, Time: time.UnixMicro(1670875200000000)},
		Kline{Amount: 0, Count: 0, Open: 17207.900000, Close: 17150.800000, Low: 17123.100000, High: 17236.900000, Volume: 308.680498, Time: time.UnixMicro(1670889600000000)},
		Kline{Amount: 0, Count: 0, Open: 17151.500000, Close: 17162.600000, Low: 17083.100000, High: 17204.800000, Volume: 842.866117, Time: time.UnixMicro(1670904000000000)},
		Kline{Amount: 0, Count: 0, Open: 17162.500000, Close: 17446.800000, Low: 17155.700000, High: 17522.000000, Volume: 2513.322084, Time: time.UnixMicro(1670918400000000)},
		Kline{Amount: 0, Count: 0, Open: 17445.600000, Close: 17753.100000, Low: 17400.000000, High: 17999.000000, Volume: 6668.260253, Time: time.UnixMicro(1670932800000000)},
		Kline{Amount: 0, Count: 0, Open: 17753.200000, Close: 17744.200000, Low: 17622.200000, High: 17834.100000, Volume: 837.882497, Time: time.UnixMicro(1670947200000000)},
		Kline{Amount: 0, Count: 0, Open: 17744.300000, Close: 17780.000000, Low: 17710.100000, High: 17818.000000, Volume: 450.975224, Time: time.UnixMicro(1670961600000000)},
		Kline{Amount: 0, Count: 0, Open: 17780.100000, Close: 17786.000000, Low: 17753.300000, High: 17846.900000, Volume: 819.425291, Time: time.UnixMicro(1670976000000000)},
		Kline{Amount: 0, Count: 0, Open: 17786.100000, Close: 17812.100000, Low: 17742.100000, High: 17820.000000, Volume: 718.329035, Time: time.UnixMicro(1670990400000000)},
		Kline{Amount: 0, Count: 0, Open: 17812.200000, Close: 17829.100000, Low: 17800.000000, High: 17886.200000, Volume: 956.883668, Time: time.UnixMicro(1671004800000000)},
		Kline{Amount: 0, Count: 0, Open: 17829.100000, Close: 18061.000000, Low: 17821.600000, High: 18110.900000, Volume: 2172.914543, Time: time.UnixMicro(1671019200000000)},
		Kline{Amount: 0, Count: 0, Open: 18061.800000, Close: 17793.800000, Low: 17713.000000, High: 18377.100000, Volume: 9826.793530, Time: time.UnixMicro(1671033600000000)},
		Kline{Amount: 0, Count: 0, Open: 17793.800000, Close: 17806.800000, Low: 17661.000000, High: 17941.600000, Volume: 1446.503418, Time: time.UnixMicro(1671048000000000)},
		Kline{Amount: 0, Count: 0, Open: 17806.800000, Close: 17730.600000, Low: 17566.200000, High: 17851.700000, Volume: 1596.449952, Time: time.UnixMicro(1671062400000000)},
		Kline{Amount: 0, Count: 0, Open: 17730.300000, Close: 17694.900000, Low: 17686.900000, High: 17747.800000, Volume: 433.264500, Time: time.UnixMicro(1671076800000000)},
		Kline{Amount: 0, Count: 0, Open: 17695.400000, Close: 17724.200000, Low: 17615.300000, High: 17727.500000, Volume: 798.690263, Time: time.UnixMicro(1671091200000000)},
		Kline{Amount: 0, Count: 0, Open: 17724.200000, Close: 17428.200000, Low: 17377.100000, High: 17728.600000, Volume: 2288.725872, Time: time.UnixMicro(1671105600000000)},
		Kline{Amount: 0, Count: 0, Open: 17428.100000, Close: 17443.300000, Low: 17325.000000, High: 17470.000000, Volume: 1479.038138, Time: time.UnixMicro(1671120000000000)},
		Kline{Amount: 0, Count: 0, Open: 17443.200000, Close: 17359.200000, Low: 17277.000000, High: 17456.800000, Volume: 586.062283, Time: time.UnixMicro(1671134400000000)},
		Kline{Amount: 0, Count: 0, Open: 17359.200000, Close: 17425.600000, Low: 17350.800000, High: 17430.600000, Volume: 306.182230, Time: time.UnixMicro(1671148800000000)},
		Kline{Amount: 0, Count: 0, Open: 17425.700000, Close: 17504.400000, Low: 17377.300000, High: 17527.500000, Volume: 1002.449952, Time: time.UnixMicro(1671163200000000)},
		Kline{Amount: 0, Count: 0, Open: 17504.400000, Close: 17025.200000, Low: 16917.600000, High: 17518.200000, Volume: 3521.194846, Time: time.UnixMicro(1671177600000000)},
		Kline{Amount: 0, Count: 0, Open: 17025.200000, Close: 16971.500000, Low: 16882.500000, High: 17080.500000, Volume: 1302.484195, Time: time.UnixMicro(1671192000000000)},
		Kline{Amount: 0, Count: 0, Open: 16971.400000, Close: 16865.500000, Low: 16724.900000, High: 16977.600000, Volume: 3180.963677, Time: time.UnixMicro(1671206400000000)},
		Kline{Amount: 0, Count: 0, Open: 16865.400000, Close: 16629.200000, Low: 16530.000000, High: 16954.600000, Volume: 2959.484218, Time: time.UnixMicro(1671220800000000)},
		Kline{Amount: 0, Count: 0, Open: 16631.400000, Close: 16674.400000, Low: 16581.800000, High: 16726.600000, Volume: 932.150911, Time: time.UnixMicro(1671235200000000)},
		Kline{Amount: 0, Count: 0, Open: 16674.400000, Close: 16741.300000, Low: 16655.000000, High: 16766.100000, Volume: 1626.233414, Time: time.UnixMicro(1671249600000000)},
		Kline{Amount: 0, Count: 0, Open: 16741.400000, Close: 16703.500000, Low: 16662.200000, High: 16752.400000, Volume: 471.164976, Time: time.UnixMicro(1671264000000000)},
		Kline{Amount: 0, Count: 0, Open: 16703.600000, Close: 16700.700000, Low: 16647.700000, High: 16721.600000, Volume: 877.085938, Time: time.UnixMicro(1671278400000000)},
		Kline{Amount: 0, Count: 0, Open: 16700.600000, Close: 16689.800000, Low: 16650.100000, High: 16718.000000, Volume: 193.799563, Time: time.UnixMicro(1671292800000000)},
		Kline{Amount: 0, Count: 0, Open: 16689.800000, Close: 16778.800000, Low: 16683.200000, High: 16795.700000, Volume: 583.780258, Time: time.UnixMicro(1671307200000000)},
		Kline{Amount: 0, Count: 0, Open: 16778.800000, Close: 16729.200000, Low: 16704.000000, High: 16783.200000, Volume: 175.370287, Time: time.UnixMicro(1671321600000000)},
		Kline{Amount: 0, Count: 0, Open: 16729.200000, Close: 16729.200000, Low: 16729.200000, High: 16729.200000, Volume: 0.000000, Time: time.UnixMicro(1671336000000000)},
		Kline{Amount: 0, Count: 0, Open: 16746.100000, Close: 16742.600000, Low: 16740.500000, High: 16781.500000, Volume: 3.309444, Time: time.UnixMicro(1671350400000000)},
		Kline{Amount: 0, Count: 0, Open: 16742.600000, Close: 16742.600000, Low: 16742.600000, High: 16742.600000, Volume: 0.000000, Time: time.UnixMicro(1671364800000000)},
		Kline{Amount: 0, Count: 0, Open: 16749.000000, Close: 16751.800000, Low: 16718.000000, High: 16769.500000, Volume: 95.086825, Time: time.UnixMicro(1671379200000000)},
		Kline{Amount: 0, Count: 0, Open: 16751.800000, Close: 16738.700000, Low: 16727.200000, High: 16869.300000, Volume: 361.194821, Time: time.UnixMicro(1671393600000000)},
		Kline{Amount: 0, Count: 0, Open: 16737.500000, Close: 16676.400000, Low: 16624.700000, High: 16812.900000, Volume: 641.471164, Time: time.UnixMicro(1671408000000000)},
		Kline{Amount: 0, Count: 0, Open: 16676.400000, Close: 16719.500000, Low: 16671.300000, High: 16744.300000, Volume: 368.712945, Time: time.UnixMicro(1671422400000000)},
		Kline{Amount: 0, Count: 0, Open: 16719.500000, Close: 16735.200000, Low: 16707.300000, High: 16780.700000, Volume: 549.646643, Time: time.UnixMicro(1671436800000000)},
		Kline{Amount: 0, Count: 0, Open: 16735.000000, Close: 16670.100000, Low: 16666.000000, High: 16760.900000, Volume: 624.486423, Time: time.UnixMicro(1671451200000000)},
		Kline{Amount: 0, Count: 0, Open: 16670.000000, Close: 16557.600000, Low: 16520.200000, High: 16699.800000, Volume: 1901.994989, Time: time.UnixMicro(1671465600000000)},
		Kline{Amount: 0, Count: 0, Open: 16557.000000, Close: 16440.000000, Low: 16248.600000, High: 16646.400000, Volume: 1773.547216, Time: time.UnixMicro(1671480000000000)},
		Kline{Amount: 0, Count: 0, Open: 16440.000000, Close: 16726.700000, Low: 16400.000000, High: 16860.000000, Volume: 1907.776805, Time: time.UnixMicro(1671494400000000)},
		Kline{Amount: 0, Count: 0, Open: 16726.800000, Close: 16801.800000, Low: 16691.800000, High: 16878.200000, Volume: 1064.472830, Time: time.UnixMicro(1671508800000000)},
		Kline{Amount: 0, Count: 0, Open: 16801.600000, Close: 16813.500000, Low: 16765.100000, High: 16855.900000, Volume: 1562.076160, Time: time.UnixMicro(1671523200000000)},
		Kline{Amount: 0, Count: 0, Open: 16813.500000, Close: 16919.800000, Low: 16709.400000, High: 17061.400000, Volume: 1889.581205, Time: time.UnixMicro(1671537600000000)},
		Kline{Amount: 0, Count: 0, Open: 16920.300000, Close: 16899.700000, Low: 16773.200000, High: 16935.200000, Volume: 653.419741, Time: time.UnixMicro(1671552000000000)},
		Kline{Amount: 0, Count: 0, Open: 16899.700000, Close: 16893.800000, Low: 16839.400000, High: 16925.200000, Volume: 267.937948, Time: time.UnixMicro(1671566400000000)},
		Kline{Amount: 0, Count: 0, Open: 16894.900000, Close: 16838.700000, Low: 16794.700000, High: 16920.600000, Volume: 292.663623, Time: time.UnixMicro(1671580800000000)},
		Kline{Amount: 0, Count: 0, Open: 16838.700000, Close: 16816.100000, Low: 16770.300000, High: 16854.900000, Volume: 358.619883, Time: time.UnixMicro(1671595200000000)},
		Kline{Amount: 0, Count: 0, Open: 16816.200000, Close: 16883.400000, Low: 16799.300000, High: 16895.200000, Volume: 295.098231, Time: time.UnixMicro(1671609600000000)},
		Kline{Amount: 0, Count: 0, Open: 16883.400000, Close: 16860.700000, Low: 16725.500000, High: 16921.100000, Volume: 1037.311620, Time: time.UnixMicro(1671624000000000)},
		Kline{Amount: 0, Count: 0, Open: 16860.700000, Close: 16756.600000, Low: 16733.700000, High: 16866.100000, Volume: 398.242172, Time: time.UnixMicro(1671638400000000)},
		Kline{Amount: 0, Count: 0, Open: 16755.800000, Close: 16830.700000, Low: 16735.200000, High: 16833.300000, Volume: 223.828477, Time: time.UnixMicro(1671652800000000)},
		Kline{Amount: 0, Count: 0, Open: 16830.500000, Close: 16854.600000, Low: 16812.800000, High: 16862.300000, Volume: 228.483684, Time: time.UnixMicro(1671667200000000)},
		Kline{Amount: 0, Count: 0, Open: 16854.600000, Close: 16827.700000, Low: 16787.500000, High: 16857.000000, Volume: 223.174146, Time: time.UnixMicro(1671681600000000)},
		Kline{Amount: 0, Count: 0, Open: 16827.600000, Close: 16834.600000, Low: 16803.100000, High: 16866.600000, Volume: 490.546312, Time: time.UnixMicro(1671696000000000)},
		Kline{Amount: 0, Count: 0, Open: 16834.700000, Close: 16637.200000, Low: 16614.400000, High: 16844.400000, Volume: 1055.749759, Time: time.UnixMicro(1671710400000000)},
		Kline{Amount: 0, Count: 0, Open: 16637.500000, Close: 16664.900000, Low: 16568.200000, High: 16669.700000, Volume: 516.681937, Time: time.UnixMicro(1671724800000000)},
		Kline{Amount: 0, Count: 0, Open: 16664.900000, Close: 16824.400000, Low: 16651.500000, High: 16867.700000, Volume: 1193.795769, Time: time.UnixMicro(1671739200000000)},
		Kline{Amount: 0, Count: 0, Open: 16824.400000, Close: 16835.200000, Low: 16777.500000, High: 16897.300000, Volume: 518.988023, Time: time.UnixMicro(1671753600000000)},
		Kline{Amount: 0, Count: 0, Open: 16835.200000, Close: 16851.900000, Low: 16809.100000, High: 16859.300000, Volume: 457.497407, Time: time.UnixMicro(1671768000000000)},
		Kline{Amount: 0, Count: 0, Open: 16852.000000, Close: 16848.300000, Low: 16814.400000, High: 16892.000000, Volume: 623.236661, Time: time.UnixMicro(1671782400000000)},
		Kline{Amount: 0, Count: 0, Open: 16848.400000, Close: 16835.800000, Low: 16722.400000, High: 16965.700000, Volume: 1749.161076, Time: time.UnixMicro(1671796800000000)},
		Kline{Amount: 0, Count: 0, Open: 16835.000000, Close: 16847.700000, Low: 16818.800000, High: 16859.700000, Volume: 237.015230, Time: time.UnixMicro(1671811200000000)},
		Kline{Amount: 0, Count: 0, Open: 16847.700000, Close: 16784.700000, Low: 16771.300000, High: 16851.200000, Volume: 471.076563, Time: time.UnixMicro(1671825600000000)},
		Kline{Amount: 0, Count: 0, Open: 16784.100000, Close: 16844.600000, Low: 16783.900000, High: 16844.700000, Volume: 260.611224, Time: time.UnixMicro(1671840000000000)},
		Kline{Amount: 0, Count: 0, Open: 16844.700000, Close: 16849.800000, Low: 16810.100000, High: 16856.800000, Volume: 326.771865, Time: time.UnixMicro(1671854400000000)},
		Kline{Amount: 0, Count: 0, Open: 16849.800000, Close: 16827.000000, Low: 16822.000000, High: 16849.800000, Volume: 308.772089, Time: time.UnixMicro(1671868800000000)},
		Kline{Amount: 0, Count: 0, Open: 16827.000000, Close: 16838.300000, Low: 16816.300000, High: 16840.000000, Volume: 141.897928, Time: time.UnixMicro(1671883200000000)},
		Kline{Amount: 0, Count: 0, Open: 16838.400000, Close: 16844.600000, Low: 16825.900000, High: 16864.000000, Volume: 109.613006, Time: time.UnixMicro(1671897600000000)},
		Kline{Amount: 0, Count: 0, Open: 16845.400000, Close: 16837.500000, Low: 16815.000000, High: 16855.900000, Volume: 143.429519, Time: time.UnixMicro(1671912000000000)},
		Kline{Amount: 0, Count: 0, Open: 16837.600000, Close: 16834.500000, Low: 16827.000000, High: 16851.500000, Volume: 277.339525, Time: time.UnixMicro(1671926400000000)},
		Kline{Amount: 0, Count: 0, Open: 16834.500000, Close: 16823.500000, Low: 16803.500000, High: 16835.700000, Volume: 283.162546, Time: time.UnixMicro(1671940800000000)},
		Kline{Amount: 0, Count: 0, Open: 16823.400000, Close: 16822.400000, Low: 16817.000000, High: 16834.900000, Volume: 123.191432, Time: time.UnixMicro(1671955200000000)},
		Kline{Amount: 0, Count: 0, Open: 16822.400000, Close: 16795.200000, Low: 16717.900000, High: 16856.300000, Volume: 527.872524, Time: time.UnixMicro(1671969600000000)},
		Kline{Amount: 0, Count: 0, Open: 16795.300000, Close: 16795.000000, Low: 16731.300000, High: 16821.100000, Volume: 265.585928, Time: time.UnixMicro(1671984000000000)},
		Kline{Amount: 0, Count: 0, Open: 16795.000000, Close: 16837.700000, Low: 16758.700000, High: 16851.300000, Volume: 243.687774, Time: time.UnixMicro(1671998400000000)},
		Kline{Amount: 0, Count: 0, Open: 16837.700000, Close: 16898.800000, Low: 16824.300000, High: 16903.100000, Volume: 532.254703, Time: time.UnixMicro(1672012800000000)},
		Kline{Amount: 0, Count: 0, Open: 16898.800000, Close: 16848.600000, Low: 16843.000000, High: 16922.900000, Volume: 448.195721, Time: time.UnixMicro(1672027200000000)},
		Kline{Amount: 0, Count: 0, Open: 16848.600000, Close: 16865.600000, Low: 16822.500000, High: 16872.500000, Volume: 295.112783, Time: time.UnixMicro(1672041600000000)},
		Kline{Amount: 0, Count: 0, Open: 16865.700000, Close: 16810.500000, Low: 16793.900000, High: 16868.800000, Volume: 315.169619, Time: time.UnixMicro(1672056000000000)},
		Kline{Amount: 0, Count: 0, Open: 16810.600000, Close: 16843.100000, Low: 16802.100000, High: 16860.300000, Volume: 218.778410, Time: time.UnixMicro(1672070400000000)},
		Kline{Amount: 0, Count: 0, Open: 16843.200000, Close: 16922.100000, Low: 16823.600000, High: 16945.800000, Volume: 259.529234, Time: time.UnixMicro(1672084800000000)},
		Kline{Amount: 0, Count: 0, Open: 16922.100000, Close: 16876.500000, Low: 16841.000000, High: 16968.900000, Volume: 399.476247, Time: time.UnixMicro(1672099200000000)},
		Kline{Amount: 0, Count: 0, Open: 16876.600000, Close: 16869.000000, Low: 16855.900000, High: 16893.400000, Volume: 275.819010, Time: time.UnixMicro(1672113600000000)},
		Kline{Amount: 0, Count: 0, Open: 16869.000000, Close: 16834.800000, Low: 16812.600000, High: 16880.500000, Volume: 395.254789, Time: time.UnixMicro(1672128000000000)},
		Kline{Amount: 0, Count: 0, Open: 16834.700000, Close: 16802.000000, Low: 16738.400000, High: 16838.800000, Volume: 766.877744, Time: time.UnixMicro(1672142400000000)},
		Kline{Amount: 0, Count: 0, Open: 16801.900000, Close: 16688.300000, Low: 16594.000000, High: 16805.300000, Volume: 1134.551329, Time: time.UnixMicro(1672156800000000)},
		Kline{Amount: 0, Count: 0, Open: 16688.300000, Close: 16706.900000, Low: 16655.800000, High: 16721.800000, Volume: 471.688413, Time: time.UnixMicro(1672171200000000)},
		Kline{Amount: 0, Count: 0, Open: 16707.000000, Close: 16653.300000, Low: 16560.800000, High: 16752.900000, Volume: 713.213848, Time: time.UnixMicro(1672185600000000)},
		Kline{Amount: 0, Count: 0, Open: 16652.800000, Close: 16647.300000, Low: 16587.100000, High: 16675.400000, Volume: 488.631150, Time: time.UnixMicro(1672200000000000)},
		Kline{Amount: 0, Count: 0, Open: 16647.300000, Close: 16679.700000, Low: 16634.000000, High: 16688.900000, Volume: 445.741269, Time: time.UnixMicro(1672214400000000)},
		Kline{Amount: 0, Count: 0, Open: 16679.700000, Close: 16579.000000, Low: 16572.600000, High: 16783.800000, Volume: 1039.914002, Time: time.UnixMicro(1672228800000000)},
		Kline{Amount: 0, Count: 0, Open: 16579.000000, Close: 16596.400000, Low: 16567.800000, High: 16665.600000, Volume: 305.119495, Time: time.UnixMicro(1672243200000000)},
		Kline{Amount: 0, Count: 0, Open: 16596.300000, Close: 16548.900000, Low: 16471.300000, High: 16615.900000, Volume: 788.602985, Time: time.UnixMicro(1672257600000000)},
		Kline{Amount: 0, Count: 0, Open: 16549.000000, Close: 16559.900000, Low: 16493.900000, High: 16589.800000, Volume: 468.716553, Time: time.UnixMicro(1672272000000000)},
		Kline{Amount: 0, Count: 0, Open: 16560.000000, Close: 16558.800000, Low: 16538.000000, High: 16588.600000, Volume: 1049.208758, Time: time.UnixMicro(1672286400000000)},
		Kline{Amount: 0, Count: 0, Open: 16558.100000, Close: 16601.200000, Low: 16522.500000, High: 16629.100000, Volume: 682.740869, Time: time.UnixMicro(1672300800000000)},
		Kline{Amount: 0, Count: 0, Open: 16601.200000, Close: 16619.100000, Low: 16596.600000, High: 16665.000000, Volume: 550.214984, Time: time.UnixMicro(1672315200000000)},
		Kline{Amount: 0, Count: 0, Open: 16619.200000, Close: 16614.200000, Low: 16595.200000, High: 16656.100000, Volume: 242.155141, Time: time.UnixMicro(1672329600000000)},
		Kline{Amount: 0, Count: 0, Open: 16614.300000, Close: 16639.200000, Low: 16560.800000, High: 16654.800000, Volume: 259.679973, Time: time.UnixMicro(1672344000000000)},
		Kline{Amount: 0, Count: 0, Open: 16639.000000, Close: 16601.400000, Low: 16583.900000, High: 16649.200000, Volume: 292.186885, Time: time.UnixMicro(1672358400000000)},
		Kline{Amount: 0, Count: 0, Open: 16601.400000, Close: 16472.300000, Low: 16432.400000, High: 16614.500000, Volume: 1805.643725, Time: time.UnixMicro(1672372800000000)},
		Kline{Amount: 0, Count: 0, Open: 16472.200000, Close: 16497.600000, Low: 16462.300000, High: 16531.600000, Volume: 556.953627, Time: time.UnixMicro(1672387200000000)},
		Kline{Amount: 0, Count: 0, Open: 16497.900000, Close: 16558.800000, Low: 16340.900000, High: 16576.000000, Volume: 1492.231570, Time: time.UnixMicro(1672401600000000)},
		Kline{Amount: 0, Count: 0, Open: 16559.100000, Close: 16528.300000, Low: 16520.100000, High: 16577.100000, Volume: 358.641246, Time: time.UnixMicro(1672416000000000)},
		Kline{Amount: 0, Count: 0, Open: 16528.400000, Close: 16613.300000, Low: 16528.300000, High: 16676.300000, Volume: 707.405802, Time: time.UnixMicro(1672430400000000)},
		Kline{Amount: 0, Count: 0, Open: 16613.300000, Close: 16556.900000, Low: 16556.900000, High: 16620.000000, Volume: 201.361784, Time: time.UnixMicro(1672444800000000)},
		Kline{Amount: 0, Count: 0, Open: 16557.000000, Close: 16570.100000, Low: 16537.900000, High: 16570.900000, Volume: 226.238824, Time: time.UnixMicro(1672459200000000)},
		Kline{Amount: 0, Count: 0, Open: 16570.200000, Close: 16572.600000, Low: 16545.900000, High: 16585.000000, Volume: 331.818820, Time: time.UnixMicro(1672473600000000)},
		Kline{Amount: 0, Count: 0, Open: 16572.600000, Close: 16595.000000, Low: 16552.000000, High: 16644.900000, Volume: 727.869453, Time: time.UnixMicro(1672488000000000)},
		Kline{Amount: 0, Count: 0, Open: 16595.000000, Close: 16572.900000, Low: 16572.900000, High: 16614.800000, Volume: 250.542114, Time: time.UnixMicro(1672502400000000)},
		Kline{Amount: 0, Count: 0, Open: 16573.000000, Close: 16544.900000, Low: 16477.600000, High: 16577.000000, Volume: 404.839644, Time: time.UnixMicro(1672516800000000)},
		Kline{Amount: 0, Count: 0, Open: 16546.100000, Close: 16538.500000, Low: 16512.700000, High: 16562.800000, Volume: 360.815116, Time: time.UnixMicro(1672531200000000)},
		Kline{Amount: 0, Count: 0, Open: 16538.600000, Close: 16531.100000, Low: 16502.000000, High: 16552.000000, Volume: 277.635612, Time: time.UnixMicro(1672545600000000)},
		Kline{Amount: 0, Count: 0, Open: 16531.100000, Close: 16558.400000, Low: 16510.000000, High: 16558.500000, Volume: 140.791604, Time: time.UnixMicro(1672560000000000)},
		Kline{Amount: 0, Count: 0, Open: 16558.500000, Close: 16560.200000, Low: 16537.200000, High: 16574.000000, Volume: 134.338404, Time: time.UnixMicro(1672574400000000)},
		Kline{Amount: 0, Count: 0, Open: 16560.000000, Close: 16604.700000, Low: 16559.900000, High: 16623.800000, Volume: 165.985616, Time: time.UnixMicro(1672588800000000)},
		Kline{Amount: 0, Count: 0, Open: 16604.600000, Close: 16615.100000, Low: 16590.000000, High: 16625.000000, Volume: 104.064644, Time: time.UnixMicro(1672603200000000)},
		Kline{Amount: 0, Count: 0, Open: 16615.100000, Close: 16662.900000, Low: 16548.100000, High: 16710.700000, Volume: 1021.047866, Time: time.UnixMicro(1672617600000000)},
		Kline{Amount: 0, Count: 0, Open: 16663.000000, Close: 16720.300000, Low: 16623.800000, High: 16765.500000, Volume: 842.614797, Time: time.UnixMicro(1672632000000000)},
		Kline{Amount: 0, Count: 0, Open: 16720.200000, Close: 16736.000000, Low: 16704.300000, High: 16768.400000, Volume: 619.730676, Time: time.UnixMicro(1672646400000000)},
		Kline{Amount: 0, Count: 0, Open: 16736.100000, Close: 16735.500000, Low: 16672.300000, High: 16750.000000, Volume: 504.085908, Time: time.UnixMicro(1672660800000000)},
		Kline{Amount: 0, Count: 0, Open: 16735.400000, Close: 16735.600000, Low: 16708.400000, High: 16741.000000, Volume: 161.605036, Time: time.UnixMicro(1672675200000000)},
		Kline{Amount: 0, Count: 0, Open: 16735.700000, Close: 16673.300000, Low: 16664.400000, High: 16799.300000, Volume: 370.283537, Time: time.UnixMicro(1672689600000000)},
		Kline{Amount: 0, Count: 0, Open: 16673.300000, Close: 16692.900000, Low: 16652.300000, High: 16711.000000, Volume: 210.745040, Time: time.UnixMicro(1672704000000000)},
		Kline{Amount: 0, Count: 0, Open: 16693.000000, Close: 16730.600000, Low: 16683.200000, High: 16774.000000, Volume: 291.400248, Time: time.UnixMicro(1672718400000000)},
		Kline{Amount: 0, Count: 0, Open: 16730.700000, Close: 16722.500000, Low: 16699.800000, High: 16758.500000, Volume: 258.293699, Time: time.UnixMicro(1672732800000000)},
		Kline{Amount: 0, Count: 0, Open: 16722.600000, Close: 16679.800000, Low: 16621.300000, High: 16771.600000, Volume: 852.532111, Time: time.UnixMicro(1672747200000000)},
		Kline{Amount: 0, Count: 0, Open: 16679.800000, Close: 16646.200000, Low: 16610.400000, High: 16688.600000, Volume: 467.602223, Time: time.UnixMicro(1672761600000000)},
		Kline{Amount: 0, Count: 0, Open: 16646.200000, Close: 16676.800000, Low: 16644.700000, High: 16695.200000, Volume: 365.671472, Time: time.UnixMicro(1672776000000000)},
		Kline{Amount: 0, Count: 0, Open: 16676.800000, Close: 16865.300000, Low: 16655.400000, High: 16877.400000, Volume: 1402.730874, Time: time.UnixMicro(1672790400000000)},
		Kline{Amount: 0, Count: 0, Open: 16866.300000, Close: 16873.300000, Low: 16843.100000, High: 16910.000000, Volume: 879.467180, Time: time.UnixMicro(1672804800000000)},
		Kline{Amount: 0, Count: 0, Open: 16873.200000, Close: 16838.500000, Low: 16834.100000, High: 16915.300000, Volume: 756.573760, Time: time.UnixMicro(1672819200000000)},
		Kline{Amount: 0, Count: 0, Open: 16838.600000, Close: 16855.000000, Low: 16768.900000, High: 16864.100000, Volume: 725.391092, Time: time.UnixMicro(1672833600000000)},
		Kline{Amount: 0, Count: 0, Open: 16855.800000, Close: 16841.000000, Low: 16791.500000, High: 16988.000000, Volume: 1127.581826, Time: time.UnixMicro(1672848000000000)},
		Kline{Amount: 0, Count: 0, Open: 16841.000000, Close: 16852.000000, Low: 16780.000000, High: 16866.600000, Volume: 288.285507, Time: time.UnixMicro(1672862400000000)},
		Kline{Amount: 0, Count: 0, Open: 16852.000000, Close: 16828.000000, Low: 16812.900000, High: 16877.800000, Volume: 344.889334, Time: time.UnixMicro(1672876800000000)},
		Kline{Amount: 0, Count: 0, Open: 16828.000000, Close: 16823.000000, Low: 16793.200000, High: 16848.300000, Volume: 285.937746, Time: time.UnixMicro(1672891200000000)},
		Kline{Amount: 0, Count: 0, Open: 16822.900000, Close: 16834.200000, Low: 16785.000000, High: 16837.200000, Volume: 365.801985, Time: time.UnixMicro(1672905600000000)},
		Kline{Amount: 0, Count: 0, Open: 16834.300000, Close: 16846.800000, Low: 16760.100000, High: 16879.200000, Volume: 630.258032, Time: time.UnixMicro(1672920000000000)},
		Kline{Amount: 0, Count: 0, Open: 16846.800000, Close: 16866.100000, Low: 16801.200000, High: 16876.700000, Volume: 440.392933, Time: time.UnixMicro(1672934400000000)},
		Kline{Amount: 0, Count: 0, Open: 16866.200000, Close: 16832.500000, Low: 16823.000000, High: 16869.100000, Volume: 315.057427, Time: time.UnixMicro(1672948800000000)},
		Kline{Amount: 0, Count: 0, Open: 16832.500000, Close: 16833.500000, Low: 16804.400000, High: 16868.000000, Volume: 244.720624, Time: time.UnixMicro(1672963200000000)},
		Kline{Amount: 0, Count: 0, Open: 16833.300000, Close: 16793.300000, Low: 16780.000000, High: 16834.800000, Volume: 358.326255, Time: time.UnixMicro(1672977600000000)},
		Kline{Amount: 0, Count: 0, Open: 16793.300000, Close: 16741.500000, Low: 16708.900000, High: 16810.000000, Volume: 573.255736, Time: time.UnixMicro(1672992000000000)},
		Kline{Amount: 0, Count: 0, Open: 16741.500000, Close: 16843.200000, Low: 16680.700000, High: 16846.600000, Volume: 1580.188683, Time: time.UnixMicro(1673006400000000)},
		Kline{Amount: 0, Count: 0, Open: 16843.300000, Close: 16938.500000, Low: 16809.900000, High: 16962.100000, Volume: 611.634583, Time: time.UnixMicro(1673020800000000)},
		Kline{Amount: 0, Count: 0, Open: 16938.500000, Close: 16950.600000, Low: 16890.400000, High: 17034.500000, Volume: 895.111311, Time: time.UnixMicro(1673035200000000)},
		Kline{Amount: 0, Count: 0, Open: 16950.700000, Close: 16953.200000, Low: 16935.200000, High: 16980.000000, Volume: 315.042962, Time: time.UnixMicro(1673049600000000)},
		Kline{Amount: 0, Count: 0, Open: 16953.200000, Close: 16947.900000, Low: 16924.600000, High: 16954.000000, Volume: 150.913930, Time: time.UnixMicro(1673064000000000)},
		Kline{Amount: 0, Count: 0, Open: 16947.900000, Close: 16922.900000, Low: 16914.200000, High: 16953.100000, Volume: 216.129401, Time: time.UnixMicro(1673078400000000)},
		Kline{Amount: 0, Count: 0, Open: 16923.000000, Close: 16941.100000, Low: 16908.700000, High: 16947.000000, Volume: 383.029827, Time: time.UnixMicro(1673092800000000)},
		Kline{Amount: 0, Count: 0, Open: 16941.100000, Close: 16945.300000, Low: 16917.100000, High: 16945.300000, Volume: 110.935535, Time: time.UnixMicro(1673107200000000)},
		Kline{Amount: 0, Count: 0, Open: 16945.300000, Close: 16945.500000, Low: 16931.800000, High: 16950.000000, Volume: 107.010458, Time: time.UnixMicro(1673121600000000)},
		Kline{Amount: 0, Count: 0, Open: 16945.600000, Close: 16944.700000, Low: 16915.000000, High: 16954.900000, Volume: 202.038988, Time: time.UnixMicro(1673136000000000)},
		Kline{Amount: 0, Count: 0, Open: 16944.700000, Close: 16954.900000, Low: 16937.200000, High: 16955.700000, Volume: 167.082698, Time: time.UnixMicro(1673150400000000)},
		Kline{Amount: 0, Count: 0, Open: 16954.900000, Close: 16927.800000, Low: 16927.800000, High: 16965.700000, Volume: 190.765251, Time: time.UnixMicro(1673164800000000)},
		Kline{Amount: 0, Count: 0, Open: 16927.800000, Close: 17000.800000, Low: 16921.400000, High: 17019.200000, Volume: 458.129743, Time: time.UnixMicro(1673179200000000)},
		Kline{Amount: 0, Count: 0, Open: 17000.800000, Close: 16929.400000, Low: 16921.300000, High: 17022.900000, Volume: 329.174463, Time: time.UnixMicro(1673193600000000)},
		Kline{Amount: 0, Count: 0, Open: 16929.300000, Close: 17129.300000, Low: 16927.000000, High: 17179.700000, Volume: 1176.086077, Time: time.UnixMicro(1673208000000000)},
		Kline{Amount: 0, Count: 0, Open: 17129.300000, Close: 17198.200000, Low: 17107.200000, High: 17252.300000, Volume: 1301.402618, Time: time.UnixMicro(1673222400000000)},
		Kline{Amount: 0, Count: 0, Open: 17198.300000, Close: 17202.100000, Low: 17185.000000, High: 17255.200000, Volume: 603.244719, Time: time.UnixMicro(1673236800000000)},
		Kline{Amount: 0, Count: 0, Open: 17202.000000, Close: 17239.900000, Low: 17189.400000, High: 17283.600000, Volume: 1100.531985, Time: time.UnixMicro(1673251200000000)},
		Kline{Amount: 0, Count: 0, Open: 17240.000000, Close: 17273.500000, Low: 17196.700000, High: 17333.000000, Volume: 1154.773678, Time: time.UnixMicro(1673265600000000)},
		Kline{Amount: 0, Count: 0, Open: 17273.400000, Close: 17296.800000, Low: 17271.900000, High: 17395.700000, Volume: 1172.822741, Time: time.UnixMicro(1673280000000000)},
		Kline{Amount: 0, Count: 0, Open: 17296.800000, Close: 17181.200000, Low: 17133.300000, High: 17308.600000, Volume: 1163.556091, Time: time.UnixMicro(1673294400000000)},
		Kline{Amount: 0, Count: 0, Open: 17181.100000, Close: 17208.000000, Low: 17150.300000, High: 17234.200000, Volume: 354.471384, Time: time.UnixMicro(1673308800000000)},
		Kline{Amount: 0, Count: 0, Open: 17208.000000, Close: 17202.800000, Low: 17190.800000, High: 17233.000000, Volume: 378.534155, Time: time.UnixMicro(1673323200000000)},
		Kline{Amount: 0, Count: 0, Open: 17202.800000, Close: 17255.100000, Low: 17191.800000, High: 17282.800000, Volume: 535.525705, Time: time.UnixMicro(1673337600000000)},
		Kline{Amount: 0, Count: 0, Open: 17254.900000, Close: 17323.700000, Low: 17213.200000, High: 17371.200000, Volume: 1056.171645, Time: time.UnixMicro(1673352000000000)},
		Kline{Amount: 0, Count: 0, Open: 17324.600000, Close: 17427.800000, Low: 17282.900000, High: 17448.000000, Volume: 1272.159558, Time: time.UnixMicro(1673366400000000)},
		Kline{Amount: 0, Count: 0, Open: 17427.800000, Close: 17439.000000, Low: 17403.000000, High: 17486.400000, Volume: 881.992420, Time: time.UnixMicro(1673380800000000)},
		Kline{Amount: 0, Count: 0, Open: 17439.000000, Close: 17414.200000, Low: 17369.000000, High: 17502.800000, Volume: 668.840197, Time: time.UnixMicro(1673395200000000)},
		Kline{Amount: 0, Count: 0, Open: 17414.200000, Close: 17445.500000, Low: 17389.000000, High: 17459.000000, Volume: 601.083046, Time: time.UnixMicro(1673409600000000)},
		Kline{Amount: 0, Count: 0, Open: 17445.600000, Close: 17437.800000, Low: 17410.400000, High: 17470.400000, Volume: 354.762211, Time: time.UnixMicro(1673424000000000)},
		Kline{Amount: 0, Count: 0, Open: 17437.800000, Close: 17334.700000, Low: 17320.900000, High: 17437.800000, Volume: 651.181021, Time: time.UnixMicro(1673438400000000)},
		Kline{Amount: 0, Count: 0, Open: 17334.500000, Close: 17548.300000, Low: 17328.700000, High: 17555.000000, Volume: 807.764946, Time: time.UnixMicro(1673452800000000)},
		Kline{Amount: 0, Count: 0, Open: 17548.400000, Close: 17945.200000, Low: 17511.900000, High: 18000.000000, Volume: 2207.247739, Time: time.UnixMicro(1673467200000000)},
		Kline{Amount: 0, Count: 0, Open: 17945.100000, Close: 18222.700000, Low: 17909.000000, High: 18379.900000, Volume: 3623.841771, Time: time.UnixMicro(1673481600000000)},
		Kline{Amount: 0, Count: 0, Open: 18222.700000, Close: 18137.200000, Low: 18077.200000, High: 18237.800000, Volume: 1123.599366, Time: time.UnixMicro(1673496000000000)},
		Kline{Amount: 0, Count: 0, Open: 18137.300000, Close: 18198.900000, Low: 18073.100000, High: 18218.100000, Volume: 1151.490903, Time: time.UnixMicro(1673510400000000)},
		Kline{Amount: 0, Count: 0, Open: 18198.200000, Close: 18087.700000, Low: 17932.000000, High: 18350.000000, Volume: 5765.535639, Time: time.UnixMicro(1673524800000000)},
		Kline{Amount: 0, Count: 0, Open: 18087.000000, Close: 18883.000000, Low: 18060.400000, High: 19056.900000, Volume: 5809.878200, Time: time.UnixMicro(1673539200000000)},
		Kline{Amount: 0, Count: 0, Open: 18883.000000, Close: 18844.200000, Low: 18759.100000, High: 19112.600000, Volume: 1908.697018, Time: time.UnixMicro(1673553600000000)},
		Kline{Amount: 0, Count: 0, Open: 18844.200000, Close: 18875.700000, Low: 18715.900000, High: 18885.000000, Volume: 1024.936724, Time: time.UnixMicro(1673568000000000)},
		Kline{Amount: 0, Count: 0, Open: 18875.700000, Close: 18818.200000, Low: 18781.100000, High: 18878.600000, Volume: 789.316266, Time: time.UnixMicro(1673582400000000)},
		Kline{Amount: 0, Count: 0, Open: 18818.200000, Close: 18922.500000, Low: 18810.300000, High: 19050.000000, Volume: 1399.365177, Time: time.UnixMicro(1673596800000000)},
		Kline{Amount: 0, Count: 0, Open: 18922.600000, Close: 19260.200000, Low: 18811.500000, High: 19306.400000, Volume: 3080.935948, Time: time.UnixMicro(1673611200000000)},
		Kline{Amount: 0, Count: 0, Open: 19260.200000, Close: 19361.000000, Low: 19060.400000, High: 19400.000000, Volume: 2155.569070, Time: time.UnixMicro(1673625600000000)},
		Kline{Amount: 0, Count: 0, Open: 19362.400000, Close: 19928.400000, Low: 19319.900000, High: 20004.000000, Volume: 3915.732137, Time: time.UnixMicro(1673640000000000)},
		Kline{Amount: 0, Count: 0, Open: 19928.500000, Close: 20976.600000, Low: 19895.000000, High: 21482.800000, Volume: 10354.599862, Time: time.UnixMicro(1673654400000000)},
		Kline{Amount: 0, Count: 0, Open: 20972.800000, Close: 20905.400000, Low: 20786.800000, High: 21044.200000, Volume: 1495.597580, Time: time.UnixMicro(1673668800000000)},
		Kline{Amount: 0, Count: 0, Open: 20904.200000, Close: 20724.700000, Low: 20239.100000, High: 21050.700000, Volume: 4928.200993, Time: time.UnixMicro(1673683200000000)},
		Kline{Amount: 0, Count: 0, Open: 20724.700000, Close: 20792.400000, Low: 20536.300000, High: 21200.000000, Volume: 4807.905171, Time: time.UnixMicro(1673697600000000)},
		Kline{Amount: 0, Count: 0, Open: 20788.500000, Close: 20788.700000, Low: 20660.100000, High: 20921.400000, Volume: 2207.863902, Time: time.UnixMicro(1673712000000000)},
		Kline{Amount: 0, Count: 0, Open: 20787.600000, Close: 20950.700000, Low: 20767.700000, High: 21079.900000, Volume: 1185.100823, Time: time.UnixMicro(1673726400000000)},
		Kline{Amount: 0, Count: 0, Open: 20950.700000, Close: 20718.300000, Low: 20555.300000, High: 20996.700000, Volume: 1716.982421, Time: time.UnixMicro(1673740800000000)},
		Kline{Amount: 0, Count: 0, Open: 20720.900000, Close: 20736.200000, Low: 20666.600000, High: 20781.000000, Volume: 720.979415, Time: time.UnixMicro(1673755200000000)},
		Kline{Amount: 0, Count: 0, Open: 20736.300000, Close: 20726.300000, Low: 20572.600000, High: 20763.400000, Volume: 932.834968, Time: time.UnixMicro(1673769600000000)},
		Kline{Amount: 0, Count: 0, Open: 20726.400000, Close: 20911.700000, Low: 20657.800000, High: 20998.000000, Volume: 1536.329846, Time: time.UnixMicro(1673784000000000)},
		Kline{Amount: 0, Count: 0, Open: 20911.700000, Close: 20867.600000, Low: 20658.100000, High: 21063.800000, Volume: 2097.798964, Time: time.UnixMicro(1673798400000000)},
		Kline{Amount: 0, Count: 0, Open: 20867.700000, Close: 20874.400000, Low: 20734.200000, High: 20954.000000, Volume: 571.301498, Time: time.UnixMicro(1673812800000000)},
		Kline{Amount: 0, Count: 0, Open: 20873.100000, Close: 21076.800000, Low: 20779.000000, High: 21431.800000, Volume: 2885.616302, Time: time.UnixMicro(1673827200000000)},
		Kline{Amount: 0, Count: 0, Open: 21076.800000, Close: 21111.000000, Low: 21045.800000, High: 21266.300000, Volume: 1124.775085, Time: time.UnixMicro(1673841600000000)},
		Kline{Amount: 0, Count: 0, Open: 21111.000000, Close: 20818.300000, Low: 20640.300000, High: 21111.000000, Volume: 2519.893562, Time: time.UnixMicro(1673856000000000)},
		Kline{Amount: 0, Count: 0, Open: 20818.300000, Close: 20993.900000, Low: 20603.000000, High: 21067.900000, Volume: 2308.012727, Time: time.UnixMicro(1673870400000000)},
		Kline{Amount: 0, Count: 0, Open: 20993.900000, Close: 21307.800000, Low: 20917.100000, High: 21400.000000, Volume: 2340.605687, Time: time.UnixMicro(1673884800000000)},
		Kline{Amount: 0, Count: 0, Open: 21309.800000, Close: 21183.800000, Low: 21046.700000, High: 21465.300000, Volume: 1987.437889, Time: time.UnixMicro(1673899200000000)},
		Kline{Amount: 0, Count: 0, Open: 21183.900000, Close: 21121.900000, Low: 20839.700000, High: 21292.900000, Volume: 1399.307412, Time: time.UnixMicro(1673913600000000)},
		Kline{Amount: 0, Count: 0, Open: 21121.900000, Close: 21134.700000, Low: 21059.500000, High: 21228.700000, Volume: 598.483018, Time: time.UnixMicro(1673928000000000)},
		Kline{Amount: 0, Count: 0, Open: 21134.700000, Close: 21220.500000, Low: 21081.000000, High: 21243.500000, Volume: 734.181468, Time: time.UnixMicro(1673942400000000)},
		Kline{Amount: 0, Count: 0, Open: 21220.400000, Close: 21173.400000, Low: 21016.200000, High: 21634.500000, Volume: 5303.264013, Time: time.UnixMicro(1673956800000000)},
		Kline{Amount: 0, Count: 0, Open: 21173.400000, Close: 21332.600000, Low: 21083.400000, High: 21411.000000, Volume: 1948.303496, Time: time.UnixMicro(1673971200000000)},
		Kline{Amount: 0, Count: 0, Open: 21332.600000, Close: 21125.500000, Low: 21121.800000, High: 21403.400000, Volume: 1128.136017, Time: time.UnixMicro(1673985600000000)},
		Kline{Amount: 0, Count: 0, Open: 21127.000000, Close: 21258.800000, Low: 21100.300000, High: 21370.000000, Volume: 1152.548427, Time: time.UnixMicro(1674000000000000)},
		Kline{Amount: 0, Count: 0, Open: 21258.900000, Close: 21304.000000, Low: 21230.900000, High: 21318.100000, Volume: 720.469362, Time: time.UnixMicro(1674014400000000)},
		Kline{Amount: 0, Count: 0, Open: 21304.100000, Close: 21205.600000, Low: 21145.000000, High: 21311.900000, Volume: 1111.557985, Time: time.UnixMicro(1674028800000000)},
		Kline{Amount: 0, Count: 0, Open: 21205.700000, Close: 21035.100000, Low: 20837.900000, High: 21664.600000, Volume: 5154.703216, Time: time.UnixMicro(1674043200000000)},
		Kline{Amount: 0, Count: 0, Open: 21035.100000, Close: 20900.000000, Low: 20400.000000, High: 21135.600000, Volume: 6935.214338, Time: time.UnixMicro(1674057600000000)},
		Kline{Amount: 0, Count: 0, Open: 20901.200000, Close: 20674.200000, Low: 20623.800000, High: 20932.800000, Volume: 2849.187239, Time: time.UnixMicro(1674072000000000)},
		Kline{Amount: 0, Count: 0, Open: 20672.500000, Close: 20747.400000, Low: 20660.000000, High: 20814.600000, Volume: 1193.617694, Time: time.UnixMicro(1674086400000000)},
		Kline{Amount: 0, Count: 0, Open: 20747.500000, Close: 20816.700000, Low: 20745.000000, High: 20872.300000, Volume: 615.823263, Time: time.UnixMicro(1674100800000000)},
		Kline{Amount: 0, Count: 0, Open: 20819.000000, Close: 20743.100000, Low: 20706.700000, High: 20846.200000, Volume: 747.265590, Time: time.UnixMicro(1674115200000000)},
		Kline{Amount: 0, Count: 0, Open: 20743.100000, Close: 20873.800000, Low: 20662.300000, High: 20956.300000, Volume: 1013.779986, Time: time.UnixMicro(1674129600000000)},
		Kline{Amount: 0, Count: 0, Open: 20873.900000, Close: 21092.300000, Low: 20792.500000, High: 21137.400000, Volume: 754.172162, Time: time.UnixMicro(1674144000000000)},
		Kline{Amount: 0, Count: 0, Open: 21092.300000, Close: 21072.700000, Low: 20920.000000, High: 21179.100000, Volume: 796.570053, Time: time.UnixMicro(1674158400000000)},
		Kline{Amount: 0, Count: 0, Open: 21072.700000, Close: 21088.300000, Low: 21011.400000, High: 21215.300000, Volume: 489.698630, Time: time.UnixMicro(1674172800000000)},
		Kline{Amount: 0, Count: 0, Open: 21087.300000, Close: 20956.200000, Low: 20860.400000, High: 21123.300000, Volume: 687.687102, Time: time.UnixMicro(1674187200000000)},
		Kline{Amount: 0, Count: 0, Open: 20956.700000, Close: 20963.100000, Low: 20902.900000, High: 20990.000000, Volume: 361.783211, Time: time.UnixMicro(1674201600000000)},
		Kline{Amount: 0, Count: 0, Open: 20963.200000, Close: 21144.900000, Low: 20954.700000, High: 21245.500000, Volume: 670.842240, Time: time.UnixMicro(1674216000000000)},
		Kline{Amount: 0, Count: 0, Open: 21145.000000, Close: 21496.900000, Low: 21143.600000, High: 21521.200000, Volume: 1246.286900, Time: time.UnixMicro(1674230400000000)},
		Kline{Amount: 0, Count: 0, Open: 21495.400000, Close: 22665.600000, Low: 21495.400000, High: 22747.000000, Volume: 6252.745925, Time: time.UnixMicro(1674244800000000)},
		Kline{Amount: 0, Count: 0, Open: 22665.700000, Close: 22546.500000, Low: 22428.600000, High: 22789.200000, Volume: 1394.164944, Time: time.UnixMicro(1674259200000000)},
		Kline{Amount: 0, Count: 0, Open: 22546.400000, Close: 22637.000000, Low: 22513.400000, High: 22669.900000, Volume: 543.404408, Time: time.UnixMicro(1674273600000000)},
		Kline{Amount: 0, Count: 0, Open: 22637.100000, Close: 22903.000000, Low: 22568.900000, High: 23327.800000, Volume: 4495.288913, Time: time.UnixMicro(1674288000000000)},
		Kline{Amount: 0, Count: 0, Open: 22904.800000, Close: 22996.000000, Low: 22849.400000, High: 23274.900000, Volume: 2253.158714, Time: time.UnixMicro(1674302400000000)},
		Kline{Amount: 0, Count: 0, Open: 22995.500000, Close: 23275.100000, Low: 22982.700000, High: 23379.000000, Volume: 2372.367300, Time: time.UnixMicro(1674316800000000)},
		Kline{Amount: 0, Count: 0, Open: 23275.200000, Close: 22779.000000, Low: 22681.700000, High: 23300.900000, Volume: 2039.368356, Time: time.UnixMicro(1674331200000000)},
		Kline{Amount: 0, Count: 0, Open: 22780.000000, Close: 22760.600000, Low: 22610.300000, High: 22976.600000, Volume: 1238.769074, Time: time.UnixMicro(1674345600000000)},
		Kline{Amount: 0, Count: 0, Open: 22758.400000, Close: 22895.200000, Low: 22721.100000, High: 22962.300000, Volume: 708.903096, Time: time.UnixMicro(1674360000000000)},
		Kline{Amount: 0, Count: 0, Open: 22894.400000, Close: 22780.000000, Low: 22729.400000, High: 22930.000000, Volume: 623.593420, Time: time.UnixMicro(1674374400000000)},
		Kline{Amount: 0, Count: 0, Open: 22784.700000, Close: 22788.700000, Low: 22639.900000, High: 23074.400000, Volume: 1657.233319, Time: time.UnixMicro(1674388800000000)},
		Kline{Amount: 0, Count: 0, Open: 22788.700000, Close: 22718.400000, Low: 22678.000000, High: 22956.100000, Volume: 1170.841580, Time: time.UnixMicro(1674403200000000)},
		Kline{Amount: 0, Count: 0, Open: 22720.000000, Close: 22708.500000, Low: 22306.600000, High: 22787.000000, Volume: 2129.384838, Time: time.UnixMicro(1674417600000000)},
		Kline{Amount: 0, Count: 0, Open: 22708.400000, Close: 22661.800000, Low: 22655.600000, High: 22823.700000, Volume: 497.997372, Time: time.UnixMicro(1674432000000000)},
		Kline{Amount: 0, Count: 0, Open: 22661.800000, Close: 22692.400000, Low: 22609.500000, High: 22788.000000, Volume: 572.914949, Time: time.UnixMicro(1674446400000000)},
		Kline{Amount: 0, Count: 0, Open: 22692.500000, Close: 22903.100000, Low: 22664.200000, High: 22955.000000, Volume: 941.980035, Time: time.UnixMicro(1674460800000000)},
		Kline{Amount: 0, Count: 0, Open: 22903.200000, Close: 22852.000000, Low: 22498.100000, High: 23123.200000, Volume: 3379.537467, Time: time.UnixMicro(1674475200000000)},
		Kline{Amount: 0, Count: 0, Open: 22852.100000, Close: 22794.800000, Low: 22658.000000, High: 23185.200000, Volume: 1388.997044, Time: time.UnixMicro(1674489600000000)},
		Kline{Amount: 0, Count: 0, Open: 22794.700000, Close: 22914.800000, Low: 22760.800000, High: 23043.400000, Volume: 654.245041, Time: time.UnixMicro(1674504000000000)},
		Kline{Amount: 0, Count: 0, Open: 22914.900000, Close: 23082.600000, Low: 22865.600000, High: 23143.100000, Volume: 689.925562, Time: time.UnixMicro(1674518400000000)},
		Kline{Amount: 0, Count: 0, Open: 23082.500000, Close: 23053.700000, Low: 23010.700000, High: 23159.800000, Volume: 472.303044, Time: time.UnixMicro(1674532800000000)},
		Kline{Amount: 0, Count: 0, Open: 23053.600000, Close: 22919.300000, Low: 22773.300000, High: 23089.500000, Volume: 1465.523345, Time: time.UnixMicro(1674547200000000)},
		Kline{Amount: 0, Count: 0, Open: 22919.300000, Close: 22928.100000, Low: 22700.000000, High: 23031.300000, Volume: 1842.015280, Time: time.UnixMicro(1674561600000000)},
		Kline{Amount: 0, Count: 0, Open: 22928.200000, Close: 23016.900000, Low: 22787.200000, High: 23083.800000, Volume: 786.263381, Time: time.UnixMicro(1674576000000000)},
		Kline{Amount: 0, Count: 0, Open: 23018.600000, Close: 22629.600000, Low: 22457.100000, High: 23079.600000, Volume: 2178.152967, Time: time.UnixMicro(1674590400000000)},
		Kline{Amount: 0, Count: 0, Open: 22629.600000, Close: 22642.600000, Low: 22332.100000, High: 22700.000000, Volume: 1585.901451, Time: time.UnixMicro(1674604800000000)},
		Kline{Amount: 0, Count: 0, Open: 22642.700000, Close: 22713.100000, Low: 22604.300000, High: 22768.000000, Volume: 665.434360, Time: time.UnixMicro(1674619200000000)},
		Kline{Amount: 0, Count: 0, Open: 22713.100000, Close: 22598.000000, Low: 22468.700000, High: 22738.600000, Volume: 905.590540, Time: time.UnixMicro(1674633600000000)},
		Kline{Amount: 0, Count: 0, Open: 22597.900000, Close: 22579.100000, Low: 22344.900000, High: 22745.100000, Volume: 1446.943127, Time: time.UnixMicro(1674648000000000)},
		Kline{Amount: 0, Count: 0, Open: 22579.200000, Close: 22752.300000, Low: 22518.300000, High: 22760.700000, Volume: 1068.552599, Time: time.UnixMicro(1674662400000000)},
		Kline{Amount: 0, Count: 0, Open: 22750.600000, Close: 23052.600000, Low: 22750.600000, High: 23818.700000, Volume: 5566.880580, Time: time.UnixMicro(1674676800000000)},
		Kline{Amount: 0, Count: 0, Open: 23060.800000, Close: 23133.700000, Low: 23046.900000, High: 23284.000000, Volume: 977.205671, Time: time.UnixMicro(1674691200000000)},
		Kline{Amount: 0, Count: 0, Open: 23134.000000, Close: 22967.700000, Low: 22900.000000, High: 23192.400000, Volume: 993.103166, Time: time.UnixMicro(1674705600000000)},
		Kline{Amount: 0, Count: 0, Open: 22967.800000, Close: 22988.300000, Low: 22847.600000, High: 23046.100000, Volume: 1010.571831, Time: time.UnixMicro(1674720000000000)},
		Kline{Amount: 0, Count: 0, Open: 22988.300000, Close: 22958.100000, Low: 22860.400000, High: 23240.800000, Volume: 1825.790962, Time: time.UnixMicro(1674734400000000)},
		Kline{Amount: 0, Count: 0, Open: 22957.700000, Close: 23069.800000, Low: 22858.700000, High: 23106.400000, Volume: 824.098370, Time: time.UnixMicro(1674748800000000)},
		Kline{Amount: 0, Count: 0, Open: 23070.300000, Close: 23010.000000, Low: 22967.900000, High: 23164.700000, Volume: 537.495442, Time: time.UnixMicro(1674763200000000)},
		Kline{Amount: 0, Count: 0, Open: 23010.100000, Close: 22771.800000, Low: 22526.300000, High: 23072.800000, Volume: 2393.003886, Time: time.UnixMicro(1674777600000000)},
		Kline{Amount: 0, Count: 0, Open: 22771.900000, Close: 23072.700000, Low: 22764.200000, High: 23104.500000, Volume: 1572.132574, Time: time.UnixMicro(1674792000000000)},
		Kline{Amount: 0, Count: 0, Open: 23072.800000, Close: 22966.000000, Low: 22906.600000, High: 23080.000000, Volume: 730.915961, Time: time.UnixMicro(1674806400000000)},
		Kline{Amount: 0, Count: 0, Open: 22966.000000, Close: 22934.400000, Low: 22841.000000, High: 23142.400000, Volume: 1405.035772, Time: time.UnixMicro(1674820800000000)},
		Kline{Amount: 0, Count: 0, Open: 22933.300000, Close: 23252.800000, Low: 22925.700000, High: 23283.800000, Volume: 1331.053175, Time: time.UnixMicro(1674835200000000)},
		Kline{Amount: 0, Count: 0, Open: 23252.900000, Close: 23074.800000, Low: 22955.900000, High: 23500.000000, Volume: 1608.608884, Time: time.UnixMicro(1674849600000000)},
		Kline{Amount: 0, Count: 0, Open: 23074.800000, Close: 23120.200000, Low: 23045.200000, High: 23185.000000, Volume: 517.243128, Time: time.UnixMicro(1674864000000000)},
		Kline{Amount: 0, Count: 0, Open: 23120.200000, Close: 22991.900000, Low: 22954.600000, High: 23132.700000, Volume: 559.763061, Time: time.UnixMicro(1674878400000000)},
		Kline{Amount: 0, Count: 0, Open: 22992.000000, Close: 22981.300000, Low: 22896.100000, High: 23017.100000, Volume: 680.605849, Time: time.UnixMicro(1674892800000000)},
		Kline{Amount: 0, Count: 0, Open: 22981.400000, Close: 23013.300000, Low: 22879.400000, High: 23041.200000, Volume: 752.836650, Time: time.UnixMicro(1674907200000000)},
		Kline{Amount: 0, Count: 0, Open: 23013.800000, Close: 23038.900000, Low: 22976.500000, High: 23050.000000, Volume: 207.560843, Time: time.UnixMicro(1674921600000000)},
		Kline{Amount: 0, Count: 0, Open: 23038.900000, Close: 23025.100000, Low: 22943.100000, High: 23044.200000, Volume: 303.342248, Time: time.UnixMicro(1674936000000000)},
		Kline{Amount: 0, Count: 0, Open: 23025.100000, Close: 23192.400000, Low: 22970.700000, High: 23495.000000, Volume: 2188.822373, Time: time.UnixMicro(1674950400000000)},
		Kline{Amount: 0, Count: 0, Open: 23193.800000, Close: 23202.900000, Low: 23150.300000, High: 23265.500000, Volume: 321.227889, Time: time.UnixMicro(1674964800000000)},
		Kline{Amount: 0, Count: 0, Open: 23202.900000, Close: 23429.500000, Low: 23158.200000, High: 23619.700000, Volume: 2145.342777, Time: time.UnixMicro(1674979200000000)},
		Kline{Amount: 0, Count: 0, Open: 23430.800000, Close: 23527.000000, Low: 23353.300000, High: 23652.400000, Volume: 1484.452760, Time: time.UnixMicro(1674993600000000)},
		Kline{Amount: 0, Count: 0, Open: 23526.900000, Close: 23898.500000, Low: 23475.100000, High: 23963.000000, Volume: 1564.099976, Time: time.UnixMicro(1675008000000000)},
		Kline{Amount: 0, Count: 0, Open: 23898.500000, Close: 23740.000000, Low: 23605.000000, High: 23911.500000, Volume: 1243.333510, Time: time.UnixMicro(1675022400000000)},
		Kline{Amount: 0, Count: 0, Open: 23743.700000, Close: 23698.100000, Low: 23569.500000, High: 23798.100000, Volume: 827.969550, Time: time.UnixMicro(1675036800000000)},
		Kline{Amount: 0, Count: 0, Open: 23698.200000, Close: 23647.400000, Low: 23601.400000, High: 23761.900000, Volume: 420.604417, Time: time.UnixMicro(1675051200000000)},
		Kline{Amount: 0, Count: 0, Open: 23643.900000, Close: 23079.300000, Low: 23016.800000, High: 23672.100000, Volume: 3437.634193, Time: time.UnixMicro(1675065600000000)},
		Kline{Amount: 0, Count: 0, Open: 23077.200000, Close: 23176.200000, Low: 22959.700000, High: 23297.000000, Volume: 1657.926239, Time: time.UnixMicro(1675080000000000)},
		Kline{Amount: 0, Count: 0, Open: 23176.200000, Close: 22780.000000, Low: 22614.200000, High: 23244.700000, Volume: 2544.053277, Time: time.UnixMicro(1675094400000000)},
		Kline{Amount: 0, Count: 0, Open: 22780.000000, Close: 22825.100000, Low: 22483.000000, High: 22835.600000, Volume: 1755.580524, Time: time.UnixMicro(1675108800000000)},
		Kline{Amount: 0, Count: 0, Open: 22825.000000, Close: 22832.900000, Low: 22718.700000, High: 22912.400000, Volume: 633.708506, Time: time.UnixMicro(1675123200000000)},
		Kline{Amount: 0, Count: 0, Open: 22833.000000, Close: 22977.900000, Low: 22745.100000, High: 22987.800000, Volume: 1131.392235, Time: time.UnixMicro(1675137600000000)},
		Kline{Amount: 0, Count: 0, Open: 22977.800000, Close: 22865.800000, Low: 22812.500000, High: 22993.200000, Volume: 719.578293, Time: time.UnixMicro(1675152000000000)},
		Kline{Amount: 0, Count: 0, Open: 22866.900000, Close: 23122.300000, Low: 22853.800000, High: 23209.900000, Volume: 1180.349586, Time: time.UnixMicro(1675166400000000)},
		Kline{Amount: 0, Count: 0, Open: 23123.500000, Close: 23169.000000, Low: 23085.300000, High: 23223.000000, Volume: 490.314062, Time: time.UnixMicro(1675180800000000)},
		Kline{Amount: 0, Count: 0, Open: 23167.200000, Close: 23127.600000, Low: 22819.000000, High: 23347.700000, Volume: 1571.137467, Time: time.UnixMicro(1675195200000000)},
		Kline{Amount: 0, Count: 0, Open: 23130.000000, Close: 23127.500000, Low: 22991.100000, High: 23184.500000, Volume: 552.968063, Time: time.UnixMicro(1675209600000000)},
		Kline{Amount: 0, Count: 0, Open: 23127.500000, Close: 23077.800000, Low: 23022.900000, High: 23167.400000, Volume: 448.876819, Time: time.UnixMicro(1675224000000000)},
		Kline{Amount: 0, Count: 0, Open: 23077.700000, Close: 23074.900000, Low: 22918.900000, High: 23097.300000, Volume: 667.328335, Time: time.UnixMicro(1675238400000000)},
		Kline{Amount: 0, Count: 0, Open: 23075.000000, Close: 22984.100000, Low: 22943.000000, High: 23141.100000, Volume: 818.842938, Time: time.UnixMicro(1675252800000000)},
		Kline{Amount: 0, Count: 0, Open: 22984.000000, Close: 23381.000000, Low: 22756.000000, High: 23494.100000, Volume: 3389.657301, Time: time.UnixMicro(1675267200000000)},
		Kline{Amount: 0, Count: 0, Open: 23381.100000, Close: 23728.700000, Low: 23323.400000, High: 23850.000000, Volume: 1976.596984, Time: time.UnixMicro(1675281600000000)},
		Kline{Amount: 0, Count: 0, Open: 23729.000000, Close: 23857.100000, Low: 23694.500000, High: 24252.700000, Volume: 2595.212304, Time: time.UnixMicro(1675296000000000)},
		Kline{Amount: 0, Count: 0, Open: 23857.200000, Close: 23781.200000, Low: 23664.400000, High: 23919.000000, Volume: 1539.367299, Time: time.UnixMicro(1675310400000000)},
		Kline{Amount: 0, Count: 0, Open: 23781.200000, Close: 23822.600000, Low: 23750.100000, High: 23853.000000, Volume: 702.597381, Time: time.UnixMicro(1675324800000000)},
		Kline{Amount: 0, Count: 0, Open: 23822.600000, Close: 23838.200000, Low: 23525.000000, High: 23941.300000, Volume: 2588.767207, Time: time.UnixMicro(1675339200000000)},
		Kline{Amount: 0, Count: 0, Open: 23843.900000, Close: 23807.100000, Low: 23729.200000, High: 24148.100000, Volume: 1407.448972, Time: time.UnixMicro(1675353600000000)},
		Kline{Amount: 0, Count: 0, Open: 23807.200000, Close: 23487.800000, Low: 23362.500000, High: 23925.000000, Volume: 1557.548035, Time: time.UnixMicro(1675368000000000)},
		Kline{Amount: 0, Count: 0, Open: 23490.200000, Close: 23535.000000, Low: 23414.200000, High: 23587.400000, Volume: 809.822879, Time: time.UnixMicro(1675382400000000)},
		Kline{Amount: 0, Count: 0, Open: 23535.100000, Close: 23453.000000, Low: 23438.100000, High: 23563.300000, Volume: 500.209409, Time: time.UnixMicro(1675396800000000)},
		Kline{Amount: 0, Count: 0, Open: 23452.900000, Close: 23533.400000, Low: 23325.700000, High: 23537.400000, Volume: 1082.878852, Time: time.UnixMicro(1675411200000000)},
		Kline{Amount: 0, Count: 0, Open: 23533.400000, Close: 23601.000000, Low: 23221.000000, High: 23708.300000, Volume: 2752.953183, Time: time.UnixMicro(1675425600000000)},
		Kline{Amount: 0, Count: 0, Open: 23603.600000, Close: 23310.100000, Low: 23303.200000, High: 23722.500000, Volume: 871.656235, Time: time.UnixMicro(1675440000000000)},
		Kline{Amount: 0, Count: 0, Open: 23310.000000, Close: 23431.200000, Low: 23211.000000, High: 23469.100000, Volume: 499.590768, Time: time.UnixMicro(1675454400000000)},
		Kline{Amount: 0, Count: 0, Open: 23431.300000, Close: 23340.000000, Low: 23333.400000, High: 23460.000000, Volume: 384.002928, Time: time.UnixMicro(1675468800000000)},
		Kline{Amount: 0, Count: 0, Open: 23340.000000, Close: 23330.400000, Low: 23255.400000, High: 23357.500000, Volume: 302.294210, Time: time.UnixMicro(1675483200000000)},
		Kline{Amount: 0, Count: 0, Open: 23327.700000, Close: 23363.400000, Low: 23292.400000, High: 23372.000000, Volume: 240.753030, Time: time.UnixMicro(1675497600000000)},
		Kline{Amount: 0, Count: 0, Open: 23363.400000, Close: 23420.400000, Low: 23342.800000, High: 23583.600000, Volume: 1049.700525, Time: time.UnixMicro(1675512000000000)},
		Kline{Amount: 0, Count: 0, Open: 23420.400000, Close: 23427.200000, Low: 23362.500000, High: 23469.100000, Volume: 261.184295, Time: time.UnixMicro(1675526400000000)},
		Kline{Amount: 0, Count: 0, Open: 23427.300000, Close: 23324.700000, Low: 23263.000000, High: 23446.000000, Volume: 362.188326, Time: time.UnixMicro(1675540800000000)},
		Kline{Amount: 0, Count: 0, Open: 23324.800000, Close: 23344.200000, Low: 23229.300000, High: 23370.200000, Volume: 598.348976, Time: time.UnixMicro(1675555200000000)},
		Kline{Amount: 0, Count: 0, Open: 23344.200000, Close: 23380.800000, Low: 23329.800000, High: 23432.100000, Volume: 238.257783, Time: time.UnixMicro(1675569600000000)},
		Kline{Amount: 0, Count: 0, Open: 23380.800000, Close: 23357.800000, Low: 23346.200000, High: 23415.900000, Volume: 195.518501, Time: time.UnixMicro(1675584000000000)},
		Kline{Amount: 0, Count: 0, Open: 23357.800000, Close: 23099.200000, Low: 23008.000000, High: 23373.700000, Volume: 2130.859367, Time: time.UnixMicro(1675598400000000)},
		Kline{Amount: 0, Count: 0, Open: 23099.100000, Close: 22891.600000, Low: 22754.600000, High: 23159.600000, Volume: 1941.697571, Time: time.UnixMicro(1675612800000000)},
		Kline{Amount: 0, Count: 0, Open: 22891.700000, Close: 22933.100000, Low: 22800.400000, High: 23032.100000, Volume: 628.829584, Time: time.UnixMicro(1675627200000000)},
		Kline{Amount: 0, Count: 0, Open: 22933.300000, Close: 22895.500000, Low: 22831.400000, High: 23093.000000, Volume: 509.337244, Time: time.UnixMicro(1675641600000000)},
		Kline{Amount: 0, Count: 0, Open: 22895.500000, Close: 22878.200000, Low: 22617.600000, High: 22925.000000, Volume: 1356.150888, Time: time.UnixMicro(1675656000000000)},
		Kline{Amount: 0, Count: 0, Open: 22878.200000, Close: 22877.400000, Low: 22736.000000, High: 22934.300000, Volume: 699.980046, Time: time.UnixMicro(1675670400000000)},
		Kline{Amount: 0, Count: 0, Open: 22879.400000, Close: 23042.100000, Low: 22745.200000, High: 23053.000000, Volume: 957.542510, Time: time.UnixMicro(1675684800000000)},
		Kline{Amount: 0, Count: 0, Open: 23044.000000, Close: 23031.000000, Low: 22927.800000, High: 23158.500000, Volume: 732.323355, Time: time.UnixMicro(1675699200000000)},
		Kline{Amount: 0, Count: 0, Open: 23028.400000, Close: 22761.700000, Low: 22638.000000, High: 23029.800000, Volume: 669.278319, Time: time.UnixMicro(1675713600000000)},
		Kline{Amount: 0, Count: 0, Open: 22761.700000, Close: 22878.000000, Low: 22747.100000, High: 22897.900000, Volume: 506.996808, Time: time.UnixMicro(1675728000000000)},
		Kline{Amount: 0, Count: 0, Open: 22878.600000, Close: 22913.100000, Low: 22865.100000, High: 22984.800000, Volume: 491.389134, Time: time.UnixMicro(1675742400000000)},
		Kline{Amount: 0, Count: 0, Open: 22911.300000, Close: 22977.900000, Low: 22868.200000, High: 23059.700000, Volume: 580.171426, Time: time.UnixMicro(1675756800000000)},
		Kline{Amount: 0, Count: 0, Open: 22977.900000, Close: 22909.100000, Low: 22880.900000, High: 23037.700000, Volume: 639.400031, Time: time.UnixMicro(1675771200000000)},
		Kline{Amount: 0, Count: 0, Open: 22909.100000, Close: 23095.800000, Low: 22768.200000, High: 23349.900000, Volume: 2476.107046, Time: time.UnixMicro(1675785600000000)},
		Kline{Amount: 0, Count: 0, Open: 23093.300000, Close: 23246.100000, Low: 23029.200000, High: 23300.000000, Volume: 783.649378, Time: time.UnixMicro(1675800000000000)},
		Kline{Amount: 0, Count: 0, Open: 23246.100000, Close: 23261.000000, Low: 23231.000000, High: 23457.100000, Volume: 860.550584, Time: time.UnixMicro(1675814400000000)},
		Kline{Amount: 0, Count: 0, Open: 23257.300000, Close: 23234.200000, Low: 23187.000000, High: 23285.300000, Volume: 499.162291, Time: time.UnixMicro(1675828800000000)},
		Kline{Amount: 0, Count: 0, Open: 23234.300000, Close: 23159.900000, Low: 23136.300000, High: 23237.600000, Volume: 626.515718, Time: time.UnixMicro(1675843200000000)},
		Kline{Amount: 0, Count: 0, Open: 23159.900000, Close: 22869.900000, Low: 22860.000000, High: 23199.100000, Volume: 1180.249218, Time: time.UnixMicro(1675857600000000)},
		Kline{Amount: 0, Count: 0, Open: 22866.300000, Close: 22874.700000, Low: 22673.800000, High: 23049.900000, Volume: 1423.549021, Time: time.UnixMicro(1675872000000000)},
		Kline{Amount: 0, Count: 0, Open: 22874.700000, Close: 22960.000000, Low: 22781.600000, High: 22977.500000, Volume: 318.759674, Time: time.UnixMicro(1675886400000000)},
		Kline{Amount: 0, Count: 0, Open: 22960.000000, Close: 22516.600000, Low: 22428.000000, High: 23009.900000, Volume: 2592.540129, Time: time.UnixMicro(1675900800000000)},
		Kline{Amount: 0, Count: 0, Open: 22516.600000, Close: 22677.900000, Low: 22358.000000, High: 22731.900000, Volume: 1618.181981, Time: time.UnixMicro(1675915200000000)},
		Kline{Amount: 0, Count: 0, Open: 22677.800000, Close: 22686.700000, Low: 22657.800000, High: 22769.600000, Volume: 674.106411, Time: time.UnixMicro(1675929600000000)},
		Kline{Amount: 0, Count: 0, Open: 22686.600000, Close: 22614.700000, Low: 22561.900000, High: 22840.900000, Volume: 1120.327945, Time: time.UnixMicro(1675944000000000)},
		Kline{Amount: 0, Count: 0, Open: 22613.800000, Close: 22032.100000, Low: 21877.000000, High: 22621.400000, Volume: 3955.633141, Time: time.UnixMicro(1675958400000000)},
		Kline{Amount: 0, Count: 0, Open: 22032.000000, Close: 21795.000000, Low: 21678.000000, High: 22083.800000, Volume: 3710.490267, Time: time.UnixMicro(1675972800000000)},
		Kline{Amount: 0, Count: 0, Open: 21795.000000, Close: 21821.000000, Low: 21758.400000, High: 21933.300000, Volume: 787.350718, Time: time.UnixMicro(1675987200000000)},
		Kline{Amount: 0, Count: 0, Open: 21821.100000, Close: 21913.000000, Low: 21625.000000, High: 21925.000000, Volume: 1173.618615, Time: time.UnixMicro(1676001600000000)},
		Kline{Amount: 0, Count: 0, Open: 21912.200000, Close: 21731.700000, Low: 21715.800000, High: 21927.700000, Volume: 522.088735, Time: time.UnixMicro(1676016000000000)},
		Kline{Amount: 0, Count: 0, Open: 21731.700000, Close: 21655.700000, Low: 21518.000000, High: 21899.700000, Volume: 1712.570027, Time: time.UnixMicro(1676030400000000)},
		Kline{Amount: 0, Count: 0, Open: 21651.900000, Close: 21772.900000, Low: 21561.600000, High: 21793.900000, Volume: 923.679752, Time: time.UnixMicro(1676044800000000)},
		Kline{Amount: 0, Count: 0, Open: 21773.300000, Close: 21628.300000, Low: 21440.000000, High: 21790.400000, Volume: 926.791409, Time: time.UnixMicro(1676059200000000)},
		Kline{Amount: 0, Count: 0, Open: 21628.400000, Close: 21680.100000, Low: 21600.100000, High: 21721.500000, Volume: 499.554349, Time: time.UnixMicro(1676073600000000)},
		Kline{Amount: 0, Count: 0, Open: 21680.100000, Close: 21691.500000, Low: 21646.300000, High: 21709.000000, Volume: 199.332596, Time: time.UnixMicro(1676088000000000)},
		Kline{Amount: 0, Count: 0, Open: 21691.500000, Close: 21695.800000, Low: 21656.600000, High: 21701.000000, Volume: 198.916739, Time: time.UnixMicro(1676102400000000)},
		Kline{Amount: 0, Count: 0, Open: 21695.900000, Close: 21733.800000, Low: 21665.100000, High: 21769.800000, Volume: 425.352664, Time: time.UnixMicro(1676116800000000)},
		Kline{Amount: 0, Count: 0, Open: 21733.900000, Close: 21658.200000, Low: 21620.400000, High: 21735.900000, Volume: 452.259527, Time: time.UnixMicro(1676131200000000)},
		Kline{Amount: 0, Count: 0, Open: 21658.200000, Close: 21856.200000, Low: 21646.400000, High: 21898.500000, Volume: 484.998852, Time: time.UnixMicro(1676145600000000)},
		Kline{Amount: 0, Count: 0, Open: 21856.200000, Close: 21777.400000, Low: 21766.900000, High: 21881.000000, Volume: 424.350392, Time: time.UnixMicro(1676160000000000)},
		Kline{Amount: 0, Count: 0, Open: 21778.500000, Close: 21798.100000, Low: 21760.000000, High: 21811.300000, Volume: 188.989818, Time: time.UnixMicro(1676174400000000)},
		Kline{Amount: 0, Count: 0, Open: 21798.100000, Close: 21891.500000, Low: 21798.100000, High: 21950.000000, Volume: 514.000319, Time: time.UnixMicro(1676188800000000)},
		Kline{Amount: 0, Count: 0, Open: 21891.500000, Close: 21943.100000, Low: 21764.000000, High: 22000.000000, Volume: 1322.877442, Time: time.UnixMicro(1676203200000000)},
		Kline{Amount: 0, Count: 0, Open: 21943.100000, Close: 21992.500000, Low: 21914.200000, High: 22083.500000, Volume: 1065.798890, Time: time.UnixMicro(1676217600000000)},
		Kline{Amount: 0, Count: 0, Open: 21992.600000, Close: 21784.500000, Low: 21628.500000, High: 22023.600000, Volume: 986.099396, Time: time.UnixMicro(1676232000000000)},
		Kline{Amount: 0, Count: 0, Open: 21784.700000, Close: 21809.900000, Low: 21628.300000, High: 21867.700000, Volume: 527.134928, Time: time.UnixMicro(1676246400000000)},
		Kline{Amount: 0, Count: 0, Open: 21809.900000, Close: 21870.900000, Low: 21799.000000, High: 21895.200000, Volume: 611.989639, Time: time.UnixMicro(1676260800000000)},
		Kline{Amount: 0, Count: 0, Open: 21871.000000, Close: 21605.100000, Low: 21421.900000, High: 21876.200000, Volume: 2053.533755, Time: time.UnixMicro(1676275200000000)},
		Kline{Amount: 0, Count: 0, Open: 21607.900000, Close: 21566.400000, Low: 21561.000000, High: 21708.300000, Volume: 1568.816704, Time: time.UnixMicro(1676289600000000)},
		Kline{Amount: 0, Count: 0, Open: 21567.300000, Close: 21608.300000, Low: 21356.500000, High: 21675.800000, Volume: 2076.458136, Time: time.UnixMicro(1676304000000000)},
		Kline{Amount: 0, Count: 0, Open: 21608.400000, Close: 21776.800000, Low: 21566.300000, High: 21847.700000, Volume: 781.783073, Time: time.UnixMicro(1676318400000000)},
		Kline{Amount: 0, Count: 0, Open: 21776.900000, Close: 21698.000000, Low: 21681.700000, High: 21814.500000, Volume: 362.786747, Time: time.UnixMicro(1676332800000000)},
		Kline{Amount: 0, Count: 0, Open: 21697.900000, Close: 21765.400000, Low: 21688.000000, High: 21784.600000, Volume: 446.900904, Time: time.UnixMicro(1676347200000000)},
		Kline{Amount: 0, Count: 0, Open: 21765.500000, Close: 21805.500000, Low: 21673.700000, High: 21871.200000, Volume: 796.707597, Time: time.UnixMicro(1676361600000000)},
		Kline{Amount: 0, Count: 0, Open: 21805.400000, Close: 22039.600000, Low: 21540.000000, High: 22319.800000, Volume: 5239.370792, Time: time.UnixMicro(1676376000000000)},
		Kline{Amount: 0, Count: 0, Open: 22041.700000, Close: 22208.000000, Low: 21919.300000, High: 22242.500000, Volume: 1153.809457, Time: time.UnixMicro(1676390400000000)},
		Kline{Amount: 0, Count: 0, Open: 22211.100000, Close: 22197.900000, Low: 22125.100000, High: 22278.300000, Volume: 562.776532, Time: time.UnixMicro(1676404800000000)},
		Kline{Amount: 0, Count: 0, Open: 22197.600000, Close: 22099.000000, Low: 22048.000000, High: 22197.600000, Volume: 450.406338, Time: time.UnixMicro(1676419200000000)},
		Kline{Amount: 0, Count: 0, Open: 22101.000000, Close: 22107.800000, Low: 22075.000000, High: 22171.800000, Volume: 422.617762, Time: time.UnixMicro(1676433600000000)},
		Kline{Amount: 0, Count: 0, Open: 22107.800000, Close: 22441.900000, Low: 22089.100000, High: 22487.800000, Volume: 1234.271529, Time: time.UnixMicro(1676448000000000)},
		Kline{Amount: 0, Count: 0, Open: 22439.400000, Close: 22811.300000, Low: 22439.400000, High: 22908.100000, Volume: 2809.073270, Time: time.UnixMicro(1676462400000000)},
		Kline{Amount: 0, Count: 0, Open: 22807.700000, Close: 23319.100000, Low: 22746.100000, High: 23360.000000, Volume: 1721.693526, Time: time.UnixMicro(1676476800000000)},
		Kline{Amount: 0, Count: 0, Open: 23324.400000, Close: 24328.800000, Low: 23299.100000, High: 24381.400000, Volume: 6313.782705, Time: time.UnixMicro(1676491200000000)},
		Kline{Amount: 0, Count: 0, Open: 24329.200000, Close: 24733.200000, Low: 24285.100000, High: 24911.900000, Volume: 3033.596495, Time: time.UnixMicro(1676505600000000)},
		Kline{Amount: 0, Count: 0, Open: 24733.300000, Close: 24618.700000, Low: 24443.700000, High: 24756.100000, Volume: 1726.581430, Time: time.UnixMicro(1676520000000000)},
		Kline{Amount: 0, Count: 0, Open: 24618.600000, Close: 24572.100000, Low: 24509.700000, High: 24662.400000, Volume: 727.868950, Time: time.UnixMicro(1676534400000000)},
		Kline{Amount: 0, Count: 0, Open: 24572.100000, Close: 25070.900000, Low: 24286.800000, High: 25080.000000, Volume: 3045.342751, Time: time.UnixMicro(1676548800000000)},
		Kline{Amount: 0, Count: 0, Open: 25071.400000, Close: 24876.400000, Low: 24610.000000, High: 25249.900000, Volume: 2873.587352, Time: time.UnixMicro(1676563200000000)},
		Kline{Amount: 0, Count: 0, Open: 24871.100000, Close: 23512.200000, Low: 23502.400000, High: 24877.700000, Volume: 4262.851645, Time: time.UnixMicro(1676577600000000)},
		Kline{Amount: 0, Count: 0, Open: 23509.800000, Close: 23835.100000, Low: 23337.100000, High: 23912.800000, Volume: 2329.457814, Time: time.UnixMicro(1676592000000000)},
		Kline{Amount: 0, Count: 0, Open: 23835.100000, Close: 23634.100000, Low: 23612.300000, High: 23881.600000, Volume: 1084.603415, Time: time.UnixMicro(1676606400000000)},
		Kline{Amount: 0, Count: 0, Open: 23632.800000, Close: 23790.500000, Low: 23528.100000, High: 23861.400000, Volume: 1474.319557, Time: time.UnixMicro(1676620800000000)},
		Kline{Amount: 0, Count: 0, Open: 23790.200000, Close: 24130.200000, Low: 23708.700000, High: 24281.000000, Volume: 3632.209849, Time: time.UnixMicro(1676635200000000)},
		Kline{Amount: 0, Count: 0, Open: 24126.200000, Close: 24525.900000, Low: 23911.700000, High: 24556.800000, Volume: 3163.126176, Time: time.UnixMicro(1676649600000000)},
		Kline{Amount: 0, Count: 0, Open: 24525.900000, Close: 24571.800000, Low: 24055.000000, High: 25025.000000, Volume: 3612.170704, Time: time.UnixMicro(1676664000000000)},
		Kline{Amount: 0, Count: 0, Open: 24570.700000, Close: 24608.500000, Low: 24508.800000, High: 24784.100000, Volume: 948.642810, Time: time.UnixMicro(1676678400000000)},
		Kline{Amount: 0, Count: 0, Open: 24608.600000, Close: 24518.800000, Low: 24427.600000, High: 24673.600000, Volume: 626.310817, Time: time.UnixMicro(1676692800000000)},
		Kline{Amount: 0, Count: 0, Open: 24519.000000, Close: 24521.200000, Low: 24438.100000, High: 24604.200000, Volume: 465.418426, Time: time.UnixMicro(1676707200000000)},
		Kline{Amount: 0, Count: 0, Open: 24521.200000, Close: 24666.000000, Low: 24486.200000, High: 24729.900000, Volume: 824.653964, Time: time.UnixMicro(1676721600000000)},
		Kline{Amount: 0, Count: 0, Open: 24666.000000, Close: 24605.300000, Low: 24568.200000, High: 24875.000000, Volume: 774.843468, Time: time.UnixMicro(1676736000000000)},
		Kline{Amount: 0, Count: 0, Open: 24607.100000, Close: 24633.700000, Low: 24503.600000, High: 24664.400000, Volume: 309.302282, Time: time.UnixMicro(1676750400000000)},
		Kline{Amount: 0, Count: 0, Open: 24633.700000, Close: 24678.600000, Low: 24624.100000, High: 24765.600000, Volume: 316.633254, Time: time.UnixMicro(1676764800000000)},
		Kline{Amount: 0, Count: 0, Open: 24678.600000, Close: 24579.600000, Low: 24564.700000, High: 24846.100000, Volume: 878.907001, Time: time.UnixMicro(1676779200000000)},
		Kline{Amount: 0, Count: 0, Open: 24579.700000, Close: 24671.100000, Low: 24559.400000, High: 24708.800000, Volume: 672.644512, Time: time.UnixMicro(1676793600000000)},
		Kline{Amount: 0, Count: 0, Open: 24671.200000, Close: 24922.900000, Low: 24658.000000, High: 25046.000000, Volume: 1423.476419, Time: time.UnixMicro(1676808000000000)},
		Kline{Amount: 0, Count: 0, Open: 24923.400000, Close: 24514.300000, Low: 24279.800000, High: 25187.000000, Volume: 3908.202144, Time: time.UnixMicro(1676822400000000)},
		Kline{Amount: 0, Count: 0, Open: 24519.300000, Close: 24276.000000, Low: 24182.100000, High: 24575.800000, Volume: 949.882288, Time: time.UnixMicro(1676836800000000)},
		Kline{Amount: 0, Count: 0, Open: 24276.100000, Close: 24416.000000, Low: 23836.500000, High: 24467.100000, Volume: 2186.493919, Time: time.UnixMicro(1676851200000000)},
		Kline{Amount: 0, Count: 0, Open: 24416.100000, Close: 24502.200000, Low: 24402.100000, High: 24553.500000, Volume: 750.407714, Time: time.UnixMicro(1676865600000000)},
		Kline{Amount: 0, Count: 0, Open: 24502.300000, Close: 24884.700000, Low: 24377.200000, High: 24985.700000, Volume: 2610.053188, Time: time.UnixMicro(1676880000000000)},
		Kline{Amount: 0, Count: 0, Open: 24887.000000, Close: 24944.500000, Low: 24725.200000, High: 25129.200000, Volume: 2425.465859, Time: time.UnixMicro(1676894400000000)},
		Kline{Amount: 0, Count: 0, Open: 24944.600000, Close: 24874.800000, Low: 24620.700000, High: 24953.700000, Volume: 1213.872269, Time: time.UnixMicro(1676908800000000)},
		Kline{Amount: 0, Count: 0, Open: 24877.900000, Close: 24843.100000, Low: 24650.700000, High: 24889.500000, Volume: 419.770922, Time: time.UnixMicro(1676923200000000)},
		Kline{Amount: 0, Count: 0, Open: 24843.000000, Close: 24925.900000, Low: 24783.500000, High: 25100.000000, Volume: 985.070100, Time: time.UnixMicro(1676937600000000)},
		Kline{Amount: 0, Count: 0, Open: 24926.000000, Close: 25000.100000, Low: 24852.600000, High: 25049.300000, Volume: 728.740732, Time: time.UnixMicro(1676952000000000)},
		Kline{Amount: 0, Count: 0, Open: 25004.100000, Close: 24569.900000, Low: 24547.000000, High: 25320.000000, Volume: 3687.906851, Time: time.UnixMicro(1676966400000000)},
		Kline{Amount: 0, Count: 0, Open: 24567.900000, Close: 24557.600000, Low: 24390.500000, High: 24781.600000, Volume: 1922.104850, Time: time.UnixMicro(1676980800000000)},
		Kline{Amount: 0, Count: 0, Open: 24557.600000, Close: 24605.300000, Low: 24279.800000, High: 24750.600000, Volume: 1814.063400, Time: time.UnixMicro(1676995200000000)},
		Kline{Amount: 0, Count: 0, Open: 24605.400000, Close: 24450.000000, Low: 24150.200000, High: 24636.000000, Volume: 2023.276306, Time: time.UnixMicro(1677009600000000)},
		Kline{Amount: 0, Count: 0, Open: 24450.000000, Close: 24200.800000, Low: 23866.000000, High: 24471.000000, Volume: 1879.021126, Time: time.UnixMicro(1677024000000000)},
		Kline{Amount: 0, Count: 0, Open: 24200.900000, Close: 24061.400000, Low: 23931.000000, High: 24220.100000, Volume: 1267.395432, Time: time.UnixMicro(1677038400000000)},
		Kline{Amount: 0, Count: 0, Open: 24060.200000, Close: 24162.100000, Low: 23866.000000, High: 24206.400000, Volume: 1165.860679, Time: time.UnixMicro(1677052800000000)},
		Kline{Amount: 0, Count: 0, Open: 24162.100000, Close: 23716.800000, Low: 23636.100000, High: 24230.800000, Volume: 2186.723024, Time: time.UnixMicro(1677067200000000)},
		Kline{Amount: 0, Count: 0, Open: 23714.000000, Close: 23791.900000, Low: 23579.500000, High: 23968.800000, Volume: 1539.232872, Time: time.UnixMicro(1677081600000000)},
		Kline{Amount: 0, Count: 0, Open: 23794.600000, Close: 24190.000000, Low: 23738.400000, High: 24213.200000, Volume: 950.673406, Time: time.UnixMicro(1677096000000000)},
		Kline{Amount: 0, Count: 0, Open: 24190.000000, Close: 24455.600000, Low: 24121.900000, High: 24602.900000, Volume: 1972.250061, Time: time.UnixMicro(1677110400000000)},
		Kline{Amount: 0, Count: 0, Open: 24455.700000, Close: 24379.200000, Low: 24314.200000, High: 24530.700000, Volume: 874.795503, Time: time.UnixMicro(1677124800000000)},
		Kline{Amount: 0, Count: 0, Open: 24379.000000, Close: 23780.100000, Low: 23650.800000, High: 24493.200000, Volume: 2943.321088, Time: time.UnixMicro(1677139200000000)},
		Kline{Amount: 0, Count: 0, Open: 23780.100000, Close: 23950.200000, Low: 23605.300000, High: 24226.700000, Volume: 3515.807821, Time: time.UnixMicro(1677153600000000)},
		Kline{Amount: 0, Count: 0, Open: 23950.100000, Close: 23970.600000, Low: 23725.000000, High: 24050.000000, Volume: 1746.100574, Time: time.UnixMicro(1677168000000000)},
		Kline{Amount: 0, Count: 0, Open: 23972.800000, Close: 23938.800000, Low: 23757.800000, High: 24066.000000, Volume: 552.964601, Time: time.UnixMicro(1677182400000000)},
		Kline{Amount: 0, Count: 0, Open: 23938.900000, Close: 23954.400000, Low: 23889.800000, High: 24131.000000, Volume: 749.432337, Time: time.UnixMicro(1677196800000000)},
		Kline{Amount: 0, Count: 0, Open: 23954.500000, Close: 23832.400000, Low: 23773.400000, High: 23980.000000, Volume: 1051.549970, Time: time.UnixMicro(1677211200000000)},
		Kline{Amount: 0, Count: 0, Open: 23832.400000, Close: 23886.900000, Low: 23757.300000, High: 23970.800000, Volume: 1006.955383, Time: time.UnixMicro(1677225600000000)},
		Kline{Amount: 0, Count: 0, Open: 23887.000000, Close: 23349.000000, Low: 23333.900000, High: 24025.200000, Volume: 3378.756219, Time: time.UnixMicro(1677240000000000)},
		Kline{Amount: 0, Count: 0, Open: 23350.100000, Close: 23178.500000, Low: 22836.100000, High: 23353.800000, Volume: 4368.354125, Time: time.UnixMicro(1677254400000000)},
		Kline{Amount: 0, Count: 0, Open: 23181.800000, Close: 23190.100000, Low: 22979.500000, High: 23330.200000, Volume: 1435.504893, Time: time.UnixMicro(1677268800000000)},
		Kline{Amount: 0, Count: 0, Open: 23190.200000, Close: 23074.200000, Low: 23050.000000, High: 23217.800000, Volume: 726.572161, Time: time.UnixMicro(1677283200000000)},
		Kline{Amount: 0, Count: 0, Open: 23077.700000, Close: 23121.800000, Low: 23023.500000, High: 23124.000000, Volume: 562.692063, Time: time.UnixMicro(1677297600000000)},
		Kline{Amount: 0, Count: 0, Open: 23121.800000, Close: 22983.800000, Low: 22914.500000, High: 23152.400000, Volume: 1078.723240, Time: time.UnixMicro(1677312000000000)},
		Kline{Amount: 0, Count: 0, Open: 22983.800000, Close: 23000.500000, Low: 22863.700000, High: 23081.200000, Volume: 796.996321, Time: time.UnixMicro(1677326400000000)},
		Kline{Amount: 0, Count: 0, Open: 23002.800000, Close: 22990.200000, Low: 22917.000000, High: 23083.100000, Volume: 584.068717, Time: time.UnixMicro(1677340800000000)},
		Kline{Amount: 0, Count: 0, Open: 22986.300000, Close: 23157.200000, Low: 22729.400000, High: 23185.100000, Volume: 1269.976075, Time: time.UnixMicro(1677355200000000)},
		Kline{Amount: 0, Count: 0, Open: 23157.300000, Close: 23214.900000, Low: 23060.700000, High: 23249.600000, Volume: 513.589338, Time: time.UnixMicro(1677369600000000)},
		Kline{Amount: 0, Count: 0, Open: 23215.000000, Close: 23153.500000, Low: 23094.000000, High: 23223.500000, Volume: 358.561414, Time: time.UnixMicro(1677384000000000)},
		Kline{Amount: 0, Count: 0, Open: 23153.500000, Close: 23252.900000, Low: 23116.200000, High: 23288.900000, Volume: 629.459132, Time: time.UnixMicro(1677398400000000)},
		Kline{Amount: 0, Count: 0, Open: 23253.000000, Close: 23247.900000, Low: 23130.000000, High: 23332.900000, Volume: 790.441829, Time: time.UnixMicro(1677412800000000)},
		Kline{Amount: 0, Count: 0, Open: 23248.000000, Close: 23501.300000, Low: 23132.300000, High: 23542.700000, Volume: 1122.983034, Time: time.UnixMicro(1677427200000000)},
		Kline{Amount: 0, Count: 0, Open: 23503.600000, Close: 23557.400000, Low: 23329.900000, High: 23686.700000, Volume: 1395.405576, Time: time.UnixMicro(1677441600000000)},
		Kline{Amount: 0, Count: 0, Open: 23552.700000, Close: 23551.300000, Low: 23436.000000, High: 23638.700000, Volume: 724.751523, Time: time.UnixMicro(1677456000000000)},
		Kline{Amount: 0, Count: 0, Open: 23550.700000, Close: 23421.700000, Low: 23350.300000, High: 23556.300000, Volume: 996.923123, Time: time.UnixMicro(1677470400000000)},
		Kline{Amount: 0, Count: 0, Open: 23419.600000, Close: 23429.800000, Low: 23339.700000, High: 23448.800000, Volume: 634.977495, Time: time.UnixMicro(1677484800000000)},
		Kline{Amount: 0, Count: 0, Open: 23429.900000, Close: 23555.900000, Low: 23380.500000, High: 23891.200000, Volume: 2353.581843, Time: time.UnixMicro(1677499200000000)},
		Kline{Amount: 0, Count: 0, Open: 23556.000000, Close: 23280.100000, Low: 23162.200000, High: 23614.600000, Volume: 2577.709329, Time: time.UnixMicro(1677513600000000)},
		Kline{Amount: 0, Count: 0, Open: 23280.000000, Close: 23491.200000, Low: 23101.000000, High: 23576.000000, Volume: 1068.298178, Time: time.UnixMicro(1677528000000000)},
		Kline{Amount: 0, Count: 0, Open: 23491.200000, Close: 23463.200000, Low: 23354.000000, High: 23548.900000, Volume: 444.140370, Time: time.UnixMicro(1677542400000000)},
		Kline{Amount: 0, Count: 0, Open: 23464.400000, Close: 23235.600000, Low: 23216.800000, High: 23468.100000, Volume: 718.322697, Time: time.UnixMicro(1677556800000000)},
		Kline{Amount: 0, Count: 0, Open: 23236.400000, Close: 23396.400000, Low: 23208.000000, High: 23417.800000, Volume: 460.114024, Time: time.UnixMicro(1677571200000000)},
		Kline{Amount: 0, Count: 0, Open: 23396.400000, Close: 23521.000000, Low: 23321.600000, High: 23575.000000, Volume: 1126.019799, Time: time.UnixMicro(1677585600000000)},
		Kline{Amount: 0, Count: 0, Open: 23521.100000, Close: 23259.400000, Low: 23202.300000, High: 23600.000000, Volume: 801.030465, Time: time.UnixMicro(1677600000000000)},
		Kline{Amount: 0, Count: 0, Open: 23261.800000, Close: 23136.000000, Low: 23025.500000, High: 23349.300000, Volume: 1149.432422, Time: time.UnixMicro(1677614400000000)},
		Kline{Amount: 0, Count: 0, Open: 23143.100000, Close: 23444.100000, Low: 23009.700000, High: 23500.000000, Volume: 1106.856113, Time: time.UnixMicro(1677628800000000)},
		Kline{Amount: 0, Count: 0, Open: 23444.000000, Close: 23718.000000, Low: 23427.700000, High: 23846.000000, Volume: 1570.069948, Time: time.UnixMicro(1677643200000000)},
		Kline{Amount: 0, Count: 0, Open: 23717.900000, Close: 23739.900000, Low: 23667.600000, High: 24009.000000, Volume: 1737.977982, Time: time.UnixMicro(1677657600000000)},
		Kline{Amount: 0, Count: 0, Open: 23739.800000, Close: 23712.400000, Low: 23570.300000, High: 23882.600000, Volume: 1753.173271, Time: time.UnixMicro(1677672000000000)},
		Kline{Amount: 0, Count: 0, Open: 23712.300000, Close: 23357.500000, Low: 23331.700000, High: 23751.200000, Volume: 1345.620863, Time: time.UnixMicro(1677686400000000)},
		Kline{Amount: 0, Count: 0, Open: 23354.300000, Close: 23632.400000, Low: 23300.300000, High: 23672.000000, Volume: 842.689064, Time: time.UnixMicro(1677700800000000)},
		Kline{Amount: 0, Count: 0, Open: 23632.000000, Close: 23497.900000, Low: 23431.300000, High: 23792.700000, Volume: 870.539582, Time: time.UnixMicro(1677715200000000)},
		Kline{Amount: 0, Count: 0, Open: 23498.800000, Close: 23391.100000, Low: 23341.800000, High: 23544.500000, Volume: 740.536879, Time: time.UnixMicro(1677729600000000)},
		Kline{Amount: 0, Count: 0, Open: 23391.900000, Close: 23418.100000, Low: 23342.000000, High: 23465.500000, Volume: 693.653301, Time: time.UnixMicro(1677744000000000)},
		Kline{Amount: 0, Count: 0, Open: 23418.100000, Close: 23361.100000, Low: 23181.500000, High: 23434.200000, Volume: 1787.887711, Time: time.UnixMicro(1677758400000000)},
		Kline{Amount: 0, Count: 0, Open: 23361.000000, Close: 23469.400000, Low: 23233.900000, High: 23500.000000, Volume: 697.915317, Time: time.UnixMicro(1677772800000000)},
		Kline{Amount: 0, Count: 0, Open: 23464.600000, Close: 23466.000000, Low: 23392.200000, High: 23566.700000, Volume: 720.649788, Time: time.UnixMicro(1677787200000000)},
		Kline{Amount: 0, Count: 0, Open: 23465.900000, Close: 22319.000000, Low: 21898.000000, High: 23472.900000, Volume: 9781.062403, Time: time.UnixMicro(1677801600000000)},
		Kline{Amount: 0, Count: 0, Open: 22319.100000, Close: 22361.600000, Low: 22252.100000, High: 22407.000000, Volume: 1948.048887, Time: time.UnixMicro(1677816000000000)},
		Kline{Amount: 0, Count: 0, Open: 22361.600000, Close: 22343.800000, Low: 22305.000000, High: 22474.500000, Volume: 1261.137374, Time: time.UnixMicro(1677830400000000)},
		Kline{Amount: 0, Count: 0, Open: 22343.900000, Close: 22336.400000, Low: 22156.300000, High: 22503.700000, Volume: 2066.883204, Time: time.UnixMicro(1677844800000000)},
		Kline{Amount: 0, Count: 0, Open: 22336.500000, Close: 22358.600000, Low: 22230.200000, High: 22446.100000, Volume: 1079.669849, Time: time.UnixMicro(1677859200000000)},
		Kline{Amount: 0, Count: 0, Open: 22357.900000, Close: 22354.500000, Low: 22150.000000, High: 22377.300000, Volume: 1019.273111, Time: time.UnixMicro(1677873600000000)},
		Kline{Amount: 0, Count: 0, Open: 22354.400000, Close: 22350.600000, Low: 22325.300000, High: 22391.000000, Volume: 433.582833, Time: time.UnixMicro(1677888000000000)},
		// Kline{Amount: 0, Count: 0, Open: 22350.600000, Close: 22356.400000, Low: 22271.600000, High: 22403.800000, Volume: 453.488856, Time: time.UnixMicro(1677902400000000)},
		// Kline{Amount: 0, Count: 0, Open: 22356.500000, Close: 22349.300000, Low: 22328.500000, High: 22370.000000, Volume: 304.625859, Time: time.UnixMicro(1677916800000000)},
		// Kline{Amount: 0, Count: 0, Open: 22349.400000, Close: 22315.500000, Low: 22310.000000, High: 22396.200000, Volume: 373.105335, Time: time.UnixMicro(1677931200000000)},
		// Kline{Amount: 0, Count: 0, Open: 22315.600000, Close: 22245.400000, Low: 22229.100000, High: 22346.000000, Volume: 430.941665, Time: time.UnixMicro(1677945600000000)},
		// Kline{Amount: 0, Count: 0, Open: 22245.500000, Close: 22349.700000, Low: 22164.500000, High: 22364.700000, Volume: 962.913394, Time: time.UnixMicro(1677960000000000)},
		// Kline{Amount: 0, Count: 0, Open: 22349.600000, Close: 22386.500000, Low: 22185.500000, High: 22666.400000, Volume: 1918.833969, Time: time.UnixMicro(1677974400000000)},
		// Kline{Amount: 0, Count: 0, Open: 22386.700000, Close: 22376.700000, Low: 22368.400000, High: 22469.400000, Volume: 481.965016, Time: time.UnixMicro(1677988800000000)},
		// Kline{Amount: 0, Count: 0, Open: 22374.900000, Close: 22391.800000, Low: 22318.000000, High: 22445.700000, Volume: 466.642386, Time: time.UnixMicro(1678003200000000)},
		// Kline{Amount: 0, Count: 0, Open: 22390.100000, Close: 22438.800000, Low: 22390.000000, High: 22490.000000, Volume: 493.194504, Time: time.UnixMicro(1678017600000000)},
		// Kline{Amount: 0, Count: 0, Open: 22438.900000, Close: 22404.000000, Low: 22378.800000, High: 22456.200000, Volume: 213.038792, Time: time.UnixMicro(1678032000000000)},
		// Kline{Amount: 0, Count: 0, Open: 22403.100000, Close: 22427.600000, Low: 22329.400000, High: 22566.900000, Volume: 632.846393, Time: time.UnixMicro(1678046400000000)},
		// Kline{Amount: 0, Count: 0, Open: 22427.600000, Close: 22387.500000, Low: 22261.300000, High: 22509.100000, Volume: 879.753666, Time: time.UnixMicro(1678060800000000)},
		// Kline{Amount: 0, Count: 0, Open: 22387.600000, Close: 22388.700000, Low: 22315.000000, High: 22435.100000, Volume: 859.340008, Time: time.UnixMicro(1678075200000000)},
		// Kline{Amount: 0, Count: 0, Open: 22388.700000, Close: 22381.300000, Low: 22358.000000, High: 22448.200000, Volume: 883.143321, Time: time.UnixMicro(1678089600000000)},
		// Kline{Amount: 0, Count: 0, Open: 22381.300000, Close: 22560.000000, Low: 22360.100000, High: 22598.000000, Volume: 1511.742557, Time: time.UnixMicro(1678104000000000)},
		// Kline{Amount: 0, Count: 0, Open: 22560.000000, Close: 22395.500000, Low: 22357.200000, High: 22584.900000, Volume: 1416.135258, Time: time.UnixMicro(1678118400000000)},
		// Kline{Amount: 0, Count: 0, Open: 22393.200000, Close: 22409.200000, Low: 22317.000000, High: 22449.300000, Volume: 406.658820, Time: time.UnixMicro(1678132800000000)},
		// Kline{Amount: 0, Count: 0, Open: 22409.200000, Close: 22456.800000, Low: 22379.400000, High: 22553.300000, Volume: 662.668312, Time: time.UnixMicro(1678147200000000)},
		// Kline{Amount: 0, Count: 0, Open: 22456.900000, Close: 22441.600000, Low: 22412.500000, High: 22473.900000, Volume: 428.776190, Time: time.UnixMicro(1678161600000000)},
	}
}

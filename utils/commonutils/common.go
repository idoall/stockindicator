package commonutils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/shopspring/decimal"
)

var (
	zeroValueUnix = time.Unix(0, 0)
	// ErrDateUnset is an error for start end check calculations
	ErrDateUnset = errors.New("date unset")
	// ErrStartAfterEnd is an error for start end check calculations
	ErrStartAfterEnd = errors.New("start date after end date")
	// ErrStartEqualsEnd is an error for start end check calculations
	ErrStartEqualsEnd = errors.New("start date equals end date")
	// ErrStartAfterTimeNow is an error for start end check calculations
	ErrStartAfterTimeNow = errors.New("start date is after current time")
)

// FloatFromString format
func FloatFromString(raw interface{}) (float64, error) {
	str, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("unable to parse, value not string: %T", raw)
	}
	flt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse, value not string: %T", raw)
	}
	return flt, nil
}

// FormatDecimalFloat64 获取精度计算后的数量
//
//		 example
//			  FormatDecimalFloat64(123.456, -2)   // output: 123.46
//	          FormatDecimalFloat64(123.456, 1) // output: 120
//			  FormatDecimalFloat64(-500,-2)   // output: -500
//			  FormatDecimalFloat64(-500,0)   // output: -500
//			  FormatDecimalFloat64(1.1001, -2) // output: 1.1
//			  FormatDecimalFloat64(1.454, -1) // output: 1.5
//			  FormatDecimalFloat64(1.454, 0) // output: 1
func FormatDecimalFloat64(value float64, exp int32) float64 {
	var val, _ = decimal.NewFromFloatWithExponent(value, exp).Float64()
	return val
}

// 三元表达式简单应用，为了便于程序易读
//
//	example
//		If(a>b,a,b).(int)
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// ReadFile reads a file and returns read data as byte array.
func ReadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// JSONDecode decodes JSON data into a structure
func JSONDecode(data []byte, to interface{}) error {
	if !strings.Contains(reflect.ValueOf(to).Type().String(), "*") {
		return errors.New("json decode error - memory address not supplied")
	}
	return sonic.Unmarshal(data, to)
}

// StartEndTimeCheck provides some basic checks which occur
// frequently in the codebase
func StartEndTimeCheck(start, end time.Time) error {
	if start.IsZero() || start.Equal(zeroValueUnix) {
		return fmt.Errorf("start %w", ErrDateUnset)
	}
	if end.IsZero() || end.Equal(zeroValueUnix) {
		return fmt.Errorf("end %w", ErrDateUnset)
	}
	if start.After(end) {
		return ErrStartAfterEnd
	}
	if start.Equal(end) {
		return ErrStartEqualsEnd
	}
	if start.After(time.Now()) {
		return ErrStartAfterTimeNow
	}

	return nil
}

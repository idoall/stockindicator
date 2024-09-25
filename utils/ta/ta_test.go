package ta_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/idoall/stockindicator/helpertools"
	"github.com/idoall/stockindicator/utils/ta"
)

// go test -v ./utils/ta -run ^TestSmaT$
func TestSmaT(t *testing.T) {
	values := helpertools.Slice2Chan([]float64{
		22.27, 22.19, 22.08, 22.17, 22.18, 22.13, 22.23, 22.43, 22.24,
		22.29, 22.15, 22.39, 22.38, 22.61, 23.36, 24.05, 23.75, 23.83,
		23.95, 23.63, 23.82, 23.87, 23.65, 23.19, 23.10, 23.33, 22.68,
		23.10, 22.40, 22.17,
	})

	expected := []float64{
		2.23, 4.45, 6.65, 8.87, 11.09, 13.3, 15.52, 17.77, 19.99, 22.22, 22.21, 22.23, 22.26, 22.3, 22.42, 22.61, 22.77, 22.91, 23.08, 23.21, 23.38, 23.53, 23.65, 23.71, 23.68, 23.61, 23.5, 23.43, 23.28, 23.13,
	}

	sma := ta.SmaT[float64](10, values)

	var result = helpertools.Chan2Slice(helpertools.RoundDigits(sma, 2))

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("result %v expected %v", result, expected)
	}
}

// go test -v ./utils/ta -run ^TestSD$
func TestSD(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{
			name:     "happy path",
			values:   []float64{1, 2, 3, 4, 5},
			expected: math.Sqrt(2.0), // Standard deviation for {1, 2, 3, 4, 5}
		},
		{
			name:     "single value",
			values:   []float64{5},
			expected: 0.0, // Standard deviation for a single value
		},
		{
			name:     "two values",
			values:   []float64{1, 3},
			expected: 1.0, // Standard deviation for {1, 3}
		},
		{
			name:     "all same value",
			values:   []float64{4, 4, 4, 4},
			expected: 0.0, // All same values should have a standard deviation of 0
		},
		{
			name:     "negative values",
			values:   []float64{-1, -2, -3, -4, -5},
			expected: math.Sqrt(2.0), // Standard deviation for {-1, -2, -3, -4, -5}
		},
		{
			name:     "mixed values",
			values:   []float64{-1, 0, 1},
			expected: 0.816496580927726, // Standard deviation for {-1, 0, 1}
		},
		{
			name:     "empty slice",
			values:   []float64{},
			expected: 0.0, // Standard deviation for an empty slice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ta.SD(tt.values)
			if !almostEqual(result, tt.expected) {
				t.Errorf("Expected %v but got %v", tt.expected, result)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

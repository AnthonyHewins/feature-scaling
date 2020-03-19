package fe

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	DELTA = 0.001
)

func TestZScore(t *testing.T) {
	sharedTests(t, ZScore)

	// Test against the dataset [0, 1, 2, ..., 200]
	// where the stddev is already known (and can generate the value for each element)
	x        := make([]float64, 201)
	expected := make([]float64, 201)
	sigma := 0.017191624218032074

	for i, _ := range expected {
		x[i] = float64(i)
		expected[i] = sigma * float64(i - 100)
	}

	ZScore(x)
	assert.InDeltaSlice(t, x, expected, DELTA)
}

func TestMeanNormalization(t *testing.T) {
	sharedTests(t, MeanNormalization)

	// Test against the dataset [0, 1, 2, ..., 200]
	// where the mean/range is already known (and can generate the value for each element)
	x        := make([]float64, 201)
	expected := make([]float64, 201)

	for i, _ := range expected {
		x[i] = float64(i)
		expected[i] = (float64(i) - 101.5) / 200
	}

	MeanNormalization(x)
	assert.InDeltaSlice(t, x, expected, DELTA)
}

func sharedTests(t *testing.T, fn func([]float64)) {
	x := []float64{}
	fn(x)
	assert.Equal(t, len(x), 0)

	x = []float64{1}
	fn(x)
	assert.Equal(t, x[0], float64(1))

	x = []float64{1,1,1}
	fn(x)
	assert.Equal(t, x, []float64{0,0,0})
}

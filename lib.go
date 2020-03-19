package fe

import (
	"sync"
	"gonum.org/v1/gonum/stat"
)

const (
	ROUTINE_DIVISOR = 200
)

func MeanNormalization(x []float64) {
	n := len(x)
	if n <= 1 { return }

	ss, min, max := float64(0), x[0], x[0]
	for _, val := range x {
		ss += val

		if val > max {
			max = val
		} else if val < min {
			min = val
		}
	}

	// Make no effort for concurrency if the values are the same; just 0 it out
	if min == max {
		for i, _ := range x { x[i] = 0 }
		return
	}

	r := max - min
	mu := ss / float64(n)
	divideAndConquer(n, x, func(slice []float64) {
		for i, val := range slice {
			slice[i] = (val - mu) / r
		}
	})
}

func ZScore(x []float64) {
	n := len(x)
	if n <= 1 { return }

	mu, sigma := stat.MeanStdDev(x, nil)

	// Make no effort for concurrency if the values are the same; just 0 it out
	if sigma == 0 {
		for i, _ := range x { x[i] = 0 }
		return
	}

	divideAndConquer(n, x, func(slice []float64) {
		for i, val := range slice {
			slice[i] = (val - mu) / sigma
		}
	})
}

func divideAndConquer(n int, x []float64, fn func(slice []float64)) {
	// No concurrency if the array size doesn't merit it
	if n <= ROUTINE_DIVISOR {
		fn(x)
		return
	}

	var w sync.WaitGroup
	routines := n / ROUTINE_DIVISOR
	w.Add(routines)

	var concurrent = func(slice []float64) {
		defer w.Done()
		fn(slice)
	}

	for i := 1; i <= routines; i++ {
		go concurrent(x[ROUTINE_DIVISOR * (i-1) : ROUTINE_DIVISOR * i])
	}

	if n % ROUTINE_DIVISOR != 0 {
		fn(x[routines * ROUTINE_DIVISOR:n])
	}

	w.Wait()
}

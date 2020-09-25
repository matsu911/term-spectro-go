package main

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func toFloat64(y []float32) []float64 {
	length := len(y)
	v := make([]float64, length, length)
	for i := 0; i < length; i++ {
		v[i] = float64(y[i])
	}
	return v
}

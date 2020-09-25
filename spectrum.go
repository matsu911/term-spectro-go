package main

import (
	"math"

	"github.com/mjibson/go-dsp/fft"
)

func stft(y []float64, win_length int) [][]float64 {
	length := len(y) / win_length
	ret := make([][]float64, 0, length)
	for offset := 0; offset < length; offset++ {
		s := fft.FFTReal(y[offset*win_length : (offset+1)*win_length])
		// fmt.Printf("%v\n", s)
		// fmt.Printf("%v\n", len(s))
		tmp := make([]float64, len(s), len(s))
		for i := 0; i < len(s); i++ {
			tmp[i] = math.Abs(real(s[i]))
		}
		ret = append(ret, tmp)
	}
	return ret
}

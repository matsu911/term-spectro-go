package main

import (
	"log"
	"math"
	"os"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
)

func ReadWav(fileName string) [][]float64 {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	w, err := wav.New(f)
	v, err := w.ReadFloats(w.Samples)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	// length := 128
	window := 64
	length := w.Samples / window
	ret := make([][]float64, 0, length)
	for offset := 0; offset < length; offset++ {
		// fmt.Printf("%v\n", window)
		vv := make([]float64, window, window)
		for i := 0; i < window; i++ {
			vv[i] = float64(v[offset*window+i])
		}
		// fmt.Printf("%v\n", vv)
		s := fft.FFTReal(vv)
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

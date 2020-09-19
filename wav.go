package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
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

func PlayAudio(audioPath string) (beep.StreamSeekCloser, beep.Format) {
	f, err := os.Open(audioPath)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	fmt.Println("before")
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	speaker.Play(volume)
	return streamer, format
}

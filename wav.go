package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/mjibson/go-dsp/wav"
)

func ReadWav(fileName string) [][]float64 {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	w, err := wav.New(f)
	y, err := w.ReadFloats(w.Samples)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	win_length := 64
	return stft(toFloat64(y), win_length)
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

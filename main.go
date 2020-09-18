package main

import (
	"fmt"
	"os"
	"time"
	"log"
	"math"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
)

var style = tcell.StyleDefault
var quit chan struct{}

func main() {
	spec := ReadWav("./test.wav")

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.EnableMouse()
	s.Clear()

	quit = make(chan struct{})
	go pollEvents(s)

	s.Show()

	go func() {
		for {
			drawScreen(s, spec)
			time.Sleep(time.Millisecond * 10)
		}
	}()

	<-quit
	s.Fini()
}

type viewport struct {
	x0, x1, y0, y1 float64
}

func NewViewport(x0, x1, y0, y1 float64) viewport {
	vp := viewport{}
	vp.x0 = x0
	vp.x1 = x1
	vp.y0 = y0
	vp.y1 = y1

	return vp
}

var vp = NewViewport(-2.0, 1.0, -1.0, 1.0)

func zoom(s tcell.Screen, direction, x, y int) {
	//	w, h := s.Size()

	factorx := (vp.x0 - vp.x1) / 10.0
	factory := (vp.y0 - vp.y1) / 10.0

	if direction == 1 {
		vp.x0 -= factorx
		vp.x1 += factorx
		vp.y0 -= factory
		vp.y1 += factory
	} else {
		vp.x0 += factorx
		vp.x1 -= factorx
		vp.y0 += factory
		vp.y1 -= factory
	}
}

func MinMax(array []float64) (float64, float64) {
	max := array[0]
	min := array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func Normalize(v []float64) []float64 {
	sum2 := 0.0
	for i := 0; i < len(v); i++ {
		sum2 += v[i] * v[i]
	}
	ret := make([]float64, len(v), len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i] / math.Sqrt(sum2)
	}
	return ret
}

func drawScreen(s tcell.Screen, spec [][]float64) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	st := tcell.StyleDefault
	// const gl = '▄'
	const gl = ' '

	s.Fill(gl, st.Background(tcell.ColorRed))

	// st = st.Background(tcell.ColorBlack)

	for i := 0; i < len(spec); i++ {
		for j := 0; j < len(spec[i]); j++ {
			// c := tcell.PaletteColor(int(spec[i][j] * 0xff))
			x := int32(spec[i][j] * 0xff)
			c := tcell.NewRGBColor(x, x, x)
			s.SetCell(i*2 + 0, j, st.Background(c), gl)
			s.SetCell(i*2 + 1, j, st.Background(c), gl)
		}
	}

	s.Show()
}

func pollEvents(s tcell.Screen) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
				close(quit)
				return
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'z':
					zoom(s, 1, 1, 1)
				case 'x':
					zoom(s, 0, 1, 1)
				case 'q':
					close(quit)
					return
				}
				//s.Sync()
			case tcell.KeyUp:
				step := (vp.y0 - vp.y1) / 10
				vp.y0 += step
				vp.y1 += step
			case tcell.KeyDown:
				step := (vp.y0 - vp.y1) / 10
				vp.y0 -= step
				vp.y1 -= step
			case tcell.KeyLeft:
				step := (vp.x0 - vp.x1) / 10
				vp.x0 += step
				vp.x1 += step
			case tcell.KeyRight:
				step := (vp.x0 - vp.x1) / 10
				vp.x0 -= step
				vp.x1 -= step
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			/*if button&tcell.WheelUp != 0 {
				bstr += " WheelUp"
			}*/
			// Only buttons, not wheel events
			button &= tcell.ButtonMask(0xff)
			switch ev.Buttons() {
			case tcell.Button1:
				zoom(s, 1, x, y)
			case tcell.Button2:
				zoom(s, 0, x, y)
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

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
	length := 128
	window := 64
	ret := make([][]float64, 0, length)
	for offset := 1000; offset < 1000 + length; offset++ {
		// fmt.Printf("%v\n", window)
		vv := make([]float64, window, window)
		for i := 0; i < window; i++ {
			vv[i] = float64(v[offset * window + i])
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

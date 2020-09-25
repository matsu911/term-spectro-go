package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var app = &views.Application{}
var window = &mainWindow{}

func (m *model) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if len(m.specs) == 0 {
		return ' ', tcell.StyleDefault, nil, 1
	}
	if y >= len(m.specs[0]) {
		return ' ', tcell.StyleDefault, nil, 1
	}
	cc := int32(m.specs[x+m.offsetx][y] * 0xff)
	c := tcell.NewRGBColor(cc, cc, cc)
	return ' ', tcell.StyleDefault.Background(c), nil, 1
}

func main() {
	// streamer, _ := PlayAudio("test.mp3")
	// defer streamer.Close()

	fmt.Println("after PlayAudio")
	specs := ReadWav("./test.wav")

	window.model = &model{
		endx:    60,
		endy:    15,
		offsetx: 0,
		offsety: 0,
		specs:   specs,
	}

	title := views.NewTextBar()
	title.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorTeal).
		Foreground(tcell.ColorWhite))
	title.SetCenter("Spectrogram Viewer", tcell.StyleDefault)
	title.SetRight("v1.0", tcell.StyleDefault)

	window.keybar = views.NewSimpleStyledText()
	window.keybar.RegisterStyle('N', tcell.StyleDefault.
		Background(tcell.ColorSilver).
		Foreground(tcell.ColorBlack))
	window.keybar.RegisterStyle('A', tcell.StyleDefault.
		Background(tcell.ColorSilver).
		Foreground(tcell.ColorRed))

	window.status = views.NewSimpleStyledTextBar()
	window.status.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorYellow))
	window.status.RegisterLeftStyle('N', tcell.StyleDefault.
		Background(tcell.ColorYellow).
		Foreground(tcell.ColorBlack))

	// window.status.SetLeft("My status is here.")
	// window.status.SetRight("%UCellView%N demo!")
	// window.status.SetCenter("Cen%ST%Ner")

	window.main = views.NewCellView()
	window.main.SetModel(window.model)
	window.main.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlack))

	window.SetMenu(window.keybar)
	window.SetTitle(title)
	window.SetContent(window.main)
	window.SetStatus(window.status)

	window.updateKeys()

	app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	app.SetRootWidget(window)
	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}

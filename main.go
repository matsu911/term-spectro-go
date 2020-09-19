package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var app = &views.Application{}
var window = &mainWindow{}

type model struct {
	x         int
	y         int
	endx      int
	endy      int
	hide      bool
	enab      bool
	open      bool
	audioPath string
	loc       string
	specs     [][]float64
}

func (m *model) GetBounds() (int, int) {
	return m.endx, m.endy
}

func (m *model) MoveCursor(offx, offy int) {
	m.x += offx
	m.y += offy
	m.limitCursor()
}

func (m *model) limitCursor() {
	if m.x < 0 {
		m.x = 0
	}
	if m.x > m.endx-1 {
		m.x = m.endx - 1
	}
	if m.y < 0 {
		m.y = 0
	}
	if m.y > m.endy-1 {
		m.y = m.endy - 1
	}
	m.loc = fmt.Sprintf("Cursor is %d,%d", m.x, m.y)
}

func (m *model) GetCursor() (int, int, bool, bool) {
	return m.x, m.y, m.enab, !m.hide
}

func (m *model) SetCursor(x int, y int) {
	m.x = x
	m.y = y

	m.limitCursor()
}

func (m *model) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if y >= len(m.specs[0]) {
		return ' ', tcell.StyleDefault, nil, 1
	}
	cc := int32(m.specs[x+1000][y] * 0xff)
	c := tcell.NewRGBColor(cc, cc, cc)
	return ' ', tcell.StyleDefault.Background(c), nil, 1
}

type mainWindow struct {
	main   *views.CellView
	keybar *views.SimpleStyledText
	status *views.SimpleStyledTextBar
	model  *model

	views.Panel
}

func (a *mainWindow) HandleEvent(ev tcell.Event) bool {

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			app.Refresh()
			return true
		case tcell.KeyEsc:
			a.model.open = false
			a.updateKeys()
			return true
		case tcell.KeyRune:
			if a.model.open {
				a.model.audioPath += string(ev.Rune())
				a.updateKeys()
			} else {
				switch ev.Rune() {
				case 'Q', 'q':
					app.Quit()
					return true
				case 'O', 'o':
					a.model.open = true
					a.updateKeys()
					return true
				case 'S', 's':
					a.model.hide = false
					a.updateKeys()
					return true
				case 'H', 'h':
					a.model.hide = true
					a.updateKeys()
					return true
				case 'E', 'e':
					a.model.enab = true
					a.updateKeys()
					return true
				case 'D', 'd':
					a.model.enab = false
					a.updateKeys()
					return true
				}
			}
		}
	}
	return a.Panel.HandleEvent(ev)
}

func (a *mainWindow) Draw() {
	// a.status.SetLeft(a.model.loc)
	if a.model.open {
		a.status.SetLeft("file: " + a.model.audioPath)
	} else {
		a.status.SetLeft("")
	}
	a.Panel.Draw()
}

func (a *mainWindow) updateKeys() {
	m := a.model
	w := "[%AQ%N] Quit"
	if !m.open {
		w += "  [%AO%N] Open file"
	}
	if !m.enab {
		w += "  [%AE%N] Enable cursor"
	} else {
		w += "  [%AD%N] Disable cursor"
		if !m.hide {
			w += "  [%AH%N] Hide cursor"
		} else {
			w += "  [%AS%N] Show cursor"
		}
	}
	a.keybar.SetMarkup(w)
	app.Update()
}

func main() {
	specs := ReadWav("./test.wav")

	window.model = &model{endx: 60, endy: 15, specs: specs}

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

	window.status.SetLeft("My status is here.")
	window.status.SetRight("%UCellView%N demo!")
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

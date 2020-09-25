package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type model struct {
	x         int
	y         int
	endx      int
	endy      int
	offsetx   int
	offsety   int
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
	return m.x, m.y, m.enab, true
}

func (m *model) SetCursor(x int, y int) {
	m.x = x
	m.y = y

	m.limitCursor()
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
				case 'E', 'e':
					a.model.enab = true
					a.updateKeys()
					return true
				case 'D', 'd':
					a.model.enab = false
					a.updateKeys()
					return true
				case 'H', 'h':
					a.model.offsetx = Max(0, a.model.offsetx-1)
					a.updateKeys()
					return true
				case 'J', 'j':
					a.model.offsetx += 100
					a.updateKeys()
					return true
				case 'K', 'k':
					a.model.offsetx = Max(0, a.model.offsetx-100)
					a.updateKeys()
					return true
				case 'L', 'l':
					a.model.offsetx += 1
					a.updateKeys()
					return true
				}
			}
		}
	}
	return a.Panel.HandleEvent(ev)
}

func (a *mainWindow) Draw() {
	if a.model.open {
		a.status.SetLeft("file: " + a.model.audioPath)
	} else {
		a.status.SetLeft("")
	}
	a.status.SetRight(fmt.Sprintf("%d/%d", a.model.offsetx, len(a.model.specs)))
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
	}
	a.keybar.SetMarkup(w)
	app.Update()
}

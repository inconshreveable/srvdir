package main

import (
	"fmt"
	termbox "github.com/nsf/termbox-go"
)

type Route struct {
	addr string
	path string
}

func ui(routes []Route) {
	termbox.Init()
	defer termbox.Close()

	width, height := termbox.Size()
	a := NewArea(0, 0, width, height)

	a.Clear()
	a.APrintf(termbox.ColorCyan | termbox.AttrBold, 0, 0, "srvdir")
	a.Printf(7, 0, "by")
	a.APrintf(termbox.ColorCyan | termbox.AttrBold, 10, 0, "@inconshreveable")
	quitMessage := "(Ctl+C to quit)"
	a.Printf(width-len(quitMessage), 0, quitMessage)

	a.APrintf(termbox.AttrBold, 0, 2, "Serving")
	for i, r := range routes {
		a.Printf(0, i+3, "%s -> %s", r.addr, r.path)
	}

	termbox.Flush()

	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyCtrlC:
				return
			}
		case termbox.EventError:
			return
		case termbox.EventResize:
			// XXX
		}
	}
}

const (
	fgColor = termbox.ColorWhite
	bgColor = termbox.ColorDefault
)

type area struct {
	// top-left corner
	x, y int

	// size of the area
	w, h int

	// default colors
	fgColor, bgColor termbox.Attribute
}

func NewArea(x, y, w, h int) *area {
	return &area{x, y, w, h, fgColor, bgColor}
}

func (a *area) Clear() {
	for i := 0; i < a.w; i++ {
		for j := 0; j < a.h; j++ {
			termbox.SetCell(a.x+i, a.y+j, ' ', a.fgColor, a.bgColor)
		}
	}
}

func (a *area) APrintf(fg termbox.Attribute, x, y int, arg0 string, args ...interface{}) {
	s := fmt.Sprintf(arg0, args...)
	for i, ch := range s {
		termbox.SetCell(a.x+x+i, a.y+y, ch, fg, bgColor)
	}
}

func (a *area) Printf(x, y int, arg0 string, args ...interface{}) {
	a.APrintf(a.fgColor, x, y, arg0, args...)
}

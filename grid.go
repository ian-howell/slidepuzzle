package main

import (
	"github.com/ian-howell/gocurse/curses"
	"github.com/ian-howell/gocurse/panels"
)

type Grid struct {
	Size    int
	cells   []int
	window  *curses.Window
	numrows int
	numcols int
	bgColor int32
	fgColor int32
	panel   *panels.Panel
}

func NewGrid(size int) *Grid {
	g := &Grid{
		Size:    size,
		cells:   make([]int, size*size),
		numrows: (4 * size) + 1,
		numcols: (4 * size) + 1,
		bgColor: curses.Color_pair(1),
		fgColor: curses.Color_pair(0),
	}
	for i := 0; i < (size*size)-1; i++ {
		if i+1 < 10 {
			g.cells[i] = int('0') + i + 1
		} else {
			g.cells[i] = int('A') + i - 9
		}
	}
	g.cells[(size*size)-1] = 0
	g.window, _ = curses.Newwin(g.numrows, g.numcols, 0, 0)
	g.panel = panels.NewPanel(g.window)
	panels.UpdatePanels()
	return g
}

func (g *Grid) At(row, col int) int {
	return g.cells[row*g.Size+col]
}

func (g *Grid) Draw() {
	for r := 0; r < g.numrows; r++ {
		for c := 0; c < g.numcols; c++ {
			g.printAt(r, c)
		}
	}
	panels.UpdatePanels()
	curses.DoUpdate()
}

func (g *Grid) printAt(r, c int) {
	if g.isBorder(r, c) {
		g.window.Addch(c, r, ' ', g.bgColor)
	} else if g.isNumber(r, c) {
		num := int32(g.At((r-2)/4, (c-2)/4))
		if num > 0 {
			g.window.Addch(c, r, num, g.fgColor)
		} else {
			g.window.Addch(c, r, ' ', g.fgColor)
		}
	} else {
		g.window.Addch(c, r, ' ', g.fgColor)
	}
}

func (g *Grid) isBorder(row, col int) bool {
	return row%4 == 0 || col%4 == 0
}

func (g *Grid) isNumber(row, col int) bool {
	return (row-2)%4 == 0 && (col-2)%4 == 0
}

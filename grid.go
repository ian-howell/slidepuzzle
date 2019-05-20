package main

import (
	"math/rand"
	"time"

	"github.com/ian-howell/gocurse/curses"
	"github.com/ian-howell/gocurse/panels"
)

type Grid struct {
	Size    int
	cells   []rune
	point   int
	window  *curses.Window
	numrows int
	numcols int
	bgColor int32
	fgColor int32
	panel   *panels.Panel
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

const solvedPosition = "123456789ABCDEF0"

func NewGrid() *Grid {
	g := &Grid{
		Size:    4,
		cells:   make([]rune, 16),
		numrows: 17,
		numcols: 17,
		bgColor: curses.Color_pair(1),
		fgColor: curses.Color_pair(0),
	}
	for i, val := range solvedPosition {
		g.cells[i] = val
	}
	g.point = 15
	g.Shuffle()
	g.window, _ = curses.Newwin(g.numrows, g.numcols, 0, 0)
	g.panel = panels.NewPanel(g.window)
	panels.UpdatePanels()
	return g
}

func (g *Grid) At(row, col int) rune {
	return g.cells[g.RowColToCellNumber(row, col)]
}

func (g *Grid) CellNumberToRowCol(cell int) (row, col int) {
	row = cell / g.Size
	col = cell % g.Size
	return row, col
}

func (g *Grid) RowColToCellNumber(row, col int) (cell int) {
	cell = row*g.Size + col
	return cell
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

func (g *Grid) Shuffle() {
	// Just do a whole bunch of random moves
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 5000; i++ {
		g.Move(Direction(r.Intn(4)))
	}
}

func (g *Grid) Move(d Direction) {
	row, col := g.CellNumberToRowCol(g.Point())
	newrow, newcol := row, col
	switch d {
	case Up:
		newrow++
	case Down:
		newrow--
	case Left:
		newcol++
	case Right:
		newcol--
	}
	if !g.Valid(newrow, newcol) {
		return
	}
	g.Swap(row, col, newrow, newcol)
}

func (g *Grid) Point() int {
	if g.cells[g.point] == '0' {
		return g.point
	}
	g.point = 0
	for g.cells[g.point] != '0' {
		g.point++
	}
	return g.point
}

func (g *Grid) Valid(r, c int) bool {
	return 0 <= c && c < g.Size && 0 <= r && r < g.Size
}

func (g *Grid) Swap(r1, c1, r2, c2 int) {
	cell1 := g.RowColToCellNumber(r1, c1)
	cell2 := g.RowColToCellNumber(r2, c2)
	g.cells[cell1], g.cells[cell2] = g.cells[cell2], g.cells[cell1]
}

func (g *Grid) Solved() bool {
	for i, v := range g.cells {
		if v != rune(solvedPosition[i]) {
			return false
		}
	}
	return true
}

func (g *Grid) printAt(r, c int) {
	if g.isBorder(r, c) {
		g.window.Addch(c, r, ' ', g.bgColor)
	} else if g.isNumber(r, c) {
		num := g.At((r-2)/4, (c-2)/4)
		if num != '0' {
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

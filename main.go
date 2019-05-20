package main

import (
	"fmt"
	"os"

	"github.com/ian-howell/gocurse/curses"
)

const EOF = 4

func main() {
	screen, err := curses.Initscr()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize curses: %s\n", err.Error())
		os.Exit(1)
	}
	defer curses.Endwin()

	if err := Initialize(screen); err != nil {
		fmt.Fprintf(os.Stderr, "Could not setup common curses settings: %s\n", err.Error())
		os.Exit(1)
	}

	Play(screen)
}

func Initialize(screen *curses.Window) error {
	// Cause key presses to become immediately available
	// Raw is used here to capture all signals
	if err := curses.Raw(); err != nil {
		return err
	}
	// Suppress unnecessary echoing while taking input from the user
	if err := curses.Noecho(); err != nil {
		return err
	}
	// Enables the reading of function keys like F1, F2, arrow keys etc
	if err := screen.Keypad(true); err != nil {
		return err
	}
	// Make the cursor stop blinking
	if err := curses.Curs_set(0); err != nil {
		return err
	}
	if err := curses.Start_color(); err != nil {
		return err
	}
	curses.Init_pair(0, curses.COLOR_WHITE, curses.COLOR_BLACK)
	curses.Init_pair(1, curses.COLOR_BLACK, curses.COLOR_WHITE)
	return nil
}

func Play(screen *curses.Window) error {
	puzzle := NewGrid()
	puzzle.Draw()

	keyMap := map[int]Direction{
		curses.KEY_UP:    Up,
		curses.KEY_DOWN:  Down,
		curses.KEY_LEFT:  Left,
		curses.KEY_RIGHT: Right,
	}
forloop:
	for {
		switch ch := screen.Getch(); ch {
		case curses.KEY_UP, curses.KEY_DOWN, curses.KEY_LEFT, curses.KEY_RIGHT:
			puzzle.Move(keyMap[ch])
		case EOF:
			break forloop
		}
		puzzle.Draw()
	}

	return nil
}

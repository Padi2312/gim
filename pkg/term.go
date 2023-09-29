package pkg

import (
	"github.com/buger/goterm"
)

type Term struct {
	navigation *Navigation
}

func NewTerm(navigation *Navigation) *Term {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()
	return &Term{
		navigation: navigation,
	}
}

// Show the VIM command line at the bottom of terminal
func (t *Term) ShowCommandLine(symbol string) {
	goterm.MoveCursor(1, goterm.Height()-1)
	goterm.Print(symbol)
	goterm.Flush()
}

// Show Insert Mode Info
func (t *Term) ShowInsertModeInfo() {
	goterm.MoveCursor(1, goterm.Height()-1)
	goterm.Print("---INSERT---")
	goterm.Flush()
}

// Show Insert Mode Info
func (t *Term) HideInsertModeInfo() {
	goterm.MoveCursor(1, goterm.Height()-1)
	goterm.Print("            ")
	goterm.Flush()
}

// Clears terminal and resets cursor to (1,1)
func (t *Term) Clear() {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()
}

// Do all changes to terminal ui
func (t *Term) Render() {
	x, y := t.navigation.Pos()
	x, y = t.getTermIndex(x, y)
	goterm.MoveCursor(x, y)
	goterm.Flush()
}

// Converts the change indices to terminal cursor indices
func (t *Term) getTermIndex(x int, y int) (int, int) {
	return x + 1, y + 1
}

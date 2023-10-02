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

func (t *Term) WriteFullContent(content [][]rune) {
	goterm.Clear()
	for i, row := range content {
		x, y := t.getTermIndex(t.navigation.Pos())
		goterm.MoveCursor(x, y)
		goterm.Print(string(row))
		if i != len(content)-1 {
			// Move cursor down by one line
			t.navigation.MoveDownLineBegin()
		}
	}
	// Set Cursor back to beginning
	t.navigation.cursor.SetX(0)
	t.navigation.cursor.SetY(0)
	goterm.Flush()
}

func (t *Term) GetHeight() int {
	return goterm.Height()
}

// Appends chars to the VIM command line
func (t *Term) InsertAtBottomLine(x int, symbol string) {
	goterm.MoveCursor(x, goterm.Height())
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

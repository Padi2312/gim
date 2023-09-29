package pkg

import (
	"strings"

	"github.com/buger/goterm"
)

type Line struct {
	text       []rune
	lineNumber int
}

func NewLine(lineNumber int) *Line {
	return &Line{
		text:       make([]rune, 0),
		lineNumber: lineNumber,
	}
}

func (l *Line) AddChar(x int, char rune) {
	l.text = append(l.text, char)
	x, y := l.getTermIndices(x, l.lineNumber)
	goterm.MoveCursor(x, y)
	goterm.Print(string(char))
	goterm.Flush()
}

func (l *Line) RemoveChar(x int) {
	l.text = l.text[:len(l.text)-1]
	x, y := l.getTermIndices(x, l.lineNumber)
	goterm.MoveCursor(x, y)
	goterm.Print(" ")
	// Move cursor one to left after deleting char
	goterm.MoveCursor(x, y)
	goterm.Flush()
}

func (l *Line) Clear() {
	x, y := l.getTermIndices(0, l.lineNumber)
	goterm.MoveCursor(x, y)
	goterm.Print(strings.Repeat(" ", len(l.text)))
	goterm.Flush()
}

func (l *Line) ClearFull() {
	x, y := l.getTermIndices(0, l.lineNumber)
	goterm.MoveCursor(x, y)
	goterm.Print(strings.Repeat(" ", goterm.Width()))
	goterm.Flush()
}

func (l *Line) getTermIndices(x int, y int) (int, int) {
	return x + 1, y + 1
}

type ErrorLine struct{}

func (e ErrorLine) PrintErrorLine(text string) {
	goterm.MoveCursor(1, goterm.Height())
	goterm.Print(goterm.Color(goterm.Background(text, goterm.RED), goterm.WHITE))
	goterm.Flush()
}

type TermUtils struct{}

func (t TermUtils) InsertAt(x int, y int, lineContent []rune) {
	remaining := string(lineContent[x:])

	x, y = getTermIndices(x, y)
	goterm.MoveCursor(x, y)
	goterm.Print(remaining)

}

func (t TermUtils) RemoveCharAt(x int, y int, lineContent []rune) {
	remaining := string(lineContent[x-1:]) + " "

	x, y = getTermIndices(x, y)
	goterm.MoveCursor(x-1, y)
	goterm.Print(remaining)
}

func (t TermUtils) RemoveFromTo(xFrom int, xTo int, y int) {
	xFrom, y = getTermIndices(xFrom, y)
	xTo, y = getTermIndices(xTo, y)

	goterm.MoveCursor(xFrom, y)
	cleanLine := ""
	for i := 0; i < xTo-xFrom; i++ {
		cleanLine += " "
	}
	goterm.Print(cleanLine)
	goterm.Flush()
}

func (t TermUtils) InsertNewLine(lineNumber int, content [][]rune) {
	lineNumber++
	for index, row := range content {
		goterm.MoveCursor(1, lineNumber+index)
		for i := 0; i < goterm.Width(); i++ {
			goterm.Print(" ")
		}
		goterm.MoveCursor(1, lineNumber+index)
		goterm.Print(string(row))
		goterm.Flush()
	}
}

func getTermIndices(x int, y int) (int, int) {
	return x + 1, y + 1
}

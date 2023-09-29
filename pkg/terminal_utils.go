package pkg

import "github.com/buger/goterm"

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

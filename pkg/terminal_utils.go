package pkg

import "github.com/buger/goterm"

func ClearCharAt(x int, y int) {
	goterm.MoveCursor(x, y)
	goterm.Print(" ") // Clear residual characters
	goterm.Flush()
}

func ClearLine(lineNumber int) {
	goterm.MoveCursor(1, lineNumber)
	cleanLine := ""
	for i := 0; i < goterm.Width(); i++ {
		cleanLine += " "
	}
	goterm.Print(cleanLine)
	goterm.Flush()
}

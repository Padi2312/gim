package main

import (
	"fmt"
	"os"
	"strings"
)

type GapBuffer struct {
	buffer   []rune
	gapStart int
	gapEnd   int
}

func NewGapBuffer(initialText string, cursorCol int) *GapBuffer {
	buffer := []rune(initialText)
	return &GapBuffer{
		buffer:   buffer,
		gapStart: cursorCol,
		gapEnd:   cursorCol,
	}
}

func (gb *GapBuffer) InsertText(text string) {
	for _, ch := range text {
		gb.buffer = append(gb.buffer[:gb.gapStart], append([]rune{ch}, gb.buffer[gb.gapStart:]...)...)
		gb.gapStart++
	}
}

func (gb *GapBuffer) ToString() string {
	return string(gb.buffer)
}

type TextEditor struct {
	lines     []GapBuffer
	cursorRow int
	cursorCol int
}

func NewTextEditor(initialText string) *TextEditor {
	lines := strings.Split(initialText, "\n")
	gapBuffers := make([]GapBuffer, len(lines))

	for i, line := range lines {
		gapBuffers[i] = *NewGapBuffer(line, len(line))
	}

	return &TextEditor{
		lines:     gapBuffers,
		cursorRow: 0,
		cursorCol: 0,
	}
}

func (te *TextEditor) InsertText(text string) {
	for _, ch := range text {
		if ch == '\n' {
			te.InsertNewLine()
		} else {
			te.lines[te.cursorRow].InsertText(string(ch))
			te.cursorCol++
		}
	}
}

func (te *TextEditor) InsertNewLine() {
	beforeCursor := te.lines[te.cursorRow].buffer[:te.cursorCol]
	afterCursor := te.lines[te.cursorRow].buffer[te.cursorCol:]

	// Update the current line to only include content before the cursor.
	te.lines[te.cursorRow].buffer = beforeCursor

	// Create a new line starting with content after the cursor.
	newLine := NewGapBuffer(string(afterCursor), 0)

	// Insert the new line into the lines slice.
	te.lines = append(te.lines[:te.cursorRow+1], append([]GapBuffer{*newLine}, te.lines[te.cursorRow+1:]...)...)

	// Move the cursor to the new line and set its column to 0.
	te.cursorRow++
	te.cursorCol = 0
}

func (te *TextEditor) ToString() string {
	var lines []string
	for _, gb := range te.lines {
		lines = append(lines, gb.ToString())
	}
	return strings.Join(lines, "\n")
}

func (te *TextEditor) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	content := te.ToString()
	_, err = file.WriteString(content)
	return err
}

func main() {
	editor := NewTextEditor("Hello, world!\nSecond line.")
	editor.cursorRow = 0
	editor.cursorCol = 7
	editor.InsertText(" beautiful")
	editor.InsertNewLine()

	fmt.Println(editor.ToString())

	err := editor.SaveToFile("output.txt")
	if err != nil {
		fmt.Println("Failed to save file:", err)
	} else {
		fmt.Println("File saved successfully.")
	}
}

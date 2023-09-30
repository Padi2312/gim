package pkg

import (
	"github.com/eiannone/keyboard"
)

type InsertHandler struct {
	editor *Editor
}

func NewInsertHandler(editor *Editor) *InsertHandler {
	return &InsertHandler{
		editor: editor,
	}
}

func (i *InsertHandler) Activate(insertMode rune) {
	i.editor.Mode = INSERT
	// Depending on how insert mode was entered a different cursor position is required
	i.moveCursorDependingOnEntry(insertMode)
	i.editor.Term.ShowInsertModeInfo()
}

func (i *InsertHandler) Deactivate() {
	i.editor.Mode = NORMAL
	i.editor.Term.HideInsertModeInfo()
}

func (i *InsertHandler) Handle(keyEvent KeyEvent) {
	if i.editor.Mode != INSERT {
		return
	}

	if keyEvent.Key == keyboard.KeyEsc {
		i.Deactivate()
		i.editor.Navigation.MoveLeft()
	}

	i.InsertBeforeCursor(keyEvent)
}

func (i InsertHandler) InsertBeforeCursor(keyEvent KeyEvent) {
	// We have to handle different types of insertions
	// 1. Insertion at end of text
	// 2. Insertion in middle of text
	// 3. Insertion of a linebreak (at end of line and in middle of text)
	// 4. Deletion of a char
	x, y := i.editor.Navigation.Pos()

	switch keyEvent.Key {
	case 0, keyboard.KeySpace: // No special key is pressed
		if keyEvent.Key == keyboard.KeySpace {
			keyEvent.Char = ' '
		}
		// 1. Insertion at end of thext
		// Check cursor is at end of line
		if x == i.editor.Content.LineLength(y) {
			i.editor.Content.Buffer[y] = append(i.editor.Content.Buffer[y], keyEvent.Char)
		} else {
			// 2. Insertion in middle of text
			newLine := make([]rune, len(i.editor.Content.Buffer[y]))
			copy(newLine, i.editor.Content.Buffer[y][:x])
			newLine[x] = keyEvent.Char // Adds the new char
			copy(newLine[x+1:], i.editor.Content.Buffer[y][x:])
			i.editor.Content.Buffer[y] = newLine
		}

		TermUtils{}.InsertAt(x, y, i.editor.Content.Buffer[y])
		i.editor.Navigation.MoveRight(true)

	case keyboard.KeyEnter: // 3. Insertion of a linebreak
		// Check cursor is at end of line
		if x+1 >= i.editor.Content.LineLength(y) {
			if y < i.editor.Content.TotalLines()-1 {
				i.editor.Content.Buffer = append(i.editor.Content.Buffer[:y+1], append(make([][]rune, 0), i.editor.Content.Buffer[y+1:]...)...)
				TermUtils{}.InsertNewLine(y+1, i.editor.Content.Buffer[y+1:])
			} else {
				i.editor.Content.Buffer = append(i.editor.Content.Buffer, make([]rune, 1)) // Just add a new line buffer
			}
		} else {
			remaining := i.editor.Content.Buffer[y][:x]
			newLineContent := i.editor.Content.Buffer[y][x:]
			i.editor.Content.Buffer[y] = remaining
			i.editor.Content.Buffer = append(i.editor.Content.Buffer[:y+1], append([][]rune{newLineContent}, i.editor.Content.Buffer[y+1:]...)...)
			//content.Buffer = append(content.Buffer, newLineContent)
			TermUtils{}.RemoveFromTo(x, x+len(newLineContent), y)
			TermUtils{}.InsertNewLine(y, i.editor.Content.Buffer[y:])
		}
		i.editor.Navigation.MoveDownLineBegin()

	case keyboard.KeyBackspace, keyboard.KeyBackspace2: // 4. Deletion of a char
		i.editor.Content.Buffer[y] = append(i.editor.Content.Buffer[y][:x-1], i.editor.Content.Buffer[y][x:]...)
		TermUtils{}.RemoveCharAt(x, y, i.editor.Content.Buffer[y])
		i.editor.Navigation.MoveLeft()
	}
}

func (i InsertHandler) moveCursorDependingOnEntry(mode rune) {
	switch mode {
	case 'a':
		i.editor.Navigation.MoveRight(true)
	case 'A':
		i.editor.Navigation.JumpEndOfLine()
		i.editor.Navigation.MoveRight(true)
	case 'I':
		i.editor.Navigation.JumpBeginOfLine()
	}
}

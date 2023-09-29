package pkg

import (
	"github.com/eiannone/keyboard"
)

type InsertHandler struct {
	gim *Gim
}

func NewInsertHandler(gim *Gim) *InsertHandler {
	return &InsertHandler{
		gim: gim,
	}
}

func (i *InsertHandler) Activate(insertMode rune) {
	i.gim.Mode = INSERT
	// Depending on how insert mode was entered a different cursor position is required
	i.moveCursorDependingOnEntry(insertMode)
	i.gim.Term.ShowInsertModeInfo()
}

func (i *InsertHandler) Deactivate() {
	i.gim.Mode = NORMAL
	i.gim.Term.HideInsertModeInfo()
}

func (i *InsertHandler) Handle(keyEvent KeyEvent) {
	if i.gim.Mode != INSERT {
		return
	}

	if keyEvent.Key == keyboard.KeyEsc {
		i.Deactivate()
		i.gim.Navigation.MoveLeft()
	}

	i.InsertBeforeCursor(keyEvent)
}

func (i InsertHandler) InsertBeforeCursor(keyEvent KeyEvent) {
	// We have to handle different types of insertions
	// 1. Insertion at end of text
	// 2. Insertion in middle of text
	// 3. Insertion of a linebreak (at end of line and in middle of text)
	// 4. Deletion of a char
	x, y := i.gim.Navigation.Pos()

	switch keyEvent.Key {
	case 0, keyboard.KeySpace: // No special key is pressed
		if keyEvent.Key == keyboard.KeySpace {
			keyEvent.Char = ' '
		}
		// 1. Insertion at end of thext
		// Check cursor is at end of line
		if x == i.gim.Content.LineLength(y) {
			i.gim.Content.Buffer[y] = append(i.gim.Content.Buffer[y], keyEvent.Char)
		} else {
			// 2. Insertion in middle of text
			newLine := make([]rune, len(i.gim.Content.Buffer[y]))
			copy(newLine, i.gim.Content.Buffer[y][:x])
			newLine[x] = keyEvent.Char // Adds the new char
			copy(newLine[x+1:], i.gim.Content.Buffer[y][x:])
			i.gim.Content.Buffer[y] = newLine
		}

		TermUtils{}.InsertAt(x, y, i.gim.Content.Buffer[y])
		i.gim.Navigation.MoveRight(true)

	case keyboard.KeyEnter: // 3. Insertion of a linebreak
		// Check cursor is at end of line
		if x+1 >= i.gim.Content.LineLength(y) {
			if y < i.gim.Content.TotalLines()-1 {
				i.gim.Content.Buffer = append(i.gim.Content.Buffer[:y+1], append(make([][]rune, 0), i.gim.Content.Buffer[y+1:]...)...)
				TermUtils{}.InsertNewLine(y+1, i.gim.Content.Buffer[y+1:])
			} else {
				i.gim.Content.Buffer = append(i.gim.Content.Buffer, make([]rune, 1)) // Just add a new line buffer
			}
		} else {
			remaining := i.gim.Content.Buffer[y][:x]
			newLineContent := i.gim.Content.Buffer[y][x:]
			i.gim.Content.Buffer[y] = remaining
			i.gim.Content.Buffer = append(i.gim.Content.Buffer[:y+1], append([][]rune{newLineContent}, i.gim.Content.Buffer[y+1:]...)...)
			//content.Buffer = append(content.Buffer, newLineContent)
			TermUtils{}.RemoveFromTo(x, x+len(newLineContent), y)
			TermUtils{}.InsertNewLine(y, i.gim.Content.Buffer[y:])
		}
		i.gim.Navigation.MoveDownLineBegin()

	case keyboard.KeyBackspace: // 4. Deletion of a char
		i.gim.Content.Buffer[y] = append(i.gim.Content.Buffer[y][:x-1], i.gim.Content.Buffer[y][x:]...)
		TermUtils{}.RemoveCharAt(x, y, i.gim.Content.Buffer[y])
		i.gim.Navigation.MoveLeft()
	}
}

func (i InsertHandler) moveCursorDependingOnEntry(mode rune) {
	switch mode {
	case 'a':
		i.gim.Navigation.MoveRight(true)
	case 'A':
		i.gim.Navigation.JumpEndOfLine()
		i.gim.Navigation.MoveRight(true)
	case 'I':
		i.gim.Navigation.JumpBeginOfLine()
	}
}

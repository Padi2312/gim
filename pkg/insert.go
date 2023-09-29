package pkg

import (
	"github.com/eiannone/keyboard"
)

type InsertHandler struct {
	Active        bool
	cmd           rune
	navigation    *Navigation
	content       *Content
	normalHandler *NormalHandler
}

func NewInsertHandler(navigation *Navigation, content *Content) *InsertHandler {
	return &InsertHandler{
		Active:     false,
		navigation: navigation,
		cmd:        ' ',
		content:    content,
	}
}

func (i *InsertHandler) SetNormalHandler(normalModeHandler *NormalHandler) {
	i.normalHandler = normalModeHandler
}

func (i *InsertHandler) Activate(insertMode rune) {
	i.Active = true
	i.cmd = insertMode

	// Depending on how insert mode was entered a different cursor position is required
	switch i.cmd {
	case 'a':
		i.navigation.MoveRight(true)
	case 'A':
		i.navigation.JumpEndOfLine()
		i.navigation.MoveRight(true)
	case 'I':
		i.navigation.JumpBeginOfLine()
	}
}

func (i *InsertHandler) Deactivate() {
	i.Active = false
}

func (i *InsertHandler) Handle(keyEvent KeyEvent) {
	if !i.Active {
		return
	}

	if keyEvent.Key == keyboard.KeyEsc {
		i.Deactivate()
		i.navigation.MoveLeft()
		i.normalHandler.Activate()
	}

	i.InsertBeforeCursor(i.navigation, i.content, keyEvent)
}

func (i InsertHandler) InsertBeforeCursor(navigation *Navigation, content *Content, keyEvent KeyEvent) {
	// We have to handle different types of insertions
	// 1. Insertion at end of text
	// 2. Insertion in middle of text
	// 3. Insertion of a linebreak (at end of line and in middle of text)
	// 4. Deletion of a char
	x, y := navigation.Pos()

	switch keyEvent.Key {
	case 0, keyboard.KeySpace: // No special key is pressed
		if keyEvent.Key == keyboard.KeySpace {
			keyEvent.Char = ' '
		}
		// 1. Insertion at end of thext
		// Check cursor is at end of line
		if x == content.LineLength(y) {
			content.Buffer[y] = append(content.Buffer[y], keyEvent.Char)
		} else {
			// 2. Insertion in middle of text
			newLine := make([]rune, len(content.Buffer[y])+1)
			copy(newLine, content.Buffer[y][:x])
			newLine[x] = keyEvent.Char // Adds the new char
			copy(newLine[x+1:], content.Buffer[y][x:])
			content.Buffer[y] = newLine
		}

		TermUtils{}.InsertAt(x, y, content.Buffer[y])
		navigation.MoveRight(true)

	case keyboard.KeyEnter: // 3. Insertion of a linebreak
		// Check cursor is at end of line
		if x+1 >= content.LineLength(y) {
			if y < content.TotalLines()-1 {
				content.Buffer = append(content.Buffer[:y+1], append(make([][]rune, 1), content.Buffer[y+1:]...)...)
				TermUtils{}.InsertNewLine(y+1, content.Buffer[y+1:])
			} else {
				content.Buffer = append(content.Buffer, make([]rune, 1)) // Just add a new line buffer
			}
		} else {
			remaining := content.Buffer[y][:x]
			newLineContent := content.Buffer[y][x:]
			content.Buffer[y] = remaining
			content.Buffer = append(content.Buffer[:y+1], append([][]rune{newLineContent}, content.Buffer[y+1:]...)...)
			//content.Buffer = append(content.Buffer, newLineContent)
			TermUtils{}.RemoveFromTo(x, x+len(newLineContent), y)
			TermUtils{}.InsertNewLine(y, content.Buffer[y:])
		}
		navigation.MoveDownLineBegin()

	case keyboard.KeyBackspace: // 4. Deletion of a char
		content.Buffer[y] = append(content.Buffer[y][:x-1], content.Buffer[y][x:]...)
		TermUtils{}.RemoveCharAt(x, y, content.Buffer[y])
		navigation.MoveLeft()
	}
}

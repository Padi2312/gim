package main

import (
	"github.com/buger/goterm"
	"github.com/eiannone/keyboard"
	"github.com/padi2312/govim/pkg"
	"github.com/padi2312/govim/pkg/handlers"
)

type Gim struct {
	x, y             int
	mode             pkg.Event
	buffer           [][]rune
	commandHandler   *handlers.CommandHandler
	keyboardListener *pkg.KeyboardListener
}

func NewVimEditor() *Gim {
	buffer := make([][]rune, 1)
	buffer[0] = make([]rune, 0)
	return &Gim{
		x:                1,
		y:                1,
		mode:             pkg.Normal,
		buffer:           buffer,
		commandHandler:   handlers.NewCommandMode(),
		keyboardListener: pkg.NewKeyboardListener(),
	}
}

func (v *Gim) Run() {
	eventListener := pkg.NewEventListener()
	keyboardInputChannel := make(chan pkg.KeyEvent)

	go v.keyboardListener.Listen(keyboardInputChannel)

	for {
		goterm.MoveCursor(v.x, v.y)
		goterm.Flush()

		keyEvent := <-keyboardInputChannel
		switch v.mode {
		case pkg.Normal:
			v.handleNormalMode(keyEvent.Char, keyEvent.Key)
		case pkg.Insert:
			v.handleInsertMode(keyEvent.Char, keyEvent.Key)
		case pkg.Command:
			if !v.commandHandler.Handle(keyEvent.Char, keyEvent.Key) {
				v.mode = pkg.Normal
			}
		}

		if v.mode == pkg.Normal {
			v.mode = eventListener.CheckHandlerEvent(keyEvent)
		}
	}
}

func (v *Gim) ConvertCursorToBufferIndex() (int, int) {
	return v.x - 1, v.y - 1
}

func (v *Gim) handleNormalMode(char rune, key keyboard.Key) {
	bufX, bufY := v.ConvertCursorToBufferIndex()
	switch key {
	case keyboard.KeyEsc:
		return
	}

	switch char {
	case 'h':
		if bufX > 0 && v.buffer[bufY][bufX-1] != 0 {
			v.x--
		}
	case 'j':
		if bufY < len(v.buffer)-1 {
			v.y++
			if bufY+1 < len(v.buffer) {
				v.x = len(v.buffer[bufY+1]) + 1
			}
		}
	case 'k':
		if bufY > 0 {
			v.y--
			if bufY-1 >= 0 {
				v.x = len(v.buffer[bufY-1]) + 1
			}
		}
	case 'l':
		if bufX < len(v.buffer[bufY])-1 && v.buffer[bufY][bufX] != 0 {
			v.x++
		}
	case 'i':
		v.mode = pkg.Insert
	}
}

func (v *Gim) handleInsertMode(char rune, key keyboard.Key) {
	bufX, bufY := v.ConvertCursorToBufferIndex()
	if key == keyboard.KeyEsc {
		v.mode = pkg.Normal
	} else if key == keyboard.KeyEnter {
		v.y++
		v.x = 1
		if v.y > len(v.buffer) {
			newRow := make([]rune, 1)
			v.buffer = append(v.buffer, newRow)
		}
	} else if key == keyboard.KeyBackspace {
		if v.x > 1 {
			// Shift characters to the left starting from the cursor position
			for i := v.x - 1; i < len(v.buffer[v.y-1])-1; i++ {
				v.buffer[v.y-1][i-1] = v.buffer[v.y-1][i]
			}
			// Remove the last character in the line and minimize buffer
			v.buffer[v.y-1] = v.buffer[v.y-1][:len(v.buffer[v.y-1])-1]

			// Optionally minimize buffer length
			if len(v.buffer[v.y-1]) > 1 && len(v.buffer[v.y-1]) < cap(v.buffer[v.y-1])/2 {
				newBuffer := make([]rune, len(v.buffer[v.y-1]))
				copy(newBuffer, v.buffer[v.y-1])
				v.buffer[v.y-1] = newBuffer
			}

			// Update cursor and redraw line
			v.x--
			goterm.MoveCursor(1, v.y)
			goterm.Print(string(v.buffer[v.y-1]))
			goterm.Print(" ") // Clear residual characters
			goterm.MoveCursor(v.x, v.y)
		}
	} else {
		// Shift characters to the right starting from cursor position
		newRow := make([]rune, len(v.buffer[bufY])+1)
		copy(newRow[:bufX], v.buffer[bufY][:bufX])
		newRow[bufX] = char
		copy(newRow[bufX+1:], v.buffer[bufY][bufX:])
		v.buffer[bufY] = newRow

		// Redraw the line from cursor position to the end
		goterm.MoveCursor(v.x, v.y)
		goterm.Print(string(v.buffer[bufY][bufX:]))
		goterm.Flush()

		// Update the cursor
		v.x++
	}
}

func main() {
	goterm.Clear()
	vim := NewVimEditor()
	vim.Run()
}

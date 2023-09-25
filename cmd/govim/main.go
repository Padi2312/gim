package main

import (
	"os"

	"github.com/buger/goterm"
	"github.com/eiannone/keyboard"
)

const (
	Normal = iota
	Insert
	Command
)

type Gim struct {
	x, y   int
	mode   int
	buffer [][]rune
}

func NewVimEditor() *Gim {
	buffer := make([][]rune, 1)
	buffer[0] = make([]rune, 1)
	return &Gim{
		x:      1,
		y:      1,
		mode:   Normal,
		buffer: buffer,
	}
}

func (v *Gim) Run() {
	for {
		goterm.MoveCursor(v.x, v.y)
		goterm.Flush()

		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch v.mode {
		case Normal:
			v.handleNormalMode(char, key)
			if char == ':' {
				v.mode = Command
				v.buffer = append(v.buffer, []rune(":"))
			}
		case Insert:
			v.handleInsertMode(char, key)
		case Command:
			v.handleCommandMode(char, key)
		}
	}
}

func (v *Gim) ConvertCursorToBufferIndex() (int, int) {
	return v.x - 1, v.y - 1
}

func (v *Gim) handleCommandMode(char rune, key keyboard.Key) {
	staticPrompt := ":"
	if key == keyboard.KeyEnter {
		command := string(v.buffer[len(v.buffer)-1])[len(staticPrompt):]
		switch command {
		case "q":
			// Exit the editor
			os.Exit(1)
			return
			// Additional commands can be added here
		}
		// Clear command line and switch back to normal mode
		v.buffer[len(v.buffer)-1] = []rune(staticPrompt)
		v.mode = Normal
	} else if key == keyboard.KeyBackspace {
		if len(v.buffer[len(v.buffer)-1]) > len(staticPrompt) {
			// Remove the last character in the command
			v.buffer[len(v.buffer)-1] = v.buffer[len(v.buffer)-1][:len(v.buffer[len(v.buffer)-1])-1]
		}
	} else {
		// Append typed character to the command
		v.buffer[len(v.buffer)-1] = append(v.buffer[len(v.buffer)-1], char)
	}

	// Redraw the command line at the bottom
	goterm.MoveCursor(1, goterm.Height())
	goterm.Print(goterm.Background(goterm.Color(string(v.buffer[len(v.buffer)-1]), goterm.BLACK), goterm.WHITE))
	goterm.Flush()
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
	case 'q':
		return
	case 'i':
		v.mode = Insert
	}
}

func (v *Gim) handleInsertMode(char rune, key keyboard.Key) {
	bufX, bufY := v.ConvertCursorToBufferIndex()
	if key == keyboard.KeyEsc {
		v.mode = Normal
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
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	goterm.Clear()
	vim := NewVimEditor()
	vim.Run()
}

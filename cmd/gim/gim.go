package main

import (
	"github.com/buger/goterm"
	"github.com/eiannone/keyboard"
	"github.com/padi2312/govim/pkg"
)

type Mode = int

const (
	Normal Mode = iota
	Insert
	Command
	Visual
)

type GimV2 struct {
	mode    Mode
	content *pkg.Content
	cursor  *pkg.Cursor
}

func NewGimV2() *GimV2 {
	content := pkg.NewContent()
	return &GimV2{
		content: content,
		cursor:  pkg.NewCursor(content),
	}
}

func (g *GimV2) Run() {
	keyboadListener := *pkg.NewKeyboardListener()
	keyboardOutput := make(chan pkg.KeyEvent)
	go keyboadListener.Listen(keyboardOutput)

	for {
		keyEvent := <-keyboardOutput

		// Intercept ESC key to return back to Normal mode
		if keyEvent.Key > 0 && keyEvent.Key == keyboard.KeyEsc {
			g.switchMode(Normal)
			continue
		}

		// Depending on the curret mode the inputs are treated different
		switch g.mode {
		case Normal:
			g.HandleNormal(keyEvent)
		case Insert:
			g.HandleInsert(keyEvent)
		case Command:
			break
		case Visual:
			break
		}
	}
}

func (g *GimV2) switchMode(mode Mode) {
	g.mode = mode
}

func (g *GimV2) HandleNormal(keyEvent pkg.KeyEvent) {
	if keyEvent.Char == 'i' {
		g.switchMode(Insert)
	}

	switch keyEvent.Char {
	case 'i':
		g.switchMode(Insert)
	case 'h':
		g.cursor.Left()
	case 'j':
		g.cursor.Down()
	case 'k':
		g.cursor.Up()
	case 'l':
		g.cursor.Right()
	}
}

func (g *GimV2) HandleInsert(keyEvent pkg.KeyEvent) {
	g.content.InsertAtCursor(keyEvent, g.cursor)
}

func main() {
	goterm.Clear()
	vim := NewGimV2()
	vim.Run()
}

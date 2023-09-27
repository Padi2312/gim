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
	mode        Mode
	content     *pkg.Content
	cursor      *pkg.Cursor
	navigation  *pkg.Navigation
	display     *pkg.Display
	changeQueue *pkg.ChangeQueue
}

func NewGimV2() *GimV2 {
	changeQueue := pkg.NewChangeQueue()
	content := pkg.NewContent(changeQueue)
	cursor := pkg.NewCursor()
	display := pkg.NewDisplay(cursor, content)
	navigation := pkg.NewNavigation(content, cursor, changeQueue)
	return &GimV2{
		mode:        Normal,
		content:     content,
		cursor:      cursor,
		navigation:  navigation,
		display:     display,
		changeQueue: changeQueue,
	}
}

func (g *GimV2) Run() {
	// Setup keyboard listener
	keyboadListener := *pkg.NewKeyboardListener()
	keyboardOutput := make(chan pkg.KeyEvent)
	go keyboadListener.Listen(keyboardOutput)

	// Init display
	g.display.Init()

	for {
		keyEvent := <-keyboardOutput

		// Intercept ESC key to return back to Normal mode
		if keyEvent.Key > 0 && keyEvent.Key == keyboard.KeyEsc {
			if g.mode != Normal {
				g.navigation.MoveLeft()
			}
			g.switchMode(Normal)
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

		// Updates the terminal according to the changes
		g.display.Update(g.changeQueue)
	}
}

func (g *GimV2) switchMode(mode Mode) {
	g.mode = mode
}

func (g *GimV2) HandleNormal(keyEvent pkg.KeyEvent) {
	switch keyEvent.Char {
	case 'i':
		g.switchMode(Insert)
	case 'h':
		g.navigation.MoveLeft()
	case 'j':
		g.navigation.MoveDown()
	case 'k':
		g.navigation.MoveUp()
	case 'l':
		g.navigation.MoveRight(false)
	default:
		return
	}
}

func (g *GimV2) HandleInsert(keyEvent pkg.KeyEvent) {
	if keyEvent.Key == 0 {
		g.content.InsertBeforeCursor(keyEvent, g.cursor)
		g.navigation.MoveRight(true)
	} else {
		switch keyEvent.Key {
		case keyboard.KeyEnter:
			keyEvent.Char = '\n'
			g.content.InsertBeforeCursor(keyEvent, g.cursor)
			g.navigation.MoveDownLineBegin()
		case keyboard.KeyBackspace:
			g.content.DeleteBeforeCursor(g.cursor)
			g.navigation.MoveLeft()
			break
		case keyboard.KeySpace:
			keyEvent.Char = ' '
			g.content.InsertBeforeCursor(keyEvent, g.cursor)
			g.navigation.MoveRight(true)
		}
	}
}

func main() {
	goterm.Clear()
	vim := NewGimV2()
	vim.Run()
}

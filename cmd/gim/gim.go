package main

import (
	"github.com/eiannone/keyboard"
	"github.com/padi2312/govim/pkg"
)

type Mode = int

type Gim struct {
	content     *pkg.Content
	cursor      *pkg.Cursor
	navigation  *pkg.Navigation
	term        *pkg.Term
	normalMode  *pkg.NormalHandler
	insertMode  *pkg.InsertHandler
	commandMode *pkg.CommandHandler
}

func NewGim() *Gim {

	content := pkg.NewContent()
	cursor := pkg.NewCursor()
	navigation := pkg.NewNavigation(content, cursor)
	term := pkg.NewTerm(navigation)

	normalMode := pkg.NewNormalHandler(navigation)
	commandMode := pkg.NewCommandHandler(navigation, content)
	insertMode := pkg.NewInsertHandler(navigation, content)

	// !IMPORTANT: Connect the handlers to switch modes
	normalMode.SetInsertHandler(insertMode)
	normalMode.SetCommandHandler(commandMode)
	commandMode.SetNormalHandler(normalMode)
	insertMode.SetNormalHandler(normalMode)

	return &Gim{
		content:     content,
		cursor:      cursor,
		navigation:  navigation,
		term:        term,
		normalMode:  normalMode,
		insertMode:  insertMode,
		commandMode: commandMode,
	}
}

func (g *Gim) Run() {
	// Setup keyboard listener
	keyboadListener := *pkg.NewKeyboardListener()
	keyboardOutput := make(chan pkg.KeyEvent)
	go keyboadListener.Listen(keyboardOutput)

	// Activate normal mode by default
	g.normalMode.Activate()

	// Init display
	g.term.Clear()

	for {
		keyEvent := <-keyboardOutput

		if g.normalMode.Active {
			g.normalMode.Handle(keyEvent)
		} else if g.insertMode.Active {
			g.insertMode.Handle(keyEvent)
		} else {
			g.commandMode.Handle(keyEvent)
		}

		if g.insertMode.Active {
			g.term.ShowInsertModeInfo()
		} else {
			g.term.HideInsertModeInfo()
		}
		// Updates the terminal according to the changes
		g.term.Render()
	}
}

func main() {
	vim := NewGim()
	vim.Run()

	defer keyboard.Close()
}

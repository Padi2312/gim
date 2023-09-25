package handlers

import (
	"os"

	"github.com/buger/goterm"
	"github.com/eiannone/keyboard"
	"github.com/padi2312/govim/pkg"
)

type CommandHandler struct {
	Active bool
	buffer []rune
}

func NewCommandMode() *CommandHandler {
	return &CommandHandler{
		Active: false,
		buffer: make([]rune, 0),
	}
}

func (c *CommandHandler) Activate() {
	c.Active = true
	c.buffer = append(c.buffer, ':')
	c.drawCommandLine()
}

func (c *CommandHandler) Handle(char rune, key keyboard.Key) bool {
	// use key if char is \x00
	if char == '\x00' {
		switch key {
		case keyboard.KeyEnter:
			c.Execute()

		case keyboard.KeyBackspace:
			// First clean char on terminal and then remove
			// If you remove char first from buffer the length is not correct
			pkg.ClearCharAt(len(c.buffer), goterm.Height())
			c.RemoveChar()

		case keyboard.KeyEsc:
			c.buffer = []rune{}
			c.Active = false
			pkg.ClearLine(goterm.Height())
			return false
		}
	} else {
		c.AppendChar(char)
	}

	c.drawCommandLine()
	return true
}

func (c *CommandHandler) AppendChar(char rune) {
	c.buffer = append(c.buffer, char)
}

func (c *CommandHandler) RemoveChar() {
	if len(c.buffer) > 1 {
		c.buffer = c.buffer[:len(c.buffer)-1]
	}
}

func (c *CommandHandler) Execute() {
	cmd := string(c.buffer[1:len(c.buffer)])
	switch cmd {
	case "q":
		os.Exit(1)
	}
}

func (c *CommandHandler) drawCommandLine() {
	// Redraw the command line at the bottom
	goterm.MoveCursor(1, goterm.Height())
	goterm.Print(goterm.Background(goterm.Color(string(c.buffer), goterm.BLACK), goterm.WHITE))
	goterm.Flush()
}

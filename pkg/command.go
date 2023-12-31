package pkg

import (
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

type CommandHandler struct {
	editor      *Editor
	cmd      []rune
	termLine *Line
}

func NewCommandHandler(editor *Editor) *CommandHandler {
	return &CommandHandler{
		editor: editor,
		cmd: make([]rune, 0),
	}
}

func (c *CommandHandler) Activate() {
	c.editor.Mode = COMMAND
	c.termLine = NewLine(c.editor.Term.GetHeight() - 1)
	c.termLine.ClearFull()
	c.cmd = append(c.cmd, ':')
	c.termLine.AddChar(0, ':')
}

func (c *CommandHandler) Deactivate() {
	c.editor.Mode = NORMAL
	// Clear before removing cmd to get length of line for cleaning
	c.termLine.Clear()
	c.cmd = c.cmd[:0]
}

func (c *CommandHandler) Handle(keyEvent KeyEvent) {
	if c.editor.Mode != COMMAND {
		return
	}

	switch keyEvent.Key {
	case 0, keyboard.KeySpace:
		if keyEvent.Key == keyboard.KeySpace {
			keyEvent.Char = ' '
		}
		c.cmd = append(c.cmd, keyEvent.Char)
		c.termLine.AddChar(len(c.cmd)-1, keyEvent.Char)
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		if len(c.cmd)-1 >= 1 {
			c.cmd = c.cmd[:len(c.cmd)-1]
			c.termLine.RemoveChar(len(c.cmd))
		} else {
			c.Deactivate()
		}
	case keyboard.KeyEnter:
		c.execute()
	case keyboard.KeyEsc:
		c.Deactivate()
	}
}

// COMMAND: w
func (c *CommandHandler) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	content := c.editor.Content.ToString()
	_, err = file.WriteString(content)
	return err
}

func (c *CommandHandler) execute() {
	cmdString := string(c.cmd[1:])
	args := strings.Split(cmdString, " ")

	var additional string
	cmd := args[0]
	if len(args) > 1 {
		additional = args[1]
	}

	switch cmd {
	case "wq":
		if len(args) > 1 {
			c.editor.FileName = additional
		}

		if len(c.editor.FileName) == 0 {
			c.Deactivate()
			ErrorLine{}.PrintErrorLine("ERROR: no file name")
		} else {
			c.SaveToFile(c.editor.FileName)
			c.editor.Term.Clear()
			os.Exit(0)
		}
	case "q":
		os.Exit(0)
	default:
		c.Deactivate()
		ErrorLine{}.PrintErrorLine("ERROR: command not found")
	}

}

package pkg

import (
	"os"

	"github.com/eiannone/keyboard"
)

type CommandHandler struct {
	Active        bool
	cmd           []rune
	navigation    *Navigation
	content       *Content
	normalHandler *NormalHandler
}

func NewCommandHandler(navigation *Navigation, content *Content) *CommandHandler {
	return &CommandHandler{
		Active:     false,
		navigation: navigation,
		cmd:        make([]rune, 0),
		content:    content,
	}
}

func (c *CommandHandler) SetNormalHandler(normalModeHandler *NormalHandler) {
	c.normalHandler = normalModeHandler
}

func (c *CommandHandler) Activate() {
	c.Active = true
	c.cmd = append(c.cmd, ':')
}

func (c *CommandHandler) Deactivate() {
	c.Active = false
	c.cmd = c.cmd[:0]
}

func (c *CommandHandler) Handle(keyEvent KeyEvent) {
	if !c.Active {
		return
	}

	if keyEvent.Key == keyboard.KeyEsc {
		c.Deactivate()
		c.normalHandler.Activate()
	}

	if keyEvent.Key == keyboard.KeyEnter {
		c.execute()
	}

	if keyEvent.Key == 0 {
		c.cmd = append(c.cmd, keyEvent.Char)
	}
	// TODO: Display command line input
}

// COMMAND: w
func (c *CommandHandler) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	content := c.content.ToString()
	_, err = file.WriteString(content)
	return err
}

func (c *CommandHandler) execute() {
	for _, cmd := range c.cmd {

		switch cmd {
		case 'w':
			c.SaveToFile("test.txt")
		case 'q':
			os.Exit(0)
		}

	}
}

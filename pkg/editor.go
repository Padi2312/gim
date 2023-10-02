package pkg

import (
	"os"
	"strings"
)

type Mode = string

const (
	NORMAL  Mode = "normal"
	INSERT  Mode = "insert"
	COMMAND Mode = "command"
)

type Editor struct {
	Mode       Mode
	Navigation *Navigation
	Content    *Content
	Term       *Term

	FileName string

	// Handlers
	NormalHandler  *NormalHandler
	InsertHandler  *InsertHandler
	CommandHandler *CommandHandler
}

func NewEditor() *Editor {
	content := NewContent()
	navigation := NewNavigation(content)
	term := NewTerm(navigation)
	return &Editor{
		Mode:       NORMAL,
		Navigation: navigation,
		Content:    content,
		Term:       term,
		FileName:   "",
	}
}

func (e *Editor) ReadFile(filePath string) {
	if len(filePath) == 0 {
		return
	}

	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic("ERROR: Failed to read file" + err.Error())
		}
		lines := strings.Split(string(data), "\n")
		runeSlice := make([][]rune, len(lines))
		for i, line := range lines {
			if i != len(lines)-1 {
				// Remove linebreak character from content
				// because we dont want to print it to the terminal
				runeSlice[i] = []rune(line[:len(line)-1])
			} else {
				runeSlice[i] = []rune(line)
			}
		}
		e.Content.Buffer = runeSlice
		e.FileName = filePath
	} else if os.IsNotExist(err) {
		panic("ERROR: File not found.")
	} else {
		panic("ERROR: " + err.Error())
	}
}

func (e *Editor) Run() {
	// Setup all necessary handlers
	e.setupHandlers()

	// Init and clear terminal
	e.Term.Clear()

	// Setup keyboard listener
	keyboadListener := *NewKeyboardListener()
	keyboardOutput := make(chan KeyEvent)
	go keyboadListener.Listen(keyboardOutput)

	// Write content to term in case a file is  from ruopenend
	e.Term.WriteFullContent(e.Content.Buffer)

	for {
		// Wait for key input
		keyEvent := <-keyboardOutput

		// Depending on the current mode some keys
		// show a different behaviour
		switch e.Mode {
		case NORMAL:
			e.NormalHandler.Handle(keyEvent)
		case INSERT:
			e.InsertHandler.Handle(keyEvent)
		case COMMAND:
			e.CommandHandler.Handle(keyEvent)
		}

		if e.Mode != COMMAND {
			e.Term.Render()
		}
	}
}

func (e *Editor) setupHandlers() {
	e.NormalHandler = NewNormalHandler(e)
	e.InsertHandler = NewInsertHandler(e)
	e.CommandHandler = NewCommandHandler(e)
}

package pkg

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

func (e *Editor) Run() {
	// Setup all necessary handlers
	e.setupHandlers()

	// Init and clear terminal
	e.Term.Clear()

	// Setup keyboard listener
	keyboadListener := *NewKeyboardListener()
	keyboardOutput := make(chan KeyEvent)
	go keyboadListener.Listen(keyboardOutput)

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

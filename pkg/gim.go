package pkg

type Mode = string

const (
	NORMAL  Mode = "normal"
	INSERT  Mode = "insert"
	COMMAND Mode = "command"
)

type Gim struct {
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

func NewGim() *Gim {
	content := NewContent()
	navigation := NewNavigation(content)
	term := NewTerm(navigation)
	return &Gim{
		Mode:       NORMAL,
		Navigation: navigation,
		Content:    content,
		Term:       term,
		FileName:   "",
	}
}

func (g *Gim) Run() {
	// Setup all necessary handlers
	g.setupHandlers()

	// Init and clear terminal
	g.Term.Clear()

	// Setup keyboard listener
	keyboadListener := *NewKeyboardListener()
	keyboardOutput := make(chan KeyEvent)
	go keyboadListener.Listen(keyboardOutput)

	for {
		// Wait for key input
		keyEvent := <-keyboardOutput

		// Depending on the current mode some keys
		// show a different behaviour
		switch g.Mode {
		case NORMAL:
			g.NormalHandler.Handle(keyEvent)
		case INSERT:
			g.InsertHandler.Handle(keyEvent)
		case COMMAND:
			g.CommandHandler.Handle(keyEvent)
		}

		if g.Mode != COMMAND {

			g.Term.Render()
		}
	}
}

func (g *Gim) setupHandlers() {
	g.NormalHandler = NewNormalHandler(g)
	g.InsertHandler = NewInsertHandler(g)
	g.CommandHandler = NewCommandHandler(g)
}

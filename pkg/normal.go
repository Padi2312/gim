package pkg

type NormalHandler struct {
	Active            bool
	navigation        *Navigation
	insertModeHandler *InsertHandler
	commandHandler    *CommandHandler
}

func NewNormalHandler(navigation *Navigation) *NormalHandler {
	return &NormalHandler{
		Active:     false,
		navigation: navigation,
	}
}

func (n *NormalHandler) SetInsertHandler(insertModeHandler *InsertHandler) {
	n.insertModeHandler = insertModeHandler
}

func (n *NormalHandler) SetCommandHandler(commandHandler *CommandHandler) {
	n.commandHandler = commandHandler
}

func (n *NormalHandler) Activate() {
	n.Active = true
}

func (n *NormalHandler) Deactivate() {
	n.Active = false
}

func (n *NormalHandler) Handle(keyEvent KeyEvent) {
	if !n.Active {
		return
	}

	switch keyEvent.Char {
	case 'i', 'I', 'a', 'A':
		n.Deactivate() // Deactive self and activate InsertHandler
		n.insertModeHandler.Activate(keyEvent.Char)
	case ':':
		n.Deactivate()
		n.commandHandler.Activate()
	case 'h':
		n.navigation.MoveLeft()
	case 'j':
		n.navigation.MoveDown()
	case 'k':
		n.navigation.MoveUp()
	case 'l':
		n.navigation.MoveRight(false)
	}
}

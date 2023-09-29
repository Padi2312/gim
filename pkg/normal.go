package pkg

type NormalHandler struct {
	gim *Gim
}

func NewNormalHandler(gim *Gim) *NormalHandler {
	return &NormalHandler{
		gim: gim,
	}
}

func (n *NormalHandler) Use() {
	n.gim.Mode = NORMAL
}

func (n *NormalHandler) Handle(keyEvent KeyEvent) {
	if n.gim.Mode != NORMAL {
		return
	}

	switch keyEvent.Char {
	case 'i', 'I', 'a', 'A':
		n.gim.InsertHandler.Activate(keyEvent.Char)
	case ':':
		n.gim.CommandHandler.Activate()
	case 'h':
		n.gim.Navigation.MoveLeft()
	case 'j':
		n.gim.Navigation.MoveDown()
	case 'k':
		n.gim.Navigation.MoveUp()
	case 'l':
		n.gim.Navigation.MoveRight(false)
	}
}

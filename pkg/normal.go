package pkg

type NormalHandler struct {
	editor *Editor
}

func NewNormalHandler(editor *Editor) *NormalHandler {
	return &NormalHandler{
		editor: editor,
	}
}

func (n *NormalHandler) Use() {
	n.editor.Mode = NORMAL
}

func (n *NormalHandler) Handle(keyEvent KeyEvent) {
	if n.editor.Mode != NORMAL {
		return
	}

	switch keyEvent.Char {
	case 'i', 'I', 'a', 'A':
		n.editor.InsertHandler.Activate(keyEvent.Char)
	case ':':
		n.editor.CommandHandler.Activate()
	case 'h':
		n.editor.Navigation.MoveLeft()
	case 'j':
		n.editor.Navigation.MoveDown()
	case 'k':
		n.editor.Navigation.MoveUp()
	case 'l':
		n.editor.Navigation.MoveRight(false)
	case 'w':
		n.editor.Navigation.WordForward()
	case 'b':
		n.editor.Navigation.WordBackward()
	}
}

package pkg

import "github.com/eiannone/keyboard"

type Event int

const (
	Normal  Event = 0
	Insert  Event = 1
	Command Event = 2
)

type EventListener struct{}

func NewEventListener() *EventListener {
	return &EventListener{}
}

func (e *EventListener) CheckHandlerEvent(event KeyEvent) Event {
	if event.Char == '\x00' {
		if event.Key == keyboard.KeyEsc {
			return Normal
		}
	} else {
		switch event.Char {
		case ':':
			return Command
		case 'i':
			return Insert
		}
	}

	return Normal
}

package pkg

import "github.com/eiannone/keyboard"

type KeyEvent struct {
	Char rune
	Key  keyboard.Key
}

type KeyboardListener struct {
}

func NewKeyboardListener() *KeyboardListener {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	return &KeyboardListener{}
}

func (k *KeyboardListener) Listen(channel chan KeyEvent) {
	for {
		char, key := k.getCharAndKey()
		channel <- KeyEvent{Char: char, Key: key}
	}
}

func (k *KeyboardListener) getCharAndKey() (rune, keyboard.Key) {
	char, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}
	return char, key
}

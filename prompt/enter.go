package prompt

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

type Enter struct {
}

func (e *Enter) Run() bool {
	var b bool
	hiddenCursor()
	_ = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Enter {
			b = true
			return true, nil
		} else {
			return false, nil
		}
	})
	return b
}

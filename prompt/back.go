package prompt

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"os"
)

func Back() bool {
	var fg bool
	_ = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Backspace:
			fg = true
			return fg, nil
		case keys.CtrlC:
			Close()
			os.Exit(0)
		}
		return fg, nil
	})
	fmt.Printf(No)
	return fg
}

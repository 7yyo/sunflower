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
			fmt.Printf("%s\n", turnBack)
			return true, nil
		case keys.CtrlC:
			Close()
			os.Exit(0)
		}
		return
	})
	return fg
}

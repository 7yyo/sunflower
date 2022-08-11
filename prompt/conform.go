package prompt

import (
	"fmt"
	"github.com/chzyer/readline"
	"os"
)

func Conform() (bool, error) {
	cfg := &readline.Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	rl, err := readline.NewEx(cfg)
	if err != nil {
		return false, err
	}
	defer rl.Close()
	cfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		rl.SetPrompt(DarkGray("are you sure? [y/N] "))
		return line, 0, false
	})
	for {
		l, _ := rl.Readline()
		switch l {
		case "y":
			moveUp()
			cleanUpRow()
			fmt.Printf("%s\n", Yes)
			return true, nil
		case "N":
			moveUp()
			cleanUpRow()
			fmt.Printf("%s   %s", X, No)
			return false, nil
		default:
			continue
		}
	}
}

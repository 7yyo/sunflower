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
		moveUp()
		cleanUpRow()
		switch l {
		case "y":
			fmt.Printf("%s\n", Yes)
			return true, nil
		case "N":
			fmt.Printf("%s", No)
			return false, nil
		default:
			continue
		}
	}
}

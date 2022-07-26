package prompt

import (
	"fmt"
	"github.com/chzyer/readline"
	"os"
)

func Conform() (bool, error) {
	defer hiddenCursor()
	showCursor()
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
	pcfg := rl.GenPasswordConfig()
	pcfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		rl.SetPrompt(fmt.Sprintf("are you sure? [y/N] %v", string(line)))
		rl.Refresh()
		return line, 0, false
	})
	for {
		l, err := rl.ReadPasswordWithConfig(pcfg)
		if err != nil {
			return false, err
		}
		moveUp()
		cleanUpRow()
		switch string(l) {
		case "y":
			fmt.Printf("%s\n", yes)
			return true, nil
		case "N":
			fmt.Printf("%s\n", wrong)
			return false, nil
		}
	}
}

package prompt

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"os"
	"strings"
)

func Assign(input string) (string, error) {
	defer hiddenCursor()
	showCursor()
	cfg := &readline.Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	fmt.Printf(color.HiBlackString("type anything or `\\q` to return\n"))
	rl, err := readline.NewEx(cfg)
	if err != nil {
		return "", err
	}
	pcfg := rl.GenPasswordConfig()
	pcfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		rl.SetPrompt(fmt.Sprintf("%s %v", input, string(line)))
		rl.Refresh()
		return line, 0, false
	})
	for {
		l, err := rl.ReadPasswordWithConfig(pcfg)
		if err != nil {
			return "", err
		}
		switch strings.TrimSpace(string(l)) {
		case "\\q", "":
			clean()
		}
		return string(l), nil
	}
}

func clean() {
	for i := 0; i < 2; i++ {
		moveUp()
		cleanUpRow()
	}
	fmt.Printf("%s\n", turnBack)
}

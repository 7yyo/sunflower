package prompt

import (
	"fmt"
	"github.com/chzyer/readline"
	"reflect"
)

type Set struct {
	O     interface{}
	C     string
	D     []string
	Emoji string
}

func (se *Set) Run() (string, error) {
	v := reflect.ValueOf(se.O).Elem()
	c := v.FieldByName(se.C)
	showCursor()
	cfg := &readline.Config{}
	rl, _ := readline.NewEx(cfg)
	cfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		rl.SetPrompt(fmt.Sprintf("%s %s: ", se.Emoji, c))
		return line, pos, false
	})
	for _, dv := range se.D {
		d := v.FieldByName(dv)
		fmt.Printf(fmt.Sprintf(White("   %s: %s\n"), dv, d.String()))
	}
	for {
		return rl.Readline()
	}
}

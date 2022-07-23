package prompt

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"errors"
	"fmt"
	"github.com/chzyer/readline"
)

var up = "up"
var down = "down"
var left = "left"
var right = "right"

var Quit = "quit"

type Screen struct {
	Option  []string
	index   int
	Cfg     *Config
	r       *readline.Instance
	isFirst bool
	others  int
}

type Config struct {
	Emoji string
	Title string
}

func (s *Screen) Select() (int, string, error) {

	if err := s.init(); err != nil {
		return 0, "", err
	}
	defer func(r *readline.Instance) {
		_ = r.Close()
	}(s.r)

	s.writeString("sunflower")
	enter := false
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Up:
			s.writeString(up)
		case keys.Down:
			s.writeString(down)
		case keys.Left:
			s.writeString(left)
		case keys.Right:
			s.writeString(right)
		case keys.Enter:
			enter = true
			return true, nil
		case keys.Esc, keys.CtrlC:
			return true, errors.New(Quit)
		}
		return false, nil
	})
	if err != nil {
		return 0, "", err
	}
	if enter {
		return s.index, s.Option[s.index], nil
	}
	return 0, "", nil
}

func (s *Screen) init() error {
	s.defaultCfg()
	return s.newRl()
}

func (s *Screen) defaultCfg() {
	if s.Cfg == nil {
		s.Cfg = &Config{}
	}
	if s.Cfg.Emoji == "" {
		s.Cfg.Emoji = ">"
	}
	s.isFirst = true
	s.others = 1
	if s.Cfg.Title != "" {
		s.others++
	}
}

func (s *Screen) newRl() error {
	readLine, err := readline.NewEx(&readline.Config{})
	if err != nil {
		return err
	}
	s.r = readLine
	_, _ = s.r.Write([]byte("\033[?25l"))
	return nil
}

func (s *Screen) writeString(key string) {
	if !s.isFirst {
		s.cleanUp()
	}
	switch key {
	case up:
		if s.index > 0 {
			s.index--
		}
	case down:
		if s.index < len(s.Option)-1 {
			s.index++
		}
	case left:
		s.index = 0
	case right:
		s.index = len(s.Option) - 1
	}
	if s.Cfg.Title != "" {
		fmt.Println(White(s.Cfg.Title + ":"))
	}
	fmt.Println(White("← ↑ → ↓"))
	for i, o := range s.Option {
		if i == s.index {
			fmt.Printf("%s %s\n", s.Cfg.Emoji, o)
		} else {
			fmt.Printf(o + "\n")
		}
	}
	s.isFirst = false
}

func (s *Screen) cleanUp() {
	for i := 0; i < len(s.Option)+s.others; i++ {
		_, _ = s.r.Write([]byte("\033[1A"))
		_, _ = s.r.Write([]byte("\033[2K\r"))
	}
}

func Reset() {
	fmt.Printf("\033[?25h")
}

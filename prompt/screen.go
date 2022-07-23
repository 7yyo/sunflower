package prompt

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/chzyer/readline"
	"os"
)

type Screen struct {
	*Config
	*readline.Instance
	status
	Option []string
	index  int
	others int
}

type status struct {
	isFirst bool
	isEnter bool
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
	}(s.Instance)

	s.writeRow(keys.F7) // any key except up, down, left, right
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Enter {
			return s.enterEvent()
		} else {
			s.writeRow(key.Code)
		}
		return false, nil
	})
	if err != nil {
		return 0, "", err
	}
	if s.status.isEnter {
		return s.index, s.Option[s.index], nil
	}
	return 0, "", nil
}

func (s *Screen) init() error {
	s.defaultCfg()
	return s.newRl()
}

func (s *Screen) defaultCfg() {
	if s.Config == nil {
		s.Config = &Config{}
	}
	if s.Config.Emoji == "" {
		s.Config.Emoji = ">"
	}
	s.status.isFirst = true
	if s.Config.Title != "" {
		s.others++
	}
}

func (s *Screen) newRl() error {
	readLine, err := readline.NewEx(&readline.Config{})
	if err != nil {
		return err
	}
	s.Instance = readLine
	_, _ = s.Instance.Write([]byte(HiddenCursor))
	return nil
}

func (s *Screen) writeRow(k keys.KeyCode) {
	if !s.status.isFirst && k != keys.F7 {
		s.cleanUp()
	}
	switch k {
	case keys.Up:
		if s.index > 0 {
			s.index--
		}
	case keys.Down:
		if s.index < len(s.Option)-1 {
			s.index++
		}
	case keys.Left:
		s.index = 0
	case keys.Right:
		s.index = len(s.Option) - 1
	case keys.CtrlC:
		Close()
		os.Exit(0)
	}
	if s.Config.Title != "" {
		fmt.Println(White(s.Config.Title))
	}
	for i, o := range s.Option {
		if i == s.index {
			fmt.Printf("%s %s\n", s.Config.Emoji, o)
		} else {
			fmt.Printf(o + "\n")
		}
	}
	s.status.isFirst = false
}

func (s *Screen) enterEvent() (bool, error) {
	s.status.isEnter = true
	s.cleanUp()
	fmt.Printf("%s %s\n", CheckMark, s.Option[s.index])
	return true, nil
}

func (s *Screen) cleanUp() {
	for i := 0; i < len(s.Option)+s.others; i++ {
		_, _ = s.Instance.Write([]byte(moveUp))
		_, _ = s.Instance.Write([]byte(clean))
	}
}

func Close() {
	fmt.Printf(ShowCursor)
}

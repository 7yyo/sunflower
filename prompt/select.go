package prompt

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/chzyer/readline"
	"os"
	"reflect"
	"strings"
)

type Select struct {
	Emoji string
	Title string
	*Description
	*readline.Instance
	ch
	Option []interface{}
	index  int
	size   int
}

type ch struct {
	keyC      chan keys.Key
	keyEnterC chan bool
	bufferC   chan string
}

type Description struct {
	T string
	D []string
}

func (s *Select) Run() (int, interface{}, error) {
	if err := s.init(); err != nil {
		return 0, "", err
	}
	defer s.Instance.Close()
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			s.cleanUpScreen()
			if s.Description == nil {
				print(fmt.Sprintf("%s %s\n", CheckMark, s.Option[s.index]))
			} else {
				v := reflect.ValueOf(s.Option[s.index])
				t := v.FieldByName(s.Description.T)
				print(fmt.Sprintf("%s %s\n", CheckMark, t))
			}
			return true, nil
		default:
			s.keyC <- key
			return false, nil
		}
	})
	if err != nil {
		return 0, "", err
	} else {
		if s.Description == nil {
			return 0, s.Option[s.index].(string), nil
		} else {
			v := reflect.ValueOf(s.Option[s.index])
			t := v.FieldByName(s.Description.T)
			return 0, t.String(), nil
		}
	}
}

func (s *Select) init() error {
	if s.Emoji == "" {
		s.Emoji = ">"
	}
	s.keyC = make(chan keys.Key)
	s.keyEnterC = make(chan bool)
	s.bufferC = make(chan string)
	go s.printer()
	go s.keyEvent()
	hiddenCursor()
	return s.newReadline()
}

func (s *Select) newReadline() error {
	readLine, err := readline.NewEx(&readline.Config{})
	if err != nil {
		return err
	}
	s.Instance = readLine
	s.render()
	return nil
}

func (s *Select) keyEvent() {
	for {
		k := <-s.keyC
		switch k.Code {
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
		default:
			return
		}
		s.cleanUpScreen()
		s.render()
	}
}

func (s *Select) printer() {
	for {
		fmt.Printf(<-s.ch.bufferC)
	}
}

func (s *Select) render() {
	s.size = 0
	s.bufferC <- fmt.Sprintf("%s\n", fmt.Sprintf(White(s.Title+" "+ArrowKey)))
	s.size++
	for i, o := range s.Option {
		s.size++
		if s.Description == nil {
			if i == s.index {
				s.bufferC <- fmt.Sprintf("%s %s\n", s.Emoji, o)
			} else {
				s.bufferC <- fmt.Sprintf("%s\n", o)
			}
		} else {
			v := reflect.ValueOf(o)
			t := v.FieldByName(s.Description.T)
			if i == s.index {
				s.bufferC <- fmt.Sprintf("%s %s\n", s.Emoji, t)
				for _, dv := range s.D {
					d := v.FieldByName(dv)
					s.size += len(strings.Split(d.String(), "\n"))
					s.bufferC <- fmt.Sprintf(White("%s: %s\n"), dv, d.String())
				}
			} else {
				s.bufferC <- fmt.Sprintf("%s\n", t)
			}
		}
	}
}

func (s *Select) cleanUpScreen() {
	for i := 0; i < s.size; i++ {
		moveUp()
		clearRow()
	}
}

func Close() {
	showCursor()
}

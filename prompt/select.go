package prompt

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Select struct {
	Title  string
	Option []interface{}
	*Description
	ch
	buffer
}

type Description struct {
	T string
	D []string
}

type ch struct {
	keyC      chan keys.Key
	keyEnterC chan bool
	backC     chan bool
}

type buffer struct {
	index int
	size  int
	buf   strings.Builder
}

func (s *Select) Run() (int, interface{}, error) {
	s.init()
	defer close(s.ch.keyC)
	defer close(s.ch.keyEnterC)
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			s.cleanUpScreen()
			if s.Description == nil {
				s.printPoint(Yes)
			} else {
				v := reflect.ValueOf(s.Option[s.index])
				t := v.FieldByName(s.Description.T)
				fmt.Printf("%s %s\n", Yes, t)
			}
			return true, nil
		case keys.Backspace:
			s.cleanUpScreen()
			showCursor()
			fmt.Printf(No)
			return true, errors.New(key.Code.String())
		default:
			s.ch.keyC <- key
			return false, nil
		}
	})
	if err != nil {
		return 0, "", err
	} else {
		if s.Description == nil {
			return s.index, s.Option[s.index].(string), nil
		} else {
			v := reflect.ValueOf(s.Option[s.index])
			t := v.FieldByName(s.Description.T)
			return s.index, t.String(), nil
		}
	}
}

func (s *Select) init() {
	hiddenCursor()
	s.ch.keyC = make(chan keys.Key)
	s.ch.keyEnterC = make(chan bool)
	s.render()
	go s.keyEvent()
}

func (s *Select) keyEvent() {
	for {
		k := <-s.ch.keyC
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

func (s *Select) printPoint(emoji string) {
	fmt.Printf("%s %s\n", emoji, s.Option[s.index])
}

func (s *Select) render() {
	s.size = 0
	s.buffer.buf.Reset()
	if s.Title != "" {
		s.buffer.buf.WriteString(DarkGray(fmt.Sprintf("%s %s\n", ArrowKeys, s.Title)))
		s.size++
	}
	for i, o := range s.Option {
		if s.Description == nil {
			switch o.(type) {
			case string:
				if i == s.index {
					s.buffer.buf.WriteString(LightGreen(o.(string), 0, 1) + "\n")
				} else {
					s.buffer.buf.WriteString(fmt.Sprintf("%s\n", o))
				}
			default:
				panic("options is not a string type, which means your option is an object, please add a description attribute")
			}
		} else {
			v := reflect.ValueOf(o)
			t := v.FieldByName(s.Description.T)
			if i == s.index {
				s.buffer.buf.WriteString(Red(t.String(), 0, 1) + "\n")
				for _, dv := range s.D {
					d := v.FieldByName(dv)
					s.size += len(strings.Split(d.String(), "\n"))
					s.buffer.buf.WriteString(DarkGray(fmt.Sprintf("	%s: %s\n", dv, d.String())))
				}
			} else {
				s.buffer.buf.WriteString(fmt.Sprintf("%s\n", t.String()))
			}
		}
	}
	s.size += len(s.Option)
	fmt.Printf(s.buffer.buf.String())
}

func (s *Select) cleanUpScreen() {
	for i := 0; i < s.size; i++ {
		moveUp()
		cleanUpRow()
	}
}

func IsBackSpace(err error) bool {
	return err.Error() == keys.Backspace.String()
}

func Close() {
	showCursor()
}

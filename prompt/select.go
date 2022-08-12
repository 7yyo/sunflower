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
	Cap    int
	*Description
	gun []interface{}
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
	index         int
	cursor        int
	size          int
	upperBoundary bool
	lowerBoundary bool
	buf           strings.Builder
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
			if s.buffer.index > 0 {
				s.buffer.index--
				if s.Cap == 0 {
					s.buffer.cursor--
				}
			}
			if s.Cap != 0 {
				if s.index > 0 {
					if s.index < s.Cap-1 {
						s.buffer.cursor--
					} else {
						s.buffer.cursor = s.Cap - 1
					}
				} else if s.index == 0 {
					s.lowerBoundary = false
					s.buffer.cursor = 0
				}
			}
		case keys.Down:
			if s.index < len(s.Option)-1 {
				s.index++
				if s.Cap == 0 {
					s.buffer.cursor++
				}
			}
			if s.Cap != 0 {
				if s.buffer.lowerBoundary {
					s.buffer.cursor = s.Cap - 1
				} else {
					s.buffer.cursor++
				}
				if s.index%s.Cap == 0 {
					s.buffer.cursor = s.Cap - 1
					s.buffer.lowerBoundary = true
				}
			}
		case keys.Left:
			s.index = 0
			if s.Cap != 0 {
				s.cursor = 0
				s.upperBoundary = true
				s.lowerBoundary = false
			}
		case keys.Right:
			s.index = len(s.Option) - 1
			if s.Cap != 0 {
				s.cursor = s.Cap - 1
				s.upperBoundary = false
				s.lowerBoundary = true
			}
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
		s.buffer.buf.WriteString(LightGray(fmt.Sprintf("%s %s\n", ArrowKeys, s.Title)))
		s.size++
	}
	s.resetGun()
	for i, o := range s.gun {
		if s.Description == nil {
			switch o.(type) {
			case string:
				if i == s.cursor {
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
			if i == s.cursor {
				s.buffer.buf.WriteString(LightGreen(t.String(), 0, 1) + "\n")
				for _, dv := range s.D {
					d := v.FieldByName(dv)
					s.size += len(strings.Split(d.String(), "\n"))
					s.buffer.buf.WriteString(LightGray(fmt.Sprintf("	%s: %s\n", dv, d.String())))
				}
			} else {
				s.buffer.buf.WriteString(fmt.Sprintf("%s\n", t.String()))
			}
		}
	}
	fmt.Printf(s.buffer.buf.String())
	if s.Cap != 0 {
		print(LightGray(fmt.Sprintf("			%d/%d\n", s.index+1, len(s.Option))))
		s.size++
	}
	s.size += len(s.gun)
}

func (s *Select) resetGun() {
	s.gun = s.Option
	if len(s.Option) > s.Cap && s.Cap != 0 {
		if s.index <= s.Cap-1 {
			s.gun = s.gun[:s.Cap]
		} else if s.index > s.Cap-1 {
			s.gun = s.gun[s.index-s.Cap+1 : s.index+1]
		}
	}
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

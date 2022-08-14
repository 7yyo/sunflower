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
	*Desc

	ch
	buf
}

type Desc struct {
	T string
	D []string
}

type ch struct {
	keyC    chan keys.Key
	cursorC chan keys.Key
}

type buf struct {
	size int
	i    int             // index
	c    int             // cursor
	g    []interface{}   // clone option
	b    strings.Builder // buf
	lb   bool            // lower boundary
}

func (s *Select) Run() (int, interface{}, error) {
	s.first()
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			s.cleanUpScreen()
			if s.Desc == nil {
				s.printMark(Y)
			} else {
				v := reflect.ValueOf(s.Option[s.buf.i])
				t := v.FieldByName(s.Desc.T)
				fmt.Printf("%s %s\n", Y, t)
			}
			return true, nil
		case keys.Backspace:
			s.cleanUpScreen()
			showCursor()
			fmt.Printf(N)
			return true, errors.New(key.Code.String())
		default:
			s.ch.keyC <- key
			return false, nil
		}
	})
	if err != nil {
		return 0, "", err
	} else {
		if s.Desc == nil {
			return s.buf.i, s.Option[s.buf.i].(string), nil
		} else {
			v := reflect.ValueOf(s.Option[s.buf.i])
			t := v.FieldByName(s.Desc.T)
			return s.buf.i, t.String(), nil
		}
	}
}

func (s *Select) first() {
	hiddenCursor()
	s.ch.keyC = make(chan keys.Key)
	s.ch.cursorC = make(chan keys.Key)
	s.render()
	go s.keyEvent()
	go s.calCursor()
}

func (s *Select) keyEvent() {
	defer close(s.ch.keyC)
	for {
		k := <-s.ch.keyC
		switch k.Code {
		case keys.Up:
			if s.buf.i > 0 {
				s.buf.i--
			}
		case keys.Down:
			if s.buf.i < len(s.Option)-1 {
				s.buf.i++
			}
		case keys.Left:
			s.buf.i = 0
		case keys.Right:
			s.buf.i = len(s.Option) - 1
		case keys.CtrlC:
			Close()
			os.Exit(0)
		default:
			return
		}
		s.ch.cursorC <- k
	}
}

func (s *Select) calCursor() {
	defer close(s.ch.cursorC)
	for {
		k := <-s.ch.cursorC
		if s.Cap >= len(s.Option) || s.Cap == 0 {
			s.buf.c = s.buf.i
		} else {
			switch k.Code {
			case keys.Up:
				if s.buf.i >= 0 {
					if s.buf.i == 0 {
						s.buf.lb = false
						s.buf.c = s.buf.i
					} else {
						if s.buf.i < s.Cap-1 {
							s.buf.c--
						} else {
							s.buf.c = s.Cap - 1
						}
					}
				}
			case keys.Down:
				if s.buf.i%s.Cap == 0 || s.buf.lb {
					s.buf.c = s.Cap - 1
					s.buf.lb = true
				} else {
					s.buf.c++
				}
			case keys.Left:
				s.buf.c = s.buf.i
				s.buf.lb = false
			case keys.Right:
				s.buf.c = s.Cap - 1
				s.buf.lb = true
			}
		}
		s.cleanUpScreen()
		s.render()
	}
}

func (s *Select) resetBuf() {
	s.buf.g = s.Option
	if len(s.Option) > s.Cap && s.Cap != 0 {
		if s.buf.i <= s.Cap-1 {
			s.buf.g = s.buf.g[:s.Cap]
		} else if s.buf.i > s.Cap-1 {
			s.buf.g = s.buf.g[s.buf.i-s.Cap+1 : s.buf.i+1]
		}
	}
	s.buf.b.Reset()
}

func (s *Select) printMark(emoji string) {
	fmt.Printf("%s %s\n", emoji, s.Option[s.i])
}

func (s *Select) render() {
	s.buf.size = 0
	s.resetBuf()
	if s.Title != "" {
		s.buf.b.WriteString(LightGray(fmt.Sprintf("%s %s [%d/%d]\n", ArrowKeys, s.Title, s.buf.i+1, len(s.Option))))
		s.buf.size++
	}
	for i, o := range s.buf.g {
		if s.Desc == nil {
			switch o.(type) {
			case string:
				if i == s.buf.c {
					s.buf.b.WriteString(LightGreen(o.(string), 0, 1) + "\n")
				} else {
					s.buf.b.WriteString(fmt.Sprintf("%s\n", o))
				}
			default:
				panic("options is not a string type, which means your option is an object, please add a description attribute")
			}
		} else {
			v := reflect.ValueOf(o)
			t := v.FieldByName(s.Desc.T)
			if i == s.buf.c {
				s.buf.b.WriteString(LightGreen(t.String(), 0, 1) + "\n")
				for _, dv := range s.Desc.D {
					d := v.FieldByName(dv)
					s.buf.size += len(strings.Split(d.String(), "\n"))
					s.buf.b.WriteString(DarkGray(fmt.Sprintf("	%s: %s\n", dv, d.String())))
				}
			} else {
				s.buf.b.WriteString(fmt.Sprintf("%s\n", t.String()))
			}
		}
	}
	fmt.Printf(s.buf.b.String())
	s.buf.size += len(s.buf.g)
}

func (s *Select) cleanUpScreen() {
	for i := 0; i < s.buf.size; i++ {
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

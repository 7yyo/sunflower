package prompt

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"errors"
	"fmt"
	co "github.com/fatih/color"
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
	i    int
	c    int
	g    []interface{}
	b    strings.Builder
	lb   bool
}

func (s *Select) Run() (int, interface{}, error) {
	return s.exc()
}

func (s *Select) exc() (int, interface{}, error) {
	s.initialization()
	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			s.cleanUpScreen()
			if s.Desc == nil {
				s.printMark(yes)
			} else {
				v := reflect.ValueOf(s.Option[s.buf.i])
				t := v.FieldByName(s.Desc.T)
				fmt.Printf("%s %s\n", yes, t)
			}
			return true, nil
		case keys.Backspace:
			s.cleanUpScreen()
			fmt.Printf("%s\n", turnBack)
			return true, errors.New(key.Code.String())
		default:
			s.ch.keyC <- key
			s.ch.cursorC <- key
			return false, nil
		}
	})
	if err != nil {
		return 0, "", err
	}
	if s.Desc == nil {
		return s.buf.i, s.Option[s.buf.i].(string), nil
	} else {
		t := s.getT()
		return s.buf.i, t.String(), nil
	}
}

func (s *Select) initialization() {
	go s.open()
}

func (s *Select) open() {
	s.ch.keyC = make(chan keys.Key)
	s.ch.cursorC = make(chan keys.Key)
	hiddenCursor()
	s.render()
	for {
		select {
		case k := <-s.ch.keyC:
			s.keyEvent(k)
		case c := <-s.ch.cursorC:
			s.calCursor(c)
		}
	}
}

func (s *Select) keyEvent(k keys.Key) {
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
}

func (s *Select) calCursor(k keys.Key) {
	if s.Cap >= len(s.Option) || s.Cap == 0 {
		s.buf.c = s.buf.i
	} else {
		switch k.Code {
		case keys.Up:
			if s.buf.i == 0 {
				s.buf.lb = false
				s.buf.c = s.buf.i
			} else {
				if s.buf.i < s.Cap-1 {
					s.buf.lb = false
					s.buf.c--
				} else {
					s.buf.c = s.Cap - 1
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

func (s *Select) render() {
	s.resetBuf()
	s.setTitle()
	g := greenWithUl()
	b := hiBlackWithUl()
	for i, o := range s.buf.g {
		if s.Desc == nil {
			switch o.(type) {
			case string:
				if i == s.buf.c {
					s.bufAppend(g.Sprintf("%s\n", o.(string)))
				} else {
					s.bufAppend(fmt.Sprintf("%s\n", o))
				}
			default:
				panic("options is not a string type, which means your option is an object, please add a description attribute")
			}
		} else {
			v, t := s.reflect(o)
			if i == s.buf.c {
				s.bufAppend(g.Sprintf("    %s\n", t.String()))
				for _, dv := range s.Desc.D {
					d := v.FieldByName(dv)
					if strings.TrimSpace(d.String()) != "" {
						s.bufAppend(g.Sprintf("        [%s] %s\n", dv, d.String()))
					} else {
						s.bufAppend(b.Sprintf("        [%s] -\n", dv))
					}
				}
			} else {
				s.bufAppend(fmt.Sprintf("    %s\n", t.String()))
			}
		}
	}
	s.flush()
	s.setSize()
}

func (s *Select) resetBuf() {
	s.buf.size = 0
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

func (s *Select) bufAppend(b string) {
	s.buf.b.WriteString(b)
}

func (s *Select) setTitle() {
	if s.Title != "" {
		s.bufAppend(co.HiBlackString(fmt.Sprintf(
			"%s %s [%d/%d]\n",
			arrowKeys,
			s.Title,
			s.buf.i+1,
			len(s.Option))))
	}
}

func (s *Select) setSize() {
	s.buf.size = len(strings.Split(s.buf.b.String(), "\n")) - 1
}

func (s *Select) getT() reflect.Value {
	v := reflect.ValueOf(s.Option[s.buf.i])
	return v.FieldByName(s.Desc.T)
}

func (s *Select) reflect(o interface{}) (reflect.Value, reflect.Value) {
	v := reflect.ValueOf(o)
	t := v.FieldByName(s.Desc.T)
	return v, t
}

func (s *Select) flush() {
	fmt.Printf("%s", s.buf.b.String())
}

func (s *Select) cleanUpScreen() {
	for i := 0; i < s.buf.size; i++ {
		moveUp()
		cleanUpRow()
	}
}

func (s *Select) printMark(emoji string) {
	fmt.Printf("%s %s\n", emoji, s.Option[s.i])
}

func IsBackSpace(err error) bool {
	if err != nil {
		return err.Error() == keys.Backspace.String()
	}
	return false
}

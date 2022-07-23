package _example

import (
	"fmt"
	p "github.com/7yyo/sunflower/prompt"
)

func main() {

	defer p.Close()

	o := []string{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}

	s := p.Screen{
		Option: o,
		Cfg: &p.Config{
			Emoji: "\U0001F33B",
			Title: "Weekly Calendar",
		},
	}
	i, _, err := s.Select()
	if err != nil && err.Error() != p.Quit {
		panic(err)
	} else if err != nil && err.Error() == p.Quit {
		return
	}
	fmt.Printf("\nyou choose: %d, %s\n", i, o[i])
}

package main

import (
	"fmt"
	p "github.com/7yyo/sunflower/prompt"
)

type Demo struct {
	Name  string
	Price string
	Size  string
}

func main() {
	defer p.Close()
	xbox := Demo{Name: "Xbox", Price: "100", Size: "1TB"}
	ps5 := Demo{Name: "PS5", Price: "200", Size: "2TB"}
	nintendoSwitch := Demo{Name: "SWITCH", Price: "300", Size: "3TB"}
	os := []interface{}{xbox, ps5, nintendoSwitch}
	s := p.Screen{
		Option: os,
		Description: &p.Description{
			T: "Name",
			D: []string{"Price", "Size"},
		},
		Config: &p.Config{
			Emoji: p.VideoGame,
			Title: "Select your platform: ",
		},
	}
	i, r, _ := s.Select()

	fmt.Printf("index: %d, result: %s\n", i, r)

	n := []interface{}{"Jim", "Green", "Tom"}
	ss := p.Screen{
		Option: n,
	}
	i, r, _ = ss.Select()
	fmt.Printf("index: %d, result: %s\n", i, r)
}

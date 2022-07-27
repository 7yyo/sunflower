package main

import (
	"fmt"
	p "github.com/7yyo/sunflower/prompt"
)

type PlayForm struct {
	Name  string
	Price string
	Size  string
}

type config struct {
	Name        string
	Default     string
	Current     string
	Description string
}

func main() {

	defer p.Close()
	xbox := PlayForm{
		Name:  "Xbox",
		Price: "$299",
		Size:  "1TB",
	}
	ps5 := PlayForm{
		Name:  "PS5",
		Price: "$399",
		Size:  "2TB",
	}
	nintendoSwitch := PlayForm{
		Name:  "SWITCH",
		Price: "$299",
		Size:  "3TB",
	}
	os := []interface{}{xbox, ps5, nintendoSwitch}
	s := p.Select{
		Option: os,
		Description: &p.Description{
			T: "Name",
			D: []string{"Price", "Size"},
		},
		Emoji: p.VideoGame,
		Title: "Select your platform: ",
	}
	i, r, _ := s.Run()
	fmt.Printf("%d, %s\n", i, r)

	n := []interface{}{
		"Jim",
		"Green",
		"Tom",
	}
	ss := p.Select{
		Option: n,
	}
	i, r, _ = ss.Run()
	fmt.Printf("%d, %s\n", i, r)

	c := config{
		Name:        "max-replicas",
		Default:     "3",
		Current:     "10",
		Description: "Sets the maximum number of replicas.",
	}
	se := p.Set{
		O: &c,
		C: "Name",
		D: []string{
			"Default",
			"Current",
			"Description",
		},
		Emoji: p.VideoGame,
	}
	result := se.Run()
	fmt.Printf("%s\n", result)

	e := p.Enter{}
	fmt.Printf("return? %v", e.Run())
}

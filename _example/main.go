package main

import (
	"fmt"
	"github.com/7yyo/sunflower/prompt"
)

func main() {
	plan()
}

func plan() {
	defer prompt.Close()
	weeks := []interface{}{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}
	s := prompt.Select{
		Title:  "week plan",
		Option: weeks,
		Cap:    10,
	}
	_, _, err := s.Run()
	if err != nil {
		if prompt.IsBackSpace(err) {
			return
		} else {
			panic(err)
		}
	}
	err = showTodo()
	if err != nil {
		if prompt.IsBackSpace(err) {
			plan()
		} else {
			panic(err)
		}
	}
}

func showTodo() error {
	todo := []interface{}{
		"play",
		"read",
	}
	s := prompt.Select{
		Title:  "what do u want to do?",
		Option: todo,
	}
	_, r, err := s.Run()
	if err != nil {
		return err
	}
	err = todoInfo(r.(string))
	if err != nil {
		if prompt.IsBackSpace(err) {
			plan()
		} else {
			panic(err)
		}
	}
	return nil
}

type game struct {
	Name    string
	Details string
}

func todoInfo(todo string) error {
	switch todo {
	case "play":
		g1 := game{Name: "BloodBorne", Details: "ブラッドボーン BloodBorne https://www.playstation.com/ja-jp/games/bloodborne/"}
		g2 := game{Name: "OveredCooked!", Details: "Let"}
		games := []interface{}{
			g1,
			g2,
		}
		s := prompt.Select{
			Title:  "what game do u want to play?",
			Option: games,
			Desc: &prompt.Desc{
				T: "Name",
				D: []string{"Details"},
			},
		}
		_, r, err := s.Run()
		if err != nil {
			if prompt.IsBackSpace(err) {
				return showTodo()
			} else {
				panic(err)
			}
		}
		b, _ := prompt.Conform()
		if b {
			switch r {
			case "BloodBorne":
				print("欢迎来到亞楠...\nhttps://twitter.com/Bloodborne_PS4")
			case "OveredCooked!":
				print("Cooking.")
			}
			if prompt.Back() {
				return todoInfo("play")
			}
		} else {
			return todoInfo("play")
		}
	case "read":
		fmt.Printf("Reading")
	}
	if prompt.Back() {
		return showTodo()
	}
	return nil
}

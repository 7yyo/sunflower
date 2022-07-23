# 🌻 sunflower - prompt
<img src="https://github.com/7yyo/sunflower/blob/master/_example/example.gif" width="550px">

## example
```go
package main

import (
	p "github.com/7yyo/sunflower/prompt"
)

func main() {
	defer p.Close()
	s := p.Screen{
		Option: []string{
			"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
		},
		Config: &p.Config{
			Emoji: p.Sunflower,
			Title: "Weekly Calendar",
		},
	}
	_, _, _ = s.Select()
}
```

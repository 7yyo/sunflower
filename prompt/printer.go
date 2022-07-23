package prompt

import "fmt"

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func White(s string) string {
	return textColor(textWhite, s)
}

func textColor(c int, s string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", c, s)
}

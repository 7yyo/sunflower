package prompt

import "fmt"

var ArrowKey = "← ↑ → ↓"

var CheckMark = "\U00002714" // ✔
var Sunflower = "\U0001F33B" // 🌻
var VideoGame = "\U0001F3AE" // 🎮

var ShowCursor = "\033[?25h"
var HiddenCursor = "\u001B[?25l"

var HeadOfRow = "\u001B[1K\n"
var EndOfRow = "\u001B[0K\n"

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

func hiddenCursor() {
	print("\u001B[?25l")
}

func moveUp() {
	print("\033[1A")
}

func clearRow() {
	print("\033[2K\r")
}

func showCursor() {
	print(ShowCursor)
}

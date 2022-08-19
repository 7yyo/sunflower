package prompt

import (
	"fmt"
	"github.com/fatih/color"
)

var arrowKeys = "← ↑ → ↓"
var up = "↑"
var down = "↓"

var sunflower = "\U0001F33B"                       // 🌻
var yes = color.GreenString("\U00002714")          // ✔
var wrong = color.RedString("\U00002716")          // ✗
var turnBack = " " + color.RedString("\U000021A9") // ⏎

func showCursor() {
	print("\033[?25h")
}

func hiddenCursor() {
	print("\u001B[?25l")
}

func moveUp() {
	print("\033[1A")
}

func moveDown(size int) {
	fmt.Printf("\033[%dB", size)
}

func cleanScreenTail() {
	fmt.Printf("\033[K")
}

func cleanUpRow() {
	print("\033[2K")
}

func Close() {
	showCursor()
}

func greenWithUl() color.Color {
	return *color.New(color.FgGreen).Add(color.Bold).Add(color.Underline)
}

func hiBlackWithUl() color.Color {
	return *color.New(color.FgHiBlack).Add(color.Bold).Add(color.Underline)
}

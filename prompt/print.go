package prompt

import (
	"fmt"
	"github.com/fatih/color"
	"reflect"
	"strings"
)

var ArrowKeys = "← ↑ → ↓"
var Up = "↑"
var Down = "↓"
var Left = "←"
var Right = "→"

var Y = color.GreenString("\U00002714")            // ✔
var X = color.GreenString("\U00002716")            // ✗
var N = " " + color.RedString("\U000021A9") + "\n" // ⏎

var Sunflower = "\U0001F33B"            // 🌻
var VideoGame = "\U0001F3AE"            // 🎮
var Moai = "\U0001F5FF"                 // 🗿
var CardBox = "\U0001F5C3"              // 🗃
var VerticalTrafficLight = "\U0001F6A6" // 🚦

var ShowCursor = "\033[?25h"
var HiddenCursor = "\u001B[?25l"

var HeadOfRow = "\u001B[1K\n"
var EndOfRow = "\u001B[0K\n"

func hiddenCursor() {
	print("\u001B[?25l")
}

func moveUp() {
	print("\033[1A")
}

func cleanUpRow() {
	print("\033[2K")
}

func showCursor() {
	print(ShowCursor)
}

func Green(str string, modifier ...interface{}) string {
	return cliColorRender(str, 32, 0, modifier...)
}

func LightGreen(str string, modifier ...interface{}) string {
	return cliColorRender(str, 32, 1, modifier...)
}

func Cyan(str string, modifier ...interface{}) string {
	return cliColorRender(str, 36, 0, modifier...)
}

func LightCyan(str string, modifier ...interface{}) string {
	return cliColorRender(str, 36, 1, modifier...)
}

func Red(str string, modifier ...interface{}) string {
	return cliColorRender(str, 31, 0, modifier...)
}

func LightRed(str string, modifier ...interface{}) string {
	return cliColorRender(str, 31, 1, modifier...)
}

func Yellow(str string, modifier ...interface{}) string {
	return cliColorRender(str, 33, 0, modifier...)
}

func Black(str string, modifier ...interface{}) string {
	return cliColorRender(str, 30, 0, modifier...)
}

func DarkGray(str string, modifier ...interface{}) string {
	return cliColorRender(str, 30, 1, modifier...)
}

func LightGray(str string, modifier ...interface{}) string {
	return cliColorRender(str, 37, 0, modifier...)
}

func White(str string, modifier ...interface{}) string {
	return cliColorRender(str, 37, 1, modifier...)
}

func Blue(str string, modifier ...interface{}) string {
	return cliColorRender(str, 34, 0, modifier...)
}

func LightBlue(str string, modifier ...interface{}) string {
	return cliColorRender(str, 34, 1, modifier...)
}

func Purple(str string, modifier ...interface{}) string {
	return cliColorRender(str, 35, 0, modifier...)
}

func LightPurple(str string, modifier ...interface{}) string {
	return cliColorRender(str, 35, 1, modifier...)
}

func Brown(str string, modifier ...interface{}) string {
	return cliColorRender(str, 33, 0, modifier...)
}

func cliColorRender(str string, color int, weight int, extraArgs ...interface{}) string {
	var isBlink int64 = 0
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}
	var isUnderLine int64 = 0
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "05")
	}
	if isUnderLine > 0 {
		mo = append(mo, "04")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	return fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
}

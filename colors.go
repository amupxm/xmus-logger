package logger

import "fmt"

type Color int

const (
	Green Color = iota
	Red
	Yellow
	Blue
	Purple
	Cyan
	White
	Reset
)

var colorGreen = "\033[32m"
var colorRed = "\033[31m"
var colorYellow = "\033[33m"
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorCyan = "\033[36m"
var colorWhite = "\033[37m"
var colorReset = "\033[0m"

var colorMap = map[Color]string{
	Green:  colorGreen,
	Red:    colorRed,
	Yellow: colorYellow,
	Blue:   colorBlue,
	Purple: colorPurple,
	Cyan:   colorCyan,
	White:  colorWhite,
	Reset:  colorReset,
}

func SetColor(str string, c Color) string {
	return fmt.Sprintf("%s%s%s", colorMap[c], str, colorMap[Reset])
}

package loggerx

import "fmt"

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type Color uint8

func (c Color) Print(s string) string {
	return fmt.Sprintf("\x1b[1;%dm%s\x1b[0m", uint8(c), s)
}

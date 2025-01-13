package colorutil

import "github.com/lucasb-eyer/go-colorful"

func IsLight(c colorful.Color) bool {
	_, _, l := c.HSLuv()
	return l > 0.5
}

func IsDark(c colorful.Color) bool {
	return !IsLight(c)
}

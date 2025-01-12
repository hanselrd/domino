package colorutil

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

func IsLight(c color.Color) bool {
	cc, ok := colorful.MakeColor(c)
	if !ok {
		panic(ok)
	}
	_, _, l := cc.HSLuv()
	return l > 0.5
}

func IsDark(c color.Color) bool {
	return !IsLight(c)
}

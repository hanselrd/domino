package stringutil

import (
	"github.com/charmbracelet/x/ansi"
)

func AnsiSubstring[T ~string](s T, offset int, length uint) T {
	size := ansi.StringWidth(string(s))
	if offset < 0 {
		offset = size + offset
		if offset < 0 {
			offset = 0
		}
	}
	if offset >= size {
		return ""
	}
	if length > uint(size)-uint(offset) {
		length = uint(size - offset)
	}
	if offset == 0 {
		return T(ansi.Truncate(string(s), int(length), ""))
	}
	return T(ansi.Truncate(ansi.TruncateLeft(string(s), offset, ""), int(length), ""))
}

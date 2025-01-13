package sliceutil

import (
	"github.com/samber/lo"
)

func Rotate[T any](s []T, n int) []T {
	nn := lo.Ternary(n < 0, n*-1, n) % len(s)
	k := lo.Ternary(n <= 0, nn, len(s)-nn)
	return append(s[k:], s[:k]...)
}

func Convert[T, U any](s []T) ([]U, bool) {
	return lo.FromAnySlice[U](lo.ToAnySlice(s))
}

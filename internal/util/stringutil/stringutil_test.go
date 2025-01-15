package stringutil_test

import (
	"testing"

	"github.com/leaanthony/go-ansi-parser"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/util/stringutil"
)

func TestAnsiSubstring(t *testing.T) {
	s := "\u001b[1;31;40mHello\033[0m \u001b[0;30mWorld!\033[0m"
	if s := lo.Must(ansi.Cleanse(stringutil.AnsiSubstring(s, 2, 6))); s != "llo Wo" {
		t.Error(s)
	}
}

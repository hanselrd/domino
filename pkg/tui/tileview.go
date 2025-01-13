package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/util/sliceutil"
	"github.com/hanselrd/domino/pkg/face"
	"github.com/hanselrd/domino/pkg/tile"
)

type TileView struct {
	Hidden     bool
	Horizontal bool
	Rotate     int
	Data       *tile.Tile

	style lipgloss.Style
}

func NewTileView(data *tile.Tile, color colorful.Color) TileView {
	return TileView{
		Hidden:     false,
		Horizontal: false,
		Rotate:     0,
		Data:       data,
		style: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(color.Hex())).
			Foreground(lipgloss.Color(color.Hex())),
	}
}

func (tv TileView) View() string {
	ss := lo.Interleave(lo.Map(sliceutil.Rotate(tv.Data.Faces(), tv.Rotate), func(f face.Face, _ int) string {
		return fmt.Sprintf(" %s ", lo.Ternary(tv.Hidden, " ", f.String()))
	}), lo.Times(len(tv.Data.Faces())-1, func(_ int) string {
		d := lo.Ternary(tv.Horizontal, "|", "---")
		return lo.Ternary(tv.Hidden, strings.Repeat(" ", utf8.RuneCountInString(d)), d)
	}))
	return tv.style.Render(lo.Ternary(tv.Horizontal,
		lipgloss.JoinHorizontal(lipgloss.Center, ss...),
		lipgloss.JoinVertical(lipgloss.Center, ss...),
	))
}

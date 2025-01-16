package domino

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/util/sliceutil"
	. "github.com/hanselrd/domino/pkg/domino"
)

type TileView struct {
	Hidden     bool
	Horizontal bool
	Rotate     int
	Tile       *Tile

	style lipgloss.Style
}

func NewTileView(tile *Tile, color colorful.Color) TileView {
	return TileView{
		Hidden:     false,
		Horizontal: false,
		Rotate:     0,
		Tile:       tile,
		style: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(color.Hex())).
			Foreground(lipgloss.Color(color.Hex())),
	}
}

func (tv TileView) View() string {
	ss := lo.Interleave(
		lo.Map(
			sliceutil.Rotate(tv.Tile.Faces(), tv.Rotate),
			func(f Face, _ int) string {
				return fmt.Sprintf(
					" %s ",
					lo.Ternary(tv.Hidden, " ", f.String()),
				)
			},
		),
		lo.Times(len(tv.Tile.Faces())-1, func(_ int) string {
			d := lo.Ternary(tv.Horizontal, "|", "---")
			return lo.Ternary(
				tv.Hidden,
				strings.Repeat(" ", lo.RuneLength(d)),
				d,
			)
		}),
	)
	return tv.style.Render(lo.Ternary(tv.Horizontal,
		lipgloss.JoinHorizontal(lipgloss.Center, ss...),
		lipgloss.JoinVertical(lipgloss.Center, ss...),
	))
}

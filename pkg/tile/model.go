package tile

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"

	"github.com/hanselrd/domino/internal/util/sliceutil"
	"github.com/hanselrd/domino/pkg/face"
)

type Model[T constraints.Integer] struct {
	Hidden     bool
	Horizontal bool
	Rotate     int
	Data       *Tile[T]

	style lipgloss.Style
}

func NewModel[T constraints.Integer](data *Tile[T], color lipgloss.TerminalColor) Model[T] {
	return Model[T]{
		Hidden:     false,
		Horizontal: false,
		Rotate:     0,
		Data:       data,
		style: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(color).
			Foreground(color),
	}
}

func (m Model[T]) Init() tea.Cmd {
	return nil
}

func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
	return m, nil
}

func (m Model[T]) View() string {
	ss := lo.Interleave(lo.Map(sliceutil.Rotate(m.Data.Faces(), m.Rotate), func(f face.Face[T], _ int) string {
		return fmt.Sprintf(" %s ", lo.Ternary(m.Hidden, " ", f.String()))
	}), lo.Times(len(m.Data.faces)-1, func(_ int) string {
		d := lo.Ternary(m.Horizontal, "|", "―――")
		return lo.Ternary(m.Hidden, strings.Repeat(" ", utf8.RuneCountInString(d)), d)
	}))
	return m.style.Render(lo.Ternary(m.Horizontal,
		lipgloss.JoinHorizontal(lipgloss.Center, ss...),
		lipgloss.JoinVertical(lipgloss.Center, ss...),
	))
}

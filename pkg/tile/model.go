package tile

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/util/sliceutil"
	"github.com/hanselrd/domino/pkg/face"
)

type Model struct {
	Hidden     bool
	Horizontal bool
	Rotate     int
	Data       *Tile

	style lipgloss.Style
}

func NewModel(data *Tile, color lipgloss.TerminalColor) Model {
	return Model{
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	ss := lo.Interleave(lo.Map(sliceutil.Rotate(m.Data.Faces(), m.Rotate), func(f face.Face, _ int) string {
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

package domino

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/util/stringutil"
)

type ViewportModel struct {
	Base    viewport.Model
	XOffset int

	lines   []string
	columns int
}

func NewViewportModel(width, height int) ViewportModel {
	m := ViewportModel{
		Base: viewport.New(width, height),
	}
	m.Base.KeyMap.PageDown.Unbind()
	m.Base.KeyMap.PageDown.Unbind()
	m.Base.KeyMap.HalfPageDown.Unbind()
	m.Base.KeyMap.HalfPageUp.Unbind()
	m.Base.Style = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Margin(0, 1).
		Padding(1)
	return m
}

func (m ViewportModel) Init() tea.Cmd {
	return m.Base.Init()
}

func (m ViewportModel) Update(msg tea.Msg) (ViewportModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg, tea.MouseMsg:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, KeyMap.Left):
				m.SetXOffset(m.XOffset - 1)
			case key.Matches(msg, KeyMap.Right):
				m.SetXOffset(m.XOffset + 1)
			}
		case tea.MouseMsg:
			if !m.Base.MouseWheelEnabled || msg.Action != tea.MouseActionPress {
				break
			}
			switch msg.Button {
			case tea.MouseButtonWheelLeft:
				m.SetXOffset(m.XOffset - 1)
			case tea.MouseButtonWheelRight:
				m.SetXOffset(m.XOffset + 1)
			}
		}
	}
	m.Base.SetContent(strings.Join(lo.Map(m.lines, func(line string, _ int) string {
		return stringutil.AnsiSubstring(line, m.XOffset, uint(m.contentWidth()))
	}), "\n"))
	m.Base, cmd = m.Base.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ViewportModel) View() string {
	return m.Base.View()
}

func (m *ViewportModel) SetContent(s string) {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	m.lines = strings.Split(s, "\n")
	m.columns = lipgloss.Width(s)
}

func (m ViewportModel) maxXOffset() int {
	return max(0, m.columns-m.contentWidth())
}

func (m *ViewportModel) SetXOffset(n int) {
	m.XOffset = lo.Clamp(n, 0, m.maxXOffset())
}

func (m ViewportModel) contentWidth() int {
	w := m.Base.Width
	if sw := m.Base.Style.GetWidth(); sw != 0 {
		w = min(w, sw)
	}
	return w - m.Base.Style.GetHorizontalFrameSize()
}

func (m *ViewportModel) SetSize(width, height int) {
	m.Base.Width, m.Base.Height = width, height
}

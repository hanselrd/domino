package main

import (
	"fmt"
	"log"
	"runtime"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/build"
	"github.com/hanselrd/domino/pkg/face"
	"github.com/hanselrd/domino/pkg/tile"
	"github.com/hanselrd/domino/pkg/tileset"
)

type keyMap struct {
	left, down, up, right, help, quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.left, k.down, k.up, k.right},
		{k.help, k.quit},
	}
}

var keys = keyMap{
	left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	keys     keyMap
	help     help.Model
	ready    bool
	viewport viewport.Model
	ts       tileset.TileSet
	tvs      []tile.TileView
}

func initialModel() model {
	t2f6f := tile.NewTileFactory(face.NewUnsignedFaceFactory(6), 2)
	ts := lo.Must(tileset.NewTileSetShuffled(t2f6f))
	return model{
		keys: keys,
		help: help.New(),
		ts:   *ts,
		tvs: lo.Map(ts.Tiles(), func(t tile.Tile, _ int) tile.TileView {
			m := tile.NewTileView(&t, lipgloss.Color("#BB4712"))
			// m.Hidden = true
			m.Horizontal = t.IsMultiple()
			return m
		}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		helpHeight := lipgloss.Height(m.helpView())
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-helpHeight)
			m.viewport.Style = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#3C3C3C")).
				Margin(0, 1)
			ss := []string{}
			tss := lo.Chunk(m.tvs, len(m.tvs)/4)
			for _, ts := range tss {
				ss = append(ss, lipgloss.JoinHorizontal(lipgloss.Center, lo.Map(ts, func(t tile.TileView, _ int) string {
					return t.View()
				})...))
			}
			ss = slices.Concat(ss, []string{
				"-----------------------",
				"/// Build Metadata",
				"-----------------------",
				runtime.Version(),
				build.Version,
				build.Time,
				build.Hash,
				build.ShortHash,
				"-----------------------",
			})
			m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, ss...))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - helpHeight
		}
	case tea.KeyMsg:
		switch {
		// case key.Matches(msg, m.keys.left):
		// case key.Matches(msg, m.keys.down):
		// case key.Matches(msg, m.keys.up):
		// case key.Matches(msg, m.keys.right):
		case key.Matches(msg, m.keys.help):
			prevHelpHeight := lipgloss.Height(m.helpView())
			m.help.ShowAll = !m.help.ShowAll
			m.viewport.Height += prevHelpHeight
			m.viewport.Height -= lipgloss.Height(m.helpView())
			if m.viewport.PastBottom() {
				m.viewport.GotoBottom()
			}
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}
	case tea.MouseMsg:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return lipgloss.JoinVertical(lipgloss.Center, m.viewport.View(), m.helpView())
}

func (m model) helpView() string {
	return fmt.Sprintf("%s", m.help.View(m.keys))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

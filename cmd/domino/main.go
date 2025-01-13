package main

import (
	"log"
	"runtime"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/build"
	"github.com/hanselrd/domino/pkg/face"
	"github.com/hanselrd/domino/pkg/tile"
	"github.com/hanselrd/domino/pkg/tileset"
	"github.com/hanselrd/domino/pkg/tui"
)

type model struct {
	help     tui.HelpModel
	ready    bool
	viewport viewport.Model
	ts       tileset.TileSet
	tvs      []tui.TileView
}

func initialModel() model {
	t2f6f := tile.NewTileFactory(face.NewUnsignedFaceFactory(6), 2)
	ts := lo.Must(tileset.NewTileSet(t2f6f))
	return model{
		help: tui.NewHelpModel(),
		ts:   *ts,
		tvs: lo.Map(ts.Tiles(), func(t tile.Tile, i int) tui.TileView {
			m := tui.NewTileView(&t, colorful.HappyColor())
			m.Hidden = i%2 == 0
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
		helpHeight := lipgloss.Height(m.help.View())
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-helpHeight)
			m.viewport.Style = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#3C3C3C")).
				Margin(0, 1)
			ss := []string{}
			tvss := lo.Chunk(m.tvs, len(m.tvs)/4)
			for _, tvs := range tvss {
				ss = append(ss, lipgloss.JoinHorizontal(lipgloss.Center, lo.Map(tvs, func(tv tui.TileView, _ int) string {
					return tv.View()
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
				build.Dirty,
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
		// case key.Matches(msg, tui.Keys.Left):
		// case key.Matches(msg, tui.Keys.Down):
		// case key.Matches(msg, tui.Keys.Up):
		// case key.Matches(msg, tui.Keys.Right):
		case key.Matches(msg, tui.Keys.Help):
			prevHelpHeight := lipgloss.Height(m.help.View())
			m.help, cmd = m.help.Update(msg)
			cmds = append(cmds, cmd)
			m.viewport.Height += prevHelpHeight
			m.viewport.Height -= lipgloss.Height(m.help.View())
			if m.viewport.PastBottom() {
				m.viewport.GotoBottom()
			}
		case key.Matches(msg, tui.Keys.Quit):
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
	return lipgloss.JoinVertical(lipgloss.Center, m.viewport.View(), m.help.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

package domino

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/build"
	. "github.com/hanselrd/domino/pkg/domino"
)

type GameModel struct {
	Width, Height int
	Focus         int
	FocusStyle    lipgloss.Style
	Style         lipgloss.Style

	viewport ViewportModel
	help     HelpModel
	tileSet  TileSet
	tiles    []TileView
	ready    bool
}

func NewGameModel() GameModel {
	var m GameModel
	m.Style = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Padding(1)
	m.FocusStyle = m.Style.BorderForeground(lipgloss.Color("#FFDD99"))
	m.help = NewHelpModel()
	m.tileSet = *lo.Must(NewTileSet(NewTileFactory(NewUnsignedFaceFactory(15), 2), WithShuffle()))
	m.tiles = lo.Map(m.tileSet.Tiles(), func(t Tile, i int) TileView {
		tv := NewTileView(&t, colorful.HappyColor())
		// m.Hidden = i%2 == 0
		tv.Horizontal = !t.IsMultiple()
		return tv
	})
	return m
}

func (m GameModel) Init() tea.Cmd {
	return tea.Batch(m.viewport.Init(), m.help.Init())
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
		viewportWidth, viewportHeight := msg.Width*75/100, msg.Height*60/100
		if !m.ready {
			m.viewport = NewViewportModel(viewportWidth, viewportHeight)
			m.ready = true
		} else {
			m.viewport.SetSize(viewportWidth, viewportHeight)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyMap.Next):
			m.SetFocus(m.Focus + 1)
		case key.Matches(msg, KeyMap.Previous):
			m.SetFocus(m.Focus - 1)
		case key.Matches(msg, KeyMap.Quit):
			return m, tea.Quit
		}
	}
	content0 := lipgloss.JoinVertical(
		lipgloss.Left,
		lo.Map(
			lo.Chunk(m.tiles, len(m.tiles)/4),
			func(ts []TileView, _ int) string {
				return lipgloss.JoinHorizontal(lipgloss.Center,
					lo.Map(ts, func(t TileView, _ int) string {
						return t.View()
					})...)
			},
		)...)
	content1 := strings.Join(
		[]string{
			"-----------------------",
			"-----------------------",
			"/// DOMINO",
			"-----------------------",
			"-----------------------",
			"/// Game Information",
			"-----------------------",
			fmt.Sprintf("window.Size= %dx%d", m.Width, m.Height),
			fmt.Sprintf(
				"viewport.Size= %dx%d",
				m.viewport.Base.Width,
				m.viewport.Base.Height,
			),
			fmt.Sprintf(
				"viewport.Style.Size= %dx%d",
				m.viewport.Base.Style.GetWidth(),
				m.viewport.Base.Style.GetHeight(),
			),
			fmt.Sprintf(
				"viewport.Style.FrameSize= %dx%d",
				m.viewport.Base.Style.GetHorizontalFrameSize(),
				m.viewport.Base.Style.GetVerticalFrameSize(),
			),
			"-----------------------",
			"/// Build Metadata",
			"-----------------------",
			fmt.Sprintf("runtime.Version= %s", runtime.Version()),
			fmt.Sprintf("build.Version= %s", build.Version),
			fmt.Sprintf("build.Time= %s", build.Time),
			fmt.Sprintf("build.Hash= %s", build.Hash),
			fmt.Sprintf("build.ShortHash= %s", build.ShortHash),
			fmt.Sprintf("build.Dirty= %s", build.Dirty),
			"-----------------------",
		}, "\n")
	m.viewport.SetContent(
		lipgloss.JoinVertical(
			lipgloss.Left,
			// fmt.Sprintf(
			// 	"content0.Size= %dx%d",
			// 	lipgloss.Width(content0),
			// 	lipgloss.Height(content0),
			// ),
			// fmt.Sprintf(
			// 	"content1.Size= %dx%d",
			// 	lipgloss.Width(content1),
			// 	lipgloss.Height(content1),
			// ),
			fmt.Sprintf("domino.Tiles= %d", len(m.tileSet.Tiles())),
			content0,
			content1,
		),
	)
	// m.viewport.InnerWidth = 22
	// m.viewport.InnerHeight = 7
	model, cmd := m.viewport.Update(msg)
	if viewport, ok := model.(ViewportModel); ok {
		m.viewport = viewport
		cmds = append(cmds, cmd)
	}
	model, cmd = m.help.Update(msg)
	if help, ok := model.(HelpModel); ok {
		m.help = help
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m GameModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	// hPlayer := func(name string) string { return lipgloss.Place(35, 1, lipgloss.Center, lipgloss.Center, name) }
	// vPlayer := func(name string) string { return lipgloss.Place(4, 11, lipgloss.Center, lipgloss.Center, name) }
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			// m.Style.Render(hPlayer("Player 3")),
			lipgloss.JoinHorizontal(lipgloss.Center,
				// m.Style.Render(vPlayer("Player 4")),
				m.viewport.View(),
				// m.Style.Render(vPlayer("Player 2")),
			),
			// m.FocusStyle.Render(hPlayer("Player 1")),
			m.help.View(),
		),
	)
}

func (m *GameModel) SetFocus(focus int) {
	m.Focus = lo.Clamp(focus, 0, 1)
}

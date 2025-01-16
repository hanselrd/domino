package domino

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpModel struct {
	Base help.Model
}

func NewHelpModel() HelpModel {
	return HelpModel{
		Base: help.New(),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Base.Width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, KeyMap.Help) {
			m.Base.ShowAll = !m.Base.ShowAll
		}
	}
	return m, nil
}

func (m HelpModel) View() string {
	return m.Base.View(KeyMap)
}

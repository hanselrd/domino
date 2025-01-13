package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpModel struct {
	help help.Model
}

func NewHelpModel() HelpModel {
	return HelpModel{help: help.New()}
}

func (hm HelpModel) Init() tea.Cmd {
	return nil
}

func (hm HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		hm.help.Width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, Keys.Help) {
			hm.help.ShowAll = !hm.help.ShowAll
		}
	}
	return hm, nil
}

func (hm HelpModel) View() string {
	return hm.help.View(Keys)
}

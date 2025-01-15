package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	. "github.com/hanselrd/domino/internal/domino"
)

type model struct {
	game GameModel
}

func initialModel() model {
	return model{
		game: NewGameModel(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.game, cmd = m.game.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.game.View()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

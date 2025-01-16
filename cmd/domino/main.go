package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	. "github.com/hanselrd/domino/internal/domino"
)

func main() {
	m := NewGameModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

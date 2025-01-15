package domino

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Next     key.Binding
	Previous key.Binding
	Left     key.Binding
	Down     key.Binding
	Up       key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Previous},
		{k.Left, k.Down, k.Up, k.Right},
		{k.Help, k.Quit},
	}
}

var KeyMap = keyMap{
	Next: key.NewBinding(
		key.WithKeys("tab", "h"),
		key.WithHelp("tab", "focus next"),
	),
	Previous: key.NewBinding(
		key.WithKeys("shift+tab", "h"),
		key.WithHelp("S-tab", "focus previous"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc/C-c", "quit"),
	),
}

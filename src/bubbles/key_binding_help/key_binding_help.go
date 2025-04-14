package keybindinghelp

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Tab            key.Binding
	S              key.Binding
	P              key.Binding
	DownloadBlocks key.Binding
	Quit           key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.S, k.P, k.DownloadBlocks, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.S, k.P}, // first column
	}
}

var Keys = KeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch between views"),
	),
	DownloadBlocks: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "download blocks"),
	),
	S: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "select registry"),
	),
	P: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "change category paths"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

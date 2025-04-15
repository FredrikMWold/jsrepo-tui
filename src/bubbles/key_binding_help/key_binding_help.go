package keybindinghelp

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Tab            key.Binding
	S              key.Binding
	P              key.Binding
	DownloadBlocks key.Binding
	AddNewRegistry key.Binding
	Quit           key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.AddNewRegistry, k.DownloadBlocks, k.S, k.P}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.S, k.P}, // first column
	}
}

var Keys = KeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "focus next block list panel"),
	),
	DownloadBlocks: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "download blocks"),
	),
	AddNewRegistry: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "add registry"),
	),
	S: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "focus registries panel"),
	),
	P: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "focus categories panel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

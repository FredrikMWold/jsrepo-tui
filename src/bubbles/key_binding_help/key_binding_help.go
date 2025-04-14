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
	return []key.Binding{k.Tab, k.AddNewRegistry, k.DownloadBlocks, k.S, k.P, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.S, k.P}, // first column
	}
}

var Keys = KeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Switch between views"),
	),
	DownloadBlocks: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "Download blocks"),
	),
	AddNewRegistry: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "Add new registry"),
	),
	S: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Select registry"),
	),
	P: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Edit category paths"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}

package main

import (
	"fmt"
	config "jsrepo-tui/Config"
	registryselector "jsrepo-tui/RegistrySelector"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	registryselector registryselector.Model
}

func (m model) Init() tea.Cmd {
	return config.LoadConfig
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.registryselector, cmd = m.registryselector.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s string
	s += m.registryselector.View()
	return s
}

func main() {
	p := tea.NewProgram(model{
		registryselector: registryselector.New(),
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

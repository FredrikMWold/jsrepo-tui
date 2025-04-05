package main

import (
	"fmt"
	blocklist "jsrepo-tui/BlockList"
	config "jsrepo-tui/Config"
	dependencytable "jsrepo-tui/DependencyTable"
	manifestfetcher "jsrepo-tui/ManifestFetcher"
	registryselector "jsrepo-tui/RegistrySelector"
	selectedblocklist "jsrepo-tui/SelectedBlockList"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	selector = iota
	listView
	selectedBlocks
)

type model struct {
	registryselector registryselector.Model
	blocklist        blocklist.Model
	selectedBlocks   selectedblocklist.Model
	dependencytable  dependencytable.Model
	active           int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(config.LoadConfig, m.selectedBlocks.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			m.registryselector.Blur()
			if m.active == listView {
				m.active = selectedBlocks
				m.blocklist.Blur()
				m.selectedBlocks.Focus()
			} else {
				m.active = listView
				m.blocklist.Focus()
				m.selectedBlocks.Blur()
			}
		}
		switch msg.String() {
		case "s":
			m.active = selector
			m.registryselector.Focus()
			m.blocklist.Blur()
			m.selectedBlocks.Blur()
		}
	case blocklist.ListItem:
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.registryselector, _ = m.registryselector.Update(msg)
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		return m, nil
	case manifestfetcher.ManifestResponse:
		m.registryselector, _ = m.registryselector.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.active = listView
		return m, nil
	case selectedblocklist.RemoveBlock:
		m.dependencytable, _ = m.dependencytable.Update(msg)
		return m, nil
	}

	switch m.active {
	case selector:
		m.registryselector, cmd = m.registryselector.Update(msg)
		cmds = append(cmds, cmd)
	case listView:
		m.blocklist, cmd = m.blocklist.Update(msg)
		cmds = append(cmds, cmd)
	case selectedBlocks:
		m.selectedBlocks, cmd = m.selectedBlocks.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, lipgloss.JoinVertical(lipgloss.Left, m.registryselector.View(), m.dependencytable.View()), lipgloss.JoinHorizontal(lipgloss.Left, m.blocklist.View(), m.selectedBlocks.View()))
}

func main() {
	p := tea.NewProgram(model{
		registryselector: registryselector.New(),
		blocklist:        blocklist.New(),
		selectedBlocks:   selectedblocklist.New(),
		dependencytable:  dependencytable.New(),
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"jsrepo-tui/src/api/manifest"
	"jsrepo-tui/src/bubbles/block_list"
	"jsrepo-tui/src/bubbles/dependency_table"
	"jsrepo-tui/src/bubbles/registry_selector"
	"jsrepo-tui/src/bubbles/selected_block_list"
	"jsrepo-tui/src/config"
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
	registryselector registry_selector.Model
	blocklist        block_list.Model
	selectedBlocks   selected_block_list.Model
	dependencytable  dependency_table.Model
	width            int
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
	case block_list.ListItem:
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.registryselector, _ = m.registryselector.Update(msg)
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		return m, nil
	case manifest.ManifestResponse:
		m.registryselector, _ = m.registryselector.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.active = listView
		return m, nil
	case selected_block_list.RemoveBlock:
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
	var testHeader string
	testHeader = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(m.width - lipgloss.Width(m.dependencytable.View()) - 2).
		Render("test header \ntest\ntest\ntest\ntest\ntest\ntest\ntest")
	m.selectedBlocks.SetHeight(lipgloss.Height(m.selectedBlocks.View()) - lipgloss.Height(testHeader))
	m.blocklist.SetHeight(lipgloss.Height(m.blocklist.View()) - lipgloss.Height(testHeader))

	sidebar := lipgloss.JoinVertical(
		lipgloss.Left,
		m.registryselector.View(),
		m.dependencytable.View(),
	)

	dashboard := lipgloss.JoinVertical(lipgloss.Left,
		testHeader,
		lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.blocklist.View(),
			m.selectedBlocks.View(),
		))

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		sidebar,
		dashboard,
	)

}

func main() {
	p := tea.NewProgram(model{
		registryselector: registry_selector.New(),
		blocklist:        block_list.New(),
		selectedBlocks:   selected_block_list.New(),
		dependencytable:  dependency_table.New(),
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

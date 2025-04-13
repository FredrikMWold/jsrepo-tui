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

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

const (
	selector = iota
	listView
	selectedBlocks
	newRegistryInput
)

type model struct {
	registryselector registry_selector.Model
	blocklist        block_list.Model
	selectedBlocks   selected_block_list.Model
	dependencytable  dependency_table.Model
	newRegistryInput textinput.Model
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
		case tea.KeyEsc:
			if m.active == newRegistryInput {
				m.active = selector
				m.registryselector.Focus()
				m.blocklist.Blur()
				m.selectedBlocks.Blur()
				return m, nil
			}
		case tea.KeyEnter:
			if m.active == newRegistryInput {
				registries := viper.GetStringSlice("registries")
				registries = append(registries, m.newRegistryInput.Value())
				viper.Set("registries", registries)
				viper.WriteConfig()
				m.active = selector
				m.registryselector.Focus()
				return m, config.LoadConfig
			}
		}
		switch msg.String() {
		case "s":
			if m.active != newRegistryInput {
				m.active = selector
				m.registryselector.Focus()
				m.blocklist.Blur()
				m.selectedBlocks.Blur()
			}
		case "n":
			if m.active != newRegistryInput {
				m.active = newRegistryInput
				m.newRegistryInput.Focus()
				m.registryselector.Blur()
				m.blocklist.Blur()
				m.selectedBlocks.Blur()
				return m, nil
			}
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
	case manifest.ManifestNotFoundError:
		m.newRegistryInput.SetValue(string(msg))
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
	case newRegistryInput:
		m.newRegistryInput, cmd = m.newRegistryInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var newRegistryView string
	newRegistryHight := 0
	if m.active == newRegistryInput {
		newRegistryView = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("140")).
			Width(m.width - lipgloss.Width(m.dependencytable.View()) - 2).
			Render(m.newRegistryInput.View())
		newRegistryHight = lipgloss.Height(newRegistryView)
	}
	m.selectedBlocks.SetHeight(lipgloss.Height(m.selectedBlocks.View()) - newRegistryHight)
	m.blocklist.SetHeight(lipgloss.Height(m.blocklist.View()) - newRegistryHight)

	sidebar := lipgloss.JoinVertical(
		lipgloss.Left,
		m.registryselector.View(),
		m.dependencytable.View(),
	)

	dashboard := lipgloss.JoinVertical(lipgloss.Left,
		newRegistryView,
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
	input := textinput.New()
	input.Placeholder = "New registry"
	p := tea.NewProgram(model{
		registryselector: registry_selector.New(),
		blocklist:        block_list.New(),
		selectedBlocks:   selected_block_list.New(),
		dependencytable:  dependency_table.New(),
		newRegistryInput: input,
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

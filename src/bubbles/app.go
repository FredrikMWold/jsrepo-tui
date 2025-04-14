package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/block_list"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/categories_table"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/dependency_table"
	keybindinghelp "github.com/fredrikmwold/jsrepo-tui/src/bubbles/key_binding_help"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/registry_selector"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/selected_block_list"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
	"github.com/fredrikmwold/jsrepo-tui/src/config"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

const (
	selector = iota
	listView
	selectedBlocks
	categoriesTable
	newRegistryInput
)

type Model struct {
	registryselector registry_selector.Model
	blocklist        block_list.Model
	selectedBlocks   selected_block_list.Model
	dependencytable  dependency_table.Model
	newRegistryInput textinput.Model
	categoriestable  categories_table.Model
	keys             keybindinghelp.KeyMap
	help             help.Model
	width            int
	active           int
	error            manifest.BannerErrorMessage
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "github/<username>/<repo>"
	return Model{
		registryselector: registry_selector.New(),
		blocklist:        block_list.New(),
		selectedBlocks:   selected_block_list.New(),
		dependencytable:  dependency_table.New(),
		categoriestable:  categories_table.New(),
		newRegistryInput: input,
		keys:             keybindinghelp.Keys,
		help:             help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(config.LoadConfig, m.selectedBlocks.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.DownloadBlocks):
			m.categoriestable, cmd = m.categoriestable.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.Tab):
			m.registryselector.Blur()
			m.categoriestable.Blur()
			if m.active == listView {
				m.active = selectedBlocks
				m.blocklist.Blur()
				m.selectedBlocks.Focus()
			} else {
				m.active = listView
				m.blocklist.Focus()
				m.selectedBlocks.Blur()
			}
		case key.Matches(msg, m.keys.S):
			if m.active != newRegistryInput {
				m.active = selector
				m.registryselector.Focus()
				m.blocklist.Blur()
				m.categoriestable.Blur()
				m.selectedBlocks.Blur()
			}
		case key.Matches(msg, m.keys.P):
			if m.active != newRegistryInput {
				m.active = categoriesTable
				m.categoriestable.Focus()
				m.registryselector.Blur()
				m.blocklist.Blur()
				m.selectedBlocks.Blur()
				return m, nil
			}
		case key.Matches(msg, m.keys.AddNewRegistry):
			if m.active != newRegistryInput {
				m.active = newRegistryInput
				m.newRegistryInput.Focus()
				m.registryselector.Blur()
				m.blocklist.Blur()
				m.categoriestable.Blur()
				m.selectedBlocks.Blur()
				return m, nil
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEsc:
			if m.error != manifest.BannerErrorMessage("") {
				m.error = manifest.BannerErrorMessage("")
				return m, nil
			}
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
	case block_list.Blocks:
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.categoriestable, _ = m.categoriestable.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		msg.Height = msg.Height - 1
		m.registryselector, _ = m.registryselector.Update(msg)
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		m.categoriestable, _ = m.categoriestable.Update(msg)
		return m, nil
	case manifest.ManifestResponse:
		m.registryselector, _ = m.registryselector.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.categoriestable, _ = m.categoriestable.Update(msg)
		m.active = listView
		return m, nil
	case manifest.BannerErrorMessage:
		m.error = msg
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
	case categoriesTable:
		m.categoriestable, cmd = m.categoriestable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var newRegistryView string
	newRegistryHight := 0

	helpView := m.help.View(m.keys)

	if m.active == newRegistryInput {
		m.newRegistryInput.Width = m.width - lipgloss.Width(m.dependencytable.View()) - 2
		newRegistryView = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("140")).
			Width(m.width - lipgloss.Width(m.dependencytable.View()) - 2).
			Render(m.newRegistryInput.View())
		newRegistryHight = lipgloss.Height(newRegistryView)
	}
	if m.error != manifest.BannerErrorMessage("") {
		m.newRegistryInput.Width = m.width - lipgloss.Width(m.dependencytable.View()) - 2
		newRegistryView = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Width(m.width-lipgloss.Width(m.dependencytable.View())-2).
			Padding(0, 1).
			Render(string(m.error))
		newRegistryHight = lipgloss.Height(newRegistryView)
	}

	m.selectedBlocks.SetHeight(lipgloss.Height(m.selectedBlocks.View()) - newRegistryHight)
	m.blocklist.SetHeight(lipgloss.Height(m.blocklist.View()) - newRegistryHight)

	sidebar := lipgloss.JoinVertical(
		lipgloss.Left,
		m.registryselector.View(),
		m.categoriestable.View(),
		m.dependencytable.View(),
	)

	dashboard := lipgloss.JoinVertical(lipgloss.Left,
		newRegistryView,
		lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.blocklist.View(),
			m.selectedBlocks.View(),
		))
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			sidebar,
			dashboard,
		),
		helpView)

}

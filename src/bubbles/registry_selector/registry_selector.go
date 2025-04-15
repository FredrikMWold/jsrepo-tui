package registry_selector

import (
	"math"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/block_list"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
	"github.com/fredrikmwold/jsrepo-tui/src/config"
	"github.com/spf13/viper"
)

type Model struct {
	config config.Config
	table  table.Model
	focus  bool
}

func New() Model {
	columns := []table.Column{
		{Title: "Registries", Width: config.SidebarWidth},
	}
	t := table.New(table.WithColumns(columns), table.WithFocused(true))
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#11111b")).
		Background(lipgloss.Color("#cba6f7")).Bold(true)
	t.SetStyles(s)
	return Model{
		table: t,
		focus: true,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		margin := 5
		if msg.Height%3 != 0 {
			margin = 4
		}
		m.table.SetHeight(int(math.Ceil(float64(msg.Height-margin)/3)) - margin)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			registryName := m.table.SelectedRow()[0]
			cmds = append(cmds, manifest.GetManifest(registryName))
			cmds = append(cmds, func() tea.Msg {
				return block_list.Blocks{}
			})
		case tea.KeyDelete:
			registryName := m.table.SelectedRow()[0]
			var filteredRegistries []string
			for _, value := range m.config.Registries {
				if value != registryName {
					filteredRegistries = append(filteredRegistries, value)
				}
			}
			viper.Set("registries", filteredRegistries)
			viper.WriteConfig()
			return m, config.LoadConfig
		}

	case config.Config:
		m.config = msg
		rows := []table.Row{}
		for _, value := range m.config.Registries {
			rows = append(rows, table.Row{value})
		}
		m.table.SetRows(rows)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string
	if m.focus {
		s += lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("140")).
			Render(m.table.View())
	} else {
		s += lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(m.table.View())
	}
	return s
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

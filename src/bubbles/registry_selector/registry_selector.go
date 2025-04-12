package registry_selector

import (
	"jsrepo-tui/src/api/manifest"
	"jsrepo-tui/src/config"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const SidebarWidth = 36

type Model struct {
	config config.Config
	table  table.Model
	focus  bool
}

func New() Model {
	columns := []table.Column{
		{Title: "Registries", Width: SidebarWidth},
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
		margin := 4
		if msg.Height%2 != 0 {
			margin = 3
		}
		m.table.SetHeight((msg.Height - margin) / 2)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmds = append(cmds, manifest.GetManifest(m.table.SelectedRow()[0]))
			m.focus = false
		}
	case config.Config:
		m.config = msg
		rows := []table.Row{}
		for key := range m.config.Entries {
			rows = append(rows, table.Row{key})
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

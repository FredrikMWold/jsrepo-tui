package dependency_table

import (
	"jsrepo-tui/src/bubbles/block_list"
	"jsrepo-tui/src/bubbles/registry_selector"
	"jsrepo-tui/src/commands/manifest"
	"jsrepo-tui/src/helpers"
	"math"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	table  table.Model
	blocks []manifest.Block
}

func New() Model {
	columns := []table.Column{
		{Title: "Dependencies", Width: registry_selector.SidebarWidth},
	}
	t := table.New(table.WithColumns(columns), table.WithFocused(false))
	s := table.Styles{
		Selected: lipgloss.NewStyle().Bold(false),
		Header: lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(true).Padding(0, 1),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}

	t.SetStyles(s)
	return Model{
		table: t,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(int(math.Floor(float64(msg.Height) / 3)))
	case block_list.Blocks:
		m.blocks = msg
		var rows []table.Row
		var dependencies []string
		for _, block := range m.blocks {
			dependencies = append(dependencies, block.Dependencies...)
			dependencies = append(dependencies, block.DevDependencies...)
		}
		dependencies = helpers.UniqueStrings(dependencies)
		sort.Strings(dependencies)
		for _, dependency := range dependencies {
			rows = append(rows, table.Row{dependency})
		}

		m.table.SetRows(rows)
	}
	return m, cmd
}

func (m Model) View() string {
	var s string
	s += lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(m.table.View())
	return s
}

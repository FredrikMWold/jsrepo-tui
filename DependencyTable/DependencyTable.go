package dependencytable

import (
	blocklist "jsrepo-tui/BlockList"
	manifestfetcher "jsrepo-tui/ManifestFetcher"
	registryselector "jsrepo-tui/RegistrySelector"
	selectedblocklist "jsrepo-tui/SelectedBlockList"
	"jsrepo-tui/helpers"
	"slices"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func New() Model {
	columns := []table.Column{
		{Title: "Dependencies", Width: registryselector.SidebarWidth},
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

type Model struct {
	table table.Model
	data  []manifestfetcher.Block
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height/2 - 2)
	case blocklist.ListItem:
		isDuplicate := slices.ContainsFunc(m.data, func(item manifestfetcher.Block) bool {
			return item.Name == msg.Block.Name
		})
		if !isDuplicate {
			m.data = append(m.data, msg.Block)
		}
		var rows []table.Row
		var dependencies []string
		for _, block := range m.data {
			dependencies = append(dependencies, block.Dependencies...)
			dependencies = append(dependencies, block.DevDependencies...)
		}
		dependencies = helpers.UniqueStrings(dependencies)
		sort.Strings(dependencies)
		for _, dependency := range dependencies {
			rows = append(rows, table.Row{dependency})
		}

		m.table.SetRows(rows)
	case selectedblocklist.RemoveBlock:
		var blocks []manifestfetcher.Block
		for _, block := range m.data {
			if block.Name != msg.Name {
				blocks = append(blocks, block)
			}
		}
		m.data = blocks
		var rows []table.Row
		var dependencies []string
		for _, block := range m.data {
			dependencies = append(dependencies, block.Dependencies...)
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

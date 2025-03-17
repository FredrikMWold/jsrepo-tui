package registryselector

import (
	"fmt"
	config "jsrepo-tui/Config"
	manifestfetcher "jsrepo-tui/ManifestFetcher"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	config config.Config
	table  table.Model
	repo   manifestfetcher.Response
}

func New() Model {
	columns := []table.Column{
		{Title: "Registries", Width: 36},
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
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height/2 - 2)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmds = append(cmds, manifestfetcher.GetManifest(m.table.SelectedRow()[0]))
		}
	case config.Config:
		m.config = msg
		rows := []table.Row{}
		for key := range m.config.Entries {
			rows = append(rows, table.Row{key})
		}
		m.table.SetRows(rows)
		fmt.Printf("Rows: %v\n", rows)
	case manifestfetcher.Response:
		m.repo = msg
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns the view.
func (m Model) View() string {
	var s string
	s += lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(m.table.View())
	s += "\n"
	for _, category := range m.repo.Categories {
		for _, block := range category.Blocks {
			s += fmt.Sprintf("%s/%s\n", category.Name, block.Name)
		}
	}
	return s
}

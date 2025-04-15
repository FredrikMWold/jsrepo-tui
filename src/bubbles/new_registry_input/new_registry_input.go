package newregistryinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikmwold/jsrepo-tui/src/config"
	"github.com/spf13/viper"
)

const (
	newRegistryInput = iota
	errorMessage
	successMessage
)

type Model struct {
	newRegistryInput textinput.Model
	width            int
}

func New() Model {
	input := textinput.New()
	input.Focus()
	input.Placeholder = "github/<username>/<repo>@<branch>"
	return Model{
		newRegistryInput: input,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			registries := viper.GetStringSlice("registries")
			registries = append(registries, m.newRegistryInput.Value())
			viper.Set("registries", registries)
			viper.WriteConfig()
			return m, config.LoadConfig

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.newRegistryInput.Width = m.width - config.SidebarWidth - 11
	}
	m.newRegistryInput, cmd = m.newRegistryInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("140")).
		Padding(0, 1).
		Render(m.newRegistryInput.View())
}

package bannermessage

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	downloadblocks "github.com/fredrikmwold/jsrepo-tui/src/commands/download_blocks"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
	"github.com/fredrikmwold/jsrepo-tui/src/config"
)

const (
	errorMessage = iota
	successMessage
)

type Model struct {
	message     string
	messageType int
	width       int
}

func New() Model {
	return Model{
		message: "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case manifest.ManifestErrorMessage:
		m.message = string(msg)
		m.messageType = errorMessage
	case downloadblocks.DownloadBlocksErrorMessage:
		m.message = string(msg)
		m.messageType = errorMessage
	case downloadblocks.SuccessMessage:
		m.message = string(msg)
		m.messageType = successMessage
	}
	return m, cmd
}

func (m Model) View() string {
	color := lipgloss.Color("9")
	if m.messageType == successMessage {
		color = lipgloss.Color("6")
	}
	return lipgloss.NewStyle().
		Width(m.width-config.SidebarWidth-6).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(0, 1).
		Render(m.message)
}

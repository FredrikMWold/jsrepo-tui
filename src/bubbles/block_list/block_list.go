package block_list

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
	"github.com/fredrikmwold/jsrepo-tui/src/config"
)

type Model struct {
	repo           manifest.ManifestResponse
	focus          bool
	list           list.Model
	selectedBlocks []manifest.Block
}

type ListItem struct {
	Name     string
	Category string
	Block    manifest.Block
}

type Blocks []manifest.Block

func (i ListItem) Title() string       { return i.Name }
func (i ListItem) Description() string { return i.Category }
func (i ListItem) FilterValue() string { return i.Name }

func New() Model {
	list := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	list.DisableQuitKeybindings()
	list.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#11111b")).
		Background(lipgloss.Color("#cba6f7")).Padding(0, 1).Bold(true)
	list.Title = "Blocks"
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "add block")),
		}
	}

	return Model{
		list:  list,
		focus: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.list.SelectedItem() == nil {
				return m, nil
			}
			selectedItem := m.list.SelectedItem().(ListItem)
			isDuplicate := slices.ContainsFunc(m.selectedBlocks, func(block manifest.Block) bool {
				return block.Name == selectedItem.Title()
			})
			if !isDuplicate {
				m.selectedBlocks = append(m.selectedBlocks, selectedItem.Block)
			}
			return m, UpdateBlocks(m.selectedBlocks)
		}
	case manifest.ManifestResponse:
		m.repo = msg
		items := []list.Item{}
		for _, value := range m.repo.Categories {
			for _, block := range value.Blocks {
				items = append(items, ListItem{
					Name:     block.Name,
					Category: block.Category,
					Block:    block,
				})
			}
		}
		m.list.SetItems(items)
	case tea.WindowSizeMsg:
		m.list.SetWidth((msg.Width-config.SidebarWidth)/2 - 4)
		m.list.SetHeight(msg.Height - 2)
		return m, nil
	case Blocks:
		m.selectedBlocks = msg
	}
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string
	if m.focus {
		s += lipgloss.NewStyle().
			Width(m.list.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("140")).
			Render(m.list.View())
	} else {
		s += lipgloss.NewStyle().
			Width(m.list.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(m.list.View())
	}
	return s
}

func UpdateBlocks(blocks []manifest.Block) tea.Cmd {
	return func() tea.Msg {
		return Blocks(blocks)
	}
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m *Model) SetHeight(height int) {
	m.list.SetHeight(height - 2)
}

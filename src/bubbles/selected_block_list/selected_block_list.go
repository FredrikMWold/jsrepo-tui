package selected_block_list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/block_list"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/registry_selector"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
)

type Model struct {
	listView list.Model
	blocks   []manifest.Block
	focus    bool
}

func New() Model {
	list := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	list.DisableQuitKeybindings()
	list.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#11111b")).
		Background(lipgloss.Color("#cba6f7")).Padding(0, 1).Bold(true)
	list.Title = "Selected Blocks"
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("backspace"), key.WithHelp("backspace", "remove block")),
		}
	}

	return Model{
		listView: list,
		blocks:   []manifest.Block{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:

			listHasItems := m.listView.Items() != nil && len(m.listView.Items()) > 0
			if !listHasItems {
				return m, nil
			}
			selectedItem := m.listView.SelectedItem().(block_list.ListItem)
			var blocks []manifest.Block
			for _, block := range m.blocks {
				if block.Name != selectedItem.Name {
					blocks = append(blocks, block)
				}
			}
			return m, block_list.UpdateBlocks(blocks)
		case tea.KeyEnter:
			listHasItems := m.listView.Items() != nil && len(m.listView.Items()) > 0
			if !listHasItems {
				return m, nil
			}
		}

	case block_list.Blocks:
		m.blocks = msg
		items := []list.Item{}
		for _, block := range m.blocks {
			items = append(items, block_list.ListItem{
				Name:     block.Name,
				Category: block.Category,
				Block:    block,
			})

		}
		m.listView.SetItems(items)
		return m, cmd

	case tea.WindowSizeMsg:
		margin := 4
		if msg.Width%2 != 0 {
			margin = 3
		}
		m.listView.SetWidth((msg.Width-registry_selector.SidebarWidth)/2 - margin)
		m.listView.SetHeight(msg.Height - 2)
		return m, nil

	}

	m.listView, cmd = m.listView.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.focus {
		return lipgloss.NewStyle().
			Width(m.listView.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("140")).
			Render(m.listView.View())
	} else {
		return lipgloss.NewStyle().
			Width(m.listView.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(m.listView.View())
	}
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m *Model) SetHeight(height int) {
	m.listView.SetHeight(height - 2)
}

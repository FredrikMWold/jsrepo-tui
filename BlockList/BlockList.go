package blocklist

import (
	manifestfetcher "jsrepo-tui/ManifestFetcher"
	registryselector "jsrepo-tui/RegistrySelector"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	repo  manifestfetcher.ManifestResponse
	focus bool
	list  list.Model
}

type ListItem struct {
	Name  string
	Path  string
	Block manifestfetcher.Block
}

func (i ListItem) Title() string       { return i.Name }
func (i ListItem) Description() string { return i.Path }
func (i ListItem) FilterValue() string { return i.Name }

func New() Model {
	list := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	list.DisableQuitKeybindings()
	list.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#11111b")).
		Background(lipgloss.Color("#cba6f7")).Padding(0, 1).Bold(true)
	list.Title = "Blocks"

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
			for _, item := range m.getLocalDependencies() {
				cmds = append(cmds, addBlock(item))
			}
			cmds = append(cmds, addBlock(m.list.SelectedItem().(ListItem)))
			return m, tea.Batch(cmds...)
		}
	case manifestfetcher.ManifestResponse:
		m.repo = msg
		items := []list.Item{}
		for _, value := range m.repo.Categories {
			for _, block := range value.Blocks {
				items = append(items, ListItem{
					Name:  block.Name,
					Path:  block.Directory,
					Block: block,
				})
			}
		}
		m.list.SetItems(items)
		m.focus = true
	case tea.WindowSizeMsg:
		m.list.SetWidth((msg.Width - registryselector.SidebarWidth - 8) / 2)
		m.list.SetHeight(msg.Height - 2)
		return m, nil

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

func addBlock(block ListItem) tea.Cmd {
	return func() tea.Msg {
		return ListItem{
			Name:  block.Name,
			Path:  "." + string(os.PathSeparator) + block.Path,
			Block: block.Block,
		}
	}
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m Model) getLocalDependencies() []ListItem {
	var localDependencies []ListItem
	for _, localDependency := range m.list.SelectedItem().(ListItem).Block.LocalDependencies {
		blockName := strings.Split(localDependency, "/")[len(strings.Split(localDependency, "/"))-1]
		for _, category := range m.repo.Categories {
			for _, block := range category.Blocks {
				if block.Name == blockName {
					localDependencies = append(localDependencies, ListItem{
						Name:  block.Name,
						Path:  block.Directory,
						Block: block,
					})
				}
			}
		}
	}
	return localDependencies
}

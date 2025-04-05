package selectedblocklist

import (
	blocklist "jsrepo-tui/BlockList"
	manifestfetcher "jsrepo-tui/ManifestFetcher"
	registryselector "jsrepo-tui/RegistrySelector"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RemoveBlock manifestfetcher.Block

const (
	listView = iota
	filePickerView
)

type Model struct {
	listView       list.Model
	filePickerView filepicker.Model
	active         int
	focus          bool
}

func New() Model {
	list := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	list.DisableQuitKeybindings()
	list.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#11111b")).
		Background(lipgloss.Color("#cba6f7")).Padding(0, 1).Bold(true)
	list.Title = "Selected Blocks"

	filePicker := filepicker.New()
	filePicker.CurrentDirectory, _ = os.Getwd()
	filePicker.DirAllowed = true
	filePicker.FileAllowed = false

	return Model{
		listView:       list,
		filePickerView: filePicker,
		active:         listView,
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
			if m.active == filePickerView {
				m.active = listView
				m.filePickerView.CurrentDirectory, _ = os.Getwd()
				return m, nil
			}
			selectedItem := m.listView.SelectedItem().(blocklist.ListItem)
			m.listView.RemoveItem(m.listView.Index())
			return m, RemoveItem(selectedItem.Block)
		case tea.KeyEnter:
			if m.active == listView {
				m.active = filePickerView
				return m, m.filePickerView.Init()
			}
		case tea.KeyEsc:
			if m.active == filePickerView {
				m.active = listView
				m.filePickerView.CurrentDirectory, _ = os.Getwd()
				return m, nil
			}
		}

	case blocklist.ListItem:
		isDuplicate := slices.ContainsFunc(m.listView.Items(), func(item list.Item) bool {
			return item.(blocklist.ListItem).Title() == msg.Title()
		})
		if !isDuplicate {
			cmd = m.listView.InsertItem(-1, msg)
		}
		return m, cmd

	case tea.WindowSizeMsg:
		margin := 4
		if msg.Width%2 != 0 {
			margin = 3
		}
		m.listView.SetWidth((msg.Width-registryselector.SidebarWidth)/2 - margin)
		m.listView.SetHeight(msg.Height - 2)
		m.filePickerView.Height = msg.Height - 3
		return m, nil

	}

	switch m.active {
	case filePickerView:
		cwd, err := os.Getwd()
		if err != nil {
			return m, nil
		}
		m.filePickerView, cmd = m.filePickerView.Update(msg)
		if didSelect, path := m.filePickerView.DidSelectFile(msg); didSelect {
			relativePath, err := filepath.Rel(cwd, path)
			if err != nil {
				return m, nil
			}
			currentIndex := m.listView.Index()
			currentItem := m.listView.SelectedItem().(blocklist.ListItem)
			m.listView.RemoveItem(m.listView.Index())
			m.listView.InsertItem(currentIndex, blocklist.ListItem{
				Block: currentItem.Block,
				Name:  currentItem.Name,
				Path:  "." + string(os.PathSeparator) + relativePath,
			})
			m.filePickerView.CurrentDirectory, _ = os.Getwd()
			m.active = listView
		}
		cmds = append(cmds, cmd)
	case listView:
		m.listView, cmd = m.listView.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var activeView string

	switch m.active {
	case listView:
		activeView = m.listView.View()
	case filePickerView:
		activeView = m.filePickerView.View()
	}

	if m.focus {
		return lipgloss.NewStyle().
			Width(m.listView.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("140")).
			Render(activeView)
	} else {
		return lipgloss.NewStyle().
			Width(m.listView.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(activeView)
	}
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func RemoveItem(item manifestfetcher.Block) tea.Cmd {
	return func() tea.Msg {
		return RemoveBlock(item)
	}
}

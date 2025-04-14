package categories_table

import (
	"fmt"
	"jsrepo-tui/src/api/manifest"
	"jsrepo-tui/src/bubbles/block_list"
	"jsrepo-tui/src/bubbles/registry_selector"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	tableView = iota
	filePickerView
)

type Model struct {
	filePicker    filepicker.Model
	table         table.Model
	active        int
	focus         bool
	categoryPaths map[string]string
	blocks        []manifest.Block
	repo          manifest.ManifestResponse
}

func New() Model {
	columns := []table.Column{
		{Title: "Category", Width: registry_selector.SidebarWidth/6*2 - 4},
		{Title: "Path", Width: registry_selector.SidebarWidth/6*4 + 4},
	}
	t := table.New(table.WithColumns(columns), table.WithFocused(false))
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

	f := filepicker.New()
	f.CurrentDirectory, _ = os.Getwd()
	f.DirAllowed = true
	f.FileAllowed = false
	f.ShowPermissions = false
	f.ShowSize = false

	return Model{
		table:         t,
		filePicker:    f,
		focus:         false,
		categoryPaths: map[string]string{},
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
		case tea.KeyEnter:
			if m.active == tableView {
				m.active = filePickerView
				return m, m.filePicker.Init()
			}

		}
		switch msg.String() {
		case "d":
			return m, downloadBlocks(m.blocks, m.categoryPaths, m.repo.RegistryName)
		}
	case tea.WindowSizeMsg:
		m.table.SetHeight(int(math.Floor(float64(msg.Height) / 3)))
		m.filePicker.SetHeight(int(math.Floor(float64(msg.Height)/3)) - 1)
		m.table.SetWidth(registry_selector.SidebarWidth + 2)
	case manifest.ManifestResponse:
		m.repo = msg
	case block_list.Blocks:
		m.blocks = msg
		for _, block := range m.blocks {
			if _, ok := m.categoryPaths[block.Category]; !ok {
				m.categoryPaths[block.Category] = "./" + block.Category
			}
			for _, item := range m.getLocalDependencies(block) {
				if _, ok := m.categoryPaths[item.Category]; !ok {
					m.categoryPaths[item.Category] = "./" + item.Category
				}
			}
		}
		for category, _ := range m.categoryPaths {
			if !slices.ContainsFunc(m.blocks, func(block manifest.Block) bool {
				return block.Category == category
			}) && !slices.ContainsFunc(m.blocks, func(block manifest.Block) bool {
				for _, item := range m.getLocalDependencies(block) {
					if item.Category == category {
						return true
					}
				}
				return false
			}) {
				delete(m.categoryPaths, category)
			}
		}

		var rows []table.Row
		for category, path := range m.categoryPaths {
			rows = append(rows, table.Row{category, path})
		}
		m.table.SetRows(rows)
	}
	switch m.active {
	case filePickerView:
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			cwd, err := os.Getwd()
			if err != nil {
				return m, nil
			}
			m.filePicker, cmd = m.filePicker.Update(msg)
			if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
				relativePath, err := filepath.Rel(cwd, path)
				if err != nil {
					return m, nil
				}
				m.categoryPaths[m.table.SelectedRow()[0]] = "./" + relativePath
				var rows []table.Row
				for category, path := range m.categoryPaths {
					rows = append(rows, table.Row{category, path})
				}
				m.table.SetRows(rows)
				m.active = tableView
			}
		}
		cmds = append(cmds, cmd)
	case tableView:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var activeView string
	switch m.active {
	case tableView:
		activeView = m.table.View()
	case filePickerView:
		activeView = m.filePicker.View()
	}
	if m.focus {
		return lipgloss.NewStyle().
			Width(m.table.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("140")).
			Render(activeView)
	} else {
		return lipgloss.NewStyle().
			Width(m.table.Width()).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(activeView)
	}
}

func (m *Model) Focus() {
	m.focus = true
	m.table.Focus()
}

func (m *Model) Blur() {
	m.focus = false
	m.table.Blur()
}

func (m Model) getLocalDependencies(selectedBlock manifest.Block) []manifest.Block {
	var localDependencies []manifest.Block
	for _, localDependency := range selectedBlock.LocalDependencies {
		blockName := strings.Split(localDependency, "/")[len(strings.Split(localDependency, "/"))-1]
		for _, category := range m.repo.Categories {
			for _, block := range category.Blocks {
				if block.Name == blockName {
					localDependencies = append(localDependencies, block)
				}
			}
		}
	}
	return localDependencies
}

func downloadBlocks(blocks []manifest.Block, categoryPath map[string]string, registryName string) tea.Cmd {
	return func() tea.Msg {
		var commandString string
		commandString += "npx jsrepo add --tests false --formatter prettier --allow --yes --paths "
		idx := 0
		for categroy, path := range categoryPath {
			//do not add comma on last item
			if len(categoryPath)-1 == idx {
				commandString += fmt.Sprintf("%s=%s ", categroy, path)
			} else {
				commandString += fmt.Sprintf("%s=%s,", categroy, path)
			}
			idx++
		}
		for _, block := range blocks {
			commandString += fmt.Sprintf("%s/%s/%s ", registryName, block.Category, block.Name)
		}
		cmd := exec.Command("sh", "-c", commandString)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return manifest.BannerErrorMessage(err.Error())
		}
		return manifest.BannerErrorMessage(string(output))

	}

}

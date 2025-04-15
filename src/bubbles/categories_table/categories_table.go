package categories_table

import (
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/block_list"
	downloadblocks "github.com/fredrikmwold/jsrepo-tui/src/commands/download_blocks"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
	"github.com/fredrikmwold/jsrepo-tui/src/config"
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
	categoryPaths []downloadblocks.CategoryPath
	blocks        []manifest.Block
	repo          manifest.ManifestResponse
}

func New() Model {
	columns := []table.Column{
		{Title: "Category", Width: config.SidebarWidth/6*2 - 4},
		{Title: "Path", Width: config.SidebarWidth/6*4 + 4},
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
		categoryPaths: []downloadblocks.CategoryPath{},
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
		case "ctrl+d":
			repoName := strings.Split(m.repo.RegistryName, "@")[0]
			return m, downloadblocks.DownloadBlocks(m.blocks, m.categoryPaths, repoName)
		}
	case tea.WindowSizeMsg:
		m.table.SetHeight(int(math.Floor(float64(msg.Height) / 3)))
		m.filePicker.SetHeight(int(math.Floor(float64(msg.Height)/3)) - 1)
		m.table.SetWidth(config.SidebarWidth + 2)
	case manifest.ManifestResponse:
		m.repo = msg
	case block_list.Blocks:
		m.handleBlocks(msg)
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
				for i, category := range m.categoryPaths {
					if category.Category == m.table.SelectedRow()[0] {
						m.categoryPaths[i].Path = "./" + relativePath
						break
					}
				}
				var rows []table.Row
				for _, category := range m.categoryPaths {
					rows = append(rows, table.Row{category.Category, category.Path})
				}
				m.table.SetRows(rows)
				m.filePicker.CurrentDirectory, _ = os.Getwd()
				m.active = tableView
			}
		}
		cmds = append(cmds, cmd)
	case tableView:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

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

func (m *Model) handleBlocks(blocks block_list.Blocks) {
	m.blocks = blocks
	var categoryPaths []downloadblocks.CategoryPath

	for _, block := range m.blocks {
		if !slices.ContainsFunc(categoryPaths, func(cp downloadblocks.CategoryPath) bool {
			return cp.Category == block.Category
		}) {
			categoryPaths = append(categoryPaths, downloadblocks.CategoryPath{
				Category: block.Category,
				Path:     "./" + block.Category,
			})
		}

		for _, item := range m.getLocalDependencies(block) {
			if !slices.ContainsFunc(categoryPaths, func(cp downloadblocks.CategoryPath) bool {
				return cp.Category == item.Category
			}) {
				categoryPaths = append(categoryPaths, downloadblocks.CategoryPath{
					Category: item.Category,
					Path:     "./" + item.Category,
				})
			}
		}
	}

	var rows []table.Row
	for _, cp := range categoryPaths {
		rows = append(rows, table.Row{cp.Category, cp.Path})
	}
	m.table.SetRows(rows)
	m.categoryPaths = categoryPaths
}

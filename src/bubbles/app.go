package app

import (
	tea "github.com/charmbracelet/bubbletea"
	bannermessage "github.com/fredrikmwold/jsrepo-tui/src/bubbles/banner_message"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/block_list"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/categories_table"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/dependency_table"
	downloadspinner "github.com/fredrikmwold/jsrepo-tui/src/bubbles/download_spinner"
	keybindinghelp "github.com/fredrikmwold/jsrepo-tui/src/bubbles/key_binding_help"
	newregistryinput "github.com/fredrikmwold/jsrepo-tui/src/bubbles/new_registry_input"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/registry_selector"
	"github.com/fredrikmwold/jsrepo-tui/src/bubbles/selected_block_list"
	downloadblocks "github.com/fredrikmwold/jsrepo-tui/src/commands/download_blocks"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
	"github.com/fredrikmwold/jsrepo-tui/src/config"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

const (
	selector = iota
	listView
	selectedBlocks
	categoriesTable
	newRegistryInput
)

type Model struct {
	registryselector registry_selector.Model
	blocklist        block_list.Model
	selectedBlocks   selected_block_list.Model
	dependencytable  dependency_table.Model
	newRegistryInput newregistryinput.Model
	categoriestable  categories_table.Model
	bannermessage    bannermessage.Model
	downloadspinner  downloadspinner.Model
	help             help.Model

	keys             keybindinghelp.KeyMap
	isDownloading    bool
	hasBannerMessage bool
	width            int
	height           int
	active           int
}

func New() Model {
	return Model{
		registryselector: registry_selector.New(),
		blocklist:        block_list.New(),
		selectedBlocks:   selected_block_list.New(),
		dependencytable:  dependency_table.New(),
		categoriestable:  categories_table.New(),
		newRegistryInput: newregistryinput.New(),
		bannermessage:    bannermessage.New(),
		downloadspinner:  downloadspinner.New(),
		keys:             keybindinghelp.Keys,
		help:             help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(config.LoadConfig, m.selectedBlocks.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.DownloadBlocks):
			m.isDownloading = true
			m.downloadspinner, cmd = m.downloadspinner.Update(msg)
			cmds = append(cmds, cmd)
			m.categoriestable, cmd = m.categoriestable.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.Tab):
			m.registryselector.Blur()
			m.categoriestable.Blur()
			if m.active == listView {
				m.active = selectedBlocks
				m.blocklist.Blur()
				m.selectedBlocks.Focus()
			} else {
				m.active = listView
				m.blocklist.Focus()
				m.selectedBlocks.Blur()
			}
		case key.Matches(msg, m.keys.S):
			if m.active != newRegistryInput {
				m.active = selector
				m.registryselector.Focus()
				m.blocklist.Blur()
				m.categoriestable.Blur()
				m.selectedBlocks.Blur()
			}
		case key.Matches(msg, m.keys.P):
			if m.active != newRegistryInput {
				m.active = categoriesTable
				m.categoriestable.Focus()
				m.registryselector.Blur()
				m.blocklist.Blur()
				m.selectedBlocks.Blur()
				return m, nil
			}
		case key.Matches(msg, m.keys.AddNewRegistry):
			if m.active != newRegistryInput {
				m.active = newRegistryInput
				m.active = newRegistryInput
				m.registryselector.Blur()
				m.blocklist.Blur()
				m.categoriestable.Blur()
				m.selectedBlocks.Blur()
				return m, nil
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEsc:
			if m.hasBannerMessage {
				m.hasBannerMessage = false
				return m, nil
			}
			if m.active == newRegistryInput {
				m.active = selector
				m.registryselector.Focus()
				m.blocklist.Blur()
				m.selectedBlocks.Blur()
				return m, nil
			}
		}
	case block_list.Blocks:
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.categoriestable, _ = m.categoriestable.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		msg.Height = msg.Height - 1
		m.registryselector, _ = m.registryselector.Update(msg)
		m.selectedBlocks, _ = m.selectedBlocks.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.dependencytable, _ = m.dependencytable.Update(msg)
		m.categoriestable, _ = m.categoriestable.Update(msg)
		m.newRegistryInput, _ = m.newRegistryInput.Update(msg)
		m.bannermessage, _ = m.bannermessage.Update(msg)
		m.downloadspinner, _ = m.downloadspinner.Update(msg)
		return m, nil
	case manifest.ManifestResponse:
		m.registryselector, _ = m.registryselector.Update(msg)
		m.blocklist, _ = m.blocklist.Update(msg)
		m.categoriestable, _ = m.categoriestable.Update(msg)
		if msg.Categories == nil {
			return m, nil
		}
		m.active = listView
		m.blocklist.Focus()
		m.registryselector.Blur()
		return m, nil
	case manifest.ManifestErrorMessage:
		m.hasBannerMessage = true
		m.bannermessage, _ = m.bannermessage.Update(msg)
		return m, func() tea.Msg {
			return manifest.ManifestResponse{}
		}
	case downloadblocks.DownloadBlocksErrorMessage:
		m.hasBannerMessage = true
		m.isDownloading = false
		m.bannermessage, _ = m.bannermessage.Update(msg)
	case downloadblocks.SuccessMessage:
		m.hasBannerMessage = true
		m.isDownloading = false
		m.bannermessage, _ = m.bannermessage.Update(msg)
		return m, nil
	case config.Config:
		m.registryselector, _ = m.registryselector.Update(msg)
		m.active = selector
		m.registryselector.Focus()
		return m, nil
	}

	m.downloadspinner, cmd = m.downloadspinner.Update(msg)
	cmds = append(cmds, cmd)

	switch m.active {
	case selector:
		m.registryselector, cmd = m.registryselector.Update(msg)
		cmds = append(cmds, cmd)
	case listView:
		m.blocklist, cmd = m.blocklist.Update(msg)
		cmds = append(cmds, cmd)
	case selectedBlocks:
		m.selectedBlocks, cmd = m.selectedBlocks.Update(msg)
		cmds = append(cmds, cmd)
	case newRegistryInput:
		m.newRegistryInput, cmd = m.newRegistryInput.Update(msg)
		cmds = append(cmds, cmd)
	case categoriesTable:
		m.categoriestable, cmd = m.categoriestable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	if m.width <= 110 || m.height <= 25 {
		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"Terminal window is too small \n",
					"Please resize the terminal window to at least 110x25"),
			)

	}

	var bannerView string
	bannerViewHeight := 0

	helpView := m.help.View(m.keys)

	if m.active == newRegistryInput {
		bannerView = m.newRegistryInput.View()
		bannerViewHeight = lipgloss.Height(bannerView)
	}
	if m.hasBannerMessage {
		bannerView = m.bannermessage.View()
		bannerViewHeight = lipgloss.Height(bannerView)
	}
	if m.isDownloading {
		bannerView = m.downloadspinner.View()
		bannerViewHeight = lipgloss.Height(bannerView)
	}

	m.selectedBlocks.SetHeight(lipgloss.Height(m.selectedBlocks.View()) - bannerViewHeight)
	m.blocklist.SetHeight(lipgloss.Height(m.blocklist.View()) - bannerViewHeight)

	sidebar := lipgloss.JoinVertical(
		lipgloss.Left,
		m.registryselector.View(),
		m.categoriestable.View(),
		m.dependencytable.View(),
	)

	dashboard := lipgloss.JoinVertical(lipgloss.Left,
		bannerView,
		lipgloss.JoinHorizontal(lipgloss.Bottom,
			m.blocklist.View(),
			m.selectedBlocks.View(),
		))
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			sidebar,
			dashboard,
		),
		helpView)

}

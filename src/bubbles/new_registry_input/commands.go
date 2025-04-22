package newregistryinput

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
)

func duplicateRegistryErrorMessage() tea.Msg {
	return manifest.ManifestErrorMessage("You have already added this registry")
}

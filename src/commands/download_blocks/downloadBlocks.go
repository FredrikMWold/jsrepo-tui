package downloadblocks

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikmwold/jsrepo-tui/src/commands/manifest"
)

type SuccessMessage string
type DownloadBlocksErrorMessage string

type CategoryPath struct {
	Category string
	Path     string
}

func DownloadBlocks(blocks []manifest.Block, categoryPaths []CategoryPath, registryName string) tea.Cmd {
	return func() tea.Msg {
		var commandString string
		commandString += "npx --no-install jsrepo add --tests false --formatter prettier --allow --yes --paths "
		idx := 0
		for _, category := range categoryPaths {
			//do not add comma on last item
			if len(categoryPaths)-1 == idx {
				commandString += fmt.Sprintf("%s=%s ", category.Category, category.Path)
			} else {
				commandString += fmt.Sprintf("%s=%s,", category.Category, category.Path)
			}
			idx++
		}
		for _, block := range blocks {
			commandString += fmt.Sprintf("%s/%s/%s ", registryName, block.Category, block.Name)
		}
		cmd := exec.Command("sh", "-c", commandString)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return DownloadBlocksErrorMessage(err.Error())
		}
		return SuccessMessage("Downloaded blocks successfully")

	}

}

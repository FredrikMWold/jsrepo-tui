package downloadblocks

import (
	"fmt"
	"jsrepo-tui/src/commands/manifest"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func DownloadBlocks(blocks []manifest.Block, categoryPath map[string]string, registryName string) tea.Cmd {
	return func() tea.Msg {
		var commandString string
		commandString += "npx --no-install jsrepo add --tests false --formatter prettier --allow --yes --paths "
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

package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func GetManifest(provider string) tea.Cmd {
	return func() tea.Msg {
		url := strings.Replace(provider, "github", "https://raw.githubusercontent.com/", 1)
		url = fmt.Sprintf("%s%s", url, "/refs/heads/master/jsrepo-manifest.json")

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error getting manifest: %v\n", err)
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var response ManifestResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return err
		}
		return response
	}
}
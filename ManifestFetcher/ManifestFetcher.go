package manifestfetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Close struct{}

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

type Block struct {
	Name              string            `json:"name"`
	Directory         string            `json:"directory"`
	Category          string            `json:"category"`
	Tests             bool              `json:"tests"`
	Subdirectory      bool              `json:"subdirectory"`
	List              bool              `json:"list"`
	Files             []string          `json:"files"`
	LocalDependencies []string          `json:"localDependencies"`
	Dependencies      []string          `json:"dependencies"`
	DevDependencies   []string          `json:"devDependencies"`
	Imports           map[string]string `json:"_imports_"`
}

type Category struct {
	Name   string  `json:"name"`
	Blocks []Block `json:"blocks"`
}

type ManifestResponse struct {
	Categories []Category `json:"categories"`
}

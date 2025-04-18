package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ManifestErrorMessage string

func GetManifest(provider string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Split(provider, "@")
		providerURL := parts[0]
		branch := "master" // default branch
		if len(parts) > 1 {
			branch = parts[1]
		}
		url := strings.Replace(providerURL, "github", "https://raw.githubusercontent.com", 1)
		url = fmt.Sprintf("%s%s", url, "/refs/heads/"+branch+"/jsrepo-manifest.json")
		resp, err := http.Get(url)
		if err != nil {
			return ManifestErrorMessage("Error getting manifest from: " + provider)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err)
			return ManifestErrorMessage(err.Error())
		}

		if resp.StatusCode == 404 {
			return ManifestErrorMessage("Error getting manifest from: " + provider)
		}

		var response ManifestResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return err
		}
		response.RegistryName = provider
		return response
	}
}

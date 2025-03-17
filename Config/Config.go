package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ConfigEntry struct {
	Schema       string            `json:"$schema"`
	IncludeTests bool              `json:"includeTests"`
	Watermark    bool              `json:"watermark"`
	Paths        map[string]string `json:"paths"`
	Repos        []string          `json:"repos"`
}

type Config struct {
	Entries map[string]ConfigEntry `json:"entries"`
}

func LoadConfig() tea.Msg {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config dir:", err)
		return err
	}

	data, err := os.ReadFile(configDir + "/jsrepo-nodejs/config.json")

	if err != nil {
		fmt.Println("Error reading config file:", err)
		return err
	}

	var rawConfig map[string]json.RawMessage

	if err := json.Unmarshal(data, &rawConfig); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	loadedconfig := Config{
		Entries: make(map[string]ConfigEntry),
	}

	for key, value := range rawConfig {
		if key == "latest-version" {
			continue
		}
		if !strings.Contains(key, "-state") {
			continue
		}
		var entry ConfigEntry
		if err := json.Unmarshal(value, &entry); err != nil {
			fmt.Printf("Error unmarshalling entry for key %s: %v\n", key, err)
			return err
		}
		loadedconfig.Entries[strings.Replace(key, "-state", "", 1)] = entry
	}
	fmt.Printf("Loaded config: %v\n", loadedconfig)
	return loadedconfig
}

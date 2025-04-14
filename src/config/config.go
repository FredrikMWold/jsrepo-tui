package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

const SidebarWidth = 38

type Config struct {
	Registries []string `mapstructure:"registries"`
}

func LoadConfig() tea.Msg {
	var config Config
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}

	configPath := filepath.Join(configDir)
	viper.AddConfigPath(configPath)
	viper.SetConfigName("jsrepo-tui.json")
	viper.SetConfigType("json")
	viper.SetDefault("registries", []string{})

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatalf("Error writing default config to file: %v", err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}
	return config
}

type JsrepoCacheEntry struct {
	Schema       string            `json:"$schema"`
	IncludeTests bool              `json:"includeTests"`
	Watermark    bool              `json:"watermark"`
	Paths        map[string]string `json:"paths"`
	Repos        []string          `json:"repos"`
}

type JsrepoCache struct {
	Entries map[string]JsrepoCacheEntry `json:"entries"`
}

func LoadConfig2() tea.Msg {
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

	loadedconfig := JsrepoCache{
		Entries: make(map[string]JsrepoCacheEntry),
	}

	for key, value := range rawConfig {
		if key == "latest-version" {
			continue
		}
		if !strings.Contains(key, "-state") {
			continue
		}
		var entry JsrepoCacheEntry
		if err := json.Unmarshal(value, &entry); err != nil {
			fmt.Printf("Error unmarshalling entry for key %s: %v\n", key, err)
			return err
		}
		loadedconfig.Entries[strings.Replace(key, "-state", "", 1)] = entry
	}
	return loadedconfig
}

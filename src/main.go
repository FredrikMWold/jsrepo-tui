package main

import (
	"fmt"
	app "jsrepo-tui/src/bubbles"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := app.New()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	app "github.com/fredrikmwold/jsrepo-tui/src/bubbles"
)

func main() {
	model := app.New()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

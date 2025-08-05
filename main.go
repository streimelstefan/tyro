package main

import (
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/streimelstefan/tyro/ui"
)

func main() {
	// Check if directory argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: tyro <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	// Initialize the Bubble Tea program
	app := ui.NewApp(dir)
	p := tea.NewProgram(app, tea.WithAltScreen())

	log.SetOutput(io.Discard)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

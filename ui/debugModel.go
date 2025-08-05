package ui

import (
	"fmt"
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type debugModel struct {
	debug  bool
	dump   io.Writer
	handle *os.File
}

func NewDebugModel() *debugModel {
	if _, ok := os.LookupEnv("DEBUG"); !ok {
		return &debugModel{
			debug: false,
		}
	}

	dump, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Printf("failed to open debug log file: %v\n", err.Error())
		os.Exit(1)
	}

	spew.Fdump(dump, time.Now().Format(time.RFC3339))

	return &debugModel{
		debug:  true,
		dump:   dump,
		handle: dump,
	}
}

func (m *debugModel) Init() tea.Cmd {
	return nil
}

func (m *debugModel) Update(msg tea.Msg) (*debugModel, tea.Cmd) {
	if m.debug {
		spew.Fdump(m.dump, msg)
	}
	return m, nil
}

func (m *debugModel) View() string {
	return ""
}

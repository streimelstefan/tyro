package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	defaults "github.com/streimelstefan/tyro/ui/defaults"
)

type Model struct {
	Folder string

	Style *StatusBarStyle

	width int
}

type StatusBarStyle struct {
	FolderStyle lipgloss.Style
}

func New(folder string) *Model {
	return &Model{
		Folder: folder,
		Style: &StatusBarStyle{
			FolderStyle: lipgloss.NewStyle().
				Foreground(defaults.TextColor).
				Background(defaults.AccentColor).
				Padding(0, 0, 0, 0),
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return ""
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	return m, nil
}

package ui

import tea "github.com/charmbracelet/bubbletea"

type FileTreeItemModel struct {
	Part string
}

func NewFileTreeItemModel(part string) FileTreeItemModel {
	return FileTreeItemModel{
		Part: part,
	}
}

func (m FileTreeItemModel) Init() tea.Cmd {
	return nil
}

func (m FileTreeItemModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m FileTreeItemModel) View() string {
	return m.Part
}

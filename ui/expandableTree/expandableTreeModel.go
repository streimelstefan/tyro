package expandableTree

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ExpandableTree *ExpandableTree

	Spacer    string
	Expanded  string
	Branch    string
	Collapsed string
	BranchEnd string
}

func New() *Model {
	return &Model{
		ExpandableTree: NewExpandableTree(),
		Spacer:         "│  ",
		Branch:         "├─ ",
		Expanded:       "",
		Collapsed:      "+ ",
		BranchEnd:      "└─ ",
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	b := strings.Builder{}

	m.renderTreeNode(m.ExpandableTree.Root, m.ExpandableTree.Root.HasChildren(), &b)

	return b.String()
}

func (m Model) renderTreeNode(node *node, isLast bool, b *strings.Builder) {
	if node.IsFilteredOut {
		return
	}

	if node.level > 1 {
		b.WriteString(strings.Repeat(m.Spacer, node.level-1))
		if isLast && (!node.IsExpanded || !node.HasChildren()) {
			b.WriteString(m.BranchEnd)
		} else {
			b.WriteString(m.Branch)
		}
	} else if node.level == 1 {
		if isLast && (!node.IsExpanded || !node.HasChildren()) {
			b.WriteString(m.BranchEnd)
		} else {
			b.WriteString(m.Branch)
		}
	}

	if node.HasChildren() && !node.isRoot {
		if node.IsExpanded {
			b.WriteString(m.Expanded)
		} else {
			b.WriteString(m.Collapsed)
		}
	}

	b.WriteString(node.Model.View())
	b.WriteRune('\n')

	if node.HasChildren() && node.IsExpanded {
		for i, child := range node.Children {
			m.renderTreeNode(child, i == len(node.Children)-1, b)
		}
	}
}

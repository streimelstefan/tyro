package expandableTree

import (
	tea "github.com/charmbracelet/bubbletea"
)

type rootNodeModel struct {
}

func (m rootNodeModel) Init() tea.Cmd {
	return nil
}

func (m rootNodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m rootNodeModel) View() string {
	return ""
}

type ExpandableTree struct {
	Root *node
}

type node struct {
	Identifier string
	Model      tea.Model
	Children   []*node

	IsExpanded bool
	IsSelected bool

	IsFilteredOut bool

	isRoot bool
	level  int
}

func NewExpandableTree() *ExpandableTree {
	return &ExpandableTree{
		Root: &node{
			Identifier:    "",
			Model:         rootNodeModel{},
			Children:      make([]*node, 0),
			IsExpanded:    true,
			IsSelected:    false,
			IsFilteredOut: false,
			isRoot:        true,
			level:         -1,
		},
	}
}

func newNode(identifier string, model tea.Model, level int) *node {
	return &node{
		Model:      model,
		Children:   make([]*node, 0),
		IsExpanded: true,
		IsSelected: false,
		isRoot:     false,
		level:      level,
		Identifier: identifier,
	}
}

func (e *ExpandableTree) AddNode(parent *node, identifier string, model tea.Model) *node {
	newNode := newNode(identifier, model, parent.level+1)
	parent.Children = append(parent.Children, newNode)
	return newNode
}

func (n node) HasChildren() bool {
	return len(n.Children) > 0
}

func (n node) HasChild(identifier string) bool {
	return n.GetChild(identifier) != nil
}

func (n node) GetChild(identifier string) *node {
	for _, child := range n.Children {
		if child.Identifier == identifier {
			return child
		}
	}
	return nil
}

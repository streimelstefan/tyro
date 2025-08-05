package ui

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/streimelstefan/tyro/ui/expandableTree"
	"github.com/streimelstefan/tyro/ui/statusbar"
)

// App represents the main application state
type App struct {
	statusBar *statusbar.Model

	discovery *discoveryModel

	width  int
	height int

	fileTree         *expandableTree.Model
	fileTreeViewPort viewport.Model

	debug *debugModel
}

// NewApp creates a new application instance
func NewApp(folder string) App {
	return App{
		statusBar: statusbar.New(folder),
		discovery: NewDiscoveryModel(folder, 100*time.Millisecond),
		fileTree:  expandableTree.New(),
		debug:     NewDebugModel(),
	}
}

// Init is called when the program starts
func (m App) Init() tea.Cmd {
	return m.discovery.Init()
}

// Update handles messages and user input
func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.fileTreeViewPort.Width = msg.Width
		m.fileTreeViewPort.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case CollectedDICOMFiles:
		m.addNewFilesToTrees(msg)
		m.fileTreeViewPort.SetContent(m.fileTree.View())
	}

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	m.discovery, cmd = m.discovery.Update(msg)
	cmds = append(cmds, cmd)

	m.fileTreeViewPort, cmd = m.fileTreeViewPort.Update(msg)
	cmds = append(cmds, cmd)

	m.debug, cmd = m.debug.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m App) View() string {
	return m.fileTreeViewPort.View()
}

func (m App) addNewFilesToTrees(files CollectedDICOMFiles) {
	for _, file := range files {
		rel, err := filepath.Rel(m.discovery.rootDir, file.Path)
		if err != nil {
			continue
		}

		parts := strings.Split(rel, string(filepath.Separator))

		currentNode := m.fileTree.ExpandableTree.Root
		for _, part := range parts {
			tmpChild := currentNode.GetChild(part)
			if tmpChild == nil {
				tmpChild = m.fileTree.ExpandableTree.AddNode(currentNode, part, NewFileTreeItemModel(part))
			}
			currentNode = tmpChild
		}
	}
}

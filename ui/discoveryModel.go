package ui

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/streimelstefan/tyro/operations"
)

type DiscoveryCollectMsg struct{}

type CollectedDICOMFiles []*operations.ParsedDicomFile

func NewDiscoveryModel(rootDir string, batchDelay time.Duration) *discoveryModel {
	return &discoveryModel{
		rootDir:                 rootDir,
		batchDelay:              batchDelay,
		collectedDiscoveryFiles: make([]*operations.ParsedDicomFile, 0),
		discoveryErrors:         make([]error, 0),
		discoveryInProgress:     false,
	}
}

type discoveryModel struct {
	rootDir    string
	batchDelay time.Duration

	collectedDiscoveryFiles []*operations.ParsedDicomFile
	discoveryErrors         []error

	discoveryInProgress bool
	collectionFinished  bool
	discoveryMutex      sync.Mutex
	discoveryErrorMutex sync.Mutex
}

func (s *discoveryModel) Init() tea.Cmd {
	return s.discoverFiles()
}

func (s *discoveryModel) Update(msg tea.Msg) (*discoveryModel, tea.Cmd) {
	switch msg.(type) {
	case DiscoveryCollectMsg:
		if s.discoveryInProgress || !s.collectionFinished {
			return s, tea.Batch(s.tickDiscovery(), s.collectFiles)
		}
	}
	return s, nil
}

func (s *discoveryModel) View() string {
	return ""
}

func (s *discoveryModel) discoverFiles() tea.Cmd {
	if s.discoveryInProgress {
		return nil
	}

	discoveryResult := operations.DiscoverDICOMFiles(s.rootDir, 8)

	parseResults := operations.ParseDICOMFiles(discoveryResult.Files, 8)

	// go routine to accumulate all errors
	go func() {
		discoveryFinished := false
		parseFinished := false

		for !discoveryFinished && !parseFinished {
			select {
			case err, ok := <-discoveryResult.Errors:
				if !ok {
					discoveryFinished = true
					continue
				}

				s.addDiscoveryError(err)
			case err, ok := <-parseResults.Errors:
				if !ok {
					parseFinished = true
					continue
				}

				s.addDiscoveryError(err)
			}
		}
	}()

	// go routine to accumulate all files
	go func() {
		for file := range parseResults.Files {
			s.addFileToCollection(file)
		}
	}()

	return s.tickDiscovery()
}

func (s *discoveryModel) addFileToCollection(file *operations.ParsedDicomFile) {
	s.discoveryMutex.Lock()
	s.collectedDiscoveryFiles = append(s.collectedDiscoveryFiles, file)
	s.discoveryMutex.Unlock()
}

func (s *discoveryModel) addDiscoveryError(err error) {
	s.discoveryErrorMutex.Lock()
	s.discoveryErrors = append(s.discoveryErrors, err)
	s.discoveryErrorMutex.Unlock()
}

func (s *discoveryModel) collectFiles() tea.Msg {
	collectedFiles := s.collectedDiscoveryFiles

	s.discoveryMutex.Lock()
	s.collectedDiscoveryFiles = make([]*operations.ParsedDicomFile, 0)
	s.discoveryMutex.Unlock()

	if len(collectedFiles) == 0 {
		s.collectionFinished = true
	}

	return CollectedDICOMFiles(collectedFiles)
}

func (s *discoveryModel) tickDiscovery() tea.Cmd {
	return tea.Tick(s.batchDelay, func(t time.Time) tea.Msg {
		return DiscoveryCollectMsg{}
	})
}

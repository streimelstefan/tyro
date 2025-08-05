// Package main provides utilities for discovering DICOM files in a directory tree.
//
// This package includes functions to identify valid DICOM files and recursively search directories
// with configurable concurrency. It is designed for efficient and robust DICOM file discovery in large
// file systems, handling errors gracefully and supporting concurrent processing.
package operations

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	// ErrorFileTooSmallToBeDICOM is returned when a file is too small to be a valid DICOM file.
	ErrorFileTooSmallToBeDICOM = errors.New("file too small to be a valid DICOM")
	// ErrorInvalidMagicNumber is returned when a file does not have the DICOM magic number.
	ErrorInvalidMagicNumber = errors.New("invalid magic number")
)

// DicomFile represents a discovered DICOM file and its open file handle.
type DicomFile struct {
	// Path is the filesystem location of the DICOM file.
	Path string
	// Handle is the open file handle for the DICOM file.
	Handle *os.File
}

// DiscoveryResult contains the channels for discovered DICOM files and errors.
type DiscoveryResult struct {
	// Files is a channel that will receive discovered DicomFile objects.
	Files <-chan DicomFile
	// Errors is a channel that will receive errors encountered during discovery.
	Errors <-chan error
}

// isValidDICOM checks if the file at the given path is a valid DICOM file.
//
// It returns true and an open file handle if the file is a valid DICOM file, otherwise false.
// If the file is not valid, the returned file handle will be nil. If an error occurs during
// reading, it is returned. The caller is responsible for closing the returned file handle.
func isValidDICOM(path string) (bool, *os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return false, nil, err
	}

	// DICOM files have a 128-byte preamble followed by "DICM"
	header := make([]byte, 132)
	n, err := file.Read(header)
	if err != nil {
		file.Close()
		return false, nil, err
	}
	if n < 132 {
		file.Close()
		return false, nil, ErrorFileTooSmallToBeDICOM
	}

	if string(header[128:132]) != "DICM" {
		file.Close()
		return false, nil, ErrorInvalidMagicNumber
	}

	// reset the file pointer to the beginning of the file
	// if this is not done, the parse will break with an eof error
	file.Seek(0, io.SeekStart)
	return true, file, nil
}

// DiscoverDICOMFiles scans the given directory and returns channels for discovered DICOM files and errors.
// This function allows for streaming processing of discovered files without waiting for all files to be found.
//
// dir specifies the root directory to search for DICOM files.
// maxConcurrency sets the maximum number of concurrent goroutines allowed (if 0, defaults to 8).
//
// Returns a DiscoveryResult containing channels for discovered files and errors.
// The caller is responsible for reading from both channels until they are closed.
func DiscoverDICOMFiles(dir string, maxConcurrency int) DiscoveryResult {
	if maxConcurrency <= 0 {
		maxConcurrency = 8
	}

	fileCh := make(chan string, maxConcurrency*2)
	resultCh := make(chan DicomFile, maxConcurrency*2)
	errCh := make(chan error, maxConcurrency*2)
	var wg sync.WaitGroup

	// Start the directory traversal goroutine.
	go fileWalker(dir, fileCh, errCh)

	// Start the worker pool for DICOM validation.
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dicomCheckerWorker(fileCh, resultCh, errCh)
		}()
	}

	// Close resultCh and errCh when all workers are done.
	go func() {
		wg.Wait()
		close(errCh)
		close(resultCh)
	}()

	return DiscoveryResult{
		Files:  resultCh,
		Errors: errCh,
	}
}

// fileWalker walks the directory tree rooted at dir and sends file paths to fileCh.
//
// Any errors encountered during traversal are sent to errCh. fileCh is closed when traversal is complete.
func fileWalker(dir string, fileCh chan<- string, errCh chan<- error) {
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			errCh <- err
		}
		if !d.IsDir() {
			fileCh <- path
		}
		return nil
	})
	if err != nil {
		errCh <- err
	}
	close(fileCh)
}

// dicomCheckerWorker receives file paths from fileCh, checks if they are valid DICOM files,
// and sends valid DicomFile objects to resultCh. Errors encountered during validation are sent to errCh.
func dicomCheckerWorker(fileCh <-chan string, resultCh chan<- DicomFile, errCh chan<- error) {
	for path := range fileCh {
		isValid, handle, err := isValidDICOM(path)
		if err != nil {
			errCh <- err
			continue
		}
		if isValid {
			resultCh <- DicomFile{Path: path, Handle: handle}
		}
	}
}

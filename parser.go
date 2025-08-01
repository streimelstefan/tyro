// Package main provides utilities for parsing DICOM files discovered by the discovery module.
//
// This package includes functions to parse DICOM files concurrently using worker pools,
// designed to work with the channel-based discovery system for efficient streaming processing.
// The parser automatically handles file handle cleanup, closing handles on both successful
// parsing and errors to prevent resource leaks.
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/suyashkumar/dicom"
)

// ParsedDicomFile represents a successfully parsed DICOM file with its dataset.
//
// This struct contains the parsed DICOM dataset along with the original file path and
// an open file handle. The caller is responsible for closing the file handle after
// processing the parsed data.
type ParsedDicomFile struct {
	// Path is the filesystem location of the DICOM file.
	Path string
	// Dataset contains the parsed DICOM dataset with all elements and metadata.
	Dataset dicom.Dataset
	// handle is the open file handle for the DICOM file.
	// The caller must close this handle after processing.
	handle *os.File

	// isOpen is a flag to indicate if the file handle is still open.
	isOpen bool
}

func (p *ParsedDicomFile) String() string {
	return fmt.Sprintf("ParsedDicomFile{Path: %s, Dataset: %d elements}", p.Path, len(p.Dataset.Elements))
}

func (p *ParsedDicomFile) GetHandle() (*os.File, error) {
	if !p.isOpen {
		handle, err := os.Open(p.Path)
		if err != nil {
			return nil, err
		}
		p.handle = handle
		p.isOpen = true
	}

	return p.handle, nil
}

func (p *ParsedDicomFile) Close() error {
	if !p.isOpen {
		return nil
	}
	p.isOpen = false
	return p.handle.Close()
}

// ParsingResult contains the channels for parsed DICOM files and errors.
//
// This struct provides access to the output channels from the parsing process,
// allowing the caller to receive parsed files and errors as they become available.
type ParsingResult struct {
	// Files is a channel that will receive parsed ParsedDicomFile objects.
	// This channel will be closed when all parsing is complete.
	Files <-chan *ParsedDicomFile
	// Errors is a channel that will receive errors encountered during parsing.
	// This channel will be closed when all parsing is complete.
	Errors <-chan error
}

// ParseDICOMFiles takes a channel of discovered DICOM files and returns channels
// for parsed DICOM files and parsing errors. This function allows for parallel parsing
// of discovered files using a configurable worker pool.
//
// dicomChannel supplies DicomFile objects from the discovery process.
// maxConcurrency sets the maximum number of concurrent parsing goroutines (if 0, defaults to 8).
//
// Returns a ParsingResult containing channels for parsed files and parsing errors.
// The caller is responsible for reading from both channels until they are closed.
// The function will close the output channels when all input channels are closed and all parsing is complete.
// File handles are automatically closed on parsing errors to prevent resource leaks.
func ParseDICOMFiles(dicomChannel <-chan DicomFile, maxConcurrency int) ParsingResult {
	if maxConcurrency <= 0 {
		maxConcurrency = 8
	}

	// Increased buffer sizes for better performance with large datasets
	resultCh := make(chan *ParsedDicomFile, maxConcurrency*4)
	errCh := make(chan error, maxConcurrency*4)
	var wg sync.WaitGroup

	// Start the worker pool for DICOM parsing.
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dicomParserWorker(dicomChannel, resultCh, errCh)
		}()
	}

	// Close resultCh and errCh when all workers are done.
	go func() {
		wg.Wait()
		close(errCh)
		close(resultCh)
	}()

	return ParsingResult{
		Files:  resultCh,
		Errors: errCh,
	}
}

// dicomParserWorker receives DicomFile objects from fileCh, parses them, and sends
// ParsedDicomFile objects to resultCh. Errors encountered during parsing are sent to errCh.
//
// This worker function runs in a goroutine and processes DICOM files concurrently.
// It automatically closes file handles on parsing errors to prevent resource leaks.
// Successfully parsed files are sent to resultCh with their handles still open
// for further processing by the caller.
// The function uses saveParseUntilEOF to handle any panics from the DICOM parsing library.
func dicomParserWorker(fileCh <-chan DicomFile, resultCh chan<- *ParsedDicomFile, errCh chan<- error) {
	for file := range fileCh {
		// Use a panic recovery wrapper to handle any panics from ParseUntilEOF
		dataset, err := saveParseUntilEOF(file.Handle)

		if err != nil {
			errCh <- err
			file.Handle.Close()
			continue
		}
		resultCh <- &ParsedDicomFile{
			Path:    file.Path,
			Dataset: dataset,
			handle:  file.Handle,
			isOpen:  true,
		}
	}
}

// saveParseUntilEOF safely parses a DICOM file with panic recovery.
//
// This function wraps the dicom.ParseUntilEOF call with panic recovery to handle
// any unexpected panics from the DICOM parsing library. Panics are converted to
// regular errors that can be handled by the calling code.
//
// file is the open file handle to parse. The file should be positioned at the beginning.
//
// Returns the parsed DICOM dataset and any error encountered during parsing.
// If a panic occurs, it is converted to an error with a descriptive message.
func saveParseUntilEOF(file *os.File) (dataset dicom.Dataset, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			if panicErr, ok := r.(error); ok {
				err = panicErr
			} else {
				err = fmt.Errorf("panic during DICOM parsing: %v", r)
			}
		}
	}()

	dataset, err = dicom.ParseUntilEOF(file, nil, dicom.ParseOption(dicom.SkipPixelData()))
	if err != nil {
		return dicom.Dataset{}, err
	}
	return dataset, nil
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// get the first argument as this will be the directory to search
	if len(os.Args) < 2 {
		fmt.Println("Usage: tyro <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	// Use buffered output for better performance
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Discover DICOM files asynchronously
	discoveryResult := DiscoverDICOMFiles(dir, 16) // Increased concurrency

	// Parse discovered DICOM files asynchronously
	parsingResult := ParseDICOMFiles(discoveryResult.Files, 16) // Increased concurrency

	// Process results from both channels
	parsedCount := 0
	errorCount := 0

	// Use a select to read from both channels until they're closed
	filesClosed := false
	discoveryErrorsClosed := false
	parsingErrorsClosed := false

	// multiErr := multierror.New()

	for !filesClosed || !discoveryErrorsClosed || !parsingErrorsClosed {
		select {
		case parsedFile, ok := <-parsingResult.Files:
			if !ok {
				filesClosed = true
				continue
			}
			parsedCount++

			// Use buffered output for better performance
			fmt.Fprintf(writer, "Parsed %s: %d elements\n", parsedFile.Path, len(parsedFile.Dataset.Elements))

			// Close the file handle immediately after processing
			parsedFile.Close()

		case _, ok := <-parsingResult.Errors:
			if !ok {
				parsingErrorsClosed = true
				continue
			}
			errorCount++
			// multiErr.Add(err)

		case _, ok := <-discoveryResult.Errors:
			if !ok {
				discoveryErrorsClosed = true
				continue
			}
			errorCount++
			// multiErr.Add(err)
		}
	}

	writer.Flush()
	fmt.Printf("\nProcessing complete. Parsed %d files with %d errors.\n", parsedCount, errorCount)
}

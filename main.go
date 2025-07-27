package main

import (
	"fmt"
	"os"

	"github.com/suyashkumar/dicom"
)

func main() {
	// get the first argument as this will be the directory to search
	if len(os.Args) < 2 {
		fmt.Println("Usage: tyro <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	// get all files in the directory
	dicomFiles, err := DiscoverDICOMFiles(dir, 8)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for _, file := range dicomFiles {
		dataset, err := dicom.ParseUntilEOF(file.Handle, nil, dicom.ParseOption(dicom.SkipPixelData()))
		if err != nil {
			fmt.Println("Error parsing DICOM file:", err)
			continue
		}
		fmt.Printf("Found %d elements in %s\n", len(dataset.Elements), file.Path)
	}

}

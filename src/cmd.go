package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func convertPdfToText(filePath string) (string, error) {
	textFile, err := createTrackedTempFile(tempFilePatternPDFToText)

	if err != nil {
		log.Printf("Error while creating temporary .txt file: %v\n", err)
		return "", err
	}

	defer textFile.Close()

	output, err := exec.Command("pdftotext", "-layout", filePath, textFile.Name()).CombinedOutput()

	if err != nil {
		log.Printf("Error while converting PDF to text: %v. %s\n", err, output)
		return "", err
	}

	return textFile.Name(), nil
}

func openFileWithDefaultApp(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

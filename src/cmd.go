package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func convertPdfToText(filePath string) string {
	textFilePath := strings.Replace(filepath.Base(filePath), ".pdf", ".txt", 1)
	output, err := exec.Command("pdftotext", "-layout", filePath, textFilePath).CombinedOutput()

	if err != nil {
		panic(fmt.Sprintf("Error while converting PDF to text: %v. %s", err, output))
	}

	return textFilePath
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

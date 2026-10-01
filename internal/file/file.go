package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

const (
	pdfToTextTempFilePattern  = "banksight-converted-pdf-to-text-*.txt"
	StatementTempFilePattern  = "banksight-statement-*.html"
	MT940ToUTFTempFilePattern = "banksight-converted-mt940-to-utf8-*.txt"
	MBankToUTFTempFilePattern = "banksight-converted-mbank-to-utf8-*.csv"
)

func ReadAllLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error while opening the file: %v\n", err)
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error while reading the file: %v\n", err)
		return nil, err
	}

	if len(lines) == 0 {
		log.Println("Error while reading the file: No lines were read from the file.")
		return nil, err
	}

	f.Close()
	RemoveTrackedTempFile(filePath)

	return lines, nil
}

func ConvertToUtf8(filePath string, decoder *encoding.Decoder, tempFilePattern string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error while opening the file: %v\n", err)
		return "", err
	}

	defer file.Close()

	reader := transform.NewReader(file, decoder)
	outputFile, err := CreateTrackedTempFile(tempFilePattern)

	if err != nil {
		log.Printf("Error while creating the file: %v\n", err)
		return "", err
	}

	defer outputFile.Close()
	outputFilePath := outputFile.Name()

	if _, err := io.Copy(outputFile, reader); err != nil {
		log.Printf("Error while copying the file content: %v\n", err)
		return "", err
	}

	return outputFilePath, nil
}

func ConvertPdfToText(filePath string) (string, error) {
	textFile, err := CreateTrackedTempFile(pdfToTextTempFilePattern)

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

func OpenWithDefaultApp(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath) // this doesn't work on windows
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// or something here doesn't work on windows
func CreateTrackedTempFile(pattern string) (*os.File, error) {
	return os.CreateTemp("", pattern)
}

func RemoveTrackedTempFile(filePath string) {
	err := os.Remove(filePath)

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Printf("WARNING: could not delete temporary file '%s': %v\n", filePath, err)
	}
}

func CleanupTempFiles() {
	for _, pattern := range tempFilePatterns() {
		cleanupByPattern(pattern)
	}
}

func tempFilePatterns() []string {
	return []string{
		filepath.Join(os.TempDir(), pdfToTextTempFilePattern),
		filepath.Join(os.TempDir(), StatementTempFilePattern),
		filepath.Join(os.TempDir(), MT940ToUTFTempFilePattern),
		filepath.Join(os.TempDir(), MBankToUTFTempFilePattern),
	}
}

func cleanupByPattern(pattern string) {
	matches, err := filepath.Glob(pattern)

	if err != nil {
		log.Printf("WARNING: could not scan temporary files for pattern '%s': %v\n", pattern, err)
		return
	}

	for _, match := range matches {
		err = os.Remove(match)

		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			log.Printf("WARNING: could not delete stale temporary file '%s': %v\n", match, err)
		}
	}
}

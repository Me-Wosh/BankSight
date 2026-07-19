package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	tempFilePatternPDFToText  = "banksight-converted-pdf-to-text-*.txt"
	tempFilePatternStatement  = "banksight-statement-*.html"
	tempFilePatternMT940ToUTF = "banksight-converted-mt940-to-utf8-*.txt"
	tempFilePatternMBankToUTF = "banksight-converted-mbank-to-utf8-*.csv"
)

func createTrackedTempFile(pattern string) (*os.File, error) {
	return os.CreateTemp("", pattern)
}

func removeTrackedTempFile(filePath string) {
	err := os.Remove(filePath)

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Printf("WARNING: could not delete temporary file '%s': %v\n", filePath, err)
	}
}

func cleanupTempFiles() {
	for _, pattern := range tempFilePatterns() {
		cleanupByPattern(pattern)
	}
}

func tempFilePatterns() []string {
	return []string{
		filepath.Join(os.TempDir(), tempFilePatternPDFToText),
		filepath.Join(os.TempDir(), tempFilePatternStatement),
		filepath.Join(os.TempDir(), tempFilePatternMT940ToUTF),
		filepath.Join(os.TempDir(), tempFilePatternMBankToUTF),
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

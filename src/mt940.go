package main

import (
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func convertFileToUtf8(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error while opening the file: %v\n", err)
		return "", err
	}

	defer file.Close()

	reader := transform.NewReader(file, charmap.CodePage852.NewDecoder())
	outputFile, err := createTrackedTempFile(tempFilePatternMT940ToUTF)

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

func extractDatePart(line string) string {
	// prefix, date of transaction 6dF
	return line[4 : 4+6]
}

func extractMoneyPart(line string) string {
	var moneyPart strings.Builder

	// prefix, date of transaction 6dF, date of post 4dF
	for i := 4 + 6 + 4; i < len(line); i++ {
		if line[i] == 'S' {
			break
		}

		moneyPart.WriteByte(line[i])
	}

	return moneyPart.String()
}

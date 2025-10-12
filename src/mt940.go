package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func convertFileToUtf8(filePath string) string {
	file, err := os.Open(filePath)

	if err != nil {
		panic(fmt.Sprintf("Error while opening the file: %v", err))
	}

	defer file.Close()

	reader := transform.NewReader(file, charmap.CodePage852.NewDecoder())

	outputFilePath := "utf-8_" + strings.Replace(filepath.Base(filePath), ".sta", ".txt", 1)

	outputFile, err := os.Create(outputFilePath)

	if err != nil {
		panic(fmt.Sprintf("Error while creating the file: %v", err))
	}

	defer outputFile.Close()

	if _, err := io.Copy(outputFile, reader); err != nil {
		panic(fmt.Sprintf("Error while copying the file content: %v", err))
	}

	return outputFilePath
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

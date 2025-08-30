package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

func convertPdfToText(filePath string) string {
	textFilePath := strings.Replace(filepath.Base(filePath), ".pdf", ".txt", 1)
	output, err := exec.Command("pdftotext", "-layout", filePath, textFilePath).CombinedOutput()

	if err != nil {
		panic(fmt.Sprintf("Error while converting PDF to text: %v. %s", err, output))
	}

	return textFilePath
}

func readAllLines(filePath string) []string {
	file, err := os.Open(filePath)

	if err != nil {
		panic(fmt.Sprintf("Error while opening the file: %v", err))
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("Error while reading the file: %v", err))
	}

	if len(lines) == 0 {
		panic("Error while reading the file: No lines were read from the file.")
	}

	return lines
}

func getShopCategories() map[string]string {
	yamlFile, err := os.ReadFile(shopCategoriesFilePath)

	if err != nil {
		panic(fmt.Sprintf("Error while reading YAML file: %v", err))
	}

	var categories map[string]string

	if err := yaml.Unmarshal(yamlFile, &categories); err != nil {
		panic(fmt.Sprintf("Error while unmarshalling YAML file: %v", err))
	}

	return categories
}

func getFilePath(filePathFlag string) string {
	var filePath string

	if filePathFlag != "" {
		filePath = filePathFlag
	} else {
		panic("File path was not provided.")
	}

	if !strings.HasSuffix(filePath, ".pdf") {
		panic("File must be a PDF file.")
	}

	return filePath
}

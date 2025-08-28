package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/goccy/go-yaml"
)

func convertPdfToText(fileName string) string {
	output, err := exec.Command("pdftotext", "-layout", fileName).CombinedOutput()

	if err != nil {
		panic(fmt.Sprintf("Error while converting PDF to text: %v. %s", err, output))
	}

	return strings.Replace(fileName, ".pdf", ".txt", 1)
}

func readAllLines(fileName string) []string {
	file, err := os.Open(fileName)

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

func getFileName(fileNameFlag string) string {
	var fileName string

	if fileNameFlag != "" {
		fileName = fileNameFlag
	} else {
		panic("File name was not provided.")
	}

	if !strings.HasSuffix(fileName, ".pdf") {
		panic("File must be a PDF file.")
	}

	return fileName
}

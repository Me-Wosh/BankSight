package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func readAllLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error while opening the file: %v\n", err)
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)

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

	file.Close()
	removeTrackedTempFile(filePath)

	return lines, nil
}

func getShopCategories() (map[string]string, error) {
	yamlFile, err := os.ReadFile(shopCategoriesFilePath)

	if err != nil {
		log.Printf("Error while reading YAML file: %v\n", err)
		return nil, err
	}

	var categories map[string]string

	if err := yaml.Unmarshal(yamlFile, &categories); err != nil {
		log.Printf("Error while unmarshalling YAML file: %v\n", err)
		return nil, err
	}

	return categories, nil
}

func getBank(bank string) (bank, error) {
	if strings.EqualFold(bank, "pko") {
		return &pko{}, nil
	}

	if strings.EqualFold(bank, "ing") {
		return &ing{}, nil
	}

	if strings.EqualFold(bank, "mbank") {
		return &mbank{}, nil
	}

	log.Printf("ERROR: Unsupported bank: %s\n", bank)
	return nil, fmt.Errorf("ERROR: Unsupported bank: %s", bank)
}

func convertFileToUtf8(filePath string, decoder *encoding.Decoder, tempFilePattern string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error while opening the file: %v\n", err)
		return "", err
	}

	defer file.Close()

	reader := transform.NewReader(file, decoder)
	outputFile, err := createTrackedTempFile(tempFilePattern)

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

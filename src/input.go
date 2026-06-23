package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
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

	log.Printf("ERROR: Unsupported bank: %s\n", bank)
	return nil, fmt.Errorf("ERROR: Unsupported bank: %s", bank)
}

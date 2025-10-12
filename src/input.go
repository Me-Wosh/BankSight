package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

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

func getBank(bankFlag string) bank {
	if strings.EqualFold(bankFlag, "pko") {
		return &pko{}
	}

	if strings.EqualFold(bankFlag, "ing") {
		return &ing{}
	}

	panic(fmt.Sprintf("Unsupported bank: %s", bankFlag))
}

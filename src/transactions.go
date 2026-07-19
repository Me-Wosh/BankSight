package main

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"math"
)

type transactions struct {
	timePeriod                      string
	spendings, incomes              float64
	previousBalance, closingBalance float64
	difference                      float64
	categorizedBalance              map[string]map[string]float64
}

type bank interface {
	prepareFile(filePath string) (string, error)
	calculateTotalTransactions(lines []string, shopCategories map[string]string) (transactions, error)
	getSupportedFileExtensions() []string
}

func validateCalculatedBalance(previousBalance, spendings, incomes, closingBalance float64) error {
	calculatedBalance := previousBalance + spendings + incomes

	if math.Abs(calculatedBalance-closingBalance) > 0.01 {
		errorMessage := fmt.Sprintf(
			"ERROR: Calculated balance (%.2f) does not match closing balance (%.2f)",
			calculatedBalance,
			closingBalance,
		)
		log.Println(errorMessage)
		return errors.New(errorMessage)
	}

	log.Println("INFO: Calculated balance matches closing balance")
	return nil
}

func validateCategorizedSpendings(
	categorizedBalance map[string]map[string]float64,
	spendings float64,
	totalKey string,
) error {
	var categorizedSpendings float64
	categories := maps.Values(categorizedBalance)

	for category := range categories {
		categorizedSpendings += category[totalKey]
	}

	log.Printf("INFO: Sum of categorized spendings: %.2f\n", categorizedSpendings)

	if math.Abs(categorizedSpendings-spendings) > 0.01 {
		errorMessage := fmt.Sprintf(
			"ERROR: Sum of categorized spendings (%.2f) don't add up to calculated spendings (%.2f)",
			categorizedSpendings,
			spendings,
		)
		log.Println(errorMessage)
		return errors.New(errorMessage)
	}

	return nil
}

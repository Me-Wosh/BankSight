package bank

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"math"
	"strings"
)

type Bank interface {
	PrepareFile(filePath string) (string, error)
	CalculateTotalTransactions(lines []string, shopCategories map[string]string) (Transactions, error)
	GetSupportedFileExtensions() []string
}

type Transactions struct {
	TimePeriod                      string
	Spendings, Incomes              float64
	PreviousBalance, ClosingBalance float64
	Difference                      float64
	CategorizedBalance              map[string]map[string]float64
}

func GetBank(selectedBank string) (Bank, error) {
	if strings.EqualFold(selectedBank, "pko") {
		return &PKO{}, nil
	}

	if strings.EqualFold(selectedBank, "ing") {
		return &ING{}, nil
	}

	if strings.EqualFold(selectedBank, "mbank") {
		return &MBank{}, nil
	}

	log.Printf("ERROR: Unsupported bank: %s\n", selectedBank)
	return nil, fmt.Errorf("ERROR: Unsupported bank: %s", selectedBank)
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

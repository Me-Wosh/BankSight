package main

import (
	"fmt"
	"log"
	"maps"
	"math"
	"strings"
)

type pko struct{}

func (pko *pko) prepareFile(filePath string) (string, error) {
	if !strings.HasSuffix(filePath, ".pdf") {
		log.Println("ERROR: PKO only supports PDF files.")
		return "", fmt.Errorf("ERROR: PKO only supports PDF files.")
	}

	return convertPdfToText(filePath)
}

func (pko *pko) calculateTotalTransactions(lines []string) (transactions, error) {
	shopCategories, err := getShopCategories()

	if err != nil {
		return transactions{}, err
	}

	categorizedBalance := make(map[string]map[string]float64)
	operations := map[string]struct{}{
		"zakup":    {},
		"przelew":  {},
		"płatność": {},
		"opłata":   {},
	}

	var (
		spendings       float64
		incomes         float64
		previousBalance float64
		closingBalance  float64
		lastOperation   string
		lastAmount      float64
	)

	const (
		otherCategory        = "Other"
		totalKey             = "total"
		previousBalanceLabel = "Saldo poprzednie"
		closingBalanceLabel  = "Saldo końcowe"
	)

	for _, line := range lines {
		sections := divideLineIntoSections(line)
		sectionsLength := len(sections)

		if sectionsLength == 2 {
			if sections[0] == previousBalanceLabel {
				validFloat, err := convertToValidFloat(sections[1])

				if err != nil {
					log.Println("WARNING: Failed to convert value. Expected a number, but got:", sections[1])
					continue
				}

				previousBalance = validFloat
				log.Printf("INFO: Scanned initial balance: %.2f\n", previousBalance)
			} else if sections[0] == closingBalanceLabel {
				validFloat, err := convertToValidFloat(sections[1])

				if err != nil {
					log.Println("WARNING: Failed to convert value. Expected a number, but got:", sections[1])
					continue
				}

				closingBalance = validFloat
				log.Printf("INFO: Scanned closing balance: %.2f\n", closingBalance)
			}
		}

		if sectionsLength == 2 || sectionsLength == 3 {
			_, contains := operations[lastOperation]

			if !contains || lastAmount >= 0 {
				continue
			}

			var combinedSections string

			if sectionsLength == 3 {
				combinedSections = strings.ToLower(sections[1] + " " + sections[2])

			} else {
				combinedSections = strings.ToLower(sections[1])
			}

			categoryFound := false
			keys := maps.Keys(shopCategories)

			for key := range keys {
				if strings.Contains(combinedSections, key) {
					categoryFound = true
					category := shopCategories[key]

					if categorizedBalance[category] == nil {
						categorizedBalance[category] = make(map[string]float64)
					}

					categorizedBalance[category][key] += lastAmount
					categorizedBalance[category][totalKey] += lastAmount

					log.Printf(
						"INFO: Found category '%s' for section '%s', amount: %.2f\n",
						category,
						combinedSections,
						lastAmount,
					)

					break
				}
			}

			if !categoryFound {
				if categorizedBalance[otherCategory] == nil {
					categorizedBalance[otherCategory] = make(map[string]float64)
				}

				categorizedBalance[otherCategory][combinedSections] += lastAmount
				categorizedBalance[otherCategory][totalKey] += lastAmount

				log.Printf(
					"WARNING: No category found for '%s', adding to '%s', amount: %.2f\n",
					combinedSections,
					otherCategory,
					lastAmount,
				)
			}

			lastOperation = ""
			lastAmount = 0
		}

		if sectionsLength == 5 {
			validFloat, err := convertToValidFloat(sections[3])

			if err != nil {
				log.Println("WARNING: Failed to convert value. Expected a number, but got:", sections[3])
				continue
			}

			words := strings.Fields(sections[2])
			lastOperation = strings.ToLower(words[0])
			lastAmount = validFloat

			log.Println("INFO: Scanned number:", validFloat)

			if validFloat < 0 {
				spendings += validFloat
			} else {
				incomes += validFloat
			}
		}
	}

	calculatedBalance := previousBalance + spendings + incomes

	if math.Abs(calculatedBalance-closingBalance) > 0.01 {
		errorMessage := fmt.Sprintf(
			"ERROR: Calculated balance (%.2f) does not match closing balance (%.2f)",
			calculatedBalance,
			closingBalance,
		)
		log.Println(errorMessage)
		return transactions{}, fmt.Errorf(errorMessage)
	}

	log.Println("INFO: Calculated balance matches closing balance")

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
		return transactions{}, fmt.Errorf(errorMessage)
	}

	return transactions{
		timePeriod:         pko.getTimePeriod(lines),
		spendings:          spendings,
		incomes:            incomes,
		previousBalance:    previousBalance,
		closingBalance:     closingBalance,
		categorizedBalance: categorizedBalance,
		difference:         closingBalance - previousBalance,
	}, nil
}

func (pko *pko) getTimePeriod(lines []string) string {
	periodLine := strings.ToLower(lines[4])
	_, after, _ := strings.Cut(periodLine, "okres ")
	return after
}

func (pko *pko) getSupportedFileExtensions() []string {
	return []string{".pdf"}
}

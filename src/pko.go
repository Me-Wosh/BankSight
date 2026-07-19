package main

import (
	"errors"
	"log"
	"maps"
	"strings"
)

type pko struct{}

func (pko *pko) prepareFile(filePath string) (string, error) {
	if !strings.HasSuffix(filePath, ".pdf") {
		log.Println("ERROR: PKO only supports PDF files.")
		return "", errors.New("ERROR: PKO only supports PDF files.")
	}

	return convertPdfToText(filePath)
}

func (pko *pko) calculateTotalTransactions(lines []string, shopCategories map[string]string) (transactions, error) {
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
					log.Printf(
						"WARNING: Failed to convert previous balance. Expected a number, but got: %s. Error: %v\n",
						sections[1],
						err,
					)
					continue
				}

				previousBalance = validFloat
				log.Printf("INFO: Scanned initial balance: %.2f\n", previousBalance)
			} else if sections[0] == closingBalanceLabel {
				validFloat, err := convertToValidFloat(sections[1])

				if err != nil {
					log.Printf(
						"WARNING: Failed to convert closing balance. Expected a number, but got: %s. Error: %v\n",
						sections[1],
						err,
					)
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
				log.Printf(
					"WARNING: Failed to convert value. Expected a number, but got: %s. Error: %v\n",
					sections[3],
					err,
				)
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

	err := validateCalculatedBalance(previousBalance, spendings, incomes, closingBalance)

	if err != nil {
		return transactions{}, err
	}

	err = validateCategorizedSpendings(categorizedBalance, spendings, totalKey)

	if err != nil {
		return transactions{}, err
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

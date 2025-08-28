package main

import (
	"fmt"
	"maps"
	"math"
	"strings"
)

func calculateTotalTransactions(lines []string, debugFlag bool) (
	spendings, incomes float64,
	categoriesBalance map[string]float64,
) {
	categories := getShopCategories()
	categoriesBalance = make(map[string]float64)
	operations := map[string]struct{}{
		"zakup":    {},
		"przelew":  {},
		"płatność": {},
		"opłata":   {},
	}
	var lastOperation string
	var lastAmount float64
	var previousBalance float64
	var closingBalance float64
	const previousBalanceLabel = "Saldo poprzednie"
	const closingBalanceLabel = "Saldo końcowe"

	for _, line := range lines {
		sections := divideLineIntoSections(line)
		sectionsLength := len(sections)

		if sectionsLength == 2 {
			if sections[0] == previousBalanceLabel {
				validFloat, err := convertToValidFloat(sections[1])

				if err != nil {
					if debugFlag {
						fmt.Println("WARNING: Failed converting value. Expected a number, but got:", sections[1])
					}

					continue
				}

				previousBalance = validFloat

				if debugFlag {
					fmt.Printf("INFO: Scanned initial balance: %.2f\n", previousBalance)
				}
			} else if sections[0] == closingBalanceLabel {
				validFloat, err := convertToValidFloat(sections[1])

				if err != nil {
					if debugFlag {
						fmt.Println("WARNING: Failed converting value. Expected a number, but got:", sections[1])
					}

					continue
				}

				closingBalance = validFloat

				if debugFlag {
					fmt.Printf("INFO: Scanned closing balance: %.2f\n", closingBalance)
				}
			}
		}

		if sectionsLength == 2 || sectionsLength == 3 {
			_, contains := operations[lastOperation]

			if !contains || lastAmount >= 0 {
				continue
			}

			var combinedSections string

			if sectionsLength == 3 {
				combinedSections = strings.ToLower(sections[1] + sections[2])

			} else {
				combinedSections = strings.ToLower(sections[1])
			}

			categoryFound := false
			keys := maps.Keys(categories)

			for key := range keys {
				if strings.Contains(combinedSections, key) {
					categoryFound = true
					category := categories[key]
					categoriesBalance[category] += lastAmount

					if debugFlag {
						fmt.Printf(
							"INFO: Found category '%s' for section '%s', amount: %.2f\n",
							category,
							combinedSections,
							lastAmount,
						)
					}

					break
				}
			}

			if !categoryFound {
				categoriesBalance["Other"] += lastAmount
				if debugFlag {
					fmt.Printf(
						"WARNING: No category found for '%s', adding to 'Other', amount: %.2f\n",
						combinedSections,
						lastAmount,
					)
				}
			}

			lastOperation = ""
			lastAmount = 0
		}

		if sectionsLength == 5 {
			validFloat, err := convertToValidFloat(sections[3])

			if err != nil {
				if debugFlag {
					fmt.Println("WARNING: Failed converting value. Expected a number, but got:", sections[3])
				}

				continue
			}

			words := strings.Fields(sections[2])
			lastOperation = strings.ToLower(words[0])
			lastAmount = validFloat

			if debugFlag {
				fmt.Println("INFO: Scanned number:", validFloat)
			}

			if validFloat < 0 {
				spendings += validFloat
			} else {
				incomes += validFloat
			}
		}
	}

	calculatedBalance := previousBalance + spendings + incomes

	if math.Abs(calculatedBalance-closingBalance) > 0.01 {
		panic(fmt.Sprintf(
			"ERROR: Calculated balance (%.2f) does not match closing balance (%.2f)",
			calculatedBalance,
			closingBalance,
		))
	}

	if debugFlag {
		fmt.Println("INFO: Calculated balance matches closing balance")
	}

	var categorizedSpendings float64
	values := maps.Values(categoriesBalance)
	for value := range values {
		categorizedSpendings += value
	}

	if debugFlag {
		fmt.Printf("INFO: Categorized spendings: %.2f\n", categorizedSpendings)
	}

	if math.Abs(categorizedSpendings-spendings) > 0.01 {
		panic(fmt.Sprintf(
			"ERROR: Categorized spendings (%.2f) don't add up to calculated spendings (%.2f)",
			categorizedSpendings,
			spendings,
		))
	}

	return spendings, incomes, categoriesBalance
}

package main

import (
	"fmt"
	"maps"
	"math"
	"strings"
)

type transactions struct {
	spendings, incomes              float64
	previousBalance, closingBalance float64
	categoriesBalance               map[string]map[string]float64
}

func calculateTotalTransactions(lines []string, debugFlag bool) transactions {
	shopCategories := getShopCategories()
	categoriesBalance := make(map[string]map[string]float64)
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
			keys := maps.Keys(shopCategories)

			for key := range keys {
				if strings.Contains(combinedSections, key) {
					categoryFound = true
					category := shopCategories[key]

					if categoriesBalance[category] == nil {
						categoriesBalance[category] = make(map[string]float64)
					}

					categoriesBalance[category][key] += lastAmount
					categoriesBalance[category][totalKey] += lastAmount

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
				if categoriesBalance[otherCategory] == nil {
					categoriesBalance[otherCategory] = make(map[string]float64)
				}

				categoriesBalance[otherCategory][combinedSections] += lastAmount
				categoriesBalance[otherCategory][totalKey] += lastAmount

				if debugFlag {
					fmt.Printf(
						"WARNING: No category found for '%s', adding to '%s', amount: %.2f\n",
						combinedSections,
						otherCategory,
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
	categories := maps.Values(categoriesBalance)

	for category := range categories {
		categorizedSpendings += category[totalKey]
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

	return transactions{
		spendings:         spendings,
		incomes:           incomes,
		previousBalance:   previousBalance,
		closingBalance:    closingBalance,
		categoriesBalance: categoriesBalance,
	}
}

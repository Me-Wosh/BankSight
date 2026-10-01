package bank

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Me-Wosh/BankSight/internal/file"
	"github.com/Me-Wosh/BankSight/internal/input"

	"golang.org/x/text/encoding/charmap"
)

type ING struct{}

func (ing *ING) PrepareFile(filePath string) (string, error) {
	if !strings.HasSuffix(filePath, ".sta") {
		log.Println("ERROR: ING only supports MT940 (.sta) files.")
		return "", errors.New("ERROR: ING only supports MT940 (.sta) files.")
	}

	return file.ConvertToUtf8(filePath, charmap.CodePage852.NewDecoder(), file.MT940ToUTFTempFilePattern)
}

func (ing *ING) CalculateTotalTransactions(lines []string, shopCategories map[string]string) (Transactions, error) {
	categorizedBalance := make(map[string]map[string]float64)

	var (
		startDate       string
		endDate         string
		spendings       float64
		incomes         float64
		previousBalance float64
		closingBalance  float64
		lastAmount      float64
	)

	const (
		otherCategory = "Other"
		totalKey      = "total"
	)

	for i, line := range lines {
		if strings.HasPrefix(line, ":60F:") {
			// prefix, balance sign 1f, date 6dF, currency 3a
			moneyPart := line[4+1+6+3+1:]
			value, err := input.ConvertToValidFloat(moneyPart)

			if err != nil {
				log.Printf(
					"WARNING: Failed to convert previous balance. Expected a number, but got: %s. Error: %v\n",
					moneyPart,
					err,
				)
				continue
			}

			previousBalance = value
			log.Printf("INFO: Scanned previous balance: %.2f\n", previousBalance)
		} else if strings.HasPrefix(line, ":61:") {
			if startDate == "" {
				startDate = extractDatePart(line)
			} else {
				endDate = extractDatePart(line)
			}

			moneyPart := extractMoneyPart(line)

			value, err := input.ConvertToValidFloat(moneyPart[1:])

			if err != nil {
				log.Printf(
					"WARNING: Failed to convert transaction amount. Expected a number, but got: %s. Error: %v\n",
					moneyPart,
					err,
				)
				continue
			}

			switch moneyPart[0] {
			case 'D':
				lastAmount = -value
				spendings -= value
			case 'C':
				lastAmount = value
				incomes += value
			default:
				log.Printf(
					"WARNING: Could not determine transaction type, expected 'D' or 'C' but got: %c instead\n",
					moneyPart[0],
				)
			}
		} else if strings.HasPrefix(line, "~21") {
			_, after, _ := strings.Cut(line, "~22")
			shopPart := after
			joinedShopPartLines := strings.Replace(shopPart, "~33", "", 1)

			key := strings.ToLower(joinedShopPartLines)
			category, categoryFound := shopCategories[key]

			if categoryFound {
				if categorizedBalance[category] == nil {
					categorizedBalance[category] = make(map[string]float64)
				}

				categorizedBalance[category][key] += lastAmount
				categorizedBalance[category][totalKey] += lastAmount

				log.Printf(
					"INFO: Found category '%s' for section '%s', amount: %.2f\n",
					category,
					joinedShopPartLines,
					lastAmount,
				)

				continue
			}

			key = strings.TrimSpace(strings.ToLower(lines[i-1][3:])) // previous line without the ~21 prefix
			category, categoryFound = shopCategories[key]

			if categoryFound {
				if categorizedBalance[category] == nil {
					categorizedBalance[category] = make(map[string]float64)
				}

				categorizedBalance[category][key] += lastAmount
				categorizedBalance[category][totalKey] += lastAmount

				log.Printf(
					"INFO: Found category '%s' for section '%s', amount: %.2f\n",
					category,
					key,
					lastAmount,
				)

				continue
			}

			if !categoryFound && lastAmount < 0 {
				if categorizedBalance[otherCategory] == nil {
					categorizedBalance[otherCategory] = make(map[string]float64)
				}

				categorizedBalance[otherCategory][joinedShopPartLines] += lastAmount
				categorizedBalance[otherCategory][totalKey] += lastAmount

				log.Printf(
					"WARNING: No category found for '%s', nor '%s' adding to '%s', amount: %.2f\n",
					joinedShopPartLines,
					key,
					otherCategory,
					lastAmount,
				)
			}

			lastAmount = 0
		} else if strings.HasPrefix(line, ":62F:") {
			// prefix, balance sign 1f, date 6dF, currency 3a
			moneyPart := line[4+1+6+3+1:]
			value, err := input.ConvertToValidFloat(moneyPart)

			if err != nil {
				log.Printf(
					"WARNING: Failed to convert closing balance. Expected a number, but got: %s. Error: %v\n",
					moneyPart,
					err,
				)
				continue
			}

			closingBalance = value
			log.Printf("INFO: Scanned closing balance: %.2f\n", closingBalance)
		}
	}

	err := validateCalculatedBalance(previousBalance, spendings, incomes, closingBalance)

	if err != nil {
		return Transactions{}, err
	}

	err = validateCategorizedSpendings(categorizedBalance, spendings, totalKey)

	if err != nil {
		return Transactions{}, err
	}

	timePeriod, err := ing.getTimePeriod(startDate, endDate)

	if err != nil {
		return Transactions{}, err
	}

	return Transactions{
		TimePeriod:         timePeriod,
		Spendings:          spendings,
		Incomes:            incomes,
		PreviousBalance:    previousBalance,
		ClosingBalance:     closingBalance,
		CategorizedBalance: categorizedBalance,
		Difference:         closingBalance - previousBalance,
	}, nil
}

func (ing *ING) GetSupportedFileExtensions() []string {
	return []string{".sta"}
}

func (ing *ING) getTimePeriod(startDate, endDate string) (string, error) {
	startYear := startDate[0:2]
	startMonth := startDate[2:4]

	endYear := endDate[0:2]
	endMonth := endDate[2:4]
	var endDay string

	if endMonth == "02" {
		yearNumber, err := strconv.Atoi(endYear)

		if err != nil {
			log.Printf("Error while converting year: %s to number\n", endYear)
			return "", fmt.Errorf("Error while converting year: %s to number\n", endYear)
		}

		if yearNumber%4 == 0 {
			endDay = "29"
		} else {
			endDay = "28"
		}
	} else {
		monthNumber, err := strconv.Atoi(endMonth)

		if err != nil {
			log.Printf("Error while converting month: %s to number\n", endMonth)
			return "", fmt.Errorf("Error while converting month: %s to number\n", endMonth)
		}

		if monthNumber%2 == 0 {
			endDay = "30"
		} else {
			endDay = "31"
		}
	}

	return fmt.Sprintf("20%s.%s.%s - 20%s.%s.%s", startYear, startMonth, "01", endYear, endMonth, endDay), nil
}

func extractDatePart(line string) string {
	// prefix, date of transaction 6dF
	return line[4 : 4+6]
}

func extractMoneyPart(line string) string {
	var moneyPart strings.Builder

	// prefix, date of transaction 6dF, date of post 4dF
	for i := 4 + 6 + 4; i < len(line); i++ {
		if line[i] == 'S' {
			break
		}

		moneyPart.WriteByte(line[i])
	}

	return moneyPart.String()
}

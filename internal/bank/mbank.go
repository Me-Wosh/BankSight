package bank

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"math"
	"strings"

	"github.com/Me-Wosh/BankSight/internal/file"
	"github.com/Me-Wosh/BankSight/internal/input"
	"golang.org/x/text/encoding/charmap"
)

type MBank struct{}

func (mBank *MBank) PrepareFile(filePath string) (string, error) {
	if !strings.HasSuffix(filePath, ".csv") {
		log.Println("ERROR: mBank only supports CSV files.")
		return "", errors.New("ERROR: mBank only supports CSV files.")
	}

	return file.ConvertToUtf8(filePath, charmap.Windows1250.NewDecoder(), file.MBankToUTFTempFilePattern)
}

func (mBank *MBank) CalculateTotalTransactions(lines []string, shopCategories map[string]string) (Transactions, error) {
	categorizedBalance := make(map[string]map[string]float64)

	var (
		timePeriod        string
		spendings         float64
		incomes           float64
		previousBalance   float64
		closingBalance    float64
		expectedSpendings float64
		expectedIncomes   float64
	)

	const (
		otherCategory        = "Other"
		totalKey             = "total"
		periodLabel          = "#Za okres:"
		incomesLabel         = "Uznania"
		spendingsLabel       = "Obciążenia"
		previousBalanceLabel = "#Saldo początkowe"
		closingBalanceLabel  = "#Saldo końcowe"
		transactionDateLabel = "Data księgowania"
		savingsCategory      = "Savings"
		savingsDescription   = "PRZELEW NA TWOJE CELE"
	)

	for i, line := range lines {
		sections := divideMBankLineIntoSections(line)
		sectionsLength := len(sections)

		if sectionsLength == 0 {
			continue
		}

		switch sections[0] {
		case periodLabel:
			if i+1 < len(lines) {
				timePeriod = mBank.parseTimePeriod(lines[i+1])
			}

		case incomesLabel:
			if sectionsLength != 3 {
				continue
			}

			validFloat, err := input.ConvertToValidFloat(strings.TrimSuffix(sections[2], " PLN"))

			if err != nil {
				log.Printf(
					"WARNING: Failed to convert expected incomes. Expected a number, but got: %s. Error: %v\n",
					sections[2],
					err,
				)
				continue
			}

			expectedIncomes = validFloat

		case spendingsLabel:
			if sectionsLength != 3 {
				continue
			}

			validFloat, err := input.ConvertToValidFloat(strings.TrimSuffix(strings.TrimSpace(sections[2]), " PLN"))

			if err != nil {
				log.Printf(
					"WARNING: Failed to convert expected spendings. Expected a number, but got: %s. Error: %v\n",
					sections[2],
					err,
				)
				continue
			}

			expectedSpendings = validFloat

		case previousBalanceLabel:
			if sectionsLength != 2 {
				continue
			}

			validFloat, err := input.ConvertToValidFloat(strings.TrimSuffix(sections[1], " PLN"))

			if err != nil {
				log.Printf(
					"WARNING: Failed to convert previous balance. Expected a number, but got: %s. Error: %v\n",
					sections[1],
					err,
				)
				continue
			}

			previousBalance = validFloat
			log.Printf("INFO: Scanned previous balance: %.2f\n", previousBalance)
		}

		// Closing balance is at position [6] in ";;;;;;#Saldo końcowe;345,78 PLN;"
		if sectionsLength == 8 && sections[6] == closingBalanceLabel {
			closingBalanceValue := strings.TrimSuffix(sections[7], " PLN")
			validFloat, err := input.ConvertToValidFloat(closingBalanceValue)

			if err != nil {
				log.Printf(
					"WARNING: Failed to convert closing balance. Expected a number, but got: %s. Error: %v\n",
					closingBalanceValue,
					err,
				)
				continue
			}

			closingBalance = validFloat
			log.Printf("INFO: Scanned closing balance: %.2f\n", closingBalance)

			continue
		}

		// Skip other rows that are not transaction rows and transaction header line
		if sectionsLength != 8 || sections[0] == transactionDateLabel {
			continue
		}

		validFloat, err := input.ConvertToValidFloat(strings.TrimSpace(sections[6]))

		if err != nil {
			log.Printf(
				"WARNING: Failed to convert transaction amount. Expected a number, but got: %s. Error: %v\n",
				sections[6],
				err,
			)
			continue
		}

		if validFloat < 0 {
			spendings += validFloat
			transactionTitle := strings.ToLower(sections[3])
			categoryFound := false
			transactionDescription := sections[2]

			if strings.EqualFold(transactionDescription, savingsDescription) {
				categoryFound = true

				if categorizedBalance[savingsCategory] == nil {
					categorizedBalance[savingsCategory] = make(map[string]float64)
				}

				categorizedBalance[savingsCategory][transactionTitle] += validFloat
				categorizedBalance[savingsCategory][totalKey] += validFloat

				log.Printf(
					"INFO: Found category '%s' for transaction '%s', amount: %.2f\n",
					savingsCategory,
					transactionTitle,
					validFloat,
				)

				continue
			}

			keys := maps.Keys(shopCategories)

			for key := range keys {
				if strings.Contains(transactionTitle, key) {
					categoryFound = true
					category := shopCategories[key]

					if categorizedBalance[category] == nil {
						categorizedBalance[category] = make(map[string]float64)
					}

					categorizedBalance[category][key] += validFloat
					categorizedBalance[category][totalKey] += validFloat

					log.Printf(
						"INFO: Found category '%s' for transaction '%s', amount: %.2f\n",
						category,
						transactionTitle,
						validFloat,
					)

					break
				}
			}

			if !categoryFound {
				if categorizedBalance[otherCategory] == nil {
					categorizedBalance[otherCategory] = make(map[string]float64)
				}

				categorizedBalance[otherCategory][transactionTitle] += validFloat
				categorizedBalance[otherCategory][totalKey] += validFloat

				log.Printf(
					"WARNING: No category found for '%s', adding to '%s', amount: %.2f\n",
					transactionTitle,
					otherCategory,
					validFloat,
				)
			}
		} else {
			incomes += validFloat
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

	if math.Abs(math.Abs(spendings)-expectedSpendings) > 0.01 {
		errorMessage := fmt.Sprintf(
			"ERROR: Calculated spendings (%.2f) do not match expected spendings from summary (%.2f)",
			math.Abs(spendings),
			expectedSpendings,
		)
		log.Println(errorMessage)

		return Transactions{}, errors.New(errorMessage)
	}

	if math.Abs(incomes-expectedIncomes) > 0.01 {
		errorMessage := fmt.Sprintf(
			"ERROR: Calculated incomes (%.2f) do not match expected incomes from summary (%.2f)",
			incomes,
			expectedIncomes,
		)
		log.Println(errorMessage)

		return Transactions{}, errors.New(errorMessage)
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

func (mBank *MBank) GetSupportedFileExtensions() []string {
	return []string{".csv"}
}

func (mBank *MBank) parseTimePeriod(line string) string {
	fields := strings.SplitN(line, ";", 3)

	if len(fields) < 2 {
		return ""
	}

	return fmt.Sprintf("%s - %s", fields[0], fields[1])
}

// Splits a semicolon-delimited line respecting double-quoted fields.
// An unescaped '"' inside a quoted field is treated as a literal character when not followed by ';'
func divideMBankLineIntoSections(line string) []string {
	var sections []string
	var section strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		switch {
		case ch == '"' && !inQuotes:
			inQuotes = true

		case ch == '"' && inQuotes:
			if i >= len(line)-1 || line[i+1] == ';' {
				inQuotes = false
			} else {
				section.WriteByte(ch)
			}

		case ch == ';' && !inQuotes:
			sections = append(sections, section.String())
			section.Reset()

		default:
			section.WriteByte(ch)
		}
	}

	if section.Len() > 0 {
		sections = append(sections, section.String())
	}

	return sections
}

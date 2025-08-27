package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/goccy/go-yaml"
)

var fileNameFlag string
var debugFlag bool

func main() {
	setupCommandLineFlags()

	fileName := getFileName()
	convertPdfToText(fileName)
	lines := readAllLines(strings.Replace(fileName, ".pdf", ".txt", 1))
	timePeriod := getTimePeriod(lines)
	spendings, incomes, categoriesBalance := calculateTotalTransactions(lines)
	drawPieChart(spendings, incomes, categoriesBalance, timePeriod, fileName)

	fmt.Printf("Spendings: %.2f\n", spendings)
	fmt.Printf("Incomes: %.2f", incomes)
}

func calculateTotalTransactions(lines []string) (spendings, incomes float64, categoriesBalance map[string]float64) {
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

func drawPieChart(spendings, incomes float64, categoriesBalance map[string]float64, subtitle, fileName string) {
	spendingsIncomesData := []opts.PieData{
		{Name: "Spendings", Value: fmt.Sprintf("%.2f", spendings*-1)},
		{Name: "Incomes", Value: fmt.Sprintf("%.2f", incomes)},
	}

	var categorizedSpendingsData []opts.PieData
	categories := maps.Keys(categoriesBalance)
	for category := range categories {
		categorizedSpendingsData = append(categorizedSpendingsData, opts.PieData{
			Name:  category,
			Value: fmt.Sprintf("%.2f", categoriesBalance[category]*-1),
		})
	}

	pieChart := charts.NewPie()

	pieChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Total transactions", Subtitle: "Time period: " + subtitle}),
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{Bottom: "0"}),
	)

	pieChart.AddSeries(
		"",
		spendingsIncomesData,
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"0", "35%"}}),
		charts.WithLabelOpts(opts.Label{Position: "inside", Formatter: "{b}:\n{c} zł"}),
	)

	pieChart.AddSeries(
		"",
		categorizedSpendingsData,
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"50%", "70%"}}),
		charts.WithLabelOpts(opts.Label{Formatter: "{b}: {c} zł"}),
	)

	page := components.NewPage()
	page.SetPageTitle("BankSight")
	page.SetLayout(components.PageFullLayout)
	page.AddCustomizedCSSAssets("styles.css")
	page.AddCharts(pieChart)

	file, err := os.Create(strings.Replace(fileName, ".pdf", ".html", 1))

	if err != nil {
		panic(fmt.Sprintf("Error while creating HTML file: %v", err))
	}

	if err := page.Render(io.MultiWriter(file)); err != nil {
		panic(fmt.Sprintf("Error while rendering HTML file: %v", err))
	}
}

func convertToValidFloat(str string) (float64, error) {
	removedSpaces := strings.ReplaceAll(str, " ", "")
	validFloatFormat := strings.ReplaceAll(removedSpaces, ",", ".")
	validFloat, err := strconv.ParseFloat(validFloatFormat, 64)

	if err != nil {
		return 0, err
	}

	return validFloat, nil
}

func divideLineIntoSections(line string) []string {
	line = strings.TrimSpace(line)
	sections := regexp.MustCompile(" {2,}").Split(line, -1) // divide line into sections on two or more spaces

	return sections
}

func convertPdfToText(fileName string) {
	output, err := exec.Command("pdftotext", "-layout", fileName).CombinedOutput()

	if err != nil {
		panic(fmt.Sprintf("Error while converting PDF to text: %v. %s", err, output))
	}
}

func readAllLines(fileName string) []string {
	file, err := os.Open(fileName)

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

	return lines
}

func getFileName() string {
	var fileName string

	if fileNameFlag != "" {
		fileName = fileNameFlag
	} else {
		panic("File name was not provided.")
	}

	if !strings.HasSuffix(fileName, ".pdf") {
		panic("File must be a PDF file.")
	}

	return fileName
}

func getTimePeriod(lines []string) string {
	periodLine := strings.ToLower(lines[4])
	_, after, _ := strings.Cut(periodLine, "okres ")
	return after
}

func getShopCategories() map[string]string {
	yamlFile, err := os.ReadFile("../shopcategories.yaml")

	if err != nil {
		panic(fmt.Sprintf("Error while reading YAML file: %v", err))
	}

	var categories map[string]string

	if err := yaml.Unmarshal(yamlFile, &categories); err != nil {
		panic(fmt.Sprintf("Error while unmarshalling YAML file: %v", err))
	}

	return categories
}

func setupCommandLineFlags() {
	flag.StringVar(&fileNameFlag, "file", "", "(Required) Path to the file containing bank statement lines")
	flag.StringVar(&fileNameFlag, "f", "", "Alias for -file")
	flag.BoolVar(&debugFlag, "debug", false, "(Optional) Enable debugging info")
	flag.BoolVar(&debugFlag, "d", false, "Alias for -debug")
	flag.Parse()
}

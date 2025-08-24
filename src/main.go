package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var fileNameFlag string
var debugFlag bool

func main() {
	setupCommandLineFlags()

	fileName := getFileName()
	convertPdfToText(fileName)
	lines := readAllLines(strings.Replace(fileName, ".pdf", ".txt", 1))
	spendings, incomes := calculateTotalTransactions(lines)
	drawPieChart(fileName, spendings, incomes)

	fmt.Printf("Spendings: %.2f\n", spendings)
	fmt.Printf("Incomes: %.2f", incomes)
}

func calculateTotalTransactions(lines []string) (spendings, incomes float64) {
	var previousBalance float64
	var closingBalance float64
	const previousBalanceLabel = "Saldo poprzednie"
	const closingBalanceLabel = "Saldo końcowe"

	for _, line := range lines {
		sections := divideLineIntoSections(line)

		if len(sections) == 2 {
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

		if len(sections) == 5 {
			validFloat, err := convertToValidFloat(sections[3])

			if err != nil {
				if debugFlag {
					fmt.Println("WARNING: Failed converting value. Expected a number, but got:", sections[3])
				}

				continue
			}

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

	if calculatedBalance != closingBalance {
		panic(fmt.Sprintf(
			"WARNING: Calculated balance (%.2f) does not match closing balance (%.2f)",
			calculatedBalance,
			closingBalance,
		))
	}

	if debugFlag {
		fmt.Println("INFO: Calculated balance matches closing balance")
	}

	return spendings, incomes
}

func drawPieChart(fileName string, spendings float64, incomes float64) {
	data := []opts.PieData{
		{Name: "Spendings", Value: fmt.Sprintf("%.2f", spendings*-1)},
		{Name: "Incomes", Value: fmt.Sprintf("%.2f", incomes)},
	}

	pieChart := charts.NewPie()
	pieChart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Total transactions"}))
	pieChart.AddSeries("", data)

	page := components.NewPage()
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

func setupCommandLineFlags() {
	flag.StringVar(&fileNameFlag, "file", "", "(Required) Path to the file containing bank statement lines")
	flag.StringVar(&fileNameFlag, "f", "", "Alias for -file")
	flag.BoolVar(&debugFlag, "debug", false, "(Optional) Enable debugging info")
	flag.BoolVar(&debugFlag, "d", false, "Alias for -debug")
	flag.Parse()
}

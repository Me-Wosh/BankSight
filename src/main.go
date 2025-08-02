package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var fileNameFlag string
var debugFlag bool

func main() {
	setupCommandLineFlags()

	fileName := getFileName()
	convertPdfToText(fileName)
	lines := readAllLines(strings.Replace(fileName, ".pdf", ".txt", 1))
	spendings, incomes := calculateTotalTransactions(lines)

	fmt.Printf("Spendings: %.2f\n", spendings)
	fmt.Printf("Incomes: %.2f", incomes)
}

func calculateTotalTransactions(lines []string) (spendings, incomes float64) {
	for _, line := range lines {
		sections := divideLineIntoSections(line)

		if len(sections) >= 4 {
			validFloat, err := convertToValidFloat(sections[3])

			if err != nil {
				if debugFlag {
					fmt.Println("INFO: Failed converting value. Expected a number, but got:", sections[3])
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

	return spendings, incomes
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
	} else if nonFlagArgs := flag.Args(); len(nonFlagArgs) > 1 {
		fileName = nonFlagArgs[1]
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

package main

import (
	"regexp"
	"strconv"
	"strings"
)

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

func getTimePeriod(lines []string) string {
	periodLine := strings.ToLower(lines[4])
	_, after, _ := strings.Cut(periodLine, "okres ")
	return after
}

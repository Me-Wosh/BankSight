package main

import (
	"strings"
)

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

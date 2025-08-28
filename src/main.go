package main

func main() {
	flags := setupCommandLineFlags()
	fileName := getFileName(flags.fileName)
	textFileName := convertPdfToText(fileName)
	lines := readAllLines(textFileName)
	timePeriod := getTimePeriod(lines)

	spendings, incomes, categoriesBalance := calculateTotalTransactions(lines, flags.debug)
	drawPieChart(spendings, incomes, categoriesBalance, timePeriod, fileName)
}

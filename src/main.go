package main

func main() {
	flags := setupCommandLineFlags()
	filePath := getFilePath(flags.filePath)
	textFilePath := convertPdfToText(filePath)
	lines := readAllLines(textFilePath)
	timePeriod := getTimePeriod(lines)

	transactions := calculateTotalTransactions(lines, flags.debug)
	drawPieChart(transactions, timePeriod, textFilePath)
}

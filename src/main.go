package main

func main() {
	flags := setupCommandLineFlags()
	fileName := getFileName(flags.fileName)
	textFileName := convertPdfToText(fileName)
	lines := readAllLines(textFileName)
	timePeriod := getTimePeriod(lines)

	transactions := calculateTotalTransactions(lines, flags.debug)
	drawPieChart(transactions, timePeriod, fileName)
}

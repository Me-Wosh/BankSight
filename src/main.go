package main

func main() {
	flags := setupCommandLineFlags()
	bank := getBank(flags.bank)
	textFilePath := bank.prepareFile(flags.filePath)
	lines := readAllLines(textFilePath)
	transactions := bank.calculateTotalTransactions(lines, flags.debug)
	drawPieChart(transactions, textFilePath, flags.openDefaultApp)
}

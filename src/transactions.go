package main

type transactions struct {
	timePeriod                      string
	spendings, incomes              float64
	previousBalance, closingBalance float64
	difference                      float64
	categorizedBalance              map[string]map[string]float64
}

type bank interface {
	prepareFile(filePath string) (string, error)
	calculateTotalTransactions(lines []string) (transactions, error)
	getSupportedFileExtensions() []string
}

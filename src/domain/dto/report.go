package dto

type Report struct {
	TotalBalance      float64
	TransactionsCount map[string]int
	AverageDebit      map[string]float64
	AverageCredit     map[string]float64
}

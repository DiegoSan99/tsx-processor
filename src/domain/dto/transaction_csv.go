package dto

import "time"

type TransactionCSV struct {
	Id          string
	Date        time.Time
	Transaction float64
	AccountId   string
	Type        string
}

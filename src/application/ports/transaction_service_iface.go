package ports

import "github.com/aws/aws-lambda-go/events"

type TransactionService interface {
	ProcessTransactions(event events.S3Event) error
}

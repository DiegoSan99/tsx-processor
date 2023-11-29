package s3adapter

import (
	"context"

	"github.com/DiegoSan99/transaction-processor/src/application/ports"
	"github.com/DiegoSan99/transaction-processor/src/application/repository/account"
	"github.com/DiegoSan99/transaction-processor/src/application/repository/transaction"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

type S3Adapter struct {
	accountRepo        account.AccountRepository
	transactionRepo    transaction.TransactionRepository
	logger             *zap.SugaredLogger
	transactionService ports.TransactionService
}

func NewS3Adapter(ts ports.TransactionService, accountRepo account.AccountRepository, transactionRepo transaction.TransactionRepository, logger *zap.SugaredLogger) *S3Adapter {
	return &S3Adapter{
		transactionService: ts,
		accountRepo:        accountRepo,
		transactionRepo:    transactionRepo,
		logger:             logger,
	}
}

func (adapter *S3Adapter) HandleS3Event(ctx context.Context, s3Event events.S3Event) error {
	adapter.logger.Info("S3 event received", s3Event)
	err := adapter.transactionService.ProcessTransactions(s3Event)
	if err != nil {
		adapter.logger.Errorw("Error processing transactions", "error", err)
		return err
	}
	return nil
}

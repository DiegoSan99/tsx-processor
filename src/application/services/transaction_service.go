package services

import (
	"context"

	"github.com/DiegoSan99/transaction-processor/src/application/ports"
	"github.com/DiegoSan99/transaction-processor/src/application/repository/account"
	"github.com/DiegoSan99/transaction-processor/src/application/repository/transaction"
	"github.com/DiegoSan99/transaction-processor/src/interfaces/email"
	"github.com/DiegoSan99/transaction-processor/src/interfaces/s3"
	"go.uber.org/zap"
)

type TransactionService struct {
	globalCtx       context.Context
	accountRepo     account.AccountRepository
	transactionRepo transaction.TransactionRepository
	s3Client        *s3.S3Client
	log             *zap.SugaredLogger
	emailClient     *email.SmtpClient
}

func NewTransactionService(globalCtx context.Context, accountRepo account.AccountRepository, transactionRepo transaction.TransactionRepository, s3Client *s3.S3Client, log *zap.SugaredLogger, email *email.SmtpClient) ports.TransactionService {
	return &TransactionService{
		globalCtx:       globalCtx,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		s3Client:        s3Client,
		log:             log,
		emailClient:     email,
	}
}

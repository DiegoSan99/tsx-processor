package main

import (
	"context"

	"github.com/DiegoSan99/transaction-processor/src/application/repository/account"
	"github.com/DiegoSan99/transaction-processor/src/application/repository/transaction"
	"github.com/DiegoSan99/transaction-processor/src/application/services"
	"github.com/DiegoSan99/transaction-processor/src/config"
	"github.com/DiegoSan99/transaction-processor/src/domain/entity"
	"github.com/DiegoSan99/transaction-processor/src/interfaces/db"
	"github.com/DiegoSan99/transaction-processor/src/interfaces/email"
	"github.com/DiegoSan99/transaction-processor/src/interfaces/envs"
	s3adapter "github.com/DiegoSan99/transaction-processor/src/interfaces/handlers/s3_adapter"
	"github.com/DiegoSan99/transaction-processor/src/interfaces/s3"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

func main() {

	var cfg config.AppConfig

	ctx := context.Background()
	ctx = envs.WithEnvs(ctx, &cfg)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	db.Connect(&cfg, sugar)
	defer db.Disconnect()

	smtpClient := email.NewSmtpClient(&cfg)

	db := db.GetDB()
	err := db.AutoMigrate(&entity.Account{}, &entity.Transaction{})
	if err != nil {
		sugar.Fatal(err)
	}

	accountRepo := account.NewGormAccountRepository(db)
	transactionRepo := transaction.NewGormTransactionRepository(db)

	s3Client, err := s3.New()
	if err != nil {
		sugar.Fatal(err)
	}

	ts := services.NewTransactionService(ctx, accountRepo, transactionRepo, s3Client, sugar, smtpClient)

	s3Adapter := s3adapter.NewS3Adapter(ts, accountRepo, transactionRepo, sugar)
	lambda.Start(s3Adapter.HandleS3Event)

}

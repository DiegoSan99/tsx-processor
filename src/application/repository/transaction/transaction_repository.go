package transaction

import "github.com/DiegoSan99/transaction-processor/src/domain/entity"

type TransactionRepository interface {
	Create(transaction *entity.Transaction) error
	GetByUserID(userID uint) ([]entity.Transaction, error)
}

package transaction

import (
	"github.com/DiegoSan99/transaction-processor/src/domain/entity"
	"gorm.io/gorm"
)

type GormTransactionRepository struct {
	Db *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) *GormTransactionRepository {
	return &GormTransactionRepository{Db: db}
}
func (repo *GormTransactionRepository) Create(transaction *entity.Transaction) error {
	return repo.Db.Create(transaction).Error
}

func (repo *GormTransactionRepository) GetByUserID(userID uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := repo.Db.Preload("Account").Where("account_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

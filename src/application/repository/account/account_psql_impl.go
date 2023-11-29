package account

import (
	"github.com/DiegoSan99/transaction-processor/src/domain/entity"
	"gorm.io/gorm"
)

type GormAccountRepository struct {
	Db *gorm.DB
}

func NewGormAccountRepository(db *gorm.DB) *GormAccountRepository {
	return &GormAccountRepository{Db: db}
}

func (repo *GormAccountRepository) Create(account *entity.Account) error {
	return repo.Db.Create(account).Error
}

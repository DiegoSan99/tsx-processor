package account

import "github.com/DiegoSan99/transaction-processor/src/domain/entity"

type AccountRepository interface {
	Create(account *entity.Account) error
}

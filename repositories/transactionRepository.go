package repositories

import (
	"gorm.io/gorm"
	"mnc-finance/entity"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) error
	FindByUserID(userID string) ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(transaction *entity.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByUserID(userID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

package services

import (
	"errors"
	"github.com/google/uuid"
	"mnc-finance/entity"
	"mnc-finance/models"
	"mnc-finance/queue"
	"mnc-finance/repositories"
	"time"
)

type TransactionService interface {
	TopUp(param *models.TopUp) (*entity.Transaction, error)
	Payment(userID uuid.UUID, amount float64, remarks string) (*entity.Transaction, error)
	Transfer(userID uuid.UUID, recipientID uuid.UUID, amount float64, remarks string) (*entity.Transaction, error)
	TransactionsReport(userID uuid.UUID) ([]entity.Transaction, error)
}

type transactionService struct {
	transactionRepository repositories.TransactionRepository
	userRepository        repositories.UserRepository
	queue                 queue.PublishDefinition
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, userRepo repositories.UserRepository, queue queue.PublishDefinition) TransactionService {
	return &transactionService{transactionRepo, userRepo, queue}
}

func (s *transactionService) TopUp(param *models.TopUp) (*entity.Transaction, error) {
	user, err := s.userRepository.FindByPhoneNumber(param.UserID)
	if err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "TopUp",
		Amount:        param.Amount,
		Remarks:       param.Remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance + param.Amount,
		CreatedAt:     time.Now(),
	}

	user.Balance += param.Amount
	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	if err := s.transactionRepository.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) Payment(userID uuid.UUID, amount float64, remarks string) (*entity.Transaction, error) {
	user, err := s.userRepository.FindByPhoneNumber(userID.String())
	if err != nil {
		return nil, err
	}

	if user.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	transaction := &entity.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "Payment",
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance - amount,
		CreatedAt:     time.Now(),
	}

	user.Balance -= amount
	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	if err := s.transactionRepository.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) Transfer(userID uuid.UUID, recipientID uuid.UUID, amount float64, remarks string) (*entity.Transaction, error) {
	user, err := s.userRepository.FindByPhoneNumber(userID.String())
	if err != nil {
		return nil, err
	}

	if user.Balance < amount {
		return nil, errors.New("Balance is not enough")
	}

	recipient, err := s.userRepository.FindByPhoneNumber(recipientID.String())
	if err != nil {
		return nil, errors.New("recipient not found")
	}

	transaction := &entity.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "Transfer",
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance - amount,
		CreatedAt:     time.Now(),
	}

	recipient.Balance += amount
	if err := s.userRepository.Update(recipient); err != nil {
		return nil, err
	}

	user.Balance -= amount
	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	if err := s.transactionRepository.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) TransactionsReport(userID uuid.UUID) ([]entity.Transaction, error) {
	return s.transactionRepository.FindByUserID(userID.String())
}

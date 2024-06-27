package services

import (
	"errors"
	"github.com/google/uuid"
	"mnc-finance/models"
	"mnc-finance/repositories"
)

type TransactionService interface {
	TopUp(userID uuid.UUID, amount float64, remarks string) (*models.Transaction, error)
	Payment(userID uuid.UUID, amount float64, remarks string) (*models.Transaction, error)
	Transfer(userID uuid.UUID, recipientID uuid.UUID, amount float64, remarks string) (*models.Transaction, error)
	TransactionsReport(userID uuid.UUID) ([]models.Transaction, error)
}

type transactionService struct {
	transactionRepository repositories.TransactionRepository
	userRepository        repositories.UserRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, userRepo repositories.UserRepository) TransactionService {
	return &transactionService{transactionRepo, userRepo}
}

func (s *transactionService) TopUp(userID uuid.UUID, amount float64, remarks string) (*models.Transaction, error) {
	user, err := s.userRepository.FindByPhoneNumber(userID.String())
	if err != nil {
		return nil, err
	}

	transaction := &models.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "TopUp",
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance + amount,
	}

	user.Balance += amount
	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	if err := s.transactionRepository.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) Payment(userID uuid.UUID, amount float64, remarks string) (*models.Transaction, error) {
	user, err := s.userRepository.FindByPhoneNumber(userID.String())
	if err != nil {
		return nil, err
	}

	if user.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	transaction := &models.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "Payment",
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance - amount,
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

func (s *transactionService) Transfer(userID uuid.UUID, recipientID uuid.UUID, amount float64, remarks string) (*models.Transaction, error) {
	user, err := s.userRepository.FindByPhoneNumber(userID.String())
	if err != nil {
		return nil, err
	}

	if user.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	recipient, err := s.userRepository.FindByPhoneNumber(recipientID.String())
	if err != nil {
		return nil, errors.New("recipient not found")
	}

	transaction := &models.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "Transfer",
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance - amount,
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

func (s *transactionService) TransactionsReport(userID uuid.UUID) ([]models.Transaction, error) {
	return s.transactionRepository.FindByUserID(userID.String())
}

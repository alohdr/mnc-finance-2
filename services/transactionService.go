package services

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mnc-finance/entity"
	"mnc-finance/models"
	"mnc-finance/queue"
	"mnc-finance/repositories"
	"mnc-finance/utils"
	"mnc-finance/utils/errorMessage"
	"time"
)

type TransactionService interface {
	TopUp(ctx *gin.Context, param *models.TopUp) (*entity.Transaction, error)
	Payment(ctx *gin.Context, param *models.Payment) (*models.Payment, error)
	Transfer(ctx *gin.Context, param *models.Transfer) (*models.Transaction, error)
	TransactionsReport(userID string) ([]entity.Transaction, error)
}

type transactionService struct {
	db                    *gorm.DB
	transactionRepository repositories.TransactionRepository
	userRepository        repositories.UserRepository
	queue                 queue.PublishDefinition
}

func NewTransactionService(db *gorm.DB, transactionRepo repositories.TransactionRepository, userRepo repositories.UserRepository, queue queue.PublishDefinition) TransactionService {
	return &transactionService{db, transactionRepo, userRepo, queue}
}

func (s *transactionService) TopUp(ctx *gin.Context, param *models.TopUp) (*entity.Transaction, error) {
	user, err := s.userRepository.FindByID(ctx.GetString("user_id"))
	if err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "TopUp",
		Amount:        param.Amount,
		Remarks:       "TopUp",
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance + param.Amount,
		Status:        utils.StatusSuccess,
		CreatedAt:     time.Now(),
	}

	tx := s.db.Begin()

	user.Balance += param.Amount
	if err := s.userRepository.Update(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.transactionRepository.Create(tx, transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	return transaction, tx.Commit().Error
}

func (s *transactionService) Payment(ctx *gin.Context, param *models.Payment) (*models.Payment, error) {
	user, err := s.userRepository.FindByID(ctx.GetString("user_id"))
	if err != nil {
		return nil, errorMessage.ErrFailedPayment
	}

	if user.Balance < param.Amount {
		return nil, errorMessage.ErrBalancePayment
	}

	transaction := &entity.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		Type:          "Payment",
		Amount:        param.Amount,
		Remarks:       param.Remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance - param.Amount,
		Status:        utils.StatusSuccess,
		CreatedAt:     time.Now(),
	}

	tx := s.db.Begin()

	user.Balance -= param.Amount
	if err := s.userRepository.Update(tx, user); err != nil {
		tx.Rollback()
		return nil, errorMessage.ErrFailedPayment
	}

	if err := s.transactionRepository.Create(tx, transaction); err != nil {
		tx.Rollback()
		return nil, errorMessage.ErrFailedPayment
	}

	param.PaymentID = transaction.ID.String()
	param.CreatedDate = transaction.CreatedAt
	param.BalanceBefore = transaction.BalanceBefore
	param.BalanceAfter = transaction.BalanceAfter

	return param, tx.Commit().Error
}

func (s *transactionService) Transfer(ctx *gin.Context, param *models.Transfer) (*models.Transaction, error) {
	userID := ctx.GetString("user_id")
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Balance < param.Amount {
		return nil, errors.New("Balance is not enough")
	}

	recipient, err := s.userRepository.FindByID(param.RecipientID)
	if err != nil {
		return nil, errors.New("recipient not found")
	}

	transaction := &entity.Transaction{
		ID:            uuid.New(),
		UserID:        user.ID,
		RecipientID:   recipient.ID,
		Type:          "Transfer",
		Amount:        param.Amount,
		Remarks:       param.Remarks,
		BalanceBefore: user.Balance,
		BalanceAfter:  user.Balance - param.Amount,
		CreatedAt:     time.Now(),
	}

	transactionComplete := entity.UserTransaction{
		UserObj:        *user,
		RecipientObj:   *recipient,
		TransactionObj: *transaction,
	}

	bytesParam, err := json.Marshal(transactionComplete)
	if err != nil {
		return nil, errorMessage.ErrInternalServerError
	}

	go func() {
		err := s.queue.ProduceMessage(utils.EventMncTransfer, utils.RouteMncTransfer, bytesParam)
		if err != nil {
			return
		}
	}()

	return &models.Transaction{
		TransferID:    transaction.ID.String(),
		Amount:        transaction.Amount,
		Remarks:       transaction.Remarks,
		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
		CreatedDate:   time.Time{},
	}, nil
}

func (s *transactionService) TransactionsReport(userID string) ([]entity.Transaction, error) {
	return s.transactionRepository.FindByUserID(userID)
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mnc-finance/models"
	"mnc-finance/services"
	"mnc-finance/utils"
	"net/http"
)

type TransactionController struct {
	transactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) *TransactionController {
	return &TransactionController{transactionService}
}

func (ctrl *TransactionController) TopUp(c *gin.Context) {
	input := new(models.TopUp)
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	transaction, err := ctrl.transactionService.TopUp(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to top up")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transaction})
}

func (ctrl *TransactionController) Payment(c *gin.Context) {
	var input struct {
		UserID  uuid.UUID `json:"user_id"`
		Amount  float64   `json:"amount"`
		Remarks string    `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	transaction, err := ctrl.transactionService.Payment(input.UserID, input.Amount, input.Remarks)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to make payment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transaction})
}

func (ctrl *TransactionController) Transfer(c *gin.Context) {
	var input struct {
		UserID      uuid.UUID `json:"user_id"`
		RecipientID uuid.UUID `json:"recipient_id"`
		Amount      float64   `json:"amount"`
		Remarks     string    `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	transaction, err := ctrl.transactionService.Transfer(input.UserID, input.RecipientID, input.Amount, input.Remarks)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to transfer")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transaction})
}

func (ctrl *TransactionController) TransactionsReport(c *gin.Context) {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	transactions, err := ctrl.transactionService.TransactionsReport(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transactions})
}

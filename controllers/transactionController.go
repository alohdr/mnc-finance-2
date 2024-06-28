package controllers

import (
	"github.com/gin-gonic/gin"
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

	transaction, err := ctrl.transactionService.TopUp(c, input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to top up")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transaction})
}

func (ctrl *TransactionController) Payment(c *gin.Context) {
	input := new(models.Payment)
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	transaction, err := ctrl.transactionService.Payment(c, input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transaction})
}

func (ctrl *TransactionController) Transfer(c *gin.Context) {
	input := new(models.Transfer)
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	transaction, err := ctrl.transactionService.Transfer(c, input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to transfer")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transaction})
}

func (ctrl *TransactionController) TransactionsReport(c *gin.Context) {
	transactions, err := ctrl.transactionService.TransactionsReport(c.GetString("user_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": transactions})
}

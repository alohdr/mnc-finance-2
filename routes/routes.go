package routes

import (
	"github.com/gin-gonic/gin"
	"mnc-finance/config"
	"mnc-finance/controllers"
	"mnc-finance/repositories"
	"mnc-finance/services"
)

func SetupRoutes(router *gin.Engine) {
	db := config.SetupDatabase()

	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	authService := services.NewAuthService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo, userRepo)

	authController := controllers.NewAuthController(authService)
	transactionController := controllers.NewTransactionController(transactionService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/refresh-token", authController.RefreshToken)
	router.POST("/topup", transactionController.TopUp)
	router.POST("/pay", transactionController.Payment)
	router.POST("/transfer", transactionController.Transfer)
	router.GET("/transactions", transactionController.TransactionsReport)
	router.PUT("/profile", authController.UpdateProfile) // Assuming you have UpdateProfile in authController
}

package routes

import (
	"github.com/gin-gonic/gin"
	"mnc-finance/config"
	"mnc-finance/controllers"
	"mnc-finance/queue"
	"mnc-finance/repositories"
	"mnc-finance/services"
)

func SetupRoutes(router *gin.Engine) {
	db := config.SetupDatabase()

	mq := config.SetUpRabbitMQ()

	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	rabbit := queue.NewPublishService(mq)

	authService := services.NewAuthService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo, userRepo, rabbit)

	authController := controllers.NewAuthController(authService)
	transactionController := controllers.NewTransactionController(transactionService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/refresh-token", authController.RefreshToken)
	router.POST("/topup", transactionController.TopUp)
	router.POST("/pay", transactionController.Payment)
	router.POST("/transfer", transactionController.Transfer)
	router.GET("/transactions", transactionController.TransactionsReport)
	router.PUT("/profile/:id", authController.UpdateProfile)
}

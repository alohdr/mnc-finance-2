package routes

import (
	"github.com/gin-gonic/gin"
	"mnc-finance/config"
	"mnc-finance/controllers"
	"mnc-finance/middlewares"
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
	transactionService := services.NewTransactionService(db, transactionRepo, userRepo, rabbit)

	authController := controllers.NewAuthController(authService)
	transactionController := controllers.NewTransactionController(transactionService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/refresh-token", authController.RefreshToken)

	// Use AuthMiddleware for routes requiring authentication
	auth := router.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("/topup", transactionController.TopUp)
		auth.POST("/pay", transactionController.Payment)
		auth.POST("/transfer", transactionController.Transfer)
		auth.GET("/transactions", transactionController.TransactionsReport)
		auth.PUT("/profile", authController.UpdateProfile)
	}
}

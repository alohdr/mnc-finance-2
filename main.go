package main

import (
	"github.com/gin-gonic/gin"
	"mnc-finance/config"
	"mnc-finance/routes"
)

func main() {
	router := gin.Default()
	config.SetupDatabase()
	routes.SetupRoutes(router)
	router.Run(":8080")
}

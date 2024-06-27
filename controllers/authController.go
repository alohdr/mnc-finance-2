package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mnc-finance/config"
	"mnc-finance/models"
	"mnc-finance/utils"
	"net/http"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}
	input.ID = uuid.New()
	input.Balance = 0

	if err := config.DB.Create(&input).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": input})
}

func Login(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number"`
		PIN         string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	var user models.User
	if err := config.DB.Where("phone_number = ?", input.PhoneNumber).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Phone Number and PIN doesn’t match.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(input.PIN)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Phone Number and PIN doesn’t match.")
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	user.RefreshToken = refreshToken
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": gin.H{"access_token": accessToken, "refresh_token": refreshToken}})
}

func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	var user models.User
	if err := config.DB.Where("refresh_token = ?", input.RefreshToken).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	user.RefreshToken = refreshToken
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": gin.H{"access_token": accessToken, "refresh_token": refreshToken}})
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"mnc-finance/models"
	"mnc-finance/services"
	"mnc-finance/utils"
	"net/http"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	input := new(models.User)
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	newUser, err := ctrl.authService.Register(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": newUser})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	input := new(models.Login)
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	accessToken, refreshToken, err := ctrl.authService.Login(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Phone Number and PIN doesn’t match")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": gin.H{"access_token": accessToken, "refresh_token": refreshToken}})
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	accessToken, refreshToken, err := ctrl.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": gin.H{"access_token": accessToken, "refresh_token": refreshToken}})
}

func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	input := new(models.User)
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	updatedUser, err := ctrl.authService.Update(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": updatedUser})
}

package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mnc-finance/entity"
	"mnc-finance/models"
	"mnc-finance/repositories"
	"mnc-finance/utils"
	"mnc-finance/utils/errorMessage"
	"strings"
	"time"
)

type AuthService interface {
	Register(user *models.User) (*entity.User, error)
	Login(param *models.Login) (string, string, error)
	Update(ctx *gin.Context, user *models.Profile) (*models.Profile, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(user *models.User) (*entity.User, error) {
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(user.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, errorMessage.ErrFailedRegister
	}

	param := &entity.User{
		ID:          uuid.New(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		PIN:         string(hashedPIN),
		CreatedAt:   time.Now(),
	}

	if err := s.userRepository.Create(param); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errorMessage.ErrUserExist
		}
		return nil, errorMessage.ErrFailedRegister
	}
	return param, nil
}

func (s *authService) Login(param *models.Login) (string, string, error) {
	user, err := s.userRepository.FindByPhoneNumber(param.PhoneNumber)
	if err != nil {
		return "", "", errors.New("phone number and PIN doesn’t match")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(param.PIN)); err != nil {
		return "", "", errors.New("phone number and PIN doesn’t match")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = refreshToken
	if err := s.userRepository.UpdateUser(user); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	user, err := s.userRepository.FindByRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	accessToken, newRefreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = newRefreshToken
	if err := s.userRepository.UpdateUser(user); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *authService) Update(ctx *gin.Context, user *models.Profile) (*models.Profile, error) {
	userID := ctx.GetString("user_id")
	existingUser, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Address = user.Address
	user.UpdateDate = time.Now()
	user.UserID = userID

	err = s.userRepository.UpdateUser(existingUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

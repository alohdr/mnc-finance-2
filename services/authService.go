package services

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mnc-finance/entity"
	"mnc-finance/models"
	"mnc-finance/repositories"
	"mnc-finance/utils"
	"time"
)

type AuthService interface {
	Register(user *models.UserRegister) (*entity.User, error)
	Login(phoneNumber, pin string) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(user *models.UserRegister) (*entity.User, error) {
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(user.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return param, nil
}

func (s *authService) Login(phoneNumber, pin string) (string, string, error) {
	user, err := s.userRepository.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return "", "", errors.New("phone number and PIN doesn’t match")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(pin)); err != nil {
		return "", "", errors.New("phone number and PIN doesn’t match")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = refreshToken
	if err := s.userRepository.Update(user); err != nil {
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
	if err := s.userRepository.Update(user); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

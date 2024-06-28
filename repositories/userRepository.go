package repositories

import (
	"gorm.io/gorm"
	"mnc-finance/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id string) (*entity.User, error)
	FindByPhoneNumber(phoneNumber string) (*entity.User, error)
	FindByRefreshToken(refreshToken string) (*entity.User, error)
	Update(tx *gorm.DB, user *entity.User) error
	UpdateUser(user *entity.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhoneNumber(phoneNumber string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByRefreshToken(refreshToken string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(tx *gorm.DB, user *entity.User) error {
	return tx.Save(user).Error
}

func (r *userRepository) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}

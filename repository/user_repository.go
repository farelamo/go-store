package repository

import (
	"store/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(username string) (model.User, error)
	Create(user model.User) error
}

type userRepository struct {
	mysql *gorm.DB
}

func NewUserRepository(mysql *gorm.DB) UserRepository {
	return &userRepository{
		mysql: mysql,
	}
}

func (u *userRepository) GetByUsername(username string) (model.User, error) {
	var user model.User
	if err := u.mysql.Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userRepository) Create(user model.User) error {
	if err := u.mysql.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
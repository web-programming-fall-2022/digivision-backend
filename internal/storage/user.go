package storage

import (
	"errors"
	"gorm.io/gorm"
)

type UserAccount struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex:idx_email"`
	PhoneNumber  string `gorm:"uniqueIndex:idx_phone_number"`
	Gender       string `gorm:"type:VARCHAR(1)"`
	FirstName    string
	LastName     string
	PasswordHash string
}

func (storage *Storage) CreateUser(user *UserAccount) error {
	if err := storage.DB.Create(user).Error; err != nil {
		return errors.New("couldn't create user in postgres storage")
	}
	return nil
}

func (storage *Storage) GetUserByID(id uint) (*UserAccount, error) {
	user := UserAccount{}
	storage.DB.First(&user, id)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (storage *Storage) GetUserByEmail(email string) (*UserAccount, error) {
	user := UserAccount{Email: email}
	storage.DB.Where(&user).First(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (storage *Storage) GetUserByPhoneNumber(phoneNumber string) (*UserAccount, error) {
	user := UserAccount{PhoneNumber: phoneNumber}
	storage.DB.Where(&user).First(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

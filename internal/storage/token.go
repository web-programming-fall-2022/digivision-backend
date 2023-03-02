package storage

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type UnauthorizedToken struct {
	gorm.Model
	UserID     uint
	User       UserAccount `gorm:"ONDELETE:CASCADE"`
	Token      string
	Expiration time.Time `gorm:"index"`
}

func (storage *Storage) CreateUnauthorizedToken(token *UnauthorizedToken) error {
	if err := storage.DB.Create(token).Error; err != nil {
		return errors.New("couldn't create unauthorized toen in postgres storage")
	}
	return nil
}

func (storage *Storage) GetUnauthorizedToken(token string) (*UnauthorizedToken, error) {
	unauthorizedToken := UnauthorizedToken{}
	storage.DB.Where("token = ?", token).First(&unauthorizedToken)
	if unauthorizedToken.ID == 0 {
		return nil, errors.New("unauthorized token not found")
	}
	return &unauthorizedToken, nil
}

func (storage *Storage) DeleteUnauthorizedTokensBeforeNow() error {
	if err := storage.DB.Where("expiration < ?", time.Now()).Delete(&UnauthorizedToken{}).Error; err != nil {
		return errors.New("couldn't delete unauthorized tokens in postgres storage")
	}
	return nil
}

func (storage *Storage) GetAllUnauthorizedTokens() ([]*UnauthorizedToken, error) {
	unauthorizedTokens := make([]*UnauthorizedToken, 0)
	if err := storage.DB.Find(&unauthorizedTokens).Error; err != nil {
		return nil, errors.New("couldn't get all unauthorized tokens in postgres storage")
	}
	return unauthorizedTokens, nil
}

package storage

import (
	"errors"
	"gorm.io/gorm"
)

type SearchHistory struct {
	gorm.Model
	UserID       uint
	User         UserAccount `gorm:"ONDELETE:CASCADE"`
	QueryAddress string
	Results      []SearchHistoryResult
}

type SearchHistoryResult struct {
	gorm.Model
	SearchHistoryID uint
	SearchHistory   SearchHistory `gorm:"ONDELETE:CASCADE"`
	ProductID       uint
}

func (storage *Storage) CreateSearchHistory(history *SearchHistory) error {
	if err := storage.DB.Create(history).Error; err != nil {
		return err
	}
	return nil
}

func (storage *Storage) GetSearchHistoryByID(id uint) (*SearchHistory, error) {
	history := SearchHistory{}
	storage.DB.First(&history, id).Preload("Results")
	if history.ID == 0 {
		return nil, errors.New("history not found")
	}
	return &history, nil
}

func (storage *Storage) GetSearchHistoryByUserID(userID uint, offset int, limit int) ([]SearchHistory, error) {
	history := make([]SearchHistory, 0)
	storage.DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&history)
	if len(history) == 0 {
		return nil, errors.New("history not found")
	}
	return history, nil
}

func (storage *Storage) CreateSearchHistoryResult(result *SearchHistoryResult) error {
	if err := storage.DB.Create(result).Error; err != nil {
		return err
	}
	return nil
}

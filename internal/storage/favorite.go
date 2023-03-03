package storage

import (
	"errors"
	"gorm.io/gorm"
)

type FavoriteList struct {
	gorm.Model
	UserID uint
	User   UserAccount `gorm:"ONDELETE:CASCADE"`
	Name   string
	items  []FavoriteListItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type FavoriteListItem struct {
	gorm.Model
	FavoritesListID uint
	FavoritesList   FavoriteList `gorm:"ONDELETE:CASCADE"`
	ProductID       uint
}

func (storage *Storage) CreateFavoriteList(list *FavoriteList) error {
	if err := storage.DB.Create(list).Error; err != nil {
		return err
	}
	return nil
}

func (storage *Storage) GetFavoriteListByID(id uint) (*FavoriteList, error) {
	list := FavoriteList{}
	storage.DB.First(&list, id).Preload("items")
	if list.ID == 0 {
		return nil, errors.New("list not found")
	}
	return &list, nil
}

func (storage *Storage) GetFavoriteListsByUserID(userID uint) (*[]FavoriteList, error) {
	list := make([]FavoriteList, 0)
	storage.DB.Where("user_id = ?", userID).Find(&list)
	if len(list) == 0 {
		return nil, errors.New("list not found")
	}
	return &list, nil
}

func (storage *Storage) AddItemToList(listID uint, productID uint) error {
	item := FavoriteListItem{
		FavoritesListID: listID,
		ProductID:       productID,
	}
	if err := storage.DB.Create(&item).Error; err != nil {
		return err
	}
	return nil
}

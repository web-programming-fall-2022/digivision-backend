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
	Items  []FavoriteListItem
}

type FavoriteListItem struct {
	gorm.Model
	FavoriteListID uint
	FavoriteList   FavoriteList `gorm:"ONDELETE:CASCADE"`
	ProductID      uint
}

func (storage *Storage) CreateFavoriteList(list *FavoriteList) error {
	if err := storage.DB.Create(list).Error; err != nil {
		return err
	}
	return nil
}

func (storage *Storage) GetFavoriteListByID(id uint) (*FavoriteList, error) {
	list := FavoriteList{}
	storage.DB.First(&list, id).Preload("Items")
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

func (storage *Storage) GetFavoriteListByUserIDAndName(userID uint, name string) (*FavoriteList, error) {
	list := FavoriteList{}
	storage.DB.Where("user_id = ? AND name = ?", userID, name).Preload("Items").First(&list)
	if list.ID == 0 {
		return nil, errors.New("list not found")
	}
	return &list, nil
}

func (storage *Storage) AddItemToList(listID uint, productID uint) error {
	item := FavoriteListItem{
		FavoriteListID: listID,
		ProductID:      productID,
	}
	if err := storage.DB.Create(&item).Error; err != nil {
		return err
	}
	return nil
}

func (storage *Storage) RemoveItemFromList(listID uint, productID uint) error {
	item := FavoriteListItem{}
	storage.DB.Where("favorite_list_id = ? AND product_id = ?", listID, productID).First(&item)
	if item.ID == 0 {
		return errors.New("item not found")
	}
	item = FavoriteListItem{}
	if err := storage.DB.Where(
		"favorite_list_id = ? AND product_id = ?",
		listID, productID,
	).Delete(&item).Error; err != nil {
		return err
	}
	return nil
}

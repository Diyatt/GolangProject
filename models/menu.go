package models

import "github.com/Diyatt/GolangProject/database"

type MenuItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func CreateMenuItem(item *MenuItem) error {
	return database.DB.Create(item).Error
}

func GetMenuItemByID(id uint) (*MenuItem, error) {
	var item MenuItem
	if err := database.DB.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func UpdateMenuItem(item *MenuItem) error {
	return database.DB.Save(item).Error
}

func DeleteMenuItem(id uint) error {
	return database.DB.Delete(&MenuItem{}, id).Error
}

func GetAllMenuItems() ([]*MenuItem, error) {
	var menuItems []*MenuItem
	if err := database.DB.Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

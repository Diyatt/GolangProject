package models

import (
	"github.com/Diyatt/GolangProject/database"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

func GetOrdersByUserID(userID uint) ([]Order, error) {
	var orders []Order
	if err := database.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func CreateOrder(db *gorm.DB, order *Order) error {
	return database.DB.Create(order).Error
}

func GetOrderByID(orderID uint) (*Order, error) {
	var order Order
	if err := database.DB.Where("id =?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func UpdateOrder(order *Order) error {
	return database.DB.Save(order).Error
}

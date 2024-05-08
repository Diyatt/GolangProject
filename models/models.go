package models

import (
	"github.com/Diyatt/GolangProject/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID      uint        `json:"user_id"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Orders   []Order
}

type Product struct {
	gorm.Model
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
}

func (user *User) CreateUserRecord() error {
	result := database.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateOrder(order *Order) error {
	if err := database.DB.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func GetOrdersByUserID(userID uint) ([]Order, error) {
	var orders []Order
	if err := database.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func GetOrderDetailsByID(orderID uint) (*Order, error) {
	var order Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func ModifyOrder(orderID uint, updatedOrder *Order) error {

	var existingOrder Order
	if err := database.DB.First(&existingOrder, orderID).Error; err != nil {
		return err
	}

	if err := database.DB.Model(&existingOrder).Updates(updatedOrder).Error; err != nil {
		return err
	}

	return nil
}

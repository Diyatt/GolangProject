package models

import (
	"errors"

	"github.com/Diyatt/GolangProject/database"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	Admin  Role = "admin"
	Client Role = "client"
)

type Status string

const (
	Canceled   Status = "canceled"
	Inproccess Status = "inproccess"
	Ready      Status = "ready"
	Done       Status = "done"
)

type Order struct {
	gorm.Model
	UserID      uint         `json:"user_id"`
	Status      Status       `json:"status"`
	TotalAmount float64      `json:"total_amount"`
	Items       []*OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID    uint            `json:"order_id"`
	MenuItemID uint            `json:"menu_item_id"`
	Quantity   int             `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
	Order      Order           `gorm:"foreignKey:OrderID"`
	MenuItem   MenuItem        `gorm:"foreignKey:MenuItemID"`
}

type MenuItem struct {
	gorm.Model
	Name     string       `json:"name"`
	Quantity uint         `json:"quantity"`
	Price    float64      `json:"price"`
	Orders   []*OrderItem `gorm:"foreignKey:MenuItemID"`
}

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required"`
	Role     Role
}

type OrderRequest struct {
	OrderItems []struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	} `json:"order_items"`
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

func CreateMenuItem(menuItem *MenuItem) error {
	if err := database.DB.Create(menuItem).Error; err != nil {
		return err
	}
	return nil
}

func GetMenuItems() ([]MenuItem, error) {
	var menuItems []MenuItem
	if err := database.DB.Find(&menuItems).Error; err != nil {
		return nil, err
	}
	return menuItems, nil
}

func GetMenuItemByID(id uint) (*MenuItem, error) {
	var menuItem MenuItem
	if err := database.DB.First(&menuItem, id).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

func UpdateMenuItem(menuItem *MenuItem) error {
	if err := database.DB.Save(menuItem).Error; err != nil {
		return err
	}
	return nil
}

func DeleteMenuItem(id uint) error {
	if err := database.DB.Delete(&MenuItem{}, id).Error; err != nil {
		return err
	}
	return nil
}

func GetOrderItemsByOrderID(orderID uint) ([]*OrderItem, error) {
	var orderItems []*OrderItem
	if err := database.DB.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	switch u.Role {
	case Admin, Client:
		return nil
	default:
		return errors.New("invalid user role")
	}
}
func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	switch o.Status {
	case Canceled, Inproccess, Ready, Done:
		return nil
	default:
		return errors.New("invalid order status")
	}
}

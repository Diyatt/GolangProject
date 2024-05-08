package models

import (
	"github.com/Diyatt/GolangProject/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Order структура для представления заказа.
type Order struct {
	gorm.Model
	UserID      uint        `json:"user_id"`                         // ID пользователя, совершившего заказ
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"` // Список элементов заказа
	TotalAmount float64     `json:"total_amount"`                    // Общая сумма заказа
	Status      string      `json:"status"`                          // Статус заказа (например, "placed", "processing", "completed")
}

// OrderItem структура для представления элемента заказа.
type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`   // ID заказа, к которому относится элемент
	ProductID uint    `json:"product_id"` // ID продукта
	Quantity  uint    `json:"quantity"`   // Количество продукта
	Price     float64 `json:"price"`      // Цена за единицу продукта
}

// User структура для представления пользователя.
type User struct {
	gorm.Model
	Name     string  `json:"name"`                // Имя пользователя
	Email    string  `json:"email" gorm:"unique"` // Email пользователя (должен быть уникальным)
	Password string  `json:"-"`                   // Пароль пользователя (не будет отображаться в JSON)
	Orders   []Order // Список заказов, совершенных пользователем
}

// Product структура для представления продукта.
type Product struct {
	gorm.Model
	Name     string  `json:"name"`     // Название продукта
	Price    float64 `json:"price"`    // Цена продукта
	Quantity uint    `json:"quantity"` // Количество доступных продуктов
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

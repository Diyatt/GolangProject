package models

import "gorm.io/gorm"

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

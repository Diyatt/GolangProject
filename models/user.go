package models

import (
	"github.com/Diyatt/GolangProject/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required"`
	Role     string `gorm:"not null"`
}

func (user *User) CreateUserRecord() error {
	result := database.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

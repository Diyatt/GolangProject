package utils

import (
	"github.com/Diyatt/GolangProject/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(user *models.User, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err

	}
	user.Password = string(bytes)
	return nil
}

func CheckPassword(providedPassword string, user *models.User) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

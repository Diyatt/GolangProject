package main

import (
	"fmt"
	"log"

	"github.com/Diyatt/GolangProject/database"
	"github.com/Diyatt/GolangProject/models"
	"github.com/Diyatt/GolangProject/routes"
)

func main() {
	err := database.InitDb()

	fmt.Println("Connecting to database")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")

	err = database.DB.AutoMigrate(&models.User{}, &models.Order{})

	if err != nil {
		panic("Failed to migrate database")
	}

	routes.RunServer()

}

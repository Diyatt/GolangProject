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

	err = database.DB.AutoMigrate(&models.User{}, &models.Order{}, &models.MenuItem{})

	if err != nil {
		panic("Failed to migrate database")
	}

	createSampleItems()

	routes.RunServer()

}

func createSampleItems() {

	items := []models.MenuItem{
		{Name: "Product A", Price: 19.99},
		{Name: "Product B", Price: 29.99},
		{Name: "Product C", Price: 9.99},
	}

	for _, item := range items {
		if err := database.DB.Create(&item).Error; err != nil {
			log.Fatal("Error creating item:", err)
		}
	}
	fmt.Println("Sample items created successfully")
}

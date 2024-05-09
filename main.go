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

	menu := []models.MenuItem{
		{Name: "Gamburger", Price: 5.99},
		{Name: "Cheeseburger", Price: 6.49},
		{Name: "Chicken Sandwich", Price: 7.99},
		{Name: "French Fries", Price: 2.99},
		{Name: "Onion Rings", Price: 3.49},
		{Name: "Soda", Price: 1.99},
		{Name: "Iced Tea", Price: 1.99},
		{Name: "Milkshake", Price: 4.99},
		// Add more menu items as needed
	}
	fmt.Println("Menu:", menu)
	for _, item := range menu {
		if err := database.DB.Create(&item).Error; err != nil {
			panic("failed to insert menu items")
		}
	}
	// var menuItems []models.MenuItem
	// if err := database.DB.Find(&menuItems).Error; err != nil {
	// 	panic("failed to fetch menu items")
	// }

	// // Display menu items
	// fmt.Println("Menu:")
	// for i, item := range menuItems {
	// 	fmt.Printf("%d. %s - $%.2f\n", i+1, item.Name, item.Price)
	// }

	routes.RunServer()

}

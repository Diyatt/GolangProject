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

<<<<<<< HEAD
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
=======
	createSampleItems()
>>>>>>> 0be092e4135aaf3f20dfbc7d8a5c2b4a6ea03c97

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

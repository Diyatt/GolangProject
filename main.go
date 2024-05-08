package main

import (
	"fmt"
	"log"

	"github.com/Diyatt/GolangProject/controllers"
	"github.com/Diyatt/GolangProject/database"
	"github.com/Diyatt/GolangProject/middleware"
	"github.com/Diyatt/GolangProject/models"
	"github.com/Diyatt/GolangProject/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	err := database.InitDb()

	fmt.Println("Connecting to database")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")

	err = database.DB.AutoMigrate(&models.User{}, &models.Order{}, &models.OrderItem{}, &models.Product{})

	if err != nil {
		panic("Failed to migrate database")
	}

	routes.RunServer() // run server with gin

}
func setupRouter() *gin.Engine {
	// Create a new router
	r := gin.Default()
	// Add a welcome route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome To This Website")
	})
	// Create a new group for the API
	api := r.Group("/api")
	{
		// Create a new group for the public routes
		public := api.Group("/public")
		{
			// Add the login route
			public.POST("/login", controllers.Login)
			// Add the signup route
			public.POST("/signup", controllers.Signup)
		}
		// Add the signup route
		protected := api.Group("/protected").Use(middlewares.Authz())
		{
			// Add the profile route
			protected.GET("/profile", controllers.Profile)
		}
	}
	// Return the router
	return r
}

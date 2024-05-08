package routes

import (
	"github.com/Diyatt/GolangProject/controllers"
	"github.com/Diyatt/GolangProject/middleware"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	router.POST("/signup", controllers.SignUp)
	router.POST("/signin", controllers.SignIn)

	authenticated := router.Group("/")
	authenticated.Use(middleware.Authz())
	{
		authenticated.GET("/orders", controllers.GetUsersOrders)
		authenticated.POST("/order", controllers.PlaceOrder)
		authenticated.PUT("/orders/:id", controllers.ModifyOrder)
		authenticated.GET("/orders/:id", controllers.GetOrderDetails)
	}

	return router
}

func RunServer() {
	router := setupRouter()
	router.Run(":8080")
}

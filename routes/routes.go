package routes

import (
	"github.com/Diyatt/GolangProject/controllers"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	router.POST("/signup", controllers.SignUp)
	router.POST("/signin", controllers.SignIn)

	router.GET("/orders", controllers.GetUsersOrders)
	router.POST("/order", controllers.PlaceOrder)
	router.PUT("/orders/:id", controllers.ModifyOrder)
	router.GET("/orders/:id", controllers.GetOrderDetails)

	return router
}

func RunServer() {
	router := setupRouter()
	router.Run(":8080")
}

package routes

import (
	"github.com/Diyatt/GolangProject/controllers"
	"github.com/Diyatt/GolangProject/middleware"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)

	authenticated := router.Group("/")
	authenticated.Use(middleware.Authz())
	{
		authenticated.GET("/order", controllers.GetUsersOrders)
		authenticated.POST("/order", controllers.PostUsersOrders)
		authenticated.DELETE("/order/:id", controllers.DeleteOrder)
		authenticated.PATCH("/order/:id", controllers.UpdateOrder)
		authenticated.POST("/menu", controllers.CreateMenuItem)

		authenticated.GET("/menu", controllers.GetMenuItems)

		authenticated.PUT("/menu/:id", controllers.UpdateMenuItem)

		authenticated.DELETE("/menu/:id", controllers.DeleteMenuItem)
	}

	return router
}

func RunServer() {
	router := setupRouter()
	router.Run(":8080")
}

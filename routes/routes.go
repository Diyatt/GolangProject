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
	authenticated.Use(middleware.RequireAuth)
	{
		authenticated.GET("/orders", controllers.GetOrders)
		authenticated.POST("/order", controllers.CreateOrder)
		authenticated.PUT("/orders/:id", controllers.ModifyOrder)
		authenticated.GET("/orders/:id", controllers.GetOrderDetails)
		authenticated.GET("/menu", controllers.GetMenuItems)

		admin := authenticated.Group("admin")
		admin.Use(middleware.IsAdmin)

		{
			admin.POST("/menu", controllers.CreateMenuItem)
			admin.PUT("/menu/:id", controllers.ModifyMenuItem)
			admin.DELETE("/menu/:id", controllers.DeleteMenuItem)
		}
	}

	return router
}

func RunServer() {
	router := setupRouter()
	router.Run(":8080")
}

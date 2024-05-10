package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Diyatt/GolangProject/auth"
	"github.com/Diyatt/GolangProject/database"
	"github.com/Diyatt/GolangProject/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func PostUsersOrders(c *gin.Context) {

	userEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user email from context"})
		return
	}
	user, err := models.GetUserByEmail(userEmail.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	var req models.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		UserID: user.ID,
		Status: models.Inproccess,
		Items:  []*models.OrderItem{},
	}

	var totalAmount float64

	for _, item := range req.OrderItems {
		menuItem, err := models.GetMenuItemByID(item.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve menu item"})
			return
		}
		if int(menuItem.Quantity) < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough quantity available for menu item"})
			return
		}

		orderItemPrice := menuItem.Price * float64(item.Quantity)
		orderItemPriceDec := decimal.NewFromFloat(orderItemPrice)

		orderItem := &models.OrderItem{
			MenuItemID: item.ProductID,
			Quantity:   item.Quantity,
			Price:      orderItemPriceDec,
		}
		order.Items = append(order.Items, orderItem)

		menuItem.Quantity -= uint(item.Quantity)
		if err := database.DB.Save(&menuItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu item quantity"})
			return
		}

		totalAmount += orderItemPrice
	}

	order.TotalAmount = totalAmount

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully!", "order_id": order.ID})
}

func sendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

func GetUsersOrders(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		sendError(c, http.StatusInternalServerError, "Failed to get user email from context")
		return
	}

	user, err := models.GetUserByEmail(email.(string))
	if err != nil {
		sendError(c, http.StatusInternalServerError, "Failed to find user")
		return
	}

	orders, err := models.GetOrdersByUserID(user.ID)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "Failed to fetch user orders")
		return
	}

	var responseOrders []gin.H
	for _, order := range orders {
		orderItems, err := models.GetOrderItemsByOrderID(order.ID)
		if err != nil {
			sendError(c, http.StatusInternalServerError, "Failed to fetch order items")
			return
		}

		var items []gin.H
		for _, item := range orderItems {
			totalPrice := item.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))

			formattedItem := gin.H{
				"MenuItem": gin.H{
					"ID":    item.MenuItemID,
					"Price": item.Price,
				},
				"Quantity": item.Quantity,
				"Price":    totalPrice,
			}
			items = append(items, formattedItem)
		}

		responseOrder := gin.H{
			"CreatedAt":   order.CreatedAt,
			"ID":          order.ID,
			"UpdatedAt":   order.UpdatedAt,
			"Status":      order.Status,
			"TotalAmount": fmt.Sprintf("%.2f", order.TotalAmount),
			"Items":       items,
		}
		responseOrders = append(responseOrders, responseOrder)
	}

	c.JSON(http.StatusOK, responseOrders)
}

func GetOrderDetails(c *gin.Context) {
	_, exists := c.Get("email")
	if !exists {
		sendError(c, http.StatusInternalServerError, "Failed to get user email from context")
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sendError(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := models.GetOrderDetailsByID(uint(orderID))
	if err != nil {
		sendError(c, http.StatusInternalServerError, "Failed to fetch order details")
		return
	}

	c.JSON(http.StatusOK, order)
}
func DeleteOrder(c *gin.Context) {

	userEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user email from context"})
		return
	}

	user, err := models.GetUserByEmail(userEmail.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	result := database.DB.Preload("Items").First(&order, orderID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if order.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own orders"})
		return
	}

	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
func UpdateOrder(c *gin.Context) {

	var updateInfo struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userRole, _ := c.Get("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update order status"})
		return
	}

	orderID := c.Param("id")
	newStatus := models.Status(updateInfo.Status)

	switch newStatus {
	case models.Canceled, models.Inproccess, models.Ready, models.Done:
		var order models.Order
		result := database.DB.First(&order, orderID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		order.Status = newStatus
		if err := database.DB.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
	}

}

type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}

func Signup(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs ",
		})
		c.Abort()
		return
	}
	err = user.HashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{
			"Error": "Error Hashing Password",
		})
		c.Abort()
		return
	}
	err = user.CreateUserRecord()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Creating User",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Message": "Sucessfully Register",
	})
}

func Login(c *gin.Context) {
	var payload LoginPayload
	var user models.User
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs",
		})
		c.Abort()
		return
	}
	result := database.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(401, gin.H{
			"Error": "Invalid User Credentials",
		})
		c.Abort()
		return
	}
	err = user.CheckPassword(payload.Password)
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{
			"Error": "Invalid User Credentials",
		})
		c.Abort()
		return
	}
	jwtWrapper := auth.JwtWrapper{
		SecretKey:         "verysecretkey",
		Issuer:            "AuthService",
		ExpirationMinutes: 1,
		ExpirationHours:   12,
	}
	signedToken, err := jwtWrapper.GenerateToken(user.Email, string(user.Role))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}
	signedtoken, err := jwtWrapper.RefreshToken(user.Email, string(user.Role))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}
	tokenResponse := LoginResponse{
		Token:        signedToken,
		RefreshToken: signedtoken,
	}
	c.JSON(200, tokenResponse)
}

func CreateMenuItem(c *gin.Context) {
	userRole, exists := c.Get("role")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user role is admin
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create menu items"})
		return
	}

	// Proceed with creating menu item
	var menuItem models.MenuItem
	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateMenuItem(&menuItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item"})
		return
	}

	c.JSON(http.StatusCreated, menuItem)
}

func GetMenuItems(c *gin.Context) {
	menuItems, err := models.GetMenuItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items"})
		return
	}

	c.JSON(http.StatusOK, menuItems)
}

func UpdateMenuItem(c *gin.Context) {
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update menu items"})
		return
	}

	menuItemID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var menuItem models.MenuItem
	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menuItem.ID = uint(menuItemID)
	if err := models.UpdateMenuItem(&menuItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu item"})
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

func DeleteMenuItem(c *gin.Context) {
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can delete menu items"})
		return
	}

	menuItemID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := models.DeleteMenuItem(uint(menuItemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}

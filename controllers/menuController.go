package controllers

import (
	"net/http"
	"strconv"

	"github.com/Diyatt/GolangProject/models"
	"github.com/gin-gonic/gin"
)

func GetMenuItems(c *gin.Context) {
	menuItems, err := models.GetAllMenuItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch menu items"})
		return
	}
	c.JSON(http.StatusOK, menuItems)
}

func CreateMenuItem(c *gin.Context) {
	var newItem models.MenuItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateMenuItem(&newItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create menu item"})
		return
	}

	c.JSON(http.StatusCreated, newItem)
}

func ModifyMenuItem(c *gin.Context) {
	menuItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	existingMenuItem, err := models.GetMenuItemByID(uint(menuItemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch existing menu item"})
		return
	}

	var updatedMenuItem models.MenuItem
	if err := c.ShouldBindJSON(&updatedMenuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	existingMenuItem.Name = updatedMenuItem.Name
	existingMenuItem.Description = updatedMenuItem.Description
	existingMenuItem.Price = updatedMenuItem.Price

	if err := models.UpdateMenuItem(existingMenuItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update menu item"})
		return
	}

	c.JSON(http.StatusOK, existingMenuItem)

}

func DeleteMenuItem(c *gin.Context) {
	menuItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID"})
		return
	}

	if err := models.DeleteMenuItem(uint(menuItemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete menu item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "menu item deleted successfully"})
}

package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Diyatt/GolangProject/database"
	"github.com/Diyatt/GolangProject/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if err != nil || tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
		return
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		database.DB.First(&user, claims["userID"])
		c.Set("user", user)
		c.Set("role", user.Role)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

func IsAdmin(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	if role.(string) != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

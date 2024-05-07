package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	secretKey       = "your_secret_key"
	tokenExpiration = 3600 // in seconds (1 hour)
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	storedHashedPassword, ok := users[req.Username]
	if !ok {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token := generateJWTToken(req.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func generateJWTToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Second * tokenExpiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic("Failed to generate JWT token")
	}

	return tokenString
}

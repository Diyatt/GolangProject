// handlers/signup.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]string)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	type SignUpRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	users[req.Username] = string(hashedPassword)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}
